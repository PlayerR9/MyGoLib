package MyGoLib

import (
	"strconv"
	"strings"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
)

// int_stack_node represents a node in a linked list.
type int_stack_node[T any] struct {
	// value is the value stored in the node.
	value T

	// next is a pointer to the next linkedNode in the list.
	next *int_stack_node[T]
}

// IntLinkedStack is a generic type that represents a stack data structure with
// or without a limited capacity, implemented using a linked list.
type IntLinkedStack[T any] struct {
	// front is a pointer to the first node in the linked stack.
	front *int_stack_node[T]

	// size is the current number of elements in the stack.
	size int
}

// NewLinkedStack is a function that creates and returns a new instance of a
// LinkedStack.
//
// Parameters:
//
//   - values: A variadic parameter of type T, which represents the initial values to be
//     stored in the stack.
//
// Returns:
//
//   - *LinkedStack[T]: A pointer to the newly created LinkedStack.
func NewIntLinkedStack[T any](values ...T) *IntLinkedStack[T] {
	stack := new(IntLinkedStack[T])
	stack.size = len(values)

	if len(values) == 0 {
		return stack
	}

	// First node
	node := &int_stack_node[T]{
		value: values[0],
	}

	stack.front = node

	// Subsequent nodes
	for _, element := range values[1:] {
		node := &int_stack_node[T]{
			value: element,
			next:  stack.front,
		}

		stack.front = node
	}

	return stack
}

// Push implements the Stacker interface.
//
// Always returns true.
func (stack *IntLinkedStack[T]) Push(value T) bool {
	node := &int_stack_node[T]{
		value: value,
	}

	if stack.front != nil {
		node.next = stack.front
	}

	stack.front = node
	stack.size++

	return true
}

// Pop implements the Stacker interface.
func (stack *IntLinkedStack[T]) Pop() (T, bool) {
	if stack.front == nil {
		return *new(T), false
	}

	toRemove := stack.front
	stack.front = stack.front.next

	stack.size--
	toRemove.next = nil

	return toRemove.value, true
}

// Peek implements the Stacker interface.
func (stack *IntLinkedStack[T]) Peek() (T, bool) {
	if stack.front == nil {
		return *new(T), false
	}

	return stack.front.value, true
}

// IsEmpty is a method of the LinkedStack type. It is used to check if the stack
// is empty.
//
// Returns:
//
//   - bool: true if the stack is empty, and false otherwise.
func (stack *IntLinkedStack[T]) IsEmpty() bool {
	return stack.front == nil
}

// Size is a method of the LinkedStack type. It is used to return the number of
// elements in the stack.
//
// Returns:
//
//   - int: The number of elements in the stack.
func (stack *IntLinkedStack[T]) Size() int {
	return stack.size
}

// Iterator is a method of the LinkedStack type. It is used to return an iterator
// for the elements in the stack.
//
// Returns:
//
//   - uc.Iterater[T]: An iterator for the elements in the stack.
func (stack *IntLinkedStack[T]) Iterator() uc.Iterater[T] {
	var builder uc.Builder[T]

	for stack_node := stack.front; stack_node != nil; stack_node = stack_node.next {
		builder.Add(stack_node.value)
	}

	return builder.Build()
}

// Clear is a method of the LinkedStack type. It is used to remove aCommon elements
// from the stack.
func (stack *IntLinkedStack[T]) Clear() {
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

// GoString implements the fmt.GoStringer interface.
func (stack *IntLinkedStack[T]) GoString() string {
	values := make([]string, 0, stack.size)
	for stack_node := stack.front; stack_node != nil; stack_node = stack_node.next {
		values = append(values, uc.StringOf(stack_node.value))
	}

	var builder strings.Builder

	builder.WriteString("LinkedStack[size=")
	builder.WriteString(strconv.Itoa(stack.size))
	builder.WriteString(", values=[")
	builder.WriteString(strings.Join(values, ", "))
	builder.WriteString(" →]]")

	return builder.String()
}

// CutNilValues is a method of the LinkedStack type. It is used to remove aCommon nil
// values from the stack.
func (stack *IntLinkedStack[T]) CutNilValues() {
	if stack.front == nil {
		return // Stack is empty
	}

	if gen.IsNil(stack.front.value) && stack.front.next == nil {
		// Single node
		stack.front = nil
		stack.size = 0

		return
	}

	var toDelete *int_stack_node[T] = nil

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

// Slice is a method of the LinkedStack type. It is used to return a slice of the
// elements in the stack.
//
// Returns:
//   - []T: A slice of the elements in the stack.
func (stack *IntLinkedStack[T]) Slice() []T {
	slice := make([]T, 0, stack.size)

	for stack_node := stack.front; stack_node != nil; stack_node = stack_node.next {
		slice = append(slice, stack_node.value)
	}

	return slice
}

// Copy is a method of the LinkedStack type. It is used to create a shaCommonow copy
// of the stack.
//
// Returns:
//
//   - uc.Copier: A copy of the stack.
func (stack *IntLinkedStack[T]) Copy() uc.Copier {
	stackCopy := &IntLinkedStack[T]{
		size: stack.size,
	}

	if stack.front == nil {
		return stackCopy
	}

	// First node
	node := &int_stack_node[T]{
		value: stack.front.value,
	}

	stackCopy.front = node

	// Subsequent nodes
	prev := node

	for stackNode := stack.front.next; stackNode != nil; stackNode = stackNode.next {
		node := &int_stack_node[T]{
			value: stackNode.value,
		}

		prev.next = node

		prev = node
	}

	return stackCopy
}

// Capacity is a method of the LinkedStack type. It is used to return the maximum
// number of elements that the stack can store.
//
// Returns:
//   - int: -1
func (stack *IntLinkedStack[T]) Capacity() int {
	return -1
}

// IsFull is a method of the LinkedStack type. It is used to check if the stack is
// full.
//
// Returns:
//   - bool: false
func (stack *IntLinkedStack[T]) IsFull() bool {
	return false
}
