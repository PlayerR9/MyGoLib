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

// ErrGT represents an error when a value is less than or equal to a specified value.
type ErrGT struct {
	// Value is the value that caused the error.
	Value int
}

// Error is a method of the error interface that returns the error message.
//
// Returns:
//
//   - string: The error message.
func (e *ErrGT) Error() string {
	if e.Value == 0 {
		return "value must be positive"
	} else {
		return fmt.Sprintf("value must be greater than %d", e.Value)
	}
}

// NewErrGT creates a new ErrGT error.
//
// Parameters:
//
//   - value: The value that caused the error.
//
// Returns:
//
//   - *ErrGT: A pointer to the newly created ErrGT.
func NewErrGT(value int) *ErrGT {
	return &ErrGT{Value: value}
}

// ErrLT represents an error when a value is greater than or equal to a specified value.
type ErrLT struct {
	// Value is the value that caused the error.
	Value int
}

// Error is a method of the error interface that returns the error message.
//
// Returns:
//
//   - string: The error message.
func (e *ErrLT) Error() string {
	if e.Value == 0 {
		return "value must be negative"
	} else {
		return fmt.Sprintf("value must be less than %d", e.Value)
	}
}

// NewErrLT creates a new ErrLT error.
//
// Parameters:
//
//   - value: The value that caused the error.
//
// Returns:
//
//   - *ErrLT: A pointer to the newly created ErrLT.
func NewErrLT(value int) *ErrLT {
	return &ErrLT{Value: value}
}

// ErrGTE represents an error when a value is less than a specified value.
type ErrGTE struct {
	// Value is the value that caused the error.
	Value int
}

// Error is a method of the error interface that returns the error message.
//
// Returns:
//
//   - string: The error message.
func (e *ErrGTE) Error() string {
	if e.Value == 0 {
		return "value must be non-negative"
	} else {
		return fmt.Sprintf("value must be greater than %d", e.Value)
	}
}

// NewErrGTE creates a new ErrGTE error.
//
// Parameters:
//
//   - value: The value that caused the error.
//
// Returns:
//
//   - *ErrGTE: A pointer to the newly created ErrGTE.
func NewErrGTE(value int) *ErrGTE {
	return &ErrGTE{Value: value}
}

// ErrLTE represents an error when a value is greater than a specified value.
type ErrLTE struct {
	// Value is the value that caused the error.
	Value int
}

// Error is a method of the error interface that returns the error message.
//
// Returns:
//
//   - string: The error message.
func (e *ErrLTE) Error() string {
	if e.Value == 0 {
		return "value must be non-positive"
	} else {
		return fmt.Sprintf("value must be less than or equal to %d", e.Value)
	}
}

// NewErrLTE creates a new ErrLTE error.
//
// Parameters:
//
//   - value: The value that caused the error.
//
// Returns:
//
//   - *ErrLTE: A pointer to the newly created ErrLTE.
func NewErrLTE(value int) *ErrLTE {
	return &ErrLTE{Value: value}
}

// ErrInvalidValue represents an error when a value is invalid.
type ErrInvalidValue struct {
	// Values is the list of invalid values.
	Values []int
}

// Error is a method of the error interface that returns the error message.
//
// Returns:
//
//   - string: The error message.
func (e *ErrInvalidValue) Error() string {
	if len(e.Values) == 0 {
		return "value is invalid"
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "value must not be %d", e.Values[0])

	if len(e.Values) > 1 {
		for i := 1; i < len(e.Values)-1; i++ {
			builder.WriteRune(',')
			builder.WriteRune(' ')
			fmt.Fprintf(&builder, "%d", e.Values[i])
		}

		builder.WriteRune(',')
		builder.WriteRune(' ')
		builder.WriteString("or ")
		fmt.Fprintf(&builder, "%d", e.Values[len(e.Values)-1])
	}

	return builder.String()
}

// NewErrInvalidValue creates a new ErrInvalidValue error.
//
// Parameters:
//
//   - values: The list of invalid values.
//
// Returns:
//
//   - *ErrInvalidValue: A pointer to the newly created ErrInvalidValue.
func NewErrInvalidValue(values ...int) *ErrInvalidValue {
	return &ErrInvalidValue{Values: values}
}

// ErrEmptyString represents an error when a string is empty.
type ErrEmptyString struct{}

// Error is a method of the error interface that returns the error message.
//
// Returns:
//
//   - string: The error message.
func (e *ErrEmptyString) Error() string {
	return "value must not be empty"
}

// NewErrEmptyString creates a new ErrEmptyString error.
//
// Returns:
//
//   - *ErrEmptyString: A pointer to the newly created ErrEmptyString.
func NewErrEmptyString() *ErrEmptyString {
	return &ErrEmptyString{}
}

// ErrEmptySlice represents an error when a slice is empty.
type ErrEmptySlice struct{}

// Error is a method of the error interface that returns the error message.
//
// Returns:
//   - string: The error message.
func (e *ErrEmptySlice) Error() string {
	return "slice must not be empty"
}

// NewErrEmptySlice creates a new ErrEmptySlice error.
//
// Returns:
//   - *ErrEmptySlice: A pointer to the newly created ErrEmptySlice.
func NewErrEmptySlice() *ErrEmptySlice {
	return &ErrEmptySlice{}
}
