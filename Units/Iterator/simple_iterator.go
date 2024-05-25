// Package Iterators provides a set of types that allow iterating over
// collections of elements in a generic and procedural manner.
package Iterators

// SimpleIterator is a struct that allows iterating over a slice of
// elements of any type.
type SimpleIterator[T any] struct {
	// The slice of elements.
	values *[]T

	// The current index in the values slice.
	index int // 0 means not initialized
}

// Consume is a method of the GenericIterator type that advances the
// iterator to the next element in the collection and returns the current
// element.
//
// Errors:
//   - *ErrNotInitialized: If the iterator is not initialized.
//   - *ErrExhaustedIter: If the iterator is exhausted.
//
// Returns:
//   - T: The current element in the collection.
//   - error: An error if it is not possible to consume the next element.
func (iter *SimpleIterator[T]) Consume() (T, error) {
	if iter.values == nil {
		return *new(T), NewErrNotInitialized()
	} else if iter.index >= len(*iter.values) {
		return *new(T), NewErrExhaustedIter()
	}

	value := (*iter.values)[iter.index]

	iter.index++

	return value, nil
}

// Restart is a method of the GenericIterator type that resets the iterator to the
// beginning of the collection.
func (iter *SimpleIterator[T]) Restart() {
	iter.index = 0
}

// NewSimpleIterator creates a new iterator over a slice of elements of type T.
//
// Parameters:
//   - values: The slice of elements to iterate over.
//
// Return:
//   - *GenericIterator[T]: A new iterator over the given slice of elements.
//
// Behaviors:
//   - If values is nil, the iterator is initialized with an empty slice.
//   - Modifications to the slice of elements after creating the iterator will
//     affect the values seen by the iterator.
func NewSimpleIterator[T any](values []T) *SimpleIterator[T] {
	if len(values) == 0 {
		values = make([]T, 0)
	}

	return &SimpleIterator[T]{
		values: &values,
		index:  0,
	}
}
