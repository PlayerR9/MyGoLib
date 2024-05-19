package Stacker

import "github.com/PlayerR9/MyGoLib/ListLike"

// Stacker is an interface that defines methods for a stack data structure.
type Stacker[T any] interface {
	// Push is a method that adds a value of type T to the end of the stack.
	// If the stack is full, it will panic.
	//
	// Parameters:
	//
	//   - value: The value of type T to add to the stack.
	Push(value T) error

	// Pop is a method that pops an element from the stack and returns it.
	// If the stack is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T that was popped.
	Pop() (T, error)

	// Peek is a method that returns the value at the front of the stack without removing
	// it.
	// If the stack is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T at the front of the stack.
	Peek() (T, error)

	ListLike.ListLike[T]
}

// StackNode represents a node in a linked list.
type StackNode[T any] struct {
	// value is the value stored in the node.
	Value T

	// next is a pointer to the next linkedNode in the list.
	next *StackNode[T]
}

// NewStackNode is a constructor function that creates a new linkedNode with the given value.
//
// Parameters:
//  	- value: The value to store in the linkedNode.
//
// Returns:
//  	- *linkedNode: A pointer to the newly created linkedNode.
func NewStackNode[T any](value T) *StackNode[T] {
	return &StackNode[T]{
		Value: value,
	}
}

// SetNext is a method that sets the next linkedNode in the list.
//
// Parameters:
//  	- next: A pointer to the next linkedNode in the list.
func (node *StackNode[T]) SetNext(next *StackNode[T]) {
	node.next = next
}

// Next is a method that returns the next linkedNode in the list.
//
// Returns:
//  	- *linkedNode: A pointer to the next linkedNode in the list.
func (node *StackNode[T]) Next() *StackNode[T] {
	return node.next
}
