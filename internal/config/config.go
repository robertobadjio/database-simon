package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const configFileNameEnvName = "CONFIG_FILE_NAME"

type Config interface {
	TCPAddress() string
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

func NewConfig() (Config, error) {
	configFileName := os.Getenv(configFileNameEnvName)
	if len(configFileName) != 0 {
		data, err := os.ReadFile(configFileName)
		if err != nil {
			log.Fatal(err)
		}

		reader := bytes.NewReader(data)
		cfg, err := load(reader)
		if err != nil {
			log.Fatal(err)
		}

		return cfg, nil
	}

	return &config{}, nil
}

func load(reader io.Reader) (Config, error) {
	if reader == nil {
		return nil, errors.New("incorrect reader")
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.New("failed to read buffer")
	}

	var c config
	if err = yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &c, nil
}
