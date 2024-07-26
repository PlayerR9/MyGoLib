package slice

import (
	"strings"
)

// ErrTokenNotFound is a struct that represents an error when a token is not
// found in the content.
type ErrTokenNotFound struct {
	// IsOpening is the type of the token (opening or closing).
	IsOpening bool
}

// Error implements the error interface.
//
// Message: "{Type} token is not in the content"
func (e *ErrTokenNotFound) Error() string {
	var builder strings.Builder

	if e.IsOpening {
		builder.WriteString("opening")
	} else {
		builder.WriteString("closing")
	}

	builder.WriteString(" token is not in the content")

	return builder.String()
}

// NewErrTokenNotFound is a constructor of ErrTokenNotFound.
//
// Parameters:
//   - is_opening: The type of the token (opening or closing).
//
// Returns:
//   - *ErrTokenNotFound: A pointer to the newly created error.
func NewErrTokenNotFound(is_opening bool) *ErrTokenNotFound {
	return &ErrTokenNotFound{
		IsOpening: is_opening,
	}
}

// ErrNeverOpened is a struct that represents an error when a closing
// token is found without a corresponding opening token.
type ErrNeverOpened struct{}

// Error implements the error interface.
//
// Message:
//   - "closing token found without a corresponding opening token".
func (e *ErrNeverOpened) Error() string {
	return "closing token is found without a corresponding opening token"
}

// NewErrNeverOpened is a constructor of ErrNeverOpened.
//
// Returns:
//   - *ErrNeverOpened: A pointer to the newly created error.
func NewErrNeverOpened() *ErrNeverOpened {
	return &ErrNeverOpened{}
}
