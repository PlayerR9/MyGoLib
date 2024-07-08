package common

import (
	"fmt"
)

// Enumer is an interface representing an enumeration.
type Enumer interface {
	~int

	fmt.Stringer
}

// Copier is an interface that provides a method to create a copy of an element.
type Copier interface {
	// Copy creates a copy of the element.
	//
	// Returns:
	//   - Copier: The copy of the element.
	Copy() Copier
}

// SliceCopy creates a copy of a slice of elements by calling the Copy method
// of each element if the element implements the Copier interface.
//
// Whether or not the copy is a shallow or deep copy depends on the implementation of
// the Copy method of the element.
//
// Parameters:
//   - s: The slice of elements to copy.
//
// Returns:
//   - []T: The copy of the slice of elements.
func SliceCopy[T Copier](s []T) []T {
	sCopy := make([]T, 0, len(s))

	for _, elem := range s {
		elemCopy := elem.Copy().(T)

		sCopy = append(sCopy, elemCopy)
	}

	return sCopy
}
