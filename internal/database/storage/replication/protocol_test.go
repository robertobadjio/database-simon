package replication

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

const sName = "wal_1000.log"

func TestRequest(t *testing.T) {
	t.Parallel()

	lastSegmentName := sName
	initialRequest := NewRequest(lastSegmentName)
	data, err := Encode(&initialRequest)
	require.NoError(t, err)
	require.NotNil(t, data)

	var decodedRequest Request
	err = Decode(&decodedRequest, data)
	require.NoError(t, err)

	require.Equal(t, initialRequest, decodedRequest)
}

func TestResponse(t *testing.T) {
	t.Parallel()

	succeed := true
	segmentName := sName
	segmentData := []byte{'s', 'y', 'n', 'c'}
	initialResponse := NewResponse(succeed, segmentName, segmentData)
	data, err := Encode(&initialResponse)
	require.NoError(t, err)
	require.NotNil(t, data)

	var decodedResponse Response
	err = Decode(&decodedResponse, data)
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(initialResponse, decodedResponse))
}
