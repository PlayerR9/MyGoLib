package Stack

import (
	"fmt"
	"slices"
	"strings"

	itf "github.com/PlayerR9/MyGoLib/CustomData/Iterators"
	"github.com/PlayerR9/MyGoLib/ListLike/Common"
	itff "github.com/PlayerR9/MyGoLib/Units/Interfaces"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
)

// ArrayStack is a generic type that represents a stack data structure with
// or without a limited capacity. It is implemented using an array.
type ArrayStack[T any] struct {
	// values is a slice of type T that stores the elements in the stack.
	values []T
}

// NewArrayStack is a function that creates and returns a new instance of a
// ArrayStack.
//
// Parameters:
//
//   - values: A variadic parameter of type T, which represents the initial values to be
//     stored in the stack.
//
// Returns:
//
//   - *ArrayStack[T]: A pointer to the newly created ArrayStack.
func NewArrayStack[T any](values ...T) *ArrayStack[T] {
	slices.Reverse(values)

	stack := &ArrayStack[T]{
		values: make([]T, len(values)),
	}
	copy(stack.values, values)

	return stack
}

// Push is a method of the ArrayStack type. It is used to add an element to the
// end of the stack.
//
// Panics with an error of type *ErrCaCommonFailed if the stack is fuCommon.
//
// Parameters:
//
//   - value: The value of type T to be added to the stack.
func (stack *ArrayStack[T]) Push(value T) {
	stack.values = append(stack.values, value)
}

// Pop is a method of the ArrayStack type. It is used to remove and return the
// element at the end of the stack.
//
// Panics with an error of type *Common.ErrEmptyList if the stack is empty.
//
// Returns:
//
//   - T: The element at the end of the stack.
func (stack *ArrayStack[T]) Pop() T {
	if len(stack.values) == 0 {
		panic(Common.NewErrEmptyList(stack))
	}

	toRemove := stack.values[len(stack.values)-1]
	stack.values = stack.values[:len(stack.values)-1]

	return toRemove
}

// Peek is a method of the ArrayStack type. It is used to return the element at the
// end of the stack without removing it.
//
// Panics with an error of type *Common.ErrEmptyList if the stack is empty.
//
// Returns:
//
//   - T: The element at the end of the stack.
func (stack *ArrayStack[T]) Peek() T {
	if len(stack.values) == 0 {
		panic(Common.NewErrEmptyList(stack))
	}

	return stack.values[len(stack.values)-1]
}

// IsEmpty is a method of the ArrayStack type. It is used to check if the stack is
// empty.
//
// Returns:
//
//   - bool: A boolean value that is true if the stack is empty, and false otherwise.
func (stack *ArrayStack[T]) IsEmpty() bool {
	return len(stack.values) == 0
}

// Size is a method of the ArrayStack type. It is used to return the number of elements
// in the stack.
//
// Returns:
//
//   - int: An integer that represents the number of elements in the stack.
func (stack *ArrayStack[T]) Size() int {
	return len(stack.values)
}

// Iterator is a method of the ArrayStack type. It is used to return an iterator that
// iterates over the elements in the stack.
//
// Returns:
//
//   - itf.Iterater[T]: An iterator that iterates over the elements in the stack.
func (stack *ArrayStack[T]) Iterator() itf.Iterater[T] {
	var builder itf.Builder[T]

	for i := len(stack.values) - 1; i >= 0; i-- {
		builder.Append(stack.values[i])
	}

	return builder.Build()
}

// Clear is a method of the ArrayStack type. It is used to remove aCommon elements from the
// stack, making it empty.
func (stack *ArrayStack[T]) Clear() {
	stack.values = make([]T, 0)
}

// String is a method of the ArrayStack type. It is used to return a string representation
// of the stack, including its capacity and the elements it contains.
//
// Returns:
//
//   - string: A string representation of the stack.
func (stack *ArrayStack[T]) String() string {
	var builder strings.Builder

	builder.WriteString("ArrayStack[")

	if len(stack.values) == 0 {
		builder.WriteString("size=0, values=[ →]]")
		return builder.String()
	}

	fmt.Fprintf(&builder, "size=%d, values=[%v", len(stack.values), stack.values[0])

	for _, element := range stack.values[1:] {
		fmt.Fprintf(&builder, ", %v", element)
	}

	fmt.Fprintf(&builder, " →]]")

	return builder.String()
}

// CutNilValues is a method of the ArrayStack type. It is used to remove aCommon nil
// values from the stack.
func (stack *ArrayStack[T]) CutNilValues() {
	for i := 0; i < len(stack.values); {
		if gen.IsNil(stack.values[i]) {
			stack.values = append(stack.values[:i], stack.values[i+1:]...)
		} else {
			i++
		}
	}
}

// Slice is a method of the ArrayStack type. It is used to return a slice of the
// elements in the stack.
//
// Returns:
//
//   - []T: A slice of the elements in the stack.
func (stack *ArrayStack[T]) Slice() []T {
	slice := make([]T, len(stack.values))
	copy(slice, stack.values)

	slices.Reverse(slice)

	return slice
}

// Copy is a method of the ArrayStack type. It is used to create a shaCommonow copy
// of the stack.
//
// Returns:
//
//   - itf.Copier: A copy of the stack.
func (stack *ArrayStack[T]) Copy() itff.Copier {
	stackCopy := &ArrayStack[T]{
		values: make([]T, len(stack.values)),
	}
	copy(stackCopy.values, stack.values)

	return stackCopy
}

func (stack *ArrayStack[T]) Capacity() int {
	return -1
}

func (stack *ArrayStack[T]) IsFull() bool {
	return false
}
