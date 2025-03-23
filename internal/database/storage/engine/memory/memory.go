package memory

import (
	"context"

	"concurrency/internal/database/storage"
)

type Memory struct {
	hashTable *HashTable
}

func NewMemory() storage.Engine {
	return &Memory{
		hashTable: NewHashTable(),
	}
}

func (m *Memory) Set(ctx context.Context, key, value string) {
	m.hashTable.Set(key, value)
}

func (m *Memory) Get(ctx context.Context, key string) (string, bool) {
	return m.hashTable.Get(key)
}

func (m *Memory) Del(ctx context.Context, key string) {
	m.hashTable.Del(key)
}
