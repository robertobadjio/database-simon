package memory

// EngineOption ...
type EngineOption func(memory *Memory)

// WithPartitions ...
func WithPartitions(partitionsNumber int) EngineOption {
	return func(engine *Memory) {
		engine.partitions = make([]*HashTable, partitionsNumber)
		for i := 0; i < partitionsNumber; i++ {
			engine.partitions[i] = NewHashTable()
		}
	}
}
