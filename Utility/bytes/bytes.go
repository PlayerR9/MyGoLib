package bytes

import (
	"bytes"
	"unicode/utf8"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// ToUtf8 is a function that converts bytes to runes.
//
// Parameters:
//   - data: The bytes to convert.
//
// Returns:
//   - []rune: The runes.
//   - int: The index of the error rune. -1 if there is no error.
//
// This function also converts '\r\n' to '\n'. Plus, whenever an error occurs, it returns the runes
// decoded so far and the index of the error rune.
func ToUtf8(data []byte) ([]rune, int) {
	if len(data) == 0 {
		return nil, -1
	}

	var chars []rune
	var i int

	for len(data) > 0 {
		c, size := utf8.DecodeRune(data)
		if c == utf8.RuneError {
			return chars, i
		}

		data = data[size:]
		i += size

		if c != '\r' {
			chars = append(chars, c)
			continue
		}

		if len(data) == 0 {
			return chars, i
		}

		c, size = utf8.DecodeRune(data)
		if c == utf8.RuneError {
			return chars, i
		}

		data = data[size:]
		i += size

		if c != '\n' {
			return chars, i
		}

		chars = append(chars, '\n')
	}

	return chars, -1
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

// FirstIndex returns the first index of the target in the tokens.
//
// Parameters:
//   - tokens: The slice of tokens in which to search for the target.
//   - target: The target to search for.
//
// Returns:
//   - int: The index of the target. -1 if the target is not found.
//
// If either tokens or the target are empty, it returns -1.
func FirstIndex(tokens [][]byte, target []byte) int {
	if len(tokens) == 0 || len(target) == 0 {
		return -1
	}

	for i, token := range tokens {
		if bytes.Equal(token, target) {
			return i
		}
	}

	return -1
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
func FindContentIndexes(op_token, cl_token []byte, tokens [][]byte) (result [2]int, err error) {
	result[0] = -1
	result[1] = -1

	if len(cl_token) == 0 {
		err = uc.NewErrInvalidParameter("cl_token", uc.NewErrEmpty(cl_token))
		return
	}

	op_tok_idx := FirstIndex(tokens, op_token)
	if op_tok_idx == -1 {
		err = NewErrTokenNotFound(op_token, true)
		return
	}

	result[0] = op_tok_idx + 1

	balance := 1
	cl_tok_idx := -1

	for i := result[0]; i < len(tokens) && cl_tok_idx == -1; i++ {
		curr_tok := tokens[i]

		if bytes.Equal(curr_tok, cl_token) {
			balance--

			if balance == 0 {
				cl_tok_idx = i
			}
		} else if bytes.Equal(curr_tok, op_token) {
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
	} else if balance != 1 || bytes.Equal(cl_token, []byte("\n")) {
		err = NewErrTokenNotFound(cl_token, false)
		return
	}

	result[1] = len(tokens)
	return
}
