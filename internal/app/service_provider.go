package app

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.uber.org/zap"

	"database-simon/internal/common"
	"database-simon/internal/config"
	"database-simon/internal/database"
	"database-simon/internal/database/compute"
	"database-simon/internal/database/filesystem"
	"database-simon/internal/database/storage"
	"database-simon/internal/database/storage/engine/memory"
	"database-simon/internal/database/storage/replication"
	"database-simon/internal/database/storage/wal"
	"database-simon/internal/network/client"
	"database-simon/internal/network/server"
)

type serviceProvider struct {
	logger *zap.Logger

	wal      *wal.WAL
	slave    *replication.Slave
	master   *replication.Master
	database *database.Database

	configFileName string
	config         *config.Config

	network *server.TCPServer
}

func newServiceProvider(configFileName string) (*serviceProvider, error) {
	if configFileName == "" {
		return nil, errors.New("config file name is required")
	}

	return &serviceProvider{configFileName: configFileName}, nil
}

// Database ...
func (sp *serviceProvider) Database(ctx context.Context) *database.Database {
	if sp.database == nil {
		comp := compute.NewCompute()

		var memoryOptions []memory.EngineOption
		if sp.Config(ctx).Engine.PartitionsNumber != 0 {
			memoryOptions = append(memoryOptions, memory.WithPartitions(sp.Config(ctx).Engine.PartitionsNumber))
		}

		memoryEngine, err := memory.NewMemory(sp.Logger(ctx), memoryOptions...)
		if err != nil {
			log.Fatal("init memory engine error")
		}

		replica, err := sp.Replica(ctx)
		if err != nil {
			log.Fatalf("init replica error: %v", err)
		}

		switch v := replica.(type) {
		case *replication.Slave:
			sp.slave = v
		case *replication.Master:
			sp.master = v
		}

		var storageOptions []storage.Option
		if sp.WAL(ctx) != nil {
			storageOptions = append(storageOptions, storage.WithWAL(sp.WAL(ctx)))
		}

		if sp.master != nil {
			storageOptions = append(storageOptions, storage.WithReplication(sp.master))
		} else if sp.slave != nil {
			storageOptions = append(storageOptions, storage.WithReplication(sp.slave))
			storageOptions = append(storageOptions, storage.WithReplicationStream(sp.slave.ReplicationStream()))
		}

		stor, err := storage.NewStorage(memoryEngine, sp.Logger(ctx), storageOptions...)
		if err != nil {
			log.Fatal("init storage error")
		}

		db, err := database.NewDatabase(sp.Logger(ctx), comp, stor)
		if err != nil {
			log.Fatal("init db error")
		}
		sp.database = db
	}

	return sp.database
}

// Logger ...
func (sp *serviceProvider) Logger(_ context.Context) *zap.Logger {
	if sp.logger == nil {
		logger, err := zap.NewProduction()
		if err != nil {
			log.Fatal("init zap logger error")
		}
		sp.logger = logger
	}

	return sp.logger
}

// Config ...
func (sp *serviceProvider) Config(_ context.Context) *config.Config {
	if sp.config == nil {
		c := config.NewConfig()

		env := config.NewEnvironment()

		err := c.Load(sp.configFileName, env)
		if err != nil {
			log.Fatal(fmt.Sprintf("load config error: %s", err.Error()))
		}

		sp.config = c
	}

	return sp.config
}

// Network ...
func (sp *serviceProvider) Network(ctx context.Context) *server.TCPServer {
	if sp.network == nil {
		var options []server.TCPServerOption
		var err error

		if sp.Config(ctx).TCP.MaxConnections != 0 {
			options = append(options, server.WithServerMaxConnectionsNumber(uint(sp.Config(ctx).TCP.MaxConnections))) // nolint : G115: integer overflow conversion int -> uint (gosec)
		}

		if sp.Config(ctx).TCP.MaxMessageSize != "" {
			size, errParseSize := common.ParseSize(sp.Config(ctx).TCP.MaxMessageSize)
			if errParseSize != nil {
				log.Fatal("incorrect max message size")
			}

			options = append(options, server.WithServerBufferSize(uint(size))) // nolint : G115: integer overflow conversion int -> uint (gosec)
		}

		if sp.Config(ctx).TCP.IdleTimeout != 0 {
			options = append(options, server.WithServerIdleTimeout(sp.Config(ctx).TCP.IdleTimeout))
		}

		// TODO: Адрес по умолчанию
		fmt.Println(sp.Config(ctx).TCP.Address()) // TODO: Удалить

		sp.network, err = server.NewTCPServer(sp.Config(ctx).TCP.Address(), sp.Logger(ctx), options...)
		if err != nil {
			log.Fatalf("init network error: %v", err)
		}
	}

	return sp.network
}

