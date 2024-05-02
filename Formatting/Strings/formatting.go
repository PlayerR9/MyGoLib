package Strings

import (
	"fmt"
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
		return fmt.Sprintf("%s and %s", vals[0], vals[1])
	default:
		var builder strings.Builder

		builder.WriteString(vals[0])

		for i := 1; i < len(vals)-1; i++ {
			builder.WriteRune(',')
			builder.WriteRune(' ')
			builder.WriteString(vals[i])
		}

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

		builder.WriteString(vals[0])

		for i := 1; i < len(vals)-1; i++ {
			builder.WriteRune(',')
			builder.WriteRune(' ')
			builder.WriteString(vals[i])
		}

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
