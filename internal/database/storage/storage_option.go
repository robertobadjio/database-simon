package storage

// Option ...
type Option func(*Storage)

// WithWAL ...
func WithWAL(wal walI) Option {
	return func(storage *Storage) {
		storage.wal = wal
	}
}
