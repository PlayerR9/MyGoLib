package ListLike

import (
	"fmt"
	"strings"

	itf "github.com/PlayerR9/MyGoLib/Interfaces"
	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
	"github.com/markphelps/optional"
)

// LinkedList is a generic type that represents a list data structure with
// or without a limited capacity, implemented using a linked list.
type LinkedList[T any] struct {
	// front and back are pointers to the first and last nodes in the linked list,
	// respectively.
	front, back *linkedNode[T]

	// size is the current number of elements in the list.
	size int

	// capacity is the maximum number of elements the list can hold.
	capacity optional.Int
}

// NewLinkedList is a function that creates and returns a new instance of a
// LinkedList.
//
// Parameters:
//
//   - values: A variadic parameter of type T, which represents the initial values to
//     be stored in the list.
//
// Returns:
//
//   - *LinkedList[T]: A pointer to the newly created LinkedList.
func NewLinkedList[T any](values ...T) *LinkedList[T] {
	list := new(LinkedList[T])

	if len(values) == 0 {
		return list
	}

	list.size = len(values)

	// First node
	list_node := &linkedNode[T]{
		value: &values[0],
	}

	list.front = list_node
	list.back = list_node

	// Subsequent nodes
	for _, element := range values {
		list_node := &linkedNode[T]{
			value: &element,
			prev:  list.back,
		}

		list.back.next = list_node
		list.back = list_node
	}

	return list
}

// WithCapacity is a method of the LinkedList type. It is used to set the maximum
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
func (list *LinkedList[T]) WithCapacity(capacity int) (Lister[T], error) {
	if list.capacity.Present() {
		return nil, fmt.Errorf("capacity is already set to %d", list.capacity.MustGet())
	}

	if capacity < 0 {
		return nil, ers.NewErrInvalidParameter("capacity").
			Wrap(fmt.Errorf("negative capacity (%d) is not allowed", capacity))
	} else if list.size > capacity {
		return nil, ers.NewErrInvalidParameter("capacity").
			Wrap(fmt.Errorf("capacity (%d) is not big enough to hold %d elements",
				capacity, list.size))
	}

	list.capacity = optional.NewInt(capacity)

	return list, nil
}

// Append is a method of the LinkedList type. It is used to add an element to
// the end of the list.
//
// Panics with an error of type *ErrFullList if the list is full.
//
// Parameters:
//
//   - value: An element of type T to be added to the list.
func (list *LinkedList[T]) Append(value T) error {
	if list.capacity.Present() && list.size >= list.capacity.MustGet() {
		return NewErrFullList(list)
	}

	list_node := &linkedNode[T]{
		value: &value,
	}

	if list.back == nil {
		list.front = list_node
	} else {
		list.back.next = list_node
		list_node.prev = list.back
	}

	list.back = list_node

	list.size++

	return nil
}

// DeleteFirst is a method of the LinkedList type. It is used to remove and return
// the first element in the list.
//
// Panics with an error of type *ErrCallFailed if the list is empty.
//
// Returns:
//
//   - T: The first element in the list.
func (list *LinkedList[T]) DeleteFirst() (T, error) {
	if list.front == nil {
		return *new(T), NewErrEmptyList(list)
	}

	toRemove := list.front
	list.front = list.front.next

	if list.front == nil {
		list.back = nil
	} else {
		list.front.prev = nil
	}

	list.size--

	toRemove.next = nil

	return *toRemove.value, nil
}

// PeekFirst is a method of the LinkedList type. It is used to return the first
// element in the list without removing it.
//
// Panics with an error of type *ErrCallFailed if the list is empty.
//
// Returns:
//
//   - value: The first element in the list.
func (list *LinkedList[T]) PeekFirst() (T, error) {
	if list.front == nil {
		return *new(T), NewErrEmptyList(list)
	}

	return *list.front.value, nil
}

// IsEmpty is a method of the LinkedList type. It is used to check if the list is
// empty.
//
// Returns:
//
//   - bool: A boolean value that is true if the list is empty, and false otherwise.
func (list *LinkedList[T]) IsEmpty() bool {
	return list.front == nil
}

// Size is a method of the LinkedList type. It is used to return the current number
// of elements in the list.
//
// Returns:
//
//   - int: An integer that represents the current number of elements in the list.
func (list *LinkedList[T]) Size() int {
	return list.size
}

// Capacity is a method of the LinkedList type. It is used to return the maximum
// number of elements the list can hold.
//
// Returns:
//
//   - optional.Int: An optional integer that represents the maximum number of elements
//     the list can hold.
func (list *LinkedList[T]) Capacity() optional.Int {
	return list.capacity
}

// Iterator is a method of the LinkedList type. It is used to return an iterator
// for the list.
//
// Returns:
//
//   - itf.Iterater[T]: An iterator for the list.
func (list *LinkedList[T]) Iterator() itf.Iterater[T] {
	var builder itf.Builder[T]

	for list_node := list.front; list_node != nil; list_node = list_node.next {
		builder.Append(*list_node.value)
	}

	return builder.Build()
}

