package Stream

import (
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

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

	// Get returns qty of items from the stream starting from the given index.
	//
	// Parameters:
	//   - from: The index of the first item to get.
	//   - qty: The number of items to get.
	//
	// Returns:
	//   - []T: The items from the stream.
	//   - error: An error of type *ers.ErrInvalidParameter if from or qty is negative.
	//
	// Behaviors:
	//   - If there are not enough items in the stream, no error is returned
	// 	but the number of items returned will be less than qty.
	Get(from int, qty int) ([]T, error)

	// IsDone returns true if from + qty is greater than the number of items in the stream.
	//
	// Parameters:
	//   - from: The index of the first item to check.
	//   - qty: The number of items to check.
	//
	// Returns:
	//   - bool: True if the stream has been fully consumed.
	IsDone(from int, qty int) bool

	// GetItems returns the items in the stream.
	//
	// Returns:
	//   - []T: The items in the stream.
	GetItems() []T
}

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
	if err := ers.NewErrGTE(0).ErrorIf(from); err != nil {
		return nil, ers.NewErrInvalidParameter("from", err)
	}

	if err := ers.NewErrGTE(-1).ErrorIf(qty); err != nil {
		return nil, ers.NewErrInvalidParameter("qty", err)
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
