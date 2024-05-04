package Stacker

import (
	"fmt"
)

type ErrEmptyStack[T any] struct {
	Stack Stacker[T]
}

func NewErrEmptyStack[T any](stack Stacker[T]) *ErrEmptyStack[T] {
	return &ErrEmptyStack[T]{Stack: stack}
}

func (e *ErrEmptyStack[T]) Error() string {
	return fmt.Sprintf("stack (%T) is empty", e.Stack)
}

type ErrFullStack[T any] struct {
	Stack Stacker[T]
}

func NewErrFullList[T any](stack Stacker[T]) *ErrFullStack[T] {
	return &ErrFullStack[T]{Stack: stack}
}

func (e *ErrFullStack[T]) Error() string {
	return fmt.Sprintf("stack (%T) is full", e.Stack)
}