// Clear is a method of the LinkedList type. It is used to remove all elements from
// the list.
func (list *LinkedList[T]) Clear() {
	if list.front == nil {
		return // List is already empty
	}

	// 1. First node
	list.front.value = nil
	list.front.prev = nil
	prev := list.front

	// 2. Subsequent nodes
	for node := list.front.next; node != nil; node = node.next {
		node.value = nil
		node.prev = nil

		prev = node
		prev.next = nil
	}

	prev.next = nil

	// 3. Reset list fields
	list.front = nil
	list.back = nil
	list.size = 0
}

// IsFull is a method of the LinkedList type. It is used to check if the list is full.
//
// Returns:
//
//   - isFull: A boolean value that is true if the list is full, and false otherwise.
func (list *LinkedList[T]) IsFull() (isFull bool) {
	list.capacity.If(func(cap int) {
		isFull = list.size >= cap
	})

	return
}

// String is a method of the LinkedList type. It returns a string representation of
// the list with information about its size, capacity, and elements.
//
// Returns:
//
//   - string: A string representation of the list.
func (list *LinkedList[T]) String() string {
	var builder strings.Builder

	builder.WriteString("LinkedList[")

	list.capacity.If(func(cap int) {
		fmt.Fprintf(&builder, "capacity=%d, ", cap)
	})

	if list.front == nil {
		builder.WriteString("size=0, values=[]]")

		return builder.String()
	}

	fmt.Fprintf(&builder, "size=%d, values=[%v", list.size, *list.front.value)

	for list_node := list.front.next; list_node != nil; list_node = list_node.next {
		fmt.Fprintf(&builder, ", %v", *list_node.value)
	}

	fmt.Fprintf(&builder, "]]")

	return builder.String()
}

// Prepend is a method of the LinkedList type. It is used to add an element to
// the end of the list.
//
// Panics with an error of type *ErrInvalidOperation if the list is full.
//
// Parameters:
//
//   - value: An element of type T to be added to the list.
func (list *LinkedList[T]) Prepend(value T) error {
	if list.capacity.Present() && list.size >= list.capacity.MustGet() {
		return NewErrFullList(list)
	}

	list_node := &linkedNode[T]{
		value: &value,
	}

	if list.front == nil {
		list.back = list_node
	} else {
		list_node.next = list.front
		list.front.prev = list_node
	}

	list.front = list_node

	list.size++

	return nil
}

// DeleteLast is a method of the LinkedList type. It is used to remove and return
// the last element in the list.
//
// Panics with an error of type *ErrCallFailed if the list is empty.
//
// Returns:
//
//   - T: The last element in the list.
func (list *LinkedList[T]) DeleteLast() (T, error) {
	if list.front == nil {
		return *new(T), NewErrEmptyList(list)
	}

	toRemove := list.back
	list.back = list.back.prev

	if list.back == nil {
		list.front = nil
	} else {
		list.back.next = nil
	}

	list.size--

	toRemove.prev = nil

	return *toRemove.value, nil
}

// PeekLast is a method of the LinkedList type. It is used to return the last
// element in the list without removing it.
//
// Panics with an error of type *ErrCallFailed if the list is empty.
//
// Returns:
//
//   - value: The last element in the list.
func (list *LinkedList[T]) PeekLast() (T, error) {
	if list.front == nil {
		return *new(T), NewErrEmptyList(list)
	}

	return *list.back.value, nil
}

// CutNilValues is a method of the LinkedList type. It is used to remove all nil
// values from the list.
func (list *LinkedList[T]) CutNilValues() {
	if list.front == nil {
		return // List is empty
	}

	if gen.IsNil(*list.front.value) && list.front == list.back {
		// Single node
		list.front = nil
		list.back = nil
		list.size = 0

		return
	}

	var toDelete *linkedNode[T] = nil

	// 1. First node
	if gen.IsNil(*list.front.value) {
		toDelete = list.front

		list.front = list.front.next
		list.front.prev = nil

		toDelete.next = nil
		list.size--
	}

	prev := list.front

	// 2. Subsequent nodes (except last)
	for node := list.front.next; node.next != nil; node = node.next {
		if !gen.IsNil(*node.value) {
			prev = node
		} else {
			prev.next = node.next
			node.next.prev = prev
			list.size--

			if toDelete != nil {
				toDelete.next = nil
			}

			toDelete = node
		}
	}

	if toDelete != nil {
		toDelete.next = nil
	}

	// 3. Last node
	if list.back.value == nil {
		list.back = prev
		list.back.next = nil
		list.size--
	}
}

// Slice is a method of the LinkedList type that returns a slice of type T
//
// Returns:
//
//   - []T: A slice of type T.
func (list *LinkedList[T]) Slice() []T {
	slice := make([]T, 0, list.size)

	for list_node := list.front; list_node != nil; list_node = list_node.next {
		slice = append(slice, *list_node.value)
	}

	return slice
}
