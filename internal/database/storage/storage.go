package storage

import "context"

type Storage interface {
	Set(context.Context, string, string)
	Get(context.Context, string) string
	Del(context.Context, string)
}

type storage struct {
	engine Engine
}

func NewStorage(engine Engine) Storage {
	return &storage{
		engine: engine,
	}
}

func (s *storage) Set(ctx context.Context, key, value string) {
	s.engine.Set(ctx, key, value)
}

func (s *storage) Get(ctx context.Context, key string) string {
	return s.engine.Get(ctx, key)
}

func (s *storage) Del(ctx context.Context, key string) {
	s.engine.Del(ctx, key)
}
