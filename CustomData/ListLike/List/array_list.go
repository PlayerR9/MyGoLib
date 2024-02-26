package ListLike

import (
	"fmt"
	"strings"

	itf "github.com/PlayerR9/MyGoLib/Interfaces"
	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
	"github.com/markphelps/optional"
)

// ArrayList is a generic type that represents a list data structure with
// or without a limited capacity. It is implemented using an array.
type ArrayList[T any] struct {
	// values is a slice of type T that stores the elements in the list.
	values []*T

	// capacity is the maximum number of elements the list can hold.
	capacity optional.Int
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
		values: make([]*T, 0, len(values)),
	}

	for _, value := range values {
		list.values = append(list.values, &value)
	}

	return list
}

// WithCapacity is a method of the ArrayList type. It is used to set the maximum
// number of elements the list can hold.
//
// Panics with an error of type *ErrCallFailed if the capacity is already set
// or with an error of type *ErrInvalidParameter if the provided capacity is negative
// or less than the current number of elements in the list.
//
// Parameters:
//
//   - capacity: An integer that represents the maximum number of elements the list
//     can hold.
//
// Returns:
//
//   - Lister[T]: A pointer to the list with the new capacity set.
func (list *ArrayList[T]) WithCapacity(capacity int) (Lister[T], error) {
	if list.capacity.Present() {
		return nil, ers.NewErrInvalidParameter("capacity").
			Wrap(fmt.Errorf("capacity is already set with a value of %d", list.capacity.MustGet()))
	}

	if capacity < 0 {
		return nil, ers.NewErrInvalidParameter("capacity").
			Wrap(fmt.Errorf("negative capacity (%d) is not allowed", capacity))
	} else if len(list.values) > capacity {
		return nil, ers.NewErrInvalidParameter("capacity").
			Wrap(fmt.Errorf("capacity (%d) is less than the current number of elements (%d)",
				capacity, len(list.values)))
	}

	list.capacity = optional.NewInt(capacity)

	newValues := make([]*T, len(list.values), capacity)
	copy(newValues, list.values)

	list.values = newValues

	return list, nil
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
	if list.capacity.Present() && len(list.values) >= list.capacity.MustGet() {
		return NewErrFullList(list)
	}

	list.values = append(list.values, &value)

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
		return *new(T), NewErrEmptyList(list)
	}

	toRemove := list.values[0]
	list.values[0], list.values = nil, list.values[1:]
	return *toRemove, nil
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
		return *new(T), NewErrEmptyList(list)
	}

	return *list.values[0], nil
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
func (list *ArrayList[T]) Capacity() optional.Int {
	return list.capacity
}

// Iterator is a method of the ArrayList type. It returns an iterator for the list.
//
// Returns:
//
//   - itf.Iterater[T]: An iterator for the list.
func (list *ArrayList[T]) Iterator() itf.Iterater[T] {
	var builder itf.Builder[T]

	for _, v := range list.values {
		builder.Append(*v)
	}

	return builder.Build()
}

// Clear is a method of the ArrayList type. It is used to remove all elements from
// the list.
func (list *ArrayList[T]) Clear() {
	if len(list.values) == 0 {
		return // nothing to clear
	}

	for i := range list.values {
		list.values[i] = nil
	}

	if list.capacity.Present() {
		list.values = make([]*T, 0, list.capacity.MustGet())
	} else {
		list.values = make([]*T, 0)
	}
}

// IsFull is a method of the ArrayList type. It checks if the list is full.
//
// Returns:
//
//   - isFull: A boolean value that is true if the list is full, and false otherwise.
func (list *ArrayList[T]) IsFull() (isFull bool) {
	list.capacity.If(func(cap int) {
		isFull = cap <= len(list.values)
	})

	return
}

// String is a method of the ArrayList type. It returns a string representation of
// the list with information about its size, capacity, and elements.
//
// Returns:
//
//   - string: A string representation of the list.
func (list *ArrayList[T]) String() string {
	var builder strings.Builder

	builder.WriteString("ArrayList[")

	list.capacity.If(func(cap int) {
		fmt.Fprintf(&builder, "capacity=%d, ", cap)
	})

	if len(list.values) == 0 {
		builder.WriteString("size=0, values=[]]")

		return builder.String()
	}

	fmt.Fprintf(&builder, "size=%d, values=[%v", len(list.values), *list.values[0])

	for _, element := range list.values[1:] {
		fmt.Fprintf(&builder, ", %v", *element)
	}

	fmt.Fprintf(&builder, "]]")

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
	if list.capacity.Present() && len(list.values) >= list.capacity.MustGet() {
		return NewErrFullList(list)
	}

	list.values = append([]*T{&value}, list.values...)

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
		return *new(T), NewErrEmptyList(list)
	}

	toRemove := list.values[len(list.values)-1]
	list.values[len(list.values)-1], list.values = nil, list.values[:len(list.values)-1]
	return *toRemove, nil
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
		return *new(T), NewErrEmptyList(list)
	}

	return *list.values[len(list.values)-1], nil
}

// CutNilValues is a method of the ArrayList type. It is used to remove all nil
// values from the list.
func (list *ArrayList[T]) CutNilValues() {
	for i := 0; i < len(list.values); {
		if gen.IsNil(*list.values[i]) {
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
	slice := make([]T, 0, len(list.values))

	for _, v := range list.values {
		slice = append(slice, *v)
	}

	return slice
}
