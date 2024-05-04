package Queuer

import (
	"fmt"
)

type ErrEmptyQueue[T any] struct {
	Queue Queuer[T]
}

func NewErrEmptyList[T any](queue Queuer[T]) *ErrEmptyQueue[T] {
	return &ErrEmptyQueue[T]{Queue: queue}
}

func (e *ErrEmptyQueue[T]) Error() string {
	return fmt.Sprintf("queue (%T) is empty", e.Queue)
}

type ErrFullQueue[T any] struct {
	Queue Queuer[T]
}

func NewErrFullList[T any](queue Queuer[T]) *ErrFullQueue[T] {
	return &ErrFullQueue[T]{Queue: queue}
}

func (e *ErrFullQueue[T]) Error() string {
	return fmt.Sprintf("queue (%T) is full", e.Queue)
}
