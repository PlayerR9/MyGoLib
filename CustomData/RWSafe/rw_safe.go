package RWSafe

import (
	"sync"
)

// RWSafe is a thread-safe, generic data structure that allows multiple
// goroutines to read and write to a value in a synchronized manner.
type RWSafe[T any] struct {
	// The value of type T that is being protected.
	value T

	// A RWMutex to synchronize access to the value.
	mutex sync.RWMutex
}

// NewRWSafe is a function that creates and returns a new instance of a
// RWSafe.
//
// Parameters:
//
//   - value: The initial value to be stored in the RWSafe.
//
// Returns:
//
//   - *RWSafe[T]: A pointer to the newly created RWSafe.
func NewRWSafe[T any](value T) *RWSafe[T] {
	return &RWSafe[T]{
		value: value,
	}
}

// Get is a method of the RWSafe type. It is used to retrieve the value
// stored in the RWSafe.
//
// Returns:
//
//   - T: The value stored in the RWSafe.
func (rw *RWSafe[T]) Get() T {
	rw.mutex.RLock()
	defer rw.mutex.RUnlock()

	return rw.value
}

// Set is a method of the RWSafe type. It is used to set the value stored
// in the RWSafe.
//
// Parameters:
//
//   - value: The new value to be stored in the RWSafe.
func (rw *RWSafe[T]) Set(value T) {
	rw.mutex.Lock()
	rw.value = value
	rw.mutex.Unlock()
}

// DoRead is a method of the RWSafe type. It is used to perform a read
// operation on the value stored in the RWSafe.
// Through the function parameter, the caller can access the value in a
// read-only manner.
//
// Parameters:
//
//   - f: A function that takes a value of type T as a parameter and
//     returns nothing.
func (rw *RWSafe[T]) DoRead(f func(T)) {
	rw.mutex.RLock()
	defer rw.mutex.RUnlock()

	f(rw.value)
}

// DoWrite is a method of the RWSafe type. It is used to perform a write
// operation on the value stored in the RWSafe.
// Through the function parameter, the caller can access the value in a
// read-write manner.
//
// Parameters:
//
//   - f: A function that takes a pointer to a value of type T as a
//     parameter and returns nothing.
func (rw *RWSafe[T]) DoWrite(f func(T)) {
	rw.mutex.Lock()
	defer rw.mutex.Unlock()

	f(rw.value)
}
