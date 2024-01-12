package Queue

import (
	"fmt"
	"strings"
)

// LinkedQueue is a generic type in Go that represents a queue data structure implemented
// using a linked list.
type LinkedQueue[T any] struct {
	// front and back are pointers to the first and last nodes in the linked queue,
	// respectively.
	front, back *linkedNode[T]

	// size is the current number of elements in the queue.
	size int
}

// NewLinkedQueue is a function that creates and returns a new instance of a LinkedQueue.
// It takes a variadic parameter of type T, which represents the initial values to be
// stored in the queue.
//
// If no initial values are provided, the function simply returns a new LinkedQueue with
// all its fields set to their zero values.
//
// If initial values are provided, the function creates a new LinkedQueue and initializes
// its size. It then creates a linked list of nodes
// from the initial values, with each node holding one value, and sets the front and back
// pointers of the queue. The new LinkedQueue is then returned.
func NewLinkedQueue[T any](values ...T) *LinkedQueue[T] {
	if len(values) == 0 {
		return new(LinkedQueue[T])
	}

	queue := new(LinkedQueue[T])
	queue.size = len(values)

	// First node
	node := linkedNode[T]{
		value: values[0],
	}

	queue.front = &node
	queue.back = &node

	// Subsequent nodes
	for _, element := range values[1:] {
		node = linkedNode[T]{
			value: element,
		}

		queue.back.next = &node
		queue.back = &node
	}

	return queue
}

func (queue *LinkedQueue[T]) Enqueue(value T) {
	node := linkedNode[T]{
		value: value,
	}

	if queue.back == nil {
		queue.front = &node
	} else {
		queue.back.next = &node
	}

	queue.back = &node
	queue.size++
}

func (queue *LinkedQueue[T]) Dequeue() (T, error) {
	if queue.front == nil {
		return *new(T), &ErrEmptyQueue{Dequeue}
	}

	var value T

	value, queue.front = queue.front.value, queue.front.next
	if queue.front == nil {
		queue.back = nil
	}

	queue.size--

	return value, nil
}

func (queue *LinkedQueue[T]) MustDequeue() T {
	if queue.front == nil {
		panic(&ErrEmptyQueue{Dequeue})
	}

	var value T

	value, queue.front = queue.front.value, queue.front.next
	if queue.front == nil {
		queue.back = nil
	}

	queue.size--

	return value
}

func (queue *LinkedQueue[T]) Peek() (T, error) {
	if queue.front == nil {
		return *new(T), &ErrEmptyQueue{Peek}
	}

	return queue.front.value, nil
}

func (queue *LinkedQueue[T]) MustPeek() T {
	if queue.front == nil {
		panic(&ErrEmptyQueue{Peek})
	}

	return queue.front.value
}

func (queue *LinkedQueue[T]) IsEmpty() bool {
	return queue.front == nil
}

func (queue *LinkedQueue[T]) Size() int {
	return queue.size
}

func (queue *LinkedQueue[T]) ToSlice() []T {
	slice := make([]T, 0, queue.size)

	for node := queue.front; node != nil; node = node.next {
		slice = append(slice, node.value)
	}

	return slice
}

func (queue *LinkedQueue[T]) Clear() {
	queue.front = nil
	queue.back = nil
	queue.size = 0
}

func (queue *LinkedQueue[T]) IsFull() bool {
	return false
}

func (queue *LinkedQueue[T]) String() string {
	if queue.front == nil {
		return QueueHead
	}

	var builder strings.Builder

	builder.WriteString(QueueHead)
	builder.WriteString(fmt.Sprintf("%v", queue.front.value))

	for node := queue.front.next; node != nil; node = node.next {
		builder.WriteString(QueueSep)
		builder.WriteString(fmt.Sprintf("%v", node.value))
	}

	return builder.String()
}
