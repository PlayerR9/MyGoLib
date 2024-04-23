package StringExt

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
	return fmt.Sprintf("opening token (%q) not found in content", e.token)
}

// NewErrOpeningTokenNotFound is a constructor of ErrOpeningTokenNotFound.
//
// Parameters:
//
//   - token: The opening token that was not found in the content.
//
// Returns:
//
//   - *ErrOpeningTokenNotFound: A pointer to the newly created error.
func NewErrOpeningTokenNotFound(token string) *ErrOpeningTokenNotFound {
	return &ErrOpeningTokenNotFound{token: token}
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
	return fmt.Sprintf("closing token (%q) not found in content", e.token)
}

// NewErrClosingTokenNotFound is a constructor of ErrClosingTokenNotFound.
//
// Parameters:
//
//   - token: The closing token that was not found in the content.
//
// Returns:
//
//   - *ErrClosingTokenNotFound: A pointer to the newly created error.
func NewErrClosingTokenNotFound(token string) *ErrClosingTokenNotFound {
	return &ErrClosingTokenNotFound{token: token}
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
	return fmt.Sprintf("closing token (%q) found without a corresponding opening token (%q)",
		e.closingToken, e.openingToken)
}

// NewErrNeverOpened is a constructor of ErrNeverOpened.
//
// Parameters:
//
//   - openingToken: The opening token that was never closed.
//   - closingToken: The closing token that was found without a corresponding opening token.
//
// Returns:
//
//   - *ErrNeverOpened: A pointer to the newly created error.
func NewErrNeverOpened(openingToken, closingToken string) *ErrNeverOpened {
	return &ErrNeverOpened{openingToken: openingToken, closingToken: closingToken}
}
