package ListLike

import (
	"fmt"
	"strings"

	itf "github.com/PlayerR9/MyGoLib/Interfaces"
	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
	"github.com/markphelps/optional"
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
func NewLinkedStack[T any](values ...T) *LinkedStack[T] {
	stack := new(LinkedStack[T])
	stack.size = len(values)

	if len(values) == 0 {
		return stack
	}

	// First node
	node := &linkedNode[T]{
		value: &values[0],
	}

	stack.front = node

	// Subsequent nodes
	for _, element := range values[1:] {
		node = &linkedNode[T]{
			value: &element,
			next:  stack.front,
		}

		stack.front = node
	}

	return stack
}

// WithCapacity is a method of the LinkedStack type. It is used to set the maximum
// number of elements the stack can hold.
//
// Panics with an error of type *ErrCallFailed if the capacity is already set,
// or with an error of type *ErrInvalidParameter if the provided capacity is negative
// or less than the current number of elements in the stack.
//
// Parameters:
//
//   - capacity: An integer representing the maximum number of elements the stack can hold.
//
// Returns:
//
//   - Stacker[T]: A pointer to the current instance of the LinkedStack.
func (stack *LinkedStack[T]) WithCapacity(capacity int) (Stacker[T], error) {
	if stack.capacity.Present() {
		return nil, fmt.Errorf("capacity is already set to %d", stack.capacity.MustGet())
	}

	if capacity < 0 {
		return nil, ers.NewErrInvalidParameter("capacity").
			Wrap(fmt.Errorf("negative capacity (%d) is not allowed", capacity))
	} else if stack.size > capacity {
		return nil, ers.NewErrInvalidParameter("capacity").Wrap(
			fmt.Errorf("provided capacity (%d) is less than the current number of values (%d)",
				capacity, stack.size),
		)
	}

	stack.capacity = optional.NewInt(capacity)

	return stack, nil
}

// Push is a method of the LinkedStack type. It is used to add an element to
// the end of the stack.
//
// Panics with an error of type *ErrCallFailed if the stack is full.
//
// Parameters:
//
//   - value: The value to be added to the stack.
func (stack *LinkedStack[T]) Push(value T) error {
	if stack.capacity.Present() && stack.size >= stack.capacity.MustGet() {
		return NewErrFullStack(stack)
	}

	node := &linkedNode[T]{
		value: &value,
	}

	if stack.front != nil {
		node.next = stack.front
	}

	stack.front = node
	stack.size++

	return nil
}

// Pop is a method of the LinkedStack type. It is used to remove and return the
// last element in the stack.
//
// Panics with an error of type *ErrCallFailed if the stack is empty.
//
// Returns:
//
//   - T: The value of the last element in the stack.
func (stack *LinkedStack[T]) Pop() (T, error) {
	if stack.front == nil {
		return *new(T), NewErrEmptyStack(stack)
	}

	toRemove := stack.front
	stack.front = stack.front.next

	stack.size--
	toRemove.next = nil

	return *toRemove.value, nil
}

// Peek is a method of the LinkedStack type. It is used to return the last element
// in the stack without removing it.
//
// Panics with an error of type *ErrCallFailed if the stack is empty.
//
// Returns:
//
//   - T: The value of the last element in the stack.
func (stack *LinkedStack[T]) Peek() (T, error) {
	if stack.front == nil {
		return *new(T), NewErrEmptyStack(stack)

	}

	return *stack.front.value, nil
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

// Iterator is a method of the LinkedStack type. It is used to return an iterator
// for the elements in the stack.
//
// Returns:
//
//   - itf.Iterater[T]: An iterator for the elements in the stack.
func (stack *LinkedStack[T]) Iterator() itf.Iterater[T] {
	var builder itf.Builder[T]

	for stack_node := stack.front; stack_node != nil; stack_node = stack_node.next {
		builder.Append(*stack_node.value)
	}

	return builder.Build()
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

	fmt.Fprintf(&builder, "size=%d, values=[%v", stack.size, *stack.front.value)

	for stack_node := stack.front.next; stack_node != nil; stack_node = stack_node.next {
		fmt.Fprintf(&builder, ", %v", *stack_node.value)
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

	if gen.IsNil(*stack.front.value) && stack.front.next == nil {
		// Single node
		stack.front = nil
		stack.size = 0

		return
	}

	var toDelete *linkedNode[T] = nil

	// 1. First node
	if gen.IsNil(*stack.front.value) {
		toDelete = stack.front

		stack.front = stack.front.next

		toDelete.next = nil
		stack.size--
	}

	prev := stack.front

	// 2. Subsequent nodes (except last)
	node := stack.front.next
	for ; node.next != nil; node = node.next {
		if !gen.IsNil(*node.value) {
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

// Slice is a method of the LinkedStack type. It is used to return a slice of the
// elements in the stack.
//
// Returns:
//
//   - []T: A slice of the elements in the stack.
func (stack *LinkedStack[T]) Slice() []T {
	slice := make([]T, 0, stack.size)

	for stack_node := stack.front; stack_node != nil; stack_node = stack_node.next {
		slice = append(slice, *stack_node.value)
	}

	return slice
}

// Copy is a method of the LinkedStack type. It is used to create a shallow copy
// of the stack.
//
// Returns:
//
//   - itf.Copier: A copy of the stack.
func (stack *LinkedStack[T]) Copy() itf.Copier {
	stackCopy := LinkedStack[T]{
		size:     stack.size,
		capacity: stack.capacity,
	}

	if stack.front == nil {
		return &stackCopy
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

	return &stackCopy
}
