package Lister

import (
	"github.com/PlayerR9/MyGoLib/ListLike"
)

// Lister is an interface that defines methods for a list data structure.
type Lister[T any] interface {
	// Append is a method that adds a value of type T to the end of the list.
	//
	// Parameters:
	//
	//   - value: The value of type T to add to the list.
	//
	// Returns:
	//
	//   - error: An error if the list is full.
	Append(value T) error

	// DeleteFirst is a method that deletes an element from the front of the list and
	// returns it. If the list is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T that was deleted.
	DeleteFirst() (T, error)

	// PeekFirst is a method that returns the value at the front of the list without
	// removing it. If the list is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T at the front of the list.
	PeekFirst() (T, error)

	// Prepend is a method that adds a value of type T to the end of the list.
	//
	// Parameters:
	//
	//   - value: The value of type T to add to the list.
	//
	// Returns:
	//
	//   - error: An error if the list is full.
	Prepend(value T) error

	// DeleteLast is a method that deletes an element from the end of the list and
	// returns it. If the list is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T that was deleted.
	DeleteLast() (T, error)

	// PeekLast is a method that returns the value at the end of the list without
	// removing it. If the list is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T at the end of the list.
	PeekLast() (T, error)

	ListLike.ListLike[T]
}

// ListNode represents a node in a linked list. It holds a value of a generic type
// and a reference to the next node in the list.
type ListNode[T any] struct {
	// The Value stored in the node.
	Value T

	// A reference to the previous and next nodes in the list, respectively.
	prev, next *ListNode[T]
}

// NewListNode creates a new LinkedNode with the given value.
func NewListNode[T any](value T) *ListNode[T] {
	return &ListNode[T]{Value: value}
}

func (node *ListNode[T]) SetNext(next *ListNode[T]) {
	node.next = next
}

func (node *ListNode[T]) Next() *ListNode[T] {
	return node.next
}

func (node *ListNode[T]) SetPrev(prev *ListNode[T]) {
	node.prev = prev
}

func (node *ListNode[T]) Prev() *ListNode[T] {
	return node.prev
}

// Lister is an interface that defines methods for a list data structure.
type SafeLister[T any] interface {
	// Append is a method that adds a value of type T to the end of the list.
	//
	// Parameters:
	//
	//   - value: The value of type T to add to the list.
	//
	// Returns:
	//
	//   - error: An error if the list is full.
	Append(value T) error

	// DeleteFirst is a method that deletes an element from the front of the list and
	// returns it. If the list is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T that was deleted.
	DeleteFirst() (T, error)

	// PeekFirst is a method that returns the value at the front of the list without
	// removing it. If the list is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T at the front of the list.
	PeekFirst() (T, error)

	// Prepend is a method that adds a value of type T to the end of the list.
	//
	// Parameters:
	//
	//   - value: The value of type T to add to the list.
	//
	// Returns:
	//
	//   - error: An error if the list is full.
	Prepend(value T) error

	// DeleteLast is a method that deletes an element from the end of the list and
	// returns it. If the list is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T that was deleted.
	DeleteLast() (T, error)

	// PeekLast is a method that returns the value at the end of the list without
	// removing it. If the list is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T at the end of the list.
	PeekLast() (T, error)

	ListLike.ListLike[T]
}

// ListSafeNode represents a node in a linked list. It holds a value of a
// generic type and a reference to the next and previous nodes in the list.
type ListSafeNode[T any] struct {
	// The Value stored in the node.
	Value T

	// A reference to the previous and next nodes in the list, respectively.
	prev, next *ListSafeNode[T]
}

// NewListSafeNode creates a new ListSafeNode with the given value.
func NewListSafeNode[T any](value T) *ListSafeNode[T] {
	return &ListSafeNode[T]{Value: value}
}

// SetNext sets the next node in the list.
func (node *ListSafeNode[T]) SetNext(next *ListSafeNode[T]) {
	node.next = next
}

// Next returns the next node in the list.
func (node *ListSafeNode[T]) Next() *ListSafeNode[T] {
	return node.next
}

// SetPrev sets the previous node in the list.
func (node *ListSafeNode[T]) SetPrev(prev *ListSafeNode[T]) {
	node.prev = prev
}

// Prev returns the previous node in the list.
func (node *ListSafeNode[T]) Prev() *ListSafeNode[T] {
	return node.prev
}