package storage

import (
	"context"
	"database-simon/internal/common"
	"database-simon/internal/database/compute"
	wal2 "database-simon/internal/database/storage/wal"
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
	engine    Engine
	logger    *zap.Logger
	wal       wal2.WAL
	generator *IDGenerator
}

// NewStorage ...
func NewStorage(engine Engine, logger *zap.Logger) (Storage, error) {
	if engine == nil {
		return nil, fmt.Errorf("engine must be set")
	}

	if logger == nil {
		return nil, fmt.Errorf("logger must be set")
	}

	st := &storage{
		engine: engine,
		logger: logger,
	}

	var lastLSN int64
	if st.wal != nil {
		logs, err := st.wal.Recover()
		if err != nil {
			logger.Error("failed to recover data from WAL", zap.Error(err))
		} else {
			lastLSN = st.applyData(logs)
		}
	}

	st.generator = NewIDGenerator(lastLSN)

	return st, nil
}

// Set ...
func (s *storage) Set(ctx context.Context, key, value string) error {
	txID := s.generator.Generate()
	ctx = common.ContextWithTxID(ctx, txID)

	if s.wal != nil {
		futureResponse := s.wal.Set(ctx, key, value)
		if err := futureResponse.Get(); err != nil {
			return err
		}
	}

	s.engine.Set(ctx, key, value)

	return nil
}

// Get ...
func (s *storage) Get(ctx context.Context, key string) (string, error) {
	txID := s.generator.Generate()
	ctx = common.ContextWithTxID(ctx, txID)

	val, found := s.engine.Get(ctx, key)
	if !found {
		return "", fmt.Errorf("not found")
	}

	return val, nil
}

// Del ...
func (s *storage) Del(ctx context.Context, key string) error {
	txID := s.generator.Generate()
	ctx = common.ContextWithTxID(ctx, txID)

	if s.wal != nil {
		futureResponse := s.wal.Del(ctx, key)
		if err := futureResponse.Get(); err != nil {
			return err
		}
	}

	s.engine.Del(ctx, key)

	return nil
}

func (s *storage) applyData(logs []wal2.Log) int64 {
	var lastLSN int64
	for _, log := range logs {
		lastLSN = max(lastLSN, log.LSN)
		ctx := common.ContextWithTxID(context.Background(), log.LSN)
		switch log.CommandID {
		case compute.SetCommand:
			s.engine.Set(ctx, log.Arguments[0], log.Arguments[1])
		case compute.DelCommand:
			s.engine.Del(ctx, log.Arguments[0])
		}
	}

	return lastLSN
}
