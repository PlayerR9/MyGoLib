// Package errors provides a custom error type for out-of-bound errors.
package Errors

import (
	"fmt"
	"strings"
)

// ErrPanic represents an error when a panic occurs.
type ErrPanic struct {
	// Value is the value that caused the panic.
	Value any
}

// ErrorIf returns the error if the value is not nil.
//
// Returns:
//   - error: The error if the value is not nil, nil otherwise.
func (e *ErrPanic) ErrorIf() error {
	if e.Value != nil {
		return e
	}

	return nil
}

// Error is a method of the error interface.
//
// Returns:
//   - string: The error message of the panic error.
func (e *ErrPanic) Error() string {
	return fmt.Sprintf("panic: %v", e.Value)
}

// NewErrPanic creates a new ErrPanic error.
//
// Parameters:
//   - value: The value that caused the panic.
//
// Returns:
//   - *ErrPanic: A pointer to the newly created ErrPanic.
func NewErrPanic(value any) *ErrPanic {
	return &ErrPanic{Value: value}
}

// ErrOutOfBounds represents an error when a value is out of a specified range.
type ErrOutOfBounds struct {
	// LowerBound and UpperBound are the lower and upper bounds of the range,
	// respectively.
	LowerBound, UpperBound int

	// LowerInclusive and UpperInclusive are flags indicating whether the lower
	// and upper bounds are inclusive, respectively.
	LowerInclusive, UpperInclusive bool

	// Value is the value that caused the error.
	Value int
}

// ErrorIf returns the error if the value is out of the specified range.
//
// Returns:
//   - error: The error if the value is out of the specified range, nil otherwise.
func (e *ErrOutOfBounds) ErrorIf() error {
	if e.Value < e.LowerBound || (e.Value == e.LowerBound && !e.LowerInclusive) {
		return e
	}

	if e.Value > e.UpperBound || (!e.UpperInclusive && e.Value == e.UpperBound) {
		return e
	}

	return nil
}

// Error is a method of the error interface.
//
// The format of the error message is as follows:
//
//	value (value) not in range [lowerBound, upperBound]
//
// The square brackets indicate that the lower bound is inclusive, while the
// parentheses indicate that the upper bound is exclusive.
//
// Returns:
//
//   - string: The error message of the out-of-bound error.
func (e *ErrOutOfBounds) Error() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "value (%d) not in range ", e.Value)
	if e.LowerInclusive {
		builder.WriteRune('[')
	} else {
		builder.WriteRune('(')
	}

	fmt.Fprintf(&builder, "%d, %d", e.LowerBound, e.UpperBound)
	if e.UpperInclusive {
		builder.WriteRune(']')
	} else {
		builder.WriteRune(')')
	}

	return builder.String()
}

// WithLowerBound sets the inclusivity of the lower bound.
//
// Parameters:
//
//   - isInclusive: A boolean indicating whether the lower bound is inclusive.
//
// Returns:
//
//   - *ErrOutOfBound: The error instance for chaining.
func (e *ErrOutOfBounds) WithLowerBound(isInclusive bool) *ErrOutOfBounds {
	e.LowerInclusive = isInclusive

	return e
}

// WithUpperBound sets the inclusivity of the upper bound.
//
// Parameters:
//
//   - isInclusive: A boolean indicating whether the upper bound is inclusive.
//
// Returns:
//
//   - *ErrOutOfBound: The error instance for chaining.
func (e *ErrOutOfBounds) WithUpperBound(isInclusive bool) *ErrOutOfBounds {
	e.UpperInclusive = isInclusive

	return e
}

// NewOutOfBounds creates a new ErrOutOfBound error. If no inclusivity flags are
// provided, the lower bound is inclusive and the upper bound is exclusive.
//
// Parameters:
//
//   - lowerBound, upperbound: The lower and upper bounds of the range,
//     respectively.
//   - value: The value that caused the error.
//
// Returns:
//
//   - *ErrOutOfBounds: A pointer to the newly created ErrOutOfBound.
func NewErrOutOfBounds(value int, lowerBound, upperBound int) *ErrOutOfBounds {
	return &ErrOutOfBounds{
		LowerBound:     lowerBound,
		UpperBound:     upperBound,
		LowerInclusive: true,
		UpperInclusive: false,
		Value:          value,
	}
}

// ErrUnexpected represents an error that occurs when an unexpected value is
// encountered.
type ErrUnexpected struct {
	// Expected is the list of expected values.
	Expected []string

	// Actual is the actual value encountered.
	Actual fmt.Stringer
}

// ErrorIf returns the error if the actual value is not one of the expected
// values.
//
// Returns:
//   - error: The error if the actual value is not one of the expected values,
//     nil otherwise.
func (e *ErrUnexpected) ErrorIf() error {
	if len(e.Expected) == 0 {
		return nil
	} else if e.Actual == nil {
		return e
	}

	actual := e.Actual.String()

	for _, expected := range e.Expected {
		if actual == expected {
			return nil
		}
	}

	return e
}

// Error is a method of the error interface.
//
// Returns:
//
//   - string: The error message.
func (e *ErrUnexpected) Error() string {
	var expected, got string

	switch len(e.Expected) {
	case 0:
		expected = "nothing"
	case 1:
		expected = fmt.Sprintf("%q", e.Expected[0])
	case 2:
		expected = fmt.Sprintf("%q or %q", e.Expected[0], e.Expected[1])
	default:
		var builder strings.Builder

		fmt.Fprintf(&builder, "%q", e.Expected[0])

		for i := 1; i < len(e.Expected)-1; i++ {
			builder.WriteRune(',')
			builder.WriteRune(' ')
			fmt.Fprintf(&builder, "%q", e.Expected[i])
		}

		builder.WriteRune(',')
		builder.WriteRune(' ')
		builder.WriteString("or ")

		fmt.Fprintf(&builder, "%q", e.Expected[len(e.Expected)-1])

		expected = builder.String()
	}

	if e.Actual == nil {
		got = "nothing"
	} else {
		got = e.Actual.String()
	}

	return fmt.Sprintf("expected %s, got %s instead", expected, got)
}

// NewErrUnexpected creates a new ErrUnexpected error.
//
// Parameters:
//
//   - got: The actual value encountered.
//   - expected: The list of expected values.
//
// Returns:
//
//   - *ErrUnexpected: A pointer to the newly created ErrUnexpected.
func NewErrUnexpected(got fmt.Stringer, expected ...string) *ErrUnexpected {
	return &ErrUnexpected{Expected: expected, Actual: got}
}
