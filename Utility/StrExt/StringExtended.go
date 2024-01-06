package StrExt

import (
	"fmt"
	"slices"
)

// ReplaceSuffix replaces the end of the given string with the provided suffix.
// The function checks the lengths of the string and the suffix to determine the appropriate action:
//   - If the length of the suffix is greater than the length of the string, the function returns an empty string and an ErrSuffixTooLong error.
//   - If the length of the suffix is equal to the length of the string, the function returns the suffix.
//   - If the length of the suffix is zero, the function returns the original string.
//   - Otherwise, the function replaces the end of the string with the suffix and returns the result.
//
// Parameters:
//
//   - str: The original string.
//   - suffix: The suffix to replace the end of the string.
//
// Returns:
//
//   - The modified string, or an error if the suffix is too long.
func ReplaceSuffix(str, suffix string) (string, error) {
	if len(str) < len(suffix) {
		return "", &ErrSuffixTooLong{}
	} else if len(str) == len(suffix) {
		return suffix, nil
	} else if len(suffix) == 0 {
		return str, nil
	}

	return str[:len(str)-len(suffix)] + suffix, nil
}

// FindContentIndexes searches for the positions of opening and closing tokens
// in a slice of strings.
// It returns the start and end indexes of the content between the tokens, and
// an error if any.
//
// Parameters:
//
//   - openingToken: The string that marks the beginning of the content.
//   - closingToken: The string that marks the end of the content.
//   - contentTokens: The slice of strings in which to search for the tokens.
//
// Returns:
//   - The start index of the content (inclusive).
//   - The end index of the content (exclusive).
//   - An error if the opening or closing tokens are empty, or if they are not
//     found in the correct order in the contentTokens.
//
// If the openingToken is found but the closingToken is not, the function will
// return an error.
// If the closingToken is found before the openingToken, the function will
// return an error.
// If the closingToken is a newline ("\n") and it is not found, the function will
// return the length of the contentTokens as the end index.
//
// Errors returned can be of type ErrOpeningTokenEmpty, ErrClosingTokenEmpty,
// ErrOpeningTokenNotFound, or a generic error with a message about the closing token.
func FindContentIndexes(openingToken, closingToken string, contentTokens []string) (int, int, error) {
	if openingToken == "" {
		return 0, 0, &ErrOpeningTokenEmpty{}
	}
	if closingToken == "" {
		return 0, 0, &ErrClosingTokenEmpty{}
	}

	openingTokenIndex := slices.Index(contentTokens, openingToken)
	if openingTokenIndex < 0 {
		return 0, 0, &ErrOpeningTokenNotFound{}
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
	} else if tokenBalance < 0 {
		return 0, 0, fmt.Errorf("closing token '%s' not opened", closingToken)
	} else if tokenBalance == 1 && closingToken == "\n" {
		return tokenStartIndex, len(contentTokens), nil
	} else {
		return 0, 0, fmt.Errorf("closing token '%s' not found", closingToken)
	}
}

// GetOrdinalSuffix returns the ordinal suffix for a given integer.
//
// Parameters:
//   - number: The integer for which to get the ordinal suffix.
//
// The function returns a string that represents the number with its ordinal suffix.
//
// For example, for the number 1, the function returns "1st"; for the number 2, it returns "2nd"; and so on.
// For numbers ending in 11, 12, or 13, the function returns the number with the suffix "th" (e.g., "11th", "12th", "13th").
// For negative numbers, the function also returns the number with the suffix "th".
//
// If the last digit of the number is 0 or greater than 3 (and the last two digits are not 11, 12, or 13), the function returns the number with the suffix "th".
// If the last digit of the number is 1, 2, or 3 (and the last two digits are not 11, 12, or 13), the function returns the number with the corresponding ordinal suffix ("st", "nd", or "rd").
func GetOrdinalSuffix(number int) string {
	if number < 0 {
		return fmt.Sprintf("%dth", number)
	}

	lastTwoDigits := number % 100
	lastDigit := lastTwoDigits % 10

	if lastTwoDigits >= 11 && lastTwoDigits <= 13 {
		return fmt.Sprintf("%dth", number)
	}

	if lastDigit == 0 || lastDigit > 3 {
		return fmt.Sprintf("%dth", number)
	}

	switch lastDigit {
	case 1:
		return fmt.Sprintf("%dst", number)
	case 2:
		return fmt.Sprintf("%dnd", number)
	case 3:
		return fmt.Sprintf("%drd", number)
	}

	return ""
}
