package Errors

import (
	"errors"
	"fmt"
	"strings"
)

type ErrOutOfBound struct {
	lowerBound, upperBound, value  int
	lowerInclusive, upperInclusive bool
}

func NewErrOutOfBound(lowerBound, upperBound, value int) *ErrOutOfBound {
	return &ErrOutOfBound{
		lowerBound:     lowerBound,
		lowerInclusive: false,
		upperBound:     upperBound,
		upperInclusive: true,
		value:          value,
	}
}

func (e *ErrOutOfBound) LowerBound(isInclusive bool) *ErrOutOfBound {
	e.lowerInclusive = isInclusive

	return e
}

func (e *ErrOutOfBound) UpperBound(isInclusive bool) *ErrOutOfBound {
	e.upperInclusive = isInclusive

	return e
}

func (e *ErrOutOfBound) Error() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "value (%d) not in range ", e.value)

	if e.lowerInclusive {
		builder.WriteString("[")
	} else {
		builder.WriteString("(")
	}

	fmt.Fprintf(&builder, "%d, %d", e.lowerBound, e.upperBound)

	if e.upperInclusive {
		builder.WriteString("]")
	} else {
		builder.WriteString(")")
	}

	return builder.String()
}

type ErrInvalidParameter struct {
	parameter string
	reason    error
}

func NewErrInvalidParameter(parameter string) *ErrInvalidParameter {
	return &ErrInvalidParameter{
		parameter: parameter,
		reason:    errors.New("parameter is invalid"),
	}
}

func (e *ErrInvalidParameter) WithReason(reason error) *ErrInvalidParameter {
	if reason == nil {
		return e
	}

	e.reason = reason

	return e
}

func (e *ErrInvalidParameter) Error() string {
	var builder strings.Builder

	builder.WriteString("invalid parameter: ")
	fmt.Fprintf(&builder, "parameter=%s", e.parameter)
	fmt.Fprintf(&builder, ", reason=%v", e.reason)

	return builder.String()
}

func (e *ErrInvalidParameter) Unwrap() error {
	return e.reason
}
