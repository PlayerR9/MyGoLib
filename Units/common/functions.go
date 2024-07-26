package common

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// DoFunc is a generic type that represents a function that takes a value
// and does something with it.
//
// Parameters:
//   - T: The type of the value.
type DoFunc[T any] func(T)

// DualDoFunc is a generic type that represents a function that takes two
// values and does something with them.
//
// Parameters:
//   - T: The type of the first value.
//   - U: The type of the second value.
type DualDoFunc[T any, U any] func(T, U)

// EvalOneFunc is a function that evaluates one element.
//
// Parameters:
//   - elem: The element to evaluate.
//
// Returns:
//   - R: The result of the evaluation.
//   - error: An error if the evaluation failed.
type EvalOneFunc[E, R any] func(elem E) (R, error)

// EvalManyFunc is a function that evaluates many elements.
//
// Parameters:
//   - elem: The element to evaluate.
//
// Returns:
//   - []R: The results of the evaluation.
//   - error: An error if the evaluation failed.
type EvalManyFunc[E, R any] func(elem E) ([]R, error)

// MainFunc is a function type that takes no parameters and returns an error.
// It is used to represent things such as the main function of a program.
//
// Returns:
//   - error: An error if the function failed.
type MainFunc func() error

// Routine is a function type used to represent a go routine.
type RoutineFunc func()

// ErrorIfFunc is a function type that takes an element and returns an error
// if the element is invalid.
//
// Parameters:
//   - elem: The element to check.
//
// Returns:
//   - error: An error if the element is invalid.
type ErrorIfFunc[T any] func(elem T) error

// TypeOf returns the type of the value as a string.
//
// Parameters:
//   - value: The value to get the type of.
//
// Returns:
//   - string: The type of the value.
func TypeOf(value any) string {
	if value == nil {
		return "nil"
	}

	return reflect.TypeOf(value).String()
}

// IsEmpty returns true if the element is empty.
//
// Parameters:
//   - elem: The element to check.
//
// Returns:
//   - bool: True if the element is empty, false otherwise.
func IsEmpty(elem any) bool {
	if elem == nil {
		return true
	}

	value := reflect.ValueOf(elem)
	return value.IsZero()
}

// GetOrdinalSuffix returns the ordinal suffix for a given integer.
//
// Parameters:
//   - number: The integer for which to get the ordinal suffix. Negative
//     numbers are treated as positive.
//
// Returns:
//   - string: The ordinal suffix for the number.
//
// Example:
//   - GetOrdinalSuffix(1) returns "1st"
//   - GetOrdinalSuffix(2) returns "2nd"
func GetOrdinalSuffix(number int) string {
	var builder strings.Builder

	builder.WriteString(strconv.Itoa(number))

	if number < 0 {
		number = -number
	}

	lastTwoDigits := number % 100
	lastDigit := lastTwoDigits % 10

	if lastTwoDigits >= 11 && lastTwoDigits <= 13 {
		builder.WriteString("th")
	} else {
		switch lastDigit {
		case 1:
			builder.WriteString("st")
		case 2:
			builder.WriteString("nd")
		case 3:
			builder.WriteString("rd")
		default:
			builder.WriteString("th")
		}
	}

	return builder.String()
}

// GoStringOf returns a string representation of the element.
//
// Parameters:
//   - elem: The element to get the string representation of.
//
// Returns:
//   - string: The string representation of the element.
//
// Behaviors:
//   - If the element is nil, the function returns "nil".
//   - If the element implements the fmt.GoStringer interface, the function
//     returns the result of the GoString method.
//   - If the element implements the fmt.Stringer interface, the function
//     returns the result of the String method.
//   - If the element is a string, the function returns the string enclosed in
//     double quotes.
//   - If the element is an error, the function returns the error message
//     enclosed in double quotes.
//   - Otherwise, the function returns the result of the %#v format specifier.
func GoStringOf(elem any) string {
	if elem == nil {
		return "nil"
	}

	switch elem := elem.(type) {
	case fmt.GoStringer:
		return elem.GoString()
	case fmt.Stringer:
		return strconv.Quote(elem.String())
	case string:
		return strconv.Quote(elem)
	case error:
		return strconv.Quote(elem.Error())
	default:
		return fmt.Sprintf("%#v", elem)
	}
}

// Min is a function that takes two parameters, a and b, of any type T
// according to the uc.CompareOf function and returns the smaller of the two values.
//
// Parameters:
//   - a, b: The two values to compare.
//
// Return:
//   - T: The smaller of the two values.
func Min[T Comparable](a, b T) T {
	if a < b {
		return a
	} else {
		return b
	}
}

// Max is a function that takes two parameters, a and b, of any type T
// according to the uc.CompareOf function and returns the larger of the two values.
//
// Parameters:
//   - a, b: The two values to compare.
//
// Return:
//   - T: The larger of the two values.
func Max[T Comparable](a, b T) T {
	if a > b {
		return a
	} else {
		return b
	}
}
