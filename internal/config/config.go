package config

import (
	"bytes"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

// Config ...
type Config struct {
	Engine      *Engine      `yaml:"engine"`
	TCP         *TCP         `yaml:"network"`
	WAL         *WAL         `yaml:"wal"`
	Replication *Replication `yaml:"replication"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{}
}

// Load ...
func (cfg *Config) Load(configFileName string, osCustom OS) error {
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
