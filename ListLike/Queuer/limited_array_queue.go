package Queuer

import (
	"strconv"
	"strings"

	itf "github.com/PlayerR9/MyGoLib/Units/Iterators"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
)

// LimitedArrayQueue is a generic type that represents a queue data structure with
// or without a limited capacity. It is implemented using an array.
type LimitedArrayQueue[T any] struct {
	// values is a slice of type T that stores the elements in the queue.
	values []T

	// capacity is the maximum number of elements the queue can hold.
	capacity int
}

// NewLimitedArrayQueue is a function that creates and returns a new instance of a
// LimitedArrayQueue.
//
// Parameters:
//
//   - values: A variadic parameter of type T, which represents the initial values to be
//     stored in the queue.
//
// Returns:
//
//   - *LimitedArrayQueue[T]: A pointer to the newly created LimitedArrayQueue.
func NewLimitedArrayQueue[T any](values ...T) *LimitedArrayQueue[T] {
	queue := &LimitedArrayQueue[T]{
		values: make([]T, len(values)),
	}
	copy(queue.values, values)

	return queue
}

// Enqueue implements the Queuer interface.
func (queue *LimitedArrayQueue[T]) Enqueue(value T) bool {
	if len(queue.values) >= queue.capacity {
		return false
	}

	queue.values = append(queue.values, value)

	return true
}

// Dequeue implements the Queuer interface.
func (queue *LimitedArrayQueue[T]) Dequeue() (T, bool) {
	if len(queue.values) == 0 {
		return *new(T), false
	}

	toRemove := queue.values[0]
	queue.values = queue.values[1:]
	return toRemove, true
}

// Peek implements the Queuer interface.
func (queue *LimitedArrayQueue[T]) Peek() (T, bool) {
	if len(queue.values) == 0 {
		return *new(T), false
	}

	return queue.values[0], true
}

// IsEmpty is a method of the LimitedArrayQueue type. It is used to check if the queue is
// empty.
//
// Returns:
//
//   - bool: A boolean value that is true if the queue is empty, and false otherwise.
func (queue *LimitedArrayQueue[T]) IsEmpty() bool {
	return len(queue.values) == 0
}

// Size is a method of the LimitedArrayQueue type. It is used to return the number of
// elements in the queue.
//
// Returns:
//
//   - int: An integer that represents the number of elements in the queue.
func (queue *LimitedArrayQueue[T]) Size() int {
	return len(queue.values)
}

// Capacity is a method of the LimitedArrayQueue type. It is used to return the maximum
// number of elements the queue can hold.
//
// Returns:
//
//   - optional.Int: An optional integer that represents the maximum number of elements
//     the queue can hold.
func (queue *LimitedArrayQueue[T]) Capacity() int {
	return queue.capacity
}

// Iterator is a method of the LimitedArrayQueue type. It is used to return an iterator
// that can be used to iterate over the elements in the queue.
//
// Returns:
//
//   - itf.Iterater[T]: An iterator that can be used to iterate over the elements
//     in the queue.
func (queue *LimitedArrayQueue[T]) Iterator() itf.Iterater[T] {
	return itf.NewSimpleIterator[T](queue.values)
}

// Clear is a method of the LimitedArrayQueue type. It is used to remove all the elements
// from the queue, making it empty.
func (queue *LimitedArrayQueue[T]) Clear() {
	queue.values = make([]T, 0, queue.capacity)
}

// IsFull is a method of the LimitedArrayQueue type. It is used to check if the queue is
// full.
//
// Returns:
//
//   - isFull: A boolean value that is true if the queue is full, and false otherwise.
func (queue *LimitedArrayQueue[T]) IsFull() bool {
	return len(queue.values) >= queue.capacity
}

// GoString implements the fmt.GoStringer interface.
func (queue *LimitedArrayQueue[T]) GoString() string {
	values := make([]string, 0, len(queue.values))
	for _, value := range queue.values {
		values = append(values, uc.StringOf(value))
	}

	var builder strings.Builder

	builder.WriteString("LimitedArrayQueue[capacity=")
	builder.WriteString(strconv.Itoa(queue.capacity))
	builder.WriteString(", size=")
	builder.WriteString(strconv.Itoa(len(queue.values)))
	builder.WriteString(", values=[← ")
	builder.WriteString(strings.Join(values, ", "))
	builder.WriteString("]]")

	return builder.String()
}

// CutNilValues is a method of the LimitedArrayQueue type. It is used to remove all nil
// values from the queue.
func (queue *LimitedArrayQueue[T]) CutNilValues() {
	for i := 0; i < len(queue.values); {
		if gen.IsNil(queue.values[i]) {
			queue.values = append(queue.values[:i], queue.values[i+1:]...)
		} else {
			i++
		}
	}
}

// Slice is a method of the LimitedArrayQueue type. It is used to return a slice of the
// elements in the queue.
//
// Returns:
//
//   - []T: A slice of the elements in the queue.
func (queue *LimitedArrayQueue[T]) Slice() []T {
	slice := make([]T, len(queue.values))
	copy(slice, queue.values)

	return slice
}

// Copy is a method of the LimitedArrayQueue type. It is used to create a shallow copy
// of the queue.
//
// Returns:
//
//   - itf.Copier: A copy of the queue.
func (queue *LimitedArrayQueue[T]) Copy() uc.Copier {
	queueCopy := &LimitedArrayQueue[T]{
		values:   make([]T, len(queue.values)),
		capacity: queue.capacity,
	}
	copy(queueCopy.values, queue.values)

	return queueCopy
}
