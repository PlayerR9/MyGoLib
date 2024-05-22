package Sorting

import (
	uc "github.com/PlayerR9/MyGoLib/Units/Common"
)

// Slice is a slice that uses the Compare method to compare elements.
type Slice[T uc.Comparer[T]] struct {
	// elems is the slice of elements in the sorted slice.
	elems []T
}

// NewSlice creates a new sorted slice.
//
// Parameters:
//   - elems: The elements to add to the sorted slice.
//
// Returns:
//   - *Slice[T]: The new sorted slice.
func NewSlice[T uc.Comparer[T]](elems []T) *Slice[T] {
	if len(elems) == 0 {
		return &Slice[T]{elems: make([]T, 0)}
	} else {
		uc.Sort(elems)

		return &Slice[T]{elems: elems}
	}
}

// Insert inserts an element into the sorted slice.
//
// Parameters:
//   - elem: The element to insert.
//   - force: If true, the element is inserted even if it is already in the slice.
//
// Returns:
//   - int: The index where the element was inserted.
//
// Behaviors:
//   - The function inserts the element in the correct position to maintain the order.
func (s *Slice[T]) Insert(elem T, force bool) int {
	if len(s.elems) == 0 {
		s.elems = append(s.elems, elem)
		return 0
	}

	pos, ok := uc.Find(s.elems, elem)
	if !ok || force {
		s.elems = append(s.elems[:pos], append([]T{elem}, s.elems[pos:]...)...)
	}

	return pos
}

// TryInsert tries to insert an element into the sorted slice.
//
// Parameters:
//   - elem: The element to insert.
//
// Returns:
//   - int: The index where the element would be inserted.
//   - bool: True if the element is already in the slice, false otherwise.
func (s *Slice[T]) TryInsert(elem T) (int, bool) {
	return uc.Find(s.elems, elem)
}

// Find is the same as Find but uses the Compare method of the elements.
//
// Parameters:
//   - elem: element to find.
//
// Returns:
//   - int: index of the first occurrence of the element or -1 if not found.
//
// Behaviors:
//   - The values must be sorted in ascending order for the Compare method to work.
func (s *Slice[T]) Find(elem T) int {
	if len(s.elems) == 0 {
		return -1
	}

	pos, ok := uc.Find(s.elems, elem)
	if ok {
		return pos
	}

	return -1
}

// Uniquefy is the same as Uniquefy but uses the Compare method of the elements.
//
// Parameters:
//   - S: slice of elements.
//
// Returns:
//   - []T: slice of elements with duplicates removed.
//
// Behavior:
//   - The function preserves the order of the elements in the slice.
//   - The values must be sorted in ascending order for the Compare method to work.
func (s *Slice[T]) Uniquefy() *Slice[T] {
	if len(s.elems) < 2 {
		return &Slice[T]{
			elems: s.elems,
		}
	}

	newS := &Slice[T]{
		elems: make([]T, 0),
	}

	for _, e := range s.elems {
		newS.Insert(e, false)
	}

	return newS
}

// MergeUnique merges two slices and removes duplicate elements.
//
// Parameters:
//   - S1: first slice of elements.
//   - S2: second slice of elements.
//
// Returns:
//   - []T: slice of elements with duplicates removed.
//
// Behaviors:
//   - The function does not preserve the order of the elements in the slices.
func (s *Slice[T]) MergeUnique(other *Slice[T]) *Slice[T] {
	if other == nil {
		return &Slice[T]{
			elems: s.elems,
		}
	}

	newS := &Slice[T]{
		elems: make([]T, len(s.elems)),
	}
	copy(newS.elems, s.elems)

	for _, e := range other.elems {
		newS.Insert(e, false)
	}

	return newS
}

// IndexOfDuplicate returns the index of the first duplicate element in the slice.
//
// Parameters:
//   - S: slice of elements.
//
// Returns:
//   - int: index of the first duplicate element or -1 if there are no duplicates.
func (s *Slice[T]) IndexOfDuplicate() int {
	if len(s.elems) < 2 {
		return -1
	}

	for i := 0; i < len(s.elems)-1; i++ {
		if s.elems[i].Compare(s.elems[i+1]) == 0 {
			return i
		}
	}

	return -1
}

// computeLPSArray is a helper function that computes the Longest Prefix
// Suffix (LPS) array for the Knuth-Morris-Pratt algorithm.
//
// Parameters:
//   - subS: The subslice to compute the LPS array for.
//   - lps: The LPS array to store the results in.
//
// Behavior:
//   - The function modifies the lps array in place.
//   - The lps array is initialized with zeros.
//   - The lps array is used to store the length of the longest prefix
//     that is also a suffix for each index in the subslice.
//   - The first element of the lps array is always 0.
func (s *Slice[T]) computeLPSArray(lps []int) {
	length := 0
	i := 1
	lps[0] = 0 // lps[0] is always 0

	// the loop calculates lps[i] for i = 1 to len(subS)-1
	for i < len(s.elems) {
		if s.elems[i].Compare(s.elems[length]) == 0 {
			length++
			lps[i] = length
			i++
		} else {
			if length != 0 {
				length = lps[length-1]
			} else {
				lps[i] = 0
				i++
			}
		}
	}
}

// FindSubBytesFrom finds the first occurrence of a subslice in a byte
// slice starting from a given index.
//
// Parameters:
//   - S: The byte slice to search in.
//   - subS: The byte slice to search for.
//   - at: The index to start searching from.
//
// Returns:
//   - int: The index of the first occurrence of the subslice.
//
// Behavior:
//   - The function uses the Knuth-Morris-Pratt algorithm to find the subslice.
//   - If S or subS is empty, the function returns -1.
//   - If the subslice is not found, the function returns -1.
//   - If at is negative, it is set to 0.
func (s *Slice[T]) FindSubsliceFrom(other *Slice[T], at int) int {
	if other == nil || len(other.elems) == 0 || len(s.elems) == 0 || at+len(other.elems) > len(s.elems) {
		return -1
	}

	if at < 0 {
		at = 0
	}

	lps := make([]int, len(other.elems))
	other.computeLPSArray(lps)

	i := at
	j := 0
	for i < len(s.elems) {
		if s.elems[i].Compare(other.elems[j]) == 0 {
			i++
			j++
		}

		if j == len(other.elems) {
			return i - j
		} else if i < len(s.elems) && s.elems[i].Compare(other.elems[j]) != 0 {
			if j != 0 {
				j = lps[j-1]
			} else {
				i = i + 1
			}
		}
	}

	return -1
}

// Difference returns the elements that are in S1 but not in S2.
//
// Parameters:
//   - S1: The first slice of elements.
//   - S2: The second slice of elements.
func (s *Slice[T]) Difference(other *Slice[T]) *Slice[T] {
	if other == nil {
		return s
	}

	newS := &Slice[T]{
		elems: make([]T, 0),
	}

	for _, e := range s.elems {
		pos := other.Find(e)
		if pos == -1 {
			newS.Insert(e, false)
		}
	}

	return newS
}
