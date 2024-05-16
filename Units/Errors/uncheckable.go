// Package errors provides a custom error type for out-of-bound errors.
package Errors

import (
	"fmt"
	"strings"
)

// ErrGT represents an error when a value is less than or equal to a specified value.
type ErrGT struct {
	// Value is the value that caused the error.
	Value int
}

// Error is a method of the error interface that returns the error message.
//
// Returns:
//
//   - string: The error message.
func (e *ErrGT) Error() string {
	if e.Value == 0 {
		return "value must be positive"
	} else {
		return fmt.Sprintf("value must be greater than %d", e.Value)
	}
}

// NewErrGT creates a new ErrGT error.
//
// Parameters:
//
//   - value: The value that caused the error.
//
// Returns:
//
//   - *ErrGT: A pointer to the newly created ErrGT.
func NewErrGT(value int) *ErrGT {
	return &ErrGT{Value: value}
}

// ErrLT represents an error when a value is greater than or equal to a specified value.
type ErrLT struct {
	// Value is the value that caused the error.
	Value int
}

// Error is a method of the error interface that returns the error message.
//
// Returns:
//
//   - string: The error message.
func (e *ErrLT) Error() string {
	if e.Value == 0 {
		return "value must be negative"
	} else {
		return fmt.Sprintf("value must be less than %d", e.Value)
	}
}

// NewErrLT creates a new ErrLT error.
//
// Parameters:
//
//   - value: The value that caused the error.
//
// Returns:
//
//   - *ErrLT: A pointer to the newly created ErrLT.
func NewErrLT(value int) *ErrLT {
	return &ErrLT{Value: value}
}

// ErrGTE represents an error when a value is less than a specified value.
type ErrGTE struct {
	// Value is the value that caused the error.
	Value int
}

// Error is a method of the error interface that returns the error message.
//
// Returns:
//
//   - string: The error message.
func (e *ErrGTE) Error() string {
	if e.Value == 0 {
		return "value must be non-negative"
	} else {
		return fmt.Sprintf("value must be greater than %d", e.Value)
	}
}

// NewErrGTE creates a new ErrGTE error.
//
// Parameters:
//
//   - value: The value that caused the error.
//
// Returns:
//
//   - *ErrGTE: A pointer to the newly created ErrGTE.
func NewErrGTE(value int) *ErrGTE {
	return &ErrGTE{Value: value}
}

// ErrLTE represents an error when a value is greater than a specified value.
type ErrLTE struct {
	// Value is the value that caused the error.
	Value int
}

// Error is a method of the error interface that returns the error message.
//
// Returns:
//
//   - string: The error message.
func (e *ErrLTE) Error() string {
	if e.Value == 0 {
		return "value must be non-positive"
	} else {
		return fmt.Sprintf("value must be less than or equal to %d", e.Value)
	}
}

// NewErrLTE creates a new ErrLTE error.
//
// Parameters:
//
//   - value: The value that caused the error.
//
// Returns:
//
//   - *ErrLTE: A pointer to the newly created ErrLTE.
func NewErrLTE(value int) *ErrLTE {
	return &ErrLTE{Value: value}
}

// ErrInvalidValue represents an error when a value is invalid.
type ErrInvalidValue struct {
	// Values is the list of invalid values.
	Values []int
}

// Error is a method of the error interface that returns the error message.
//
// Returns:
//
//   - string: The error message.
func (e *ErrInvalidValue) Error() string {
	if len(e.Values) == 0 {
		return "value is invalid"
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "value must not be %d", e.Values[0])

	if len(e.Values) > 1 {
		for i := 1; i < len(e.Values)-1; i++ {
			builder.WriteRune(',')
			builder.WriteRune(' ')
			fmt.Fprintf(&builder, "%d", e.Values[i])
		}

		builder.WriteRune(',')
		builder.WriteRune(' ')
		builder.WriteString("or ")
		fmt.Fprintf(&builder, "%d", e.Values[len(e.Values)-1])
	}

	return builder.String()
}

