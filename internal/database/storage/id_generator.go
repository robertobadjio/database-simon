package storage

import (
	"math"
	"sync/atomic"
)

// IDGenerator ...
type IDGenerator struct {
	counter atomic.Int64
}

// NewIDGenerator ...
func NewIDGenerator(previousID int64) *IDGenerator {
	generator := &IDGenerator{}
	generator.counter.Store(previousID)

	return generator
}

// Generate ...
func (g *IDGenerator) Generate() int64 {
	g.counter.CompareAndSwap(math.MaxInt64, 0)
	return g.counter.Add(1)
}
