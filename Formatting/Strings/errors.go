package Strings

import "fmt"

// ErrLinesGreaterThanWords is an error type that is returned when the
// number of lines in a text is greater than the number of words.
type ErrLinesGreaterThanWords struct {
	// NumberOfLines is the number of lines in the text.
	NumberOfLines int

	// NumberOfWords is the number of words in the text.
	NumberOfWords int
}

// Error is a method of the error interface that returns the error message.
//
// Returns:
//  	- string: The error message.
func (e *ErrLinesGreaterThanWords) Error() string {
	return fmt.Sprintf(
		"number of lines (%d) is greater than the number of words (%d)",
		e.NumberOfLines, e.NumberOfWords,
	)
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
	return &ErrLinesGreaterThanWords{NumberOfLines: numberOfLines, NumberOfWords: numberOfWords}
}

// ErrNoCandidateFound is an error type that is returned when no candidate is found.
type ErrNoCandidateFound struct{}

// Error is a method of the error interface that returns the error message.
//
// Returns:
//   - string: The error message.
func (e *ErrNoCandidateFound) Error() string {
	return "no candidate found"
}

// NewErrNoCandidateFound is a constructor of ErrNoCandidateFound.
//
// Returns:
//   - *ErrNoCandidateFound: A pointer to the newly created error.
func NewErrNoCandidateFound() *ErrNoCandidateFound {
	return &ErrNoCandidateFound{}
}
