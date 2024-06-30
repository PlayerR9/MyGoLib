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
