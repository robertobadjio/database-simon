package config

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"time"

	"gopkg.in/yaml.v3"
)

// Config ...
type Config interface {
	TCPAddress() string
	WALS() *WAL
	Load(string, OS) error
	TCPConfigS() *TCPConfig
}

type config struct {
	TCPConfig *TCPConfig `yaml:"network"`
	WAL       *WAL       `yaml:"wal"`
}

// TCPConfig ...
type TCPConfig struct {
	Host           string        `yaml:"host"`
	Port           string        `yaml:"port"`
	MaxConnections int           `yaml:"max_connections"`
	MaxMessageSize string        `yaml:"max_message_size"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
}

// WAL ...
type WAL struct {
	FlushingBatchSize    int           `yaml:"flushing_batch_size"`
	FlushingBatchTimeout time.Duration `yaml:"flushing_batch_timeout"`
	MaxSegmentSize       string        `yaml:"max_segment_size"`
	DataDirectory        string        `yaml:"data_directory"`
}

// WALS ...
func (cfg *config) WALS() *WAL {
	return cfg.WAL
}

// TCPConfigS ...
func (cfg *config) TCPConfigS() *TCPConfig {
	return cfg.TCPConfig
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
