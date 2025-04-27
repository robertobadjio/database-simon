package replication

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// Request ...
type Request struct {
	LastSegmentName string
}

// NewRequest ...
func NewRequest(lastSegmentName string) Request {
	return Request{
		LastSegmentName: lastSegmentName,
	}
}

// Response ...
type Response struct {
	Succeed     bool
	SegmentName string
	SegmentData []byte
}

// NewResponse ...
func NewResponse(succeed bool, segmentName string, segmentData []byte) Response {
	return Response{
		Succeed:     succeed,
		SegmentName: segmentName,
		SegmentData: segmentData,
	}
}

// Encode ...
func Encode[ProtocolObject Request | Response](object *ProtocolObject) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(object); err != nil {
		return nil, fmt.Errorf("failed to encode object: %w", err)
	}

	return buffer.Bytes(), nil
}

// Decode ...
func Decode[ProtocolObject Request | Response](object *ProtocolObject, data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&object); err != nil {
		return fmt.Errorf("failed to decode object: %w", err)
	}

	return nil
}
