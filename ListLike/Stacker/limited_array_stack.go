package Stacker

import (
	"slices"
	"strconv"
	"strings"

	itf "github.com/PlayerR9/MyGoLib/ListLike/Iterator"
	uc "github.com/PlayerR9/MyGoLib/Units/Common"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
)

// LimitedArrayStack is a generic type that represents a stack data structure with
// or without a limited capacity. It is implemented using an array.
type LimitedArrayStack[T any] struct {
	// values is a slice of type T that stores the elements in the stack.
	values []T

	// capacity is the maximum number of elements the stack can hold.
	capacity int
}

// NewLimitedArrayStack is a function that creates and returns a new instance of a
// LimitedArrayStack.
//
// Parameters:
//
//   - values: A variadic parameter of type T, which represents the initial values to be
//     stored in the stack.
//
// Returns:
//
//   - *LimitedArrayStack[T]: A pointer to the newly created LimitedArrayStack.
func NewLimitedArrayStack[T any](values ...T) *LimitedArrayStack[T] {
	slices.Reverse(values)

	stack := &LimitedArrayStack[T]{
		values: make([]T, len(values)),
	}
	copy(stack.values, values)

	return stack
}

// Push is a method of the LimitedArrayStack type. It is used to add an element to the
// end of the stack.
//
// Panics with an error of type *ErrCallFailed if the stack is full.
//
// Parameters:
//
//   - value: The value of type T to be added to the stack.
func (stack *LimitedArrayStack[T]) Push(value T) error {
	if len(stack.values) == stack.capacity {
		return NewErrFullStack(stack)
	}

	stack.values = append(stack.values, value)

	return nil
}

// Pop is a method of the LimitedArrayStack type. It is used to remove and return the
// element at the end of the stack.
//
// Panics with an error of type *ErrCallFailed if the stack is empty.
//
// Returns:
//
//   - T: The element at the end of the stack.
func (stack *LimitedArrayStack[T]) Pop() (T, error) {
	if len(stack.values) == 0 {
		return *new(T), ers.NewErrEmpty(stack)
	}

	toRemove := stack.values[len(stack.values)-1]
	stack.values = stack.values[:len(stack.values)-1]

	return toRemove, nil
}

// Peek is a method of the LimitedArrayStack type. It is used to return the element at the
// end of the stack without removing it.
//
// Panics with an error of type *ErrCallFailed if the stack is empty.
//
// Returns:
//
//   - T: The element at the end of the stack.
func (stack *LimitedArrayStack[T]) Peek() (T, error) {
	if len(stack.values) == 0 {
		return *new(T), ers.NewErrEmpty(stack)
	}

	return stack.values[len(stack.values)-1], nil
}

// IsEmpty is a method of the LimitedArrayStack type. It is used to check if the stack is
// empty.
//
// Returns:
//
//   - bool: A boolean value that is true if the stack is empty, and false otherwise.
func (stack *LimitedArrayStack[T]) IsEmpty() bool {
	return len(stack.values) == 0
}

// Size is a method of the LimitedArrayStack type. It is used to return the number of elements
// in the stack.
//
// Returns:
//
//   - int: An integer that represents the number of elements in the stack.
func (stack *LimitedArrayStack[T]) Size() int {
	return len(stack.values)
}

// Capacity is a method of the LimitedArrayStack type. It is used to return the maximum number
// of elements the stack can hold.
//
// Returns:
//
//   - optional.Int: An optional integer that represents the maximum number of elements
//     the stack can hold.
func (stack *LimitedArrayStack[T]) Capacity() int {
	return stack.capacity
}

// Iterator is a method of the LimitedArrayStack type. It is used to return an iterator that
// iterates over the elements in the stack.
//
// Returns:
//
//   - itf.Iterater[T]: An iterator that iterates over the elements in the stack.
func (stack *LimitedArrayStack[T]) Iterator() itf.Iterater[T] {
	var builder itf.Builder[T]

	for i := len(stack.values) - 1; i >= 0; i-- {
		builder.Append(stack.values[i])
	}

	return builder.Build()
}

// Clear is a method of the LimitedArrayStack type. It is used to remove all elements from the
// stack, making it empty.
func (stack *LimitedArrayStack[T]) Clear() {
	stack.values = make([]T, 0, stack.capacity)
}

// IsFull is a method of the LimitedArrayStack type. It is used to check if the stack is full,
// i.e., if it has reached its maximum capacity.
//
// Returns:
//
//   - isFull: A boolean value that is true if the stack is full, and false otherwise.
func (stack *LimitedArrayStack[T]) IsFull() (isFull bool) {
	return len(stack.values) == stack.capacity
}

// String is a method of the LimitedArrayStack type. It is used to return a string representation
// of the stack, including its capacity and the elements it contains.
//
// Returns:
//
//   - string: A string representation of the stack.
func (stack *LimitedArrayStack[T]) String() string {
	values := make([]string, 0, len(stack.values))
	for _, value := range stack.values {
		values = append(values, uc.StringOf(value))
	}

	var builder strings.Builder

	builder.WriteString("LimitedArrayStack[capacity=")
	builder.WriteString(strconv.Itoa(stack.capacity))
	builder.WriteString(", size=")
	builder.WriteString(strconv.Itoa(len(stack.values)))
	builder.WriteString(", values=[")
	builder.WriteString(strings.Join(values, ", "))
	builder.WriteString(" →]]")

	return builder.String()
}

// CutNilValues is a method of the LimitedArrayStack type. It is used to remove all nil
// values from the stack.
func (stack *LimitedArrayStack[T]) CutNilValues() {
	for i := 0; i < len(stack.values); {
		if gen.IsNil(stack.values[i]) {
			stack.values = append(stack.values[:i], stack.values[i+1:]...)
		} else {
			i++
		}
	}
}

// Slice is a method of the LimitedArrayStack type. It is used to return a slice of the
// elements in the stack.
//
// Returns:
//
//   - []T: A slice of the elements in the stack.
func (stack *LimitedArrayStack[T]) Slice() []T {
	slice := make([]T, len(stack.values))
	copy(slice, stack.values)

	slices.Reverse(slice)

	return slice
}

// Copy is a method of the LimitedArrayStack type. It is used to create a shallow copy
// of the stack.
//
// Returns:
//
//   - itf.Copier: A copy of the stack.
func (stack *LimitedArrayStack[T]) Copy() uc.Copier {
	stackCopy := &LimitedArrayStack[T]{
		values:   make([]T, len(stack.values)),
		capacity: stack.capacity,
	}
	copy(stackCopy.values, stack.values)

	return stackCopy
}
