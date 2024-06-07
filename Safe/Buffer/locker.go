package Buffer

import (
	"sync"

	uc "github.com/PlayerR9/MyGoLib/Units/Common"
)

// Locker is a thread-safe locker that allows multiple goroutines to wait for a condition.
type Locker[T uc.Enumer] struct {
	// cond is the condition variable.
	cond *sync.Cond

	// elems is the list of elements.
	subjects map[T]bool

	// mu is the mutex to synchronize map access.
	mu sync.RWMutex
}

// NewLocker creates a new Locker.
//
// Use Locker.Set for observer boolean predicates.
//
// Parameters:
//   - keys: The keys to initialize the locker.
//
// Returns:
//   - *Locker[T]: A new Locker.
//
// Behaviors:
//   - All the predicates are initialized to true.
func NewLocker[T uc.Enumer](keys ...T) *Locker[T] {
	l := &Locker[T]{
		cond:     sync.NewCond(&sync.Mutex{}),
		subjects: make(map[T]bool),
	}

	for _, key := range keys {
		l.subjects[key] = true
	}

	return l
}

// DoFunc is a function that executes a function while waiting for the condition to be false.
//
// Parameters:
//   - l: The SafeMap to use.
//
// Returns:
//   - bool: True if the function should exit, false otherwise.
type DoFunc[T uc.Enumer] func(l map[T]bool) bool

// Do executes a function while waiting for at least one of the conditions to be false.
//
// Parameters:
//   - f: The function to execute.
//
// Returns:
//   - bool: True if the function should exit, false otherwise.
func (l *Locker[T]) Do(f DoFunc[T]) bool {
	l.cond.L.Lock()
	defer l.cond.L.Unlock()

	for !l.hasFalse() {
		l.cond.Wait()
	}

	shouldExit := f(l.subjects)

	return shouldExit
}

func (l *Locker[T]) hasFalse() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	for _, value := range l.subjects {
		if !value {
			return true
		}
	}

	return false
}

// DoUntill executes a function while waiting for the condition to be false.
//
// The function will be executed until the condition returned by the function is true.
//
// Parameters:
//   - f: The function to execute.
func (l *Locker[T]) DoUntill(f DoFunc[T]) {
	for {
		shouldExit := l.Do(f)
		if shouldExit {
			break
		}
	}
}

// Broadcast broadcasts the condition to all waiting goroutines.
//
// Parameters:
//   - key: The key to broadcast.
//   - value: The value to broadcast.
func (l *Locker[T]) Broadcast(key T, value bool) {
	l.cond.L.Lock()
	defer l.cond.L.Unlock()

	l.mu.Lock()
	defer l.mu.Unlock()

	l.subjects[key] = value

	l.cond.Broadcast()
}

// Signal signals the condition to a single waiting goroutine.
//
// Parameters:
//   - key: The key to signal.
//   - value: The value to signal.
func (l *Locker[T]) Signal(key T, value bool) {
	l.cond.L.Lock()
	defer l.cond.L.Unlock()

	l.mu.Lock()
	defer l.mu.Unlock()

	l.subjects[key] = value

	l.cond.Signal()
}

// Get returns the value of a predicate.
//
// Parameters:
//   - key: The key to get the value.
//
// Returns:
//   - bool: The value of the predicate or false if the key does not exist.
func (l *Locker[T]) Get(key T) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	val, ok := l.subjects[key]
	if !ok {
		return false
	}

	return val
}
