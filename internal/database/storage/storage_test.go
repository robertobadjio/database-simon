package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"database-simon/internal/database/storage/wal"
)

func TestStorage_New(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)

	writeAheadLog := NewMockwalI(controller)
	writeAheadLog.EXPECT().
		Recover().
		Return(nil, nil)

	tests := map[string]struct {
		engine  engine
		logger  *zap.Logger
		options []Option

		expectedErr    error
		expectedNilObj bool
	}{
		"create storage without engine": {
			expectedErr:    errors.New("engine must be set"),
			expectedNilObj: true,
		},
		"create storage without logger": {
			engine: NewMockengine(controller),

			expectedErr:    errors.New("logger must be set"),
			expectedNilObj: true,
		},
		"create engine without options": {
			engine: NewMockengine(controller),
			logger: zap.NewNop(),

			expectedErr:    nil,
			expectedNilObj: false,
		},
		"create engine with wal": {
			engine:  NewMockengine(controller),
			logger:  zap.NewNop(),
			options: []Option{WithWAL(writeAheadLog)},

			expectedErr: nil,
		},
		"create engine with replica": {
			engine:      NewMockengine(controller),
			logger:      zap.NewNop(),
			options:     []Option{WithReplication(NewMockreplica(controller))},
			expectedErr: nil,
		},
		"create engine with replication stream": {
			engine:      NewMockengine(controller),
			logger:      zap.NewNop(),
			options:     []Option{WithReplicationStream(make(<-chan []wal.Log))},
			expectedErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			stor, err := NewStorage(test.engine, test.logger, test.options...)
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

			stor, err := NewStorage(test.engine(), zap.NewNop())
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

		expectedValue string
		expectedErr   error
	}{
		"get with exiting key": {
			engine: func() engine {
				eng := NewMockengine(controller)
				eng.EXPECT().
					Get(gomock.Any(), "key").
					Return("value", true)
				return eng
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

			stor, err := NewStorage(test.engine(), zap.NewNop())
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

			stor, err := NewStorage(test.engine(), zap.NewNop())
			require.NoError(t, err)

			err = stor.Del(context.Background(), "key")
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
