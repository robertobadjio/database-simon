package memory

import (
	"context"
	"fmt"
	"hash/fnv"

	"go.uber.org/zap"

	"database-simon/internal/common"
)

// Memory ...
type Memory struct {
	partitions []*HashTable
	logger     *zap.Logger
}

// NewMemory ...
func NewMemory(logger *zap.Logger, options ...EngineOption) (*Memory, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger must be set")
	}

	memoryEngine := &Memory{
		logger: logger,
	}

	for _, option := range options {
		option(memoryEngine)
	}

	if len(memoryEngine.partitions) == 0 {
		memoryEngine.partitions = make([]*HashTable, 1)
		memoryEngine.partitions[0] = NewHashTable()
	}

	return memoryEngine, nil
}

// Set ...
func (m *Memory) Set(ctx context.Context, key, value string) {
	partitionIdx := 0
	if len(m.partitions) > 1 {
		partitionIdx = m.partitionIdx(key)
	}

	partition := m.partitions[partitionIdx]
	partition.Set(key, value)

	//m.hashTable.Set(key, value)
	txID := common.GetTxIDFromContext(ctx)
	m.logger.Debug("successful get query", zap.Int64("tx", txID))
}

// Get ...
func (m *Memory) Get(ctx context.Context, key string) (string, bool) {
	partitionIdx := 0
	if len(m.partitions) > 1 {
		partitionIdx = m.partitionIdx(key)
	}

	partition := m.partitions[partitionIdx]
	value, found := partition.Get(key)

	txID := common.GetTxIDFromContext(ctx)
	m.logger.Debug("successful get query", zap.Int64("tx", txID))

	return value, found
}

// Del ...
func (m *Memory) Del(ctx context.Context, key string) {
	partitionIdx := 0
	if len(m.partitions) > 1 {
		partitionIdx = m.partitionIdx(key)
	}

	partition := m.partitions[partitionIdx]
	partition.Del(key)

	txID := common.GetTxIDFromContext(ctx)
	m.logger.Debug("successful del query", zap.Int64("tx", txID))
}

func (m *Memory) partitionIdx(key string) int {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(key))
	return int(hash.Sum32()) % len(m.partitions)
}
