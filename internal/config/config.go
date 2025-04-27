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
// TODO: Remove interface
type Config interface {
	TCPAddress() string
	WALS() *WAL
	Load(string, OS) error
	TCPConfigS() *TCPConfig
	ReplicationS() *Replication
}

type config struct {
	TCPConfig   *TCPConfig   `yaml:"network"`
	WAL         *WAL         `yaml:"wal"`
	Replication *Replication `yaml:"replication"`
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

// Replication ...
type Replication struct {
	ReplicaType       string        `yaml:"replica_type"`
	MasterAddress     string        `yaml:"master_address"`
	SyncInterval      time.Duration `yaml:"sync_interval"`
	MaxReplicasNumber int           `yaml:"max_replicas_number"`
}

// WALS ...
func (cfg *config) WALS() *WAL {
	return cfg.WAL
}

// TCPConfigS ...
func (cfg *config) TCPConfigS() *TCPConfig {
	return cfg.TCPConfig
}

// ReplicationS ...
func (cfg *config) ReplicationS() *Replication {
	return cfg.Replication
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
