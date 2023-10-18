package StrExt

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode"
	"unicode/utf8"
)

// GLOBAL VARIABLES
var (
	// DebugMode is a boolean that is used to enable or disable debug mode. When debug mode is enabled, the package will print debug messages.
	// **Note:** Debug mode is disabled by default.
	DebugMode bool = false

	debugger *log.Logger = log.New(os.Stdout, "[StrExt] ", log.LstdFlags) // Debugger
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
	characters := make([]rune, 0) // The characters in the string

	// Loop through the string
	char_index := 1 // The index of the character in the string

	for current := 0; current < len(str); char_index++ {
		character, size := utf8.DecodeRuneInString(str[current:]) // The character and its size in bytes

		if character == utf8.RuneError {
			// The character is not a valid UTF-8 character
			return nil, fmt.Errorf("%s character is not a valid UTF-8 character", PrintOrdinalNumber(char_index))
		}

		characters = append(characters, character) // Add the character to the slice of characters
		current += size                            // Move to the next character
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
	// Loop through the string
	for current := 0; current < len(str); {
		character, size := utf8.DecodeRuneInString(str)

		if character == utf8.RuneError || !unicode.IsLetter(character) {
			return false
		}

		current += size
	}

	return true
}

// FindContentIndexes finds the first occurrence of a block of tokens (defined as 'content') between the opening and closing token so that
// they match; which means they close each other. Best for extracting content between parentheses, brackets, etc. Panics if the opening or closing tokens are empty.
//
// Parameters:
//   - opening: The opening token.
//   - closing: The closing token.
//   - tokens: The tokens to extract from.
//
// Returns:
//   - int: The index of the first character of the content. (The opening token is included in the content)
//   - int: The index of the last character of the content. (The closing token is included in the content)
//   - error: An error if one occurred.
//
// Example:
//   - FindContentIndexes("(", ")", ["a", "(", "hello", ")"]) returns 1, 3, nil; because the content is "(hello)" and the last character of the content is at index 3.
func FindContentIndexes(opening, closing string, tokens []string) (int, int, error) {
	if opening == "" {
		// Opening character is empty, panic.
		if DebugMode {
			debugger.Panic("opening character cannot be empty")
		} else {
			panic("opening character cannot be empty")
		}
	}

	if closing == "" {
		// Closing character is empty, panic.
		if DebugMode {
			debugger.Panic("closing character cannot be empty")
		} else {
			panic("closing character cannot be empty")
		}
	}

	// Initialize the starting index
	starting_index := 0

	// Find the first opening character
	for starting_index < len(tokens) && tokens[starting_index] != opening {
		starting_index++
	}

	if starting_index == len(tokens) {
		// No opening character found, return an error
		return 0, 0, fmt.Errorf("no opening character found; expected %s", opening)
	}

	// Find the closing character

	counter := 1 // The number of opening characters minus the number of closing characters

	for i := starting_index + 1; i < len(tokens); i++ {
		if tokens[i] == closing {
			counter--
		} else if tokens[i] == opening {
			counter++
		}

		if counter == 0 {
			// Found the closing character, return the indexes
			return starting_index, i, nil
		}
	}

	if closing != "\n" {
		// No closing character found, return an error
		return 0, 0, fmt.Errorf("no closing character found; expected %s", closing)
	}

	return starting_index, len(tokens) - 1, nil
}

// Similar to FindContentIndexes, but does not require an opening token. It returns the index of the first occurrence of the closing token.
// Best for extracting content up to a closing token, like a newline character. Panics if the closing token is empty.
//
// Parameters:
//   - closing: The closing token.
//   - tokens: The tokens to extract from.
//
// Returns:
//   - int: The index of the first occurrence of the closing token.
//   - error: An error if one occurred.
func FindIndexesUpTo(closing string, tokens []string) (int, error) {
	if closing == "" {
		// Closing character is empty, panic.
		if DebugMode {
			debugger.Panic("closing character cannot be empty")
		} else {
			panic("closing character cannot be empty")
		}
	}

	// Find the closing character
	for i := 0; i < len(tokens); i++ {
		if tokens[i] == closing {
			// Found the closing character, return the index
			return i, nil
		}
	}

	if closing != "\n" {
		// No closing character found, return an error
		return 0, fmt.Errorf("no closing character found; expected %s", closing)
	}

	return len(tokens) - 1, nil
}

// PrintOrdinalNumber prints the ordinal number of a position. For example, if the position is 1, it will print "1st". Panics if the position is negative.
//
// Parameters:
//   - position: The position to print.
//
// Returns:
//   - string: The ordinal number of the position.
func PrintOrdinalNumber(position int) string {
	if position < 0 {
		// Position is negative, panic.
		if DebugMode {
			debugger.Panic("position cannot be negative")
		} else {
			panic("position cannot be negative")
		}
	}

	if position >= 11 && position < 20 {
		// 11th, 12th, 13th, 14th, 15th, 16th, 17th, 18th, or 19th
		return fmt.Sprintf("%dth", position)
	}

	// Get the last digit of the position
	last_digit := position % 10

	switch last_digit {
	case 1:
		// 1st
		return fmt.Sprintf("%dst", position)
	case 2:
		// 2nd
		return fmt.Sprintf("%dnd", position)
	case 3:
		// 3rd
		return fmt.Sprintf("%drd", position)
	default:
		// 4th, 5th, 6th, 7th, 8th, 9th, 0th
		return fmt.Sprintf("%dth", position)
	}
}

// DeEscapeSequence converts an escape sequence to its character. If the character is not an escape sequence, it returns the character as a string.
// For example, if the character is '\n', it will return "\\n".
//
// Parameters:
//   - char: The character to convert.
//
// Returns:
//   - string: The character as a string.
//
// Escape sequences: [\a \b \f \n \r \t \v]
func DeEscapeSequence(char rune) string {
	switch char {
	case '\a':
		return "\\a"
	case '\b':
		return "\\b"
	case '\f':
		return "\\f"
	case '\n':
		return "\\n"
	case '\r':
		return "\\r"
	case '\t':
		return "\\t"
	case '\v':
		return "\\v"
	default:
		return string(char)
	}
}
