package memory

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"concurrency/internal/database/storage"
)

func TestNewMemory(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		logger         *zap.Logger
		expectedErr    error
		expectedNilObj bool
	}{
		"create engine without logger": {
			expectedErr:    errors.New("logger must be set"),
			expectedNilObj: true,
		},
		"create engine without options": {
			logger:      zap.NewNop(),
			expectedErr: nil,
		},
		"create engine with partitions": {
			logger:      zap.NewNop(),
			expectedErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			engine, err := NewMemory(test.logger)
			assert.Equal(t, test.expectedErr, err)
			if test.expectedNilObj {
				assert.Nil(t, engine)
			} else {
				assert.NotNil(t, engine)
			}
		})
	}
}

func TestEngineSet(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		engine storage.Engine
		key    string
		value  string
	}{
		"set": {
			engine: func() storage.Engine {
				engine, err := NewMemory(zap.NewNop())
				require.NoError(t, err)
				return engine
			}(),
			key: "key",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			test.engine.Set(ctx, test.key, test.value)
			value, found := test.engine.Get(ctx, test.key)
			assert.True(t, found)
			assert.Equal(t, test.value, value)
		})
	}
}

func TestEngineGet(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		engine storage.Engine
		key    string
	}{
		"get with single partition": {
			engine: func() storage.Engine {
				engine, err := NewMemory(zap.NewNop())
				require.NoError(t, err)
				return engine
			}(),
			key: "key",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			value, found := test.engine.Get(ctx, test.key)
			assert.False(t, found)
			assert.Empty(t, value)
		})
	}
}

func TestEngineDel(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		engine storage.Engine
		key    string
	}{
		"del with single partition": {
			engine: func() storage.Engine {
				engine, err := NewMemory(zap.NewNop())
				require.NoError(t, err)
				return engine
			}(),
			key: "key",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			test.engine.Del(ctx, test.key)
			value, found := test.engine.Get(ctx, test.key)
			assert.False(t, found)
			assert.Empty(t, value)
		})
	}
}
