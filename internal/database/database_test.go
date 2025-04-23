package database

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"database-simon/internal/database/compute"
)

func TestNewDatabase(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)

	tests := map[string]struct {
		compute computeLayer
		storage storageLayer
		logger  *zap.Logger

		expectedErr    error
		expectedNilObj bool
	}{
		"create database without compute layer": {
			expectedErr:    errors.New("logger is invalid"),
			expectedNilObj: true,
		},
		"create database without storage layer": {
			compute:        NewMockcomputeLayer(controller),
			expectedErr:    errors.New("logger is invalid"),
			expectedNilObj: true,
		},
		"create database without logger": {
			compute:        NewMockcomputeLayer(controller),
			storage:        NewMockstorageLayer(controller),
			expectedErr:    errors.New("logger is invalid"),
			expectedNilObj: true,
		},
		"create database": {
			compute: NewMockcomputeLayer(controller),
			storage: NewMockstorageLayer(controller),
			logger:  zap.NewNop(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			stor, err := NewDatabase(test.logger, test.compute, test.storage)
			assert.Equal(t, test.expectedErr, err)

			if test.expectedNilObj {
				assert.Nil(t, stor)
			} else {
				assert.NotNil(t, stor)
			}
		})
	}
}

func TestHandleQuery(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)

	tests := map[string]struct {
		query string

		comp func() computeLayer
		stor func() storageLayer

		expectedResponse string
	}{
		"handle incorrect query": {
			query: "TRUNCATE",
			comp: func() computeLayer {
				comp := NewMockcomputeLayer(controller)
				comp.EXPECT().
					Parse(gomock.Any(), "TRUNCATE").
					Return(nil, errors.New("unknown command"))
				return comp
			},
			stor:             func() storageLayer { return NewMockstorageLayer(controller) },
			expectedResponse: "[error]",
		},
		"handle set query with error from storage": {
			query: "SET key value",
			comp: func() computeLayer {
				comp := NewMockcomputeLayer(controller)
				comp.EXPECT().
					Parse(gomock.Any(), "SET key value").
					Return(compute.NewQuery(
						compute.SetCommand,
						[]string{"key", "value"},
					), nil)
				return comp
			},
			stor: func() storageLayer {
				stop := NewMockstorageLayer(controller)
				stop.EXPECT().
					Set(gomock.Any(), "key", "value").
					Return(errors.New("storage error"))
				return stop
			},
			expectedResponse: "[error]",
		},
		"handle set query": {
			query: "SET key value",
			comp: func() computeLayer {
				comp := NewMockcomputeLayer(controller)
				comp.EXPECT().
					Parse(context.Background(), "SET key value").
					Return(compute.NewQuery(
						compute.SetCommand,
						[]string{"key", "value"},
					), nil)
				return comp
			},
			stor: func() storageLayer {
				stor := NewMockstorageLayer(controller)
				stor.EXPECT().
					Set(gomock.Any(), "key", "value").
					Return(nil)
				return stor
			},
			expectedResponse: "[ok]",
		},
		"handle del query with error from storage": {
			query: "DEL key",
			comp: func() computeLayer {
				comp := NewMockcomputeLayer(controller)
				comp.EXPECT().
					Parse(context.Background(), "DEL key").
					Return(compute.NewQuery(
						compute.DelCommand,
						[]string{"key"},
					), nil)
				return comp
			},
			stor: func() storageLayer {
				stor := NewMockstorageLayer(controller)
				stor.EXPECT().
					Del(gomock.Any(), "key").
					Return(errors.New("storage error"))
				return stor
			},
			expectedResponse: "[error]",
		},
		"handle del query": {
			query: "DEL key",
			comp: func() computeLayer {
				comp := NewMockcomputeLayer(controller)
				comp.EXPECT().
					Parse(context.Background(), "DEL key").
					Return(compute.NewQuery(
						compute.DelCommand,
						[]string{"key"},
					), nil)
				return comp
			},
			stor: func() storageLayer {
				stor := NewMockstorageLayer(controller)
				stor.EXPECT().
					Del(gomock.Any(), "key").
					Return(nil)
				return stor
			},
			expectedResponse: "[ok]",
		},
		"handle get query with error from storage": {
			query: "GET key",
			comp: func() computeLayer {
				comp := NewMockcomputeLayer(controller)
				comp.EXPECT().
					Parse(context.Background(), "GET key").
					Return(compute.NewQuery(
						compute.GetCommand,
						[]string{"key"},
					), nil)
				return comp
			},
			stor: func() storageLayer {
				stor := NewMockstorageLayer(controller)
				stor.EXPECT().
					Get(gomock.Any(), "key").
					Return("", errors.New("storage error"))
				return stor
			},
			expectedResponse: "[not found]",
		},
		"handle get query with not found error from storage": {
			query: "GET key",
			comp: func() computeLayer {
				comp := NewMockcomputeLayer(controller)
				comp.EXPECT().
					Parse(gomock.Any(), "GET key").
					Return(compute.NewQuery(
						compute.GetCommand,
						[]string{"key"},
					), nil)
				return comp
			},
			stor: func() storageLayer {
				stor := NewMockstorageLayer(controller)
				stor.EXPECT().
					Get(gomock.Any(), "key").
					Return("", errors.New("not found"))
				return stor
			},
			expectedResponse: "[not found]",
		},
		"handle get query": {
			query: "GET key",
			comp: func() computeLayer {
				comp := NewMockcomputeLayer(controller)
				comp.EXPECT().
					Parse(gomock.Any(), "GET key").
					Return(compute.NewQuery(
						compute.GetCommand,
						[]string{"key"},
					), nil)
				return comp
			},
			stor: func() storageLayer {
				stor := NewMockstorageLayer(controller)
				stor.EXPECT().
					Get(gomock.Any(), "key").
					Return("value", nil)
				return stor
			},
			expectedResponse: "value",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			stor, err := NewDatabase(zap.NewNop(), test.comp(), test.stor())
			require.NoError(t, err)

			response, _ := stor.HandleQuery(context.Background(), test.query)
			assert.Equal(t, test.expectedResponse, response)
		})
	}
}
