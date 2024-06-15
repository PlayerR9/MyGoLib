package Queuer

import (
	"github.com/PlayerR9/MyGoLib/ListLike"
)

// Queuer is an interface that defines methods for a queue data structure.
type Queuer[T any] interface {
	// Enqueue is a method that adds a value of type T to the end of the queue.
	//
	// Parameters:
	//   - value: The value of type T to add to the queue.
	//
	// Returns:
	//   - bool: True if the value was successfully added to the queue, false otherwise.
	Enqueue(value T) bool

	// Dequeue is a method that dequeues an element from the queue and returns it.
	//
	// Returns:
	//   - T: The value of type T that was dequeued.
	//   - bool: True if the value was successfully dequeued, false otherwise.
	Dequeue() (T, bool)

	// Peek is a method that returns the value at the front of the queue without
	// removing it.
	//
	// Returns:
	//   - T: The value of type T at the front of the queue.
	//   - bool: True if the value was successfully peeked, false otherwise.
	Peek() (T, bool)

	ListLike.ListLike[T]
}

// QueueNode represents a node in a linked queue.
type QueueNode[T any] struct {
	// Value is the value stored in the node.
	Value T

	// next is a pointer to the next linkedNode in the queue.
	next *QueueNode[T]
}

// NewQueueNode creates a new QueueNode with the given value.
//
// Parameters:
//   - value: The value to store in the node.
//
// Returns:
//   - *QueueNode: A pointer to the new node.
func NewQueueNode[T any](value T) *QueueNode[T] {
	return &QueueNode[T]{
		Value: value,
	}
}

// SetNext sets the next node in the queue.
//
// Parameters:
//   - next: The next node in the queue.
func (node *QueueNode[T]) SetNext(next *QueueNode[T]) {
	node.next = next
}

// Next returns the next node in the queue.
//
// Returns:
//   - *QueueNode: The next node in the queue.
func (node *QueueNode[T]) Next() *QueueNode[T] {
	return node.next
}

// QueueSafeNode represents a node in a linked queue.
type QueueSafeNode[T any] struct {
	// Value is the Value stored in the node.
	Value T

	// next is a pointer to the next queueLinkedNode in the queue.
	next *QueueSafeNode[T]
}

// NewQueueSafeNode creates a new QueueSafeNode with the given value.
//
// Parameters:
//   - value: The value to store in the node.
//
// Returns:
//   - *QueueSafeNode: A pointer to the new node.
func NewQueueSafeNode[T any](value T) *QueueSafeNode[T] {
	return &QueueSafeNode[T]{Value: value}
}

// SetNext sets the next node in the queue.
//
// Parameters:
//   - next: The next node in the queue.
func (node *QueueSafeNode[T]) SetNext(next *QueueSafeNode[T]) {
	node.next = next
}

// Next returns the next node in the queue.
//
// Returns:
//   - *QueueSafeNode: The next node in the queue.
func (node *QueueSafeNode[T]) Next() *QueueSafeNode[T] {
	return node.next
}
