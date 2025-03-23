package storage

import (
	"context"
	"fmt"
)

type Storage interface {
	Set(context.Context, string, string) error
	Get(context.Context, string) (string, error)
	Del(context.Context, string) error
}

type storage struct {
	engine Engine
}

func NewStorage(engine Engine) Storage {
	return &storage{
		engine: engine,
	}
}

func (s *storage) Set(ctx context.Context, key, value string) error {
	s.engine.Set(ctx, key, value)
	return nil
}

func (s *storage) Get(ctx context.Context, key string) (string, error) {
	val, found := s.engine.Get(ctx, key)
	if !found {
		return "", fmt.Errorf("not found") // TODO: custom error
	}

	return val, nil
}

func (s *storage) Del(ctx context.Context, key string) error {
	s.engine.Del(ctx, key)
	return nil
}
