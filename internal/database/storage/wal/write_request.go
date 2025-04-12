package wal

import "database-simon/internal/concurrency"

// WriteRequest ...
type WriteRequest struct {
	log     Log
	promise concurrency.PromiseError
}

// NewWriteRequest ...
func NewWriteRequest(lsn int64, commandID string, args []string) WriteRequest {
	return WriteRequest{
		log: Log{
			LSN:       lsn,
			CommandID: commandID,
			Arguments: args,
		},
		promise: concurrency.NewPromise[error](),
	}
}

// Log ...
func (l *WriteRequest) Log() Log {
	return l.log
}

// SetResponse ...
func (l *WriteRequest) SetResponse(err error) {
	l.promise.Set(err)
}

// FutureResponse ...
func (l *WriteRequest) FutureResponse() concurrency.FutureError {
	return l.promise.GetFuture()
}
