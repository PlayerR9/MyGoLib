// Package Iterators provides a set of types that allow iterating over
// collections of elements in a generic and procedural manner.
package Iterators

// GenericIterator is a struct that allows iterating over a slice of
// elements of any type.
type GenericIterator[T any] struct {
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
func (iter *GenericIterator[T]) Consume() (T, error) {
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
func (iter *GenericIterator[T]) Restart() {
	iter.index = 0
}
