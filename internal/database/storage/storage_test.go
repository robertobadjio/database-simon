package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestNewStorage(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)

	tests := map[string]struct {
		engine Engine
		logger *zap.Logger

		expectedErr    error
		expectedNilObj bool
	}{
		"create storage without engine": {
			expectedErr:    errors.New("engine must be set"),
			expectedNilObj: true,
		},
		"create storage without logger": {
			engine:         NewMockEngine(controller),
			expectedErr:    errors.New("logger must be set"),
			expectedNilObj: true,
		},
		"create storage": {
			engine:         NewMockEngine(controller),
			logger:         zap.NewNop(),
			expectedErr:    nil,
			expectedNilObj: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			stor, err := NewStorage(test.engine, test.logger)
			assert.Equal(t, test.expectedErr, err)
			if test.expectedNilObj {
				assert.Nil(t, stor)
			} else {
				assert.NotNil(t, stor)
			}
		})
	}
}

func TestStorageSet(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)

	tests := map[string]struct {
		engine      func() Engine
		expectedErr error
	}{
		"set": {
			engine: func() Engine {
				engine := NewMockEngine(controller)
				engine.EXPECT().
					Set(gomock.Any(), "key", "value")
				return engine
			},
			expectedErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			stor, err := NewStorage(test.engine(), zap.NewNop())
			require.NoError(t, err)

			err = stor.Set(context.Background(), "key", "value")
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func TestStorageGet(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)

	tests := map[string]struct {
		engine        func() Engine
		expectedValue string
		expectedErr   error
	}{
		"get with exiting key": {
			engine: func() Engine {
				engine := NewMockEngine(controller)
				engine.EXPECT().
					Get(gomock.Any(), "key").
					Return("value", true)
				return engine
			},
			expectedErr:   nil,
			expectedValue: "value",
		},
		"get with non-existent key": {
			engine: func() Engine {
				engine := NewMockEngine(controller)
				engine.EXPECT().
					Get(gomock.Any(), "key").
					Return("", false)
				return engine
			},
			expectedErr:   errors.New("not found"),
			expectedValue: "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			stor, err := NewStorage(test.engine(), zap.NewNop())
			require.NoError(t, err)

			val, err := stor.Get(context.Background(), "key")
			assert.Equal(t, test.expectedValue, val)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func TestStorageDel(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)

	tests := map[string]struct {
		engine      func() Engine
		expectedErr error
	}{
		"del": {
			engine: func() Engine {
				engine := NewMockEngine(controller)
				engine.EXPECT().
					Del(gomock.Any(), "key")
				return engine
			},
			expectedErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			stor, err := NewStorage(test.engine(), zap.NewNop())
			require.NoError(t, err)

			err = stor.Del(context.Background(), "key")
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
