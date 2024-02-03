package StrExt

import "fmt"

// ErrWordTooLong is an error that occurs when a word is too long to fit
// within a certain width.
type ErrWordTooLong struct {
	word string
}

// Error method for ErrWordTooLong. It returns a formatted string indicating
// the word that was too long.
func (e *ErrWordTooLong) Error() string {
	return fmt.Sprintf("word '%s' is too long", e.word)
}

// ErrWidthTooSmall is an error that occurs when the width is too small to fit
// the text.
type ErrWidthTooSmall struct{}

// Error method for ErrWidthTooSmall. It returns a string indicating that the
// width was too small to fit the text.
func (e *ErrWidthTooSmall) Error() string {
	return "width is too small to fit the text"
}

// ErrSuffixTooLong is a struct that represents an error when a suffix
// is too long.
// It does not have any fields as the error condition is solely based
// on the length of the suffix.
type ErrSuffixTooLong struct{}

// Error is a method of the ErrSuffixTooLong type that implements the
// error interface. It returns a string representation of the error,
// that is, the string "suffix is too long".
func (e *ErrSuffixTooLong) Error() string {
	return "suffix is too long"
}

// ErrOpeningTokenEmpty is a struct that represents an error when an
// opening token is empty.
// It does not have any fields as the error condition is solely based
// on the absence of an opening token.
type ErrOpeningTokenEmpty struct{}

// Error is a method of the ErrOpeningTokenEmpty type that implements
// the error interface. It returns a string representation of the error,
// that is, the string "opening token is empty".
func (e *ErrOpeningTokenEmpty) Error() string {
	return "opening token is empty"
}

// ErrClosingTokenEmpty is a struct that represents an error when a
// closing token is empty.
// It does not have any fields as the error condition is solely based
// on the absence of a closing token.
type ErrClosingTokenEmpty struct{}

// Error is a method of the ErrClosingTokenEmpty type that implements
// the error interface. It returns a string representation of the error,
// that is, the string "closing token is empty".
func (e *ErrClosingTokenEmpty) Error() string {
	return "closing token is empty"
}

// ErrOpeningTokenNotFound is a struct that represents an error when an
// opening token is not found in the content.
// It does not have any fields as the error condition is solely based
// on the absence of an opening token in the content.
type ErrOpeningTokenNotFound struct {
	token string
}

// Error is a method of the ErrOpeningTokenNotFound type that implements
// the error interface. It returns a string representation of the error,
// that is, the string "opening token not found in content".
func (e *ErrOpeningTokenNotFound) Error() string {
	return fmt.Sprintf("opening token (%s) not found in content", e.token)
}

// ErrClosingTokenNotFound is a struct that represents an error when a
// closing token is not found in the content.
// It does not have any fields as the error condition is solely based
// on the absence of a closing token in the content.
type ErrClosingTokenNotFound struct {
	token string
}

// Error is a method of the ErrClosingTokenNotFound type that implements
// the error interface. It returns a string representation of the error,
// that is, the string "closing token not found in content".
func (e *ErrClosingTokenNotFound) Error() string {
	return fmt.Sprintf("closing token (%s) not found in content", e.token)
}

type ErrNeverOpened struct {
	openingToken, closingToken string
}

func (e *ErrNeverOpened) Error() string {
	return fmt.Sprintf("closing token (%s) found without a corresponding opening token (%s)", e.closingToken, e.openingToken)
}

// ErrEmptyText represents an error when a text input is empty.
type ErrEmptyText struct{}

// Error generates the error message for the ErrEmptyText error.
// The message indicates that the text input cannot be empty.
func (e *ErrEmptyText) Error() string {
	return "text cannot be empty"
}

// ErrHeightTooSmall represents an error when a height value is less than 1.
type ErrHeightTooSmall struct{}

// Error generates the error message for the ErrHeightTooSmall error.
// The message indicates that the height value must be at least 1.
func (e *ErrHeightTooSmall) Error() string {
	return "height must be at least 1"
}
