package Queue

import (
	"fmt"
	"strings"

	itf "github.com/PlayerR9/MyGoLib/CustomData/Iterators"
	"github.com/PlayerR9/MyGoLib/ListLike/Common"
	itff "github.com/PlayerR9/MyGoLib/Units/Interfaces"
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

// Enqueue is a method of the LimitedArrayQueue type. It is used to add an element to the
// end of the queue.
//
// Panics with an error of type *ErrCallFailed if the queue is full.
//
// Parameters:
//
//   - value: The value of type T to be added to the queue.
func (queue *LimitedArrayQueue[T]) Enqueue(value T) {
	if len(queue.values) >= queue.capacity {
		panic(Common.NewErrFullList(queue))
	}

	queue.values = append(queue.values, value)
}

// Dequeue is a method of the LimitedArrayQueue type. It is used to remove and return the
// element at the front of the queue.
//
// Panics with an error of type *ErrCallFailed if the queue is empty.
//
// Returns:
//
//   - T: The element at the front of the queue.
func (queue *LimitedArrayQueue[T]) Dequeue() T {
	if len(queue.values) == 0 {
		panic(Common.NewErrEmptyList(queue))
	}

	toRemove := queue.values[0]
	queue.values = queue.values[1:]
	return toRemove
}

// Peek is a method of the LimitedArrayQueue type. It is used to return the element at the
// front of the queue without removing it.
//
// Panics with an error of type *ErrCallFailed if the queue is empty.
//
// Returns:
//
//   - T: The element at the front of the queue.
func (queue *LimitedArrayQueue[T]) Peek() T {
	if len(queue.values) == 0 {
		panic(Common.NewErrEmptyList(queue))
	}

	return queue.values[0]
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
	var builder itf.Builder[T]

	for _, value := range queue.values {
		builder.Append(value)
	}

	return builder.Build()
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

// String is a method of the LimitedArrayQueue type. It returns a string representation of
// the queue, including its capacity and the elements it contains.
//
// Returns:
//
//   - string: A string representation of the queue.
func (queue *LimitedArrayQueue[T]) String() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "LimitedArrayQueue[capacity=%d, ", queue.capacity)

	if len(queue.values) == 0 {
		builder.WriteString("size=0, values=[← ]]")
		return builder.String()
	}

	fmt.Fprintf(&builder, "size=%d, values=[← %v", len(queue.values), queue.values[0])

	for _, element := range queue.values[1:] {
		fmt.Fprintf(&builder, ", %v", element)
	}

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
func (queue *LimitedArrayQueue[T]) Copy() itff.Copier {
	queueCopy := &LimitedArrayQueue[T]{
		values:   make([]T, len(queue.values)),
		capacity: queue.capacity,
	}
	copy(queueCopy.values, queue.values)

	return queueCopy
}
