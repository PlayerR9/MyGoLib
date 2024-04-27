package StringExt

import (
	"slices"
	"strings"
	"testing"
	"unicode/utf8"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// ReplaceSuffix replaces the end of the given string with the provided suffix.
//
// Parameters:
//   - str: The original string.
//   - suffix: The suffix to replace the end of the string.
//
// Returns:
//   - string: The resulting string after replacing the end with the suffix.
//   - error: An error of type *ErrLongerSuffix if the suffix is longer than
//     the string.
//
// Examples:
//
//	const (
//		str    string = "hello world"
//		suffix string = "Bob"
//	)
//
//	result, err := ReplaceSuffix(str, suffix)
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(result) // Output: hello woBob
//	}
func TestReplaceSuffix(t *testing.T) {
	const (
		Str    string = "hello world"
		Suffix string = "Bob"
	)

	result, err := ReplaceSuffix(Str, Suffix)
	if err != nil {
		t.Errorf("expected no error, got %s instead", err.Error())
	}

	if result != "hello woBob" {
		t.Errorf("expected 'hello woBob', got %s instead", result)
	}
}

// FindContentIndexes searches for the positions of opening and closing
// tokens in a slice of strings.
//
// Parameters:
//   - openingToken: The string that marks the beginning of the content.
//   - closingToken: The string that marks the end of the content.
//   - contentTokens: The slice of strings in which to search for the tokens.
//
// Returns:
//   - int: The start index of the content (inclusive).
//   - int: The end index of the content (exclusive).
//   - error: Any error that occurred while searching for the tokens.
//
// Errors:
//   - *ers.ErrInvalidParameter: If the openingToken or closingToken is an
//     empty string.
//   - *ErrTokenNotFound: If the opening or closing token is not found in the
//     content.
//   - *ErrNeverOpened: If the closing token is found without any
//     corresponding opening token.
func TestFindContentIndexes(openingToken, closingToken string, contentTokens []string) (int, int, error) {
	var (
		ContentTokens []string = []string{
			"(", "(", "a", "+", "b", ")", "*", "c", ")", "+", "d",
		}
	)

	if openingToken == "" {
		return 0, 0, ers.NewErrInvalidParameter(
			"openingToken",
			ers.NewErrEmptyString(),
		)
	}

	if closingToken == "" {
		return 0, 0, ers.NewErrInvalidParameter(
			"closingToken",
			ers.NewErrEmptyString(),
		)
	}

	openingTokenIndex := slices.Index(contentTokens, openingToken)
	if openingTokenIndex < 0 {
		return 0, 0, NewErrTokenNotFound(openingToken, OpToken)
	}

	tokenStartIndex := openingTokenIndex + 1

	tokenBalance := 1
	closingTokenIndex := slices.IndexFunc(contentTokens[tokenStartIndex:], func(token string) bool {
		if token == closingToken {
			tokenBalance--
		} else if token == openingToken {
			tokenBalance++
		}

		return tokenBalance == 0
	})

	if closingTokenIndex != -1 {
		return tokenStartIndex, tokenStartIndex + closingTokenIndex + 1, nil
	}

	if tokenBalance < 0 {
		return 0, 0, NewErrNeverOpened(openingToken, closingToken)
	} else if tokenBalance == 1 && closingToken == "\n" {
		return tokenStartIndex, len(contentTokens), nil
	}

	return 0, 0, NewErrTokenNotFound(closingToken, ClToken)
}

// SplitSentenceIntoFields splits the string into fields, where each field is a
// substring separated by one or more whitespace characters.
// The function also handles special characters such as tabs, vertical tabs,
// carriage returns, line feeds, and form feeds.
//
// Parameters:
//   - sentence: The string to split into fields.
//   - indentLevel: The number of spaces that a tab character is replaced with.
//
// Returns:
//   - [][]string: A two-dimensional slice of strings, where each inner slice
//     represents the fields of a line from the input string.
//   - error: An error of type *ers.ErrInvalidRuneAt if an invalid rune is found in
//     the sentence.
//
// Notes:
//   - Negative indentLevel values are converted to positive values.
//   - Empty sentences return an empty slice with no errors.
func TestSplitSentenceIntoFields(sentence string, indentLevel int) ([][]string, error) {
	if sentence == "" {
		return [][]string{}, nil
	}

	if indentLevel < 0 {
		indentLevel *= -1
	}

	lines := make([][]string, 0)
	words := make([]string, 0)

	var builder strings.Builder

	for j := 0; len(sentence) > 0; j++ {
		char, size := utf8.DecodeRuneInString(sentence)
		sentence = sentence[size:]

		if char == utf8.RuneError {
			return nil, ers.NewErrInvalidRuneAt(j, nil)
		}

		switch char {
		case '\t':
			builder.WriteString(strings.Repeat(" ", indentLevel)) // 3 spaces
		case '\v':
			// Do nothing
		case '\r':
			if utf8.RuneCountInString(sentence) > 0 {
				nextRune, size := utf8.DecodeRuneInString(sentence)

				if nextRune == '\n' {
					sentence = sentence[size:]
				}
			}

			fallthrough
		case '\n', '\u0085':
			if builder.Len() != 0 {
				words = append(words, builder.String())
				builder.Reset()
			}

			lines = append(lines, words)
			words = make([]string, 0)
		case ' ':
			if builder.Len() != 0 {
				words = append(words, builder.String())
				builder.Reset()
			}
		case '\u00A0':
			builder.WriteRune(' ')
		case '\f':
			if builder.Len() != 0 {
				words = append(words, builder.String())
				builder.Reset()
			}

			lines = append(lines, words)
			words = make([]string, 0)

			lines = append(lines, []string{string(char)})
		default:
			builder.WriteRune(char)
		}
	}

	if builder.Len() != 0 {
		words = append(words, builder.String())
	}

	if len(words) > 0 {
		lines = append(lines, words)
	}

	return lines, nil
}

// ByteSplitter splits a byte slice into multiple slices based on a separator byte.
// The separator byte is not included in the resulting slices.
//
// If the input slice is empty, the function returns nil.
//
// Parameters:
//   - data: The byte slice to split.
//   - sep: The separator byte.
//
// Returns:
//   - [][]byte: A slice of byte slices.
func TestByteSplitter(data []byte, sep byte) [][]byte {
	if len(data) == 0 {
		return [][]byte{}
	}

	slices := make([][]byte, 0)

	start := 0

	for i := 0; i < len(data); i++ {
		if data[i] == sep {
			slices = append(slices, data[start:i])
			start = i + 1
		}
	}

	slices = append(slices, data[start:])

	return slices
}

// JoinBytes joins multiple byte slices into a single string using a separator byte.
//
// If the input slice is empty, the function returns an empty string.
//
// Parameters:
//   - slices: A slice of byte slices to join.
//   - sep: The separator byte.
//
// Returns:
//   - string: The joined string.
func TestJoinBytes(slices [][]byte, sep byte) string {
	if len(slices) == 0 {
		return ""
	}

	var builder strings.Builder

	builder.Write(slices[0])

	for _, slice := range slices[1:] {
		builder.WriteByte(sep)
		builder.Write(slice)
	}

	return builder.String()
}
