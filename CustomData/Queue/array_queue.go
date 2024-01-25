package Queue

import (
	"fmt"
	"strings"
)

// ArrayQueue is a generic type in Go that represents a queue data structure implemented
// using an array.
// It has a single field, values, which is a slice of type T. This slice stores the
// elements in the queue.
type ArrayQueue[T any] struct {
	values []T
}

// NewArrayQueue is a function that creates and returns a new instance of an ArrayQueue.
// It takes a variadic parameter of type T, which represents the initial values to be
// stored in the queue.
// The function creates a new ArrayQueue, initializes its values field with a slice of
// the same length as the input values, and then copies the input values into the new
// slice. The new ArrayQueue is then returned.
func NewArrayQueue[T any](values ...T) *ArrayQueue[T] {
	queue := &ArrayQueue[T]{
		values: make([]T, len(values)),
	}

	copy(queue.values, values)

	return queue
}

func (queue *ArrayQueue[T]) Enqueue(value T) {
	queue.values = append(queue.values, value)
}

func (queue *ArrayQueue[T]) Dequeue() T {
	if len(queue.values) == 0 {
		panic(NewErrEmptyQueue(Dequeue))
	}

	var value T

	value, queue.values = queue.values[0], queue.values[1:]

	return value
}

func (queue *ArrayQueue[T]) Peek() T {
	if len(queue.values) == 0 {
		panic(NewErrEmptyQueue(Peek))
	}

	return queue.values[0]
}

func (queue *ArrayQueue[T]) IsEmpty() bool {
	return len(queue.values) == 0
}

func (queue *ArrayQueue[T]) Size() int {
	return len(queue.values)
}

func (queue *ArrayQueue[T]) ToSlice() []T {
	slice := make([]T, len(queue.values))

	copy(slice, queue.values)

	return slice
}

func (queue *ArrayQueue[T]) Clear() {
	queue.values = make([]T, 0)
}

// IsFull is a method of the ArrayQueue type. It checks if the queue is full.
//
// In this implementation, the method always returns false. This is because an
// ArrayQueue, implemented with a slice, can dynamically grow and shrink in size
// as elements are added or removed. Therefore, it is never considered full,
// and elements can always be added to it.
func (queue *ArrayQueue[T]) IsFull() bool {
	return false
}

func (queue *ArrayQueue[T]) String() string {
	if len(queue.values) == 0 {
		return QueueHead
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "%s%v", QueueHead, queue.values[0])

	for _, element := range queue.values[1:] {
		fmt.Fprintf(&builder, "%s%v", QueueSep, element)
	}

	return builder.String()
}
