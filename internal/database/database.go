package database

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"database-simon/internal/database/compute"
)

const (
	errorResult    = "[error]"
	okResult       = "[ok]"
	notFoundResult = "[not found]"
)

type computeLayer interface {
	Parse(ctx context.Context, queryStr string) (compute.Query, error)
}

type storageLayer interface {
	Set(context.Context, string, string) error
	Get(context.Context, string) (string, error)
	Del(context.Context, string) error
}

// Database ...
type Database struct {
	comp   computeLayer
	stor   storageLayer
	logger *zap.Logger
}

// NewDatabase ...
func NewDatabase(logger *zap.Logger, comp compute.Compute, stor storageLayer) (*Database, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is invalid")
	}

	if comp == nil {
		return nil, fmt.Errorf("compute is invalid")
	}

	if stor == nil {
		return nil, fmt.Errorf("storage is invalid")
	}

	return &Database{
		logger: logger,
		comp:   comp,
		stor:   stor,
	}, nil
}

// HandleQuery ...
func (db *Database) HandleQuery(ctx context.Context, queryStr string) (string, error) {
	query, err := db.comp.Parse(ctx, queryStr)
	if err != nil {
		return errorResult, fmt.Errorf("error parsing: %w", err)
	}

	switch query.Command() {
	case compute.SetCommand:
		_, errSet := db.handlerSetQuery(ctx, query)
		if errSet != nil {
			return errorResult, errSet
		}
		return okResult, nil
	case compute.GetCommand:
		res, errGet := db.handlerGetQuery(ctx, query)
		if errGet != nil {
			return notFoundResult, errGet
		}
		return res, nil
	case compute.DelCommand:
		_, errDel := db.handlerDelQuery(ctx, query)
		if errDel != nil {
			return errorResult, errDel
		}
		return okResult, nil
	}

	return errorResult, fmt.Errorf("error handle query")
}

func (db *Database) handlerSetQuery(ctx context.Context, query compute.Query) (string, error) {
	err := db.stor.Set(ctx, query.Arguments()[0], query.Arguments()[1])
	if err != nil {
		return "", err
	}
	return "", nil
}

func (db *Database) handlerGetQuery(ctx context.Context, query compute.Query) (string, error) {
	value, err := db.stor.Get(ctx, query.Arguments()[0])
	if err != nil {
		return "", fmt.Errorf("error handle get query: %w", err)
	}
	return value, nil
}

func (db *Database) handlerDelQuery(ctx context.Context, query compute.Query) (string, error) {
	err := db.stor.Del(ctx, query.Arguments()[0])
	if err != nil {
		return "", err
	}
	return "", nil
}
