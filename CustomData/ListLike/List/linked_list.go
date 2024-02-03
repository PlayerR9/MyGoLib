package ListLike

import (
	"fmt"
	"strings"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

// LinkedList is a generic type in Go that represents a list data structure implemented
// using a linked list.
type LinkedList[T any] struct {
	// front and back are pointers to the first and last nodes in the linked list,
	// respectively.
	front, back *linkedNode[T]

	// size is the current number of elements in the list.
	size int
}

// NewLinkedList is a function that creates and returns a new instance of a LinkedList.
// It takes a variadic parameter of type T, which represents the initial values to be
// stored in the list.
//
// If no initial values are provided, the function simply returns a new LinkedList with
// all its fields set to their zero values.
//
// If initial values are provided, the function creates a new LinkedList and initializes
// its size. It then creates a linked list of nodes
// from the initial values, with each node holding one value, and sets the front and back
// pointers of the list. The new LinkedList is then returned.
func NewLinkedList[T any](values ...*T) *LinkedList[T] {
	if len(values) == 0 {
		return new(LinkedList[T])
	}

	list := new(LinkedList[T])
	list.size = len(values)

	// First node
	node := &linkedNode[T]{
		value: values[0],
	}

	list.front = node
	list.back = node

	// Subsequent nodes
	for _, element := range values[1:] {
		node = &linkedNode[T]{
			value: element,
			prev:  list.back,
		}

		list.back.next = node
		list.back = node
	}

	return list
}

func (list *LinkedList[T]) Append(value *T) {
	node := &linkedNode[T]{
		value: value,
	}

	if list.front == nil {
		list.front = node
	} else {
		list.back.next = node
		node.prev = list.back
	}

	list.back = node
	list.size++
}

func (list *LinkedList[T]) DeleteFirst() *T {
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

func (list *LinkedList[T]) PeekFirst() *T {
	if list.front == nil {
		panic(ers.NewErrOperationFailed(
			"peek first element", NewErrEmptyList(list),
		))
	}

	return list.front.value
}

func (list *LinkedList[T]) IsEmpty() bool {
	return list.front == nil
}

func (list *LinkedList[T]) Size() int {
	return list.size
}

func (list *LinkedList[T]) ToSlice() []*T {
	slice := make([]*T, 0, list.size)

	for node := list.front; node != nil; node = node.next {
		slice = append(slice, node.value)
	}

	return slice
}

func (list *LinkedList[T]) Clear() {
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

func (list *LinkedList[T]) IsFull() bool {
	return false
}

func (list *LinkedList[T]) String() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "LinkedList[size=%d, values=[", list.size)

	if list.front != nil {
		fmt.Fprintf(&builder, "%v", *list.front.value)

		for node := list.front.next; node != nil; node = node.next {
			fmt.Fprintf(&builder, ", %v", *node.value)
		}
	}

	fmt.Fprintf(&builder, "]]")

	return builder.String()
}

func (list *LinkedList[T]) Prepend(value *T) {
	node := &linkedNode[T]{
		value: value,
	}

	if list.front == nil {
		list.back = node
	} else {
		list.front.prev = node
		node.next = list.front
	}

	list.front = node
	list.size++
}

func (list *LinkedList[T]) DeleteLast() *T {
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

func (list *LinkedList[T]) PeekLast() *T {
	if list.front == nil {
		panic(ers.NewErrOperationFailed(
			"peek last element", NewErrEmptyList(list),
		))
	}

	return list.back.value
}
