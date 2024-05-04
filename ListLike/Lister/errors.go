package Lister

import "fmt"

type ErrEmptyList[T any] struct {
	List Lister[T]
}

func NewErrEmptyList[T any](list Lister[T]) *ErrEmptyList[T] {
	return &ErrEmptyList[T]{List: list}
}

func (e *ErrEmptyList[T]) Error() string {
	return fmt.Sprintf("list (%T) is empty", e.List)
}

type ErrFullList[T any] struct {
	List Lister[T]
}

func NewErrFullList[T any](list Lister[T]) *ErrFullList[T] {
	return &ErrFullList[T]{List: list}
}

func (e *ErrFullList[T]) Error() string {
	return fmt.Sprintf("list (%T) is full", e.List)
}
