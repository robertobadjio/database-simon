package client

import (
	"time"
)

const defaultBufferSize = 4 << 10

// TCPClientOption ...
type TCPClientOption func(*TCPClient)

// WithClientIdleTimeout ...
func WithClientIdleTimeout(timeout time.Duration) TCPClientOption {
	return func(client *TCPClient) {
		client.idleTimeout = timeout
	}
}

// WithClientBufferSize ...
func WithClientBufferSize(size uint) TCPClientOption {
	return func(client *TCPClient) {
		client.bufferSize = int(size) // nolint : G115: integer overflow conversion uint -> int
	}
}
