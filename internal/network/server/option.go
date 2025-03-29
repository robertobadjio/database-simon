package server

import (
	"time"
)

// TCPServerOption ...
type TCPServerOption func(*TCPServer)

// WithServerIdleTimeout ...
func WithServerIdleTimeout(timeout time.Duration) TCPServerOption {
	return func(server *TCPServer) {
		server.idleTimeout = timeout
	}
}

// WithServerBufferSize ...
func WithServerBufferSize(size uint) TCPServerOption {
	return func(server *TCPServer) {
		server.bufferSize = int(size) // nolint : G115: integer overflow conversion uint -> int
	}
}

// WithServerMaxConnectionsNumber ...
func WithServerMaxConnectionsNumber(count uint) TCPServerOption {
	return func(server *TCPServer) {
		server.maxConnections = int(count) // nolint : G115: integer overflow conversion uint -> int
	}
}
