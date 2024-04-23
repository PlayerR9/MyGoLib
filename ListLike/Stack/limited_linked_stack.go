package Stack

import (
	"fmt"
	"strings"

	itf "github.com/PlayerR9/MyGoLib/CustomData/Iterators"
	"github.com/PlayerR9/MyGoLib/ListLike/Common"
	itff "github.com/PlayerR9/MyGoLib/Units/Interfaces"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
)

// LimitedLinkedStack is a generic type that represents a stack data structure with
// or without a limited capacity, implemented using a linked list.
type LimitedLinkedStack[T any] struct {
	// front is a pointer to the first node in the linked stack.
	front *Common.StackNode[T]

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
	node := Common.NewStackNode(values[0])

	stack.front = &node

	// Subsequent nodes
	for _, element := range values[1:] {
		node = Common.NewStackNode(element)
		node.SetNext(stack.front)

		stack.front = &node
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
func (stack *LimitedLinkedStack[T]) Push(value T) {
	if stack.size >= stack.capacity {
		panic(Common.NewErrFullList(stack))
	}

	node := Common.NewStackNode(value)

	if stack.front != nil {
		node.SetNext(stack.front)
	}

	stack.front = &node
	stack.size++
}

// Pop is a method of the LimitedLinkedStack type. It is used to remove and return the
// last element in the stack.
//
// Panics with an error of type *ErrCallFailed if the stack is empty.
//
// Returns:
//
//   - T: The value of the last element in the stack.
func (stack *LimitedLinkedStack[T]) Pop() T {
	if stack.front == nil {
		panic(Common.NewErrEmptyList(stack))
	}

	toRemove := stack.front
	stack.front = stack.front.Next()

	stack.size--
	toRemove.SetNext(nil)

	return toRemove.Value
}

// Peek is a method of the LimitedLinkedStack type. It is used to return the last element
// in the stack without removing it.
//
// Panics with an error of type *ErrCallFailed if the stack is empty.
//
// Returns:
//
//   - T: The value of the last element in the stack.
func (stack *LimitedLinkedStack[T]) Peek() T {
	if stack.front == nil {
		panic(Common.NewErrEmptyList(stack))
	}

	return stack.front.Value
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
func (stack *LimitedLinkedStack[T]) Capacity() int {
	return stack.capacity
}

// Iterator is a method of the LimitedLinkedStack type. It is used to return an iterator
// for the elements in the stack.
//
// Returns:
//
//   - itf.Iterater[T]: An iterator for the elements in the stack.
func (stack *LimitedLinkedStack[T]) Iterator() itf.Iterater[T] {
	var builder itf.Builder[T]

	for stack_node := stack.front; stack_node != nil; stack_node = stack_node.Next() {
		builder.Append(stack_node.Value)
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
	for node := stack.front.Next(); node != nil; node = node.Next() {
		prev = node
		prev.SetNext(nil)
	}

	prev.SetNext(nil)

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

	fmt.Fprintf(&builder, "size=%d, values=[%v", stack.size, stack.front.Value)

	for stack_node := stack.front.Next(); stack_node != nil; stack_node = stack_node.Next() {
		fmt.Fprintf(&builder, ", %v", stack_node.Value)
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

	if gen.IsNil(stack.front.Value) && stack.front.Next() == nil {
		// Single node
		stack.front = nil
		stack.size = 0

		return
	}

	var toDelete *Common.StackNode[T] = nil

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

// Slice is a method of the LimitedLinkedStack type. It is used to return a slice of the
// elements in the stack.
//
// Returns:
//
//   - []T: A slice of the elements in the stack.
func (stack *LimitedLinkedStack[T]) Slice() []T {
	slice := make([]T, 0, stack.size)

	for stack_node := stack.front; stack_node != nil; stack_node = stack_node.Next() {
		slice = append(slice, stack_node.Value)
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
	// FIXME: This doesn't work: Node.SetNext(Node)!!!

	stackCopy := &LimitedLinkedStack[T]{
		size:     stack.size,
		capacity: stack.capacity,
	}

	if stack.front == nil {
		return stackCopy
	}

	// First node
	node := Common.NewStackNode(stack.front.Value)

	stackCopy.front = &node

	// Subsequent nodes
	for stack_node := stack.front.Next(); stack_node != nil; stack_node = stack_node.Next() {
		node := Common.NewStackNode(stack_node.Value)
		node.SetNext(&node)
	}

	return stackCopy
}
