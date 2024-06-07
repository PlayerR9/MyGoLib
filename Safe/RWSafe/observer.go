package RWSafe

import (
	"sync"

	uc "github.com/PlayerR9/MyGoLib/Units/Common"
)

// Observer is the interface that wraps the Notify method.
type Observer[T any] interface {
	// Notify notifies the observer of a change.
	//
	// Parameters:
	//   - change: The change that occurred.
	Notify(change T)
}

// Subject is the subject that observers observe.
type Subject[T any] struct {
	// observers is the list of observers.
	observers []Observer[T]

	// state is the state of the subject.
	state T

	// mu is the mutex to synchronize access to the subject.
	mu sync.RWMutex
}

// Copy implements the Copier interface.
//
// However, the obsevers are not copied.
func (s *Subject[T]) Copy() uc.Copier {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return &Subject[T]{
		observers: make([]Observer[T], 0),
		state:     s.state,
	}
}

// NewSubject creates a new subject.
//
// Parameters:
//   - state: The initial state of the subject.
//
// Returns:
//
//	*Subject[T]: A new subject.
func NewSubject[T any](state T) *Subject[T] {
	return &Subject[T]{
		observers: make([]Observer[T], 0),
		state:     state,
	}
}

// Attach attaches an observer to the subject.
//
// Parameters:
//   - o: The observer to attach.
func (s *Subject[T]) Attach(o Observer[T]) {
	if o == nil {
		return
	}

	s.observers = append(s.observers, o)
}

// SetState sets the state of the subject.
//
// Parameters:
//   - state: The new state of the subject.
func (s *Subject[T]) SetState(state T) {
	s.mu.Lock()
	s.state = state
	s.mu.Unlock()

	s.NotifyAll()
}

// GetState gets the state of the subject.
//
// Returns:
//   - T: The state of the subject.
func (s *Subject[T]) GetState() T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.state
}

// ModifyState modifies the state of the subject.
//
// Parameters:
//   - f: The function to modify the state of the subject.
func (s *Subject[T]) ModifyState(f func(T) T) {
	s.mu.RLock()
	curr := s.state
	s.mu.RUnlock()

	new := f(curr)

	s.mu.Lock()
	s.state = new
	s.mu.Unlock()

	s.NotifyAll()
}

// NotifyAll notifies all observers of a change.
func (s *Subject[T]) NotifyAll() {
	var wg sync.WaitGroup

	wg.Add(len(s.observers))

	s.mu.RLock()
	state := s.state
	s.mu.RUnlock()

	for _, observer := range s.observers {
		go func(observer Observer[T]) {
			defer wg.Done()

			observer.Notify(state)
		}(observer)
	}

	wg.Wait()
}
