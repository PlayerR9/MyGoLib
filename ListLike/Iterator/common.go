package Iterators

import (
	intf "github.com/PlayerR9/MyGoLib/Units/Common"
)

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

// Iterable is an interface that defines a method to get an iterator over a
// collection of elements of type T. It is implemented by data structures that
// can be iterated over.
type Iterable[T any] interface {
	Iterator() Iterater[T]
}

// IteratorFromSlice creates a new iterator over a slice of elements of type T.
//
// Parameters:
//   - values: The slice of elements to iterate over.
//
// Return:
//   - Iterater[T]: A new iterator over the given slice of elements.
func IteratorFromSlice[T any](values []T) Iterater[T] {
	return &GenericIterator[T]{
		values: &values,
		index:  0,
	}
}

// IteratorFromSlicer creates a new iterator over a data structure that implements
// the Slicer interface. It uses the Slice method of the data structure to get the
// slice of elements to iterate over.
//
// Parameters:
//   - slicer: The data structure that implements the Slicer interface.
//
// Return:
//   - Iterater[T]: A new iterator over the slice of elements returned by the slicer.
func IteratorFromSlicer[T any](slicer intf.Slicer[T]) Iterater[T] {
	elements := slicer.Slice()

	return &GenericIterator[T]{
		values: &elements,
		index:  0,
	}
}

// IteratorFromValues creates a new iterator over a variadic list of elements of
// type T.
//
// Parameters:
//   - values: The variadic list of elements to iterate over.
//
// Return:
//   - Iterater[T]: The new iterator over the given elements.
func IteratorFromValues[T any](values ...T) Iterater[T] {
	// Create a copy of the values slice.
	valuesCopy := make([]T, len(values))
	copy(valuesCopy, values)

	return &GenericIterator[T]{
		values: &valuesCopy,
		index:  0,
	}
}
