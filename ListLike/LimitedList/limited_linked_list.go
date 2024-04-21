package ListLike

import (
	"fmt"
	"strings"

	itff "github.com/PlayerR9/MyGoLib/Common/Interfaces"
	itf "github.com/PlayerR9/MyGoLib/CustomData/Iterators"
	ll "github.com/PlayerR9/MyGoLib/ListLike"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
)

// LimitedLinkedList is a generic type that represents a list data structure with
// or without a limited capacity, implemented using a linked list.
type LimitedLinkedList[T any] struct {
	// front and back are pointers to the first and last nodes in the linked list,
	// respectively.
	front, back *linkedNode[T]

	// size is the current number of elements in the list.
	size int

	// capacity is the maximum number of elements the list can hold.
	capacity int
}

// NewLimitedLinkedList is a function that creates and returns a new instance of a
// LimitedLinkedList.
//
// Parameters:
//
//   - values: A variadic parameter of type T, which represents the initial values to
//     be stored in the list.
//
// Returns:
//
//   - *LimitedLinkedList[T]: A pointer to the newly created LimitedLinkedList.
func NewLimitedLinkedList[T any](values ...T) *LimitedLinkedList[T] {
	list := new(LimitedLinkedList[T])

	if len(values) == 0 {
		return list
	}

	list.size = len(values)

	// First node
	list_node := &linkedNode[T]{
		value: values[0],
	}

	list.front = list_node
	list.back = list_node

	// Subsequent nodes
	for _, element := range values {
		list_node := &linkedNode[T]{
			value: element,
			prev:  list.back,
		}

		list.back.next = list_node
		list.back = list_node
	}

	return list
}

