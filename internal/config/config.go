package config

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"time"

	"gopkg.in/yaml.v3"
)

// FileNameEnvName ...
const FileNameEnvName = "CONFIG_FILE_NAME"

// Config ...
type Config interface {
	TCPAddress() string
	Load(string, OS) error
}

type config struct {
	TCPConfig *TCPConfig `yaml:"network"`
}

// TCPConfig ...
type TCPConfig struct {
	Host           string        `yaml:"host"`
	Port           string        `yaml:"port"`
	MaxConnections int           `yaml:"max_connections"`
	MaxMessageSize string        `yaml:"max_message_size"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
}

// TCPAddress ...
func (cfg *config) TCPAddress() string {
	return net.JoinHostPort(cfg.TCPConfig.Host, cfg.TCPConfig.Port)
}

// NewConfig ...
func NewConfig() Config {
	return &config{}
}

// Load ...
func (cfg *config) Load(configFileName string, osCustom OS) error {
	dataRaw, err := osCustom.ReadFile(configFileName)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	reader := bytes.NewReader(dataRaw)
	if reader == nil {
		return fmt.Errorf("incorrect reader")
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read buffer: %w", err)
	}

	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	return nil
}
