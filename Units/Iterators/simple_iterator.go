// Package Iterators provides a set of types that allow iterating over
// collections of elements in a generic and procedural manner.
package Iterators

// SimpleIterator is a struct that allows iterating over a slice of
// elements of any type.
type SimpleIterator[T any] struct {
	// values is a slice of elements of type T.
	values *[]T

	// index is the current index of the iterator.
	// 0 means not initialized.
	index int
}

// Size implements the Iterater interface.
func (iter *SimpleIterator[T]) Size() int {
	if iter.values == nil {
		return 0
	}

	return len(*iter.values)
}

// Consume implements the Iterater interface.
func (iter *SimpleIterator[T]) Consume() (T, error) {
	if iter.index >= len(*iter.values) {
		return *new(T), NewErrExhaustedIter()
	}

	value := (*iter.values)[iter.index]

	iter.index++

	return value, nil
}

// Restart implements the Iterater interface.
func (iter *SimpleIterator[T]) Restart() {
	iter.index = 0
}

// NewSimpleIterator creates a new iterator over a slice of elements of type T.
//
// Parameters:
//   - values: The slice of elements to iterate over.
//
// Return:
//   - *SimpleIterator[T]: A new iterator over the given slice of elements.
//
// Behaviors:
//   - If values is nil, the iterator is initialized with an empty slice.
//   - Modifications to the slice of elements after creating the iterator will
//     affect the values seen by the iterator.
func NewSimpleIterator[T any](values []T) *SimpleIterator[T] {
	if len(values) == 0 {
		values = make([]T, 0)
	}

	si := &SimpleIterator[T]{
		values: &values,
		index:  0,
	}
	return si
}
