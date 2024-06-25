package Sets

import (
	"fmt"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// Seter is an interface for a set.
type Seter[T any] interface {
	// IsEmpty checks if the set is empty.
	//
	// Returns:
	//  - bool: True if the set is empty, false otherwise.
	IsEmpty() bool

	// Size returns the number of elements in the set.
	//
	// Returns:
	//  - int: The number of elements in the set.
	Size() int

	// HasElem checks if the set has the element.
	//
	// Parameters:
	//  - elem: The element to check.
	//
	// Returns:
	//  - bool: True if the set has the element, false otherwise.
	HasElem(elem T) bool

	// Add adds an element to the set.
	//
	// Parameters:
	//  - elem: The element to add.
	//
	// Behaviors:
	//  - If the element is already in the set, the function does nothing.
	Add(elem T)

	// Remove removes an element from the set.
	//
	// Parameters:
	//  - elem: The element to remove.
	//
	// Behaviors:
	//  - If the element is not in the set, the function does nothing.
	Remove(elem T)

	// Union returns the union of the set with another set.
	//
	// Parameters:
	//  - other: The other set.
	//
	// Returns:
	//  - Seter[T]: The union of the set with the other set.
	Union(other Seter[T]) Seter[T]

	// Intersection returns the intersection of the set with another set.
	//
	// Parameters:
	//  - other: The other set.
	//
	// Returns:
	//  - Seter[T]: The intersection of the set with the other set.
	Intersection(other Seter[T]) Seter[T]

	// Difference returns the difference of the set with another set.
	//
	// Parameters:
	//  - other: The other set.
	//
	// Returns:
	//  - Seter[T]: The difference of the set with the other set.
	Difference(other Seter[T]) Seter[T]

	// SymmetricDifference returns the symmetric difference of the set with another set.
	//
	// Parameters:
	//  - other: The other set.
	//
	// Returns:
	//  - Seter[T]: The symmetric difference of the set with the other set.
	SymmetricDifference(other Seter[T]) Seter[T]

	// IsSubset checks if the set is a subset of another set.
	//
	// Parameters:
	//  - other: The other set.
	//
	// Returns:
	//  - bool: True if the set is a subset of the other set, false otherwise.
	IsSubset(other Seter[T]) bool

	// Clear removes all elements from the set.
	Clear()

	fmt.Stringer
	uc.Objecter
	uc.Slicer[T]
}
