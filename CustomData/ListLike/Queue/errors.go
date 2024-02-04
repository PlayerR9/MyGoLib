package ListLike

import "fmt"

// ErrEmptyQueue is a struct that represents an error when attempting to perform a queue
// operation on an empty queue.
type ErrEmptyQueue struct {
	// queue is the type of queue that caused the error.
	queue string
}

// NewErrEmptyQueue creates a new ErrEmptyQueue.
//
// Parameters:
//
//   - queue: the type of queue that caused the error.
//
// Returns:
//
//   - *ErrEmptyQueue: a pointer to the new ErrEmptyQueue.
func NewErrEmptyQueue[T any](queue Queuer[T]) *ErrEmptyQueue {
	return &ErrEmptyQueue{queue: fmt.Sprintf("%T", queue)}
}

// Error is a method of the ErrEmptyQueue type that implements the error interface. It
// returns a string representation of the error.
//
// Returns:
//
//   - string: a string representation of the error.
func (e *ErrEmptyQueue) Error() string {
	return fmt.Sprintf("queue (%v) is empty", e.queue)
}

// ErrFullQueue is a struct that represents an error when attempting to enqueue an
// element into a full queue.
type ErrFullQueue struct {
	// queue is the type of queue that caused the error.
	queue string
}

// NewErrFullQueue creates a new ErrFullQueue.
//
// Parameters:
//
//   - queue: the type of queue that caused the error.
//
// Returns:
//
//   - *ErrFullQueue: a pointer to the new ErrFullQueue.
func NewErrFullQueue[T any](queue Queuer[T]) *ErrFullQueue {
	return &ErrFullQueue{queue: fmt.Sprintf("%T", queue)}
}

// Error is a method of the ErrFullQueue type that implements the error interface. It
// returns a string representation of the error.
//
// Returns:
//
//   - string: a string representation of the error.
func (e *ErrFullQueue) Error() string {
	return fmt.Sprintf("queue (%v) is full", e.queue)
}
