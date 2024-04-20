// Package errors provides a custom error type for out-of-bound errors.
package Errors

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// ErrNoError represents an error when no error occurs.
type ErrNoError struct {
	// reason is the reason for the no error error.
	reason error
}

// Error generates the error message for the ErrNoError error, that is, the reason
// for the no error error.
//
// Returns:
//
//   - string: The error message.
func (e *ErrNoError) Error() string {
	return fmt.Sprintf("%v", e.reason)
}

// Unwrap returns the reason for the no error error.
// It is used for error unwrapping.
//
// Returns:
//
//   - error: The reason for the no error error.
func (e *ErrNoError) Unwrap() error {
	return e.reason
}

// NewErrNoError creates a new ErrNoError error.
//
// Parameters:
//
//   - reason: The reason for the no error error.
//
// Returns:
//
//   - *ErrNoError: A pointer to the newly created ErrNoError.
func NewErrNoError(reason error) *ErrNoError {
	return &ErrNoError{reason: reason}
}

// Wrap sets the reason for the no error error.
//
// Parameters:
//
//   - reason: The reason for the no error error.
//
// Returns:
//
//   - *ErrNoError: The error instance for chaining.
func (e *ErrNoError) Wrap(reason error) *ErrNoError {
	e.reason = reason

	return e
}

// IsNoError checks if an error is a no error error or if it is nil.
//
// Parameters:
//
//   - err: The error to check.
//
// Returns:
//
//   - bool: True if the error is a no error error or if it is nil, otherwise false.
func IsNoError(err error) bool {
	if err == nil {
		return true
	}

	var errNoError *ErrNoError

	return errors.As(err, &errNoError)
}

// ErrPanic represents an error when a panic occurs.
type ErrPanic struct {
	// value is the value that caused the panic.
	value any
}

// Error generates the error message for the ErrPanic error, including the value
// that caused the panic.
//
// Returns:
//
//   - string: The error message.
func (e *ErrPanic) Error() string {
	return fmt.Sprintf("panic: %v", e.value)
}

// NewErrPanic creates a new ErrPanic error.
//
// Parameters:
//
//   - value: The value that caused the panic.
//
// Returns:
//
//   - error: A pointer to the newly created ErrPanic.
func NewErrPanic(value any) *ErrPanic {
	return &ErrPanic{value: value}
}

// ErrOutOfBound represents an error when a value is out of a specified range.
type ErrOutOfBound struct {
	// lowerBound and upperBound are the lower and upper bounds of the range,
	// respectively.
	lowerBound, upperBound int

	// lowerInclusive and upperInclusive are flags indicating whether the lower
	// and upper bounds are inclusive, respectively.
	lowerInclusive, upperInclusive bool

	// value is the value that caused the error.
	value int
}

