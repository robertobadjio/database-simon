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

	"concurrency/internal/network/client"
)

func main() {
	address := flag.String("address", "localhost:8081", "Address of the spider")
	idleTimeout := flag.Duration("idle_timeout", time.Minute, "Idle timeout for connection")
	maxMessageSizeStr := flag.String("max_message_size", "4KB", "Max message size for connection")
	flag.Parse()

	logger, _ := zap.NewProduction()
	maxMessageSize, err := ParseSize(*maxMessageSizeStr)
	if err != nil {
		logger.Fatal("failed to parse max message size", zap.Error(err))
	}

	var options []client.TCPClientOption
	options = append(options, client.WithClientIdleTimeout(*idleTimeout))
	options = append(options, client.WithClientBufferSize(uint(maxMessageSize)))

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

func ParseSize(text string) (int, error) {
	if len(text) == 0 || text[0] < '0' || text[0] > '9' {
		return 0, errors.New("incorrect size")
	}

	idx := 0
	size := 0
	for idx < len(text) && text[idx] >= '0' && text[idx] <= '9' {
		number := int(text[idx] - '0')
		size = size*10 + number
		idx++
	}

	parameter := text[idx:]
	switch parameter {
	case "GB", "Gb", "gb":
		return size << 30, nil
	case "MB", "Mb", "mb":
		return size << 20, nil
	case "KB", "Kb", "kb":
		return size << 10, nil
	case "B", "b", "":
		return size, nil
	default:
		return 0, errors.New("incorrect size")
	}
}
