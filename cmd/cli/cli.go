package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"syscall"
	"time"

	"go.uber.org/zap"

	"database-simon/internal/common"
	"database-simon/internal/network/client"
)

func main() {
	//address := flag.String("address", "localhost:3232", "Address of the spider")
	address := flag.String("address", "localhost:8081", "Address of the spider")
	idleTimeout := flag.Duration("idle_timeout", time.Minute, "Idle timeout for connection")
	maxMessageSizeStr := flag.String("max_message_size", "4KB", "Max message size for connection")
	flag.Parse()

	logger, _ := zap.NewProduction()
	maxMessageSize, err := common.ParseSize(*maxMessageSizeStr)
	if err != nil {
		logger.Fatal("failed to parse max message size", zap.Error(err))
	}

	var options []client.TCPClientOption
	options = append(options, client.WithClientIdleTimeout(*idleTimeout))
	options = append(options, client.WithClientBufferSize(uint(maxMessageSize))) // nolint : G115: integer overflow conversion uint -> int

	reader := bufio.NewReader(os.Stdin)
	c, err := client.NewTCPClient(*address, options...)
	if err != nil {
		logger.Fatal("failed to connect with server", zap.Error(err))
	}

	for {
		fmt.Print("simon > ")
		request, err := reader.ReadString('\n')
		if errors.Is(err, syscall.EPIPE) {
			logger.Fatal("connection was closed", zap.Error(err))
		} else if err != nil {
			logger.Error("failed to read query", zap.Error(err))
		}

		response, err := c.Send([]byte(request))
		if errors.Is(err, syscall.EPIPE) {
			logger.Fatal("connection was closed", zap.Error(err))
		} else if err != nil {
			logger.Error("failed to send query", zap.Error(err))
		}

		fmt.Println(string(response))
	}
}
