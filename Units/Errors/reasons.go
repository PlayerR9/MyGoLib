// Package errors provides a custom error type for out-of-bound errors.
package Errors

import (
	"fmt"
	"strconv"
	"strings"

	com "github.com/PlayerR9/MyGoLib/Units/Common"
)

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

	fmt.Fprintf(&builder, "value (%d) not in range ", e.Value)
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
			builder.WriteRune(',')
			builder.WriteRune(' ')
			fmt.Fprintf(&builder, "%q", e.Expected[i])
		}

		builder.WriteRune(',')
		builder.WriteRune(' ')
		builder.WriteString("or ")

		fmt.Fprintf(&builder, "%q", e.Expected[len(e.Expected)-1])

		expected = builder.String()
	}

	if e.Actual == nil {
		got = "nothing"
	} else {
		got = e.Actual.String()
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

// ErrEmpty represents an error when a value is empty.
type ErrEmpty[T any] struct {
	// Value is the value that caused the error.
	Value T
}

// Error returns the error message: "value must not be empty".
//
// Returns:
//   - string: The error message.
func (e *ErrEmpty[T]) Error() string {
	var builder strings.Builder

	builder.WriteString(com.TypeOf(e.Value))
	builder.WriteString(" must not be empty")

	return builder.String()
}

// NewErrEmpty creates a new ErrEmpty error.
//
// Parameters:
//   - value: The value that caused the error.
//
// Returns:
//   - *ErrEmpty: A pointer to the newly created ErrEmpty.
func NewErrEmpty[T any](value T) *ErrEmpty[T] {
	return &ErrEmpty[T]{Value: value}
}

// ErrNotComparable represents an error when a value is not comparable.
type ErrNotComparable[T any] struct {
	// Value is the value that caused the error.
	Value T
}

// Error returns the error message: "type <type> does not support comparison".
//
// Returns:
//   - string: The error message.
func (e *ErrNotComparable[T]) Error() string {
	var builder strings.Builder

	builder.WriteString("type ")
	builder.WriteString(fmt.Sprintf("%T", e.Value))
	builder.WriteString(" does not support comparison")

	return builder.String()
}

// NewErrNotComparable creates a new ErrNotComparable error.
//
// Returns:
//   - *ErrNotComparable: A pointer to the newly created ErrNotComparable.
func NewErrNotComparable[T any](value T) *ErrNotComparable[T] {
	return &ErrNotComparable[T]{
		Value: value,
	}
}

// ErrGT represents an error when a value is less than or equal to a specified value.
type ErrGT struct {
	// Value is the value that caused the error.
	Value int
}

// Error returns the error message: "value must be greater than <value>"
//
// Returns:
//   - string: The error message.
//
// Behaviors:
//   - If the value is 0, the error message is "value must be positive".
func (e *ErrGT) Error() string {
	if e.Value == 0 {
		return "value must be positive"
	}

	var builder strings.Builder

	builder.WriteString("value must be greater than ")
	builder.WriteString(strconv.Itoa(e.Value))

	return builder.String()
}

// NewErrGT creates a new ErrGT error with the specified value.
//
// Parameters:
//   - value: The minimum value that is not allowed.
//
// Returns:
//   - *ErrGT: A pointer to the newly created ErrGT.
func NewErrGT(value int) *ErrGT {
	return &ErrGT{Value: value}
}

// ErrLT represents an error when a value is greater than or equal to a specified value.
type ErrLT struct {
	// Value is the value that caused the error.
	Value int
}

// Error returns the error message: "value must be less than <value>"
//
// Returns:
//   - string: The error message.
//
// Behaviors:
//   - If the value is 0, the error message is "value must be negative".
func (e *ErrLT) Error() string {
	if e.Value == 0 {
		return "value must be negative"
	}

	var builder strings.Builder

	builder.WriteString("value must be less than ")
	builder.WriteString(strconv.Itoa(e.Value))

	return builder.String()
}

// NewErrLT creates a new ErrLT error with the specified value.
//
// Parameters:
//   - value: The maximum value that is not allowed.
//
// Returns:
//   - *ErrLT: A pointer to the newly created ErrLT.
func NewErrLT(value int) *ErrLT {
	return &ErrLT{Value: value}
}

// ErrGTE represents an error when a value is less than a specified value.
type ErrGTE struct {
	// Value is the value that caused the error.
	Value int
}

// Error returns the error message: "value must be greater than or equal to <value>"
//
// Returns:
//   - string: The error message.
//
// Behaviors:
//   - If the value is 0, the error message is "value must be non-negative".
func (e *ErrGTE) Error() string {
	if e.Value == 0 {
		return "value must be non-negative"
	}

	var builder strings.Builder

	builder.WriteString("value must be greater than or equal to ")
	builder.WriteString(strconv.Itoa(e.Value))

	return builder.String()
}

// NewErrGTE creates a new ErrGTE error with the specified value.
//
// Parameters:
//   - value: The minimum value that is allowed.
//
// Returns:
//   - *ErrGTE: A pointer to the newly created ErrGTE.
func NewErrGTE(value int) *ErrGTE {
	return &ErrGTE{Value: value}
}

// ErrLTE represents an error when a value is greater than a specified value.
type ErrLTE struct {
	// Value is the value that caused the error.
	Value int
}

// Error returns the error message: "value must be less than or equal to <value>"
//
// Returns:
//   - string: The error message.
//
// Behaviors:
//   - If the value is 0, the error message is "value must be non-positive".
func (e *ErrLTE) Error() string {
	if e.Value == 0 {
		return "value must be non-positive"
	}

	var builder strings.Builder

	builder.WriteString("value must be less than or equal to ")
	builder.WriteString(strconv.Itoa(e.Value))

	return builder.String()
}

// NewErrLTE creates a new ErrLTE error with the specified value.
//
// Parameters:
//   - value: The maximum value that is allowed.
//
// Returns:
//   - *ErrLTE: A pointer to the newly created ErrLTE.
func NewErrLTE(value int) *ErrLTE {
	return &ErrLTE{Value: value}
}

// ErrInvalidValues represents an error when a value is in a list of invalid values.
type ErrInvalidValues[T comparable] struct {
	// Values is the list of invalid values.
	Values []T
}

// Error returns the error message: "value must not be <values>"
// according to the following format:
//
//	<value 0>, <value 1>, <value 2>, ..., or <value n>
//
// Returns:
//   - string: The error message.
//
// Behaviors:
//   - If there are no values, the error message is "value is invalid".
//   - If there is one value, the error message is "value must not be
//     <value 0>".
//   - If there are two values, the error message is "value must not be
//     either <value 0> or <value 1>".
func (e *ErrInvalidValues[T]) Error() string {
	switch len(e.Values) {
	case 0:
		return "value is invalid"
	case 1:
		var builder strings.Builder

		builder.WriteString("value must not be ")
		builder.WriteString(com.StringOf(e.Values[0]))

		return builder.String()
	case 2:
		var builder strings.Builder

		builder.WriteString("value must not be either ")
		builder.WriteString(com.StringOf(e.Values[0]))
		builder.WriteString(" or ")
		builder.WriteString(com.StringOf(e.Values[1]))

		return builder.String()
	default:
		values := make([]string, 0, len(e.Values))

		for _, v := range e.Values {
			values = append(values, com.StringOf(v))
		}

		var builder strings.Builder

		builder.WriteString("value must not be ")
		builder.WriteString(strings.Join(values[:len(values)-1], ", "))

		builder.WriteRune(',')
		builder.WriteRune(' ')
		builder.WriteString("or ")
		builder.WriteString(values[len(values)-1])

		return builder.String()
	}
}

// NewErrInvalidValues creates a new ErrInvalidValues error.
//
// Parameters:
//   - values: The list of invalid values.
//
// Returns:
//   - *ErrInvalidValues: A pointer to the newly created ErrInvalidValues.
func NewErrInvalidValues[T comparable](values []T) *ErrInvalidValues[T] {
	return &ErrInvalidValues[T]{
		Values: values,
	}
}

// NewErrUnexpectedValue is a function that creates a new ErrInvalidValues error.
//
// Parameters:
//   - value: The value that was unexpected.
//
// Returns:
//   - *ErrInvalidValues: A pointer to the newly created ErrInvalidValues.
func NewErrUnexpectedValue[T comparable](value T) *ErrInvalidValues[T] {
	return &ErrInvalidValues[T]{
		Values: []T{value},
	}
}

// ErrUnexpectedType represents an error when a value has an invalid type.
type ErrUnexpectedType[T any] struct {
	// Elem is the element that caused the error.
	Elem T

	// Kind is the category of the type that was expected.
	Kind string
}

// Error returns the error message: "type <type> is not a valid <kind> type".
//
// Returns:
//   - string: The error message.
func (e *ErrUnexpectedType[T]) Error() string {
	var builder strings.Builder

	builder.WriteString("type ")
	builder.WriteString(fmt.Sprintf("%T", e.Elem))
	builder.WriteString(" is not a valid ")
	builder.WriteString(e.Kind)
	builder.WriteString(" type")

	return builder.String()
}

// NewErrUnexpectedType creates a new ErrUnexpectedType error.
//
// Parameters:
//   - typeName: The name of the type that was expected.
//   - elem: The element that caused the error.
//
// Returns:
//   - *ErrUnexpectedType: A pointer to the newly created ErrUnexpectedType.
func NewErrUnexpectedType[T any](kind string, elem T) *ErrUnexpectedType[T] {
	return &ErrUnexpectedType[T]{Elem: elem, Kind: kind}
}

// ErrInvalidCharacter represents an error when an invalid character is found.
type ErrInvalidCharacter struct {
	// Character is the invalid character.
	Character rune
}

// Error returns the error message: "character (<character>) is invalid".
//
// Returns:
//   - string: The error message.
func (e *ErrInvalidCharacter) Error() string {
	var builder strings.Builder

	builder.WriteString("character (")
	builder.WriteRune(e.Character)
	builder.WriteRune(')')
	builder.WriteString(" is invalid")

	return builder.String()
}

// NewErrInvalidCharacter creates a new ErrInvalidCharacter error.
//
// Parameters:
//   - character: The invalid character.
//
// Returns:
//   - *ErrInvalidCharacter: A pointer to the newly created ErrInvalidCharacter.
func NewErrInvalidCharacter(character rune) *ErrInvalidCharacter {
	return &ErrInvalidCharacter{Character: character}
}

// ErrNilValue represents an error when a value is nil.
type ErrNilValue struct{}

// Error returns the error message: "value must not be nil".
//
// Returns:
//   - string: The error message.
func (e *ErrNilValue) Error() string {
	return "pointer must not be nil"
}

// NewErrNilValue creates a new ErrNilValue error.
//
// Returns:
//   - *ErrNilValue: The new ErrNilValue error.
func NewErrNilValue() *ErrNilValue {
	return &ErrNilValue{}
}
