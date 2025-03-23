package database

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"concurrency/internal/database/compute"
	"concurrency/internal/database/storage"
)

type Database interface {
	HandleQuery(ctx context.Context, queryStr string) (string, error)
}

type database struct {
	comp   compute.Compute
	stor   storage.Storage
	logger *zap.Logger
}

func NewDatabase(logger *zap.Logger, comp compute.Compute, stor storage.Storage) (Database, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is invalid")
	}

	if comp == nil {
		return nil, fmt.Errorf("compute is invalid")
	}

	if stor == nil {
		return nil, fmt.Errorf("storage is invalid")
	}

	return &database{
		logger: logger,
		comp:   comp,
		stor:   stor,
	}, nil
}

func (db *database) HandleQuery(ctx context.Context, queryStr string) (string, error) {
	query, err := db.comp.Parse(ctx, queryStr)
	if err != nil {
		return "", fmt.Errorf("error parsing")
	}

	switch query.Command() {
	case compute.SetCommand:
		return db.handlerSetQuery(ctx, query)
	case compute.GetCommand:
		return db.handlerGetQuery(ctx, query)
	case compute.DelCommand:
		return db.handlerDelQuery(ctx, query)
	}

	return "", fmt.Errorf("error handle query")
}

func (db *database) handlerSetQuery(ctx context.Context, query compute.Query) (string, error) {
	err := db.stor.Set(ctx, query.Arguments()[0], query.Arguments()[1])
	if err != nil {
		return "", err
	}
	return "", nil
}

func (db *database) handlerGetQuery(ctx context.Context, query compute.Query) (string, error) {
	value, err := db.stor.Get(ctx, query.Arguments()[0])
	if err != nil {
		return "", fmt.Errorf("error hadnle get query: %w", err)
	}
	return value, nil
}

func (db *database) handlerDelQuery(ctx context.Context, query compute.Query) (string, error) {
	err := db.stor.Del(ctx, query.Arguments()[0])
	if err != nil {
		return "", err
	}
	return "", nil
}
