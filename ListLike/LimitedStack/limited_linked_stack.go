package ListLike

import (
	"fmt"
	"strings"

	itff "github.com/PlayerR9/MyGoLib/Common/Interfaces"
	itf "github.com/PlayerR9/MyGoLib/CustomData/Iterators"
	ll "github.com/PlayerR9/MyGoLib/ListLike"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
)

// LimitedLinkedStack is a generic type that represents a stack data structure with
// or without a limited capacity, implemented using a linked list.
type LimitedLinkedStack[T any] struct {
	// front is a pointer to the first node in the linked stack.
	front *linkedNode[T]

	// size is the current number of elements in the stack.
	size int

	// capacity is the maximum number of elements the stack can hold.
	capacity int
}

// NewLimitedLinkedStack is a function that creates and returns a new instance of a
// LimitedLinkedStack.
//
// Parameters:
//
//   - values: A variadic parameter of type T, which represents the initial values to be
//     stored in the stack.
//
// Returns:
//
//   - *LimitedLinkedStack[T]: A pointer to the newly created LimitedLinkedStack.
func NewLimitedLinkedStack[T any](values ...T) *LimitedLinkedStack[T] {
	stack := new(LimitedLinkedStack[T])
	stack.size = len(values)

	if len(values) == 0 {
		return stack
	}

	// First node
	node := &linkedNode[T]{
		value: values[0],
	}

	stack.front = node

	// Subsequent nodes
	for _, element := range values[1:] {
		node = &linkedNode[T]{
			value: element,
			next:  stack.front,
		}

		stack.front = node
	}

	return stack
}

// Push is a method of the LimitedLinkedStack type. It is used to add an element to
// the end of the stack.
//
// Panics with an error of type *ErrCallFailed if the stack is full.
//
// Parameters:
//
//   - value: The value to be added to the stack.
func (stack *LimitedLinkedStack[T]) Push(value T) error {
	if stack.size >= stack.capacity {
		return ll.NewErrFullList(stack)
	}

	node := &linkedNode[T]{
		value: value,
	}

	if stack.front != nil {
		node.next = stack.front
	}

	stack.front = node
	stack.size++

	return nil
}

// Pop is a method of the LimitedLinkedStack type. It is used to remove and return the
// last element in the stack.
//
// Panics with an error of type *ErrCallFailed if the stack is empty.
//
// Returns:
//
//   - T: The value of the last element in the stack.
func (stack *LimitedLinkedStack[T]) Pop() (T, error) {
	if stack.front == nil {
		return *new(T), ll.NewErrEmptyList(stack)
	}

	toRemove := stack.front
	stack.front = stack.front.next

	stack.size--
	toRemove.next = nil

	return toRemove.value, nil
}

// Peek is a method of the LimitedLinkedStack type. It is used to return the last element
// in the stack without removing it.
//
// Panics with an error of type *ErrCallFailed if the stack is empty.
//
// Returns:
//
//   - T: The value of the last element in the stack.
func (stack *LimitedLinkedStack[T]) Peek() (T, error) {
	if stack.front == nil {
		return *new(T), ll.NewErrEmptyList(stack)

	}

	return stack.front.value, nil
}

// IsEmpty is a method of the LimitedLinkedStack type. It is used to check if the stack
// is empty.
//
// Returns:
//
//   - bool: true if the stack is empty, and false otherwise.
func (stack *LimitedLinkedStack[T]) IsEmpty() bool {
	return stack.front == nil
}

// Size is a method of the LimitedLinkedStack type. It is used to return the number of
// elements in the stack.
//
// Returns:
//
//   - int: The number of elements in the stack.
func (stack *LimitedLinkedStack[T]) Size() int {
	return stack.size
}

// Capacity is a method of the LimitedLinkedStack type. It is used to return the maximum
// number of elements the stack can hold.
//
// Returns:
//
//   - optional.Int: The maximum number of elements the stack can hold.
func (stack *LimitedLinkedStack[T]) Capacity() (int, bool) {
	return stack.capacity, true
}

