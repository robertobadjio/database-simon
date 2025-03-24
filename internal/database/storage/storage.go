package storage

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// Storage ...
type Storage interface {
	Set(context.Context, string, string) error
	Get(context.Context, string) (string, error)
	Del(context.Context, string) error
}

type storage struct {
	engine Engine
	logger *zap.Logger
}

// NewStorage ...
func NewStorage(engine Engine, logger *zap.Logger) (Storage, error) {
	if engine == nil {
		return nil, fmt.Errorf("engine must be set")
	}

	if logger == nil {
		return nil, fmt.Errorf("logger must be set")
	}

	return &storage{
		engine: engine,
		logger: logger,
	}, nil
}

// Set ...
func (s *storage) Set(ctx context.Context, key, value string) error {
	s.engine.Set(ctx, key, value)
	return nil
}

// Get ...
func (s *storage) Get(ctx context.Context, key string) (string, error) {
	val, found := s.engine.Get(ctx, key)
	if !found {
		return "", fmt.Errorf("not found")
	}

	return val, nil
}

// Del ...
func (s *storage) Del(ctx context.Context, key string) error {
	s.engine.Del(ctx, key)
	return nil
}
