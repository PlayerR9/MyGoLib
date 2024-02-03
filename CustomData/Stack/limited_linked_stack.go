package Stack

import (
	"fmt"
	"strings"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
	"golang.org/x/exp/slices"
)

// LimitedLinkedStack is a generic type in Go that represents a stack data structure with
// a limited capacity, implemented using a linked list.
type LimitedLinkedStack[T any] struct {
	// front is a pointer to the first node in the linked stack.
	front *linkedNode[T]

	// size is the current number of elements in the stack. capacity is the maximum
	// number of elements the stack can hold.
	size, capacity int
}

// NewLimitedLinkedStack is a function that creates and returns a new instance of a
// LimitedLinkedStack.
// It takes an integer capacity, which represents the maximum number of elements the
// stack can hold, and a variadic parameter of type T,
// which represents the initial values to be stored in the stack.
//
// The function first checks if the provided capacity is negative. If it is, it returns
// an error of type ErrNegativeCapacity.
// It then checks if the number of initial values exceeds the provided capacity. If it
// does, it returns an error of type ErrTooManyValues.
//
// If the provided capacity and initial values are valid, the function creates a new
// LimitedLinkedStack and initializes its size and capacity.
// It then creates a linked list of nodes from the initial values, with each node
// holding one value, and sets the front and back pointers of the stack.
// The new LimitedLinkedStack is then returned.
func NewLimitedLinkedStack[T any](capacity int, values ...*T) (*LimitedLinkedStack[T], error) {
	if capacity < 0 {
		return nil, ers.NewErrInvalidParameter(
			"capacity", new(ErrNegativeCapacity),
		)
	} else if len(values) > capacity {
		return nil, ers.NewErrInvalidParameter(
			"values", new(ErrTooManyValues),
		)
	}

	stack := new(LimitedLinkedStack[T])
	stack.size = len(values)
	stack.capacity = capacity

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

	return stack, nil
}

// Push is a method of the LimitedLinkedStack type. It is used to add an element to
// the end of the stack.
//
// The method takes a parameter, value, of a generic type T, which is the element to be
// added to the stack.
//
// Before adding the element, the method checks if the current size of the stack is equal
// to or greater than its capacity.
// If it is, it means the stack is full, and the method panics by throwing an ErrFullStack
// error.
//
// If the stack is not full, the method creates a new linkedNode with the provided value.
// If the stack is currently empty (i.e., stack.back is nil),
// the new node is set as both the front and back of the stack. If the stack is not empty,
// the new node is added to the end of the stack by setting it
// as the next node of the current back node, and then updating the back pointer of the
// stack to the new node.
//
// Finally, the size of the stack is incremented by 1 to reflect the addition of the
// new element.
func (stack *LimitedLinkedStack[T]) Push(value *T) {
	if stack.size >= stack.capacity {
		panic(new(ErrFullStack))
	}

	stack_node := linkedNode[T]{
		value: value,
	}

	if stack.front != nil {
		stack_node.next = stack.front
	}

	stack.front = &stack_node

	stack.size++
}

func (stack *LimitedLinkedStack[T]) Pop() *T {
	if stack.front == nil {
		panic(NewErrEmptyStack(Pop))
	}

	var value *T

	value, stack.front = stack.front.value, stack.front.next

	stack.size--

	return value
}

func (stack *LimitedLinkedStack[T]) Peek() *T {
	if stack.front == nil {
		panic(NewErrEmptyStack(Peek))
	}

	return stack.front.value
}

func (stack *LimitedLinkedStack[T]) IsEmpty() bool {
	return stack.front == nil
}

func (stack *LimitedLinkedStack[T]) Size() int {
	return stack.size
}

func (stack *LimitedLinkedStack[T]) ToSlice() []*T {
	slice := make([]*T, 0, stack.size)

	for stack_node := stack.front; stack_node != nil; stack_node = stack_node.next {
		slice = append(slice, stack_node.value)
	}

	slices.Reverse(slice)

	return slice
}

func (stack *LimitedLinkedStack[T]) Clear() {
	stack.front = nil
	stack.size = 0
}

func (stack *LimitedLinkedStack[T]) IsFull() bool {
	return stack.size >= stack.capacity
}

func (stack *LimitedLinkedStack[T]) String() string {
	if stack.front == nil {
		return StackHead
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "%s%v", StackHead, stack.front.value)

	for stack_node := stack.front.next; stack_node != nil; stack_node = stack_node.next {
		fmt.Fprintf(&builder, "%s%v", StackSep, stack_node.value)
	}

	return builder.String()
}
