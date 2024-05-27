package Set

import (
	"strings"

	uc "github.com/PlayerR9/MyGoLib/Units/Common"
	ui "github.com/PlayerR9/MyGoLib/Units/Iterator"
	us "github.com/PlayerR9/MyGoLib/Units/Slices"
	"golang.org/x/exp/slices"
)

// LessSet is a set that uses the Equals method to compare elements.
type LessSet[T uc.Comparer] struct {
	// elems is the slice of elements in the set.
	elems []T
}

// IsEmpty checks if the set is empty.
//
// Returns:
//   - bool: True if the set is empty, false otherwise.
func (s *LessSet[T]) IsEmpty() bool {
	return len(s.elems) == 0
}

// Size returns the number of elements in the set.
//
// Returns:
//   - int: The number of elements in the set.
func (s *LessSet[T]) Size() int {
	return len(s.elems)
}

// HasElem checks if the set has the element.
//
// Parameters:
//   - elem: The element to check.
//
// Returns:
//   - bool: True if the set has the element, false otherwise.
func (s *LessSet[T]) HasElem(elem T) bool {
	if len(s.elems) == 0 {
		return false
	}

	_, ok := uc.Find(s.elems, elem)

	return ok
}

// Add adds an element to the set.
//
// Parameters:
//   - elem: The element to add.
//
// Behaviors:
//   - If the element is already in the set, the function does nothing.
func (s *LessSet[T]) Add(elem T) {
	if len(s.elems) == 0 {
		s.elems = append(s.elems, elem)
		return
	}

	pos, ok := uc.Find(s.elems, elem)
	if ok {
		return
	}

	s.elems = slices.Insert(s.elems, pos, elem)
}

// Remove removes an element from the set.
//
// Parameters:
//   - elem: The element to remove.
//
// Behaviors:
//   - If the element is not in the set, the function does nothing.
func (s *LessSet[T]) Remove(elem T) {
	if len(s.elems) == 0 {
		return
	}

	pos, ok := uc.Find(s.elems, elem)
	if !ok {
		return
	}

	s.elems = slices.Delete(s.elems, pos, pos+1)
}

// Union returns the union of the set with another set.
//
// Parameters:
//   - other: The other set.
//
// Returns:
//   - *LessSet[T]: The union of the set with the other set.
func (s *LessSet[T]) Union(other *LessSet[T]) *LessSet[T] {
	if other == nil {
		return s
	}

	elems := make([]T, len(s.elems))
	copy(elems, s.elems)

	for _, e := range other.elems {
		pos, ok := uc.Find(elems, e)
		if ok {
			continue
		}

		elems = slices.Insert(elems, pos, e)
	}

	return &LessSet[T]{
		elems: elems,
	}
}

// Intersection returns the intersection of the set with another set.
//
// Parameters:
//   - other: The other set.
//
// Returns:
//   - *LessSet[T]: The intersection of the set with the other set.
func (s *LessSet[T]) Intersection(other *LessSet[T]) *LessSet[T] {
	if other == nil {
		return &LessSet[T]{
			elems: make([]T, 0),
		}
	}

	newElems := make([]T, 0)

	for _, e := range s.elems {
		if other.HasElem(e) {
			newElems = append(newElems, e)
		}
	}

	return &LessSet[T]{
		elems: newElems,
	}
}

// Difference returns the difference of the set with another set.
//
// Parameters:
//   - other: The other set.
//
// Returns:
//   - *LessSet[T]: The difference of the set with the other set.
func (s *LessSet[T]) Difference(other *LessSet[T]) *LessSet[T] {
	if other == nil {
		return s
	}

	newElems := make([]T, 0)

	for _, e := range s.elems {
		if !other.HasElem(e) {
			newElems = append(newElems, e)
		}
	}

	return &LessSet[T]{
		elems: newElems,
	}
}

// SymmetricDifference returns the symmetric difference of the set with another set.
//
// Parameters:
//   - other: The other set.
//
// Returns:
//   - *LessSet[T]: The symmetric difference of the set with the other set.
func (s *LessSet[T]) SymmetricDifference(other *LessSet[T]) *LessSet[T] {
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

	elems := make([]T, len(diff1))
	copy(elems, diff1)

	for _, e := range diff2 {
		pos, ok := uc.Find(elems, e)
		if ok {
			continue
		}

		elems = slices.Insert(elems, pos, e)
	}

	return &LessSet[T]{
		elems: elems,
	}
}

// IsSubset checks if the set is a subset of another set.
//
// Parameters:
//   - other: The other set to check.
//
// Returns:
//   - bool: True if the set is a subset of the other set, false otherwise.
func (s *LessSet[T]) IsSubset(other *LessSet[T]) bool {
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
func (s *LessSet[T]) Clear() {
	for i := 0; i < len(s.elems); i++ {
		s.elems[i] = *new(T)
	}

	s.elems = s.elems[:0]
}

// String returns a string representation of the set.
//
// Returns:
//   - string: The string representation of the set.
func (s *LessSet[T]) String() string {
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
func (s *LessSet[T]) Equals(other uc.Objecter) bool {
	if other == nil {
		return false
	}

	otherSet, ok := other.(*LessSet[T])
	if !ok {
		return false
	}

	if len(s.elems) != len(otherSet.elems) {
		return false
	}

	for _, k := range s.elems {
		if !otherSet.HasElem(k) {
			return false
		}
	}

	return true
}

// Copy returns a copy of the set.
//
// Returns:
//   - *LessSet[T]: A copy of the set.
func (s *LessSet[T]) Copy() uc.Objecter {
	newElems := make([]T, len(s.elems))
	copy(newElems, s.elems)

	return &LessSet[T]{
		elems: newElems,
	}
}

// Slice returns a slice of the elements in the set.
//
// Returns:
//   - []T: A slice of the elements in the set.
func (s *LessSet[T]) Slice() []T {
	return s.elems
}

// Iterator returns an iterator for the set.
//
// Returns:
//   - ui.Iterater[T]: An iterator for the set.
func (s *LessSet[T]) Iterator() ui.Iterater[T] {
	return ui.NewSimpleIterator(s.elems)
}

// NewLessSet creates a new LessSet.
//
// Returns:
//   - *LessSet: A new LessSet.
func NewLessSet[T uc.Comparer](elems []T) *LessSet[T] {
	elems = us.UniquefyEquals(elems)
	uc.Sort(elems)

	return &LessSet[T]{
		elems: elems,
	}
}

func (s *LessSet[T]) Find(elem T) (int, bool) {
	return uc.Find(s.elems, elem)
}
