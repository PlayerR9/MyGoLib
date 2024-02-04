package ListLike

import (
	"fmt"
	"strings"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
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
func NewLinkedList[T any](values ...*T) *LinkedList[T] {
	list := new(LinkedList[T])

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

// WithCapacity is a method of the LinkedList type. It is used to set the maximum
// number of elements the list can hold.
//
// Panics with an error of type *ErrOperationFailed if the capacity is already set
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
func (list *LinkedList[T]) WithCapacity(capacity int) *LinkedList[T] {
	list.capacity.If(func(cap int) {
		panic(ers.NewErrOperationFailed(
			"set capacity", fmt.Errorf("capacity is already set to %d", cap),
		))
	})

	if capacity < 0 {
		panic(ers.NewErrInvalidParameter(
			"capacity", fmt.Errorf("negative capacity (%d) is not allowed", capacity),
		))
	} else if list.size > capacity {
		panic(ers.NewErrInvalidParameter(
			"values", fmt.Errorf("capacity (%d) is not big enough to hold %d elements",
				capacity, list.size),
		))
	}

	list.capacity = optional.NewInt(capacity)

	return list
}

// Append is a method of the LinkedList type. It is used to add an element to
// the end of the list.
//
// Panics with an error of type *ErrFullList if the list is full.
//
// Parameters:
//
//   - value: A pointer to an element of type T to be added to the list.
func (list *LinkedList[T]) Append(value *T) {
	list.capacity.If(func(cap int) {
		ers.Check(list.size < cap, ers.NewErrOperationFailed(
			"append element", NewErrFullList(list),
		))
	})

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
}

// DeleteFirst is a method of the LinkedList type. It is used to remove and return
// the first element in the list.
//
// Panics with an error of type *ErrOperationFailed if the list is empty.
//
// Returns:
//
//   - *T: A pointer to the value of the first element in the list.
func (list *LinkedList[T]) DeleteFirst() *T {
	if list.front == nil {
		panic(ers.NewErrOperationFailed(
			"delete first element", NewErrEmptyList(list),
		))
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

	return toRemove.value
}

// PeekFirst is a method of the LinkedList type. It is used to return the first
// element in the list without removing it.
//
// Panics with an error of type *ErrOperationFailed if the list is empty.
//
// Returns:
//
//   - value: A pointer to the value of the first element in the list.
func (list *LinkedList[T]) PeekFirst() *T {
	if list.front != nil {
		return list.front.value
	}

	panic(ers.NewErrOperationFailed(
		"peek first element", NewErrEmptyList(list),
	))
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

// ToSlice is a method of the LinkedList type. It is used to convert the list to a
// slice of pointers to its elements.
//
// Returns:
//
//   - []*T: A slice of pointers to the elements in the list.
func (list *LinkedList[T]) ToSlice() []*T {
	slice := make([]*T, 0, list.size)

	for list_node := list.front; list_node != nil; list_node = list_node.next {
		slice = append(slice, list_node.value)
	}

	return slice
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

	fmt.Fprintf(&builder, "size=%d, values=[%v", list.size, list.front.value)

	for list_node := list.front.next; list_node != nil; list_node = list_node.next {
		fmt.Fprintf(&builder, ", %v", list_node.value)
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
//   - value: A pointer to an element of type T to be added to the list.
func (list *LinkedList[T]) Prepend(value *T) {
	list.capacity.If(func(cap int) {
		ers.Check(list.size < cap, ers.NewErrOperationFailed(
			"prepend element", NewErrFullList(list),
		))
	})

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
}

// DeleteLast is a method of the LinkedList type. It is used to remove and return
// the last element in the list.
//
// Panics with an error of type *ErrOperationFailed if the list is empty.
//
// Returns:
//
//   - *T: A pointer to the value of the last element in the list.
func (list *LinkedList[T]) DeleteLast() *T {
	if list.front == nil {
		panic(ers.NewErrOperationFailed(
			"delete last element", NewErrEmptyList(list),
		))
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

	return toRemove.value
}

// PeekLast is a method of the LinkedList type. It is used to return the last
// element in the list without removing it.
//
// Panics with an error of type *ErrOperationFailed if the list is empty.
//
// Returns:
//
//   - value: A pointer to the value of the last element in the list.
func (list *LinkedList[T]) PeekLast() *T {
	if list.front != nil {
		return list.back.value

	}

	panic(ers.NewErrOperationFailed(
		"peek last element", NewErrEmptyList(list),
	))
}

// CutNilValues is a method of the LinkedList type. It is used to remove all nil
// values from the list.
func (list *LinkedList[T]) CutNilValues() {
	if list.front == nil {
		return // List is empty
	}

	if list.front.value == nil && list.front == list.back {
		// Single node
		list.front = nil
		list.back = nil
		list.size = 0

		return
	}

	var toDelete *linkedNode[T] = nil

	// 1. First node
	if list.front.value == nil {
		toDelete = list.front

		list.front = list.front.next
		list.front.prev = nil

		toDelete.next = nil
		list.size--
	}

	prev := list.front

	// 2. Subsequent nodes (except last)
	for node := list.front.next; node.next != nil; node = node.next {
		if node.value != nil {
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
