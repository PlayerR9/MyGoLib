package Stack

import "fmt"

// StackHead and StackSep are constants used in the String() method of the Stacker interface
// to format the string representation of a stack.
const (
	// StackHead is a string constant that represents the start of the stack. It is used to
	// indicate where elements are removed from the stack.
	// The value of StackHead is " | →", which visually indicates the direction of element
	// removal.
	StackHead string = " | →"

	// StackSep is a string constant that is used as a separator between elements in the string
	// representation of the stack.
	// The value of StackSep is " | ", which provides a clear visual separation between individual
	// elements in the stack.
	StackSep string = " | "
)

// Package stack provides a Stacker interface that defines methods for a stack data structure.
//
// Stacker is an interface that defines methods for a stack data structure.
// It includes methods to add and remove elements, check if the stack is empty or full,
// get the size of the stack, convert the stack to a slice, clear the stack, and get a string
// representation of the stack.
type Stacker[T any] interface {
	// The Push method adds a value of type T to the end of the stack.
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

	// The ToSlice method returns a slice containing all the elements in the stack.
	ToSlice() []*T

	// The Clear method is used to remove all elements from the stack, making it empty.
	Clear()

	// The IsFull method checks if the stack is full, meaning it has reached its maximum
	// capacity and cannot accept any more elements.
	IsFull() bool

	fmt.Stringer
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
	value *T
	next  *linkedNode[T]
}

// StackIterator is a generic type in Go that represents an iterator for a stack.
//
// The values field is a slice of type T, which represents the elements stored in the stack.
//
// The currentIndex field is an integer that keeps track of the current index position of the
// iterator in the stack.
// It is used to iterate over the elements in the stack.
type StackIterator[T any] struct {
	values       []*T
	currentIndex int
}

// NewStackIterator is a function that creates and returns a new StackIterator object for a
// given stack.
// It takes a stack of type Stacker[T] as an argument, where T can be any type.
//
// The function uses the ToSlice method of the stack to get a slice of its values, and
// initializes the currentIndex to -1, indicating that the iterator is at the start of the
// stack.
//
// The returned StackIterator can be used to iterate over the elements in the stack.
func NewStackIterator[T any](stack Stacker[T]) *StackIterator[T] {
	return &StackIterator[T]{
		values:       stack.ToSlice(),
		currentIndex: 0,
	}
}

// GetNext is a method of the StackIterator type. It is used to move the iterator to the next
// element in the stack and return the value of that element.
// If the iterator is at the end of the stack, the method panics by throwing an
// ErrOutOfBoundsIterator error.
//
// This method is typically used in a loop to iterate over all the elements in a stack.
func (iterator *StackIterator[T]) GetNext() *T {
	if len(iterator.values) <= iterator.currentIndex {
		panic(new(ErrOutOfBoundsIterator))
	}

	value := iterator.values[iterator.currentIndex]
	iterator.currentIndex++

	return value
}

// HasNext is a method of the StackIterator type. It returns true if there are more elements
// that the iterator is pointing to, and false otherwise.
//
// This method is typically used in conjunction with the GetNext method to iterate over and
// access all the elements in a stack.
func (iterator *StackIterator[T]) HasNext() bool {
	return iterator.currentIndex < len(iterator.values)
}
