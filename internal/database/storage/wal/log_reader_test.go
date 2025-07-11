package wal

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.uber.org/mock/gomock"
)

func TestNewLogsReader(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)

	tests := map[string]struct {
		directory segmentsDirectory

		expectedErr    error
		expectedNilObj bool
	}{
		"create logs reader without segments directory": {
			expectedErr:    errors.New("segments directory is invalid"),
			expectedNilObj: true,
		},
		"create logs reader": {
			directory: NewMocksegmentsDirectory(controller),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			reader, err := NewLogsReader(test.directory)
			assert.Equal(t, test.expectedErr, err)
			if test.expectedNilObj {
				assert.Nil(t, reader)
			} else {
				assert.NotNil(t, reader)
			}
		})
	}
}

func TestReadWithError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("read error")

	controller := gomock.NewController(t)
	directory := NewMocksegmentsDirectory(controller)
	directory.EXPECT().
		ForEach(gomock.Any()).
		Return(expectedErr)

	reader, err := NewLogsReader(directory)
	require.NoError(t, err)
	logs, err := reader.Read()
	assert.True(t, errors.Is(err, expectedErr))
	assert.Nil(t, logs)
}

func TestRead(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)
	directory := NewMocksegmentsDirectory(controller)
	directory.EXPECT().
		ForEach(gomock.Any()).
		Return(nil)

	reader, err := NewLogsReader(directory)
	require.NoError(t, err)
	logs, err := reader.Read()
	assert.Nil(t, err)
	assert.Nil(t, logs)
}
