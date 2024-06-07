package RWSafe

import (
	"sync"

	uc "github.com/PlayerR9/MyGoLib/Units/Common"
)

// Data is a thread-safe variable that allows multiple goroutines to
// read and write to a variable.
type Data[T any] struct {
	// state is the state of the subject.
	state T

	// mu is the mutex to synchronize access to the subject.
	mu sync.RWMutex
}

// Copy implements the Copier interface.
//
// However, the obsevers are not copied.
func (s *Data[T]) Copy() uc.Copier {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return &Data[T]{
		state: s.state,
	}
}

// NewData creates a new subject.
//
// Parameters:
//   - state: The initial state of the subject.
//
// Returns:
//
//	*Data[T]: A new subject.
func NewData[T any](state T) *Data[T] {
	return &Data[T]{
		state: state,
	}
}

// Set sets the state of the subject.
//
// Parameters:
//   - state: The new state of the subject.
func (s *Data[T]) Set(state T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.state = state
}

// Get gets the state of the subject.
//
// Returns:
//   - T: The state of the subject.
func (s *Data[T]) Get() T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.state
}

// Modify modifies the state of the subject.
//
// Parameters:
//   - f: The function to modify the state of the subject.
func (s *Data[T]) Modify(f func(T) T) {
	s.mu.RLock()
	curr := s.state
	s.mu.RUnlock()

	new := f(curr)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.state = new
}

// DoRead is a method of the Data type. It is used to perform a read
// operation on the value stored in the Data.
// Through the function parameter, the caller can access the value in a
// read-only manner.
//
// Parameters:
//
//   - f: A function that takes a value of type T as a parameter and
//     returns nothing.
func (s *Data[T]) DoRead(f func(T)) {
	s.mu.RLock()
	value := s.state
	s.mu.RUnlock()

	f(value)
}
