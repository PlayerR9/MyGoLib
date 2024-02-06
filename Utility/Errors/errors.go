// Package errors provides a custom error type for out-of-bound errors.
package Errors

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	itf "github.com/PlayerR9/MyGoLib/Interfaces"
)

// ConditionCheck is a function that checks a condition and returns an error
// if the condition is not met.
type ConditionCheck func(any) error

// ErrOutOfBound represents an error when a value is out of a specified range.
type ErrOutOfBound struct {
	// lowerBound and upperBound are the lower and upper bounds of the range,
	// respectively.
	lowerBound, upperBound int

	// lowerInclusive and upperInclusive are flags indicating whether the lower
	// and upper bounds are inclusive, respectively.
	lowerInclusive, upperInclusive bool
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

	fmt.Fprintf(&builder, "value not in range ")
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
func NewErrOutOfBound(lowerBound, upperBound int) error {
	return &ErrOutOfBound{
		lowerBound:     lowerBound,
		upperBound:     upperBound,
		lowerInclusive: true,
		upperInclusive: false,
	}
}

// OutOfBound is a constructor function that creates a new ErrOutOfBound error.
// If no inclusivity flags are provided, the lower bound is inclusive and the
// upper bound is exclusive.
// The difference between this function and NewErrOutOfBound is that this
// function is used in conjuction with Check to create a fluent interface for
// error handling.
//
// Parameters:
//
//   - lowerBound, upperBound: The lower and upper bounds of the range,
//     respectively.
//
// Returns:
//
//   - *ErrOutOfBound: A pointer to the newly created ErrOutOfBound.
func OutOfBound(lowerBound, upperBound int) *ErrOutOfBound {
	return &ErrOutOfBound{
		lowerBound:     lowerBound,
		upperBound:     upperBound,
		lowerInclusive: true,
		upperInclusive: false,
	}
}

// IsIn creates a condition check that validates if a value is within the range
// specified by the ErrOutOfBound error.
//
// Returns:
//
//   - ConditionCheck: A function that checks if a value is within the range.
func (e *ErrOutOfBound) IsIn() ConditionCheck {
	return func(v any) error {
		if v.(int) >= e.lowerBound && v.(int) <= e.upperBound {
			return nil
		}

		return e
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

// NewInvalidParameter creates a new ErrInvalidParameter error.
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
func NewInvalidParameter(parameter string, reason error) error {
	if reason == nil {
		reason = errors.New("parameter is invalid")
	}

	return &ErrInvalidParameter{
		parameter: parameter,
		reason:    reason,
	}
}

// InvalidParameter creates a new ErrInvalidParameter error.
// The difference between this function and NewInvalidParameter is that this
// function is used in conjuction with Check to create a fluent interface for
// error handling.
//
// Parameters:
//
//   - parameter: The name of the parameter.
//
// Returns:
//
//   - *ErrInvalidParameter: A pointer to the new ErrInvalidParameter.
func InvalidParameter(parameter string) *ErrInvalidParameter {
	return &ErrInvalidParameter{
		parameter: parameter,
	}
}

// Not creates a condition check that validates if a value is not equal to the
// provided value. Errors if the value is equal to the provided value.
//
// Parameters:
//
//   - other: The value to compare against.
//
// Returns:
//
//   - ConditionCheck: A function that checks if a value is not equal to the
//     provided value.
func (e *ErrInvalidParameter) Not(other any) ConditionCheck {
	return func(v any) error {
		if !reflect.DeepEqual(v, other) {
			return nil
		}

		return &ErrInvalidParameter{
			parameter: e.parameter,
			reason:    fmt.Errorf("value (%v) is not allowed", v),
		}
	}
}

// NotNil creates a condition check that validates if a value is not nil.
//
// Returns:
//
//   - ConditionCheck: A function that checks if a value is not nil.
func (e *ErrInvalidParameter) NotNil() ConditionCheck {
	return func(v any) error {
		if v != nil {
			return nil
		}

		return &ErrInvalidParameter{
			parameter: e.parameter,
			reason:    fmt.Errorf("%T was found to be nil", v),
		}
	}
}

// NotZero creates a condition check that validates if a value is not the zero
// value of its type.
//
// Returns:
//
//   - ConditionCheck: A function that checks if a value is not the zero value
//     of its type.
func (e *ErrInvalidParameter) NotZero() ConditionCheck {
	return func(v any) error {
		if !reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface()) {
			return nil
		}

		return &ErrInvalidParameter{
			parameter: e.parameter,
			reason:    fmt.Errorf("value (%v) is the zero value of its type", v),
		}
	}
}

func (c *ErrInvalidParameter) Lt(other any) ConditionCheck {
	return func(v any) error {
		x, ok := v.(interface{ Less(other any) bool })
		if ok {
			if x.Less(other) {
				return nil
			}

			return &ErrInvalidParameter{
				parameter: c.parameter,
				reason:    fmt.Errorf("value (%v) is not less than %v", v, other),
			}
		}

		y, ok := other.(itf.InequalityComparable)

		if itf.Less(v, other) {
			return nil
		}

		return &ErrInvalidParameter{
			parameter: c.parameter,
			reason:    fmt.Errorf("value (%v) is not less than %v", v, other),
		}
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
func NewCallFailed(functionName string, function any, reason error) error {
	if reason == nil {
		reason = errors.New("an error occurred while calling the function")
	}

	return &ErrCallFailed{
		fnName:    functionName,
		signature: reflect.ValueOf(function).Type(),
		reason:    reason,
	}
}

// CallFailed creates a new ErrCallFailed.
// The difference between this function and NewCallFailed is that this function
// is used in conjuction with Check to create a fluent interface for error
// handling.
//
// Parameters:
//
//   - functionName is the name of the function.
//   - function is the function that failed.
//
// Returns:
//
//   - *ErrCallFailed: A pointer to the new ErrCallFailed.
func CallFailed(functionName string, function any) *ErrCallFailed {
	return &ErrCallFailed{
		fnName:    functionName,
		signature: reflect.ValueOf(function).Type(),
	}
}

// Success creates a condition check that validates if a function call is
// successful. This is useful for functions that do not return a value, but
// can fail.
//
// Parameters:
//
//   - f: The function to call.
//
// Returns:
//
//   - ConditionCheck: A function that checks if a function call is successful.
func (e *ErrCallFailed) Success(f func()) ConditionCheck {
	return func(v any) (err error) {
		defer func() {
			if r := recover(); r != nil {
				if recErr, ok := r.(error); ok {
					err = recErr
				} else {
					err = fmt.Errorf("panic: %v", r)
				}
			}
		}()

		f()

		return
	}
}
