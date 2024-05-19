package StringExt

import "fmt"

// ErrInvalidUTF8Encoding is an error type for invalid UTF-8 encoding.
type ErrInvalidUTF8Encoding struct{}

// Error is a method of the error interface that returns the error message.
//
// Returns:
//   - string: The error message.
func (e *ErrInvalidUTF8Encoding) Error() string {
	return "invalid UTF-8 encoding"
}

// NewErrInvalidUTF8Encoding creates a new ErrInvalidUTF8Encoding error.
//
// Returns:
//   - *ErrInvalidUTF8Encoding: A pointer to the newly created error.
func NewErrInvalidUTF8Encoding() *ErrInvalidUTF8Encoding {
	return &ErrInvalidUTF8Encoding{}
}

// ErrLongerSuffix is a struct that represents an error when the suffix is
// longer than the string.
type ErrLongerSuffix struct {
	// Str is the string that is shorter than the suffix.
	Str string

	// Suffix is the Suffix that is longer than the string.
	Suffix string
}

// Error is a method of error interface that returns the error message.
//
// Returns:
//  	- string: The error message.
func (e *ErrLongerSuffix) Error() string {
	return fmt.Sprintf("suffix (%s) is longer than the string (%s)", e.Suffix, e.Str)
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
	return &ErrLongerSuffix{Str: str, Suffix: suffix}
}

// ErrTokenNotFound is a struct that represents an error when a token is not
// found in the content.
type ErrTokenNotFound struct {
	// Token is the token that was not found in the content.
	Token string

	// Type is the type of the token (opening or closing).
	Type TokenType
}

// Error is a method of the error interface that returns the error message.
//
// Returns:
//   - string: The error message.
func (e *ErrTokenNotFound) Error() string {
	return fmt.Sprintf("%s token (%q) is not in the content", e.Type.String(), e.Token)
}

// NewErrTokenNotFound is a constructor of ErrTokenNotFound.
//
// Parameters:
//   - token: The token that was not found in the content.
//   - tokenType: The type of the token (opening or closing).
//
// Returns:
//   - *ErrTokenNotFound: A pointer to the newly created error.
func NewErrTokenNotFound(token string, tokenType TokenType) *ErrTokenNotFound {
	return &ErrTokenNotFound{Token: token, Type: tokenType}
}

// ErrNeverOpened is a struct that represents an error when a closing
// token is found without a corresponding opening token.
type ErrNeverOpened struct {
	// OpeningToken and ClosingToken are the opening and closing tokens,
	// respectively.
	OpeningToken, ClosingToken string
}

// Error is a method of the error interface that returns the error message.
//
// Returns:
//   - string: The error message.
func (e *ErrNeverOpened) Error() string {
	return fmt.Sprintf("closing token (%q) found without a corresponding opening token (%q)",
		e.ClosingToken, e.OpeningToken)
}

// NewErrNeverOpened is a constructor of ErrNeverOpened.
//
// Parameters:
//   - openingToken: The opening token that was never closed.
//   - closingToken: The closing token that was found without a corresponding opening token.
//
// Returns:
//   - *ErrNeverOpened: A pointer to the newly created error.
func NewErrNeverOpened(openingToken, closingToken string) *ErrNeverOpened {
	return &ErrNeverOpened{OpeningToken: openingToken, ClosingToken: closingToken}
}

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
