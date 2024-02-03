package ListLike

import (
	"fmt"
	"slices"
	"strings"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

// LinkedStack is a generic type in Go that represents a stack data structure implemented
// using a linked list.
type LinkedStack[T any] struct {
	// front is a pointer to the first node in the linked stack.
	front *linkedNode[T]

	// size is the current number of elements in the stack.
	size int
}

// NewLinkedStack is a function that creates and returns a new instance of a LinkedStack.
// It takes a variadic parameter of type T, which represents the initial values to be
// stored in the stack.
//
// If no initial values are provided, the function simply returns a new LinkedStack with
// all its fields set to their zero values.
//
// If initial values are provided, the function creates a new LinkedStack and initializes
// its size. It then creates a linked list of nodes
// from the initial values, with each node holding one value, and sets the front and back
// pointers of the stack. The new LinkedStack is then returned.
func NewLinkedStack[T any](values ...*T) *LinkedStack[T] {
	if len(values) == 0 {
		return new(LinkedStack[T])
	}

	stack := new(LinkedStack[T])
	stack.size = len(values)

	// First node
	node := linkedNode[T]{
		value: values[0],
	}

	stack.front = &node

	// Subsequent nodes
	for _, element := range values[1:] {
		node = linkedNode[T]{
			value: element,
			next:  stack.front,
		}

		stack.front = &node
	}

	return stack
}

func (stack *LinkedStack[T]) Push(value *T) {
	node := linkedNode[T]{
		value: value,
	}

	if stack.front != nil {
		node.next = stack.front
	}

	stack.front = &node
	stack.size++
}

func (stack *LinkedStack[T]) Pop() *T {
	if stack.front == nil {
		panic(ers.NewErrOperationFailed(
			"pop", NewErrEmptyStack(stack),
		))
	}

	var value *T

	value, stack.front = stack.front.value, stack.front.next
	stack.size--

	return value
}

func (stack *LinkedStack[T]) Peek() *T {
	if stack.front == nil {
		panic(ers.NewErrOperationFailed(
			"peek", NewErrEmptyStack(stack),
		))
	}

	return stack.front.value
}

func (stack *LinkedStack[T]) IsEmpty() bool {
	return stack.front == nil
}

func (stack *LinkedStack[T]) Size() int {
	return stack.size
}

func (stack *LinkedStack[T]) ToSlice() []*T {
	slice := make([]*T, 0, stack.size)

	for node := stack.front; node != nil; node = node.next {
		slice = append(slice, node.value)
	}

	slices.Reverse(slice)

	return slice
}

func (stack *LinkedStack[T]) Clear() {
	stack.front = nil
	stack.size = 0
}

func (stack *LinkedStack[T]) IsFull() bool {
	return false
}

func (stack *LinkedStack[T]) String() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "LinkedStack[size=%d, values=[", stack.size)

	if stack.front != nil {
		fmt.Fprintf(&builder, "%v", stack.front.value)

		for stack_node := stack.front.next; stack_node != nil; stack_node = stack_node.next {
			fmt.Fprintf(&builder, ", %v", stack_node.value)
		}
	}

	fmt.Fprintf(&builder, "→]]")

	return builder.String()
}
