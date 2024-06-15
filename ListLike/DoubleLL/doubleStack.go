package DoubleLL

import (
	"errors"
	"fmt"

	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	ui "github.com/PlayerR9/MyGoLib/Units/Iterators"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

// DoubleStack is a stack that can be accepted or refused.
// It is used when we want to pop elements from a stack and then decide whether
// to accept or refuse the elements.
//
// If the elements are accepted, they are removed from the stack.
// If the elements are refused, they are pushed back onto the stack as they were.
type DoubleStack[T any] struct {
	// mainStack represents the main stack.
	mainStack Stacker.Stacker[T]

	// auxStack represents the auxiliary stack that is used to store the elements
	// that are popped from the main stack.
	auxStack *Stacker.ArrayStack[T]
}

// Clear clears the double stack.
func (ds *DoubleStack[T]) Clear() {
	ds.mainStack.Clear()
	ds.auxStack.Clear()
}

// IsEmpty returns whether the double stack is empty.
//
// Returns:
//
//   - bool: Whether the double stack is empty.
func (ds *DoubleStack[T]) IsEmpty() bool {
	ok := ds.mainStack.IsEmpty()
	return ok
}

// Push pushes a value onto the double stack.
//
// (WARNING): If the auxiliary stack is not empty, values wiCommon be pushed in
// the middle of the stack.
//
// Parameters:
//   - value (T): The value to push onto the double stack.
//
// Returns:
//   - bool: True if the value was successfully pushed onto the double stack,
//     false otherwise.
func (ds *DoubleStack[T]) Push(value T) bool {
	ok := ds.mainStack.Push(value)
	return ok
}

// Pop pops a value from the double stack.
//
// It stores the popped value in the auxiliary stack.
//
// Returns:
//   - T: The value that was popped from the double stack.
//   - bool: True if the value was successfully popped from the double stack,
//     false otherwise.
func (ds *DoubleStack[T]) Pop() (T, bool) {
	top, ok := ds.mainStack.Pop()
	if !ok {
		return *new(T), false
	}

	ds.auxStack.Push(top)

	return top, true
}

// Peek returns the value at the top of the double stack without removing it.
//
// Returns:
//   - T: The value at the top of the double stack.
//   - bool: True if the value was successfully peeked, false otherwise.
func (ds *DoubleStack[T]) Peek() (T, bool) {
	top, ok := ds.mainStack.Peek()
	if !ok {
		return *new(T), false
	}

	return top, true
}

// Size returns the number of elements in the double stack.
//
// Returns:
//
//   - int: The number of elements in the double stack.
func (ds *DoubleStack[T]) Size() int {
	return ds.mainStack.Size()
}

// GoString implements the fmt.GoStringer interface.
func (ds *DoubleStack[T]) GoString() string {
	return fmt.Sprintf("DoubleStack[mainStack=%s, auxStack=%s]",
		ds.mainStack.GoString(), ds.auxStack.GoString())
}

// CutNilValues is a method that removes aCommon nil values from the double stack.
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
func (ds *DoubleStack[T]) Iterator() ui.Iterater[T] {
	return ds.mainStack.Iterator()
}

// Copy returns a copy of the double stack.
//
// Returns:
//
//   - *DoubleStack: A pointer to a new double stack that is a copy of the original.
func (ds *DoubleStack[T]) Copy() uc.Copier {
	return &DoubleStack[T]{
		mainStack: ds.mainStack.Copy().(Stacker.Stacker[T]),
		auxStack:  ds.auxStack.Copy().(*Stacker.ArrayStack[T]),
	}
}

// Capacity returns the capacity of the double stack.
//
// Returns:
//   - int: The capacity of the double stack.
func (ds *DoubleStack[T]) Capacity() int {
	return ds.mainStack.Capacity()
}

// IsFull returns whether the double stack is full.
//
// Returns:
//   - bool: Whether the double stack is full.
func (ds *DoubleStack[T]) IsFull() bool {
	return ds.mainStack.IsFull()
}

// NewDoubleStack creates a new double stack that uses a specified stack as
// the main stack.
//
// Parameters:
//   - stack: The stack to use as the main stack.
//   - values: The values to push onto the double stack.
//
// Returns:
//   - *DoubleStack: A pointer to the new double stack.
//
// Behaviors:
//   - The stack parameter is used as the main stack while values are pushed onto
//     the end of the specified stack.
func NewDoubleStack[T any](stack Stacker.Stacker[T], values ...T) (*DoubleStack[T], error) {
	if stack == nil {
		return nil, ue.NewErrNilParameter("stack")
	}

	ds := &DoubleStack[T]{
		mainStack: stack,
		auxStack:  Stacker.NewArrayStack[T](),
	}

	for i, value := range values {
		ok := ds.mainStack.Push(value)
		if !ok {
			return nil, ue.NewErrWhileAt("pushing", i+1, "value", nil)
		}
	}

	return ds, nil
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
//
// Returns:
//   - bool: True if the elements were successfully refused, false otherwise.
func (ds *DoubleStack[T]) Refuse() bool {
	for {
		top, ok := ds.auxStack.Pop()
		if !ok {
			break
		}

		ok = ds.mainStack.Push(top)
		if !ok {
			return false
		}
	}

	return true
}

// RefuseOne refuses one element that has been popped from the stack.
// The element is pushed back onto the stack.
//
// Returns:
//   - error: An error if the element could not be refused, nil otherwise.
//
// Errors:
//   - *ErrNoElementsHaveBeenPopped: The auxiliary stack is empty.
//   - *errors.ErrUnexpectedError: The main stack is full.
func (ds *DoubleStack[T]) RefuseOne() error {
	top, ok := ds.auxStack.Pop()
	if !ok {
		return NewErrNoElementsHaveBeenPopped()
	}

	ok = ds.mainStack.Push(top)
	if !ok {
		return ue.NewErrUnexpectedError(errors.New("main stack is full"))
	}

	return nil
}
