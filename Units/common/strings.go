package common

import (
	"strconv"
	"strings"
)

// QuoteInt returns a quoted string of an integer prefixed and suffixed with
// square brackets.
//
// Parameters:
//   - value: The integer to quote.
//
// Returns:
//   - string: The quoted integer.
func QuoteInt(value int) string {
	var builder strings.Builder

	builder.WriteRune('[')
	builder.WriteString(strconv.Itoa(value))
	builder.WriteRune(']')

	str := builder.String()

	return str
}

// TrimEmpty removes empty strings from a slice of strings.
// Empty spaces at the beginning and end of the strings are also removed from
// the strings.
//
// Parameters:
//   - values: The slice of strings to trim.
//
// Returns:
//   - []string: The slice of strings with empty strings removed.
func TrimEmpty(values []string) []string {
	if len(values) == 0 {
		return values
	}

	var top int

	for i := 0; i < len(values); i++ {
		current_value := values[i]

		str := strings.TrimSpace(current_value)
		if str != "" {
			values[top] = str
			top++
		}
	}

	values = values[:top]

	return values
}

// OrString returns a string that represents a list of strings in a human-readable
// format. The strings are separated by commas and the last string is preceded by
// the word "or".
//
// Empty strings are removed from the list according to the TrimEmpty function.
//
// Parameters:
//   - values: The list of strings to format.
//   - is_negative: Whether the list of strings is negative. (e.g. "a, b, nor c")
//
// Returns:
//   - string: The formatted string.
//
// Examples:
//   - OrString([]string{"a", "b", "c"}) => "a, b, or c"
func OrString(values []string, is_negative bool) string {
	values = TrimEmpty(values)

	var sep string

	if is_negative {
		sep = " nor "
	} else {
		sep = " or "
	}

	var str string

	switch len(values) {
	case 0:
		// Do nothing.
	case 1:
		// a
		str = values[0]
	case 2:
		// a or b
		var builder strings.Builder

		builder.WriteString(values[0])
		builder.WriteString(sep)
		builder.WriteString(values[1])

		str = builder.String()
	default:
		// a, b, or c
		var builder strings.Builder

		builder.WriteString(values[0])

		for i := 1; i < len(values)-1; i++ {
			builder.WriteString(", ")
			builder.WriteString(values[i])
		}
		builder.WriteRune(',')

		builder.WriteString(sep)
		builder.WriteString(values[len(values)-1])

		str = builder.String()
	}

	return str
}

// OrQuoteString returns a string that represents a list of strings in a human-readable
// format. The strings are quoted and separated by commas and the last string is preceded
// by the word "or".
//
// Empty strings are removed from the list according to the TrimEmpty function before
// being quoted.
//
// Parameters:
//   - values: The list of strings to format.
//   - is_negative: Whether the list of strings is negative. (e.g. "a, b, nor c")
//
// Returns:
//   - string: The formatted string.
//
// Examples:
//   - OrQuoteString([]string{"a", "b", "c"}, false) => "\"a\", \"b\", or \"c\""
func OrQuoteString(values []string, is_negative bool) string {
	values = TrimEmpty(values)
	if len(values) == 0 {
		return ""
	}

	for i := 0; i < len(values); i++ {
		current_value := values[i]

		values[i] = strconv.Quote(current_value)
	}

	var sep string

	if is_negative {
		sep = " nor "
	} else {
		sep = " or "
	}

	var str string

	switch len(values) {
	case 1:
		// "a"
		str = values[0]
	case 2:
		// "a" or "b"
		var builder strings.Builder

		builder.WriteString(values[0])
		builder.WriteString(sep)
		builder.WriteString(values[1])

		str = builder.String()
	default:
		// "a", "b", or "c"
		var builder strings.Builder

		builder.WriteString(values[0])

		for i := 1; i < len(values)-1; i++ {
			builder.WriteString(", ")
			builder.WriteString(values[i])
		}
		builder.WriteRune(',')
		builder.WriteString(sep)
		builder.WriteString(values[len(values)-1])

		str = builder.String()
	}

	return str
}
