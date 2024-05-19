package Strings

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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
		var builder strings.Builder

		builder.WriteString(vals[0])
		builder.WriteString(" and ")
		builder.WriteString(vals[1])

		return builder.String()
	default:
		var builder strings.Builder

		builder.WriteString(strings.Join(vals[0:len(vals)-1], ", "))
		builder.WriteString(", and ")
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
		var builder strings.Builder

		builder.WriteString(vals[0])
		builder.WriteString(" or ")
		builder.WriteString(vals[1])

		return builder.String()
	default:
		var builder strings.Builder

		builder.WriteString(strings.Join(vals[0:len(vals)-1], ", "))
		builder.WriteString(", or ")
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
	var builder strings.Builder

	builder.WriteString(GetOrdinalSuffix(date.Day()))
	builder.WriteRune(' ')
	builder.WriteString(date.Month().String())
	builder.WriteString(", ")
	builder.WriteString(strconv.Itoa(date.Year()))

	return builder.String()
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

// StringsJoiner joins a list of fmt.Stringer values using a separator.
//
// Parameters:
//   - values: The list of fmt.Stringer values to join.
//   - sep: The separator to use when joining the strings.
//
// Returns:
//   - string: The string representation of the values.
func StringsJoiner[T fmt.Stringer](values []T, sep string) string {
	stringValues := make([]string, 0, len(values))

	for _, value := range values {
		stringValues = append(stringValues, value.String())
	}

	return strings.Join(stringValues, sep)
}
