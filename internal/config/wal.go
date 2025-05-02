package config

import (
	"errors"
	"log"
	"time"

	"database-simon/internal/common"
)

const (
	defaultFlushingBatchSize    = 100
	defaultFlushingBatchTimeout = time.Millisecond * 10
	defaultMaxSegmentSize       = 10 << 20
	defaultWALDataDirectory     = "./data/wal"
)

// WAL ...
type WAL struct {
	FlushingBatchSize    int           `yaml:"flushing_batch_size"`
	FlushingBatchTimeout time.Duration `yaml:"flushing_batch_timeout"`
	MaxSegmentSize       string        `yaml:"max_segment_size"`
	DataDirectory        string        `yaml:"data_directory"`
}

// GetFlushingBatchSize ...
func (w WAL) GetFlushingBatchSize() int {
	flushingBatchSize := defaultFlushingBatchSize
	if w.FlushingBatchSize != 0 {
		flushingBatchSize = w.FlushingBatchSize
	}

	return flushingBatchSize
}

// GetFlushingBatchTimeout ...
func (w WAL) GetFlushingBatchTimeout() time.Duration {
	flushingBatchTimeout := defaultFlushingBatchTimeout
	if w.FlushingBatchTimeout != 0 {
		flushingBatchTimeout = w.FlushingBatchTimeout
	}

	return flushingBatchTimeout
}

// GetMaxSegmentSize ...
func (w WAL) GetMaxSegmentSize() int {
	maxSegmentSize := defaultMaxSegmentSize
	if w.MaxSegmentSize != "" {
		size, err := common.ParseSize(w.MaxSegmentSize)
		if err != nil {
			log.Fatal(errors.New("max segment size is incorrect"))
		}

		maxSegmentSize = size
	}

	return maxSegmentSize
}

// GetDataDirectory ...
func (w WAL) GetDataDirectory() string {
	dataDirectory := defaultWALDataDirectory
	if w.DataDirectory != "" {
		dataDirectory = w.DataDirectory
	}

	return dataDirectory
}
