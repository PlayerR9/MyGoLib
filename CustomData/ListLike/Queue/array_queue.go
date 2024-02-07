package ListLike

import (
	"fmt"
	"strings"

	itf "github.com/PlayerR9/MyGoLib/Interfaces"
	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
	"github.com/markphelps/optional"
)

// ArrayQueue is a generic type that represents a queue data structure with
// or without a limited capacity. It is implemented using an array.
type ArrayQueue[T any] struct {
	// values is a slice of type T that stores the elements in the queue.
	values []*T

	// capacity is the maximum number of elements the queue can hold.
	capacity optional.Int
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
func NewArrayQueue[T any](values ...*T) *ArrayQueue[T] {
	queue := &ArrayQueue[T]{
		values: make([]*T, len(values)),
	}
	copy(queue.values, values)

	return queue
}

// WithCapacity is a method of the ArrayList type. It is used to set the maximum
// number of elements the list can hold.
//
// Panics with an error of type *ErrCallFailed if the capacity is already set,
// or with an error of type *ErrInvalidParameter if the provided capacity is negative
// or less than the current number of elements in the list.
//
// Parameters:
//
//   - capacity: An integer that represents the maximum number of elements the list can hold.
//
// Returns:
//
//   - Queuer[T]: A pointer to the list.
func (queue *ArrayQueue[T]) WithCapacity(capacity int) Queuer[T] {
	defer ers.PropagatePanic(ers.NewErrCallFailed("WithCapacity", queue.WithCapacity))

	queue.capacity.If(func(cap int) {
		panic(fmt.Errorf("capacity is already set to %d", cap))
	})

	if capacity < 0 {
		panic(ers.NewErrInvalidParameter("capacity").
			WithReason(fmt.Errorf("negative capacity (%d) is not allowed", capacity)),
		)
	} else if len(queue.values) > capacity {
		panic(ers.NewErrInvalidParameter("capacity").WithReason(
			fmt.Errorf("capacity (%d) is less than the current number of elements in the queue (%d)",
				capacity, len(queue.values))),
		)
	}

	newValues := make([]*T, len(queue.values), capacity)
	copy(newValues, queue.values)

	queue.values = newValues

	return queue
}

// Enqueue is a method of the ArrayQueue type. It is used to add an element to the
// end of the queue.
//
// Panics with an error of type *ErrCallFailed if the queue is full.
//
// Parameters:
//
//   - value: The value of type T to be added to the queue.
func (queue *ArrayQueue[T]) Enqueue(value T) {
	queue.capacity.If(func(cap int) {
		if len(queue.values) >= cap {
			panic(ers.NewErrCallFailed("Enqueue", queue.Enqueue).
				WithReason(NewErrFullQueue(queue)),
			)
		}
	})

	queue.values = append(queue.values, &value)
}

// Dequeue is a method of the ArrayQueue type. It is used to remove and return the
// element at the front of the queue.
//
// Panics with an error of type *ErrCallFailed if the queue is empty.
//
// Returns:
//
//   - T: The element at the front of the queue.
func (queue *ArrayQueue[T]) Dequeue() T {
	if len(queue.values) <= 0 {
		panic(ers.NewErrCallFailed("Dequeue", queue.Dequeue).
			WithReason(NewErrEmptyQueue(queue)),
		)
	}

	toRemove := queue.values[0]
	queue.values[0], queue.values = nil, queue.values[1:]
	return *toRemove
}

// Peek is a method of the ArrayQueue type. It is used to return the element at the
// front of the queue without removing it.
//
// Panics with an error of type *ErrCallFailed if the queue is empty.
//
// Returns:
//
//   - T: The element at the front of the queue.
func (queue *ArrayQueue[T]) Peek() T {
	if len(queue.values) == 0 {
		return *queue.values[0]
	}

	panic(ers.NewErrCallFailed("Peek", queue.Peek).
		WithReason(NewErrEmptyQueue(queue)),
	)
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

// Capacity is a method of the ArrayQueue type. It is used to return the maximum
// number of elements the queue can hold.
//
// Returns:
//
//   - optional.Int: An optional integer that represents the maximum number of elements
//     the queue can hold.
func (queue *ArrayQueue[T]) Capacity() optional.Int {
	return queue.capacity
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
		builder.Append(*value)
	}

	return builder.Build()
}

// Clear is a method of the ArrayQueue type. It is used to remove all the elements
// from the queue, making it empty.
func (queue *ArrayQueue[T]) Clear() {
	for i := range queue.values {
		queue.values[i] = nil
	}

	if queue.capacity.Present() {
		queue.values = make([]*T, 0, queue.capacity.MustGet())
	} else {
		queue.values = make([]*T, 0)
	}
}

// IsFull is a method of the ArrayQueue type. It is used to check if the queue is
// full.
//
// Returns:
//
//   - isFull: A boolean value that is true if the queue is full, and false otherwise.
func (queue *ArrayQueue[T]) IsFull() (isFull bool) {
	queue.capacity.If(func(cap int) {
		isFull = len(queue.values) >= cap
	})

	return
}

// String is a method of the ArrayQueue type. It returns a string representation of
// the queue, including its capacity and the elements it contains.
//
// Returns:
//
//   - string: A string representation of the queue.
func (queue *ArrayQueue[T]) String() string {
	var builder strings.Builder

	builder.WriteString("ArrayQueue[")

	queue.capacity.If(func(cap int) {
		fmt.Fprintf(&builder, "capacity=%d, ", cap)
	})

	if len(queue.values) == 0 {
		builder.WriteString("size=0, values=[← ]]")
		return builder.String()
	}

	fmt.Fprintf(&builder, "size=%d, values=[← %v", len(queue.values), *queue.values[0])

	for _, element := range queue.values[1:] {
		fmt.Fprintf(&builder, ", %v", *element)
	}

	builder.WriteString("]]")

	return builder.String()
}

// CutNilValues is a method of the ArrayQueue type. It is used to remove all nil
// values from the queue.
func (queue *ArrayQueue[T]) CutNilValues() {
	for i := 0; i < len(queue.values); {
		if gen.IsNil(*queue.values[i]) {
			queue.values = append(queue.values[:i], queue.values[i+1:]...)
		} else {
			i++
		}
	}
}
