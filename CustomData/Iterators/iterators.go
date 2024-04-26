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

// ProceduralIterator is a struct that allows iterating over a collection of iterators
// of type Iterater[T].
// The major difference between this and the GenericIterator is that this iterator is
// designed to iterate over a collection of elements in a progressive manner; reducing
// the need to store the entire collection in memory.
type ProceduralIterator[E, T any] struct {
	// The iterator over the collection of iterators.
	source Iterater[E]

	// The current iterator in the collection.
	current Iterater[T]

	// Transition function between iterators.
	transition func(E) Iterater[T]
}

// Consume is a method of the ProceduralIterator type that advances the
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
func (iter *ProceduralIterator[E, T]) Consume() (T, error) {
	if iter.source == nil {
		return *new(T), NewErrNotInitialized()
	}

	if iter.current != nil {
		next2, err := iter.current.Consume()
		if err == nil {
			return next2, nil
		}
	}

	next1, err := iter.source.Consume()
	if err != nil {
		return *new(T), NewErrExhaustedIter()
	}

	iter.current = iter.transition(next1)

	return iter.current.Consume()
}

// Restart is a method of the ProceduralIterator type that resets the
// iterator to the beginning of the collection.
func (iter *ProceduralIterator[E, T]) Restart() {
	iter.current = nil
	iter.source.Restart()
}
