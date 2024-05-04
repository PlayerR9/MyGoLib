package Queuer

import (
	"github.com/PlayerR9/MyGoLib/ListLike"
)

// Queuer is an interface that defines methods for a queue data structure.
type Queuer[T any] interface {
	// Enqueue is a method that adds a value of type T to the end of the queue.
	// If the queue is full, it will panic.
	//
	// Parameters:
	//
	//   - value: The value of type T to add to the queue.
	Enqueue(value T) error

	// Dequeue is a method that dequeues an element from the queue and returns it.
	// If the queue is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T that was dequeued.
	Dequeue() (T, error)

	// Peek is a method that returns the value at the front of the queue without
	// removing it.
	// If the queue is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T at the front of the queue.
	Peek() (T, error)

	ListLike.ListLike[T]
}

// QueueNode represents a node in a linked list.
type QueueNode[T any] struct {
	// Value is the value stored in the node.
	Value T

	// next is a pointer to the next linkedNode in the list.
	next *QueueNode[T]
}

// NewQueueNode creates a new LinkedNode with the given value.
func NewQueueNode[T any](value T) *QueueNode[T] {
	return &QueueNode[T]{
		Value: value,
	}
}

func (node *QueueNode[T]) SetNext(next *QueueNode[T]) {
	node.next = next
}

func (node *QueueNode[T]) Next() *QueueNode[T] {
	return node.next
}

// Queuer is an interface that defines methods for a queue data structure.
type SafeQueuer[T any] interface {
	// Enqueue is a method that adds a value of type T to the end of the queue.
	// If the queue is full, it will panic.
	//
	// Parameters:
	//
	//   - value: The value of type T to add to the queue.
	Enqueue(value T) error

	// Dequeue is a method that dequeues an element from the queue and returns it.
	// If the queue is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T that was dequeued.
	Dequeue() (T, error)

	// Peek is a method that returns the value at the front of the queue without
	// removing it.
	// If the queue is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T at the front of the queue.
	Peek() (T, error)

	ListLike.ListLike[T]
}

// QueueSafeNode represents a node in a linked list.
type QueueSafeNode[T any] struct {
	// Value is the Value stored in the node.
	Value T

	// next is a pointer to the next queueLinkedNode in the list.
	next *QueueSafeNode[T]
}

// NewQueueSafeNode creates a new QueueSafeNode with the given value.
func NewQueueSafeNode[T any](value T) *QueueSafeNode[T] {
	return &QueueSafeNode[T]{Value: value}
}

// SetNext sets the next node in the list.
func (node *QueueSafeNode[T]) SetNext(next *QueueSafeNode[T]) {
	node.next = next
}

// Next returns the next node in the list.
func (node *QueueSafeNode[T]) Next() *QueueSafeNode[T] {
	return node.next
}
