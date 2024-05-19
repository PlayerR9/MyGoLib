package Stacker

import (
	"fmt"
)

// ErrFullStack is an error type for a full stack.
type ErrFullStack[T any] struct {
	// Stack is the stack that is full.
	Stack Stacker[T]
}

// Error returns the error message: "stack (%T) is full".
//
// Returns:
//   - string: The error message.
//
// Behaviors:
//   - If the stack is nil, the error message is "stack is full".
func (e *ErrFullStack[T]) Error() string {
	if e.Stack == nil {
		return "stack is full"
	} else {
		return fmt.Sprintf("stack (%T) is full", e.Stack)
	}
}

// NewErrFullStack is a constructor for ErrFullStack.
//
// Parameters:
//   - stack: The stack that is full.
//
// Returns:
//   - *ErrFullStack: The error.
func NewErrFullStack[T any](stack Stacker[T]) *ErrFullStack[T] {
	return &ErrFullStack[T]{Stack: stack}
}