// Error generates the error message for the ErrOutOfBound error, including the
// value, the range, and whether the bounds are inclusive.
// The message includes the range, and whether the bounds are inclusive.
//
// Returns:
//
//   - string: The error message.
func (e *ErrOutOfBound) Error() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "value (%v) not in range ", e.value)
	if e.lowerInclusive {
		builder.WriteRune('[')
	} else {
		builder.WriteRune('(')
	}

	fmt.Fprintf(&builder, "%d, %d", e.lowerBound, e.upperBound)
	if e.upperInclusive {
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
func (e *ErrOutOfBound) WithLowerBound(isInclusive bool) *ErrOutOfBound {
	e.lowerInclusive = isInclusive

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
func (e *ErrOutOfBound) WithUpperBound(isInclusive bool) *ErrOutOfBound {
	e.upperInclusive = isInclusive

	return e
}

// NewOutOfBound creates a new ErrOutOfBound error. If no inclusivity flags are
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
//   - error: A pointer to the newly created ErrOutOfBound.
func NewErrOutOfBound(value int, lowerBound, upperBound int) *ErrOutOfBound {
	return &ErrOutOfBound{
		lowerBound:     lowerBound,
		upperBound:     upperBound,
		lowerInclusive: true,
		upperInclusive: false,
		value:          value,
	}
}

// ErrInvalidParameter represents an error when a parameter is invalid.
type ErrInvalidParameter struct {
	// parameter is the name of the parameter.
	parameter string

	// reason is the reason for the invalidity of the parameter.
	reason error
}

// Error generates the error message for the ErrInvalidParameter error, including
// the parameter name and the reason for its invalidity.
//
// Returns:
//
//   - string: The error message.
func (e *ErrInvalidParameter) Error() string {
	return fmt.Sprintf("parameter (%s) is invalid: %v", e.parameter, e.reason)
}

// Unwrap returns the reason for the invalidity of the parameter.
// It is used for error unwrapping.
//
// Returns:
//
//   - error: The reason for the invalidity of the parameter.
func (e *ErrInvalidParameter) Unwrap() error {
	return e.reason
}

// NewErrInvalidParameter creates a new ErrInvalidParameter error.
// If the reason is not provided (nil), the reason is set to
// "parameter is invalid" by default.
//
// Parameters:
//
//   - parameter: The name of the parameter.
//   - reason: The reason for the invalidity.
//
// Returns:
//
//   - error: A pointer to the newly created ErrInvalidParameter.
func NewErrInvalidParameter(parameter string) *ErrInvalidParameter {
	return &ErrInvalidParameter{
		parameter: parameter,
		reason:    errors.New("parameter is invalid"),
	}
}

// Wrap sets the reason for the invalidity of the parameter.
// If the reason is not provided (nil), the reason is set to "parameter is invalid"
// by default.
func (e *ErrInvalidParameter) Wrap(reason error) *ErrInvalidParameter {
	e.reason = reason

	return e
}

// ErrCallFailed represents an error that occurs when a function call fails.
type ErrCallFailed struct {
	// fnName is the name of the function.
	fnName string

	// signature is the signature of the function.
	signature reflect.Type

	// reason is the reason for the failure.
	reason error
}

// Error generates a string representation of an ErrCallFailed, including the
// function name, signature, and the reason for the failure.
//
// Returns:
//
//   - string: The error message.
func (e *ErrCallFailed) Error() string {
	return fmt.Sprintf("call to %s%v failed: %v",
		e.fnName, e.signature, e.reason)
}

// Unwrap returns the underlying error that caused the ErrCallFailed.
// It is used for error unwrapping.
//
// Returns:
//
//   - error: The reason for the failure.
func (e *ErrCallFailed) Unwrap() error {
	return e.reason
}

// NewCallFailed creates a new ErrCallFailed. If the reason is not provided (nil),
// the reason is set to "an error occurred while calling the function" by default.
//
// Parameters:
//
//   - functionName: The name of the function.
//   - function: The function that failed.
//   - reason: The reason for the failure.
//
// Returns:
//
//   - error: A pointer to the new ErrCallFailed.
func NewErrCallFailed(functionName string, function any) *ErrCallFailed {
	return &ErrCallFailed{
		fnName:    functionName,
		signature: reflect.ValueOf(function).Type(),
		reason:    errors.New("an error occurred while calling the function"),
	}
}

// Wrap sets the reason for the failure.
// If the reason is not provided (nil), the reason is set to "an error occurred
// while calling the function" by default.
//
// Parameters:
//
//   - reason: The reason for the failure.
//
// Returns:
//
//   - *ErrCallFailed: The error instance for chaining.
func (e *ErrCallFailed) Wrap(reason error) *ErrCallFailed {
	e.reason = reason

	return e
}

// ErrUnexpected represents an error that occurs when an unexpected value is
// encountered.
type ErrUnexpected struct {
	// expected is the list of expected values.
	expected []string

	// actual is the actual value encountered.
	actual fmt.Stringer
}

// Error generates a string representation of an ErrUnexpected, including the
// expected values and the actual value encountered.
func (e *ErrUnexpected) Error() string {
	var expected, got string

	switch len(e.expected) {
	case 0:
		expected = "nothing"
	case 1:
		expected = fmt.Sprintf("%q", e.expected[0])
	case 2:
		expected = fmt.Sprintf("%q or %q", e.expected[0], e.expected[1])
	default:
		var builder strings.Builder

		fmt.Fprintf(&builder, "%q", e.expected[0])

		for i := 1; i < len(e.expected)-1; i++ {
			fmt.Fprintf(&builder, ", %q", e.expected[i])
		}

		fmt.Fprintf(&builder, ", or %q", e.expected[len(e.expected)-1])

		expected = builder.String()
	}

	if e.actual == nil {
		got = "nothing"
	} else {
		got = fmt.Sprintf("%v", e.actual)
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
	return &ErrUnexpected{expected: expected, actual: got}
}

// ErrNilParameter represents an error when a parameter is nil.
// This is a shorthand for NewErrInvalidParameter(parameter).Wrap(errors.New("value is nil")).
//
// Parameters:
//
//   - parameter: The name of the parameter.
//
// Returns:
//
//   - *ErrInvalidParameter: A pointer to the newly created ErrInvalidParameter.
func NewErrNilParameter(parameter string) *ErrInvalidParameter {
	return &ErrInvalidParameter{parameter: parameter, reason: errors.New("value is nil")}
}
