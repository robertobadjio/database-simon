package storage

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"database-simon/internal/common"
	"database-simon/internal/concurrency"
	"database-simon/internal/database/compute"
	"database-simon/internal/database/storage/wal"
)

var (
	// ErrorNotFound ...
	ErrorNotFound = errors.New("not found")
	// ErrorMutableTX ...
	ErrorMutableTX = errors.New("mutable transaction on slave")
)

type walI interface {
	Recover() ([]wal.Log, error)
	Set(context.Context, string, string) concurrency.FutureError
	Del(context.Context, string) concurrency.FutureError
}

type engine interface {
	Set(context.Context, string, string)
	Get(context.Context, string) (string, bool)
	Del(context.Context, string)
}

type replica interface {
	IsMaster() bool
}

// Storage ...
type Storage struct {
	engine    engine
	replica   replica
	logger    *zap.Logger
	wal       walI
	stream    <-chan []wal.Log
	generator *IDGenerator
}

// NewStorage ...
func NewStorage(engine engine, logger *zap.Logger, options ...Option) (*Storage, error) {
	if engine == nil {
		return nil, fmt.Errorf("engine must be set")
	}

	if logger == nil {
		return nil, fmt.Errorf("logger must be set")
	}

	st := &Storage{
		engine: engine,
		logger: logger,
	}

	for _, option := range options {
		option(st)
	}

	var lastLSN int64
	// TODO: Nil is not nil: https://yourbasic.org/golang/gotcha-why-nil-error-not-equal-nil/
	if st.wal != (*wal.WAL)(nil) {
		logs, err := st.wal.Recover()
		if err != nil {
			logger.Error("failed to recover data from WAL", zap.Error(err))
		} else {
			lastLSN = st.applyData(logs)
		}
	}

	if st.stream != nil {
		go func() {
			for logs := range st.stream {
				_ = st.applyData(logs)
			}
		}()
	}

	st.generator = NewIDGenerator(lastLSN)

	return st, nil
}

// Set ...
func (s *Storage) Set(ctx context.Context, key, value string) error {
	if s.replica != nil && !s.replica.IsMaster() {
		return ErrorMutableTX
	} else if ctx.Err() != nil {
		return ctx.Err()
	}

	txID := s.generator.Generate()
	ctx = common.ContextWithTxID(ctx, txID)

	if s.wal != (*wal.WAL)(nil) {
		futureResponse := s.wal.Set(ctx, key, value)
		if err := futureResponse.Get(); err != nil {
			return err
		}
	}

	s.engine.Set(ctx, key, value)

	return nil
}

// Get ...
func (s *Storage) Get(ctx context.Context, key string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	txID := s.generator.Generate()
	ctx = common.ContextWithTxID(ctx, txID)

	val, found := s.engine.Get(ctx, key)
	if !found {
		return "", fmt.Errorf("not found")
	}

	return val, nil
}

// Del ...
func (s *Storage) Del(ctx context.Context, key string) error {
	if s.replica != nil && !s.replica.IsMaster() {
		return ErrorMutableTX
	} else if ctx.Err() != nil {
		return ctx.Err()
	}

	txID := s.generator.Generate()
	ctx = common.ContextWithTxID(ctx, txID)

	if s.wal != (*wal.WAL)(nil) {
		futureResponse := s.wal.Del(ctx, key)
		if err := futureResponse.Get(); err != nil {
			return err
		}
	}

	s.engine.Del(ctx, key)

	return nil
}

func (s *Storage) applyData(logs []wal.Log) int64 {
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
