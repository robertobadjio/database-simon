package main

import (
	"context"
	"flag"
	"log"

	"database-simon/internal/app"
)

func main() {
	ctx := context.Background()

	config := flag.String("config", "./config.yml", "Server config")
	flag.Parse()

	a, err := app.NewApp(ctx, *config)
	if err != nil {
		log.Fatal("app", "init", "msg", "failed to init app", "err", err.Error())
	}

	err = a.Run(ctx)
	if err != nil {
		log.Fatal("app", "run", "msg", "failed to run app", "err", err.Error())
	}
}
