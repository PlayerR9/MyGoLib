package Iterators

import (
	"errors"
)

// GenericIterator is a struct that allows iterating over a slice of elements of any
// type.
type GenericIterator[T any] struct {
	// The slice of elements.
	values *[]T

	// The current index in the values slice.
	index int // -1 means not initialized
}

// Next is a method of the GenericIterator type that advances the iterator to the next
// element in the collection and returns true if there is a next element, otherwise false.
//
// Returns:
//
//   - bool: True if there is a next element, otherwise false.
func (iter *GenericIterator[T]) Next() (hasNext bool) {
	if iter.values == nil {
		return false
	} else if iter.index+1 >= len(*iter.values) {
		return false
	}

	iter.index++

	return true
}

// Value is a method of the GenericIterator type that returns the current element in
// the collection.
// It should be called after Next to get the current element.
//
// Panics with *ErrCallFailed if the iterator is exhausted, if it is called before the
// iterator is initialized, or if it is called before Next.
//
// Returns:
//
//   - T: The current element in the collection.
func (iter *GenericIterator[T]) Value() (T, error) {
	if iter.values == nil {
		return *new(T), errors.New("iterator was never initialized")
	} else if iter.index == -1 {
		return *new(T), errors.New("Next must be called before Value")
	} else if iter.index >= len(*iter.values) {
		return *new(T), errors.New("value called on exhausted iter")
	}

	vals := *iter.values

	return vals[iter.index], nil
}

// ValueNoErr is a method of the GenericIterator type that returns the current element in
// the collection without an error.
// It should be called after Next to get the current element.
//
// Panics with *ErrCallFailed if the iterator is exhausted, if it is called before the
// iterator is initialized, or if it is called before Next.
//
// Returns:
//
//   - T: The current element in the collection.
func (iter *GenericIterator[T]) ValueNoErr() T {
	if iter.values == nil {
		panic(errors.New("iterator was never initialized"))
	} else if iter.index == -1 {
		panic(errors.New("Next must be called before Value"))
	} else if iter.index >= len(*iter.values) {
		panic(errors.New("value called on exhausted iter"))
	}

	vals := *iter.values

	return vals[iter.index]
}

// Restart is a method of the GenericIterator type that resets the iterator to the
// beginning of the collection.
func (iter *GenericIterator[T]) Restart() {
	iter.index = -1
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

// Next is a method of the ProceduralIterator type that advances the iterator to the
// next element in the collection and returns true if there is a next element, otherwise
// false.
//
// Panics with *ErrCallFailed if the iterator is exhausted or if it is called before
// the iterator is initialized.
//
// Returns:
//
//   - bool: True if there is a next element, otherwise false.
func (iter *ProceduralIterator[E, T]) Next() bool {
	if iter.source == nil {
		return false
	}

	if iter.current == nil || !iter.current.Next() {
		if !iter.source.Next() {
			return false
		}

		val, err := iter.source.Value()
		if err != nil {
			panic(err)
		}

		iter.current = iter.transition(val)

		iter.current.Next()
	}

	return true
}

// Value is a method of the ProceduralIterator type that returns the current element in
// the collection.
// It should be called after Next to get the current element.
//
// Panics with *ErrCallFailed if the iterator is exhausted, if it is called before the
// iterator is initialized, or if it is called before Next.
//
// Returns:
//
//   - T: The current element in the collection.
func (iter *ProceduralIterator[E, T]) Value() (T, error) {
	if iter.current == nil {
		return *new(T), errors.New("Next must be called before Value")
	}

	return iter.current.Value()
}

// ValueNoErr is a method of the ProceduralIterator type that returns the current element
// in the collection without an error.
// It should be called after Next to get the current element.
//
// Panics with *ErrCallFailed if the iterator is exhausted, if it is called before the
// iterator is initialized, or if it is called before Next.
//
// Returns:
//
//   - T: The current element in the collection.
func (iter *ProceduralIterator[E, T]) ValueNoErr() T {
	if iter.current == nil {
		panic(errors.New("Next must be called before Value"))
	}

	val, err := iter.current.Value()
	if err != nil {
		panic(err)
	}

	return val
}

// Restart is a method of the ProceduralIterator type that resets the iterator to the
// beginning of the collection.
func (iter *ProceduralIterator[E, T]) Restart() {
	iter.current = nil
	iter.source.Restart()
}
