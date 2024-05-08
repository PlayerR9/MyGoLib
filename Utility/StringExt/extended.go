// Package StringExt provides a set of functions that extend the functionality of
// the built-in string type.
package StringExt

import (
	"crypto/rand"
	"encoding/hex"
	"slices"
	"strings"
	"unicode/utf8"

	ue "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// ToUTF8Runes converts a string to a slice of runes.
//
// Parameters:
//   - s: The string to convert.
//
// Returns:
//   - []rune: The slice of runes
//   - error: An error of type *ue,ErrAtIndex if the string contains
//     invalid UTF-8 encoding.
//
// Behaviors:
//   - An empty string returns a nil slice with no errors.
//   - The function stops at the first invalid UTF-8 encoding; returning an
//     error and the runes found up to that point.
func ToUTF8Runes(s string) ([]rune, error) {
	if s == "" {
		return nil, nil
	}

	solution := make([]rune, 0)

	for i := 0; len(s) > 0; i++ {
		r, size := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError {
			return solution, ue.NewErrAt(i, NewErrInvalidUTF8Encoding())
		}

		solution = append(solution, r)
		s = s[size:]
	}

	return solution, nil
}

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
func ReplaceSuffix(str, suffix string) (string, error) {
	countStr := utf8.RuneCountInString(str)
	countSuffix := utf8.RuneCountInString(suffix)

	if countStr < countSuffix {
		return "", NewErrLongerSuffix(str, suffix)
	}

	if countStr == countSuffix {
		return suffix, nil
	}

	if countSuffix == 0 {
		return str, nil
	}

	return str[:countStr-countSuffix] + suffix, nil
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
//   - result: An array of two integers representing the start and end indexes
//     of the content.
//   - err: Any error that occurred while searching for the tokens.
//
// Errors:
//   - *ue.ErrInvalidParameter: If the openingToken or closingToken is an
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
func FindContentIndexes(openingToken, closingToken string, contentTokens []string) (result [2]int, err error) {
	result[0] = -1
	result[1] = -1

	if openingToken == "" {
		err = ue.NewErrInvalidParameter(
			"openingToken",
			ue.NewErrEmptyString(),
		)

		return
	} else if closingToken == "" {
		err = ue.NewErrInvalidParameter(
			"closingToken",
			ue.NewErrEmptyString(),
		)

		return
	}

	openingTokenIndex := slices.Index(contentTokens, openingToken)
	if openingTokenIndex < 0 {
		err = NewErrTokenNotFound(openingToken, OpToken)
		return
	} else {
		result[0] = openingTokenIndex + 1
	}

	tokenBalance := 1
	tokenBalanceFunc := func(token string) bool {
		if token == closingToken {
			tokenBalance--
		} else if token == openingToken {
			tokenBalance++
		}

		return tokenBalance == 0
	}

	result[1] = slices.IndexFunc(contentTokens[result[0]:], tokenBalanceFunc) + 1
	if result[1] != 0 {
		return
	}

	if tokenBalance < 0 {
		err = NewErrNeverOpened(openingToken, closingToken)
		return
	} else if tokenBalance != 1 || closingToken != "\n" {
		err = NewErrTokenNotFound(closingToken, ClToken)
		return
	}

	result[1] = len(contentTokens)
	return
}

// SplitSentenceIntoFields splits the string into fields, where each field is a
// substring separated by one or more whitespace charactue.
//
// Parameters:
//   - sentence: The string to split into fields.
//   - indentLevel: The number of spaces that a tab character is replaced with.
//
// Returns:
//   - [][]string: A two-dimensional slice of strings, where each inner slice
//     represents the fields of a line from the input string.
//   - error: An error of type *ue.ErrInvalidRuneAt if an invalid rune is found in
//     the sentence.
//
// Behaviors:
//   - Negative indentLevel values are converted to positive values.
//   - Empty sentences return a nil slice with no errors.
//   - The function handles the following whitespace characters: space, tab,
//     vertical tab, carriage return, line feed, and form feed.
//   - The function returns a partial result if an invalid rune is found where
//     the result are the fields found up to that point.
func AdvancedFieldsSplitter(sentence string, indentLevel int) ([][]string, error) {
	if sentence == "" {
		return nil, nil
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
			if builder.Len() != 0 {
				words = append(words, builder.String())
			}

			if len(words) > 0 {
				lines = append(lines, words)
			}

			return lines, ue.NewErrAt(j, ue.NewErrInvalidRune(nil))
		}

		switch char {
		case '\t':
			// Replace tabs with N spaces
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

// GenerateID generates a random ID of the specified size (in bytes).
//
// Parameters:
//   - size: The size of the ID to generate (in bytes).
//
// Returns:
//   - string: The generated ID.
//   - error: An error if the ID cannot be generated.
//
// Errors:
//   - *ue.ErrInvalidParameter: If the size is less than 1.
//   - Any error returned by the rand.Read function.
//
// Behaviors:
//   - The function uses the crypto/rand package to generate a random ID of the
//     specified size.
//   - The ID is returned as a hexadecimal string.
func GenerateID(size int) (string, error) {
	if size < 1 {
		return "", ue.NewErrInvalidParameter(
			"size",
			ue.NewErrGT(0),
		)
	}

	b := make([]byte, size) // 128 bits

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
