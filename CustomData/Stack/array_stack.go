package Stack

import (
	"fmt"
	"slices"
	"strings"
)

// ArrayStack is a generic type in Go that represents a stack data structure implemented
// using an array.
// It has a single field, values, which is a slice of type T. This slice stores the
// elements in the stack.
type ArrayStack[T any] struct {
	values []T
}

// NewArrayStack is a function that creates and returns a new instance of an ArrayStack.
// It takes a variadic parameter of type T, which represents the initial values to be
// stored in the stack.
// The function creates a new ArrayStack, initializes its values field with a slice of
// the same length as the input values, and then copies the input values into the new
// slice. The new ArrayStack is then returned.
func NewArrayStack[T any](values ...T) *ArrayStack[T] {
	stack := &ArrayStack[T]{
		values: make([]T, len(values)),
	}

	slices.Reverse(values)

	copy(stack.values, values)

	return stack
}

func (stack *ArrayStack[T]) Push(value T) {
	stack.values = append(stack.values, value)
}

func (stack *ArrayStack[T]) Pop() T {
	if len(stack.values) == 0 {
		panic(NewErrEmptyStack(Pop))
	}

	var value T

	value, stack.values = stack.values[len(stack.values)-1], stack.values[:len(stack.values)-1]

	return value
}

func (stack *ArrayStack[T]) Peek() T {
	if len(stack.values) == 0 {
		panic(NewErrEmptyStack(Peek))
	}

	return stack.values[len(stack.values)-1]
}

func (stack *ArrayStack[T]) IsEmpty() bool {
	return len(stack.values) == 0
}

func (stack *ArrayStack[T]) Size() int {
	return len(stack.values)
}

func (stack *ArrayStack[T]) ToSlice() []T {
	slice := make([]T, len(stack.values))

	copy(slice, stack.values)
	slices.Reverse(slice)

	return slice
}

func (stack *ArrayStack[T]) Clear() {
	stack.values = make([]T, 0)
}

// IsFull is a method of the ArrayStack type. It checks if the stack is full.
//
// In this implementation, the method always returns false. This is because an
// ArrayStack, implemented with a slice, can dynamically grow and shrink in size
// as elements are added or removed. Therefore, it is never considered full,
// and elements can always be added to it.
func (stack *ArrayStack[T]) IsFull() bool {
	return false
}

func (stack *ArrayStack[T]) String() string {
	if len(stack.values) == 0 {
		return StackHead
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "%s%v", StackHead, stack.values[0])

	for _, element := range stack.values[1:] {
		fmt.Fprintf(&builder, "%s%v", StackSep, element)
	}

	return builder.String()
}
