package Common

import (
	"fmt"
	"slices"
)

// Copier is an interface that provides a method to create a copy of an element.
type Copier interface {
	// Copy creates a copy of the element.
	//
	// Returns:
	//   - Copier: The copy of the element.
	Copy() Copier
}

// Cleaner is an interface that provides a method to remove all the elements
// from a data structure.
type Cleaner interface {
	// Clean removes all the elements from the data structure.
	Clean()
}

// Clean removes all the elements from the data structure by calling the Clean method if the
// element implements the Cleaner interface.
//
// Parameters:
//   - cleaner: The data structure to clean.
func Clean(elem any) {
	target, ok := elem.(Cleaner)
	if ok {
		target.Clean()
	}
}

// Find is a function that searches for an element in a slice and returns its index.
//
// Parameters:
//   - slice: The slice to search in.
//   - elem: The element to search for.
//
// Returns:
//   - int: The index of the element in the slice. -1 if there is at least one non-
//     comparable element.
//   - bool: A flag indicating if the element was found.
//
// Behaviors:
//   - This function uses the Binary Search algorithm to find the element.
//   - The slice must be sorted in ascending order.
//   - If the element is not found, the index is the position where the element should be inserted
//     to maintain the order of the slice.
func Find[T Comparer](S []T, elem T) (int, bool) {
	if len(S) == 0 {
		return 0, false
	}

	l, r := 0, len(S)-1

	for l <= r {
		middle := (l + r) / 2

		evaluation, ok := S[middle].Compare(elem)
		if !ok {
			return -1, false
		}

		if evaluation < 0 {
			l = middle + 1
		} else if evaluation > 0 {
			r = middle - 1
		} else {
			return middle, true
		}
	}

	return l, false
}

// StableSort is a function that sorts a slice of elements in ascending order while preserving
// the order of equal elements.
//
// Parameters:
//   - slice: The slice to sort.
//   - isAsc: A flag indicating if the sort is in ascending order.
//
// Returns:
//   - error: An error of type *ErrNotComparable if the elements are not comparable.
//
// Behaviors:
//   - This function uses the Merge Sort algorithm to sort the slice.
//   - The elements in the slice must implement the Comparer interface.
func StableSort[T Comparer](S []T, isAsc bool) error {
	var sortFunc func(a, b T) int
	var err error

	if isAsc {
		sortFunc = func(a, b T) int {
			if err != nil {
				return -1
			}

			res, ok := a.Compare(b)
			if !ok {
				err = NewErrNotComparable(a, b)
				return -1
			}

			return res
		}
	} else {
		sortFunc = func(a, b T) int {
			if err != nil {
				return -1
			}

			res, ok := b.Compare(a)
			if !ok {
				err = NewErrNotComparable(b, a)
				return -1
			}

			if res == 0 {
				return 0
			} else {
				return -res
			}
		}
	}

	slices.SortStableFunc(S, sortFunc)

	return err
}

// Sort is a function that sorts a slice of elements in ascending order.
//
// Parameters:
//   - slice: The slice to sort.
//
// Behaviors:
//   - This function uses the Quick Sort algorithm to sort the slice.
//   - The elements in the slice must implement the Comparer interface.
func Sort[T Comparer](S []T) {
	if len(S) < 2 {
		return
	}

	sortQuick(S, 0, len(S)-1)
}

// sortQuick is a helper function that sorts a slice of elements using the Quick Sort algorithm.
//
// Parameters:
//   - S: The slice to sort.
//   - l: The left index of the slice.
//   - r: The right index of the slice.
//
// Returns:
//   - error: An error of type *ErrNotComparable if the elements are not comparable.
//
// Behaviors:
//   - This function sorts the slice in ascending order.
func sortQuick[T Comparer](S []T, l, r int) error {
	if l < r {
		p, err := partition(S, l, r)
		if err != nil {
			return err
		}

		err = sortQuick(S, l, p-1)
		if err != nil {
			return err
		}

		err = sortQuick(S, p+1, r)
		if err != nil {
			return err
		}
	}

	return nil
}