// NewErrInvalidValue creates a new ErrInvalidValue error.
//
// Parameters:
//
//   - values: The list of invalid values.
//
// Returns:
//
//   - *ErrInvalidValue: A pointer to the newly created ErrInvalidValue.
func NewErrInvalidValue(values ...int) *ErrInvalidValue {
	return &ErrInvalidValue{Values: values}
}

// ErrEmptyString represents an error when a string is empty.
type ErrEmptyString struct{}

// Error is a method of the error interface that returns the error message.
//
// Returns:
//
//   - string: The error message.
func (e *ErrEmptyString) Error() string {
	return "value must not be empty"
}

// NewErrEmptyString creates a new ErrEmptyString error.
//
// Returns:
//
//   - *ErrEmptyString: A pointer to the newly created ErrEmptyString.
func NewErrEmptyString() *ErrEmptyString {
	return &ErrEmptyString{}
}

// ErrEmptySlice represents an error when a slice is empty.
type ErrEmptySlice struct{}

// Error is a method of the error interface that returns the error message.
//
// Returns:
//   - string: The error message.
func (e *ErrEmptySlice) Error() string {
	return "slice must not be empty"
}

// NewErrEmptySlice creates a new ErrEmptySlice error.
//
// Returns:
//   - *ErrEmptySlice: A pointer to the newly created ErrEmptySlice.
func NewErrEmptySlice() *ErrEmptySlice {
	return &ErrEmptySlice{}
}

// ErrInvalidType represents an error when a value has an invalid type.
type ErrInvalidType[T any] struct {
	// Elem is the element that caused the error.
	Elem T

	// Type is the name of the type that was expected.
	Type string
}

// Error is a method of the error interface that returns the error message.
//
// Returns:
//   - string: The error message.
func (e *ErrInvalidType[T]) Error() string {
	return fmt.Sprintf("invalid %s type: %T", e.Type, e.Elem)
}

// NewErrInvalidType creates a new ErrInvalidType error.
//
// Parameters:
//   - typeName: The name of the type that was expected.
//   - elem: The element that caused the error.
//
// Returns:
//   - *ErrInvalidType: A pointer to the newly created ErrInvalidType.
func NewErrInvalidType[T any](typeName string, elem T) *ErrInvalidType[T] {
	return &ErrInvalidType[T]{Elem: elem, Type: typeName}
}

// ErrInvalidCharacter represents an error when an invalid character is found.
type ErrInvalidCharacter struct {
	Character rune
}

// Error returns the error message: "invalid character (character)".
//
// Returns:
//   - string: The error message.
func (e *ErrInvalidCharacter) Error() string {
	return fmt.Sprintf("invalid character (%c)", e.Character)
}

// NewErrInvalidCharacter creates a new ErrInvalidCharacter error.
//
// Parameters:
//   - character: The invalid character.
//
// Returns:
//   - *ErrInvalidCharacter: The new ErrInvalidCharacter error.
func NewErrInvalidCharacter(character rune) *ErrInvalidCharacter {
	return &ErrInvalidCharacter{Character: character}
}

// ErrNotComparable represents an error when a value is not comparable.
type ErrNotComparable struct{}

// Error returns the error message: "value is not comparable".
//
// Returns:
//   - string: The error message.
func (e *ErrNotComparable) Error() string {
	return "value is not comparable"
}

// NewErrNotComparable creates a new ErrNotComparable error.
//
// Returns:
//   - *ErrNotComparable: The new ErrNotComparable error.
func NewErrNotComparable() *ErrNotComparable {
	return &ErrNotComparable{}
}

// ErrNilValue represents an error when a value is nil.
type ErrNilValue struct{}

// Error returns the error message: "value must not be nil".
//
// Returns:
//   - string: The error message.
func (e *ErrNilValue) Error() string {
	return "value must not be nil"
}

// NewErrNilValue creates a new ErrNilValue error.
//
// Returns:
//   - *ErrNilValue: The new ErrNilValue error.
func NewErrNilValue() *ErrNilValue {
	return &ErrNilValue{}
}
