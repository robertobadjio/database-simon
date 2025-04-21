package memory

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"database-simon/internal/common"
	"database-simon/internal/database/storage"
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
func (m *Memory) Set(ctx context.Context, key, value string) {
	m.hashTable.Set(key, value)
	txID := common.GetTxIDFromContext(ctx)
	m.logger.Debug("successful get query", zap.Int64("tx", txID))
}

// Get ...
func (m *Memory) Get(ctx context.Context, key string) (string, bool) {
	txID := common.GetTxIDFromContext(ctx)
	m.logger.Debug("successful get query", zap.Int64("tx", txID))

	return m.hashTable.Get(key)
}

// Del ...
func (m *Memory) Del(ctx context.Context, key string) {
	m.hashTable.Del(key)

	txID := common.GetTxIDFromContext(ctx)
	m.logger.Debug("successful del query", zap.Int64("tx", txID))
}
