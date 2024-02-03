package ListLike

import "fmt"

// ErrEmptyQueue is a struct that represents an error when attempting to perform a queue
// operation on an empty queue.
// It has a single field, operation, of type QueueOperationType, which indicates the type
// of operation that caused the error.
type ErrEmptyQueue struct {
	queue string
}

// NewErrEmptyQueue creates a new ErrEmptyQueue.
// It takes the following parameter:
//
//   - operation is the type of operation that caused the error.
//
// The function returns the following:
//
//   - A pointer to the new ErrEmptyQueue.
func NewErrEmptyQueue[T any](queue Queuer[T]) *ErrEmptyQueue {
	return &ErrEmptyQueue{queue: fmt.Sprintf("%T", queue)}
}

// Error is a method of the ErrEmptyQueue type that implements the error interface. It
// returns a string representation of the error.
// The method constructs the error message by concatenating the string "could not ", the
// string representation of the operation that caused the error,
// and the string ": queue is empty". This provides a clear and descriptive error message
// when attempting to perform a queue operation on an empty queue.
func (e *ErrEmptyQueue) Error() string {
	return fmt.Sprintf("queue (%v) is empty", e.queue)
}

// ErrFullQueue is a struct that represents an error when attempting to enqueue an element
// into a full queue.
// It does not have any fields as the error condition is solely based on the state of the
// queue being full.
type ErrFullQueue struct {
	queue string
}

// NewErrFullQueue creates a new ErrFullQueue.
// It takes the following parameter:
//
//   - operation is the type of operation that caused the error.
//
// The function returns the following:
//
//   - A pointer to the new ErrFullQueue.
func NewErrFullQueue[T any](queue Queuer[T]) *ErrFullQueue {
	return &ErrFullQueue{queue: fmt.Sprintf("%T", queue)}
}

// Error is a method of the ErrFullQueue type that implements the error interface. It
// returns a string representation of the error.
// The method returns the string "could not enqueue: queue is full", providing a clear and
// descriptive error message when attempting to enqueue an element into a full queue.
func (e *ErrFullQueue) Error() string {
	return fmt.Sprintf("queue (%v) is full", e.queue)
}
