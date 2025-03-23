package database

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"concurrency/internal/database/compute"
	"concurrency/internal/database/storage"
)

func TestNewDatabase(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

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
			compute:        compute.NewMockCompute(ctrl),
			expectedErr:    errors.New("logger is invalid"),
			expectedNilObj: true,
		},
		"create database without logger": {
			compute:        compute.NewMockCompute(ctrl),
			storage:        storage.NewMockStorage(ctrl),
			expectedErr:    errors.New("logger is invalid"),
			expectedNilObj: true,
		},
		"create database": {
			compute: compute.NewMockCompute(ctrl),
			storage: storage.NewMockStorage(ctrl),
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
