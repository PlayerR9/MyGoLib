package Lister

import (
	"strconv"
	"strings"

	uc "github.com/PlayerR9/MyGoLib/Units/Common"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
	itf "github.com/PlayerR9/MyGoLib/Units/Iterator"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
)

// ArrayList is a generic type that represents a list data structure with
// or without a limited capacity. It is implemented using an array.
type ArrayList[T any] struct {
	// values is a slice of type T that stores the elements in the list.
	values []T
}

// Equals implements Common.Objecter.
func (list *ArrayList[T]) Equals(other uc.Objecter) bool {
	panic("unimplemented")
}

// NewArrayList is a function that creates and returns a new instance of a
// ArrayList.
//
// Parameters:
//
//   - values: A variadic parameter of type T, which represents the initial values to
//     be stored in the list.
//
// Returns:
//
//   - *ArrayList[T]: A pointer to the newly created ArrayList.
func NewArrayList[T any](values ...T) *ArrayList[T] {
	list := &ArrayList[T]{
		values: make([]T, len(values)),
	}
	copy(list.values, values)

	return list
}

// Append is a method of the ArrayList type. It is used to add an element to the
// end of the list.
//
// Panics with an error of type *ErrFullList if the list is full.
//
// Parameters:
//
//   - value: A pointer to an element of type T to be added to the list.
func (list *ArrayList[T]) Append(value T) error {
	list.values = append(list.values, value)

	return nil
}

// DeleteFirst is a method of the ArrayList type. It is used to remove and return
// the first element in the list.
//
// Panics with an error of type *ErrInvalidOperation if the list is empty.
//
// Returns:
//
//   - T: The first element in the list.
func (list *ArrayList[T]) DeleteFirst() (T, error) {
	if len(list.values) <= 0 {
		return *new(T), ers.NewErrEmpty(list)
	}

	toRemove := list.values[0]
	list.values = list.values[1:]
	return toRemove, nil
}

// PeekFirst is a method of the ArrayList type. It is used to return the first
// element in the list without removing it.
//
// Panics with an error of type *ErrInvalidOperation if the list is empty.
//
// Returns:
//
//   - T: A pointer to the first element in the list.
func (list *ArrayList[T]) PeekFirst() (T, error) {
	if len(list.values) == 0 {
		return *new(T), ers.NewErrEmpty(list)
	}

	return list.values[0], nil
}

// IsEmpty is a method of the ArrayList type. It checks if the list is empty.
//
// Returns:
//
//   - bool: A boolean value that is true if the list is empty, and false otherwise.
func (list *ArrayList[T]) IsEmpty() bool {
	return len(list.values) == 0
}

// Size is a method of the ArrayList type. It returns the number of elements in
// the list.
//
// Returns:
//
//   - int: An integer that represents the number of elements in the list.
func (list *ArrayList[T]) Size() int {
	return len(list.values)
}

// Capacity is a method of the ArrayList type. It returns the maximum number of
// elements the list can hold.
//
// Returns:
//
//   - optional.Int: An optional integer that represents the maximum number of
//     elements the list can hold.
func (list *ArrayList[T]) Capacity() int {
	return -1
}

// Iterator is a method of the ArrayList type. It returns an iterator for the list.
//
// Returns:
//
//   - itf.Iterater[T]: An iterator for the list.
func (list *ArrayList[T]) Iterator() itf.Iterater[T] {
	return itf.NewGenericIterator(list.values)
}

// Clear is a method of the ArrayList type. It is used to remove all elements from
// the list.
func (list *ArrayList[T]) Clear() {
	list.values = make([]T, 0)
}

// IsFull is a method of the ArrayList type. It checks if the list is full.
//
// Returns:
//
//   - isFull: A boolean value that is true if the list is full, and false otherwise.
func (list *ArrayList[T]) IsFull() (isFull bool) {
	return false
}

// String is a method of the ArrayList type. It returns a string representation of
// the list with information about its size, capacity, and elements.
//
// Returns:
//
//   - string: A string representation of the list.
func (list *ArrayList[T]) String() string {
	values := make([]string, 0, len(list.values))
	for _, element := range list.values {
		values = append(values, uc.StringOf(element))
	}

	var builder strings.Builder

	builder.WriteString("ArrayList[size=")
	builder.WriteString(strconv.Itoa(len(list.values)))
	builder.WriteString(", values=[")
	builder.WriteString(strings.Join(values, ", "))
	builder.WriteString("]]")

	return builder.String()
}

// Prepend is a method of the ArrayList type. It is used to add an element to the
// end of the list.
//
// Panics with an error of type *ErrFullList if the list is full.
//
// Parameters:
//
//   - value: A pointer to an element of type T to be added to the list.
func (list *ArrayList[T]) Prepend(value T) error {
	list.values = append([]T{value}, list.values...)

	return nil
}

// DeleteLast is a method of the ArrayList type. It is used to remove and return
// the last element in the list.
//
// Panics with an error of type *ErrInvalidOperation if the list is empty.
//
// Returns:
//
//   - T: The last element in the list.
func (list *ArrayList[T]) DeleteLast() (T, error) {
	if len(list.values) == 0 {
		return *new(T), ers.NewErrEmpty(list)
	}

	toRemove := list.values[len(list.values)-1]
	list.values = list.values[:len(list.values)-1]
	return toRemove, nil
}

// PeekLast is a method of the ArrayList type. It is used to return the last
// element in the list without removing it.
//
// Panics with an error of type *ErrInvalidOperation if the list is empty.
//
// Returns:
//
//   - T: The last element in the list.
func (list *ArrayList[T]) PeekLast() (T, error) {
	if len(list.values) == 0 {
		return *new(T), ers.NewErrEmpty(list)
	}

	return list.values[len(list.values)-1], nil
}

// CutNilValues is a method of the ArrayList type. It is used to remove all nil
// values from the list.
func (list *ArrayList[T]) CutNilValues() {
	for i := 0; i < len(list.values); {
		if gen.IsNil(list.values[i]) {
			list.values = append(list.values[:i], list.values[i+1:]...)
		} else {
			i++
		}
	}
}

// Slice is a method of the ArrayList type that returns a slice of type T
// containing the elements of the list.
//
// Returns:
//
//   - []T: A slice of type T containing the elements of the list.
func (list *ArrayList[T]) Slice() []T {
	slice := make([]T, len(list.values))
	copy(slice, list.values)

	return slice
}

// Copy is a method of the ArrayList type. It is used to create a shallow copy
// of the list.
//
// Returns:
//
//   - itf.Copier: A copy of the list.
func (list *ArrayList[T]) Copy() uc.Objecter {
	listCopy := &ArrayList[T]{
		values: make([]T, len(list.values)),
	}
	copy(listCopy.values, list.values)

	return listCopy
}
