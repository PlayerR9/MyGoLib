package Strings

import (
	"fmt"
	"strings"
	"time"

	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
)

// AndString concatenates a list of strings using commas and the word "and" before the last string.
//
// Parameters:
//
//   - vals: The list of strings to concatenate.
//
// Returns:
//
//   - string: The concatenated string.
func AndString(vals ...string) string {
	switch len(vals) {
	case 0:
		return ""
	case 1:
		return vals[0]
	case 2:
		return fmt.Sprintf("%s and %s", vals[0], vals[1])
	default:
		var builder strings.Builder

		builder.WriteString(strings.Join(vals[0:len(vals)-1], ", "))

		builder.WriteRune(',')
		builder.WriteRune(' ')
		builder.WriteString("and")
		builder.WriteRune(' ')
		builder.WriteString(vals[len(vals)-1])

		return builder.String()
	}
}

// OrString concatenates a list of strings using commas and the word "or" before the last string.
//
// Parameters:
//
//   - vals: The list of strings to concatenate.
//
// Returns:
//
//   - string: The concatenated string.
func OrString(vals ...string) string {
	switch len(vals) {
	case 0:
		return ""
	case 1:
		return vals[0]
	case 2:
		return fmt.Sprintf("%s or %s", vals[0], vals[1])
	default:
		var builder strings.Builder

		builder.WriteString(strings.Join(vals[0:len(vals)-1], ", "))

		builder.WriteRune(',')
		builder.WriteRune(' ')
		builder.WriteString("or")
		builder.WriteRune(' ')
		builder.WriteString(vals[len(vals)-1])

		return builder.String()
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
	if number < 0 {
		return fmt.Sprintf("%dth", number)
	}

	lastTwoDigits := number % 100
	lastDigit := lastTwoDigits % 10

	if lastTwoDigits >= 11 && lastTwoDigits <= 13 {
		return fmt.Sprintf("%dth", number)
	}

	if lastDigit == 0 || lastDigit > 3 {
		return fmt.Sprintf("%dth", number)
	}

	switch lastDigit {
	case 1:
		return fmt.Sprintf("%dst", number)
	case 2:
		return fmt.Sprintf("%dnd", number)
	case 3:
		return fmt.Sprintf("%drd", number)
	}

	return ""
}

// DateStringer prints the date in the format "1st January, 2006".
//
// Parameters:
//
//   - date: The date to print.
//
// Returns:
//
//   - string: The date in the format "1st January, 2006".
func DateStringer(date time.Time) string {
	return fmt.Sprintf("%s %s, %d",
		GetOrdinalSuffix(date.Day()),
		date.Month().String(),
		date.Year(),
	)
}

// TimeStringer prints the time in the format "3:04 PM".
//
// Parameters:
//
//   - time: The time to print.
//
// Returns:
//
//   - string: The time in the format "3:04 PM".
func TimeStringer(time time.Time) string {
	return time.Format("3:04 PM")
}

// StringOf converts any type to a string.
//
// The function converts the input element to a string using the following rules:
//   - If the element is a string, it is returned as is.
//   - If the element implements the fmt.Stringer interface, the String method is called.
//   - If the element is an error, the Error method is called.
//   - If the element is a byte slice, it is converted to a string.
//   - If the element is a rune slice, it is converted to a string.
//   - If the element implements the ffs.FStringer interface, the FString function is called.
//   - Otherwise, the element is converted to a string using the fmt.Sprintf function.
//
// Parameters:
//   - elem: The element to convert to a string.
//
// Returns:
//   - string: The string representation of the element.
func StringOf(elem any) string {
	if elem == nil {
		return ""
	}

	switch elem := elem.(type) {
	case string:
		return elem
	case fmt.Stringer:
		return elem.String()
	case error:
		return elem.Error()
	case []byte:
		return string(elem)
	case []rune:
		return string(elem)
	case ffs.FStringer:
		return ffs.FString(elem)
	default:
		return fmt.Sprintf("%v", elem)
	}
}

// StringOfFunc is a function type that converts a value of type T to a string.
//
// Parameters:
//   - elem: The element to convert to a string.
//
// Returns:
//   - string: The string representation of the element.
type StringOfFunc[T any] func(elem T) string

var (
	// StringOfStringer converts a value of type fmt.Stringer to a string.
	StringOfStringer StringOfFunc[fmt.Stringer] = func(elem fmt.Stringer) string {
		return elem.String()
	}

	// StringOfError converts a value of type error to a string.
	StringOfError StringOfFunc[error] = func(elem error) string {
		return elem.Error()
	}

	// StringOfString converts a value of type string to a string.
	StringOfByteSlice StringOfFunc[[]byte] = func(elem []byte) string {
		return string(elem)
	}

	// StringOfString converts a value of type string to a string.
	StringOfRuneSlice StringOfFunc[[]rune] = func(elem []rune) string {
		return string(elem)
	}

	// StringOfFStringer converts a value of type ffs.FStringer to a string.
	StringOfFStringer StringOfFunc[ffs.FStringer] = func(elem ffs.FStringer) string {
		return ffs.FString(elem)
	}

	// StringOfAny converts a value of any type to a string.
	StringOfAny StringOfFunc[any] = func(elem any) string {
		return fmt.Sprintf("%v", elem)
	}
)

// StringsOf converts a list of elements to a string using a separator.
//
// The function converts the input elements to strings using the StringOf function.
// The strings are then joined using the separator.
//
// Parameters:
//   - sep: The separator to use when joining the strings.
//   - f: A function that converts an element to a string.
//   - elems: The elements to convert to strings.
//
// Returns:
//   - string: The string representation of the elements.
func StringsOf[T any](sep string, f StringOfFunc[T], elems ...T) string {
	values := make([]string, 0, len(elems))

	for _, elem := range elems {
		values = append(values, f(elem))
	}

	return strings.Join(values, sep)
}
