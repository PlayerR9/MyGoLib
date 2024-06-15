package Lister

import (
	"strconv"
	"strings"

	itf "github.com/PlayerR9/MyGoLib/Units/Iterators"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
)

// LimitedLinkedList is a generic type that represents a list data structure with
// or without a limited capacity, implemented using a linked list.
type LimitedLinkedList[T any] struct {
	// front and back are pointers to the first and last nodes in the linked list,
	// respectively.
	front, back *ListNode[T]

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
	list_node := NewListNode(values[0])

	list.front = list_node
	list.back = list_node

	// Subsequent nodes
	for _, element := range values {
		list_node := NewListNode(element)
		list_node.SetPrev(list.back)

		list.back.SetNext(list_node)

		list.back = list_node
	}

	return list
}

// Append implements the Lister interface.
func (list *LimitedLinkedList[T]) Append(value T) bool {
	if list.size >= list.capacity {
		return false
	}

	list_node := NewListNode(value)

	if list.back == nil {
		list.front = list_node
	} else {
		list.back.SetNext(list_node)
		list_node.SetPrev(list.back)
	}

	list.back = list_node

	list.size++

	return true
}

// DeleteFirst implements the Lister interface.
func (list *LimitedLinkedList[T]) DeleteFirst() (T, bool) {
	if list.front == nil {
		return *new(T), false
	}

	toRemove := list.front
	list.front = list.front.Next()

	if list.front == nil {
		list.back = nil
	} else {
		list.front.SetPrev(nil)
	}

	list.size--

	toRemove.SetNext(nil)

	return toRemove.Value, true
}

// PeekFirst implements the Lister interface.
func (list *LimitedLinkedList[T]) PeekFirst() (T, bool) {
	if list.front == nil {
		return *new(T), false
	}

	return list.front.Value, true
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
func (list *LimitedLinkedList[T]) Capacity() int {
	return list.capacity
}

// Iterator is a method of the LimitedLinkedList type. It is used to return an iterator
// for the list.
//
// Returns:
//
//   - itf.Iterater[T]: An iterator for the list.
func (list *LimitedLinkedList[T]) Iterator() itf.Iterater[T] {
	var builder itf.Builder[T]

	for list_node := list.front; list_node != nil; list_node = list_node.Next() {
		builder.Add(list_node.Value)
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
	list.front.SetPrev(nil)
	prev := list.front

	// 2. Subsequent nodes
	for node := list.front.Next(); node != nil; node = node.Next() {
		node.SetPrev(nil)

		prev = node
		prev.SetNext(nil)
	}

	prev.SetNext(nil)

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

// GoString implements the fmt.GoStringer interface.
func (list *LimitedLinkedList[T]) GoString() string {
	values := make([]string, 0, list.size)

	for list_node := list.front; list_node != nil; list_node = list_node.Next() {
		values = append(values, uc.StringOf(list_node.Value))
	}

	var builder strings.Builder

	builder.WriteString("LimitedLinkedList[capacity=")
	builder.WriteString(strconv.Itoa(list.capacity))
	builder.WriteString(", size=")
	builder.WriteString(strconv.Itoa(list.size))
	builder.WriteString(", values=[")
	builder.WriteString(strings.Join(values, ", "))
	builder.WriteString("]]")

	return builder.String()
}

// Prepend implements the Lister interface.
func (list *LimitedLinkedList[T]) Prepend(value T) bool {
	if list.size >= list.capacity {
		return false
	}

	list_node := NewListNode(value)

	if list.front == nil {
		list.back = list_node
	} else {
		list_node.SetNext(list.front)
		list.front.SetPrev(list_node)
	}

	list.front = list_node

	list.size++

	return true
}

// DeleteLast implements the Lister interface.
func (list *LimitedLinkedList[T]) DeleteLast() (T, bool) {
	if list.front == nil {
		return *new(T), false
	}

	toRemove := list.back
	list.back = list.back.Prev()

	if list.back == nil {
		list.front = nil
	} else {
		list.back.SetNext(nil)
	}

	list.size--

	toRemove.SetPrev(nil)

	return toRemove.Value, true
}

// PeekLast implements the Lister interface.
func (list *LimitedLinkedList[T]) PeekLast() (T, bool) {
	if list.front == nil {
		return *new(T), false
	}

	return list.back.Value, true
}

// CutNilValues is a method of the LimitedLinkedList type. It is used to remove all nil
// values from the list.
func (list *LimitedLinkedList[T]) CutNilValues() {
	if list.front == nil {
		return // List is empty
	}

	if gen.IsNil(list.front.Value) && list.front == list.back {
		// Single node
		list.front = nil
		list.back = nil
		list.size = 0

		return
	}

	var toDelete *ListNode[T] = nil

	// 1. First node
	if gen.IsNil(list.front.Value) {
		toDelete = list.front

		list.front = list.front.Next()
		list.front.SetPrev(nil)

		toDelete.SetNext(nil)
		list.size--
	}

	prev := list.front

	// 2. Subsequent nodes (except last)
	for node := list.front.Next(); node.Next() != nil; node = node.Next() {
		if !gen.IsNil(node.Value) {
			prev = node
		} else {
			prev.SetNext(node.Next())
			node.Next().SetPrev(prev)
			list.size--

			if toDelete != nil {
				toDelete.SetNext(nil)
			}

			toDelete = node
		}
	}

	if toDelete != nil {
		toDelete.SetNext(nil)
	}

	// 3. Last node
	if gen.IsNil(list.back.Value) {
		list.back = prev
		list.back.SetNext(nil)
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

	for list_node := list.front; list_node != nil; list_node = list_node.Next() {
		slice = append(slice, list_node.Value)
	}

	return slice
}

// Copy is a method of the LimitedLinkedList type. It is used to create a shallow copy
// of the list.
//
// Returns:
//
//   - itf.Copier: A copy of the list.
func (list *LimitedLinkedList[T]) Copy() uc.Copier {
	listCopy := &LimitedLinkedList[T]{
		size:     list.size,
		capacity: list.capacity,
	}

	if list.front == nil {
		return listCopy
	}

	// First node
	node := NewListNode(list.front.Value)

	listCopy.front = node

	prev := listCopy.front

	// Subsequent nodes
	for list_node := list.front.Next(); list_node != nil; list_node = list_node.Next() {
		list_node_copy := NewListNode(list_node.Value)
		list_node_copy.SetPrev(prev)

		prev.SetNext(list_node_copy)
		prev = list_node_copy
	}

	if listCopy.front.Next() != nil {
		listCopy.front.Next().SetPrev(listCopy.front)
	}

	listCopy.back = prev

	return listCopy
}
