package Queue

import (
	"fmt"
	"strings"
)

// LimitedArrayQueue is a generic type in Go that represents a queue data structure with
// a limited capacity.
// It has a single field, values, which is a slice of type T. This slice stores the
// elements in the queue.
type LimitedArrayQueue[T any] struct {
	values []T
}

// NewLimitedArrayQueue is a function that creates and returns a new instance of a
// LimitedArrayQueue.
// It takes an integer capacity, which represents the maximum number of elements the
// queue can hold, and a variadic parameter of type T, which represents the initial values
// to be stored in the queue.
//
// The function first checks if the provided capacity is negative. If it is, it returns an
// error of type ErrNegativeCapacity.
// It then checks if the number of initial values exceeds the provided capacity. If it does,
// it returns an error of type ErrTooManyValues.
//
// If the provided capacity and initial values are valid, the function creates a new
// LimitedArrayQueue, initializes its values field with a slice
// of the same length as the input values and the provided capacity, and then copies the
// input values into the new slice. The new LimitedArrayQueue is then returned.
func NewLimitedArrayQueue[T any](capacity int, values ...T) (*LimitedArrayQueue[T], error) {
	if capacity < 0 {
		return nil, &ErrNegativeCapacity{}
	} else if len(values) > capacity {
		return nil, &ErrTooManyValues{}
	}

	queue := &LimitedArrayQueue[T]{
		values: make([]T, len(values), capacity),
	}
	copy(queue.values, values)

	return queue, nil
}

// Enqueue is a method of the LimitedArrayQueue type. It is used to add an element to the
// end of the queue.
//
// The method takes a parameter, value, of a generic type T, which is the element to be
// added to the queue.
//
// Before adding the element, the method checks if the current length of the values slice
// is equal to the capacity of the queue.
// If it is, it means the queue is full, and the method panics by throwing an ErrFullQueue
// error.
//
// If the queue is not full, the method appends the value to the end of the values slice,
// effectively adding the element to the end of the queue.
func (queue *LimitedArrayQueue[T]) Enqueue(value T) {
	if cap(queue.values) == len(queue.values) {
		panic(&ErrFullQueue{})
	}

	queue.values = append(queue.values, value)
}

func (queue *LimitedArrayQueue[T]) Dequeue() (T, error) {
	if len(queue.values) == 0 {
		return *new(T), &ErrEmptyQueue{Dequeue}
	}

	var value T

	value, queue.values = queue.values[0], queue.values[1:]

	return value, nil
}

func (queue *LimitedArrayQueue[T]) MustDequeue() T {
	if len(queue.values) == 0 {
		panic(&ErrEmptyQueue{Dequeue})
	}

	var value T

	value, queue.values = queue.values[0], queue.values[1:]

	return value
}

func (queue *LimitedArrayQueue[T]) Peek() (T, error) {
	if len(queue.values) == 0 {
		return *new(T), &ErrEmptyQueue{Peek}
	}

	return queue.values[0], nil
}

func (queue *LimitedArrayQueue[T]) MustPeek() T {
	if len(queue.values) == 0 {
		panic(&ErrEmptyQueue{Peek})
	}

	return queue.values[0]
}

func (queue *LimitedArrayQueue[T]) IsEmpty() bool {
	return len(queue.values) == 0
}

func (queue *LimitedArrayQueue[T]) Size() int {
	return len(queue.values)
}

func (queue *LimitedArrayQueue[T]) ToSlice() []T {
	slice := make([]T, len(queue.values))
	copy(slice, queue.values)

	return slice
}

func (queue *LimitedArrayQueue[T]) Clear() {
	queue.values = make([]T, 0, cap(queue.values))
}

func (queue *LimitedArrayQueue[T]) IsFull() bool {
	return cap(queue.values) == len(queue.values)
}

func (queue *LimitedArrayQueue[T]) String() string {
	if len(queue.values) == 0 {
		return QueueHead
	}

	var builder strings.Builder

	builder.WriteString(QueueHead)
	builder.WriteString(fmt.Sprintf("%v", queue.values[0]))

	for _, element := range queue.values[1:] {
		builder.WriteString(QueueSep)
		builder.WriteString(fmt.Sprintf("%v", element))
	}

	return builder.String()
}