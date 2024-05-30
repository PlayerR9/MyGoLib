package Set

import (
	"strings"

	uc "github.com/PlayerR9/MyGoLib/Units/Common"
	ui "github.com/PlayerR9/MyGoLib/Units/Iterator"
	us "github.com/PlayerR9/MyGoLib/Units/Slices"
)

// EqualSet is a set that uses the Equals method to compare elements.
type EqualSet[T uc.Objecter] struct {
	// elems is the slice of elements in the set.
	elems []T
}

// IsEmpty checks if the set is empty.
//
// Returns:
//   - bool: True if the set is empty, false otherwise.
func (s *EqualSet[T]) IsEmpty() bool {
	return len(s.elems) == 0
}

// Size returns the number of elements in the set.
//
// Returns:
//   - int: The number of elements in the set.
func (s *EqualSet[T]) Size() int {
	return len(s.elems)
}

// HasElem checks if the set has the element.
//
// Parameters:
//   - elem: The element to check.
//
// Returns:
//   - bool: True if the set has the element, false otherwise.
func (s *EqualSet[T]) HasElem(elem T) bool {
	if len(s.elems) == 0 {
		return false
	}

	for _, e := range s.elems {
		if e.Equals(elem) {
			return true
		}
	}

	return false
}

// Add adds an element to the set.
//
// Parameters:
//   - elem: The element to add.
//
// Behaviors:
//   - If the element is already in the set, the function does nothing.
func (s *EqualSet[T]) Add(elem T) {
	for _, e := range s.elems {
		if e.Equals(elem) {
			return
		}
	}

	s.elems = append(s.elems, elem)
}

// Remove removes an element from the set.
//
// Parameters:
//   - elem: The element to remove.
//
// Behaviors:
//   - If the element is not in the set, the function does nothing.
func (s *EqualSet[T]) Remove(elem T) {
	if len(s.elems) == 0 {
		return
	}

	indexOf := us.FindEquals(s.elems, elem)
	if indexOf == -1 {
		return
	}

	s.elems = append(s.elems[:indexOf], s.elems[indexOf+1:]...)
}

// Union returns the union of the set with another set.
//
// Parameters:
//   - other: The other set.
//
// Returns:
//   - *EqualSet[T]: The union of the set with the other set.
func (s *EqualSet[T]) Union(other *EqualSet[T]) *EqualSet[T] {
	if other == nil {
		return s
	}

	elems := make([]T, len(s.elems))
	copy(elems, s.elems)
	limit := len(elems)

	for _, e := range other.elems {
		found := false

		for i := 0; i < limit; i++ {
			if elems[i].Equals(e) {
				found = true
				break
			}
		}

		if !found {
			elems = append(elems, e)
		}
	}

	return &EqualSet[T]{
		elems: elems,
	}
}

// Intersection returns the intersection of the set with another set.
//
// Parameters:
//   - other: The other set.
//
// Returns:
//   - *EqualSet[T]: The intersection of the set with the other set.
func (s *EqualSet[T]) Intersection(other *EqualSet[T]) *EqualSet[T] {
	if other == nil {
		return &EqualSet[T]{
			elems: make([]T, 0),
		}
	}

	newElems := make([]T, 0)

	for _, e := range s.elems {
		if other.HasElem(e) {
			newElems = append(newElems, e)
		}
	}

	return &EqualSet[T]{
		elems: newElems,
	}
}

// Difference returns the difference of the set with another set.
//
// Parameters:
//   - other: The other set.
//
// Returns:
//   - *EqualSet[T]: The difference of the set with the other set.
func (s *EqualSet[T]) Difference(other *EqualSet[T]) *EqualSet[T] {
	if other == nil {
		return s
	}

	newElems := make([]T, 0)

	for _, e := range s.elems {
		if !other.HasElem(e) {
			newElems = append(newElems, e)
		}
	}

	return &EqualSet[T]{
		elems: newElems,
	}
}

// SymmetricDifference returns the symmetric difference of the set with another set.
//
// Parameters:
//   - other: The other set.
//
// Returns:
//   - *EqualSet[T]: The symmetric difference of the set with the other set.
func (s *EqualSet[T]) SymmetricDifference(other *EqualSet[T]) *EqualSet[T] {
	if other == nil {
		return s
	}

	diff1 := make([]T, 0)

	for _, e := range s.elems {
		if !other.HasElem(e) {
			diff1 = append(diff1, e)
		}
	}

	diff2 := make([]T, 0)

	for _, e := range other.elems {
		if !s.HasElem(e) {
			diff2 = append(diff2, e)
		}
	}

	return &EqualSet[T]{
		elems: us.MergeUniqueEquals(diff1, diff2),
	}
}

// IsSubset checks if the set is a subset of another set.
//
// Parameters:
//   - other: The other set to check.
//
// Returns:
//   - bool: True if the set is a subset of the other set, false otherwise.
func (s *EqualSet[T]) IsSubset(other *EqualSet[T]) bool {
	if other == nil || len(s.elems) > len(other.elems) {
		return false
	}

	for _, k := range s.elems {
		if !other.HasElem(k) {
			return false
		}
	}

	return true
}

// Clear removes all elements from the set.
func (s *EqualSet[T]) Clear() {
	for i := 0; i < len(s.elems); i++ {
		s.elems[i] = *new(T)
	}

	s.elems = s.elems[:0]
}

// String returns a string representation of the set.
//
// Returns:
//   - string: The string representation of the set.
func (s *EqualSet[T]) String() string {
	if len(s.elems) == 0 {
		return "{}"
	}

	if len(s.elems) == 1 {
		return "{" + uc.StringOf(s.elems[0]) + "}"
	}

	var builder strings.Builder

	builder.WriteRune('{')
	builder.WriteString(uc.StringOf(s.elems[0]))
	for _, k := range s.elems[1:] {
		builder.WriteString(", ")
		builder.WriteString(uc.StringOf(k))
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
func (s *EqualSet[T]) Equals(other uc.Equaler) bool {
	if other == nil {
		return false
	}

	otherEs, ok := other.(*EqualSet[T])
	if !ok {
		return false
	}

	if len(s.elems) != len(otherEs.elems) {
		return false
	}

	for _, k := range s.elems {
		if !otherEs.HasElem(k) {
			return false
		}
	}

	return true
}

// Copy returns a copy of the set.
//
// Returns:
//   - *EqualSet[T]: A copy of the set.
func (s *EqualSet[T]) Copy() uc.Copier {
	newElems := make([]T, len(s.elems))
	copy(newElems, s.elems)

	return &EqualSet[T]{
		elems: newElems,
	}
}

// Slice returns a slice of the elements in the set.
//
// Returns:
//   - []T: A slice of the elements in the set.
func (s *EqualSet[T]) Slice() []T {
	return s.elems
}

// Iterator returns an iterator for the set.
//
// Returns:
//   - ui.Iterater[T]: An iterator for the set.
func (s *EqualSet[T]) Iterator() ui.Iterater[T] {
	return ui.NewSimpleIterator(s.elems)
}

// NewEqualSet creates a new EqualSet.
//
// Returns:
//   - *EqualSet: A new EqualSet.
func NewEqualSet[T uc.Objecter](elems []T) *EqualSet[T] {
	return &EqualSet[T]{
		elems: us.UniquefyEquals(elems, true),
	}
}
