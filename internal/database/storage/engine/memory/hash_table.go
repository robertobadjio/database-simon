package memory

import (
	"sync"
)

// HashTable ...
type HashTable struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewHashTable ...
func NewHashTable() *HashTable {
	return &HashTable{
		data: make(map[string]string),
	}
}

// Set ...
func (ht *HashTable) Set(key, value string) {
	ht.mu.Lock()
	defer ht.mu.Unlock()

	ht.data[key] = value
}

// Get ...
func (ht *HashTable) Get(key string) (string, bool) {
	ht.mu.RLock()
	defer ht.mu.RUnlock()

	value, found := ht.data[key]
	return value, found
}

// Del ...
func (ht *HashTable) Del(key string) {
	ht.mu.Lock()
	defer ht.mu.Unlock()

	delete(ht.data, key)
}
