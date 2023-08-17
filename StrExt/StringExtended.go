package StrExt

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

// StringToUTF8 converts a string to a slice of runes.
//
// Parameters:
//   - str: The string to convert.
//
// Returns:
//   - []rune: A slice of runes, each rune representing a character in the string.
//   - error: If the string contains a character that is not a valid UTF-8 character.
func StringToUTF8(str string) ([]rune, error) {
	characters := make([]rune, 0)

	for char_index := 0; len(str) > 0; char_index++ {
		character, size := utf8.DecodeRuneInString(str)

		if character == utf8.RuneError {
			return nil, fmt.Errorf("character at the %s index is not a valid UTF-8 character", PrintOrdinalNumber(char_index))
		}

		characters = append(characters, character)
		str = str[size:]
	}

	return characters, nil
}

// IsNumber returns true if the string is a number, false otherwise.
//
// Parameters:
//   - str: The string to check.
//
// Returns:
//   - bool: True if the string is a number, false otherwise.
func IsNumber(str string) bool {
	_, err := strconv.Atoi(str)

	return err == nil
}

// IsWord returns true if the string is a word, false otherwise.
//
// Parameters:
//   - str: The string to check.
//
// Returns:
//   - bool: True if the string is a word, false otherwise.
func IsWord(str string) bool {
	for len(str) > 0 {
		character, size := utf8.DecodeRuneInString(str)

		if character == utf8.RuneError || !unicode.IsLetter(character) {
			return false
		}

		str = str[size:]
	}

	return true
}

// ExtractContent extracts the content between the opening and closing characters, and returns the content, the remaining content, and an error if one occurred.
// The opening and closing characters are not included in the content. Panics if the opening or closing character is empty.
//
// Parameters:
//   - op_char: The opening character.
//   - cl_char: The closing character.
//   - content: The content to extract from.
//
// Returns:
//   - []string: The content between the opening and closing characters.
//   - []string: The remaining content.
//   - error: An error if one occurred.
//
// Information:
//   - The first element of the content must be the opening character.
func ExtractContent(op_char, cl_char string, content []string) ([]string, []string, error) {
	if len(content) == 0 || content[0] != op_char {
		return nil, nil, fmt.Errorf("no opening character")
	}

	if op_char == "" {
		panic("opening character cannot be empty")
	}

	if cl_char == "" {
		panic("closing character cannot be empty")
	}

	counter := 1

	for i := 1; i < len(content); i++ {
		if content[i] == cl_char {
			counter--
		} else if content[i] == op_char {
			counter++
		}

		if counter == 0 {
			return content[1:i], content[i+1:], nil
		}
	}

	if cl_char != "\n" {
		return nil, nil, fmt.Errorf("no closing character found; expected %s", cl_char)
	}

	return content[1:], nil, nil
}

// ExtractUpTo extracts the content up to the closing character, and returns the content, the remaining content, and an error if one occurred.
// The closing character is not included in the content. Panics if the closing character is empty.
//
// Parameters:
//   - cl_char: The closing character.
//   - content: The content to extract from.
//
// Returns:
//   - []string: The content up to the closing character.
//   - []string: The remaining content.
//   - error: An error if one occurred.
func ExtractUpTo(cl_char string, content []string) ([]string, []string, error) {
	if len(content) == 0 {
		return nil, nil, fmt.Errorf("no opening character")
	}

	if cl_char == "" {
		panic("closing character cannot be empty")
	}

	for i := 0; i < len(content); i++ {
		if content[i] == cl_char {
			return content[:i], content[i+1:], nil
		}
	}

	if cl_char != "\n" {
		return nil, nil, fmt.Errorf("no closing character found; expected %s", cl_char)
	}

	return content[:len(content)-1], nil, nil
}

// PrintOrdinalNumber prints the ordinal number of a position. For example, if the position is 1, it will print 1st.
//
// Parameters:
//   - position: The position to print.
//
// Returns:
//   - string: The ordinal number of the position.
func PrintOrdinalNumber(position int) string {
	// Get the last digit
	last_digit := position % 10

	switch last_digit {
	case 1:
		return fmt.Sprintf("%dst", position)
	case 2:
		return fmt.Sprintf("%dnd", position)
	case 3:
		return fmt.Sprintf("%drd", position)
	default:
		return fmt.Sprintf("%dth", position)
	}
}
