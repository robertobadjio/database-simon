package app

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/sync/errgroup"

	"database-simon/internal/database/storage/replication"
)

// App ...
type App struct {
	serviceProvider *serviceProvider
	configFileName  string
}

// NewApp ...
func NewApp(ctx context.Context, configFileName string) (*App, error) {
	a := &App{
		configFileName: configFileName,
	}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initServiceProvider,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			log.Fatal("init", "deps", "error", err.Error())
			return err
		}
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	var err error
	a.serviceProvider, err = newServiceProvider(a.configFileName)
	if err != nil {
		return fmt.Errorf("init-serviceProvider: %w", err)
	}

	return nil
}

// Run ...
func (a *App) Run(ctx context.Context) error {
	fmt.Println("Start server...")
	group, groupCtx := errgroup.WithContext(ctx)

	db := a.serviceProvider.Database(ctx)

	var err error

	if a.serviceProvider.Config(ctx).WALS() != nil {
		if a.serviceProvider.slave != (*replication.Slave)(nil) { // TOOD: ?!
			a.serviceProvider.Logger(ctx).Info("start slave replication")
			group.Go(func() error {
				a.serviceProvider.slave.Start(groupCtx)
				return nil
			})
		} else {
			a.serviceProvider.Logger(ctx).Info("start WAL")
			group.Go(func() error {
				a.serviceProvider.WAL(ctx).Start(groupCtx)
				return nil
			})
		}

		if a.serviceProvider.master != nil {
			a.serviceProvider.Logger(ctx).Info("start master replication")
			group.Go(func() error {
				a.serviceProvider.master.Start(groupCtx)
				return nil
			})
		}
	}

	group.Go(func() error {
		a.serviceProvider.Network(ctx).HandleQueries(groupCtx, func(ctx context.Context, query []byte) []byte {
			response, _ := db.HandleQuery(ctx, string(query)) // TODO: Handle error?
			return []byte(response)
		})

		return nil
	})

	err = group.Wait()
	_ = a.serviceProvider.Logger(ctx).Sync() // TODO: Handle error

	return err
}
