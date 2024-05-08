package Pair

import (
	"fmt"

	intf "github.com/PlayerR9/MyGoLib/Units/Common"
)

// Pair is a pair of values.
type Pair[A any, B any] struct {
	// The first value.
	First A

	// The second value.
	Second B
}

// String returns a string in the following format: (first, second). Here, both
// first and second are determined by the intf.StringOf() function.
//
// Returns:
//   - string: The string representation of the pair.
func (p *Pair[A, B]) String() string {
	return fmt.Sprintf("(%s, %s)", intf.StringOf(p.First), intf.StringOf(p.Second))
}

// Copy returns a shallow or deep copy of the pair according to the function
// intf.CopyOf().
//
// Returns:
//   - intf.Copier: A shallow or deep copy of the pair.
func (p *Pair[A, B]) Copy() intf.Copier {
	return &Pair[A, B]{
		First:  intf.CopyOf(p.First).(A),
		Second: intf.CopyOf(p.Second).(B),
	}
}

// Equals returns true if the pair is equal to the other pair according to the
// function intf.EqualOf().
//
// Parameters:
//   - other: The other pair to compare to.
//
// Returns:
//   - bool: True if the pair is equal to the other pair.
func (p *Pair[A, B]) Equals(other *Pair[A, B]) bool {
	return intf.EqualOf(p.First, other.First) && intf.EqualOf(p.Second, other.Second)
}

// Clean cleans the pair by first calling the intf.Clean() function on both the
// first and second values and then setting them to their zero values.
func (p *Pair[A, B]) Clean() {
	intf.Clean(p.First)
	intf.Clean(p.Second)

	p.First = *new(A)
	p.Second = *new(B)
}

// NewPair creates a new pair.
//
// Parameters:
//   - first: The first value.
//   - second: The second value.
//
// Returns:
//   - Pair[A, B]: The new pair.
func NewPair[A any, B any](first A, second B) *Pair[A, B] {
	return &Pair[A, B]{
		First:  first,
		Second: second,
	}
}
