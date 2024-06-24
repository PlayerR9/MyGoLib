package Stacker

import (
	"strconv"
	"strings"

	itf "github.com/PlayerR9/MyGoLib/Units/Iterators"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
)

// LinkedStack is a generic type that represents a stack data structure with
// or without a limited capacity, implemented using a linked list.
type LinkedStack[T any] struct {
	// front is a pointer to the first node in the linked stack.
	front *StackNode[T]

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
func NewLinkedStack[T any](values ...T) *LinkedStack[T] {
	stack := new(LinkedStack[T])
	stack.size = len(values)

	if len(values) == 0 {
		return stack
	}

	// First node
	node := NewStackNode(values[0])

	stack.front = node

	// Subsequent nodes
	for _, element := range values[1:] {
		node = NewStackNode(element)
		node.SetNext(stack.front)

		stack.front = node
	}

	return stack
}

// Push implements the Stacker interface.
//
// Always returns true.
func (stack *LinkedStack[T]) Push(value T) bool {
	node := NewStackNode(value)

	if stack.front != nil {
		node.SetNext(stack.front)
	}

	stack.front = node
	stack.size++

	return true
}

// Pop implements the Stacker interface.
func (stack *LinkedStack[T]) Pop() (T, bool) {
	if stack.front == nil {
		return *new(T), false
	}

	toRemove := stack.front
	stack.front = stack.front.Next()

	stack.size--
	toRemove.SetNext(nil)

	return toRemove.Value, true
}

// Peek implements the Stacker interface.
func (stack *LinkedStack[T]) Peek() (T, bool) {
	if stack.front == nil {
		return *new(T), false
	}

	return stack.front.Value, true
}

// IsEmpty is a method of the LinkedStack type. It is used to check if the stack
// is empty.
//
// Returns:
//
//   - bool: true if the stack is empty, and false otherwise.
func (stack *LinkedStack[T]) IsEmpty() bool {
	return stack.front == nil
}

// Size is a method of the LinkedStack type. It is used to return the number of
// elements in the stack.
//
// Returns:
//
//   - int: The number of elements in the stack.
func (stack *LinkedStack[T]) Size() int {
	return stack.size
}

// Iterator is a method of the LinkedStack type. It is used to return an iterator
// for the elements in the stack.
//
// Returns:
//
//   - itf.Iterater[T]: An iterator for the elements in the stack.
func (stack *LinkedStack[T]) Iterator() itf.Iterater[T] {
	var builder itf.Builder[T]

	for stack_node := stack.front; stack_node != nil; stack_node = stack_node.Next() {
		builder.Add(stack_node.Value)
	}

	return builder.Build()
}

// Clear is a method of the LinkedStack type. It is used to remove aCommon elements
// from the stack.
func (stack *LinkedStack[T]) Clear() {
	if stack.front == nil {
		return // Stack is already empty
	}

	// 1. First node
	prev := stack.front

	// 2. Subsequent nodes
	for node := stack.front.Next(); node != nil; node = node.Next() {
		prev = node
		prev.SetNext(nil)
	}

	prev.SetNext(nil)

	// 3. Reset list fields
	stack.front = nil
	stack.size = 0
}

// GoString implements the fmt.GoStringer interface.
func (stack *LinkedStack[T]) GoString() string {
	values := make([]string, 0, stack.size)
	for stack_node := stack.front; stack_node != nil; stack_node = stack_node.Next() {
		values = append(values, uc.StringOf(stack_node.Value))
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
func (stack *LinkedStack[T]) CutNilValues() {
	if stack.front == nil {
		return // Stack is empty
	}

	if gen.IsNil(stack.front.Value) && stack.front.Next() == nil {
		// Single node
		stack.front = nil
		stack.size = 0

		return
	}

	var toDelete *StackNode[T] = nil

	// 1. First node
	if gen.IsNil(stack.front.Value) {
		toDelete = stack.front

		stack.front = stack.front.Next()

		toDelete.SetNext(nil)
		stack.size--
	}

	prev := stack.front

	// 2. Subsequent nodes (except last)
	node := stack.front.Next()
	for ; node.Next() != nil; node = node.Next() {
		if !gen.IsNil(node.Value) {
			prev = node
		} else {
			prev.SetNext(node.Next())
			stack.size--

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
	if gen.IsNil(node.Value) {
		node = prev
		node.SetNext(nil)
		stack.size--
	}
}

// Slice is a method of the LinkedStack type. It is used to return a slice of the
// elements in the stack.
//
// Returns:
//   - []T: A slice of the elements in the stack.
func (stack *LinkedStack[T]) Slice() []T {
	slice := make([]T, 0, stack.size)

	for stack_node := stack.front; stack_node != nil; stack_node = stack_node.Next() {
		slice = append(slice, stack_node.Value)
	}

	return slice
}

// Copy is a method of the LinkedStack type. It is used to create a shaCommonow copy
// of the stack.
//
// Returns:
//
//   - itf.Copier: A copy of the stack.
func (stack *LinkedStack[T]) Copy() uc.Copier {
	stackCopy := &LinkedStack[T]{
		size: stack.size,
	}

	if stack.front == nil {
		return stackCopy
	}

	// First node
	node := NewStackNode(stack.front.Value)
	stackCopy.front = node

	// Subsequent nodes
	prev := node

	for stackNode := stack.front.Next(); stackNode != nil; stackNode = stackNode.Next() {
		node := NewStackNode(stackNode.Value)
		prev.SetNext(node)

		prev = node
	}

	return stackCopy
}

// Capacity is a method of the LinkedStack type. It is used to return the maximum
// number of elements that the stack can store.
//
// Returns:
//   - int: -1
func (stack *LinkedStack[T]) Capacity() int {
	return -1
}

// IsFull is a method of the LinkedStack type. It is used to check if the stack is
// full.
//
// Returns:
//   - bool: false
func (stack *LinkedStack[T]) IsFull() bool {
	return false
}
