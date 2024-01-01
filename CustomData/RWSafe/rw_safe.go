package RWSafe

import "sync"

type RWSafe[T any] struct {
	value T
	mutex sync.RWMutex
}

func NewRWSafe[T any](value T) RWSafe[T] {
	return RWSafe[T]{
		value: value,
	}
}

func (rw *RWSafe[T]) Get() T {
	rw.mutex.RLock()
	defer rw.mutex.RUnlock()

	return rw.value
}

func (rw *RWSafe[T]) Set(value T) {
	rw.mutex.Lock()
	rw.value = value
	rw.mutex.Unlock()
}
