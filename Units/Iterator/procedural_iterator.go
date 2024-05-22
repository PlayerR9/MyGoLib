// Package Iterators provides a set of types that allow iterating over
// collections of elements in a generic and procedural manner.
package Iterators

import (
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

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

// IteratorFromIterator creates a new iterator over a collection of iterators
// of type Iterater[T].
// It uses the input iterator to iterate over the collection of iterators and
// return the elements from each iterator in turn.
//
// Parameters:
//   - source: The iterator over the collection of iterators to iterate over.
//   - f: The transition function that takes an element of type E and returns
//     an iterator.
//
// Return:
//   - Iterater[T]: The new iterator over the collection of elements.
//   - error: An error of type *ers.ErrInvalidParameter if the transition function
//     is nil.
func NewProceduralIterator[E, T any](source Iterater[E], f func(E) Iterater[T]) (*ProceduralIterator[E, T], error) {
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
