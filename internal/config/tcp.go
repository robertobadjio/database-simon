package config

import (
	"net"
	"time"
)

// TCP ...
type TCP struct {
	Host           string        `yaml:"host"`
	Port           string        `yaml:"port"`
	MaxConnections int           `yaml:"max_connections"`
	MaxMessageSize string        `yaml:"max_message_size"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
}

// Address ...
func (tcp TCP) Address() string {
	return net.JoinHostPort(tcp.Host, tcp.Port)
}
