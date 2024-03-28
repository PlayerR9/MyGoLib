// Package ListLike provides a Queuer interface that defines methods for a queue data
// structure.
package ListLike

import (
	"github.com/PlayerR9/MyGoLib/ListLike"
)

// Queuer is an interface that defines methods for a queue data structure.
type Queuer[T any] interface {
	// The Enqueue method adds a value of type T to the end of the queue.
	Enqueue(value T)

	// The Dequeue method is a convenience method that dequeues an element from the
	// queue and returns it.
	//
	// Returns:
	//
	//   - T: The value of type T that was dequeued.
	//   - error: An error if the queue is empty.
	Dequeue() (T, error)

	// The Dequeue method is a convenience method that dequeues an element from the
	// queue and returns it. If the queue is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T that was dequeued.
	MustDequeue() T

	// Peek is a method that returns the value at the front of the queue without
	// removing it.
	//
	// Returns:
	//
	//   - T: The value of type T at the front of the queue.
	//   - error: An error if the queue is empty.
	Peek() (T, error)

	// Peek is a method that returns the value at the front of the queue without
	// removing it. If the queue is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T at the front of the queue.
	MustPeek() T

	ListLike.ListLike[T]
}

// linkedNode represents a node in a linked list.
type linkedNode[T any] struct {
	// value is the value stored in the node.
	value T

	// next is a pointer to the next linkedNode in the list.
	next *linkedNode[T]
}
