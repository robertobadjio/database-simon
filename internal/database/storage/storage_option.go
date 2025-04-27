package storage

import "database-simon/internal/database/storage/wal"

// Option ...
type Option func(*Storage)

// WithWAL ...
func WithWAL(wal walI) Option {
	return func(storage *Storage) {
		storage.wal = wal
	}
}

// WithReplication ...
func WithReplication(replica replica) Option {
	return func(storage *Storage) {
		storage.replica = replica
	}
}

// WithReplicationStream ...
func WithReplicationStream(stream <-chan []wal.Log) Option {
	return func(storage *Storage) {
		storage.stream = stream
	}
}
