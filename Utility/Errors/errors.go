// Package errors provides a custom error type for out-of-bound errors.
package Errors

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// ErrOutOfBound represents an error when a value is out of a specified range.
type ErrOutOfBound struct {
	// lowerBound and upperBound are the lower and upper bounds of the range,
	// respectively.
	lowerBound, upperBound int

	// value is the value that caused the error.
	value int

	// lowerInclusive and upperInclusive are flags indicating whether the lower and upper
	// bounds are inclusive, respectively.
	lowerInclusive, upperInclusive bool
}

// NewErrOutOfBound creates a new ErrOutOfBound error.
//
// Parameters:
//
//   - lowerBound is the lower bound of the range.
//   - upperBound is the upper bound of the range.
//   - value is the value that caused the error.
//
// Returns:
//
//   - A pointer to the new ErrOutOfBound.
func NewErrOutOfBound(lowerBound, upperBound, value int) *ErrOutOfBound {
	return &ErrOutOfBound{
		lowerBound:     lowerBound,
		lowerInclusive: false,
		upperBound:     upperBound,
		upperInclusive: true,
		value:          value,
	}
}

// WithLowerBound sets the inclusivity of the lower bound.
//
// Parameters:
//
//   - isInclusive is a boolean indicating whether the lower bound is inclusive.
//
// Returns:
//
//   - The ErrOutOfBound instance for chaining.
func (e *ErrOutOfBound) WithLowerBound(isInclusive bool) *ErrOutOfBound {
	e.lowerInclusive = isInclusive
	return e
}

// WithUpperBound sets the inclusivity of the upper bound.
//
// Parameters:
//
//   - isInclusive is a boolean indicating whether the upper bound is inclusive.
//
// Returns:
//
//   - The ErrOutOfBound instance for chaining.
func (e *ErrOutOfBound) WithUpperBound(isInclusive bool) *ErrOutOfBound {
	e.upperInclusive = isInclusive
	return e
}

// Error generates the error message for the ErrOutOfBound error.
// The message includes the value, the range, and whether the bounds are inclusive.
func (e *ErrOutOfBound) Error() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "value (%d) not in range ", e.value)
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

// ErrInvalidParameter represents an error when a parameter is invalid.
type ErrInvalidParameter struct {
	// parameter is the name of the parameter.
	parameter string

	// reason is the reason for the invalidity of the parameter.
	reason error
}

// NewErrInvalidParameter creates a new ErrInvalidParameter error.
//
// Parameters:
//
//   - parameter is the name of the parameter.
//   - reason is the reason for the invalidity.
//
// Returns:
//
//   - A pointer to the new ErrInvalidParameter.
//
// If the reason is not provided (nil), the reason is set to
// "parameter is invalid" by default.
func NewErrInvalidParameter(parameter string, reason error) *ErrInvalidParameter {
	if reason == nil {
		reason = errors.New("parameter is invalid")
	}

	return &ErrInvalidParameter{
		parameter: parameter,
		reason:    reason,
	}
}

// Error generates the error message for the ErrInvalidParameter error.
// The message includes the parameter name and the reason for its invalidity.
func (e *ErrInvalidParameter) Error() string {
	return fmt.Sprintf("invalid parameter: parameter=%s, reason=%v", e.parameter, e.reason)
}

// Unwrap returns the reason for the invalidity of the parameter.
// It is used for error unwrapping.
func (e *ErrInvalidParameter) Unwrap() error {
	return e.reason
}

// NewErrNilParameter creates a new ErrInvalidParameter error with the reason
// "parameter cannot be nil".
//
// Parameters:
//
//   - parameter is the name of the parameter.
//
// Returns:
//
//   - A pointer to the new ErrInvalidParameter.
func NewErrNilParameter(parameter string) *ErrInvalidParameter {
	return &ErrInvalidParameter{
		parameter: parameter,
		reason:    errors.New("parameter cannot be nil"),
	}
}

// ErrCallFailed represents an error that occurs when a function call fails.
type ErrCallFailed struct {
	// signature is the signature of the function.
	signature string

	// reason is the reason for the failure.
	reason error
}

// NewErrCallFailed creates a new ErrCallFailed. It generates the function signature
// from the provided function name and function, and sets the reason for the failure.
//
// Parameters:
//
//   - functionName is the name of the function.
//   - function is the function that failed.
//   - reason is the reason for the failure.
//
// Returns:
//
//   - A pointer to the new ErrCallFailed.
//
// If the reason is not provided (nil), the reason is set to
// "an error occurred while calling the function" by default.
func NewErrCallFailed(functionName string, function any, reason error) *ErrCallFailed {
	if reason == nil {
		reason = errors.New("an error occurred while calling the function")
	}

	return &ErrCallFailed{
		signature: fmt.Sprintf("%s%v", functionName, reflect.ValueOf(function).Type()),
		reason:    reason,
	}
}

// Error generates a string representation of an ErrCallFailed.
// It includes the function signature and the reason for the failure.
func (e *ErrCallFailed) Error() string {
	return fmt.Sprintf("call failed: function=%s, reason=%v", e.signature, e.reason)
}

// Unwrap returns the underlying error that caused the ErrCallFailed.
// Returns: The underlying error.
func (e *ErrCallFailed) Unwrap() error {
	return e.reason
}

// RecoverFromPanic is a function that recovers from a panic and sets the provided
// error pointer to the recovered value.
// If the recovered value is an error, it sets the error pointer to that error.
// Otherwise, it creates a new error with the recovered value and sets the error
// pointer to that.
//
// The function should be called using the defer statement to recover from panics.
//
// Parameters:
//
//   - err: The error pointer to set.
//
// Example:
//
//	  func MyFunction(n int) (result int, err error) {
//	      defer RecoverFromPanic(&err)
//
//	      // ...
//
//	      panic("something went wrong")
//		}
func RecoverFromPanic(err *error) {
	r := recover()
	if r == nil {
		return
	}

	if recErr, ok := r.(error); ok {
		*err = recErr
	} else {
		*err = fmt.Errorf("panic: %v", r)
	}
}

// Check is a generic function that validates if the provided value is the
// zero value of its type. If it is, the function triggers a panic with the
// supplied error.
//
// This function is designed to be used in conjunction with RecoverFromPanic.
// This combination allows a function to handle errors without the need for
// explicit error checks at every step.
//
// Parameters:
//
//   - val: The value to validate.
//   - err: The error to trigger a panic with if the value is the zero value of its type.
//
// Example:
//
//	  func MyFunction(n int) (result int, err error) {
//	      defer RecoverFromPanic(&err)
//
//	      Check(n <= 0, NewErrInvalidParameter("n", fmt.Errorf("value (%d) must be positive", n)))
//
//	      return 60 / n, nil
//		}
//
// In this example, if n is less than or equal to 0, the Check function will trigger
// a panic with an ErrInvalidParameter error. The RecoverFromPanic function will then
// capture the panic and assign the error to the err variable.
//
// The function then returns the result and the error. If n is greater than 0, the error
// will be nil. If n is less than or equal to 0, the error will be the one captured by
// the RecoverFromPanic function.
//
// The Check function uses the reflect.DeepEqual function to compare the value to the
// zero value of its type, allowing it to be used with any type. However, it may not
// behave as expected with types that have non-exported fields.
//
// Given this, it is recommended to use the Check function primarily with boolean or
// error values, as these are the most common types to check for zero values.
//
//	Check(n <= 0, ...)
//	Check(err, ...)
func Check(val any, err error) {
	zero := reflect.Zero(reflect.TypeOf(val)).Interface()
	if reflect.DeepEqual(val, zero) {
		panic(err)
	}
}
