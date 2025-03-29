package app

import (
	"concurrency/internal/network/server"
	"context"
	"fmt"
	"log"

	"go.uber.org/zap"

	"concurrency/internal/config"
	"concurrency/internal/database"
	"concurrency/internal/database/compute"
	"concurrency/internal/database/storage"
	"concurrency/internal/database/storage/engine/memory"
)

type serviceProvider struct {
	logger *zap.Logger

	database database.Database

	config  config.Config
	network *server.TCPServer
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (sp *serviceProvider) Database(ctx context.Context) database.Database {
	if sp.database == nil {
		comp := compute.NewCompute()
		memoryEngine, err := memory.NewMemory(sp.Logger(ctx))
		if err != nil {
			log.Fatal("init memory engine error")
		}

		stor, err := storage.NewStorage(memoryEngine, sp.Logger(ctx))
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
		configFileName := env.GetEnv(config.FileNameEnvName)
		err := c.Load(configFileName, env)
		if err != nil {
			log.Fatal("load config error")
		}

		sp.config = c
	}

	return sp.config
}

func (sp *serviceProvider) Network(ctx context.Context) *server.TCPServer {
	if sp.network == nil {
		var options []server.TCPServerOption
		var err error
		fmt.Println(sp.Config(ctx).TCPAddress())
		sp.network, err = server.NewTCPServer(sp.Config(ctx).TCPAddress(), sp.Logger(ctx), options...)
		if err != nil {
			log.Fatal("init network error")
		}
	}

	return sp.network
}
