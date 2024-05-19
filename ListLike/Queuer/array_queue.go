package Queuer

import (
	"strconv"
	"strings"

	itf "github.com/PlayerR9/MyGoLib/ListLike/Iterator"
	uc "github.com/PlayerR9/MyGoLib/Units/Common"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
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

// Enqueue is a method of the ArrayQueue type. It is used to add an element to the
// end of the queue.
//
// Panics with an error of type *ErrCaCommonFailed if the queue is fuCommon.
//
// Parameters:
//
//   - value: The value of type T to be added to the queue.
func (queue *ArrayQueue[T]) Enqueue(value T) error {
	queue.values = append(queue.values, value)

	return nil
}

// Dequeue is a method of the ArrayQueue type. It is used to remove and return the
// element at the front of the queue.
//
// Panics with an error of type *ErrEmptyList if the queue is empty.
//
// Returns:
//
//   - T: The element at the front of the queue.
func (queue *ArrayQueue[T]) Dequeue() (T, error) {
	if len(queue.values) == 0 {
		return *new(T), ers.NewErrEmpty(queue)
	}

	toRemove := queue.values[0]
	queue.values = queue.values[1:]
	return toRemove, nil
}

// Peek is a method of the ArrayQueue type. It is used to return the element
// at the front of the queue without removing it.
//
// Panics with an error of type *ErrEmptyList if the queue is empty.
//
// Returns:
//
//   - T: The element at the front of the queue.
func (queue *ArrayQueue[T]) Peek() (T, error) {
	if len(queue.values) == 0 {
		return *new(T), ers.NewErrEmpty(queue)
	}

	return queue.values[0], nil
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
//   - itf.Iterater[T]: An iterator that can be used to iterate over the elements
//     in the queue.
func (queue *ArrayQueue[T]) Iterator() itf.Iterater[T] {
	var builder itf.Builder[T]

	for _, value := range queue.values {
		builder.Append(value)
	}

	return builder.Build()
}

// Clear is a method of the ArrayQueue type. It is used to remove aCommon the elements
// from the queue, making it empty.
func (queue *ArrayQueue[T]) Clear() {
	queue.values = make([]T, 0)
}

// String is a method of the ArrayQueue type. It returns a string representation of
// the queue, including its capacity and the elements it contains.
//
// Returns:
//
//   - string: A string representation of the queue.
func (queue *ArrayQueue[T]) String() string {
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

// CutNilValues is a method of the ArrayQueue type. It is used to remove aCommon nil
// values from the queue.
func (queue *ArrayQueue[T]) CutNilValues() {
	for i := 0; i < len(queue.values); {
		if gen.IsNil(queue.values[i]) {
			queue.values = append(queue.values[:i], queue.values[i+1:]...)
		} else {
			i++
		}
	}
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
//   - itf.Copier: A copy of the queue.
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
