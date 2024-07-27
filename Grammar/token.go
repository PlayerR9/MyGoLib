package Grammar

import (
	"fmt"
	"strconv"
	"strings"
)

// TokenTyper is an interface that defines the behavior of a token type.
type TokenTyper interface {
	~int

	fmt.Stringer
}

// Token is a struct that represents a generic token of type T.
type Token[T TokenTyper] struct {
	// Type is the type of the token.
	Type T

	// Data is the data of the token.
	Data string

	// Lookahead is the lookahead token.
	Lookahead *Token[T]

	// At is the position of the token in the input. It is the byte position
	// of the token in the input.
	At int
}

// String implements the fmt.Stringer interface.
//
// Format:
//
//	"Token[T][{{ .Type }} ({{ .Data }})] : {{ .At }}]"
func (t *Token[T]) String() string {
	var builder strings.Builder

	builder.WriteString("Token[T][")
	builder.WriteString(t.Type.String())

	if t.Data != "" {
		builder.WriteString(" (")
		builder.WriteString(strconv.Quote(t.Data))
		builder.WriteRune(')')
	}

	builder.WriteString(" : ")
	builder.WriteString(strconv.Itoa(t.At))
	builder.WriteRune(']')

	return builder.String()
}

// NewToken creates a new token of type T.
//
// Parameters:
//   - t: The type of the token.
//   - d: The data of the token.
//   - at: The position of the token in the input.
//   - lookahead: The lookahead token.
//
// Returns:
//   - *Token[T]: A pointer to the newly created token. Never returns nil.
func NewToken[T TokenTyper](t T, d string, at int, lookahead *Token[T]) *Token[T] {
	return &Token[T]{
		Type:      t,
		Data:      d,
		Lookahead: lookahead,
		At:        at,
	}
}
