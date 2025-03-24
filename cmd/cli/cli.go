package main

import (
	"bufio"
	"concurrency/internal/database/storage/engine/memory"
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"os"

	"concurrency/internal/database"
	"concurrency/internal/database/compute"
	"concurrency/internal/database/storage"
)

func main() {
	ctx := context.Background()
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("init zap logger error")
	}
	defer func() {
		_ = logger.Sync()
	}()

	comp := compute.NewCompute()
	stor, err := storage.NewStorage(memory.NewMemory(), logger)
	if err != nil {
		log.Fatal("init storage error")
	}

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
