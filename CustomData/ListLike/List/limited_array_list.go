package ListLike

import (
	"fmt"
	"strings"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

// LimitedArrayList is a generic type in Go that represents a list data structure with
// a limited capacity.
// It has a single field, values, which is a slice of type T. This slice stores the
// elements in the list.
type LimitedArrayList[T any] struct {
	values []*T
}

// NewLimitedArrayList is a function that creates and returns a new instance of a
// LimitedArrayList.
// It takes an integer capacity, which represents the maximum number of elements the
// list can hold, and a variadic parameter of type T, which represents the initial values
// to be stored in the list.
//
// The function first checks if the provided capacity is negative. If it is, it returns an
// error of type ErrNegativeCapacity.
// It then checks if the number of initial values exceeds the provided capacity. If it does,
// it returns an error of type ErrTooManyValues.
//
// If the provided capacity and initial values are valid, the function creates a new
// LimitedArrayList, initializes its values field with a slice
// of the same length as the input values and the provided capacity, and then copies the
// input values into the new slice. The new LimitedArrayList is then returned.
func NewLimitedArrayList[T any](capacity int, values ...*T) *LimitedArrayList[T] {
	if capacity < 0 {
		panic(ers.NewErrInvalidParameter(
			"capacity", fmt.Errorf("negative capacity (%d) is not allowed", capacity),
		))
	} else if len(values) > capacity {
		panic(ers.NewErrInvalidParameter(
			"values", fmt.Errorf("number of values (%d) exceeds the provided capacity (%d)",
				len(values),
				capacity,
			),
		))
	}

	list := &LimitedArrayList[T]{
		values: make([]*T, len(values), capacity),
	}
	copy(list.values, values)

	return list
}

// Append is a method of the LimitedArrayList type. It is used to add an element to the
// end of the list.
//
// The method takes a parameter, value, of a generic type T, which is the element to be
// added to the list.
//
// Before adding the element, the method checks if the current length of the values slice
// is equal to the capacity of the list.
// If it is, it means the list is full, and the method panics by throwing an ErrFullList
// error.
//
// If the list is not full, the method appends the value to the end of the values slice,
// effectively adding the element to the end of the list.
func (list *LimitedArrayList[T]) Append(value *T) {
	if cap(list.values) == len(list.values) {
		panic(ers.NewErrOperationFailed(
			"append element", NewErrFullList(list),
		))
	}

	list.values = append(list.values, value)
}

func (list *LimitedArrayList[T]) DeleteFirst() *T {
	if len(list.values) == 0 {
		panic(ers.NewErrOperationFailed(
			"delete first element", NewErrEmptyList(list),
		))
	}

	var value *T

	value, list.values = list.values[0], list.values[1:]

	return value
}

func (list *LimitedArrayList[T]) PeekFirst() *T {
	if len(list.values) == 0 {
		panic(ers.NewErrOperationFailed(
			"peek first element", NewErrEmptyList(list),
		))
	}

	return list.values[0]
}

func (list *LimitedArrayList[T]) IsEmpty() bool {
	return len(list.values) == 0
}

func (list *LimitedArrayList[T]) Size() int {
	return len(list.values)
}

func (list *LimitedArrayList[T]) ToSlice() []*T {
	slice := make([]*T, len(list.values))
	copy(slice, list.values)

	return slice
}

func (list *LimitedArrayList[T]) Clear() {
	if len(list.values) == 0 {
		return // nothing to clear
	}

	for i := range list.values {
		list.values[i] = nil
	}

	list.values = make([]*T, 0, cap(list.values))
}

func (list *LimitedArrayList[T]) IsFull() bool {
	return cap(list.values) == len(list.values)
}

func (list *LimitedArrayList[T]) String() string {
	var builder strings.Builder

	fmt.Fprintf(&builder,
		"LimitedArrayList[size=%d, capacity=%d, values=[", len(list.values), cap(list.values))

	if len(list.values) > 0 {
		fmt.Fprintf(&builder, "%v", list.values[0])

		for _, element := range list.values[1:] {
			fmt.Fprintf(&builder, ", %v", element)
		}
	}

	fmt.Fprintf(&builder, "]]")

	return builder.String()
}

// Prepend is a method of the LimitedArrayList type. It is used to add an element to the
// end of the list.
//
// The method takes a parameter, value, of a generic type T, which is the element to be
// added to the list.
//
// Before adding the element, the method checks if the current length of the values slice
// is equal to the capacity of the list.
// If it is, it means the list is full, and the method panics by throwing an ErrFullList
// error.
//
// If the list is not full, the method appends the value to the end of the values slice,
// effectively adding the element to the end of the list.
func (list *LimitedArrayList[T]) Prepend(value *T) {
	if cap(list.values) == len(list.values) {
		panic(ers.NewErrOperationFailed(
			"prepend element", NewErrFullList(list),
		))
	}

	list.values = append([]*T{value}, list.values...)
}

func (list *LimitedArrayList[T]) DeleteLast() *T {
	if len(list.values) == 0 {
		panic(ers.NewErrOperationFailed(
			"delete last element", NewErrEmptyList(list),
		))
	}

	var value *T

	value, list.values = list.values[len(list.values)-1], list.values[:len(list.values)-1]

	return value
}

func (list *LimitedArrayList[T]) PeekLast() *T {
	if len(list.values) == 0 {
		panic(ers.NewErrOperationFailed(
			"peek last element", NewErrEmptyList(list),
		))
	}

	return list.values[len(list.values)-1]
}
