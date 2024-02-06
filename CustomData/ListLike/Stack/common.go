// Package ListLike provides a Stacker interface that defines methods for a stack data structure.
package ListLike

import (
	"errors"
	"fmt"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
	"github.com/markphelps/optional"
)

// Stacker is an interface that defines methods for a stack data structure.
type Stacker[T any] interface {
	// The Push method adds a value of type T to the end of the stack.
	// If the stack is full, it will panic.
	Push(value *T)

	// The Pop method is a convenience method that pops an element from the stack
	// and returns it.
	// If the stack is empty, it will panic.
	Pop() *T

	// Peek is a method that returns the value at the front of the stack without removing
	// it.
	// If the stack is empty, it will panic.
	Peek() *T

	// The IsEmpty method checks if the stack is empty and returns a boolean value indicating
	// whether it is empty or not.
	IsEmpty() bool

	// The Size method returns the number of elements currently in the stack.
	Size() int

	// The Capacity method returns the maximum number of elements that the stack can hold.
	Capacity() optional.Int

	// The ToSlice method returns a slice containing all the elements in the stack.
	ToSlice() []*T

	// The Clear method is used to remove all elements from the stack, making it empty.
	Clear()

	// The IsFull method checks if the stack is full, meaning it has reached its maximum
	// capacity and cannot accept any more elements.
	IsFull() bool

	// The String method returns a string representation of the stack.
	fmt.Stringer

	// The CutNilValues method is used to remove all nil values from the stack.
	CutNilValues()
}

// linkedNode represents a node in a linked list.
type linkedNode[T any] struct {
	// value is the value stored in the node.
	value *T

	// next is a pointer to the next linkedNode in the list.
	next *linkedNode[T]
}

// StackIterator is a generic type in Go that represents an iterator for a stack.
type StackIterator[T any] struct {
	// The values stored in the stack.
	values []*T

	// The current index position of the iterator in the stack.
	currentIndex int
}

// NewStackIterator is a function that creates and returns a new StackIterator object for a
// given stack.
//
// Parameters:
//
//   - stack: A stack of type Stacker[T] that the iterator will be created for.
//
// Returns:
//
//   - *StackIterator[T]: A pointer to the newly created StackIterator object.
func NewStackIterator[T any](stack Stacker[T]) *StackIterator[T] {
	return &StackIterator[T]{
		values:       stack.ToSlice(),
		currentIndex: 0,
	}
}

// GetNext is a method of the StackIterator type. It is used to move the iterator to the next
// element in the stack and return the value of that element.
//
// Panics with an error of type *ErrCallFailed if the iterator is at the end of the stack.
//
// Returns:
//
//   - *T: A pointer to the next element in the stack.
func (iterator *StackIterator[T]) GetNext() *T {
	if len(iterator.values) <= iterator.currentIndex {
		panic(ers.NewErrCallFailed("GetNext", iterator.GetNext).WithReason(
			errors.New("iterator is out of bounds"),
		))
	}

	value := iterator.values[iterator.currentIndex]
	iterator.currentIndex++

	return value
}

// HasNext is a method of the StackIterator type. It is used to check if there are more elements
// to iterate over in the stack.
//
// Returns:
//
//   - bool: true if there are more elements, and false otherwise.
func (iterator *StackIterator[T]) HasNext() bool {
	return iterator.currentIndex < len(iterator.values)
}
