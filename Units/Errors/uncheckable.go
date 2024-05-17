// Package errors provides a custom error type for out-of-bound errors.
package Errors

import (
	"fmt"
	"strconv"
	"strings"

	com "github.com/PlayerR9/MyGoLib/Units/Common"
)

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

// ErrorIf returns the error if the target value is less than or equal to the
// specified value.
//
// Parameters:
//   - target: The value to check.
//
// Returns:
//   - error: The error if the target value is less than or equal to the specified
//     value, nil otherwise.
func (e *ErrGT) ErrorIf(target int) error {
	if target > e.Value {
		return nil
	} else {
		return e
	}
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

// ErrorIf returns the error if the target value is greater than or equal to the
// specified value.
//
// Parameters:
//   - target: The value to check.
//
// Returns:
//   - error: The error if the target value is greater than or equal to the specified
//     value, nil otherwise.
func (e *ErrLT) ErrorIf(target int) error {
	if target < e.Value {
		return nil
	} else {
		return e
	}
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

// ErrorIf returns the error if the target value is less than the specified value.
//
// Parameters:
//   - target: The value to check.
//
// Returns:
//   - error: The error if the target value is less than the specified value, nil
//     otherwise.
func (e *ErrGTE) ErrorIf(target int) error {
	if target >= e.Value {
		return nil
	} else {
		return e
	}
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

// ErrorIf returns the error if the target value is greater than the specified value.
//
// Parameters:
//   - target: The value to check.
//
// Returns:
//   - error: The error if the target value is greater than the specified value, nil
//     otherwise.
func (e *ErrLTE) ErrorIf(target int) error {
	if target <= e.Value {
		return nil
	} else {
		return e
	}
}

// ErrInvalidValues represents an error when a value is in a list of invalid values.
type ErrInvalidValues[T any] struct {
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
func NewErrInvalidValues[T any](values []T) *ErrInvalidValues[T] {
	return &ErrInvalidValues[T]{
		Values: values,
	}
}

// ErrorIf returns the error if the target value is in the list of invalid values.
//
// Parameters:
//   - equals: A function that returns true if the target value is equal to the
//     value in the list of invalid values.
//
// Returns:
//   - error: The error if the target value is in the list of invalid values, nil
//     otherwise.
func (e *ErrInvalidValues[T]) ErrorIf(equals func(T) bool) error {
	for _, v := range e.Values {
		if equals(v) {
			return e
		}
	}

	return nil
}

// NewErrUnexpectedValue is a function that creates a new ErrInvalidValues error.
//
// Parameters:
//   - value: The value that was unexpected.
//
// Returns:
//   - *ErrInvalidValues: A pointer to the newly created ErrInvalidValues.
func NewErrUnexpectedValue[T any](value T) *ErrInvalidValues[T] {
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

// ErrorIf returns the error if the target value is nil.
//
// Parameters:
//   - target: The value to check.
//
// Returns:
//   - error: The error if the target value is nil, nil otherwise.
func (e *ErrNilValue) ErrorIf(target any) error {
	if target == nil {
		return e
	} else {
		return nil
	}
}
