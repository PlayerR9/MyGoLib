// Package errors provides a custom error type for out-of-bound errors.
package Errors

import (
	"errors"
	"fmt"
	"strings"
)

// ErrOutOfBound represents an error when a value is out of a specified range.
// It holds the lower and upper bounds of the range, the value that caused the error,
// and flags indicating whether the bounds are inclusive.
type ErrOutOfBound struct {
	lowerBound, upperBound, value  int
	lowerInclusive, upperInclusive bool
}

// NewErrOutOfBound creates a new ErrOutOfBound error.
// The lower bound is exclusive and the upper bound is inclusive by default.
// The bounds and the value that caused the error are passed as arguments.
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
// If isInclusive is true, the lower bound is inclusive.
// It returns the ErrOutOfBound instance for chaining.
func (e *ErrOutOfBound) WithLowerBound(isInclusive bool) *ErrOutOfBound {
	e.lowerInclusive = isInclusive
	return e
}

// WithUpperBound sets the inclusivity of the upper bound.
// If isInclusive is true, the upper bound is inclusive.
// It returns the ErrOutOfBound instance for chaining.
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
// It holds the parameter name and the reason for its invalidity.
type ErrInvalidParameter struct {
	parameter string
	reason    error
}

// NewErrInvalidParameter creates a new ErrInvalidParameter error.
// The parameter name is passed as an argument.
// The reason for the invalidity is set to "parameter is invalid" by default.
func NewErrInvalidParameter(parameter string) *ErrInvalidParameter {
	return &ErrInvalidParameter{
		parameter: parameter,
		reason:    errors.New("parameter is invalid"),
	}
}

// WithReason sets the reason for the invalidity of the parameter.
// If the reason is nil, it does not change the existing reason.
// It returns the ErrInvalidParameter instance for chaining.
func (e *ErrInvalidParameter) WithReason(reason error) *ErrInvalidParameter {
	if reason == nil {
		return e
	}

	e.reason = reason

	return e
}

// Error generates the error message for the ErrInvalidParameter error.
// The message includes the parameter name and the reason for its invalidity.
func (e *ErrInvalidParameter) Error() string {
	var builder strings.Builder

	builder.WriteString("invalid parameter: ")
	fmt.Fprintf(&builder, "parameter=%s", e.parameter)
	fmt.Fprintf(&builder, ", reason=%v", e.reason)

	return builder.String()
}

// Unwrap returns the reason for the invalidity of the parameter.
// It is used for error unwrapping.
func (e *ErrInvalidParameter) Unwrap() error {
	return e.reason
}
