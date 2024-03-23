// Package ListLike provides a Queuer interface that defines methods for a queue data
// structure.
package ListLike

import (
	"github.com/PlayerR9/MyGoLib/ListLike"
)

// Queuer is an interface that defines methods for a queue data structure.
type LimitedQueuer[T any] interface {
	// The Enqueue method adds a value of type T to the end of the queue.
	// If the queue is full, it will panic.
	Enqueue(value T) error

	// The Dequeue method is a convenience method that dequeues an element from the
	// queue and returns it.
	// If the queue is empty, it will panic.
	Dequeue() (T, error)

	// Peek is a method that returns the value at the front of the queue without
	// removing it.
	// If the queue is empty, it will panic.
	Peek() (T, error)

	// ListLike.ListLike[T] is an interface that defines methods for a queue data structure.
	ListLike.ListLike[T]

	// The Capacity method returns the maximum number of elements that the list can hold.
	Capacity() (int, bool)

	// The IsFull method checks if the list is full, meaning it has reached its maximum
	// capacity and cannot accept any more elements.
	IsFull() bool
}

// linkedNode represents a node in a linked list.
type linkedNode[T any] struct {
	// value is the value stored in the node.
	value T

	// next is a pointer to the next linkedNode in the list.
	next *linkedNode[T]
}
