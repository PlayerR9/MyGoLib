package StrExt

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

/*
	StringToUTF8 converts a string to a slice of runes.

Parameters:

	str: The string to convert.

Returns:

	A slice of runes, each rune representing a character in the string.
	Error: If the string contains a character that is not a valid UTF-8 character.
*/
func StringToUTF8(str string) ([]rune, error) {
	characters := make([]rune, 0)

	for char_index := 0; len(str) > 0; char_index++ {
		character, size := utf8.DecodeRuneInString(str)

		if character == utf8.RuneError {
			return nil, fmt.Errorf("character at index %d is not a valid UTF-8 character", char_index)
		}

		characters = append(characters, character)
		str = str[size:]
	}

	return characters, nil
}

/*
	IsNumber returns true if the string is a number, false otherwise.

Parameters:

	str: The string to check.

Returns:

	true if the string is a number, false otherwise.
*/
func IsNumber(str string) bool {
	_, err := strconv.Atoi(str)

	return err == nil
}

/*
	IsWord returns true if the string is a word, false otherwise.

Parameters:

	str: The string to check.

Returns:

	true if the string is a word, false otherwise.
*/
func IsWord(str string) bool {
	characters, err := StringToUTF8(str)
	if err != nil {
		return false
	}

	for _, character := range characters {
		if !unicode.IsLetter(character) {
			return false
		}
	}

	return true
}

/*
func DefaultTokenization(str string) ([]string, error) {
	characters, err := StringToUTF8(str)
	if err != nil {
		return nil, fmt.Errorf("could not tokenize string: %s", err)
	}

	tokens := make([]string, 0)

	current_token := ""
	var next_character rune

	for i, character := range characters {
		if i < len(characters) - 1 {
			next_character = characters[i + 1]
		}else{
			next_character = '\n'
		}

		current_token += string(character)

		if unicode.IsLetter(character) {
			if unicode.IsLetter(next_character) {
				continue
			}else if unicode.IsDigit(next_character) {
				continue
			}else{
				tokens = append(tokens, current_token)
				current_token = ""
			}
		}else if unicode.IsDigit(character) {
			if unicode.IsDigit(next_character) {
				continue
			}else{
				tokens = append(tokens, current_token)
				current_token = ""
			}
		} else {
			tokens = append(tokens, current_token)
			current_token = ""
		}
	}

	if current_token != "" {
		tokens = append(tokens, current_token)
	}

	return tokens, nil

}*/

// ExtractContent extracts the content between the opening and closing characters, and returns the content, the remaining content, and an error if one occurred.
// The opening and closing characters are not included in the content.
//
// Parameters:
//
//	op_char: The opening character.
//	cl_char: The closing character.
//	content: The content to extract from.
//
// Returns:
//
//	[]string: The content between the opening and closing characters.
//	[]string: The remaining content.
//	error: An error if one occurred.
func ExtractContent(op_char, cl_char string, content []string) ([]string, []string, error) {
	if len(content) == 0 || content[0] != op_char {
		return nil, nil, fmt.Errorf("no opening character")
	}

	if op_char == "" {
		return nil, nil, fmt.Errorf("opening character cannot be empty")
	}

	if cl_char == "" {
		return nil, nil, fmt.Errorf("closing character cannot be empty")
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

func ExtractUpTo(cl_char string, content []string) ([]string, []string, error) {
	if len(content) == 0 {
		return nil, nil, fmt.Errorf("no opening character")
	}

	if cl_char == "" {
		return nil, nil, fmt.Errorf("closing character cannot be empty")
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
