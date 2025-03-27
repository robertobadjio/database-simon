package app

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
)

// App ...
type App struct {
	serviceProvider *serviceProvider
}

// NewApp ...
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

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
	a.serviceProvider = newServiceProvider()
	return nil
}

// Run ...
func (a *App) Run(ctx context.Context) error {
	fmt.Println("Start server...")
	group, groupCtx := errgroup.WithContext(ctx)

	var err error

	group.Go(func() error {
		a.serviceProvider.Network(ctx).HandleQueries(groupCtx, func(ctx context.Context, query []byte) []byte {
			response, err := a.serviceProvider.Database(ctx).HandleQuery(ctx, string(query))
			if err != nil {
				fmt.Println("ERROR", err.Error()) // TODO: ?!
			}

			if response == "" {
				response = "Success" // TODO: !
			}

			return []byte(response)
		})

		return nil
	})

	err = group.Wait()
	_ = a.serviceProvider.Logger(ctx).Sync() // TODO: Handle error

	return err
}
