package ListLike

import (
	"fmt"
)

type ErrEmptyList[T any] struct {
	list ListLike[T]
}

func NewErrEmptyList[T any](list ListLike[T]) *ErrEmptyList[T] {
	return &ErrEmptyList[T]{list: list}
}

func (e *ErrEmptyList[T]) Error() string {
	return fmt.Sprintf("ListLike (%T) is empty", e.list)
}

type ErrFullList[T any] struct {
	list ListLike[T]
}

func NewErrFullList[T any](list ListLike[T]) *ErrFullList[T] {
	return &ErrFullList[T]{list: list}
}

func (e *ErrFullList[T]) Error() string {
	return fmt.Sprintf("ListLike (%T) is full", e.list)
}
