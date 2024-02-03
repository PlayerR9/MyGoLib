package Stack

import (
	"fmt"
	"slices"
	"strings"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

// LimitedArrayStack is a generic type in Go that represents a stack data structure with
// a limited capacity.
// It has a single field, values, which is a slice of type T. This slice stores the
// elements in the stack.
type LimitedArrayStack[T any] struct {
	values []T
}

// NewLimitedArrayStack is a function that creates and returns a new instance of a
// LimitedArrayStack.
// It takes an integer capacity, which represents the maximum number of elements the
// stack can hold, and a variadic parameter of type T, which represents the initial values
// to be stored in the stack.
//
// The function first checks if the provided capacity is negative. If it is, it returns an
// error of type ErrNegativeCapacity.
// It then checks if the number of initial values exceeds the provided capacity. If it does,
// it returns an error of type ErrTooManyValues.
//
// If the provided capacity and initial values are valid, the function creates a new
// LimitedArrayStack, initializes its values field with a slice
// of the same length as the input values and the provided capacity, and then copies the
// input values into the new slice. The new LimitedArrayStack is then returned.
func NewLimitedArrayStack[T any](capacity int, values ...T) (*LimitedArrayStack[T], error) {
	if capacity < 0 {
		return nil, ers.NewErrInvalidParameter(
			"capacity",
			new(ErrNegativeCapacity),
		)
	} else if len(values) > capacity {
		return nil, ers.NewErrInvalidParameter(
			"values", new(ErrTooManyValues),
		)
	}

	slices.Reverse(values)

	stack := &LimitedArrayStack[T]{
		values: make([]T, len(values), capacity),
	}
	copy(stack.values, values)

	return stack, nil
}

// Push is a method of the LimitedArrayStack type. It is used to add an element to the
// end of the stack.
//
// The method takes a parameter, value, of a generic type T, which is the element to be
// added to the stack.
//
// Before adding the element, the method checks if the current length of the values slice
// is equal to the capacity of the stack.
// If it is, it means the stack is full, and the method panics by throwing an ErrFullStack
// error.
//
// If the stack is not full, the method appends the value to the end of the values slice,
// effectively adding the element to the end of the stack.
func (stack *LimitedArrayStack[T]) Push(value T) {
	if cap(stack.values) == len(stack.values) {
		panic(new(ErrFullStack))
	}

	stack.values = append(stack.values, value)
}

func (stack *LimitedArrayStack[T]) Pop() T {
	if len(stack.values) == 0 {
		panic(NewErrEmptyStack(Pop))
	}

	var value T

	value, stack.values = stack.values[len(stack.values)-1], stack.values[:len(stack.values)-1]

	return value
}

func (stack *LimitedArrayStack[T]) Peek() T {
	if len(stack.values) == 0 {
		panic(NewErrEmptyStack(Peek))
	}

	return stack.values[len(stack.values)-1]
}

func (stack *LimitedArrayStack[T]) IsEmpty() bool {
	return len(stack.values) == 0
}

func (stack *LimitedArrayStack[T]) Size() int {
	return len(stack.values)
}

func (stack *LimitedArrayStack[T]) ToSlice() []T {
	slice := make([]T, len(stack.values))

	copy(slice, stack.values)
	slices.Reverse(slice)

	return slice
}

func (stack *LimitedArrayStack[T]) Clear() {
	stack.values = make([]T, 0, cap(stack.values))
}

func (stack *LimitedArrayStack[T]) IsFull() bool {
	return cap(stack.values) == len(stack.values)
}

func (stack *LimitedArrayStack[T]) String() string {
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
