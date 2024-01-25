package Counters

import (
	"fmt"
	"strings"
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
func NewErrCannotAdvanceCounter(reason error) *ErrCannotAdvanceCounter {
	return &ErrCannotAdvanceCounter{
		counter: nil,
		reason:  reason,
	}
}

// WithCounter sets the counter that caused the error.
// It follows the builder pattern, returning the ErrCannotAdvanceCounter
// itself for chaining.
//
// The function takes the following parameter:
//
//   - c is the counter that caused the error.
//
// The function returns the following:
//
//   - A pointer to the ErrCannotAdvanceCounter itself.
func (e *ErrCannotAdvanceCounter) WithCounter(c Counter) *ErrCannotAdvanceCounter {
	e.counter = c
	return e
}

// Error returns a string representation of the ErrCannotAdvanceCounter.
// It includes the type of the counter and the reason for the error.
//
// The function returns the following:
//
//   - A string representing the ErrCannotAdvanceCounter.
func (e *ErrCannotAdvanceCounter) Error() string {
	var builder strings.Builder

	builder.WriteString("cannot advance ")

	if e.counter == nil {
		builder.WriteString("counter")
	} else {
		fmt.Fprintf(&builder, "%T", e.counter)
	}

	fmt.Fprintf(&builder, ": reason=%v", e.reason)

	return builder.String()
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
func NewErrCannotRetreatCounter(reason error) *ErrCannotRetreatCounter {
	return &ErrCannotRetreatCounter{
		counter: nil,
		reason:  reason,
	}
}

// WithCounter sets the counter that caused the error.
// It follows the builder pattern, returning the ErrCannotRetreatCounter
// itself for chaining.
//
// The function takes the following parameter:
//
//   - c is the counter that caused the error.
//
// The function returns the following:
//
//   - A pointer to the ErrCannotRetreatCounter itself.
func (e *ErrCannotRetreatCounter) WithCounter(c Counter) *ErrCannotRetreatCounter {
	e.counter = c
	return e
}

// Error returns a string representation of the ErrCannotRetreatCounter.
// It includes the type of the counter and the reason for the error.
//
// The function returns the following:
//
//   - A string representing the ErrCannotRetreatCounter.
func (e *ErrCannotRetreatCounter) Error() string {
	var builder strings.Builder

	builder.WriteString("cannot retreat ")

	if e.counter == nil {
		builder.WriteString("counter")
	} else {
		fmt.Fprintf(&builder, "%T", e.counter)
	}

	fmt.Fprintf(&builder, ": reason=%v", e.reason)

	return builder.String()
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
