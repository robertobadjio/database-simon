package wal

import (
	"bytes"
	"encoding/gob"
)

// Log ...
type Log struct {
	LSN       int64
	CommandID string
	Arguments []string
}

// Encode ...
func (l *Log) Encode(buffer *bytes.Buffer) error {
	encoder := gob.NewEncoder(buffer)
	return encoder.Encode(*l)
}

// Decode ...
func (l *Log) Decode(buffer *bytes.Buffer) error {
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(l)
}
