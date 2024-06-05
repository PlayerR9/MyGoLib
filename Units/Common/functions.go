package Common

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

// TypeOf returns the type of the value as a string.
//
// Parameters:
//   - value: The value to get the type of.
//
// Returns:
//   - string: The type of the value.
func TypeOf(value any) string {
	if value == nil {
		return "no type"
	}

	switch value.(type) {
	case string:
		return "string"
	case []any:
		return "slice"
	case map[any]any:
		return "map"
	default:
		return fmt.Sprintf("%T", value)
	}
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

	switch elem := elem.(type) {
	case int:
		return elem == 0
	case int8:
		return elem == 0
	case int16:
		return elem == 0
	case int32:
		return elem == 0
	case int64:
		return elem == 0
	case uint:
		return elem == 0
	case uint8:
		return elem == 0
	case uint16:
		return elem == 0
	case uint32:
		return elem == 0
	case uint64:
		return elem == 0
	case float32:
		return elem == 0
	case float64:
		return elem == 0
	case bool:
		return !elem
	case string:
		return elem == ""
	case error:
		return elem == nil
	case []any:
		return len(elem) == 0
	case map[any]any:
		return len(elem) == 0
	default:
		reflectValue := reflect.ValueOf(elem)
		return reflectValue.IsZero()
	}
}

// GetOrdinalSuffix returns the ordinal suffix for a given integer.
//
// Parameters:
//
//   - number: The integer for which to get the ordinal suffix.
//
// Returns:
//
//   - string: The ordinal suffix for the number.
//
// For example, for the number 1, the function returns "1st"; for the number 2,
// it returns "2nd"; and so on.
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
