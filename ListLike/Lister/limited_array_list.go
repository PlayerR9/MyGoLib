package Lister

import (
	"strconv"
	"strings"

	uc "github.com/PlayerR9/MyGoLib/Units/Common"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
	itf "github.com/PlayerR9/MyGoLib/Units/Iterator"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
)

// LimitedArrayList is a generic type that represents a list data structure with
// or without a limited capacity. It is implemented using an array.
type LimitedArrayList[T any] struct {
	// values is a slice of type T that stores the elements in the list.
	values []T

	// capacity is the maximum number of elements the list can hold.
	capacity int
}

// NewLimitedArrayList is a function that creates and returns a new instance of a
// LimitedArrayList.
//
// Parameters:
//
//   - capacity: An integer that represents the maximum number of elements the list.
//     can hold. If the capacity is negative, the value is converted to a positive
//     value.
//   - values: A variadic parameter of type T, which represents the initial values to
//     be stored in the list.
//
// Returns:
//
//   - *LimitedArrayList[T]: A pointer to the newly created LimitedArrayList.
func NewLimitedArrayList[T any](capacity int, values ...T) *LimitedArrayList[T] {
	if capacity < 0 {
		capacity *= -1
	}

	list := &LimitedArrayList[T]{
		values:   make([]T, len(values), capacity),
		capacity: capacity,
	}
	copy(list.values, values)

	return list
}

// Append is a method of the LimitedArrayList type. It is used to add an element to the
// end of the list.
//
// Parameters:
//
//   - value: An element of type T to be added to the list.
//
// Returns:
//
//   - error: An error if the list is full.
func (list *LimitedArrayList[T]) Append(value T) error {
	if len(list.values) >= list.capacity {
		return NewErrFullList(list)
	}

	list.values = append(list.values, value)

	return nil
}

// DeleteFirst is a method of the LimitedArrayList type. It is used to remove and return
// the first element in the list.
// If the list is empty, it will panic.
//
// Returns:
//
//   - T: The first element in the list.
func (list *LimitedArrayList[T]) DeleteFirst() (T, error) {
	if len(list.values) <= 0 {
		return *new(T), ers.NewErrEmpty(list)
	}

	toRemove := list.values[0]
	list.values = list.values[1:]
	return toRemove, nil
}

// PeekFirst is a method of the LimitedArrayList type. It is used to return the first
// element in the list without removing it.
// If the list is empty, it will panic.
//
// Returns:
//
//   - T: A pointer to the first element in the list.
func (list *LimitedArrayList[T]) PeekFirst() (T, error) {
	if len(list.values) == 0 {
		return *new(T), ers.NewErrEmpty(list)
	}

	return list.values[0], nil
}

// IsEmpty is a method of the LimitedArrayList type. It checks if the list is empty.
//
// Returns:
//
//   - bool: A boolean value that is true if the list is empty, and false otherwise.
func (list *LimitedArrayList[T]) IsEmpty() bool {
	return len(list.values) == 0
}

// Size is a method of the LimitedArrayList type. It returns the number of elements in
// the list.
//
// Returns:
//
//   - int: An integer that represents the number of elements in the list.
func (list *LimitedArrayList[T]) Size() int {
	return len(list.values)
}

// Capacity is a method of the LimitedArrayList type. It returns the maximum number of
// elements the list can hold.
//
// Returns:
//
//   - optional.Int: An optional integer that represents the maximum number of
//     elements the list can hold.
func (list *LimitedArrayList[T]) Capacity() int {
	return list.capacity
}

// Iterator is a method of the LimitedArrayList type. It returns an iterator for the list.
//
// Returns:
//
//   - itf.Iterater[T]: An iterator for the list.
func (list *LimitedArrayList[T]) Iterator() itf.Iterater[T] {
	return itf.NewSimpleIterator(list.values)
}

// Clear is a method of the LimitedArrayList type. It is used to remove all elements from
// the list.
func (list *LimitedArrayList[T]) Clear() {
	list.values = make([]T, 0, list.capacity)
}

// IsFull is a method of the LimitedArrayList type. It checks if the list is full.
//
// Returns:
//
//   - isFull: A boolean value that is true if the list is full, and false otherwise.
func (list *LimitedArrayList[T]) IsFull() bool {
	return list.capacity <= len(list.values)
}

// GoString implements the fmt.GoStringer interface.
func (list *LimitedArrayList[T]) GoString() string {
	values := make([]string, 0, len(list.values))
	for _, v := range list.values {
		values = append(values, uc.StringOf(v))
	}

	var builder strings.Builder

	builder.WriteString("LimitedArrayList[capacity=")
	builder.WriteString(strconv.Itoa(list.capacity))
	builder.WriteString(", size=")
	builder.WriteString(strconv.Itoa(len(list.values)))
	builder.WriteString(", values=[")
	builder.WriteString(strings.Join(values, ", "))
	builder.WriteString("]]")

	return builder.String()
}

// Prepend is a method of the LimitedArrayList type. It is used to add an element to the
// end of the list.
//
// Parameters:
//
//   - value: A pointer to an element of type T to be added to the list.
//   - error: An error if the list is full.
func (list *LimitedArrayList[T]) Prepend(value T) error {
	if len(list.values) >= list.capacity {
		return NewErrFullList(list)
	}

	list.values = append([]T{value}, list.values...)

	return nil
}

// DeleteLast is a method of the LimitedArrayList type. It is used to remove and return
// the last element in the list.
// If the list is empty, it will panic.
//
// Returns:
//
//   - T: The last element in the list.
func (list *LimitedArrayList[T]) DeleteLast() (T, error) {
	if len(list.values) == 0 {
		return *new(T), ers.NewErrEmpty(list)
	}

	toRemove := list.values[len(list.values)-1]
	list.values = list.values[:len(list.values)-1]
	return toRemove, nil
}

// PeekLast is a method of the LimitedArrayList type. It is used to return the last
// element in the list without removing it.
// If the list is empty, it will panic.
//
// Returns:
//
//   - T: The last element in the list.
func (list *LimitedArrayList[T]) PeekLast() (T, error) {
	if len(list.values) == 0 {
		return *new(T), ers.NewErrEmpty(list)
	}

	return list.values[len(list.values)-1], nil
}

// CutNilValues is a method of the LimitedArrayList type. It is used to remove all nil
// values from the list.
func (list *LimitedArrayList[T]) CutNilValues() {
	for i := 0; i < len(list.values); {
		if gen.IsNil(list.values[i]) {
			list.values = append(list.values[:i], list.values[i+1:]...)
		} else {
			i++
		}
	}
}

// Slice is a method of the LimitedArrayList type that returns a slice of type T
// containing the elements of the list.
//
// Returns:
//
//   - []T: A slice of type T containing the elements of the list.
func (list *LimitedArrayList[T]) Slice() []T {
	slice := make([]T, len(list.values))
	copy(slice, list.values)

	return slice
}

// Copy is a method of the LimitedArrayList type. It is used to create a shallow copy
// of the list.
//
// Returns:
//
//   - itf.Copier: A copy of the list.
func (list *LimitedArrayList[T]) Copy() uc.Copier {
	listCopy := &LimitedArrayList[T]{
		values:   make([]T, len(list.values)),
		capacity: list.capacity,
	}
	copy(listCopy.values, list.values)

	return listCopy
}
