package Sets

import (
	"strings"

	uc "github.com/PlayerR9/lib_units/common"
	lustr "github.com/PlayerR9/lib_units/strings"
)

// ComparableSet is a set that uses the == operator to compare elements.
type ComparableSet[T comparable] struct {
	// elems is the slice of elements in the set.
	elems map[T]bool
}

// IsEmpty checks if the set is empty.
//
// Returns:
//   - bool: True if the set is empty, false otherwise.
func (s *ComparableSet[T]) IsEmpty() bool {
	return len(s.elems) == 0
}

// Size returns the number of elements in the set.
//
// Returns:
//   - int: The number of elements in the set.
func (s *ComparableSet[T]) Size() int {
	return len(s.elems)
}

// HasElem checks if the set has the element.
//
// Parameters:
//   - elem: The element to check.
//
// Returns:
//   - bool: True if the set has the element, false otherwise.
func (s *ComparableSet[T]) HasElem(elem T) bool {
	if len(s.elems) == 0 {
		return false
	}

	_, ok := s.elems[elem]
	return ok
}

// Add adds an element to the set.
//
// Parameters:
//   - elem: The element to add.
//
// Behaviors:
//   - If the element is already in the set, the function does nothing.
func (s *ComparableSet[T]) Add(elem T) {
	s.elems[elem] = true
}

// Remove removes an element from the set.
//
// Parameters:
//   - elem: The element to remove.
//
// Behaviors:
//   - If the element is not in the set, the function does nothing.
func (s *ComparableSet[T]) Remove(elem T) {
	delete(s.elems, elem)
}

// Union returns the union of the set with another set.
//
// Parameters:
//   - other: The other set.
//
// Returns:
//   - *ComparableSet[T]: The union of the set with the other set.
func (s *ComparableSet[T]) Union(other *ComparableSet[T]) *ComparableSet[T] {
	if other == nil {
		return s
	}

	newElems := make(map[T]bool)

	for k := range s.elems {
		newElems[k] = true
	}

	for k := range other.elems {
		newElems[k] = true
	}

	return &ComparableSet[T]{
		elems: newElems,
	}
}

// Intersection returns the intersection of the set with another set.
//
// Parameters:
//   - other: The other set.
//
// Returns:
//   - *ComparableSet[T]: The intersection of the set with the other set.
func (s *ComparableSet[T]) Intersection(other *ComparableSet[T]) *ComparableSet[T] {
	if other == nil {
		return &ComparableSet[T]{
			elems: make(map[T]bool),
		}
	}

	newElems := make(map[T]bool)

	for k := range s.elems {
		_, ok := other.elems[k]
		if ok {
			newElems[k] = true
		}
	}

	return &ComparableSet[T]{
		elems: newElems,
	}
}

// Difference returns the difference of the set with another set.
//
// Parameters:
//   - other: The other set.
//
// Returns:
//   - *ComparableSet[T]: The difference of the set with the other set.
func (s *ComparableSet[T]) Difference(other *ComparableSet[T]) *ComparableSet[T] {
	if other == nil {
		return s
	}

	newElems := make(map[T]bool)

	for k := range s.elems {
		_, ok := other.elems[k]
		if !ok {
			newElems[k] = true
		}
	}

	return &ComparableSet[T]{
		elems: newElems,
	}
}

// SymmetricDifference returns the symmetric difference of the set with another set.
//
// Parameters:
//   - other: The other set.
//
// Returns:
//   - *ComparableSet[T]: The symmetric difference of the set with the other set.
func (s *ComparableSet[T]) SymmetricDifference(other *ComparableSet[T]) *ComparableSet[T] {
	if other == nil {
		return s
	}

	diff1 := make(map[T]bool)

	for k := range s.elems {
		_, ok := other.elems[k]
		if !ok {
			diff1[k] = true
		}
	}

	diff2 := make(map[T]bool)

	for k := range other.elems {
		_, ok := s.elems[k]
		if !ok {
			diff2[k] = true
		}
	}

	newElems := make(map[T]bool)

	for k := range diff1 {
		_, ok := diff2[k]
		if ok {
			newElems[k] = true
		}
	}

	return &ComparableSet[T]{
		elems: newElems,
	}
}