// WAL ...
func (sp *serviceProvider) WAL(ctx context.Context) *wal.WAL {
	if sp.wal != nil {
		return sp.wal
	}

	if sp.Config(ctx).WAL == nil {
		return nil
	}

	segmentsDirectory := filesystem.NewSegmentsDirectory(sp.Config(ctx).WAL.GetDataDirectory())
	reader, err := wal.NewLogsReader(segmentsDirectory)
	if err != nil {
		log.Fatal(err)
	}

	segment := filesystem.NewSegment(sp.Config(ctx).WAL.GetDataDirectory(), sp.Config(ctx).WAL.GetMaxSegmentSize())
	writer, err := wal.NewLogsWriter(segment, sp.Logger(ctx))
	if err != nil {
		log.Fatal(err)
	}

	w, err := wal.NewWAL(
		writer,
		reader,
		sp.Config(ctx).WAL.GetFlushingBatchTimeout(),
		sp.Config(ctx).WAL.GetFlushingBatchSize(),
	)
	if err != nil {
		log.Fatal(err)
	}

	sp.wal = w

	return sp.wal
}

// Replica ...
func (sp *serviceProvider) Replica(ctx context.Context) (interface{}, error) {
	if sp.Config(ctx).Replication == nil {
		return nil, nil
	}

	if _, found := config.SupportedTypes[sp.Config(ctx).Replication.ReplicaType]; !found {
		return nil, errors.New("replica type is incorrect")
	}

	maxMessageSize := sp.Config(ctx).WAL.GetMaxSegmentSize()
	walDirectory := sp.Config(ctx).WAL.GetDataDirectory()

	if sp.Config(ctx).WAL.DataDirectory != "" {
		walDirectory = sp.Config(ctx).WAL.DataDirectory
	}

	if sp.Config(ctx).WAL.MaxSegmentSize != "" {
		size, _ := common.ParseSize(sp.Config(ctx).WAL.MaxSegmentSize)
		maxMessageSize = size
	}

	idleTimeout := sp.Config(ctx).Replication.GetSyncInterval() * 3
	if sp.config.Replication.ReplicaType == config.MasterType {
		var options []server.TCPServerOption
		options = append(options, server.WithServerIdleTimeout(idleTimeout))
		options = append(options, server.WithServerBufferSize(uint(maxMessageSize)))                                              // nolint : G115: integer overflow conversion int -> uint (gosec)
		options = append(options, server.WithServerMaxConnectionsNumber(uint(sp.Config(ctx).Replication.GetMaxReplicasNumber()))) // nolint : G115: integer overflow conversion int -> uint (gosec)
		s, err := server.NewTCPServer(sp.Config(ctx).Replication.MasterAddress, sp.Logger(ctx), options...)
		if err != nil {
			return nil, err
		}

		return replication.NewMaster(s, walDirectory, sp.Logger(ctx))
	}

	var options []client.TCPClientOption
	//options = append(options, client.WithClientIdleTimeout(idleTimeout))
	options = append(options, client.WithClientBufferSize(uint(maxMessageSize))) // nolint : G115: integer overflow conversion int -> uint (gosec)
	c, err := client.NewTCPClient(sp.Config(ctx).Replication.MasterAddress, options...)
	if err != nil {
		return nil, err
	}

	return replication.NewSlave(c, walDirectory, sp.Config(ctx).Replication.GetSyncInterval(), sp.Logger(ctx))
}
