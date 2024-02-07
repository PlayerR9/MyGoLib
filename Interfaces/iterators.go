package Interfaces

import (
	"errors"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
	"github.com/markphelps/optional"
)

// Slicer is an interface that defines a method to convert a data structure to a slice.
// It is implemented by data structures that can be converted to a slice of elements of
// the same type.
type Slicer[T any] interface {
	// The Slice method returns a slice containing all the elements in the data structure.
	Slice() []T
}

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
	Value() T

	// The Restart method resets the iterator to the beginning of the collection.
	Restart()
}

// Iterable is an interface that defines a method to get an iterator over a collection
// of elements of type T. It is implemented by data structures that can be iterated over.
type Iterable[T any] interface {
	Iterator() Iterater[T]
}

// GenericIterator is a struct that allows iterating over a slice of elements of any
// type.
type GenericIterator[T any] struct {
	// The slice of elements.
	values *[]T

	// The current index in the values slice.
	index optional.Int
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
		index:  optional.NewInt(-1),
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
func IteratorFromSlicer[T any](slicer Slicer[T]) Iterater[T] {
	elements := slicer.Slice()

	return &GenericIterator[T]{
		values: &elements,
		index:  optional.NewInt(-1),
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
		index:  optional.NewInt(-1),
	}
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
	}

	iter.index.If(func(index int) {
		index++

		if index < len(*iter.values) {
			iter.index = optional.NewInt(index)

			hasNext = true
		} else {
			iter.index = optional.Int{}
		}
	})

	return
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
func (iter *GenericIterator[T]) Value() T {
	defer ers.PropagatePanic(ers.NewErrCallFailed("Value", iter.Value))

	if iter.values == nil {
		panic(errors.New("iterator was never initialized"))
	} else if !iter.index.Present() {
		panic(errors.New("value called on exhausted iter"))
	}

	if index := iter.index.MustGet(); index >= 0 {
		return (*iter.values)[index]
	}

	panic(errors.New("Next must be called before Value"))
}

// Restart is a method of the GenericIterator type that resets the iterator to the
// beginning of the collection.
func (iter *GenericIterator[T]) Restart() {
	iter.index = optional.NewInt(-1)
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
func IteratorFromIterator[E, T any](source Iterater[E], f func(E) Iterater[T]) Iterater[T] {
	if f == nil {
		panic(ers.NewErrInvalidParameter("f").
			WithReason(errors.New("transition function cannot be nil")),
		)
	}

	iter := &ProceduralIterator[E, T]{
		source:  source,
		current: nil,
	}

	iter.transition = f

	return iter
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

		iter.current = iter.transition(iter.source.Value())
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
func (iter *ProceduralIterator[E, T]) Value() T {
	if iter.current == nil {
		panic(ers.NewErrCallFailed("Value", iter.Value).
			WithReason(errors.New("Next must be called before Value")),
		)
	}

	return iter.current.Value()
}

// Restart is a method of the ProceduralIterator type that resets the iterator to the
// beginning of the collection.
func (iter *ProceduralIterator[E, T]) Restart() {
	iter.current = nil
	iter.source.Restart()
}

// Builder is a struct that allows building a collection of elements of type T in a
// procedural manner.
// It is used to build a collection of elements by appending elements to a buffer and
// then creating an iterator over the buffer.
type Builder[T any] struct {
	// The buffer of elements to be built into a collection.
	buffer []T
}

// Append is a method of the Builder type that appends an element to the buffer of
// elements to be built.
//
// Parameters:
//
//   - element: The element to append to the buffer.
func (b *Builder[T]) Append(element T) {
	b.buffer = append(b.buffer, element)
}

// Build is a method of the Builder type that creates an iterator over the buffer of
// elements.
// Finally, the buffer is cleared.
//
// Returns:
//
//   - Iterater[T]: A pointer to a new iterator over the buffer of elements.
func (b *Builder[T]) Build() Iterater[T] {
	iter := &GenericIterator[T]{
		values: &b.buffer,
		index:  optional.NewInt(-1),
	}

	b.buffer = make([]T, 0)

	return iter
}

// Clear is a method of the Builder type that removes all elements from the buffer.
func (b *Builder[T]) Clear() {
	b.buffer = make([]T, 0)
}
