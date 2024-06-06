package RWSafe

import (
	"sync"

	uc "github.com/PlayerR9/MyGoLib/Units/Common"
)

// Locker is a thread-safe locker that allows multiple goroutines to wait for a condition.
type Locker[T uc.Enumer] struct {
	// preds is the map of predicates.
	preds *SafeMap[T, bool]

	// cond is the condition variable.
	cond *sync.Cond
}

// NewLocker creates a new Locker.
//
// Parameters:
//   - keys: The keys to initialize the locker.
//
// Returns:
//   - *Locker[T]: A new Locker.
func NewLocker[T uc.Enumer](keys ...T) *Locker[T] {
	l := &Locker[T]{
		preds: NewSafeMap[T, bool](),
		cond:  sync.NewCond(&sync.Mutex{}),
	}

	for _, key := range keys {
		l.preds.Set(key, false)
	}

	return l
}

// is checks if any of the predicates are true.
//
// Returns:
//   - bool: True if all predicates are true, false otherwise.
func (l *Locker[T]) is() bool {
	ok, err := l.preds.Scan(func(key T, value bool) (bool, error) {
		return value, nil
	})
	if err != nil {
		return false
	}

	return ok
}

// DoFunc is a function that executes a function while waiting for the condition to be false.
//
// Parameters:
//   - sm: The SafeMap to use.
//
// Returns:
//   - bool: A flag indicating whether the waiting should end or not.
type DoFunc[T uc.Enumer] func(sm *SafeMap[T, bool]) bool

// Do executes a function while waiting for the condition to be false.
//
// Parameters:
//   - f: The function to execute.
//
// Returns:
//   - bool: A flag indicating the result of the function.
func (l *Locker[T]) Do(f DoFunc[T]) bool {
	l.cond.L.Lock()
	defer l.cond.L.Unlock()

	for l.is() {
		l.cond.Wait()
	}

	ok := f(l.preds)

	return ok
}

// DoUntill executes a function while waiting for the condition to be false.
//
// The function will be executed until the condition returned by the function is true.
//
// Parameters:
//   - f: The function to execute.
func (l *Locker[T]) DoUntill(f DoFunc[T]) {
	ok := false

	for !ok {
		l.cond.L.Lock()

		for l.is() {
			l.cond.Wait()
		}

		ok = f(l.preds)

		l.cond.L.Unlock()
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

	l.preds.Set(key, value)

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

	l.preds.Set(key, value)

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
	val, ok := l.preds.Get(key)
	if !ok {
		return false
	}

	return val
}
