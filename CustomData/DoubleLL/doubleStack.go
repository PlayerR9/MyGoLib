package DoubleLL

import (
	"errors"
	"fmt"

	itf "github.com/PlayerR9/MyGoLib/Interfaces"
	Stack "github.com/PlayerR9/MyGoLib/ListLike/Stack"
)

// DoubleStack is a stack that can be accepted or refused.
// It is used when we want to pop elements from a stack and then decide whether
// to accept or refuse the elements.
//
// If the elements are accepted, they are removed from the stack.
// If the elements are refused, they are pushed back onto the stack as they were.
type DoubleStack[T any] struct {
	// mainStack represents the main stack.
	mainStack Stack.Stacker[T]

	// auxStack represents the auxiliary stack that is used to store the elements
	// that are popped from the main stack.
	auxStack *Stack.ArrayStack[T]
}

// NewDoubleLinkedStack creates a new double stack that uses a linked stack as
// the main stack.
//
// Returns:
//
//   - *DoubleStack: A pointer to the new double stack.
func NewDoubleLinkedStack[T any]() *DoubleStack[T] {
	return &DoubleStack[T]{
		mainStack: Stack.NewLinkedStack[T](),
		auxStack:  Stack.NewArrayStack[T](),
	}
}

// NewDoubleArrayStack creates a new double stack that uses an array stack as
// the main stack.
//
// Returns:
//
//   - *DoubleStack: A pointer to the new double stack.
func NewDoubleArrayStack[T any]() *DoubleStack[T] {
	return &DoubleStack[T]{
		mainStack: Stack.NewArrayStack[T](),
		auxStack:  Stack.NewArrayStack[T](),
	}
}

// Clear clears the double stack.
func (ds *DoubleStack[T]) Clear() {
	ds.mainStack.Clear()
	ds.auxStack.Clear()
}

// Accept accepts the elements that have been popped from the stack.
func (ds *DoubleStack[T]) Accept() {
	ds.auxStack.Clear()
}

// GetExtracted returns the elements that have been popped from the stack.
//
// Returns:
//
//   - []T: The elements that have been popped from the stack.
func (ds *DoubleStack[T]) GetExtracted() []T {
	return ds.auxStack.Slice()
}

// Refuse refuses the elements that have been popped from the stack.
// The elements are pushed back onto the stack in the same order that they were
// popped.
func (ds *DoubleStack[T]) Refuse() {
	for !ds.auxStack.IsEmpty() {
		ds.mainStack.Push(ds.auxStack.MustPop())
	}
}

// IsEmpty returns whether the double stack is empty.
//
// Returns:
//
//   - bool: Whether the double stack is empty.
func (ds *DoubleStack[T]) IsEmpty() bool {
	return ds.mainStack.IsEmpty()
}

// Push pushes a value onto the double stack.
//
// (WARNING): If the auxiliary stack is not empty, values will be pushed in
// the middle of the stack.
//
// Parameters:
//
//   - value (T): The value to push onto the double stack.
func (ds *DoubleStack[T]) Push(value T) {
	ds.mainStack.Push(value)
}

// MustPop pops a value from the double stack.
//
// It stores the popped value in the auxiliary stack.
// Panics if the main stack is empty.
//
// Returns:
//
//   - T: The value that was popped from the double stack.
func (ds *DoubleStack[T]) MustPop() T {
	if ds.mainStack.IsEmpty() {
		panic(errors.New("main stack is empty"))
	}

	top := ds.mainStack.MustPop()
	ds.auxStack.Push(top)

	return top
}

// Pop pops a value from the double stack.
//
// It stores the popped value in the auxiliary stack.
// Returns an error if the main stack is empty.
//
// Returns:
//
//   - T: The value that was popped from the double stack.
//   - error: An error if the main stack is empty.
func (ds *DoubleStack[T]) Pop() (T, error) {
	if ds.mainStack.IsEmpty() {
		return *new(T), errors.New("main stack is empty")
	}

	top := ds.mainStack.MustPop()
	ds.auxStack.Push(top)

	return top, nil
}

// Peek returns the value at the top of the double stack without removing it.
//
// Returns:
//
//   - T: The value at the top of the double stack.
//   - error: An error if the main stack is empty.
func (ds *DoubleStack[T]) Peek() (T, error) {
	if ds.mainStack.IsEmpty() {
		return *new(T), errors.New("main stack is empty")
	}

	return ds.mainStack.Peek()
}

// MustPeek returns the value at the top of the double stack without removing it.
//
// Panics if the main stack is empty.
//
// Returns:
//
//   - T: The value at the top of the double stack.
func (ds *DoubleStack[T]) MustPeek() T {
	if ds.mainStack.IsEmpty() {
		panic(errors.New("main stack is empty"))
	}

	return ds.mainStack.MustPeek()
}

// Size returns the number of elements in the double stack.
//
// Returns:
//
//   - int: The number of elements in the double stack.
func (ds *DoubleStack[T]) Size() int {
	return ds.mainStack.Size()
}

// String is a method of fmt.Stringer interface.
//
// It should only be used for debugging and logging purposes.
//
// Returns:
//
//   - string: A string representation of the double stack.
func (ds *DoubleStack[T]) String() string {
	if ds == nil {
		return "DoubleStack[nil]"
	}

	return fmt.Sprintf("DoubleStack[mainStack:%v, auxStack:%v]", ds.mainStack, ds.auxStack)
}

// CutNilValues is a method that removes all nil values from the double stack.
//
// It also removes any empty or nil elements in the auxiliary stack.
func (ds *DoubleStack[T]) CutNilValues() {
	ds.mainStack.CutNilValues()
	ds.auxStack.CutNilValues()
}

// Iterator returns an iterator over the double stack.
//
// Returns:
//
//   - Iterater[T]: A pointer to a new iterator over the double stack.
func (ds *DoubleStack[T]) Iterator() itf.Iterater[T] {
	return ds.mainStack.Iterator()
}

// Slice returns a slice containing all the elements in the double stack.
//
// Returns:
//
//   - []T: A slice containing all the elements in the double stack.
func (ds *DoubleStack[T]) Slice() []T {
	return ds.mainStack.Slice()
}

// Copy returns a copy of the double stack.
//
// Returns:
//
//   - *DoubleStack: A pointer to a new double stack that is a copy of the original.
func (ds *DoubleStack[T]) Copy() itf.Copier {
	return &DoubleStack[T]{
		mainStack: ds.mainStack.Copy().(Stack.Stacker[T]),
		auxStack:  ds.auxStack.Copy().(*Stack.ArrayStack[T]),
	}
}
