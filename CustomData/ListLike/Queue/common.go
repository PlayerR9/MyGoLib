// Package ListLike provides a Queuer interface that defines methods for a queue data
// structure.
package ListLike

import (
	"github.com/PlayerR9/MyGoLib/CustomData/ListLike"
)

// Queuer is an interface that defines methods for a queue data structure.
type Queuer[T any] interface {
	// The Enqueue method adds a value of type T to the end of the queue.
	// If the queue is full, it will panic.
	Enqueue(value T)

	// The Dequeue method is a convenience method that dequeues an element from the
	// queue and returns it.
	// If the queue is empty, it will panic.
	Dequeue() T

	// Peek is a method that returns the value at the front of the queue without
	// removing it.
	// If the queue is empty, it will panic.
	Peek() T

	// WithCapacity is a special function that modifies an existing queue data
	// structure to have a specific capacity. Panics if the list already has a capacity
	// set or if the new capacity is less than the current size of the list-like data
	// structure.
	//
	// As a result, it is recommended to use this function only when creating a new
	// list-like data structure.
	WithCapacity(int) Queuer[T]

	// ListLike.ListLike[T] is an interface that defines methods for a queue data structure.
	ListLike.ListLike[T]
}

// linkedNode represents a node in a linked list.
type linkedNode[T any] struct {
	// value is the value stored in the node.
	value *T

	// next is a pointer to the next linkedNode in the list.
	next *linkedNode[T]
}
