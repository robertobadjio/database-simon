package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"database-simon/internal/database/storage/wal"
)

func TestWithReplicationStream(t *testing.T) {
	t.Parallel()

	stream := make(<-chan []wal.Log)
	option := WithReplicationStream(stream)

	var storage Storage
	option(&storage)

	assert.Equal(t, stream, storage.stream)
}

func TestWithWAL(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)
	w := NewMockwalI(controller)
	option := WithWAL(w)

	var storage Storage
	option(&storage)

	assert.Equal(t, w, storage.wal)
}

func TestWithReplication(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)
	r := NewMockreplica(controller)
	option := WithReplication(r)

	var storage Storage
	option(&storage)

	assert.Equal(t, r, storage.replica)
}
