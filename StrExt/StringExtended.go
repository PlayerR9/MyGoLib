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
