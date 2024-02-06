package ListLike

import (
	"fmt"
	"slices"
	"strings"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
	"github.com/markphelps/optional"
)

// ArrayStack is a generic type that represents a stack data structure with
// or without a limited capacity. It is implemented using an array.
type ArrayStack[T any] struct {
	// values is a slice of type T that stores the elements in the stack.
	values []*T

	// capacity is the maximum number of elements the stack can hold.
	capacity optional.Int
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
func NewArrayStack[T any](values ...*T) *ArrayStack[T] {
	slices.Reverse(values)

	stack := &ArrayStack[T]{
		values: make([]*T, len(values)),
	}
	copy(stack.values, values)

	return stack
}

// WithCapacity is a method of the ArrayStack type. It is used to set the maximum
// number of elements the stack can hold.
//
// Panics with an error of type *ErrCallFailed if the capacity is already set,
// or with an error of type *ErrInvalidParameter if the provided capacity is negative
// or less than the current number of elements in the stack.
//
// Parameters:
//
//   - capacity: An integer that represents the maximum number of elements the stack can hold.
//
// Returns:
//
//   - *ArrayStack[T]: A pointer to the stack with the new capacity set.
func (stack *ArrayStack[T]) WithCapacity(capacity int) *ArrayStack[T] {
	defer ers.PropagatePanic(ers.NewErrCallFailed("WithCapacity", stack.WithCapacity))

	stack.capacity.If(func(cap int) {
		panic(fmt.Errorf("capacity is already set to %d", cap))
	})

	if capacity < 0 {
		panic(ers.NewErrInvalidParameter("capacity").
			WithReason(fmt.Errorf("negative capacity (%d) is not allowed", capacity)),
		)
	} else if len(stack.values) > capacity {
		panic(ers.NewErrInvalidParameter("capacity").WithReason(
			fmt.Errorf("provided capacity (%d) is less than the current number of values (%d)",
				capacity, len(stack.values)),
		))
	}

	stack.capacity = optional.NewInt(capacity)

	newValues := make([]*T, len(stack.values), capacity)
	copy(newValues, stack.values)

	stack.values = newValues

	return stack
}

// Push is a method of the ArrayStack type. It is used to add an element to the
// end of the stack.
//
// Panics with an error of type *ErrCallFailed if the stack is full.
//
// Parameters:
//
//   - value: A pointer to the element to be added to the stack.
func (stack *ArrayStack[T]) Push(value *T) {
	stack.capacity.If(func(cap int) {
		if len(stack.values) <= cap {
			panic(ers.NewErrCallFailed("Push", stack.Push))
		}
	})

	stack.values = append(stack.values, value)
}

// Pop is a method of the ArrayStack type. It is used to remove and return the
// element at the end of the stack.
//
// Panics with an error of type *ErrCallFailed if the stack is empty.
//
// Returns:
//
//   - *T: A pointer to the element that was removed from the stack.
func (stack *ArrayStack[T]) Pop() *T {
	if len(stack.values) == 0 {
		panic(ers.NewErrCallFailed("Pop", stack.Pop).
			WithReason(NewErrEmptyStack(stack)),
		)
	}

	toRemove := stack.values[len(stack.values)-1]
	stack.values[len(stack.values)-1], stack.values = nil, stack.values[:len(stack.values)-1]

	return toRemove
}

// Peek is a method of the ArrayStack type. It is used to return the element at the
// end of the stack without removing it.
//
// Panics with an error of type *ErrCallFailed if the stack is empty.
//
// Returns:
//
//   - *T: A pointer to the element at the end of the stack.
func (stack *ArrayStack[T]) Peek() *T {
	if len(stack.values) == 0 {
		panic(ers.NewErrCallFailed("Peek", stack.Peek).
			WithReason(NewErrEmptyStack(stack)),
		)
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

// Capacity is a method of the ArrayStack type. It is used to return the maximum number
// of elements the stack can hold.
//
// Returns:
//
//   - optional.Int: An optional integer that represents the maximum number of elements
//     the stack can hold.
func (stack *ArrayStack[T]) Capacity() optional.Int {
	return stack.capacity
}

// ToSlice is a method of the ArrayStack type. It is used to return the elements in the
// stack as a slice.
//
// Returns:
//
//   - []*T: A slice of pointers to the elements in the stack.
func (stack *ArrayStack[T]) ToSlice() []*T {
	slice := make([]*T, len(stack.values))

	copy(slice, stack.values)
	slices.Reverse(slice)

	return slice
}

// Clear is a method of the ArrayStack type. It is used to remove all elements from the
// stack, making it empty.
func (stack *ArrayStack[T]) Clear() {
	for i := range stack.values {
		stack.values[i] = nil
	}

	if stack.capacity.Present() {
		stack.values = make([]*T, 0, stack.capacity.MustGet())
	} else {
		stack.values = make([]*T, 0)
	}
}

// IsFull is a method of the ArrayStack type. It is used to check if the stack is full,
// i.e., if it has reached its maximum capacity.
//
// Returns:
//
//   - isFull: A boolean value that is true if the stack is full, and false otherwise.
func (stack *ArrayStack[T]) IsFull() (isFull bool) {
	stack.capacity.If(func(cap int) {
		isFull = len(stack.values) == cap
	})

	return
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

	stack.capacity.If(func(cap int) {
		fmt.Fprintf(&builder, "capacity=%d, ", cap)
	})

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

// CutNilValues is a method of the ArrayStack type. It is used to remove all nil
// values from the stack.
func (stack *ArrayStack[T]) CutNilValues() {
	for i := 0; i < len(stack.values); {
		if stack.values[i] == nil {
			stack.values = append(stack.values[:i], stack.values[i+1:]...)
		} else {
			i++
		}
	}
}
