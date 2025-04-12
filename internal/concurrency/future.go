package concurrency

// FutureError ...
type FutureError = Future[error]

// Future ...
type Future[T any] struct {
	result <-chan T
}

// NewFuture ...
func NewFuture[T any](result <-chan T) Future[T] {
	return Future[T]{
		result: result,
	}
}

// Get ...
func (f *Future[T]) Get() T {
	return <-f.result
}
