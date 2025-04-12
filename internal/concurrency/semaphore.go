package concurrency

// Semaphore ...
type Semaphore struct {
	tickets chan struct{}
}

// NewSemaphore ...
func NewSemaphore(ticketsNumber int) Semaphore {
	return Semaphore{
		tickets: make(chan struct{}, ticketsNumber),
	}
}

// Acquire ...
func (s *Semaphore) Acquire() {
	if s == nil || s.tickets == nil {
		return
	}

	s.tickets <- struct{}{}
}

// Release ...
func (s *Semaphore) Release() {
	if s == nil || s.tickets == nil {
		return
	}

	<-s.tickets
}

// WithAcquire ...
func (s *Semaphore) WithAcquire(action func()) {
	if s == nil || action == nil {
		return
	}

	s.Acquire()
	action()
	s.Release()
}
