package Stream

// Streamer is an interface for streams of items.
type Streamer[T any] interface {
	// Size returns the number of items in the stream.
	//
	// Returns:
	//   - int: The number of items in the stream.
	Size() int

	// IsEmpty returns true if the stream is empty.
	//
	// Returns:
	//   - bool: True if the stream is empty.
	IsEmpty() bool

	// Peek returns the next item in the stream without consuming it.
	// It panics if there are no more items in the stream.
	//
	// Returns:
	//   - T: The next item in the stream.
	Peek() T

	// Consume consumes the next item in the stream.
	// It panics if there are no more items in the stream.
	//
	// Returns:
	//   - T: The consumed item.
	Consume() T

	// Reset resets the stream to the beginning.
	Reset()

	// IsDone returns true if the stream has been fully consumed.
	//
	// Returns:
	//   - bool: True if the stream has been fully consumed.
	IsDone() bool

	// GetItems returns the items in the stream.
	//
	// Returns:
	//   - []T: The items in the stream.
	GetItems() []T

	// GetLeftoverItems returns the items that have not been consumed.
	//
	// Returns:
	//   - []T: The leftover items in the stream.
	GetLeftoverItems() []T
}

// Stream is a stream of items.
//
// Side effects:
// 	- Modifications to the elements inserted and removed from the stream
// 	can affect the stream's values. Especially if the elements are pointers.
type Stream[T any] struct {
	// items is the slice of items in the stream.
	items []T

	// currentIndex is the current index of the stream.
	// It indicates the first non-consumed item.
	currentIndex int
}

// Size returns the number of items in the stream.
//
// Returns:
//   - int: The number of items in the stream.
func (s *Stream[T]) Size() int {
	return len(s.items)
}

// IsEmpty returns true if the stream is empty.
//
// Returns:
//   - bool: True if the stream is empty.
func (s *Stream[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Peek returns the next item in the stream without consuming it.
// It panics if there are no more items in the stream.
//
// Returns:
//   - T: The next item in the stream.
func (s *Stream[T]) Peek() T {
	if s.currentIndex >= len(s.items) {
		panic(NewErrNoMoreItems())
	}

	return s.items[s.currentIndex]
}

// Consume consumes the next item in the stream.
// It panics if there are no more items in the stream.
//
// Returns:
//   - T: The consumed item.
func (s *Stream[T]) Consume() T {
	if s.currentIndex >= len(s.items) {
		panic(NewErrNoMoreItems())
	}

	item := s.items[s.currentIndex]
	s.currentIndex++

	return item
}

// Reset resets the stream to the beginning.
func (s *Stream[T]) Reset() {
	s.currentIndex = 0
}

// IsDone returns true if the stream has been fully consumed.
//
// Returns:
//   - bool: True if the stream has been fully consumed.
func (s *Stream[T]) IsDone() bool {
	return s.currentIndex >= len(s.items)
}

// GetItems returns the items in the stream.
//
// Returns:
//   - []T: The items in the stream.
func (s *Stream[T]) GetItems() []T {
	return s.items
}

// GetLeftoverItems returns the items that have not been consumed.
//
// Returns:
//   - []T: The leftover items in the stream.
func (s *Stream[T]) GetLeftoverItems() []T {
	return s.items[s.currentIndex:]
}

// NewStream creates a new stream with the given items.
//
// Parameters:
//   - items: The items to add to the stream.
//
// Returns:
//   - Stream: The new stream.
func NewStream[T any](items []T) *Stream[T] {
	return &Stream[T]{
		items:        items,
		currentIndex: 0,
	}
}
