package common

import (
	"strings"
)

// Pair is a pair of values.
type Pair[A, B any] struct {
	// The first value.
	First A

	// The second value.
	Second B
}

// String implements the fmt.Stringer interface.
func (p Pair[A, B]) String() string {
	var builder strings.Builder

	builder.WriteRune('(')
	builder.WriteString(StringOf(p.First))
	builder.WriteString(", ")
	builder.WriteString(StringOf(p.Second))
	builder.WriteRune(')')

	return builder.String()
}

// NewPair creates a new pair.
//
// Parameters:
//   - first: The first value.
//   - second: The second value.
//
// Returns:
//   - Pair[A, B]: The new pair.
func NewPair[A any, B any](first A, second B) Pair[A, B] {
	p := Pair[A, B]{
		First:  first,
		Second: second,
	}

	return p
}
