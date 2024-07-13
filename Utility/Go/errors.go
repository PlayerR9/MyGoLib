package Go

import (
	"errors"
	"strings"
)

// ErrNotGeneric is an error type for when a type is not a generic.
type ErrNotGeneric struct {
	Reason error
}

// Error implements the error interface.
//
// Message: "not a generic type"
func (e *ErrNotGeneric) Error() string {
	if e.Reason == nil {
		return "not a generic type"
	}

	var builder strings.Builder

	builder.WriteString("not a generic type: ")
	builder.WriteString(e.Reason.Error())

	str := builder.String()

	return str
}

// NewErrNotGeneric creates a new ErrNotGeneric error.
//
// Parameters:
//   - reason: The reason for the error.
//
// Returns:
//   - *ErrNotGeneric: The new error.
func NewErrNotGeneric(reason error) *ErrNotGeneric {
	e := &ErrNotGeneric{
		Reason: reason,
	}

	return e
}

// IsErrNotGeneric checks if an error is of type ErrNotGeneric.
//
// Parameters:
//   - target: The error to check.
//
// Returns:
//   - bool: True if the error is of type ErrNotGeneric, false otherwise.
func IsErrNotGeneric(target error) bool {
	if target == nil {
		return false
	}

	var targetErr *ErrNotGeneric

	ok := errors.As(target, &targetErr)
	return ok
}
