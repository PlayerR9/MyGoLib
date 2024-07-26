package errors

import (
	"fmt"
	"strconv"
	"strings"

	utstr "github.com/PlayerR9/MyGoLib/Utility/strings"
)

// ErrUnexpected represents an error that occurs when an unexpected value is
// encountered.
type ErrUnexpected struct {
	// Expected is the list of expected values.
	Expected []string

	// Actual is the actual value encountered.
	Actual string
}

// Error implements the error interface.
//
// Message: "expected <value 0>, <value 1>, <value 2>, ..., or <value n>,
// got <actual> instead"
func (e *ErrUnexpected) Error() string {
	var expected string

	if len(e.Expected) == 0 {
		expected = "nothing"
	} else {
		expected = utstr.EitherOrString(e.Expected, true)
	}

	var actual string

	if e.Actual == "" {
		actual = "nothing"
	} else {
		actual = strconv.Quote(e.Actual)
	}

	values := []string{
		"expected",
		expected,
		", got",
		actual,
		"instead",
	}

	str := strings.Join(values, " ")

	return str
}

// NewErrUnexpected creates a new ErrUnexpected error.
//
// Parameters:
//   - got: The actual value encountered.
//   - expected: The list of expected values.
//
// Returns:
//   - *ErrUnexpected: A pointer to the newly created ErrUnexpected.
func NewErrUnexpected(got string, expected ...string) *ErrUnexpected {
	e := &ErrUnexpected{
		Expected: expected,
		Actual:   got,
	}

	return e
}

// ErrInvalidValues represents an error when a value is in a list of invalid values.
type ErrInvalidValues[T comparable] struct {
	// Values is the list of invalid values.
	Values []T
}

// Error implements the error interface.
//
// Message: "value must not be <value 0>, <value 1>, <value 2>, ..., or <value n>"
//
// If no values are provided, the message is "value is invalid".
func (e *ErrInvalidValues[T]) Error() string {
	if len(e.Values) == 0 {
		return "value is invalid"
	}

	values := make([]string, 0, len(e.Values))
	for _, v := range e.Values {
		values = append(values, fmt.Sprintf("%v", v))
	}

	value := utstr.OrString(values, false, true)

	var builder strings.Builder

	builder.WriteString("value must not be ")
	builder.WriteString(value)

	str := builder.String()
	return str
}

// NewErrInvalidValues creates a new ErrInvalidValues error.
//
// Parameters:
//   - values: The list of invalid values.
//
// Returns:
//   - *ErrInvalidValues: A pointer to the newly created ErrInvalidValues.
func NewErrInvalidValues[T comparable](values []T) *ErrInvalidValues[T] {
	e := &ErrInvalidValues[T]{
		Values: values,
	}
	return e
}

// NewErrUnexpectedValue is a function that creates a new ErrInvalidValues error.
//
// Parameters:
//   - value: The value that was unexpected.
//
// Returns:
//   - *ErrInvalidValues: A pointer to the newly created ErrInvalidValues.
func NewErrUnexpectedValue[T comparable](value T) *ErrInvalidValues[T] {
	e := &ErrInvalidValues[T]{
		Values: []T{value},
	}
	return e
}
