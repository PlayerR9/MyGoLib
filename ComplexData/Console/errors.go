package Console

import (
	"strconv"
	"strings"
)

// ErrMissingArgument is an error that is returned when a required argument is missing.
type ErrMissingArgument struct {
	// Name is the name of the missing argument.
	Name string
}

// Error returns the error message: "argument <name> is required".
//
// Returns:
//   - string: The error message.
func (e *ErrMissingArgument) Error() string {
	var builder strings.Builder

	builder.WriteString("argument ")
	builder.WriteString(strconv.Quote(e.Name))
	builder.WriteString(" is required")

	return builder.String()
}

// NewErrMissingArgument creates a new ErrMissingArgument.
//
// Parameters:
//   - name: The name of the missing argument.
//
// Returns:
//   - *ErrMissingArgument: The new ErrMissingArgument.
func NewErrMissingArgument(name string) *ErrMissingArgument {
	return &ErrMissingArgument{
		Name: name,
	}
}

// ErrArgumentNotRecognized is an error that is returned when an argument is not recognized.
type ErrArgumentNotRecognized struct {
	// Value is the value of the unrecognized argument.
	Value string
}

// Error returns the error message: "argument <value> is not recognized".
//
// Returns:
//   - string: The error message.
func (e *ErrArgumentNotRecognized) Error() string {
	var builder strings.Builder

	builder.WriteString("argument ")
	builder.WriteString(strconv.Quote(e.Value))
	builder.WriteString(" is not recognized")

	return builder.String()
}

// NewErrArgumentNotRecognized creates a new ErrArgumentNotRecognized.
//
// Parameters:
//   - value: The value of the unrecognized argument.
//
// Returns:
//   - *ErrArgumentNotRecognized: The new ErrArgumentNotRecognized.
func NewErrArgumentNotRecognized(value string) *ErrArgumentNotRecognized {
	return &ErrArgumentNotRecognized{
		Value: value,
	}
}
