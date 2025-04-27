package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const testConfigData = `
network:
  host: "127.0.0.1"
  port: "8081"
  max_connections: 100
  max_message_size: "4KB"
  idle_timeout: 5m
wal:
  flushing_batch_size: 100
  flushing_batch_timeout: "10ms"
  max_segment_size: "10MB"
  data_directory: "./data/wal"
replication:
  replica_type: "slave"
  master_address: "127.0.0.1:3232"
  sync_interval: "1s"
  max_replicas_number: 1
`

func TestNewConfig(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		expectedNilObj bool
	}{
		"load config": {
			expectedNilObj: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg := NewConfig()
			if test.expectedNilObj {
				assert.Nil(t, cfg)
			} else {
				assert.NotNil(t, cfg)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)

	tests := map[string]struct {
		os          func() OS
		cfgFileName string

		expectedErr error
		expectedCfg Config
	}{
		"load config": {
			cfgFileName: "config.yml",
			os: func() OS {
				env := NewMockOS(controller)
				env.EXPECT().ReadFile("config.yml").Return([]byte(testConfigData), nil).Times(1)

				return env
			},
			expectedErr: nil,
			expectedCfg: &config{
				&TCPConfig{
					Host:           "127.0.0.1",
					Port:           "8081",
					MaxConnections: 100,
					MaxMessageSize: "4KB",
					IdleTimeout:    5 * time.Minute,
				},
				&WAL{
					FlushingBatchSize:    100,
					FlushingBatchTimeout: 10 * time.Millisecond,
					MaxSegmentSize:       "10MB",
					DataDirectory:        "./data/wal",
				},
				&Replication{
					ReplicaType:       "slave",
					MasterAddress:     "127.0.0.1:3232",
					SyncInterval:      time.Second,
					MaxReplicasNumber: 1,
				},
			},
		},
		"load empty config": {
			cfgFileName: "",
			os: func() OS {
				env := NewMockOS(controller)
				env.EXPECT().ReadFile(gomock.Any()).Return([]byte(""), nil).Times(1)

				return env
			},
			expectedErr: nil,
			expectedCfg: &config{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg := NewConfig()
			assert.NotNil(t, cfg)

			err := cfg.Load(test.cfgFileName, test.os())
			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedCfg, cfg)
		})
	}
}