// Append is a method of the LimitedLinkedList type. It is used to add an element to
// the end of the list.
//
// Panics with an error of type *ErrFullList if the list is full.
//
// Parameters:
//
//   - value: An element of type T to be added to the list.
func (list *LimitedLinkedList[T]) Append(value T) error {
	if list.size >= list.capacity {
		return ll.NewErrFullList(list)
	}

	list_node := &linkedNode[T]{
		value: value,
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

// DeleteFirst is a method of the LimitedLinkedList type. It is used to remove and return
// the first element in the list.
//
// Panics with an error of type *ErrCallFailed if the list is empty.
//
// Returns:
//
//   - T: The first element in the list.
func (list *LimitedLinkedList[T]) DeleteFirst() (T, error) {
	if list.front == nil {
		return *new(T), ll.NewErrEmptyList(list)
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

	return toRemove.value, nil
}

// PeekFirst is a method of the LimitedLinkedList type. It is used to return the first
// element in the list without removing it.
//
// Panics with an error of type *ErrCallFailed if the list is empty.
//
// Returns:
//
//   - value: The first element in the list.
func (list *LimitedLinkedList[T]) PeekFirst() (T, error) {
	if list.front == nil {
		return *new(T), ll.NewErrEmptyList(list)
	}

	return list.front.value, nil
}

// IsEmpty is a method of the LimitedLinkedList type. It is used to check if the list is
// empty.
//
// Returns:
//
//   - bool: A boolean value that is true if the list is empty, and false otherwise.
func (list *LimitedLinkedList[T]) IsEmpty() bool {
	return list.front == nil
}

// Size is a method of the LimitedLinkedList type. It is used to return the current number
// of elements in the list.
//
// Returns:
//
//   - int: An integer that represents the current number of elements in the list.
func (list *LimitedLinkedList[T]) Size() int {
	return list.size
}

// Capacity is a method of the LimitedLinkedList type. It is used to return the maximum
// number of elements the list can hold.
//
// Returns:
//
//   - optional.Int: An optional integer that represents the maximum number of elements
//     the list can hold.
func (list *LimitedLinkedList[T]) Capacity() (int, bool) {
	return list.capacity, true
}

// Iterator is a method of the LimitedLinkedList type. It is used to return an iterator
// for the list.
//
// Returns:
//
//   - itf.Iterater[T]: An iterator for the list.
func (list *LimitedLinkedList[T]) Iterator() itf.Iterater[T] {
	var builder itf.Builder[T]

	for list_node := list.front; list_node != nil; list_node = list_node.next {
		builder.Append(list_node.value)
	}

	return builder.Build()
}

// Clear is a method of the LimitedLinkedList type. It is used to remove all elements from
// the list.
func (list *LimitedLinkedList[T]) Clear() {
	if list.front == nil {
		return // List is already empty
	}

	// 1. First node
	list.front.prev = nil
	prev := list.front

	// 2. Subsequent nodes
	for node := list.front.next; node != nil; node = node.next {
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

// IsFull is a method of the LimitedLinkedList type. It is used to check if the list is full.
//
// Returns:
//
//   - isFull: A boolean value that is true if the list is full, and false otherwise.
func (list *LimitedLinkedList[T]) IsFull() bool {
	return list.size >= list.capacity
}

// String is a method of the LimitedLinkedList type. It returns a string representation of
// the list with information about its size, capacity, and elements.
//
// Returns:
//
//   - string: A string representation of the list.
func (list *LimitedLinkedList[T]) String() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "LimitedLinkedList[capacity=%d, ", list.capacity)

	if list.front == nil {
		builder.WriteString("size=0, values=[]]")

		return builder.String()
	}

	fmt.Fprintf(&builder, "size=%d, values=[%v", list.size, list.front.value)

	for list_node := list.front.next; list_node != nil; list_node = list_node.next {
		fmt.Fprintf(&builder, ", %v", list_node.value)
	}

	fmt.Fprintf(&builder, "]]")

	return builder.String()
}

// Prepend is a method of the LimitedLinkedList type. It is used to add an element to
// the end of the list.
//
// Panics with an error of type *ErrInvalidOperation if the list is full.
//
// Parameters:
//
//   - value: An element of type T to be added to the list.
func (list *LimitedLinkedList[T]) Prepend(value T) error {
	if list.size >= list.capacity {
		return ll.NewErrFullList(list)
	}

	list_node := &linkedNode[T]{
		value: value,
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

// DeleteLast is a method of the LimitedLinkedList type. It is used to remove and return
// the last element in the list.
//
// Panics with an error of type *ErrCallFailed if the list is empty.
//
// Returns:
//
//   - T: The last element in the list.
func (list *LimitedLinkedList[T]) DeleteLast() (T, error) {
	if list.front == nil {
		return *new(T), ll.NewErrEmptyList(list)
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

	return toRemove.value, nil
}

// PeekLast is a method of the LimitedLinkedList type. It is used to return the last
// element in the list without removing it.
//
// Panics with an error of type *ErrCallFailed if the list is empty.
//
// Returns:
//
//   - value: The last element in the list.
func (list *LimitedLinkedList[T]) PeekLast() (T, error) {
	if list.front == nil {
		return *new(T), ll.NewErrEmptyList(list)
	}

	return list.back.value, nil
}

// CutNilValues is a method of the LimitedLinkedList type. It is used to remove all nil
// values from the list.
func (list *LimitedLinkedList[T]) CutNilValues() {
	if list.front == nil {
		return // List is empty
	}

	if gen.IsNil(list.front.value) && list.front == list.back {
		// Single node
		list.front = nil
		list.back = nil
		list.size = 0

		return
	}

	var toDelete *linkedNode[T] = nil

	// 1. First node
	if gen.IsNil(list.front.value) {
		toDelete = list.front

		list.front = list.front.next
		list.front.prev = nil

		toDelete.next = nil
		list.size--
	}

	prev := list.front

	// 2. Subsequent nodes (except last)
	for node := list.front.next; node.next != nil; node = node.next {
		if !gen.IsNil(node.value) {
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
	if gen.IsNil(list.back.value) {
		list.back = prev
		list.back.next = nil
		list.size--
	}
}

// Slice is a method of the LimitedLinkedList type that returns a slice of type T
//
// Returns:
//
//   - []T: A slice of type T.
func (list *LimitedLinkedList[T]) Slice() []T {
	slice := make([]T, 0, list.size)

	for list_node := list.front; list_node != nil; list_node = list_node.next {
		slice = append(slice, list_node.value)
	}

	return slice
}

// Copy is a method of the LimitedLinkedList type. It is used to create a shallow copy
// of the list.
//
// Returns:
//
//   - itf.Copier: A copy of the list.
func (list *LimitedLinkedList[T]) Copy() itff.Copier {
	listCopy := &LimitedLinkedList[T]{
		size:     list.size,
		capacity: list.capacity,
	}

	if list.front == nil {
		return listCopy
	}

	// First node
	listCopy.front = &linkedNode[T]{
		value: list.front.value,
	}

	prev := listCopy.front

	// Subsequent nodes
	for list_node := list.front.next; list_node != nil; list_node = list_node.next {
		list_node_copy := &linkedNode[T]{
			value: list_node.value,
			prev:  prev,
		}

		prev.next = list_node_copy
		prev = list_node_copy
	}

	if listCopy.front.next != nil {
		listCopy.front.next.prev = listCopy.front
	}

	listCopy.back = prev

	return listCopy
}
