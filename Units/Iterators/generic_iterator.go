package Iterators

import (
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

type GenericIterator[E Iterater[T], T any] struct {
	source Iterater[E]
	iter   Iterater[T]
}

// Size implements the Iterater interface.
//
// Size is evaluated by summing the sizes of the current iterator and the source
// iterator. Of course, this is just an estimate of the total number of elements
// in the collection.
func (gi *GenericIterator[E, T]) Size() (count int) {
	if gi.iter != nil {
		count += gi.iter.Size()
	}

	if gi.source != nil {
		count += gi.source.Size()
	}

	return
}

// Consume implements the Iterater interface.
func (gi *GenericIterator[E, T]) Consume() (T, error) {
	if gi.iter == nil {
		iter, err := gi.source.Consume()
		if err != nil {
			return *new(T), err
		}

		gi.iter = iter
	}

	var val T
	var err error

	for {
		val, err = gi.iter.Consume()
		if err == nil {
			break
		}

		ok := ue.Is[*ErrExhaustedIter](err)
		if !ok {
			return *new(T), err
		}

		iter, err := gi.source.Consume()
		if err != nil {
			return *new(T), err
		}

		gi.iter = iter
	}

	return val, nil
}

// Restart implements the Iterater interface.
func (gi *GenericIterator[E, T]) Restart() {
	gi.source.Restart()
	gi.iter = nil
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
//   - *DynamicIterator[E, T]: The new iterator. Nil if f is nil.
func NewGenericIterator[E Iterater[T], T any](source Iterater[E]) *GenericIterator[E, T] {
	if source == nil {
		return nil
	}

	iter := &GenericIterator[E, T]{
		source: source,
		iter:   nil,
	}

	return iter
}
