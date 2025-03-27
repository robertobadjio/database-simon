package server

import (
	"time"
)

const defaultBufferSize = 4 << 10

type TCPServerOption func(*TCPServer)

func WithServerIdleTimeout(timeout time.Duration) TCPServerOption {
	return func(server *TCPServer) {
		server.idleTimeout = timeout
	}
}

func WithServerBufferSize(size uint) TCPServerOption {
	return func(server *TCPServer) {
		server.bufferSize = int(size)
	}
}

func WithServerMaxConnectionsNumber(count uint) TCPServerOption {
	return func(server *TCPServer) {
		server.maxConnections = int(count)
	}
}
