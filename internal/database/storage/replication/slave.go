package replication

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"database-simon/internal/database/filesystem"
	"database-simon/internal/database/storage/wal"
)

type tcpClient interface {
	Send([]byte) ([]byte, error)
	Close()
}

// Slave ...
type Slave struct {
	client tcpClient
	stream chan []wal.Log

	syncInterval    time.Duration
	walDirectory    string
	lastSegmentName string

	logger *zap.Logger
}

// NewSlave ...
func NewSlave(
	client tcpClient,
	walDirectory string,
	syncInterval time.Duration,
	logger *zap.Logger,
) (*Slave, error) {
	if client == nil {
		return nil, errors.New("tcp client must be set")
	}

	if logger == nil {
		return nil, errors.New("logger must be set")
	}

	segmentName, err := filesystem.SegmentLast(walDirectory)
	if err != nil {
		logger.Error("failed to find last WAL segment", zap.Error(err))
	}

	return &Slave{
		client:          client,
		stream:          make(chan []wal.Log),
		syncInterval:    syncInterval,
		walDirectory:    walDirectory,
		lastSegmentName: segmentName,
		logger:          logger,
	}, nil
}

// Start ...
func (s *Slave) Start(ctx context.Context) {
	ticker := time.NewTicker(s.syncInterval)

	go func() {
		defer func() {
			ticker.Stop()
			s.client.Close()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.synchronize()
			}
		}
	}()
}

// IsMaster ...
func (s *Slave) IsMaster() bool {
	return false
}

// ReplicationStream ...
func (s *Slave) ReplicationStream() <-chan []wal.Log {
	return s.stream
}

func (s *Slave) synchronize() {
	request := NewRequest(s.lastSegmentName)
	requestData, err := Encode(&request)
	if err != nil {
		s.logger.Error("failed to encode replication request", zap.Error(err))
		return
	}

	responseData, errE := s.client.Send(requestData)
	if errE != nil {
		s.logger.Error("failed to send replication request", zap.Error(errE))
		return
	}

	var response Response
	if err = Decode(&response, responseData); err != nil {
		s.logger.Error("failed to decode replication response", zap.Error(err))
		return
	}

	if response.Succeed {
		s.handleResponse(response)
	} else {
		s.logger.Error("failed to apply replication data: master error")
	}
}

func (s *Slave) handleResponse(response Response) {
	if response.SegmentName == "" {
		s.logger.Debug("no changes from replication")
		return
	}

	if err := s.saveWALSegment(response.SegmentName, response.SegmentData); err != nil {
		s.logger.Error("failed to apply replication data", zap.Error(err))
		return
	}

	if err := s.writeDataToStream(response.SegmentData); err != nil {
		s.logger.Error("failed to write data to stream", zap.Error(err))
		return
	}

	s.lastSegmentName = response.SegmentName
}

func (s *Slave) saveWALSegment(segmentName string, segmentData []byte) error {
	filename := fmt.Sprintf("%s/%s", s.walDirectory, segmentName)
	file, err := filesystem.CreateFile(filename)
	if err != nil {
		return fmt.Errorf("failed to create wal segment: %w", err)
	}

	_, err = filesystem.WriteFile(file, segmentData)
	return err
}

func (s *Slave) writeDataToStream(segmentData []byte) error {
	var logs []wal.Log
	buffer := bytes.NewBuffer(segmentData)
	for buffer.Len() > 0 {
		var log wal.Log
		if err := log.Decode(buffer); err != nil {
			return fmt.Errorf("failed to decode data: %w", err)
		}

		logs = append(logs, log)
	}
	s.stream <- logs
	return nil
}
