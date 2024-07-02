// Package errors provides a custom error type for out-of-bound errors.
package common

import (
	"fmt"
	"strconv"
	"strings"
)

// ErrPanic represents an error when a panic occurs.
type ErrPanic struct {
	// Value is the value that caused the panic.
	Value any
}

// Error implements the error interface.
//
// Message: "panic: {value}"
func (e *ErrPanic) Error() string {
	var builder strings.Builder

	builder.WriteString("panic: ")
	fmt.Fprintf(&builder, "%v", e.Value)

	return builder.String()
}

// NewErrPanic creates a new ErrPanic error.
//
// Parameters:
//   - value: The value that caused the panic.
//
// Returns:
//   - *ErrPanic: A pointer to the newly created ErrPanic.
func NewErrPanic(value any) *ErrPanic {
	return &ErrPanic{
		Value: value,
	}
}

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

// Error implements the error interface.
//
// Message: "value (value) not in range <lowerBound, upperBound>"
//
// If the lower bound is inclusive, the message uses square brackets. If the
// upper bound is inclusive, the message uses square brackets. Otherwise, the
// message uses parentheses.
func (e *ErrOutOfBounds) Error() string {
	var builder strings.Builder

	builder.WriteString("value (")
	builder.WriteString(strconv.Itoa(e.Value))
	builder.WriteString(") not in range ")

	if e.LowerInclusive {
		builder.WriteRune('[')
	} else {
		builder.WriteRune('(')
	}

	builder.WriteString(strconv.Itoa(e.LowerBound))
	builder.WriteString(", ")
	builder.WriteString(strconv.Itoa(e.UpperBound))

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
//   - isInclusive: A boolean indicating whether the lower bound is inclusive.
//
// Returns:
//   - *ErrOutOfBound: The error instance for chaining.
func (e *ErrOutOfBounds) WithLowerBound(isInclusive bool) *ErrOutOfBounds {
	e.LowerInclusive = isInclusive

	return e
}

// WithUpperBound sets the inclusivity of the upper bound.
//
// Parameters:
//   - isInclusive: A boolean indicating whether the upper bound is inclusive.
//
// Returns:
//   - *ErrOutOfBound: The error instance for chaining.
func (e *ErrOutOfBounds) WithUpperBound(isInclusive bool) *ErrOutOfBounds {
	e.UpperInclusive = isInclusive

	return e
}

// NewOutOfBounds creates a new ErrOutOfBound error. By default, the lower bound
// is inclusive and the upper bound is exclusive.
//
// Parameters:
//   - lowerBound, upperbound: The lower and upper bounds of the range,
//     respectively.
//   - value: The value that caused the error.
//
// Returns:
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
	Actual string
}

// Error implements the error interface.
//
// Message: "expected <value 0>, <value 1>, <value 2>, ..., or <value n>,
// got <actual> instead"
func (e *ErrUnexpected) Error() string {
	var builder strings.Builder

	builder.WriteString("expected ")

	if len(e.Expected) == 0 {
		builder.WriteString("nothing")
	} else {
		builder.WriteString(strconv.Quote(e.Expected[0]))
	}

	if len(e.Expected) > 2 {
		var values []string

		for i := 1; i < len(e.Expected)-1; i++ {
			values = append(values, strconv.Quote(e.Expected[i]))
		}

		builder.WriteString(strings.Join(values, ", "))
		builder.WriteRune(',')
	}

	builder.WriteString(" or ")
	builder.WriteString(strconv.Quote(e.Expected[len(e.Expected)-1]))
	builder.WriteString(", got ")

	if e.Actual == "" {
		builder.WriteString("nothing")
	} else {
		builder.WriteString(strconv.Quote(e.Actual))
	}

	builder.WriteString(" instead")

	return builder.String()
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
	return &ErrUnexpected{
		Expected: expected,
		Actual:   got,
	}
}

// ErrEmpty represents an error when a value is empty.
type ErrEmpty[T any] struct {
	// Value is the value that caused the error.
	Value T
}

// Error implements the error interface.
//
// Message: "<type> must not be empty"
func (e *ErrEmpty[T]) Error() string {
	var builder strings.Builder

	builder.WriteString(TypeOf(e.Value))
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
	return &ErrEmpty[T]{
		Value: value,
	}
}

// ErrGT represents an error when a value is less than or equal to a specified value.
type ErrGT struct {
	// Value is the value that caused the error.
	Value int
}

// Error implements the error interface.
//
// Message: "value must be greater than <value>"
//
// If the value is 0, the message is "value must be positive".
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
	return &ErrGT{
		Value: value,
	}
}

// ErrLT represents an error when a value is greater than or equal to a specified value.
type ErrLT struct {
	// Value is the value that caused the error.
	Value int
}

// Error implements the error interface.
//
// Message: "value must be less than <value>"
//
// If the value is 0, the message is "value must be negative".
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
	return &ErrLT{
		Value: value,
	}
}

// ErrGTE represents an error when a value is less than a specified value.
type ErrGTE struct {
	// Value is the value that caused the error.
	Value int
}

// Error implements the error interface.
//
// Message: "value must be greater than or equal to <value>"
//
// If the value is 0, the message is "value must be non-negative".
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
	return &ErrGTE{
		Value: value,
	}
}

// ErrLTE represents an error when a value is greater than a specified value.
type ErrLTE struct {
	// Value is the value that caused the error.
	Value int
}

// Error implements the error interface.
//
// Message: "value must be less than or equal to <value>"
//
// If the value is 0, the message is "value must be non-positive".
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
	return &ErrLTE{
		Value: value,
	}
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
	} else if len(e.Values) == 1 {
		return "value must not be " + StringOf(e.Values[0])
	}

	values := make([]string, 0, len(e.Values))
	for _, v := range e.Values {
		values = append(values, StringOf(v))
	}

	var builder strings.Builder

	builder.WriteString("value must not be ")

	if len(e.Values) == 2 {
		builder.WriteString("either ")
		builder.WriteString(values[0])
	} else if len(e.Values) > 2 {
		builder.WriteString(strings.Join(values[:len(values)-1], ", "))
		builder.WriteRune(',')
	}

	builder.WriteString(" or ")
	builder.WriteString(values[len(values)-1])

	return builder.String()
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

// Error implements the error interface.
//
// Message: "type <type> is not a valid <kind> type"
func (e *ErrUnexpectedType[T]) Error() string {
	var builder strings.Builder

	builder.WriteString("type ")
	fmt.Fprintf(&builder, "%T", e.Elem)
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
	return &ErrUnexpectedType[T]{
		Elem: elem,
		Kind: kind,
	}
}

// ErrNilValue represents an error when a value is nil.
type ErrNilValue struct{}

// Error implements the error interface.
//
// Message: "pointer must not be nil"
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

// ErrExhaustedIter is an error type that is returned when an iterator
// is exhausted (i.e., there are no more elements to consume).
type ErrExhaustedIter struct{}

// Error implements the error interface.
//
// Message: "iterator is exhausted"
func (e *ErrExhaustedIter) Error() string {
	return "iterator is exhausted"
}

// NewErrExhaustedIter creates a new ErrExhaustedIter error.
//
// Returns:
//   - *ErrExhaustedIter: A pointer to the new error.
func NewErrExhaustedIter() *ErrExhaustedIter {
	return &ErrExhaustedIter{}
}
