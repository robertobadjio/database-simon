package database

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"concurrency/internal/database/compute"
	"concurrency/internal/database/storage"
)

func TestNewDatabase(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)

	tests := map[string]struct {
		compute compute.Compute
		storage storage.Storage
		logger  *zap.Logger

		expectedErr    error
		expectedNilObj bool
	}{
		"create database without compute layer": {
			expectedErr:    errors.New("logger is invalid"),
			expectedNilObj: true,
		},
		"create database without storage layer": {
			compute:        compute.NewMockCompute(controller),
			expectedErr:    errors.New("logger is invalid"),
			expectedNilObj: true,
		},
		"create database without logger": {
			compute:        compute.NewMockCompute(controller),
			storage:        storage.NewMockStorage(controller),
			expectedErr:    errors.New("logger is invalid"),
			expectedNilObj: true,
		},
		"create database": {
			compute: compute.NewMockCompute(controller),
			storage: storage.NewMockStorage(controller),
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

		comp func() compute.Compute
		stor func() storage.Storage

		expectedResponse string
	}{
		"handle incorrect query": {
			query: "TRUNCATE",
			comp: func() compute.Compute {
				comp := compute.NewMockCompute(controller)
				comp.EXPECT().
					Parse(gomock.Any(), "TRUNCATE").
					Return(nil, errors.New("unknown command"))
				return comp
			},
			stor:             func() storage.Storage { return storage.NewMockStorage(controller) },
			expectedResponse: "",
		},
		"handle set query with error from storage": {
			query: "SET key value",
			comp: func() compute.Compute {
				comp := compute.NewMockCompute(controller)
				comp.EXPECT().
					Parse(gomock.Any(), "SET key value").
					Return(compute.NewQuery(
						compute.SetCommand,
						[]string{"key", "value"},
					), nil)
				return comp
			},
			stor: func() storage.Storage {
				comp := storage.NewMockStorage(controller)
				comp.EXPECT().
					Set(gomock.Any(), "key", "value").
					Return(errors.New("storage error"))
				return comp
			},
			expectedResponse: "",
		},
		"handle set query": {
			query: "SET key value",
			comp: func() compute.Compute {
				comp := compute.NewMockCompute(controller)
				comp.EXPECT().
					Parse(context.Background(), "SET key value").
					Return(compute.NewQuery(
						compute.SetCommand,
						[]string{"key", "value"},
					), nil)
				return comp
			},
			stor: func() storage.Storage {
				stor := storage.NewMockStorage(controller)
				stor.EXPECT().
					Set(gomock.Any(), "key", "value").
					Return(nil)
				return stor
			},
			expectedResponse: "",
		},
		"handle del query with error from storage": {
			query: "DEL key",
			comp: func() compute.Compute {
				comp := compute.NewMockCompute(controller)
				comp.EXPECT().
					Parse(context.Background(), "DEL key").
					Return(compute.NewQuery(
						compute.DelCommand,
						[]string{"key"},
					), nil)
				return comp
			},
			stor: func() storage.Storage {
				stor := storage.NewMockStorage(controller)
				stor.EXPECT().
					Del(gomock.Any(), "key").
					Return(errors.New("storage error"))
				return stor
			},
			expectedResponse: "",
		},
		"handle del query": {
			query: "DEL key",
			comp: func() compute.Compute {
				computeLayer := compute.NewMockCompute(controller)
				computeLayer.EXPECT().
					Parse(context.Background(), "DEL key").
					Return(compute.NewQuery(
						compute.DelCommand,
						[]string{"key"},
					), nil)
				return computeLayer
			},
			stor: func() storage.Storage {
				stor := storage.NewMockStorage(controller)
				stor.EXPECT().
					Del(gomock.Any(), "key").
					Return(nil)
				return stor
			},
			expectedResponse: "",
		},
		"handle get query with error from storage": {
			query: "GET key",
			comp: func() compute.Compute {
				comp := compute.NewMockCompute(controller)
				comp.EXPECT().
					Parse(context.Background(), "GET key").
					Return(compute.NewQuery(
						compute.GetCommand,
						[]string{"key"},
					), nil)
				return comp
			},
			stor: func() storage.Storage {
				storageLayer := storage.NewMockStorage(controller)
				storageLayer.EXPECT().
					Get(gomock.Any(), "key").
					Return("", errors.New("storage error"))
				return storageLayer
			},
			expectedResponse: "",
		},
		"handle get query with not found error from storage": {
			query: "GET key",
			comp: func() compute.Compute {
				computeLayer := compute.NewMockCompute(controller)
				computeLayer.EXPECT().
					Parse(gomock.Any(), "GET key").
					Return(compute.NewQuery(
						compute.GetCommand,
						[]string{"key"},
					), nil)
				return computeLayer
			},
			stor: func() storage.Storage {
				stor := storage.NewMockStorage(controller)
				stor.EXPECT().
					Get(gomock.Any(), "key").
					Return("", errors.New("not found"))
				return stor
			},
			expectedResponse: "",
		},
		"handle get query": {
			query: "GET key",
			comp: func() compute.Compute {
				comp := compute.NewMockCompute(controller)
				comp.EXPECT().
					Parse(gomock.Any(), "GET key").
					Return(compute.NewQuery(
						compute.GetCommand,
						[]string{"key"},
					), nil)
				return comp
			},
			stor: func() storage.Storage {
				stor := storage.NewMockStorage(controller)
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
