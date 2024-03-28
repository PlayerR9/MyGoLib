// Package ListLike provides a Stacker interface that defines methods for a stack data structure.
package ListLike

import (
	"github.com/PlayerR9/MyGoLib/ListLike"
)

// Stacker is an interface that defines methods for a stack data structure.
type Stacker[T any] interface {
	// The Push method adds a value of type T to the end of the stack.
	// If the stack is full, it will panic.
	Push(value T)

	// The Pop method pops an element from the stack and returns it.
	// If the stack is empty, it will return an error.
	Pop() (T, error)

	// The MustPop method pops an element from the stack and returns it.
	// If the stack is empty, it will panic.
	MustPop() T

	// Peek is a method that returns the value at the front of the stack without removing
	// it.
	// If the stack is empty, it will error.
	Peek() (T, error)

	// Peek is a method that returns the value at the front of the stack without removing
	// it.
	// If the stack is empty, it will panic.
	MustPeek() T

	// .ListLike[T] is an interface that defines methods for a stack data structure.
	ListLike.ListLike[T]
}

// linkedNode represents a node in a linked list.
type linkedNode[T any] struct {
	// value is the value stored in the node.
	value T

	// next is a pointer to the next linkedNode in the list.
	next *linkedNode[T]
}
