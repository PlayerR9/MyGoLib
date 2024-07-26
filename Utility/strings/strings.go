package strings

import (
	"slices"
	"strconv"
	"strings"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// filter_equals returns the indices of the other in the data.
//
// Parameters:
//   - indices: The indices.
//   - data: The data.
//   - other: The other value.
//   - offset: The offset to start the search from.
//
// Returns:
//   - []int: The indices.
func filter_equals(indices []int, data []string, other string, offset int) []int {
	var top int

	for i := 0; i < len(indices); i++ {
		idx := indices[i]

		if data[idx+offset] == other {
			indices[top] = idx
			top++
		}
	}

	indices = indices[:top]

	return indices
}

// Indices returns the indices of the separator in the data.
//
// Parameters:
//   - data: The data.
//   - sep: The separator.
//   - exclude_sep: Whether the separator is inclusive. If set to true, the indices will point to the character right after the
//     separator. Otherwise, the indices will point to the character right before the separator.
//
// Returns:
//   - []int: The indices.
func IndicesOf(data []string, sep []string, exclude_sep bool) []int {
	if len(data) == 0 || len(sep) == 0 {
		return nil
	}

	var indices []int

	for i := 0; i < len(data)-len(sep); i++ {
		if data[i] == sep[0] {
			indices = append(indices, i)
		}
	}

	if len(indices) == 0 {
		return nil
	}

	for i := 1; i < len(sep); i++ {
		other := sep[i]

		indices = filter_equals(indices, data, other, i)

		if len(indices) == 0 {
			return nil
		}
	}

	if exclude_sep {
		for i := 0; i < len(indices); i++ {
			indices[i] += len(sep)
		}
	}

	return indices
}

// FindContentIndexes searches for the positions of opening and closing
// tokens in a slice of strings.
//
// Parameters:
//   - op_token: The string that marks the beginning of the content.
//   - cl_token: The string that marks the end of the content.
//   - tokens: The slice of strings in which to search for the tokens.
//
// Returns:
//   - result: An array of two integers representing the start and end indexes
//     of the content.
//   - err: Any error that occurred while searching for the tokens.
//
// Errors:
//   - *uc.ErrInvalidParameter: If the closingToken is an empty string.
//   - *ErrTokenNotFound: If the opening or closing token is not found in the
//     content.
//   - *ErrNeverOpened: If the closing token is found without any
//     corresponding opening token.
//
// Behaviors:
//   - The first index of the content is inclusive, while the second index is
//     exclusive.
//   - This function returns a partial result when errors occur. ([-1, -1] if
//     errors occur before finding the opening token, [index, 0] if the opening
//     token is found but the closing token is not found.
func FindContentIndexes(op_token, cl_token string, tokens []string) (result [2]int, err error) {
	result[0] = -1
	result[1] = -1

	if cl_token == "" {
		err = uc.NewErrInvalidParameter("cl_token", uc.NewErrEmpty(cl_token))
		return
	}

	op_tok_idx := slices.Index(tokens, op_token)
	if op_tok_idx < 0 {
		err = NewErrTokenNotFound(op_token, true)
		return
	} else {
		result[0] = op_tok_idx + 1
	}

	balance := 1
	cl_tok_idx := -1

	for i := result[0]; i < len(tokens) && cl_tok_idx == -1; i++ {
		curr_tok := tokens[i]

		if curr_tok == cl_token {
			balance--

			if balance == 0 {
				cl_tok_idx = i
			}
		} else if curr_tok == op_token {
			balance++
		}
	}

	if cl_tok_idx != -1 {
		result[1] = cl_tok_idx + 1
		return
	}

	if balance < 0 {
		err = NewErrNeverOpened(op_token, cl_token)
		return
	} else if balance != 1 || cl_token != "\n" {
		err = NewErrTokenNotFound(cl_token, false)
		return
	}

	result[1] = len(tokens)
	return
}

// AndString is a function that returns a string representation of a slice
// of strings.
//
// Parameters:
//   - values: The values to convert to a string.
//   - quote: Whether to quote the values.
//
// Returns:
//   - string: The string representation of the values.
func AndString(values []string, quote bool) string {
	values = TrimEmpty(values)
	if len(values) == 0 {
		return ""
	}

	if len(values) == 1 {
		if !quote {
			return values[0]
		} else {
			return strconv.Quote(values[0])
		}
	}

	var elems []string

	if quote {
		for i := 0; i < len(values); i++ {
			elems = append(elems, strconv.Quote(values[i]))
		}
	} else {
		elems = values
	}

	var builder strings.Builder

	builder.WriteString(elems[0])

	if len(elems) > 2 {
		builder.WriteString(strings.Join(elems[1:len(elems)-1], ", "))
		builder.WriteRune(',')
	}

	builder.WriteString(" and ")
	builder.WriteString(elems[len(elems)-1])

	return builder.String()
}

// EitherOrString is a function that returns a string representation of a slice
// of strings.
//
// Parameters:
//   - values: The values to convert to a string.
//   - quote: True if the values should be quoted, false otherwise.
//
// Returns:
//   - string: The string representation.
//
// Example:
//
//	EitherOrString([]string{"a", "b", "c"}, false) // "a, b or c"
func EitherOrString(values []string, quote bool) string {
	if len(values) == 0 {
		return ""
	}

	if len(values) == 1 {
		if !quote {
			return values[0]
		} else {
			return strconv.Quote(values[0])
		}
	}

	var elems []string

	if quote {
		for _, v := range values {
			elems = append(elems, strconv.Quote(v))
		}
	} else {
		elems = values
	}

	var builder strings.Builder

	builder.WriteString("either ")
	builder.WriteString(elems[0])

	if len(values) > 2 {
		builder.WriteString(strings.Join(elems[1:len(elems)-1], ", "))
		builder.WriteRune(',')
	}

	builder.WriteString(" or ")
	builder.WriteString(elems[len(elems)-1])

	return builder.String()
}

// OrString is a function that returns a string representation of a slice of
// strings.
//
// Parameters:
//   - values: The values to convert to a string.
//   - quote: True if the values should be quoted, false otherwise.
//   - is_negative: True if the string should use "nor" instead of "or", false
//     otherwise.
//
// Returns:
//   - string: The string representation.
//
// Example:
//
//	OrString([]string{"a", "b", "c"}, false, true) // "a, b, nor c"
func OrString(values []string, quote, is_negative bool) string {
	values = TrimEmpty(values)
	if len(values) == 0 {
		return ""
	}

	if len(values) == 1 {
		if !quote {
			return values[0]
		} else {
			return strconv.Quote(values[0])
		}
	}

	var sep string

	if is_negative {
		sep = " nor "
	} else {
		sep = " or "
	}

	var elems []string

	if quote {
		for _, v := range values {
			elems = append(elems, strconv.Quote(v))
		}
	} else {
		elems = values
	}

	var builder strings.Builder

	builder.WriteString(elems[0])

	if len(values) > 2 {
		builder.WriteString(strings.Join(elems[1:len(elems)-1], ", "))
		builder.WriteRune(',')
	}

	builder.WriteString(sep)
	builder.WriteString(elems[len(elems)-1])

	return builder.String()
}

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

	return builder.String()
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
