package Queuer

import (
	"strconv"
	"strings"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// ArrayQueue is a generic type that represents a queue data structure with
// or without a limited capacity. It is implemented using an array.
type ArrayQueue[T any] struct {
	// values is a slice of type T that stores the elements in the queue.
	values []T
}

// NewArrayQueue is a function that creates and returns a new instance of a
// ArrayQueue.
//
// Parameters:
//
//   - values: A variadic parameter of type T, which represents the initial values to be
//     stored in the queue.
//
// Returns:
//
//   - *ArrayQueue[T]: A pointer to the newly created ArrayQueue.
func NewArrayQueue[T any](values ...T) *ArrayQueue[T] {
	queue := &ArrayQueue[T]{
		values: make([]T, len(values)),
	}
	copy(queue.values, values)

	return queue
}

// Enqueue implements the Queuer interface.
//
// Always returns true.
func (queue *ArrayQueue[T]) Enqueue(value T) bool {
	queue.values = append(queue.values, value)

	return true
}

// Dequeue implements the Queuer interface.
func (queue *ArrayQueue[T]) Dequeue() (T, bool) {
	if len(queue.values) == 0 {
		return *new(T), false
	}

	toRemove := queue.values[0]
	queue.values = queue.values[1:]
	return toRemove, true
}

// Peek implements the Queuer interface.
func (queue *ArrayQueue[T]) Peek() (T, bool) {
	if len(queue.values) == 0 {
		return *new(T), false
	}

	return queue.values[0], true
}

// IsEmpty is a method of the ArrayQueue type. It is used to check if the queue is
// empty.
//
// Returns:
//
//   - bool: A boolean value that is true if the queue is empty, and false otherwise.
func (queue *ArrayQueue[T]) IsEmpty() bool {
	return len(queue.values) == 0
}

// Size is a method of the ArrayQueue type. It is used to return the number of
// elements in the queue.
//
// Returns:
//
//   - int: An integer that represents the number of elements in the queue.
func (queue *ArrayQueue[T]) Size() int {
	return len(queue.values)
}

// Iterator is a method of the ArrayQueue type. It is used to return an iterator
// that can be used to iterate over the elements in the queue.
//
// Returns:
//
//   - uc.Iterater[T]: An iterator that can be used to iterate over the elements
//     in the queue.
func (queue *ArrayQueue[T]) Iterator() uc.Iterater[T] {
	return uc.NewSimpleIterator(queue.values)
}

// Clear is a method of the ArrayQueue type. It is used to remove aCommon the elements
// from the queue, making it empty.
func (queue *ArrayQueue[T]) Clear() {
	queue.values = make([]T, 0)
}

// GoString implements the fmt.GoStringer interface.
func (queue *ArrayQueue[T]) GoString() string {
	values := make([]string, 0, len(queue.values))
	for _, value := range queue.values {
		values = append(values, uc.StringOf(value))
	}

	var builder strings.Builder

	builder.WriteString("ArrayQueue{size=")
	builder.WriteString(strconv.Itoa(len(queue.values)))
	builder.WriteString(", values=[← ")
	builder.WriteString(strings.Join(values, ", "))
	builder.WriteString("]}")

	return builder.String()
}

// Slice is a method of the ArrayQueue type. It is used to return a slice of the
// elements in the queue.
//
// Returns:
//
//   - []T: A slice of the elements in the queue.
func (queue *ArrayQueue[T]) Slice() []T {
	slice := make([]T, len(queue.values))
	copy(slice, queue.values)

	return slice
}

// Copy is a method of the ArrayQueue type. It is used to create a shaCommonow copy
// of the queue.
//
// Returns:
//
//   - uc.Copier: A copy of the queue.
func (queue *ArrayQueue[T]) Copy() uc.Copier {
	queueCopy := &ArrayQueue[T]{
		values: make([]T, len(queue.values)),
	}
	copy(queueCopy.values, queue.values)

	return queueCopy
}

// Capacity is a method of the ArrayQueue type. It is used to return the capacity of
// the queue.
//
// Returns:
//   - int: The capacity of the queue, which is always -1 for an ArrayQueue.
func (queue *ArrayQueue[T]) Capacity() int {
	return -1
}

// IsFull is a method of the ArrayQueue type. It is used to check if the queue is
// full.
//
// Returns:
//   - bool: A boolean value that is always false for an ArrayQueue.
func (queue *ArrayQueue[T]) IsFull() bool {
	return false
}
