package main

import (
	"bufio"
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"os"

	"concurrency/internal/database"
	"concurrency/internal/database/compute"
	"concurrency/internal/database/storage"
	"concurrency/internal/database/storage/engine"
)

func main() {
	ctx := context.Background()
	comp := compute.NewCompute()
	stor := storage.NewStorage(engine.NewMemory(1000))

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("init zap logger error")
	}
	defer func() {
		_ = logger.Sync()
	}()

	db, err := database.NewDatabase(logger, comp, stor)
	if err != nil {
		log.Fatal("init db error")
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("simon > ")
		queryRaw, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error")
			continue
		}

		res, err := db.HandleQuery(ctx, queryRaw)
		if err != nil {
			fmt.Printf("Error handle command: %s\n" + err.Error())
			continue
		}

		if res != "" {
			fmt.Println(res)
		}
	}
}
