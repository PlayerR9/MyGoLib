// Package Iterators provides a set of types that allow iterating over
// collections of elements in a generic and procedural manner.
package Iterators

// DynamicIterator is a struct that allows iterating over a collection of iterators
// of type Iterater[T].
// The major difference between this and the GenericIterator is that this iterator is
// designed to iterate over a collection of elements in a progressive manner; reducing
// the need to store the entire collection in memory.
type DynamicIterator[E, T any] struct {
	// source is the iterator over the collection of iterators.
	source Iterater[E]

	// current is the current iterator in the collection.
	current *SimpleIterator[T]

	// transition is the transition function that takes an element of type E and
	// returns an iterator.
	transition func(E) *SimpleIterator[T]
}

// Size implements the Iterater interface for the DynamicIterator type.
//
// Size is evaluated by summing the sizes of the current iterator and the source
// iterator. Of course, this is just an estimate of the total number of elements
// in the collection.
func (iter *DynamicIterator[E, T]) Size() (count int) {
	if iter.current != nil {
		count += iter.current.Size()
	}

	if iter.source != nil {
		count += iter.source.Size()
	}

	return
}

// Consume implements the Iterater interface.
func (iter *DynamicIterator[E, T]) Consume() (T, error) {
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

	elem, err := iter.current.Consume()
	return elem, err
}

// Restart implements the Iterater interface.
func (iter *DynamicIterator[E, T]) Restart() {
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
//   - *DynamicIterator[E, T]: The new iterator. Nil if f is nil.
func NewDynamicIterator[E, T any](source Iterater[E], f func(E) *SimpleIterator[T]) *DynamicIterator[E, T] {
	if f == nil {
		return nil
	}

	iter := &DynamicIterator[E, T]{
		source:  source,
		current: nil,
	}

	iter.transition = f

	return iter
}
