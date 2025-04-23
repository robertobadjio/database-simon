package storage

import (
	"context"
	wal2 "database-simon/internal/database/storage/wal"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestStorage_New(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)

	tests := map[string]struct {
		engine Engine
		logger *zap.Logger
		wal    *wal2.WAL

		expectedErr    error
		expectedNilObj bool
	}{
		"create storage without engine": {
			expectedErr:    errors.New("engine must be set"),
			expectedNilObj: true,
		},
		"create storage without logger": {
			engine:         NewMockengine(controller),
			expectedErr:    errors.New("logger must be set"),
			expectedNilObj: true,
		},
		"create storage": {
			engine:         NewMockengine(controller),
			logger:         zap.NewNop(),
			expectedErr:    nil,
			expectedNilObj: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			stor, err := NewStorage(test.engine, test.logger, WithWAL(test.wal))
			assert.Equal(t, test.expectedErr, err)
			if test.expectedNilObj {
				assert.Nil(t, stor)
			} else {
				assert.NotNil(t, stor)
			}
		})
	}
}

func TestStorage_Set(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)

	tests := map[string]struct {
		engine func() engine
		wal    *wal2.WAL

		expectedErr error
	}{
		"set": {
			engine: func() engine {
				eng := NewMockengine(controller)
				eng.EXPECT().
					Set(gomock.Any(), "key", "value")
				return eng
			},
			expectedErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			stor, err := NewStorage(test.engine(), zap.NewNop(), WithWAL(test.wal))
			require.NoError(t, err)

			err = stor.Set(context.Background(), "key", "value")
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func TestStorage_Get(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)

	tests := map[string]struct {
		engine func() engine
		wal    *wal2.WAL

		expectedValue string
		expectedErr   error
	}{
		"get with exiting key": {
			engine: func() engine {
				engine := NewMockengine(controller)
				engine.EXPECT().
					Get(gomock.Any(), "key").
					Return("value", true)
				return engine
			},
			expectedErr:   nil,
			expectedValue: "value",
		},
		"get with non-existent key": {
			engine: func() engine {
				eng := NewMockengine(controller)
				eng.EXPECT().
					Get(gomock.Any(), "key").
					Return("", false)
				return eng
			},
			expectedErr:   errors.New("not found"),
			expectedValue: "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			stor, err := NewStorage(test.engine(), zap.NewNop(), WithWAL(test.wal))
			require.NoError(t, err)

			val, err := stor.Get(context.Background(), "key")
			assert.Equal(t, test.expectedValue, val)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func TestStorage_Del(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)

	tests := map[string]struct {
		engine func() engine
		wal    *wal2.WAL

		expectedErr error
	}{
		"del": {
			engine: func() engine {
				eng := NewMockengine(controller)
				eng.EXPECT().
					Del(gomock.Any(), "key")
				return eng
			},
			expectedErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			stor, err := NewStorage(test.engine(), zap.NewNop(), WithWAL(test.wal))
			require.NoError(t, err)

			err = stor.Del(context.Background(), "key")
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
