package concurrency

// PromiseError ...
type PromiseError = Promise[error]

// Promise ...
type Promise[T any] struct {
	result   chan T
	promised bool
}

// NewPromise ...
func NewPromise[T any]() Promise[T] {
	return Promise[T]{
		result: make(chan T, 1),
	}
}

// Set ...
func (p *Promise[T]) Set(value T) {
	if p.promised {
		return
	}

	p.promised = true
	p.result <- value
	close(p.result)
}

// GetFuture ...
func (p *Promise[T]) GetFuture() Future[T] {
	return NewFuture(p.result)
}
