package config

import (
	"errors"
	"time"
)

const (
	defaultReplicationSyncInterval = time.Second
	defaultMaxReplicasNumber       = 5
)

const (
	// MasterType ...
	MasterType = "master"
	// SlaveType ...
	SlaveType = "slave"
)

// SupportedTypes ...
var SupportedTypes = map[string]struct{}{
	MasterType: {},
	SlaveType:  {},
}

// Replication ...
type Replication struct {
	ReplicaType       string        `yaml:"replica_type"`
	MasterAddress     string        `yaml:"master_address"`
	SyncInterval      time.Duration `yaml:"sync_interval"`
	MaxReplicasNumber int           `yaml:"max_replicas_number"`
}

// GetSyncInterval ...
func (r Replication) GetSyncInterval() time.Duration {
	syncInterval := defaultReplicationSyncInterval
	if r.SyncInterval != 0 {
		syncInterval = r.SyncInterval
	}

	return syncInterval
}

// GetMaxReplicasNumber ...
func (r Replication) GetMaxReplicasNumber() int {
	maxReplicasNumber := defaultMaxReplicasNumber
	if r.MaxReplicasNumber != 0 {
		maxReplicasNumber = r.MaxReplicasNumber
	}

	return maxReplicasNumber
}

// GetMasterAddress ...
func (r Replication) GetMasterAddress() (string, error) {
	if r.MasterAddress == "" {
		return "", errors.New("master address is incorrect")
	}

	return r.MasterAddress, nil
}
