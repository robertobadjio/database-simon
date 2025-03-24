package memory

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"concurrency/internal/database/storage"
)

// Memory ...
type Memory struct {
	hashTable *HashTable
	logger    *zap.Logger
}

// NewMemory ...
func NewMemory(logger *zap.Logger) (storage.Engine, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger must be set")
	}

	return &Memory{
		hashTable: NewHashTable(),
		logger:    logger,
	}, nil
}

// Set ...
func (m *Memory) Set(_ context.Context, key, value string) {
	m.hashTable.Set(key, value)
}

// Get ...
func (m *Memory) Get(_ context.Context, key string) (string, bool) {
	return m.hashTable.Get(key)
}

// Del ...
func (m *Memory) Del(_ context.Context, key string) {
	m.hashTable.Del(key)
}
