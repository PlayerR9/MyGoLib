package StrExt

import "fmt"

// ErrOpeningTokenNotFound is a struct that represents an error when an
// opening token is not found in the content.
type ErrOpeningTokenNotFound struct {
	token string
}

// Error is a method of the ErrOpeningTokenNotFound type that implements
// the error interface.
//
// Returns:
//
//   - string: The error message.
func (e *ErrOpeningTokenNotFound) Error() string {
	return fmt.Sprintf("opening token (%s) not found in content", e.token)
}

// ErrClosingTokenNotFound is a struct that represents an error when a
// closing token is not found in the content.
type ErrClosingTokenNotFound struct {
	// token is the closing token that was not found in the content.
	token string
}

// Error is a method of the ErrClosingTokenNotFound type that implements
// the error interface.
//
// Returns:
//
//   - string: The error message.
func (e *ErrClosingTokenNotFound) Error() string {
	return fmt.Sprintf("closing token (%s) not found in content", e.token)
}

// ErrNeverOpened is a struct that represents an error when a closing
// token is found without a corresponding opening token.
type ErrNeverOpened struct {
	// openingToken and closingToken are the opening and closing tokens,
	// respectively.
	openingToken, closingToken string
}

// Error is a method of the ErrNeverOpened type that implements the error
// interface.
//
// Returns:
//
//   - string: The error message.
func (e *ErrNeverOpened) Error() string {
	return fmt.Sprintf("closing token (%s) found without a corresponding opening token (%s)",
		e.closingToken, e.openingToken)
}
