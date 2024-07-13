package Go

import (
	"errors"
	"unicode"
	"unicode/utf8"
)

// IsGenericsID checks if the input string is a valid single upper case letter and returns it as a rune.
//
// Parameters:
//   - str: The input string to check.
//
// Returns:
//   - rune: The valid single upper case letter.
//   - error: An error if the input string is not a valid single upper case letter.
func IsGenericsID(str string) (rune, error) {
	if str == "" {
		err := errors.New("generic type is empty")
		return '\000', err
	}

	size := utf8.RuneCountInString(str)
	if size > 1 {
		err := errors.New("generic type is not a single character")
		return '\000', err
	}

	letter := rune(str[0])

	ok := unicode.IsLetter(letter)
	if !ok {
		return '\000', errors.New("generic type is not a letter")
	}

	ok = unicode.IsUpper(letter)
	if !ok {
		return '\000', errors.New("generic type is not an upper case letter")
	}

	return letter, nil
}
