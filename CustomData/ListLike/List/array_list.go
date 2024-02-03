package ListLike

import (
	"fmt"
	"strings"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

// ArrayList is a generic type in Go that represents a list data structure implemented
// using an array.
// It has a single field, values, which is a slice of type T. This slice stores the
// elements in the list.
type ArrayList[T any] struct {
	values []*T
}

// NewArrayList is a function that creates and returns a new instance of an ArrayList.
// It takes a variadic parameter of type T, which represents the initial values to be
// stored in the list.
// The function creates a new ArrayList, initializes its values field with a slice of
// the same length as the input values, and then copies the input values into the new
// slice. The new ArrayList is then returned.
func NewArrayList[T any](values ...*T) *ArrayList[T] {
	list := &ArrayList[T]{
		values: make([]*T, len(values)),
	}

	copy(list.values, values)

	return list
}

func (list *ArrayList[T]) Append(value *T) {
	list.values = append(list.values, value)
}

func (list *ArrayList[T]) DeleteFirst() *T {
	if len(list.values) == 0 {
		panic(ers.NewErrOperationFailed(
			"delete first element", NewErrEmptyList(list),
		))
	}

	var value *T

	value, list.values = list.values[0], list.values[1:]

	return value
}

func (list *ArrayList[T]) PeekFirst() *T {
	if len(list.values) == 0 {
		panic(ers.NewErrOperationFailed(
			"peek first element", NewErrEmptyList(list),
		))
	}

	return list.values[0]
}

func (list *ArrayList[T]) IsEmpty() bool {
	return len(list.values) == 0
}

func (list *ArrayList[T]) Size() int {
	return len(list.values)
}

func (list *ArrayList[T]) ToSlice() []*T {
	slice := make([]*T, len(list.values))

	copy(slice, list.values)

	return slice
}

func (list *ArrayList[T]) Clear() {
	if len(list.values) == 0 {
		return // List is already empty
	}

	for i := range list.values {
		list.values[i] = nil
	}

	list.values = make([]*T, 0)
}

// IsFull is a method of the ArrayList type. It checks if the list is full.
//
// In this implementation, the method always returns false. This is because an
// ArrayList, implemented with a slice, can dynamically grow and shrink in size
// as elements are added or removed. Therefore, it is never considered full,
// and elements can always be added to it.
func (list *ArrayList[T]) IsFull() bool {
	return false
}

func (list *ArrayList[T]) String() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "ArrayList[size=%d, values=[", len(list.values))

	if len(list.values) > 0 {
		fmt.Fprintf(&builder, "%v", list.values[0])

		for _, value := range list.values[1:] {
			fmt.Fprintf(&builder, ", %v", value)
		}
	}

	builder.WriteString("]]")

	return builder.String()
}

func (list *ArrayList[T]) Prepend(value *T) {
	list.values = append([]*T{value}, list.values...)
}

func (list *ArrayList[T]) DeleteLast() *T {
	if len(list.values) == 0 {
		panic(ers.NewErrOperationFailed(
			"delete last element", NewErrEmptyList(list),
		))
	}

	var value *T

	value, list.values = list.values[len(list.values)-1], list.values[:len(list.values)-1]

	return value
}

func (list *ArrayList[T]) PeekLast() *T {
	if len(list.values) == 0 {
		panic(ers.NewErrOperationFailed(
			"peek last element", NewErrEmptyList(list),
		))
	}

	return list.values[len(list.values)-1]
}
