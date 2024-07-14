package go_generator

import (
	"strconv"
	"strings"
)

// ErrInvalidID represents an error when an identifier is invalid.
type ErrInvalidID struct {
	// ID is the invalid identifier.
	ID string

	// Reason is the reason why the identifier is invalid.
	Reason error
}

// Error implements the error interface.
//
// Message: "identifier <id> is invalid: <reason>"
func (e *ErrInvalidID) Error() string {
	q_id := strconv.Quote(e.ID)

	var reason string
	var builder strings.Builder

	if e.Reason != nil {
		re := e.Reason.Error()

		builder.WriteString(": ")
		builder.WriteString(re)

		reason = builder.String()
		builder.Reset()
	}

	builder.WriteString("identifier ")
	builder.WriteString(q_id)
	builder.WriteString(" is invalid")
	builder.WriteString(reason)

	str := builder.String()
	return str
}

// NewErrInvalidID creates a new ErrInvalidID error.
//
// Parameters:
//   - id: The invalid identifier.
//   - reason: The reason for the error.
//
// Returns:
//   - *ErrInvalidID: The new error.
func NewErrInvalidID(id string, reason error) *ErrInvalidID {
	e := &ErrInvalidID{
		ID:     id,
		Reason: reason,
	}

	return e
}
