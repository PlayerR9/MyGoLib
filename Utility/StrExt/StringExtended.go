package StrExt

import (
	"fmt"
	"slices"
)

func FindContentIndexes(openingToken, closingToken string, contentTokens []string) (int, int, error) {
	if openingToken == "" {
		panic("Opening token cannot be empty")
	}

	if closingToken == "" {
		panic("Closing token cannot be empty")
	}

	startIndex := slices.Index(contentTokens, openingToken)
	if startIndex == -1 {
		return 0, 0, fmt.Errorf("No opening token found; expected %s", openingToken)
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
			return startIndex, startIndex + i + 1, nil
		}
	}

	if closingToken != "\n" {
		return 0, 0, fmt.Errorf("No closing token found; expected %s", closingToken)
	}

	return startIndex, len(contentTokens) - 1, nil
}

func GetOrdinalSuffix(number int) string {
	if number < 0 {
		return fmt.Sprintf("%dth", number)
	}

	lastDigit := number % 10

	if (number%100 >= 11 && number%100 <= 13) || lastDigit > 3 || lastDigit == 0 {
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