// Iterator is a method of the LimitedLinkedStack type. It is used to return an iterator
// for the elements in the stack.
//
// Returns:
//
//   - itf.Iterater[T]: An iterator for the elements in the stack.
func (stack *LimitedLinkedStack[T]) Iterator() itf.Iterater[T] {
	var builder itf.Builder[T]

	for stack_node := stack.front; stack_node != nil; stack_node = stack_node.next {
		builder.Append(stack_node.value)
	}

	return builder.Build()
}

// Clear is a method of the LimitedLinkedStack type. It is used to remove all elements
// from the stack.
func (stack *LimitedLinkedStack[T]) Clear() {
	if stack.front == nil {
		return // Stack is already empty
	}

	// 1. First node
	prev := stack.front

	// 2. Subsequent nodes
	for node := stack.front.next; node != nil; node = node.next {
		prev = node
		prev.next = nil
	}

	prev.next = nil

	// 3. Reset list fields
	stack.front = nil
	stack.size = 0
}

// IsFull is a method of the LimitedLinkedStack type. It is used to check if the stack is
// full.
//
// Returns:
//
//   - isFull: true if the stack is full, and false otherwise.
func (stack *LimitedLinkedStack[T]) IsFull() bool {
	return stack.size >= stack.capacity
}

// String is a method of the LimitedLinkedStack type. It is used to return a string
// representation of the stack, which includes the size, capacity, and elements
// in the stack.
//
// Returns:
//
//   - string: A string representation of the stack.
func (stack *LimitedLinkedStack[T]) String() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "LimitedLinkedStack[capacity=%d, ", stack.capacity)

	if stack.size == 0 {
		builder.WriteString("size=0, values=[ →]]")
		return builder.String()
	}

	fmt.Fprintf(&builder, "size=%d, values=[%v", stack.size, stack.front.value)

	for stack_node := stack.front.next; stack_node != nil; stack_node = stack_node.next {
		fmt.Fprintf(&builder, ", %v", stack_node.value)
	}

	fmt.Fprintf(&builder, " →]]")

	return builder.String()
}

// CutNilValues is a method of the LimitedLinkedStack type. It is used to remove all nil
// values from the stack.
func (stack *LimitedLinkedStack[T]) CutNilValues() {
	if stack.front == nil {
		return // Stack is empty
	}

	if gen.IsNil(stack.front.value) && stack.front.next == nil {
		// Single node
		stack.front = nil
		stack.size = 0

		return
	}

	var toDelete *linkedNode[T] = nil

	// 1. First node
	if gen.IsNil(stack.front.value) {
		toDelete = stack.front

		stack.front = stack.front.next

		toDelete.next = nil
		stack.size--
	}

	prev := stack.front

	// 2. Subsequent nodes (except last)
	node := stack.front.next
	for ; node.next != nil; node = node.next {
		if !gen.IsNil(node.value) {
			prev = node
		} else {
			prev.next = node.next
			stack.size--

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
	if gen.IsNil(node.value) {
		node = prev
		node.next = nil
		stack.size--
	}
}

// Slice is a method of the LimitedLinkedStack type. It is used to return a slice of the
// elements in the stack.
//
// Returns:
//
//   - []T: A slice of the elements in the stack.
func (stack *LimitedLinkedStack[T]) Slice() []T {
	slice := make([]T, 0, stack.size)

	for stack_node := stack.front; stack_node != nil; stack_node = stack_node.next {
		slice = append(slice, stack_node.value)
	}

	return slice
}

// Copy is a method of the LimitedLinkedStack type. It is used to create a shallow copy
// of the stack.
//
// Returns:
//
//   - itf.Copier: A copy of the stack.
func (stack *LimitedLinkedStack[T]) Copy() itff.Copier {
	stackCopy := &LimitedLinkedStack[T]{
		size:     stack.size,
		capacity: stack.capacity,
	}

	if stack.front == nil {
		return stackCopy
	}

	// First node
	node := &linkedNode[T]{
		value: stack.front.value,
	}

	stackCopy.front = node

	// Subsequent nodes
	for stack_node := stack.front.next; stack_node != nil; stack_node = stack_node.next {
		node.next = &linkedNode[T]{
			value: stack_node.value,
		}

		node = node.next
	}

	return stackCopy
}
