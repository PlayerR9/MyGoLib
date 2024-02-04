package ListLike

import (
	"fmt"
	"strings"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
	"github.com/markphelps/optional"
	"golang.org/x/exp/slices"
)

// LinkedStack is a generic type that represents a stack data structure with
// or without a limited capacity, implemented using a linked list.
type LinkedStack[T any] struct {
	// front is a pointer to the first node in the linked stack.
	front *linkedNode[T]

	// size is the current number of elements in the stack.
	size int

	// capacity is the maximum number of elements the stack can hold.
	capacity optional.Int
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
func NewLinkedStack[T any](values ...*T) *LinkedStack[T] {
	stack := new(LinkedStack[T])
	stack.size = len(values)

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

// WithCapacity is a method of the LinkedStack type. It is used to set the maximum
// number of elements the stack can hold.
//
// Panics with an error of type *ErrOperationFailed if the capacity is already set,
// or with an error of type *ErrInvalidParameter if the provided capacity is negative
// or less than the current number of elements in the stack.
//
// Parameters:
//
//   - capacity: An integer representing the maximum number of elements the stack can hold.
//
// Returns:
//
//   - *LinkedStack[T]: A pointer to the current instance of the LinkedStack.
func (stack *LinkedStack[T]) WithCapacity(capacity int) *LinkedStack[T] {
	stack.capacity.If(func(cap int) {
		panic(ers.NewErrOperationFailed(
			"set capacity", fmt.Errorf("capacity is already set to %d", cap),
		))
	})

	if capacity < 0 {
		panic(ers.NewErrInvalidParameter(
			"capacity", fmt.Errorf("negative capacity (%d) is not allowed", capacity),
		))
	} else if stack.size > capacity {
		panic(ers.NewErrInvalidParameter(
			"capacity", fmt.Errorf("provided capacity (%d) is less than the current number of values (%d)",
				capacity, stack.size),
		))
	}

	stack.capacity = optional.NewInt(capacity)

	return stack
}

// Push is a method of the LinkedStack type. It is used to add an element to
// the end of the stack.
//
// Panics with an error of type *ErrOperationFailed if the stack is full.
//
// Parameters:
//
//   - value: A pointer to a value of type T, which is the element to be added to the stack.
func (stack *LinkedStack[T]) Push(value *T) {
	stack.capacity.If(func(cap int) {
		ers.Check(stack.size < cap, ers.NewErrOperationFailed(
			"push element", NewErrFullStack(stack),
		))
	})

	node := &linkedNode[T]{
		value: value,
	}

	if stack.front != nil {
		node.next = stack.front
	}

	stack.front = node
	stack.size++
}

// Pop is a method of the LinkedStack type. It is used to remove and return the
// last element in the stack.
//
// Panics with an error of type *ErrOperationFailed if the stack is empty.
//
// Returns:
//
//   - *T: A pointer to the value of the last element in the stack.
func (stack *LinkedStack[T]) Pop() *T {
	if stack.front == nil {
		panic(ers.NewErrOperationFailed(
			"pop element", NewErrEmptyStack(stack),
		))
	}

	toRemove := stack.front
	stack.front = stack.front.next

	stack.size--
	toRemove.next = nil

	return toRemove.value
}

// Peek is a method of the LinkedStack type. It is used to return the last element
// in the stack without removing it.
//
// Panics with an error of type *ErrOperationFailed if the stack is empty.
//
// Returns:
//
//   - *T: A pointer to the value of the last element in the stack.
func (stack *LinkedStack[T]) Peek() *T {
	if stack.front == nil {
		return stack.front.value
	}

	panic(ers.NewErrOperationFailed(
		"peek element", NewErrEmptyStack(stack),
	))
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

// Capacity is a method of the LinkedStack type. It is used to return the maximum
// number of elements the stack can hold.
//
// Returns:
//
//   - optional.Int: The maximum number of elements the stack can hold.
func (stack *LinkedStack[T]) Capacity() optional.Int {
	return stack.capacity
}

// ToSlice is a method of the LinkedStack type. It is used to return the elements
// in the stack as a slice.
//
// Returns:
//
//   - []*T: A slice of pointers to the elements in the stack.
func (stack *LinkedStack[T]) ToSlice() []*T {
	slice := make([]*T, 0, stack.size)

	for stack_node := stack.front; stack_node != nil; stack_node = stack_node.next {
		slice = append(slice, stack_node.value)
	}

	slices.Reverse(slice)

	return slice
}

// Clear is a method of the LinkedStack type. It is used to remove all elements
// from the stack.
func (stack *LinkedStack[T]) Clear() {
	if stack.front == nil {
		return // Stack is already empty
	}

	// 1. First node
	stack.front.value = nil
	prev := stack.front

	// 2. Subsequent nodes
	for node := stack.front.next; node != nil; node = node.next {
		node.value = nil

		prev = node
		prev.next = nil
	}

	prev.next = nil

	// 3. Reset list fields
	stack.front = nil
	stack.size = 0
}

// IsFull is a method of the LinkedStack type. It is used to check if the stack is
// full.
//
// Returns:
//
//   - isFull: true if the stack is full, and false otherwise.
func (stack *LinkedStack[T]) IsFull() (isFull bool) {
	stack.capacity.If(func(cap int) {
		isFull = stack.size >= cap
	})
	return
}

// String is a method of the LinkedStack type. It is used to return a string
// representation of the stack, which includes the size, capacity, and elements
// in the stack.
//
// Returns:
//
//   - string: A string representation of the stack.
func (stack *LinkedStack[T]) String() string {
	var builder strings.Builder

	builder.WriteString("LinkedStack[")

	stack.capacity.If(func(cap int) {
		fmt.Fprintf(&builder, "capacity=%d, ", cap)
	})

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

// CutNilValues is a method of the LinkedStack type. It is used to remove all nil
// values from the stack.
func (stack *LinkedStack[T]) CutNilValues() {
	if stack.front == nil {
		return // Stack is empty
	}

	if stack.front.value == nil && stack.front.next == nil {
		// Single node
		stack.front = nil
		stack.size = 0

		return
	}

	var toDelete *linkedNode[T] = nil

	// 1. First node
	if stack.front.value == nil {
		toDelete = stack.front

		stack.front = stack.front.next

		toDelete.next = nil
		stack.size--
	}

	prev := stack.front

	// 2. Subsequent nodes (except last)
	node := stack.front.next
	for ; node.next != nil; node = node.next {
		if node.value != nil {
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
	if node.value == nil {
		node = prev
		node.next = nil
		stack.size--
	}
}
