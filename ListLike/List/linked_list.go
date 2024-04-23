package List

import (
	"fmt"
	"strings"

	itf "github.com/PlayerR9/MyGoLib/CustomData/Iterators"
	"github.com/PlayerR9/MyGoLib/ListLike/Common"
	itff "github.com/PlayerR9/MyGoLib/Units/Interfaces"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
)

// LinkedList is a generic type that represents a list data structure with
// or without a limited capacity, implemented using a linked list.
type LinkedList[T any] struct {
	// front and back are pointers to the first and last nodes in the linked list,
	// respectively.
	front, back *Common.ListNode[T]

	// size is the current number of elements in the list.
	size int
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
	list_node := Common.NewListNode(values[0])

	list.front = list_node
	list.back = list_node

	// Subsequent nodes
	for _, element := range values {
		list_node := Common.NewListNode(element)
		list_node.SetPrev(list.back)

		list.back.SetNext(list_node)
		list.back = list_node
	}

	return list
}

// Append is a method of the LinkedList type. It is used to add an element to
// the end of the list.
//
// Panics with an error of type *ErrFullList if the list is full.
//
// Parameters:
//
//   - value: An element of type T to be added to the list.
func (list *LinkedList[T]) Append(value T) {
	list_node := Common.NewListNode(value)

	if list.back == nil {
		list.front = list_node
	} else {
		list.back.SetNext(list_node)
		list_node.SetPrev(list.back)
	}

	list.back = list_node

	list.size++
}

// DeleteFirst is a method of the LinkedList type. It is used to remove and return
// the first element in the list.
//
// Panics with an error of type *ErrCallFailed if the list is empty.
//
// Returns:
//
//   - T: The first element in the list.
func (list *LinkedList[T]) DeleteFirst() T {
	if list.front == nil {
		panic(Common.NewErrEmptyList(list))
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

	return toRemove.Value
}

// PeekFirst is a method of the LinkedList type. It is used to return the first
// element in the list without removing it.
//
// Panics with an error of type *ErrCallFailed if the list is empty.
//
// Returns:
//
//   - value: The first element in the list.
func (list *LinkedList[T]) PeekFirst() T {
	if list.front == nil {
		panic(Common.NewErrEmptyList(list))
	}

	return list.front.Value
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
func (list *LinkedList[T]) Capacity() int {
	return -1
}

// Iterator is a method of the LinkedList type. It is used to return an iterator
// for the list.
//
// Returns:
//
//   - itf.Iterater[T]: An iterator for the list.
func (list *LinkedList[T]) Iterator() itf.Iterater[T] {
	var builder itf.Builder[T]

	for list_node := list.front; list_node != nil; list_node = list_node.Next() {
		builder.Append(list_node.Value)
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

// IsFull is a method of the LinkedList type. It is used to check if the list is full.
//
// Returns:
//
//   - isFull: A boolean value that is true if the list is full, and false otherwise.
func (list *LinkedList[T]) IsFull() bool {
	return false
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

	if list.front == nil {
		builder.WriteString("size=0, values=[]]")

		return builder.String()
	}

	fmt.Fprintf(&builder, "size=%d, values=[%v", list.size, list.front.Value)

	for list_node := list.front.Next(); list_node != nil; list_node = list_node.Next() {
		fmt.Fprintf(&builder, ", %v", list_node.Value)
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
func (list *LinkedList[T]) Prepend(value T) {
	list_node := Common.NewListNode(value)

	if list.front == nil {
		list.back = list_node
	} else {
		list_node.SetNext(list.front)
		list.front.SetPrev(list_node)
	}

	list.front = list_node

	list.size++
}

// DeleteLast is a method of the LinkedList type. It is used to remove and return
// the last element in the list.
//
// Panics with an error of type *ErrCallFailed if the list is empty.
//
// Returns:
//
//   - T: The last element in the list.
func (list *LinkedList[T]) DeleteLast() T {
	if list.front == nil {
		panic(Common.NewErrEmptyList(list))
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

	return toRemove.Value
}

// PeekLast is a method of the LinkedList type. It is used to return the last
// element in the list without removing it.
//
// Panics with an error of type *ErrCallFailed if the list is empty.
//
// Returns:
//
//   - value: The last element in the list.
func (list *LinkedList[T]) PeekLast() T {
	if list.front == nil {
		panic(Common.NewErrEmptyList(list))
	}

	return list.back.Value
}

// CutNilValues is a method of the LinkedList type. It is used to remove all nil
// values from the list.
func (list *LinkedList[T]) CutNilValues() {
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

	var toDelete *Common.ListNode[T] = nil

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

// Slice is a method of the LinkedList type that returns a slice of type T
//
// Returns:
//
//   - []T: A slice of type T.
func (list *LinkedList[T]) Slice() []T {
	slice := make([]T, 0, list.size)

	for list_node := list.front; list_node != nil; list_node = list_node.Next() {
		slice = append(slice, list_node.Value)
	}

	return slice
}

// Copy is a method of the LinkedList type. It is used to create a shallow copy
// of the list.
//
// Returns:
//
//   - itf.Copier: A copy of the list.
func (list *LinkedList[T]) Copy() itff.Copier {
	listCopy := &LinkedList[T]{
		size: list.size,
	}

	if list.front == nil {
		return listCopy
	}

	// First node
	listCopy.front = Common.NewListNode(list.front.Value)

	prev := listCopy.front

	// Subsequent nodes
	for list_node := list.front.Next(); list_node != nil; list_node = list_node.Next() {
		list_node_copy := Common.NewListNode(list_node.Value)
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
