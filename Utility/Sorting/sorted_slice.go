package Sorting

import "slices"

// Slice is a slice that uses the Compare method to compare elements.
type Slice[T any] struct {
	// elems is the slice of elements in the sorted slice.
	elems []T

	// sf is the sort function used to compare elements.
	sf SortFunc[T]

	// isAsc is a flag indicating if the sort is in ascending order.
	isAsc bool
}

// NewSlice creates a new sorted slice.
//
// Parameters:
//   - elems: The elements to add to the sorted slice.
//
// Returns:
//   - *Slice[T]: The new sorted slice.
//
// Behaviors:
//   - Returns nil if the sort function is nil.
func NewSlice[T any](elems []T, sf SortFunc[T], isAsc bool) *Slice[T] {
	if sf == nil {
		return nil
	}

	Sort(elems, sf, isAsc)

	res := &Slice[T]{
		elems: elems,
		sf:    sf,
		isAsc: isAsc,
	}

	return res
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

	pos, ok := slices.BinarySearchFunc(s.elems, elem, s.sf)
	if !ok || force {
		s.elems = slices.Insert(s.elems, pos, elem)
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
	pos, ok := slices.BinarySearchFunc(s.elems, elem, s.sf)
	return pos, ok
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

	pos, ok := slices.BinarySearchFunc(s.elems, elem, s.sf)
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
func (s *Slice[T]) Uniquefy() {
	if len(s.elems) < 2 {
		return
	}

	for i := 0; i < len(s.elems)-1; i++ {
		elem := s.elems[i]

		top := i
		for ; top < len(s.elems)-1; top++ {
			res := s.sf(s.elems[top+1], elem)
			if res != 0 {
				break
			}
		}

		if top != i {
			s.elems = slices.Delete(s.elems, i+1, top+1)
		}
	}
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
//   - The function does preserve the order of the elements in the slice.
func (s *Slice[T]) MergeUnique(other *Slice[T]) *Slice[T] {
	newElems := make([]T, len(s.elems))
	copy(newElems, s.elems)

	if other != nil {
		for _, e := range other.elems {
			pos, ok := slices.BinarySearchFunc(newElems, e, s.sf)
			if !ok {
				newElems = slices.Insert(newElems, pos, e)
			}
		}
	}

	unique := &Slice[T]{
		elems: newElems,
		sf:    s.sf,
		isAsc: s.isAsc,
	}

	return unique
}

// IndexOfDuplicate returns the index of the first duplicate element in the slice.
//
// Returns:
//   - int: index of the first duplicate element.
//
// Behaviors:
//   - The function returns -1 if no duplicates are found.
func (s *Slice[T]) IndexOfDuplicate() int {
	if len(s.elems) < 2 {
		return -1
	}

	for i := 0; i < len(s.elems)-1; i++ {
		res := s.sf(s.elems[i], s.elems[i+1])
		if res == 0 {
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
		res := s.sf(s.elems[i], s.elems[length])

		if res == 0 {
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
		res := s.sf(s.elems[i], other.elems[j])

		if res == 0 {
			i++
			j++
		}

		if j == len(other.elems) {
			return i - j
		} else if i < len(s.elems) {
			res := s.sf(s.elems[i], other.elems[j])

			if res != 0 {
				if j != 0 {
					j = lps[j-1]
				} else {
					i = i + 1
				}
			}
		}
	}

	return -1
}

// Difference returns the elements that are in s but not in other.
//
// Parameters:
//   - other: The slice to compare with.
//
// Returns:
//   - *Slice[T]: The slice of elements that are in s but not in other.
func (s *Slice[T]) Difference(other *Slice[T]) *Slice[T] {
	var newElems []T

	if other == nil {
		newElems = make([]T, len(s.elems))
		copy(newElems, s.elems)
	} else {
		for _, e := range s.elems {
			pos := other.Find(e)
			if pos == -1 {
				newElems = append(newElems, e)
			}
		}
	}

	diff := &Slice[T]{
		elems: newElems,
		sf:    s.sf,
		isAsc: s.isAsc,
	}

	return diff
}
