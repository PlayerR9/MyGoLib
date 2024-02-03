// Package ListLike provides a Lister interface that defines methods for a list data structure.
package ListLike

import (
	"errors"
	"fmt"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

// Lister is an interface that defines methods for a list data structure.
// It includes methods to add and remove elements, check if the list is empty or full,
// get the size of the list, convert the list to a slice, clear the list, and get a string
// representation of the list.
type Lister[T any] interface {
	// The Append method adds a value of type T to the end of the list.
	Append(value *T)

	// The DeleteFirst method is a convenience method that deletefirsts an element from the list
	// and returns it.
	// If the list is empty, it will panic.
	DeleteFirst() *T

	// PeekFirst is a method that returns the value at the front of the list without removing
	// it.
	// If the list is empty, it will panic.
	PeekFirst() *T

	// The IsEmpty method checks if the list is empty and returns a boolean value indicating
	// whether it is empty or not.
	IsEmpty() bool

	// The Size method returns the number of elements currently in the list.
	Size() int

	// The ToSlice method returns a slice containing all the elements in the list.
	ToSlice() []*T

	// The Clear method is used to remove all elements from the list, making it empty.
	Clear()

	// The IsFull method checks if the list is full, meaning it has reached its maximum
	// capacity and cannot accept any more elements.
	IsFull() bool

	fmt.Stringer

	// The Prepend method adds a value of type T to the end of the list.
	Prepend(value *T)

	// The DeleteLast method is a convenience method that deletelasts an element from the list
	// and returns it.
	// If the list is empty, it will panic.
	DeleteLast() *T

	// PeekLast is a method that returns the value at the front of the list without removing
	// it.
	// If the list is empty, it will panic.
	PeekLast() *T
}

// linkedNode represents a node in a linked list. It holds a value of a generic type and a
// reference to the next node in the list.
//
// The value field is of a generic type T, which can be any type such as int, string, or a
// custom type.
// It represents the value stored in the node.
//
// The next field is a pointer to the next linkedNode in the list. This allows for traversal
// through the linked list by pointing to the subsequent node in the sequence.
type linkedNode[T any] struct {
	value      *T
	prev, next *linkedNode[T]
}

// ListIterator is a generic type in Go that represents an iterator for a list.
//
// The values field is a slice of type T, which represents the elements stored in the list.
//
// The currentIndex field is an integer that keeps track of the current index position of the
// iterator in the list.
// It is used to iterate over the elements in the list.
type ListIterator[T any] struct {
	values       []*T
	currentIndex int
}

// NewListIterator is a function that creates and returns a new ListIterator object for a
// given list.
// It takes a list of type Lister[T] as an argument, where T can be any type.
//
// The function uses the ToSlice method of the list to get a slice of its values, and
// initializes the currentIndex to -1, indicating that the iterator is at the start of the
// list.
//
// The returned ListIterator can be used to iterate over the elements in the list.
func NewListIterator[T any](list Lister[T]) *ListIterator[T] {
	return &ListIterator[T]{
		values:       list.ToSlice(),
		currentIndex: 0,
	}
}

// GetNext is a method of the ListIterator type. It is used to move the iterator to the next
// element in the list and return the value of that element.
// If the iterator is at the end of the list, the method panics by throwing an
// ErrOutOfBoundsIterator error.
//
// This method is typically used in a loop to iterate over all the elements in a list.
func (iterator *ListIterator[T]) GetNext() *T {
	if len(iterator.values) <= iterator.currentIndex {
		panic(ers.NewErrOperationFailed(
			"get next element", errors.New("iterator is out of bounds")),
		)
	}

	value := iterator.values[iterator.currentIndex]
	iterator.currentIndex++

	return value
}

// HasNext is a method of the ListIterator type. It returns true if there are more elements
// that the iterator is pointing to, and false otherwise.
//
// This method is typically used in conjunction with the GetNext method to iterate over and
// access all the elements in a list.
func (iterator *ListIterator[T]) HasNext() bool {
	return iterator.currentIndex < len(iterator.values)
}
