// Package ListLike provides a Stacker interface that defines methods for a stack data structure.
package ListLike

import (
	"github.com/PlayerR9/MyGoLib/ListLike"
)

// Stacker is an interface that defines methods for a stack data structure.
type LimitedStacker[T any] interface {
	// The Push method adds a value of type T to the end of the stack.
	// If the stack is full, it will panic.
	Push(value T) error

	// The Pop method is a convenience method that pops an element from the stack
	// and returns it.
	// If the stack is empty, it will panic.
	Pop() (T, error)

	// Peek is a method that returns the value at the front of the stack without removing
	// it.
	// If the stack is empty, it will panic.
	Peek() (T, error)

	// .ListLike[T] is an interface that defines methods for a stack data structure.
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
