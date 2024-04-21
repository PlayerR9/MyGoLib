package Iterators

import (
	itf "github.com/PlayerR9/MyGoLib/Common/Interfaces"
	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

// Iterater is an interface that defines methods for an iterator over a collection of
// elements of type T. It includes methods to get the next element, get the current
// element, and restart the iterator.
type Iterater[T any] interface {
	// The Next method advances the iterator to the next element in the collection.
	// It returns true if there is a next element, otherwise false.
	Next() bool

	// The Value method returns the current element in the collection.
	// It should be called after Next to get the current element.
	//
	// If the iterator is exhausted, it will panic.
	Value() (T, error)

	// The Restart method resets the iterator to the beginning of the collection.
	Restart()
}

// Iterable is an interface that defines a method to get an iterator over a collection
// of elements of type T. It is implemented by data structures that can be iterated over.
type Iterable[T any] interface {
	Iterator() Iterater[T]
}

// IteratorFromSlice creates a new iterator over a slice of elements of type T.
// It creates a shallow copy of the input slice to minimize side effects.
//
// Parameters:
//
//   - values: The slice of elements to iterate over.
//
// Return:
//
//   - Iterater[T]: A pointer to a new iterator over the given slice of elements.
func IteratorFromSlice[T any](values []T) Iterater[T] {
	return &GenericIterator[T]{
		values: &values,
		index:  -1,
	}
}

// IteratorFromSlicer creates a new iterator over a data structure that implements the
// Slicer interface. It uses the Slice method of the data structure to get the slice of
// elements to iterate over.
//
// Parameters:
//
//   - slicer: The data structure that implements the Slicer interface.
//
// Return:
//
//   - Iterater[T]: A pointer to a new iterator over the slice of elements returned by the
//     Slice method of the input slicer.
func IteratorFromSlicer[T any](slicer itf.Slicer[T]) Iterater[T] {
	elements := slicer.Slice()

	return &GenericIterator[T]{
		values: &elements,
		index:  -1,
	}
}

// IteratorFromValues creates a new iterator over a variadic list of elements of type T.
// It creates a shallow copy of the input variadic list to minimize side effects.
//
// Parameters:
//
//   - values: The variadic list of elements to iterate over.
//
// Return:
//
//   - Iterater[T]: A pointer to a new iterator over the given variadic list of elements.
func IteratorFromValues[T any](values ...T) Iterater[T] {
	// Create a copy of the values slice.
	valuesCopy := make([]T, len(values))
	copy(valuesCopy, values)

	return &GenericIterator[T]{
		values: &valuesCopy,
		index:  -1,
	}
}

// IteratorFromIterator creates a new iterator over a collection of iterators of type
// Iterater[T].
// It uses the input iterator to iterate over the collection of iterators and return
// the elements from each iterator in turn.
//
// Parameters:
//
//   - source: The iterator over the collection of iterators to iterate over.
//   - f: The transition function that takes an element of type E and returns an iterator
//
// Return:
//
//   - Iterater[T]: A pointer to a new iterator over the collection of iterators.
func IteratorFromIterator[E, T any](source Iterater[E], f func(E) Iterater[T]) (Iterater[T], error) {
	if f == nil {
		return nil, ers.NewErrNilParameter("f")
	}

	iter := &ProceduralIterator[E, T]{
		source:  source,
		current: nil,
	}

	iter.transition = f

	return iter, nil
}
