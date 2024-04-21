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
	// Err is the reason for the no error error.
	Err error
}

// Error is a method of the error interface.
//
// Returns:
//
//   - string: The error message of the no error error
//     (no mention of being a no error error).
func (e *ErrNoError) Error() string {
	return e.Err.Error()
}

// Unwrap is a method of the errors interface.
//
// Returns:
//
//   - error: The reason for the no error error.
func (e *ErrNoError) Unwrap() error {
	return e.Err
}

// NewErrNoError creates a new ErrNoError error.
//
// Parameters:
//
//   - err: The reason for the no error error.
//
// Returns:
//
//   - *ErrNoError: A pointer to the newly created ErrNoError.
func NewErrNoError(err error) *ErrNoError {
	return &ErrNoError{Err: err}
}

// ErrPanic represents an error when a panic occurs.
type ErrPanic struct {
	// Value is the value that caused the panic.
	Value any
}

// Error is a method of the error interface.
//
// Returns:
//
//   - string: The error message of the panic error.
func (e *ErrPanic) Error() string {
	return fmt.Sprintf("panic: %v", e.Value)
}

// NewErrPanic creates a new ErrPanic error.
//
// Parameters:
//
//   - value: The value that caused the panic.
//
// Returns:
//
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

	fmt.Fprintf(&builder, "value (%v) not in range ", e.Value)
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

// ErrInvalidParameter represents an error when a parameter is invalid.
type ErrInvalidParameter struct {
	// Parameter is the name of the Parameter.
	Parameter string

	// Reason is the Reason for the invalidity of the parameter.
	Reason error
}

// Error is a method of the error interface.
//
// If the reason is not provided (nil), no reason is included in the error message.
//
// Returns:
//
//   - string: The error message.
func (e *ErrInvalidParameter) Error() string {
	if e.Reason == nil {
		return fmt.Sprintf("parameter (%s) is invalid", e.Parameter)
	}

	return fmt.Sprintf("parameter (%s) is invalid: %v", e.Parameter, e.Reason)
}

// Unwrap returns the reason for the invalidity of the parameter.
// It is used for error unwrapping.
//
// Returns:
//
//   - error: The reason for the invalidity of the parameter.
func (e *ErrInvalidParameter) Unwrap() error {
	return e.Reason
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
func NewErrInvalidParameter(parameter string, reason error) *ErrInvalidParameter {
	return &ErrInvalidParameter{
		Parameter: parameter,
		Reason:    reason,
	}
}

// ErrInvalidCall represents an error that occurs when a function
// is not called correctly.
type ErrInvalidCall struct {
	// FnName is the name of the function.
	FnName string

	// Signature is the Signature of the function.
	Signature reflect.Type

	// Reason is the Reason for the failure.
	Reason error
}

// Error is a method of the error interface.
//
// Returns:
//
//   - string: The error message.
func (e *ErrInvalidCall) Error() string {
	if e.Reason == nil {
		return fmt.Sprintf("call to %s%v failed", e.FnName, e.Signature)
	}

	return fmt.Sprintf("call to %s%v failed: %v",
		e.FnName, e.Signature, e.Reason)
}

// Unwrap returns the underlying error that caused the ErrInvalidCall.
// It is used for error unwrapping.
//
// Returns:
//
//   - error: The reason for the failure.
func (e *ErrInvalidCall) Unwrap() error {
	return e.Reason
}

// NewErrInvalidCall creates a new ErrInvalidCall. If the reason is not provided (nil),
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
//   - *ErrInvalidCall: A pointer to the new ErrInvalidCall.
func NewErrInvalidCall(functionName string, function any, reason error) *ErrInvalidCall {
	return &ErrInvalidCall{
		FnName:    functionName,
		Signature: reflect.ValueOf(function).Type(),
		Reason:    reason,
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
			fmt.Fprintf(&builder, ", %q", e.Expected[i])
		}

		fmt.Fprintf(&builder, ", or %q", e.Expected[len(e.Expected)-1])

		expected = builder.String()
	}

	if e.Actual == nil {
		got = "nothing"
	} else {
		got = fmt.Sprintf("%v", e.Actual)
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

// ErrNilParameter represents an error when a parameter is nil.
// This is a shorthand for NewErrInvalidParameter(parameter, errors.New("value is nil")).
//
// Parameters:
//
//   - parameter: The name of the parameter.
//
// Returns:
//
//   - *ErrInvalidParameter: A pointer to the newly created ErrInvalidParameter.
func NewErrNilParameter(parameter string) *ErrInvalidParameter {
	return &ErrInvalidParameter{Parameter: parameter, Reason: errors.New("value is nil")}
}

// ErrIgnorable represents an error that can be ignored. Useful for indicating
// that an error is ignorable.
type ErrIgnorable struct {
	// Err is the error that can be ignored.
	Err error
}

// Error is a method of the error interface.
//
// Returns:
//
//   - string: The error message of the ignorable error (no mention of being ignorable).
func (e ErrIgnorable) Error() string {
	return e.Err.Error()
}

// Unwrap returns the error that can be ignored.
// It is used for error unwrapping.
//
// Returns:
//
//   - error: The error that can be ignored.
func (e ErrIgnorable) Unwrap() error {
	return e.Err
}

// NewErrIgnorable creates a new ErrIgnorable error.
//
// Parameters:
//
//   - err: The error that can be ignored.
//
// Returns:
//
//   - ErrIgnorable: The newly created ErrIgnorable.
func NewErrIgnorable(err error) ErrIgnorable {
	return ErrIgnorable{
		Err: err,
	}
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

// IsErrIgnorable checks if an error is an ErrIgnorable or ErrInvalidParameter error.
//
// Parameters:
//
//   - err: The error to check.
//
// Returns:
//
//   - bool: True if the error is an ErrIgnorable or ErrInvalidParameter error,
//     otherwise false.
func IsErrIgnorable(err error) bool {
	var ignorable *ErrIgnorable

	if errors.As(err, &ignorable) {
		return true
	}

	var invalid *ErrInvalidParameter

	return errors.As(err, &invalid)
}