// IsSubset checks if the set is a subset of another set.
//
// Parameters:
//   - other: The other set to check.
//
// Returns:
//   - bool: True if the set is a subset of the other set, false otherwise.
func (s *ComparableSet[T]) IsSubset(other *ComparableSet[T]) bool {
	if other == nil || len(s.elems) > len(other.elems) {
		return false
	}

	for k := range s.elems {
		_, ok := other.elems[k]
		if !ok {
			return false
		}
	}

	return true
}

// Clear removes all elements from the set.
func (s *ComparableSet[T]) Clear() {
	s.elems = make(map[T]bool)
}

// String returns a string representation of the set.
//
// Returns:
//   - string: The string representation of the set.
func (s *ComparableSet[T]) String() string {
	if len(s.elems) == 0 {
		return "{}"
	}

	allKeys := make([]T, 0, len(s.elems))
	for k := range s.elems {
		allKeys = append(allKeys, k)
	}

	if len(allKeys) == 1 {
		return "{" + lustr.GoStringOf(allKeys[0]) + "}"
	}

	var builder strings.Builder

	builder.WriteRune('{')
	builder.WriteString(lustr.GoStringOf(allKeys[0]))
	for _, k := range allKeys[1:] {
		builder.WriteString(", ")
		builder.WriteString(lustr.GoStringOf(k))
	}
	builder.WriteRune('}')

	return builder.String()
}

// Equals checks if the set is equal to another set.
//
// Parameters:
//   - other: The other set to compare.
//
// Returns:
//   - bool: True if the sets are equal, false otherwise.
func (s *ComparableSet[T]) Equals(other *ComparableSet[T]) bool {
	if other == nil {
		return false
	}

	if len(s.elems) != len(other.elems) {
		return false
	}

	for k := range s.elems {
		_, ok := other.elems[k]
		if !ok {
			return false
		}
	}

	return true
}

// Copy returns a copy of the set.
//
// Returns:
//   - *ComparableSet[T]: A copy of the set.
func (s *ComparableSet[T]) Copy() *ComparableSet[T] {
	newElems := make(map[T]bool)

	for k := range s.elems {
		newElems[k] = true
	}

	return &ComparableSet[T]{
		elems: newElems,
	}
}

// Slice returns a slice of the elements in the set.
//
// Returns:
//   - []T: A slice of the elements in the set.
func (s *ComparableSet[T]) Slice() []T {
	if len(s.elems) == 0 {
		return nil
	}

	slice := make([]T, 0, len(s.elems))

	for k := range s.elems {
		slice = append(slice, k)
	}

	return slice
}

// Iterator returns an iterator for the set.
//
// Returns:
//   - uc.Iterater[T]: An iterator for the set.
func (s *ComparableSet[T]) Iterator() uc.Iterater[T] {
	var builder uc.Builder[T]

	for k, ok := range s.elems {
		if ok {
			builder.Add(k)
		}
	}

	return builder.Build()
}

// NewComparableSet creates a new ComparableSet.
//
// Parameters:
//   - elems: The elements to add to the set.
//
// Returns:
//   - *ComparableSet[T]: A new ComparableSet.
//
// Behaviors:
//   - It ignores duplicate elements.
func NewComparableSet[T comparable](elems []T) *ComparableSet[T] {
	set := &ComparableSet[T]{
		elems: make(map[T]bool),
	}

	if len(elems) == 0 {
		return set
	}

	seen := make(map[T]bool)

	for _, elem := range elems {
		seen[elem] = true
	}

	return &ComparableSet[T]{
		elems: seen,
	}
}
