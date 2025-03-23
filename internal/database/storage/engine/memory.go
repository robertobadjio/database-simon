package engine

import (
	"context"
	"sync"

	"concurrency/internal/database/storage"
)

type Memory struct {
	mu    sync.RWMutex
	items map[string]string
}

func NewMemory(amount int) storage.Engine {
	return &Memory{
		items: make(map[string]string, amount),
	}
}

func (s *Memory) Set(ctx context.Context, key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[key] = value
}

func (s *Memory) Get(ctx context.Context, key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if v, found := s.items[key]; found {
		return v
	}

	return ""
}

func (s *Memory) Del(ctx context.Context, key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.items, key)
}
