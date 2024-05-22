package Iterators

// Slicer is an interface that provides a method to convert a data structure to a slice.
type Slicer[T any] interface {
	// The Slice method returns a slice containing all the elements in the data structure.
	Slice() []T

	Iterable[T]
}

// SliceOf converts any type to a slice of elements of the same type.
//
// Parameters:
//   - elem: The element to convert to a slice.
//
// Returns:
//   - []T: The slice representation of the element.
//
// Behaviors:
//   - Nil elements are converted to nil slices.
//   - Slice elements are returned as is.
//   - Slicer elements have their Slice method called.
//   - Other elements are converted to slices containing a single element.
func SliceOf[T any](elem any) []T {
	if elem == nil {
		return nil
	}

	switch elem := elem.(type) {
	case []T:
		return elem
	case Slicer[T]:
		return elem.Slice()
	default:
		return []T{elem.(T)}
	}
}

// Iterable is an interface that defines a method to get an iterator over a
// collection of elements of type T. It is implemented by data structures that
// can be iterated over.
type Iterable[T any] interface {
	Iterator() Iterater[T]
}

// IteratorOf converts any type to an iterator over elements of the same type.
//
// Parameters:
//   - elem: The element to convert to an iterator.
//
// Returns:
//   - Iterater[T]: The iterator over the element.
//
// Behaviors:
//   - IF elem is nil, an empty iterator is returned.
//   - IF elem -implements-> Iterater[T], the element is returned as is.
//   - IF elem -implements-> Iterable[T], the element's Iterator method is called.
//   - IF elem -implements-> []T, a new iterator over the slice is created.
//   - ELSE, a new iterator over a single-element collection is created.
func IteratorOf[T any](elem any) Iterater[T] {
	if elem == nil {
		var builder Builder[T]

		return builder.Build()
	}

	switch elem := elem.(type) {
	case Iterater[T]:
		return elem
	case Iterable[T]:
		return elem.Iterator()
	case []T:
		return &GenericIterator[T]{
			values: &elem,
			index:  0,
		}
	default:
		return &GenericIterator[T]{
			values: &[]T{elem.(T)},
			index:  0,
		}
	}
}

// Iterater is an interface that defines methods for an iterator over a
// collection of elements of type T.
type Iterater[T any] interface {
	// The Consume method advances the iterator to the next element in the
	// collection and returns the current element.
	//
	// Returns:
	//  - T: The current element in the collection.
	//  - error: An error if the iterator is exhausted.
	Consume() (T, error)

	// The Restart method resets the iterator to the beginning of the
	// collection.
	Restart()
}
