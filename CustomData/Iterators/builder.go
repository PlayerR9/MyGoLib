package Iterators

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
	bufferCopy := make([]T, len(b.buffer))
	copy(bufferCopy, b.buffer)

	iter := &GenericIterator[T]{
		values: &bufferCopy,
		index:  -1,
	}

	b.buffer = make([]T, 0)

	return iter
}

// Clear is a method of the Builder type that removes all elements from the buffer.
func (b *Builder[T]) Clear() {
	b.buffer = make([]T, 0)
}
