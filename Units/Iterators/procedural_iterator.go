package Iterators

import ue "github.com/PlayerR9/MyGoLib/Units/errors"

// ProceduralIterator is a struct that allows iterating over a collection of
// iterators of type Iterater[T].
type ProceduralIterator[E Iterable[T], T any] struct {
	// source is the iterator over the collection of iterators.
	source Iterater[E]

	// iter is the iterator in the collection.
	iter Iterater[T]
}

// Size implements the Iterater interface.
//
// Size is evaluated by summing the sizes of the current iterator and the source
// iterator. Of course, this is just an estimate of the total number of elements
// in the collection.
func (pi *ProceduralIterator[E, T]) Size() (count int) {
	if pi.iter != nil {
		count += pi.iter.Size()
	}

	if pi.source != nil {
		count += pi.source.Size()
	}

	return
}

// Consume implements the Iterater interface.
func (pi *ProceduralIterator[E, T]) Consume() (T, error) {
	if pi.iter == nil {
		iter, err := pi.source.Consume()
		if err != nil {
			return *new(T), err
		}

		pi.iter = iter.Iterator().(*SimpleIterator[T])
	}

	var val T
	var err error

	for {
		val, err = pi.iter.Consume()
		if err == nil {
			break
		}

		ok := ue.Is[*ErrExhaustedIter](err)
		if !ok {
			return *new(T), err
		}

		iter, err := pi.source.Consume()
		if err != nil {
			return *new(T), err
		}

		pi.iter = iter.Iterator().(*SimpleIterator[T])
	}

	return val, nil
}

// Restart implements the Iterater interface.
func (pi *ProceduralIterator[E, T]) Restart() {
	pi.iter = nil
	pi.source.Restart()
}

// IteratorFromIterator creates a new iterator over a collection of iterators
// of type Iterater[T].
// It uses the input iterator to iterate over the collection of iterators and
// return the elements from each iterator in turn.
//
// Parameters:
//   - source: The iterator over the collection of iterators to iterate over.
//
// Return:
//   - *ProceduralIterator[E, T]: The new iterator over the collection of elements.
//     Nil if source is nil.
func NewProceduralIterator[E Iterable[T], T any](source Iterater[E]) *ProceduralIterator[E, T] {
	if source == nil {
		return nil
	}

	pi := &ProceduralIterator[E, T]{
		source: source,
		iter:   nil,
	}

	return pi
}
