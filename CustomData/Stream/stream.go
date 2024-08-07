package Stream

import (
	uc "github.com/PlayerR9/lib_units/common"
)

// Stream is a stream of items.
//
// Side effects:
//   - Modifications to the elements inserted and removed from the stream
//     can affect the stream's values. Especially if the elements are pointers.
type Stream[T any] struct {
	// items is the slice of items in the stream.
	items []T

	// size is the number of items in the stream.
	size int
}

// Iterator returns an iterator for the stream.
//
// Returns:
//   - uc.Iterater[T]: An iterator for the stream.
func (s *Stream[T]) Iterator() uc.Iterater[T] {
	return uc.NewSimpleIterator(s.items)
}

// NewStream creates a new stream with the given items.
//
// Parameters:
//   - items: The items to add to the stream.
//
// Returns:
//   - Stream: The new stream.
func NewStream[T any](items []T) *Stream[T] {
	if items == nil {
		return &Stream[T]{
			items: []T{},
			size:  0,
		}
	} else {
		return &Stream[T]{
			items: items,
			size:  len(items),
		}
	}
}

// Size returns the number of items in the stream.
//
// Returns:
//   - int: The number of items in the stream.
func (s *Stream[T]) Size() int {
	return s.size
}

// IsEmpty returns true if the stream is empty.
//
// Returns:
//   - bool: True if the stream is empty.
func (s *Stream[T]) IsEmpty() bool {
	return s.size == 0
}

// Get returns qty of items from the stream starting from the given index.
//
// Parameters:
//   - from: The index of the first item to get.
//   - qty: The number of items to get.
//
// Returns:
//   - []T: The items from the stream.
//   - error: An error if quantity or from is negative.
//
// Behaviors:
//   - Use qty -1 to get all items from 'from' to the end of the stream.
func (s *Stream[T]) Get(from int, qty int) ([]T, error) {
	if from < 0 {
		return nil, uc.NewErrInvalidParameter("from", uc.NewErrGTE(0))
	} else if qty < -1 {
		return nil, uc.NewErrInvalidParameter("qty", uc.NewErrGTE(-1))
	}

	if qty == 0 {
		return nil, nil
	} else if qty == -1 {
		qty = s.size - from
	}

	if from+qty >= s.size {
		return s.items[from:], nil
	} else {
		return s.items[from : from+qty], nil
	}
}

// GetOne returns the item at the given index.
//
// Parameters:
//   - index: The index of the item to get.
//
// Returns:
//   - T: The item at the given index.
//   - error: An error if the index is negative or out of bounds.
func (s *Stream[T]) GetOne(index int) (T, error) {
	if index < 0 {
		return *new(T), uc.NewErrInvalidParameter("index", uc.NewErrGTE(0))
	}

	if index >= s.size {
		return *new(T), uc.NewErrInvalidParameter("index", uc.NewErrLT(s.size))
	}

	return s.items[index], nil
}

// IsDone returns true if from + qty is greater than the number of items in the stream.
//
// Returns:
//   - bool: True if the stream has been fully consumed. False otherwise.
func (s *Stream[T]) IsDone(from int, qty int) bool {
	if from < 0 || qty <= 0 {
		return false
	}

	return from+qty >= len(s.items)
}

// GetItems returns the items in the stream.
//
// Returns:
//   - []T: The items in the stream.
func (s *Stream[T]) GetItems() []T {
	return s.items
}
