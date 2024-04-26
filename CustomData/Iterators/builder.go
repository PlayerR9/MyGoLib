package Iterators

// Builder is a struct that allows building iterators over a collection of
// elements.
type Builder[T any] struct {
	// buffer is the slice of elements to be built.
	buffer []T
}

// Append is a method of the Builder type that appends an element to the buffer.
//
// Parameters:
//   - element: The element to append to the buffer.
func (b *Builder[T]) Append(element T) {
	b.buffer = append(b.buffer, element)
}

// Build creates a new iterator over the buffer of elements.
//
// It clears the buffer after creating the iterator.
//
// Returns:
//   - Iterater[T]: The new iterator.
func (b *Builder[T]) Build() Iterater[T] {
	bufferCopy := make([]T, len(b.buffer))
	copy(bufferCopy, b.buffer)

	iter := &GenericIterator[T]{
		values: &bufferCopy,
		index:  0,
	}

	b.buffer = make([]T, 0)

	return iter
}

// Clear is a method of the Builder type that removes all elements from the buffer.
func (b *Builder[T]) Clear() {
	b.buffer = make([]T, 0)
}
