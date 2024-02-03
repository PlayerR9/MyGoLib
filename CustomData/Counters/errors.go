package Counters

import (
	"fmt"
)

// ErrCannotAdvanceCounter represents an error that occurs when a
// counter cannot be advanced.
// It includes the counter that caused the error and the reason
// for the error.
type ErrCannotAdvanceCounter struct {
	// counter is the counter that caused the error.
	counter Counter

	// reason is the specific error that occurred.
	reason error
}

// NewErrCannotAdvanceCounter creates a new ErrCannotAdvanceCounter.
// It takes the following parameter:
//
//   - reason is the specific error that occurred.
//
// The function returns the following:
//
//   - A pointer to the new ErrCannotAdvanceCounter.
func NewErrCannotAdvanceCounter(counter Counter, reason error) *ErrCannotAdvanceCounter {
	return &ErrCannotAdvanceCounter{
		counter: counter,
		reason:  reason,
	}
}

// Error returns a string representation of the ErrCannotAdvanceCounter.
// It includes the type of the counter and the reason for the error.
//
// The function returns the following:
//
//   - A string representing the ErrCannotAdvanceCounter.
func (e *ErrCannotAdvanceCounter) Error() string {
	return fmt.Sprintf("cannot advance %T: reason=%v", e.counter, e.reason)
}

// Unwrap returns the specific error that occurred, which is wrapped in
// the ErrCannotAdvanceCounter.
// This allows the use of errors.Is and errors.As to check for specific
// errors.
//
// The function returns the following:
//
//   - The specific error that occurred.
func (e *ErrCannotAdvanceCounter) Unwrap() error {
	return e.reason
}

// ErrCannotRetreatCounter represents an error that occurs when a
// counter cannot be retreated.
// It includes the counter that caused the error and the reason
// for the error.
type ErrCannotRetreatCounter struct {
	// counter is the counter that caused the error.
	counter Counter

	// reason is the specific error that occurred.
	reason error
}

// NewErrCannotRetreatCounter creates a new ErrCannotRetreatCounter.
// It takes the following parameter:
//
//   - reason is the specific error that occurred.
//
// The function returns the following:
//
//   - A pointer to the new ErrCannotRetreatCounter.
func NewErrCannotRetreatCounter(counter Counter, reason error) *ErrCannotRetreatCounter {
	return &ErrCannotRetreatCounter{
		counter: counter,
		reason:  reason,
	}
}

// Error returns a string representation of the ErrCannotRetreatCounter.
// It includes the type of the counter and the reason for the error.
//
// The function returns the following:
//
//   - A string representing the ErrCannotRetreatCounter.
func (e *ErrCannotRetreatCounter) Error() string {
	return fmt.Sprintf("cannot retreat %T: reason=%v", e.counter, e.reason)
}

// Unwrap returns the specific error that occurred, which is wrapped
// in the ErrCannotRetreatCounter.
// This allows the use of errors.Is and errors.As to check for specific
// errors.
//
// The function returns the following:
//
//   - The specific error that occurred.
func (e *ErrCannotRetreatCounter) Unwrap() error {
	return e.reason
}
