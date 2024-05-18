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
//   - error: An error of type *ErrInvalidUTF8Encoding if the string contains
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

	if !utf8.ValidString(s) {
		return nil, NewErrInvalidUTF8Encoding()
	}

	return []rune(s), nil
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
	if suffix == "" {
		return str, nil
	}

	countStr := utf8.RuneCountInString(str)
	countSuffix := utf8.RuneCountInString(suffix)

	if countStr < countSuffix {
		return "", NewErrLongerSuffix(str, suffix)
	}

	if countStr == countSuffix {
		return suffix, nil
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

	if reason := ue.NewErrEmpty(openingToken).ErrorIf(); reason != nil {
		err = ue.NewErrInvalidParameter("openingToken", reason)
		return
	}

	if reason := ue.NewErrEmpty(closingToken).ErrorIf(); reason != nil {
		err = ue.NewErrInvalidParameter("closingToken", reason)
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
	closingTokenIndex := -1

	for i := result[0]; i < len(contentTokens); i++ {
		if contentTokens[i] == closingToken {
			tokenBalance--

			if tokenBalance == 0 {
				closingTokenIndex = i
				break
			}
		} else if contentTokens[i] == openingToken {
			tokenBalance++
		}
	}

	if closingTokenIndex != -1 {
		result[1] = closingTokenIndex + 1
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
				builder.Reset()
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
			if size != 0 {
				nextRune, size := utf8.DecodeRuneInString(sentence)

				if nextRune == '\n' {
					sentence = sentence[size:]
				}
			}

			fallthrough
		case '\n', '\u0085', ' ', '\f':
			if builder.Len() != 0 {
				words = append(words, builder.String())
				builder.Reset()
			}

			if char != ' ' {
				if len(words) > 0 {
					lines = append(lines, words)
					words = make([]string, 0)
				}

				if char == '\f' {
					lines = append(lines, []string{string(char)})
				}
			}
		case '\u00A0':
			builder.WriteRune(' ')
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

// FitString fits a string to the specified width by adding spaces to the end
// of the string until the width is reached.
//
// Parameters:
//   - width: The width to fit the string to.
//
// Returns:
//   - string: The string with spaces added to the end to fit the width.
//   - error: An error if the width is less than 0.
//
// Behaviors:
//   - If the width is less than the length of the string, the string is
//     truncated to fit the width.
//   - If the width is greater than the length of the string, spaces are added
//     to the end of the string until the width is reached.
func FitString(s string, width int) (string, error) {
	if width < 0 {
		return "", ue.NewErrInvalidParameter(
			"width",
			ue.NewErrGTE(0),
		)
	}

	len := len([]rune(s))

	if width == 0 {
		return "", nil
	} else if len == 0 {
		return strings.Repeat(" ", width), nil
	} else if len == width {
		return s, nil
	}

	if len > width {
		return s[:width], nil
	} else {
		return s + strings.Repeat(" ", width-len), nil
	}
}
