package ListLike

import (
	"fmt"

	itf "github.com/PlayerR9/MyGoLib/CustomData/Iterators"
	itff "github.com/PlayerR9/MyGoLibUnits/Interfaces"
)

type ListLike[T any] interface {
	// The IsEmpty method checks if the list is empty and returns a boolean value
	// indicating whether it is empty or not.
	IsEmpty() bool

	// The Size method returns the number of elements currently in the list.
	Size() int

	// The Clear method is used to remove all elements from the list, making it empty.
	Clear()

	// The String method returns a string representation of the list.
	// It is useful for debugging and logging purposes.
	fmt.Stringer

	// CutNilValues is a method that removes all nil values from the list.
	// It is useful for cleaning up the list and removing any empty or nil elements.
	CutNilValues()

	// The itf.Iterable interface is used to provide an iterator for the list.
	itf.Iterable[T]

	// The itf.Slicer interface is used to provide a slicer for the list.
	itff.Slicer[T]

	// The itf.Copier interface is used to provide a method for copying the list.
	itff.Copier
}

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
