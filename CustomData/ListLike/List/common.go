// Package ListLike provides a Lister interface that defines methods for a list data structure.
package ListLike

import (
	"errors"
	"fmt"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
	"github.com/markphelps/optional"
)

// Lister is an interface that defines methods for a list data structure.
// It includes methods to add and remove elements, check if the list is empty or full,
// get the size of the list, convert the list to a slice, clear the list, and get a
// string representation of the list.
type Lister[T any] interface {
	// The Append method adds a value of type T to the end of the list.
	Append(value *T)

	// The DeleteFirst method is a convenience method that deletefirsts an element from
	// the list and returns it.
	// If the list is empty, it will panic.
	DeleteFirst() *T

	// PeekFirst is a method that returns the value at the front of the list without
	// removing it.
	// If the list is empty, it will panic.
	PeekFirst() *T

	// The IsEmpty method checks if the list is empty and returns a boolean value
	// indicating whether it is empty or not.
	IsEmpty() bool

	// The Size method returns the number of elements currently in the list.
	Size() int

	// The Capacity method returns the maximum number of elements that the list can hold.
	Capacity() optional.Int

	// The ToSlice method returns a slice containing all the elements in the list.
	ToSlice() []*T

	// The Clear method is used to remove all elements from the list, making it empty.
	Clear()

	// The IsFull method checks if the list is full, meaning it has reached its maximum
	// capacity and cannot accept any more elements.
	IsFull() bool

	// The String method returns a string representation of the list.
	// It is useful for debugging and logging purposes.
	fmt.Stringer

	// The Prepend method adds a value of type T to the end of the list.
	Prepend(value *T)

	// The DeleteLast method is a convenience method that deletelasts an element from the
	// list and returns it.
	// If the list is empty, it will panic.
	DeleteLast() *T

	// PeekLast is a method that returns the value at the front of the list without
	// removing it.
	// If the list is empty, it will panic.
	PeekLast() *T

	// CutNilValues is a method that removes all nil values from the list.
	// It is useful for cleaning up the list and removing any empty or nil elements.
	CutNilValues()
}

// linkedNode represents a node in a linked list. It holds a value of a generic type
// and a reference to the next node in the list.
type linkedNode[T any] struct {
	// The value stored in the node.
	value *T

	// A reference to the previous and next nodes in the list, respectively.
	prev, next *linkedNode[T]
}

// ListIterator is a generic type in Go that represents an iterator for a list.
type ListIterator[T any] struct {
	// The values stored in the list.
	values []*T

	// The current index position of the iterator in the list.
	currentIndex int
}

// NewListIterator is a function that creates and returns a new ListIterator object
// for a given list.
//
// Parameters:
//
//   - list: A list of type Lister[T] that the iterator will iterate over.
//
// Returns:
//
//   - *ListIterator[T]: A pointer to a new ListIterator object.
func NewListIterator[T any](list Lister[T]) *ListIterator[T] {
	return &ListIterator[T]{
		values:       list.ToSlice(),
		currentIndex: 0,
	}
}

// GetNext is a method of the ListIterator type that moves the iterator to
// the next element in the list and return the value of that element.
//
// Panics with *ErrOperationFailed if the iterator is at the end of the list.
//
// Returns:
//
//   - *T: A pointer to the next element in the list.
func (iterator *ListIterator[T]) GetNext() *T {
	if len(iterator.values) > iterator.currentIndex {
		value := iterator.values[iterator.currentIndex]
		iterator.currentIndex++

		return value
	}

	panic(ers.NewErrOperationFailed(
		"get next element", errors.New("iterator is out of bounds")),
	)
}

// HasNext is a method of the ListIterator type that checks if there are more elements
// in the list that the iterator can access.
//
// Returns:
//
//   - bool: true if there are more elements, and false otherwise.
func (iterator *ListIterator[T]) HasNext() bool {
	return iterator.currentIndex < len(iterator.values)
}
