package runes

import (
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// BytesToUtf8 is a function that converts bytes to runes.
//
// Parameters:
//   - data: The bytes to convert.
//
// Returns:
//   - []rune: The runes.
//   - error: An error of type *ErrInvalidUTF8Encoding if the bytes are not
//     valid UTF-8.
//
// This function also converts '\r\n' to '\n'. Plus, whenever an error occurs, it returns the runes
// decoded so far and the index of the error rune.
func BytesToUtf8(data []byte) ([]rune, error) {
	if len(data) == 0 {
		return nil, nil
	}

	var chars []rune
	var i int

	for len(data) > 0 {
		c, size := utf8.DecodeRune(data)
		if c == utf8.RuneError {
			return chars, NewErrInvalidUTF8Encoding(i)
		}

		data = data[size:]
		i += size

		if c != '\r' {
			chars = append(chars, c)
			continue
		}

		if len(data) == 0 {
			return chars, NewErrInvalidUTF8Encoding(i)
		}

		c, size = utf8.DecodeRune(data)
		if c == utf8.RuneError {
			return chars, NewErrInvalidUTF8Encoding(i)
		}

		data = data[size:]
		i += size

		if c != '\n' {
			return chars, NewErrInvalidUTF8Encoding(i)
		}

		chars = append(chars, '\n')
	}

	return chars, nil
}

// StringToUtf8 converts a string to a slice of runes.
//
// Parameters:
//   - str: The string to convert.
//
// Returns:
//   - runes: The slice of runes.
//   - error: An error of type *ErrInvalidUTF8Encoding if the string is not
//     valid UTF-8.
//
// Behaviors:
//   - An empty string returns a nil slice with no errors.
//   - The function stops at the first invalid UTF-8 encoding; returning an
//     error and the runes found up to that point.
//   - The function converts '\r\n' to '\n'.
func StringToUtf8(str string) ([]rune, error) {
	if str == "" {
		return nil, nil
	}

	var chars []rune
	var i int

	for len(str) > 0 {
		c, size := utf8.DecodeRuneInString(str)
		if c == utf8.RuneError {
			return chars, NewErrInvalidUTF8Encoding(i)
		}

		str = str[size:]
		i += size

		if c != '\r' {
			chars = append(chars, c)
			continue
		}

		if len(str) == 0 {
			return chars, NewErrInvalidUTF8Encoding(i)
		}

		c, size = utf8.DecodeRuneInString(str)
		if c == utf8.RuneError {
			return chars, NewErrInvalidUTF8Encoding(i)
		}

		str = str[size:]
		i += size

		if c != '\n' {
			return chars, NewErrInvalidUTF8Encoding(i)
		}

		chars = append(chars, '\n')
	}

	return chars, nil
}

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
func filter_equals(indices []int, data []byte, other byte, offset int) []int {
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
func IndicesOf(data []byte, sep []byte, exclude_sep bool) []int {
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
//   - *uc.ErrInvalidParameter: If the openingToken or closingToken is an
//     empty string.
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
func FindContentIndexes(op_token, cl_token rune, tokens []rune) (result [2]int, err error) {
	result[0] = -1
	result[1] = -1

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
	} else if balance != 1 || cl_token != '\n' {
		err = NewErrTokenNotFound(cl_token, false)
		return
	}

	result[1] = len(tokens)
	return
}

// AndString is a function that returns a string representation of a slice
// of runes.
//
// Parameters:
//   - values: The values to convert to a string.
//   - quote: Whether to quote the values.
//
// Returns:
//   - string: The string representation of the values.
func AndString(values []rune, quote bool) string {
	if len(values) == 0 {
		return ""
	}

	if len(values) == 1 {
		if !quote {
			return string(values[0])
		} else {
			return strconv.QuoteRune(values[0])
		}
	}

	elems := make([]string, 0, len(values))

	if quote {
		for i := 0; i < len(values); i++ {
			elems = append(elems, strconv.QuoteRune(values[i]))
		}
	} else {
		for i := 0; i < len(values); i++ {
			elems = append(elems, string(values[i]))
		}
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
// of runes.
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
//	EitherOrString([]string{'a', 'b', 'c'}, false) // "a, b or c"
func EitherOrString(values []rune, quote bool) string {
	if len(values) == 0 {
		return ""
	}

	if len(values) == 1 {
		if !quote {
			return string(values[0])
		} else {
			return strconv.QuoteRune(values[0])
		}
	}

	elems := make([]string, 0, len(values))

	if quote {
		for _, v := range values {
			elems = append(elems, strconv.QuoteRune(v))
		}
	} else {
		for _, v := range values {
			elems = append(elems, string(v))
		}
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
// runes.
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
//	OrString([]string{'a', 'b', 'c'}, false, true) // "a, b, nor c"
func OrString(values []rune, quote, is_negative bool) string {
	if len(values) == 0 {
		return ""
	}

	if len(values) == 1 {
		if !quote {
			return string(values[0])
		} else {
			return strconv.QuoteRune(values[0])
		}
	}

	var sep string

	if is_negative {
		sep = " nor "
	} else {
		sep = " or "
	}

	elems := make([]string, 0, len(values))

	if quote {
		for _, v := range values {
			elems = append(elems, strconv.QuoteRune(v))
		}
	} else {
		for _, v := range values {
			elems = append(elems, string(v))
		}
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

// LevenshteinTable is a table of words for the Levenshtein distance.
type LavenshteinTable struct {
	// words is the list of words.
	word_list [][]rune

	// word_length_list is the list of word lengths.
	word_length_list []int
}

// NewLevenshteinTable creates a new Levenshtein table
// with the given words.
//
// Parameters:
//   - words: The words to add to the table.
//
// Returns:
//   - *LevenshteinTable: The new Levenshtein table.
//   - error: An error if any of the words cannot be added to the table.
//
// Errors:
//   - *common.ErrAt: Whenever a word is not valid UTF-8.
//
// It is the same as creating an empty table and then adding the words to it.
func NewLevenshteinTable(words ...string) (*LavenshteinTable, error) {
	lt := &LavenshteinTable{}

	for i, word := range words {
		err := lt.AddWord(word)
		if err != nil {
			return nil, uc.NewErrAt(i+1, word, err)
		}
	}

	return lt, nil
}

// AddWord adds a word to the table.
//
// Parameters:
//   - word: The word to add.
//
// Returns:
//   - error: An error of type *ErrInvalidUTF8Encoding if the word is not
//     valid UTF-8.
func (lt *LavenshteinTable) AddWord(word string) error {
	chars, err := StringToUtf8(word)
	if err != nil {
		return err
	}

	lt.word_list = append(lt.word_list, chars)
	lt.word_length_list = append(lt.word_length_list, len(chars))

	return nil
}

// GetClosest gets the closest word to a target.
//
// Parameters:
//   - target: The target.
//
// Returns:
//   - string: The closest word.
//   - error: The error if any occurs.
//
// Errors:
//   - *common.ErrInvalidParameter: If the target is empty.
//   - *ErrNoClosestWordFound: If no closest word is found.
func (lt *LavenshteinTable) GetClosest(target []rune) (string, error) {
	if len(target) == 0 {
		return "", uc.NewErrInvalidParameter("target", uc.NewErrEmpty(target))
	}

	target_len := len(target)

	closest_idx := -1
	var min int

	for i, word := range lt.word_list {
		d := levenshtein_distance(target, target_len, word, lt.word_length_list[i])

		if closest_idx == -1 || d < min {
			min = d
			closest_idx = i
		}
	}

	if closest_idx == -1 {
		return "", NewErrNoClosestWordFound()
	}

	word := lt.word_list[closest_idx]

	return string(word), nil
}

// levenshteinDistance calculates the Levenshtein distance between two strings.
//
// Parameters:
//   - target: The target.
//   - target_len: The target length.
//   - other: The other.
//   - other_len: The other length.
//
// Returns:
//   - int: The Levenshtein distance.
func levenshtein_distance(target []rune, target_len int, other []rune, other_len int) int {
	matrix := make([][]int, 0, target_len+1)

	for i := 0; i <= target_len; i++ {
		row := make([]int, other_len+1)

		matrix = append(matrix, row)
	}

	// Initialize the matrix
	for i := 0; i <= target_len; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= other_len; j++ {
		matrix[0][j] = j
	}

	// Compute the distances
	for i := 1; i <= target_len; i++ {
		for j := 1; j <= other_len; j++ {
			if target[i-1] == other[j-1] {
				matrix[i][j] = matrix[i-1][j-1] // No operation needed
			} else {
				deletion := matrix[i-1][j] + 1
				insertion := matrix[i][j-1] + 1
				substitution := matrix[i-1][j-1] + 1

				min_first := uc.Min(deletion, insertion)
				min_second := uc.Min(min_first, substitution)
				matrix[i][j] = min_second
			}
		}
	}

	d := matrix[target_len][other_len]

	return d
}
