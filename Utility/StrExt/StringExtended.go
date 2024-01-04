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

func FindContentIndexes(openingToken, closingToken string, contentTokens []string) (int, int, bool) {
	if openingToken == "" {
		panic("Invalid input: Opening token cannot be empty")
	}

	if closingToken == "" {
		panic("Invalid input: Closing token cannot be empty")
	}

	startIndex := slices.Index(contentTokens, openingToken)
	if startIndex == -1 {
		panic("Invalid input: Opening token not found in the content")
	}

	tokenCounter := 1
	for i, token := range contentTokens[startIndex+1:] {
		switch token {
		case closingToken:
			tokenCounter--
		case openingToken:
			tokenCounter++
		}

		if tokenCounter == 0 {
			return startIndex, startIndex + i + 1, true
		}
	}

	return startIndex, len(contentTokens) - 1, closingToken == "\n"
}

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
