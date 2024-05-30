package pkg

import "strings"

// ErrFormatError is an error type that is returned when a format error occurs.
type ErrFormatError struct {
	// Reason is the reason for the error.
	Reason error
}

// Error returns the error message: "format error: <reason>".
//
// Returns:
//   - string: The error message.
//
// Behaviors:
//   - If Reason is nil, the error message is: "format is invalid" instead.
func (e *ErrFormatError) Error() string {
	if e.Reason == nil {
		return "format is invalid"
	}

	var builder strings.Builder

	builder.WriteString("format error: ")
	builder.WriteString(e.Reason.Error())

	return builder.String()
}

// NewErrFormatError creates a new ErrFormatError error.
//
// Parameters:
//   - reason: The reason for the error.
//
// Returns:
//   - *ErrFormatError: A pointer to the new error.
func NewErrFormatError(reason error) *ErrFormatError {
	return &ErrFormatError{
		Reason: reason,
	}
}
