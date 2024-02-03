package ListLike

import (
	"fmt"
	"strings"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

// LimitedLinkedList is a generic type in Go that represents a list data structure with
// a limited capacity, implemented using a linked list.
type LimitedLinkedList[T any] struct {
	// front and back are pointers to the first and last nodes in the linked list,
	// respectively.
	front, back *linkedNode[T]

	// size is the current number of elements in the list. capacity is the maximum
	// number of elements the list can hold.
	size, capacity int
}

// NewLimitedLinkedList is a function that creates and returns a new instance of a
// LimitedLinkedList.
// It takes an integer capacity, which represents the maximum number of elements the
// list can hold, and a variadic parameter of type T,
// which represents the initial values to be stored in the list.
//
// The function first checks if the provided capacity is negative. If it is, it returns
// an error of type ErrNegativeCapacity.
// It then checks if the number of initial values exceeds the provided capacity. If it
// does, it returns an error of type ErrTooManyValues.
//
// If the provided capacity and initial values are valid, the function creates a new
// LimitedLinkedList and initializes its size and capacity.
// It then creates a linked list of nodes from the initial values, with each node
// holding one value, and sets the front and back pointers of the list.
// The new LimitedLinkedList is then returned.
func NewLimitedLinkedList[T any](capacity int, values ...*T) (*LimitedLinkedList[T], error) {
	if capacity < 0 {
		return nil, ers.NewErrInvalidParameter(
			"capacity", fmt.Errorf("negative capacity (%d) is not allowed", capacity),
		)
	} else if len(values) > capacity {
		return nil, ers.NewErrInvalidParameter(
			"values", fmt.Errorf("number of values (%d) exceeds the provided capacity (%d)",
				len(values), capacity),
		)
	}

	list := new(LimitedLinkedList[T])
	list.size = len(values)
	list.capacity = capacity

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

	return list, nil
}

// Append is a method of the LimitedLinkedList type. It is used to add an element to
// the end of the list.
//
// The method takes a parameter, value, of a generic type T, which is the element to be
// added to the list.
//
// Before adding the element, the method checks if the current size of the list is equal
// to or greater than its capacity.
// If it is, it means the list is full, and the method panics by throwing an ErrFullList
// error.
//
// If the list is not full, the method creates a new linkedNode with the provided value.
// If the list is currently empty (i.e., list.back is nil),
// the new node is set as both the front and back of the list. If the list is not empty,
// the new node is added to the end of the list by setting it
// as the next node of the current back node, and then updating the back pointer of the
// list to the new node.
//
// Finally, the size of the list is incremented by 1 to reflect the addition of the
// new element.
func (list *LimitedLinkedList[T]) Append(value *T) {
	if list.size >= list.capacity {
		panic(ers.NewErrOperationFailed(
			"append element", NewErrFullList(list),
		))
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
}

func (list *LimitedLinkedList[T]) DeleteFirst() *T {
	if list.front == nil {
		panic(ers.NewErrOperationFailed(
			"delete first element", NewErrEmptyList(list),
		))
	}

	var value *T

	value, list.front = list.front.value, list.front.next
	if list.front == nil {
		list.back = nil
	} else {
		list.front.prev = nil
	}

	list.size--

	return value
}

func (list *LimitedLinkedList[T]) PeekFirst() *T {
	if list.front == nil {
		panic(ers.NewErrOperationFailed(
			"peek first element", NewErrEmptyList(list),
		))
	}

	return list.front.value
}

func (list *LimitedLinkedList[T]) IsEmpty() bool {
	return list.front == nil
}

func (list *LimitedLinkedList[T]) Size() int {
	return list.size
}

func (list *LimitedLinkedList[T]) ToSlice() []*T {
	slice := make([]*T, 0, list.size)

	for list_node := list.front; list_node != nil; list_node = list_node.next {
		slice = append(slice, list_node.value)
	}

	return slice
}

func (list *LimitedLinkedList[T]) Clear() {
	if list.front == nil {
		return // List is already empty
	}

	// 1. First node
	node := list.front
	node.value = nil

	// 2. Subsequent nodes
	for node = node.next; node != nil; node = node.next {
		node.value = nil
		node.prev.next = nil
		node.prev = nil
	}

	// 3. Reset list fields
	list.front = nil
	list.back = nil
	list.size = 0
}

func (list *LimitedLinkedList[T]) IsFull() bool {
	return list.size >= list.capacity
}

func (list *LimitedLinkedList[T]) String() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "LimitedLinkedList[size=%d, capacity=%d, values=[", list.size, list.capacity)

	if list.front != nil {
		fmt.Fprintf(&builder, "%v", list.front.value)

		for list_node := list.front.next; list_node != nil; list_node = list_node.next {
			fmt.Fprintf(&builder, ", %v", list_node.value)
		}
	}

	fmt.Fprintf(&builder, "]]")

	return builder.String()
}

// Prepend is a method of the LimitedLinkedList type. It is used to add an element to
// the end of the list.
//
// The method takes a parameter, value, of a generic type T, which is the element to be
// added to the list.
//
// Before adding the element, the method checks if the current size of the list is equal
// to or greater than its capacity.
// If it is, it means the list is full, and the method panics by throwing an ErrFullList
// error.
//
// If the list is not full, the method creates a new linkedNode with the provided value.
// If the list is currently empty (i.e., list.back is nil),
// the new node is set as both the front and back of the list. If the list is not empty,
// the new node is added to the end of the list by setting it
// as the next node of the current back node, and then updating the back pointer of the
// list to the new node.
//
// Finally, the size of the list is incremented by 1 to reflect the addition of the
// new element.
func (list *LimitedLinkedList[T]) Prepend(value *T) {
	if list.size >= list.capacity {
		panic(ers.NewErrOperationFailed(
			"prepend element", NewErrFullList(list),
		))
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
}

func (list *LimitedLinkedList[T]) DeleteLast() *T {
	if list.front == nil {
		panic(ers.NewErrOperationFailed(
			"delete last element", NewErrEmptyList(list),
		))
	}

	var value *T

	value, list.back = list.back.value, list.back.prev

	if list.back == nil {
		list.front = nil
	} else {
		list.back.next = nil
	}

	list.size--

	return value
}

func (list *LimitedLinkedList[T]) PeekLast() *T {
	if list.front == nil {
		panic(ers.NewErrOperationFailed(
			"peek last element", NewErrEmptyList(list),
		))
	}

	return list.back.value
}
