package Queuer

import (
	"fmt"
)

// ErrFullQueue is an error type for a full queue.
type ErrFullQueue[T any] struct {
	// Queue is the queue that is full.
	Queue Queuer[T]
}

// Error returns the error message: "queue (%T) is full".
//
// Returns:
//   - string: The error message.
//
// Behaviors:
//   - If the queue is nil, the error message is "queue is full".
func (e *ErrFullQueue[T]) Error() string {
	if e.Queue == nil {
		return "queue is full"
	} else {
		return fmt.Sprintf("queue (%T) is full", e.Queue)
	}
}

// NewErrFullQueue is a constructor for ErrFullQueue.
//
// Parameters:
//   - queue: The queue that is full.
//
// Returns:
//   - *ErrFullQueue: The error.
func NewErrFullQueue[T any](queue Queuer[T]) *ErrFullQueue[T] {
	return &ErrFullQueue[T]{Queue: queue}
}
