package StrExt

import (
	"fmt"
	"slices"
)

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
