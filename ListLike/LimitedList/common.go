// Package ListLike provides a Lister interface that defines methods for a list data structure.
package ListLike

import (
	"github.com/PlayerR9/MyGoLib/ListLike"
)

// Lister is an interface that defines methods for a list data structure.
// It includes methods to add and remove elements, check if the list is empty or full,
// get the size of the list, convert the list to a slice, clear the list, and get a
// string representation of the list.
type LimitedLister[T any] interface {
	// The Append method adds a value of type T to the end of the list.
	Append(value T) error

	// The DeleteFirst method is a convenience method that deletefirsts an element from
	// the list and returns it.
	// If the list is empty, it will panic.
	DeleteFirst() (T, error)

	// PeekFirst is a method that returns the value at the front of the list without
	// removing it.
	// If the list is empty, it will panic.
	PeekFirst() (T, error)

	// The Prepend method adds a value of type T to the end of the list.
	Prepend(value T) error

	// The DeleteLast method is a convenience method that deletelasts an element from the
	// list and returns it.
	// If the list is empty, it will panic.
	DeleteLast() (T, error)

	// PeekLast is a method that returns the value at the front of the list without
	// removing it.
	// If the list is empty, it will panic.
	PeekLast() (T, error)

	// ListLike[T] is an interface that defines methods for a list data structure.
	ListLike.ListLike[T]

	// The Capacity method returns the maximum number of elements that the list can hold.
	Capacity() (int, bool)

	// The IsFull method checks if the list is full, meaning it has reached its maximum
	// capacity and cannot accept any more elements.
	IsFull() bool
}

// linkedNode represents a node in a linked list. It holds a value of a generic type
// and a reference to the next node in the list.
type linkedNode[T any] struct {
	// The value stored in the node.
	value T

	// A reference to the previous and next nodes in the list, respectively.
	prev, next *linkedNode[T]
}
