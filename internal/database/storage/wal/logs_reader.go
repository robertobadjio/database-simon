package wal

import (
	"bytes"
	"fmt"
	"sort"
)

type segmentsDirectory interface {
	ForEach(func([]byte) error) error
}

// LogsReader ...
type LogsReader struct {
	segmentsDirectory segmentsDirectory
}

// NewLogsReader ...
func NewLogsReader(segmentsDirectory segmentsDirectory) (*LogsReader, error) {
	if segmentsDirectory == nil {
		return nil, fmt.Errorf("segments directory is invalid")
	}

	return &LogsReader{
		segmentsDirectory: segmentsDirectory,
	}, nil
}

// Read ...
func (r *LogsReader) Read() ([]Log, error) {
	var logs []Log
	err := r.segmentsDirectory.ForEach(func(data []byte) error {
		var err error
		logs, err = r.readSegment(logs, data)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to read segments: %w", err)
	}

	// TODO: need to check invariant for sorting
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].LSN < logs[j].LSN
	})

	return logs, nil
}

func (r *LogsReader) readSegment(logs []Log, data []byte) ([]Log, error) {
	buffer := bytes.NewBuffer(data)
	for buffer.Len() > 0 {
		var log Log
		if err := log.Decode(buffer); err != nil {
			return nil, fmt.Errorf("failed to parse logs data: %w", err)
		}

		logs = append(logs, log)
	}

	return logs, nil
}
