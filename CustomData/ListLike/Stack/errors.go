package ListLike

import "fmt"

// ErrEmptyStack is a struct that represents an error when attempting to perform a stack
// operation on an empty stack.
// It has a single field, operation, of type StackOperationType, which indicates the type
// of operation that caused the error.
type ErrEmptyStack struct {
	stack string
}

// NewErrEmptyStack creates a new ErrEmptyStack.
// It takes the following parameter:
//
//   - operation is the type of operation that caused the error.
//
// The function returns the following:
//
//   - A pointer to the new ErrEmptyStack.
func NewErrEmptyStack[T any](stack Stacker[T]) *ErrEmptyStack {
	return &ErrEmptyStack{stack: fmt.Sprintf("%T", stack)}
}

// Error is a method of the ErrEmptyStack type that implements the error interface. It
// returns a string representation of the error.
// The method constructs the error message by concatenating the string "could not ", the
// string representation of the operation that caused the error,
// and the string ": stack is empty". This provides a clear and descriptive error message
// when attempting to perform a stack operation on an empty stack.
func (e *ErrEmptyStack) Error() string {
	return fmt.Sprintf("stack (%v) is empty", e.stack)
}

// ErrFullStack is a struct that represents an error when attempting to push an element
// into a full stack.
// It does not have any fields as the error condition is solely based on the state of the
// stack being full.
type ErrFullStack struct {
	stack string
}

// NewErrFullStack creates a new ErrFullStack.
// It takes the following parameter:
//
//   - operation is the type of operation that caused the error.
//
// The function returns the following:
//
//   - A pointer to the new ErrFullStack.
func NewErrFullStack[T any](stack Stacker[T]) *ErrFullStack {
	return &ErrFullStack{stack: fmt.Sprintf("%T", stack)}
}

// Error is a method of the ErrFullStack type that implements the error interface. It
// returns a string representation of the error.
// The method returns the string "could not push: stack is full", providing a clear and
// descriptive error message when attempting to push an element into a full stack.
func (e *ErrFullStack) Error() string {
	return fmt.Sprintf("stack (%v) is full", e.stack)
}