// partition is a helper function that partitions a slice of elements for the Quick Sort algorithm.
//
// Parameters:
//   - S: The slice to partition.
//   - l: The left index of the slice.
//   - r: The right index of the slice.
//
// Returns:
//   - int: The index of the pivot element.
//   - error: An error of type *ErrNotComparable if the elements are not comparable.
func partition[T Comparer](S []T, l, r int) (int, error) {
	pivot := S[r]
	i := l

	for j := l; j < r; j++ {
		res, ok := S[j].Compare(pivot)
		if !ok {
			return 0, NewErrNotComparable(S[j], pivot)
		}

		if res < 0 {
			S[i], S[j] = S[j], S[i]
			i++
		}
	}

	S[i], S[r] = S[r], S[i]
	return i, nil
}

// Runer is an interface that provides a method to get the runes of a string.
type Runer interface {
	// Runes returns the runes of the string.
	//
	// Returns:
	//   - []rune: The runes of the string.
	Runes() []rune
}

// RunesOf returns the runes of a string. If the string implements the Runer interface,
// the Runes method is called. Otherwise, the string is converted to a slice of runes.
//
// Parameters:
//   - str: The string to get the runes from.
//
// Returns:
//   - []rune: The runes of the string.
func RunesOf(str any) []rune {
	if str == nil {
		return nil
	}

	switch str := str.(type) {
	case Runer:
		return str.Runes()
	case []rune:
		return str
	case string:
		return []rune(str)
	case []byte:
		return []rune(string(str))
	case error:
		return []rune(str.Error())
	case fmt.Stringer:
		return []rune(str.String())
	default:
		return []rune(fmt.Sprintf("%v", str))
	}
}

// Hashable is an interface that provides a method to get the hash code of an element.
type Hashable interface {
	// Hash returns the hash code of the element.
	//
	// Returns:
	//   - int: The hash code of the element.
	Hash() int
}

// HashCode returns the hash code of an element. If the element implements the Hashable interface,
// the Hash method is called. Otherwise, the hash code is calculated using the default hash function.
//
// Parameters:
//   - elem: The element to get the hash code from.
//
// Returns:
//   - int: The hash code of the element.
func HashCode(elem any) int {
	if elem == nil {
		return 0
	}

	switch elem := elem.(type) {
	case Hashable:
		return elem.Hash()
	case int:
		return elem
	case int8:
		return int(elem)
	case int16:
		return int(elem)
	case int32:
		return int(elem)
	case int64:
		return int(elem)
	case uint:
		return int(elem)
	case uint8:
		return int(elem)
	case uint16:
		return int(elem)
	case uint32:
		return int(elem)
	case uint64:
		return int(elem)
	case float32:
		return int(elem)
	case float64:
		return int(elem)
	case bool:
		if elem {
			return 1
		} else {
			return 0
		}
	case string:
		return hashString(elem)
	case []byte:
		return hashBytes(elem)
	case error:
		return hashString(elem.Error())
	case fmt.Stringer:
		return hashString(elem.String())
	default:
		return hashString(fmt.Sprintf("%v", elem))
	}
}

// hashString is a helper function that calculates the hash code of a string.
//
// Parameters:
//   - str: The string to calculate the hash code from.
//
// Returns:
//   - int: The hash code of the string.
func hashString(str string) int {
	hash := 0

	for _, r := range str {
		hash = 31*hash + int(r)
	}

	return hash
}

// hashBytes is a helper function that calculates the hash code of a byte slice.
//
// Parameters:
//   - bytes: The byte slice to calculate the hash code from.
//
// Returns:
//   - int: The hash code of the byte slice.
func hashBytes(bytes []byte) int {
	hash := 0

	for _, b := range bytes {
		hash = 31*hash + int(b)
	}

	return hash
}
