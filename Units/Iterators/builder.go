package Iterators

// Builder is a struct that allows building iterators over a collection of
// elements.
type Builder[T any] struct {
	// buffer is the slice of elements to be built.
	buffer []T
}

// Add is a method of the Builder type that appends an element to the buffer.
//
// Parameters:
//   - element: The element to append to the buffer.
func (b *Builder[T]) Add(element T) {
	b.buffer = append(b.buffer, element)
}

// AddMany is a method of the Builder type that appends multiple elements to
// the buffer.
//
// Parameters:
//   - elements: The elements to append to the buffer.
func (b *Builder[T]) AddMany(elements []T) {
	if len(elements) == 0 {
		return
	}

	b.buffer = append(b.buffer, elements...)
}

// Build creates a new iterator over the buffer of elements.
//
// It clears the buffer after creating the iterator.
//
// Returns:
//   - *SimpleIterator[T]: The new iterator.
func (b *Builder[T]) Build() *SimpleIterator[T] {
	bufferCopy := make([]T, len(b.buffer))
	copy(bufferCopy, b.buffer)

	iter := &SimpleIterator[T]{
		values: &bufferCopy,
		index:  0,
	}

	b.Clear()

	return iter
}

// Clear is a method of the Builder type that removes all elements from the buffer.
func (b *Builder[T]) Clear() {
	for i := 0; i < len(b.buffer); i++ {
		b.buffer[i] = *new(T)
	}

	b.buffer = b.buffer[:0]
}
