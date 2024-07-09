// Package errors provides a custom error type for out-of-bound errors.
package common

import (
	"errors"
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

	str := builder.String()

	return str
}

// NewErrPanic creates a new ErrPanic error.
//
// Parameters:
//   - value: The value that caused the panic.
//
// Returns:
//   - *ErrPanic: A pointer to the newly created ErrPanic.
func NewErrPanic(value any) *ErrPanic {
	e := &ErrPanic{
		Value: value,
	}

	return e
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
	qi := QuoteInt(e.Value)
	left_bound := strconv.Itoa(e.LowerBound)
	right_bound := strconv.Itoa(e.UpperBound)

	var open, close string

	if e.LowerInclusive {
		open = "["
	} else {
		open = "("
	}

	if e.UpperInclusive {
		close = "]"
	} else {
		close = ")"
	}

	values := []string{
		"value",
		qi,
		"not in range",
		open,
		left_bound,
		",",
		right_bound,
		close,
	}

	str := strings.Join(values, " ")

	return str
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
	e := &ErrOutOfBounds{
		LowerBound:     lowerBound,
		UpperBound:     upperBound,
		LowerInclusive: true,
		UpperInclusive: false,
		Value:          value,
	}
	return e
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
	var expected string

	if len(e.Expected) == 0 {
		expected = "nothing"
	} else {
		expected = OrQuoteString(e.Expected, false)
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

// ErrEmpty represents an error when a value is empty.
type ErrEmpty[T any] struct {
	// Value is the value that caused the error.
	Value T
}

// Error implements the error interface.
//
// Message: "<type> must not be empty"
func (e *ErrEmpty[T]) Error() string {
	type_str := TypeOf(e.Value)

	var builder strings.Builder

	builder.WriteString(type_str)
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
	e := &ErrEmpty[T]{
		Value: value,
	}
	return e
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

	value := strconv.Itoa(e.Value)

	var builder strings.Builder

	builder.WriteString("value must be greater than ")
	builder.WriteString(value)

	str := builder.String()

	return str
}

// NewErrGT creates a new ErrGT error with the specified value.
//
// Parameters:
//   - value: The minimum value that is not allowed.
//
// Returns:
//   - *ErrGT: A pointer to the newly created ErrGT.
func NewErrGT(value int) *ErrGT {
	e := &ErrGT{
		Value: value,
	}
	return e
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

	value := strconv.Itoa(e.Value)

	var builder strings.Builder

	builder.WriteString("value must be less than ")
	builder.WriteString(value)

	str := builder.String()
	return str
}

// NewErrLT creates a new ErrLT error with the specified value.
//
// Parameters:
//   - value: The maximum value that is not allowed.
//
// Returns:
//   - *ErrLT: A pointer to the newly created ErrLT.
func NewErrLT(value int) *ErrLT {
	e := &ErrLT{
		Value: value,
	}
	return e
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

	value := strconv.Itoa(e.Value)

	var builder strings.Builder

	builder.WriteString("value must be greater than or equal to ")
	builder.WriteString(value)

	str := builder.String()
	return str
}

// NewErrGTE creates a new ErrGTE error with the specified value.
//
// Parameters:
//   - value: The minimum value that is allowed.
//
// Returns:
//   - *ErrGTE: A pointer to the newly created ErrGTE.
func NewErrGTE(value int) *ErrGTE {
	e := &ErrGTE{
		Value: value,
	}
	return e
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

	value := strconv.Itoa(e.Value)

	var builder strings.Builder

	builder.WriteString("value must be less than or equal to ")
	builder.WriteString(value)

	str := builder.String()
	return str
}

// NewErrLTE creates a new ErrLTE error with the specified value.
//
// Parameters:
//   - value: The maximum value that is allowed.
//
// Returns:
//   - *ErrLTE: A pointer to the newly created ErrLTE.
func NewErrLTE(value int) *ErrLTE {
	e := &ErrLTE{
		Value: value,
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

	value := OrString(values, true)

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
	values := []string{
		"type",
		fmt.Sprintf("%T", e.Elem),
		"is not a valid",
		e.Kind,
		"type",
	}

	str := strings.Join(values, " ")

	return str
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
	e := &ErrUnexpectedType[T]{
		Elem: elem,
		Kind: kind,
	}
	return e
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
	e := &ErrNilValue{}
	return e
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
	e := &ErrExhaustedIter{}
	return e
}

// ErrNotFound is an error type that is returned when a value is not found.
type ErrNotFound struct{}

// Error implements the error interface.
//
// Message: "not found"
func (e *ErrNotFound) Error() string {
	return "not found"
}

// NewErrNotFound creates a new ErrNotFound error.
//
// Returns:
//   - *ErrNotFound: A pointer to the new error.
func NewErrNotFound() *ErrNotFound {
	e := &ErrNotFound{}
	return e
}

// IsNotFound checks if the error is an ErrNotFound.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - bool: True if the error is an ErrNotFound, false otherwise.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}

	var notFound *ErrNotFound

	ok := errors.As(err, &notFound)
	return ok
}
