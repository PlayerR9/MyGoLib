package Lister

import (
	"fmt"

	uc "github.com/PlayerR9/lib_units/common"
)

// Lister is an interface that defines methods for a list data structure.
type Lister[T any] interface {
	// Append is a method that adds a value of type T to the end of the list.
	//
	// Parameters:
	//   - value: The value of type T to add to the list.
	//
	// Returns:
	//   - bool: True if the value was successfully added to the list, false otherwise.
	Append(value T) bool

	// DeleteFirst is a method that deletes an element from the front of the list and
	// returns it.
	//
	// Returns:
	//   - T: The value of type T that was deleted.
	//   - bool: True if the value was successfully deleted, false otherwise.
	DeleteFirst() (T, bool)

	// PeekFirst is a method that returns the value at the front of the list without
	// removing it.
	//
	// Returns:
	//   - T: The value of type T at the front of the list.
	//   - bool: True if the value was successfully peeked, false otherwise.
	PeekFirst() (T, bool)

	// Prepend is a method that adds a value of type T to the end of the list.
	//
	// Parameters:
	//   - value: The value of type T to add to the list.
	//
	// Returns:
	//   - bool: True if the value was successfully added to the list, false otherwise.
	Prepend(value T) bool

	// DeleteLast is a method that deletes an element from the end of the list and
	// returns it.
	//
	// Returns:
	//   - T: The value of type T that was deleted.
	//   - bool: True if the value was successfully deleted, false otherwise.
	DeleteLast() (T, bool)

	// PeekLast is a method that returns the value at the end of the list without
	// removing it.
	//
	// Returns:
	//   - T: The value of type T at the end of the list.
	//   - bool: True if the value was successfully peeked, false otherwise.
	PeekLast() (T, bool)

	// IsEmpty is a method that checks whether the list is empty.
	//
	// Returns:
	//
	//   - bool: True if the list is empty, false otherwise.
	IsEmpty() bool

	// Size method returns the number of elements currently in the list.
	//
	// Returns:
	//
	//   - int: The number of elements in the list.
	Size() int

	// Clear method is used to remove all elements from the list, making it empty.
	Clear()

	// Capacity is a method that returns the maximum number of elements that the list can hold.
	//
	// Returns:
	//
	//   - int: The maximum number of elements that the list can hold. -1 if there is no limit.
	Capacity() int

	// IsFull is a method that checks whether the list is full.
	//
	// Returns:
	//
	//   - bool: True if the list is full, false otherwise.
	IsFull() bool

	// Slice is a method that returns a slice of the elements in the list.
	//
	// Returns:
	//  	- []T: A slice of the elements in the list.
	Slice() []T

	uc.Iterable[T]

	uc.Copier
	fmt.GoStringer
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
//
// Parameters:
//   - value: The value to store in the node.
//
// Returns:
//   - *ListNode: A pointer to the new node.
func NewListNode[T any](value T) *ListNode[T] {
	return &ListNode[T]{Value: value}
}

// SetNext sets the next node in the list.
//
// Parameters:
//   - next: The next node in the list.
func (node *ListNode[T]) SetNext(next *ListNode[T]) {
	node.next = next
}

// Next returns the next node in the list.
//
// Returns:
//   - *ListNode: The next node in the list.
func (node *ListNode[T]) Next() *ListNode[T] {
	return node.next
}

// SetPrev sets the previous node in the list.
//
// Parameters:
//   - prev: The previous node in the list.
func (node *ListNode[T]) SetPrev(prev *ListNode[T]) {
	node.prev = prev
}

// Prev returns the previous node in the list.
//
// Returns:
//   - *ListNode: The previous node in the list.
func (node *ListNode[T]) Prev() *ListNode[T] {
	return node.prev
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
//
// Parameters:
//   - value: The value to store in the node.
//
// Returns:
//   - *ListSafeNode: A pointer to the new node.
func NewListSafeNode[T any](value T) *ListSafeNode[T] {
	return &ListSafeNode[T]{Value: value}
}

// SetNext sets the next node in the list.
//
// Parameters:
//   - next: The next node in the list.
func (node *ListSafeNode[T]) SetNext(next *ListSafeNode[T]) {
	node.next = next
}

// Next returns the next node in the list.
//
// Returns:
//   - *ListSafeNode: The next node in the list.
func (node *ListSafeNode[T]) Next() *ListSafeNode[T] {
	return node.next
}

// SetPrev sets the previous node in the list.
//
// Parameters:
//   - prev: The previous node in the list.
func (node *ListSafeNode[T]) SetPrev(prev *ListSafeNode[T]) {
	node.prev = prev
}

// Prev returns the previous node in the list.
//
// Returns:
//   - *ListSafeNode: The previous node in the list.
func (node *ListSafeNode[T]) Prev() *ListSafeNode[T] {
	return node.prev
}
