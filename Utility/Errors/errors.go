// Package errors provides a custom error type for out-of-bound errors.
package Errors

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type PropagableError interface {
	error
	Unwrap() error

	// Wrap sets the reason for the error.
	Wrap(error) error
}

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
func (e *ErrNoError) Wrap(reason error) error {
	return &ErrNoError{reason: reason}
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

// WithReason sets the reason for the invalidity of the parameter.
// If the reason is not provided (nil), the reason is set to "parameter is invalid"
// by default.
func (e *ErrInvalidParameter) Wrap(reason error) error {
	return &ErrInvalidParameter{
		parameter: e.parameter,
		reason:    reason,
	}
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

// WithReason sets the reason for the failure.
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
func (e *ErrCallFailed) Wrap(reason error) error {
	return &ErrCallFailed{
		fnName:    e.fnName,
		signature: e.signature,
		reason:    reason,
	}
}
