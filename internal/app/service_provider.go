package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"

	"database-simon/internal/common"
	"database-simon/internal/config"
	"database-simon/internal/database"
	"database-simon/internal/database/compute"
	"database-simon/internal/database/filesystem"
	"database-simon/internal/database/storage"
	"database-simon/internal/database/storage/engine/memory"
	"database-simon/internal/database/storage/wal"
	"database-simon/internal/network/server"
)

type serviceProvider struct {
	logger *zap.Logger

	wal      *wal.WAL
	database *database.Database

	configFileName string
	config         config.Config

	network *server.TCPServer
}

func newServiceProvider(configFileName string) (*serviceProvider, error) {
	if configFileName == "" {
		return nil, errors.New("config file name is required")
	}

	return &serviceProvider{configFileName: configFileName}, nil
}

func (sp *serviceProvider) Database(ctx context.Context) *database.Database {
	if sp.database == nil {
		comp := compute.NewCompute()
		memoryEngine, err := memory.NewMemory(sp.Logger(ctx))
		if err != nil {
			log.Fatal("init memory engine error")
		}

		stor, err := storage.NewStorage(memoryEngine, sp.Logger(ctx), storage.WithWAL(sp.wal))
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

func (sp *serviceProvider) Config(_ context.Context) config.Config {
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

		if sp.Config(ctx).TCPConfigS().MaxConnections != 0 {
			options = append(options, server.WithServerMaxConnectionsNumber(uint(sp.Config(ctx).TCPConfigS().MaxConnections))) // nolint : G115: integer overflow conversion int -> uint (gosec)
		}

		if sp.Config(ctx).TCPConfigS().MaxMessageSize != "" {
			size, errParseSize := common.ParseSize(sp.Config(ctx).TCPConfigS().MaxMessageSize)
			if errParseSize != nil {
				log.Fatal("incorrect max message size")
			}

			options = append(options, server.WithServerBufferSize(uint(size))) // nolint : G115: integer overflow conversion int -> uint (gosec)
		}

		if sp.Config(ctx).TCPConfigS().IdleTimeout != 0 {
			options = append(options, server.WithServerIdleTimeout(sp.Config(ctx).TCPConfigS().IdleTimeout))
		}

		// TODO: Адрес по умолчанию
		fmt.Println(sp.Config(ctx).TCPAddress()) // TODO: Удалить

		sp.network, err = server.NewTCPServer(sp.Config(ctx).TCPAddress(), sp.Logger(ctx), options...)
		if err != nil {
			log.Fatalf("init network error: %v", err)
		}
	}

	return sp.network
}

const (
	defaultFlushingBatchSize    = 100
	defaultFlushingBatchTimeout = time.Millisecond * 10
	defaultMaxSegmentSize       = 10 << 20
	defaultWALDataDirectory     = "./data/wal"
)

func (sp *serviceProvider) WAL(ctx context.Context) *wal.WAL {
	if sp.wal != nil {
		return sp.wal
	}

	if sp.Config(ctx).WALS() == nil {
		return nil
	}

	flushingBatchSize := defaultFlushingBatchSize
	flushingBatchTimeout := defaultFlushingBatchTimeout
	maxSegmentSize := defaultMaxSegmentSize
	dataDirectory := defaultWALDataDirectory

	if sp.Config(ctx).WALS().FlushingBatchSize != 0 {
		flushingBatchSize = sp.Config(ctx).WALS().FlushingBatchSize
	}

	if sp.Config(ctx).WALS().FlushingBatchTimeout != 0 {
		flushingBatchTimeout = sp.Config(ctx).WALS().FlushingBatchTimeout
	}

	if sp.Config(ctx).WALS().MaxSegmentSize != "" {
		size, err := common.ParseSize(sp.Config(ctx).WALS().MaxSegmentSize)
		if err != nil {
			log.Fatal(errors.New("max segment size is incorrect"))
		}

		maxSegmentSize = size
	}

	if sp.Config(ctx).WALS().DataDirectory != "" {
		dataDirectory = sp.Config(ctx).WALS().DataDirectory
	}

	segmentsDirectory := filesystem.NewSegmentsDirectory(dataDirectory)
	reader, err := wal.NewLogsReader(segmentsDirectory)
	if err != nil {
		log.Fatal(err)
	}

	segment := filesystem.NewSegment(dataDirectory, maxSegmentSize)
	writer, err := wal.NewLogsWriter(segment, sp.Logger(ctx))
	if err != nil {
		log.Fatal(err)
	}

	w, err := wal.NewWAL(writer, reader, flushingBatchTimeout, flushingBatchSize)
	if err != nil {
		log.Fatal(err)
	}

	sp.wal = w

	return sp.wal
}
