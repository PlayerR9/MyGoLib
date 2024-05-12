package RWSafe

import (
	"sync"

	uc "github.com/PlayerR9/MyGoLib/Units/Common"
)

// Locker is a struct that represents a locker.
type Locker struct {
	// cond is the condition of the locker.
	cond bool

	// signal is the signal of the locker.
	signal sync.Cond
}

// NewLocker creates a new Locker with the given condition.
//
// Parameters:
//   - cond: The condition of the locker.
//
// Returns:
//   - *Locker: A pointer to the new Locker.
func NewLocker(cond bool) *Locker {
	return &Locker{
		cond:   cond,
		signal: sync.Cond{L: new(sync.Mutex)},
	}
}

// ModifyCond modifies the condition of the locker.
//
// Parameters:
//   - cond: The new condition of the locker.
func (l *Locker) ModifyCond(cond bool) {
	l.signal.L.Lock()
	defer l.signal.L.Unlock()

	l.cond = cond
	l.signal.Broadcast()
}

// WaitForCond waits for the condition of the locker to be true.
//
// Parameters:
//   - do: The function to execute when the condition is true.
func (l *Locker) WaitForCond(do uc.RoutineFunc) {
	l.signal.L.Lock()
	defer l.signal.L.Unlock()

	for !l.cond {
		l.signal.Wait()
	}

	// Make use of the condition variable.
	do()
}
