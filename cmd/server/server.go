package main

import (
	"context"
	"log"

	"concurrency/internal/app"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatal("app", "init", "msg", "failed to init app", "err", err.Error())
	}

	err = a.Run(ctx)
	if err != nil {
		log.Fatal("app", "run", "msg", "failed to run app", "err", err.Error())
	}
}
