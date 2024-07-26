package StringExt

import (
	"strconv"
	"strings"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// ErrInvalidUTF8Encoding is an error type for invalid UTF-8 encoding.
type ErrInvalidUTF8Encoding struct{}

// Error implements the error interface.
//
// Message: "invalid UTF-8 encoding"
func (e *ErrInvalidUTF8Encoding) Error() string {
	return "invalid UTF-8 encoding"
}

// NewErrInvalidUTF8Encoding creates a new ErrInvalidUTF8Encoding error.
//
// Returns:
//   - *ErrInvalidUTF8Encoding: A pointer to the newly created error.
func NewErrInvalidUTF8Encoding() *ErrInvalidUTF8Encoding {
	e := &ErrInvalidUTF8Encoding{}
	return e
}

// ErrLongerSuffix is a struct that represents an error when the suffix is
// longer than the string.
type ErrLongerSuffix struct {
	// Str is the string that is shorter than the suffix.
	Str string

	// Suffix is the Suffix that is longer than the string.
	Suffix string
}

// Error implements the error interface.
//
// Message: "suffix {Suffix} is longer than the string {Str}"
func (e *ErrLongerSuffix) Error() string {
	values := []string{
		"suffix",
		strconv.Quote(e.Suffix),
		"is longer than the string",
		strconv.Quote(e.Str),
	}

	msg := strings.Join(values, " ")

	return msg
}

// NewErrLongerSuffix is a constructor of ErrLongerSuffix.
//
// Parameters:
//   - str: The string that is shorter than the suffix.
//   - suffix: The suffix that is longer than the string.
//
// Returns:
//   - *ErrLongerSuffix: A pointer to the newly created error.
func NewErrLongerSuffix(str, suffix string) *ErrLongerSuffix {
	e := &ErrLongerSuffix{
		Str:    str,
		Suffix: suffix,
	}
	return e
}

// ErrLinesGreaterThanWords is an error type that is returned when the
// number of lines in a text is greater than the number of words.
type ErrLinesGreaterThanWords struct {
	// NumberOfLines is the number of lines in the text.
	NumberOfLines int

	// NumberOfWords is the number of words in the text.
	NumberOfWords int
}

// Error implements the error interface.
//
// Message: "number of lines ({NumberOfLines}) is greater than the number of words ({NumberOfWords})"
func (e *ErrLinesGreaterThanWords) Error() string {
	values := []string{
		"number of lines",
		uc.QuoteInt(e.NumberOfLines),
		"is greater than the number of words",
		uc.QuoteInt(e.NumberOfWords),
	}

	msg := strings.Join(values, " ")

	return msg
}

// NewErrLinesGreaterThanWords is a constructor of ErrLinesGreaterThanWords.
//
// Parameters:
//   - numberOfLines: The number of lines in the text.
//   - numberOfWords: The number of words in the text.
//
// Returns:
//   - *ErrLinesGreaterThanWords: A pointer to the newly created error.
func NewErrLinesGreaterThanWords(numberOfLines, numberOfWords int) *ErrLinesGreaterThanWords {
	e := &ErrLinesGreaterThanWords{
		NumberOfLines: numberOfLines,
		NumberOfWords: numberOfWords,
	}
	return e
}

// ErrNoCandidateFound is an error type that is returned when no candidate is found.
type ErrNoCandidateFound struct{}

// Error implements the error interface.
//
// Message: "no candidate found"
func (e *ErrNoCandidateFound) Error() string {
	return "no candidate found"
}

// NewErrNoCandidateFound is a constructor of ErrNoCandidateFound.
//
// Returns:
//   - *ErrNoCandidateFound: A pointer to the newly created error.
func NewErrNoCandidateFound() *ErrNoCandidateFound {
	e := &ErrNoCandidateFound{}
	return e
}

// ErrNoClosestWordFound is an error when no closest word is found.
type ErrNoClosestWordFound struct{}

// Error implements the error interface.
//
// Message: "no closest word was found"
func (e *ErrNoClosestWordFound) Error() string {
	return "no closest word was found"
}

// NewErrNoClosestWordFound creates a new ErrNoClosestWordFound.
//
// Returns:
//   - *ErrNoClosestWordFound: The new ErrNoClosestWordFound.
func NewErrNoClosestWordFound() *ErrNoClosestWordFound {
	e := &ErrNoClosestWordFound{}
	return e
}
