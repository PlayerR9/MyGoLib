package Common

import "fmt"

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
//   - int: The index of the element in the slice.
//   - bool: A flag indicating if the element was found.
//
// Behaviors:
//   - This function uses the Binary Search algorithm to find the element.
//   - The slice must be sorted in ascending order.
//   - If the element is not found, the index is the position where the element should be inserted
//     to maintain the order of the slice.
func Find[T Comparer[T]](S []T, elem T) (int, bool) {
	if len(S) == 0 {
		return 0, false
	}

	l, r := 0, len(S)-1

	for l <= r {
		middle := (l + r) / 2

		evaluation := S[middle].Compare(elem)
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

// Sort is a function that sorts a slice of elements in ascending order.
//
// Parameters:
//   - slice: The slice to sort.
//
// Behaviors:
//   - This function uses the Quick Sort algorithm to sort the slice.
//   - The elements in the slice must implement the Comparer interface.
func Sort[T Comparer[T]](S []T) {
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
// Behaviors:
//   - This function sorts the slice in ascending order.
func sortQuick[T Comparer[T]](S []T, l, r int) {
	if l < r {
		p := partition(S, l, r)
		sortQuick(S, l, p-1)
		sortQuick(S, p+1, r)
	}
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
func partition[T Comparer[T]](S []T, l, r int) int {
	pivot := S[r]
	i := l

	for j := l; j < r; j++ {
		if S[j].Compare(pivot) < 0 {
			S[i], S[j] = S[j], S[i]
			i++
		}
	}

	S[i], S[r] = S[r], S[i]
	return i
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
