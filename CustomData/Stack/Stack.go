package Stack

import (
	"fmt"
	"slices"
	"strings"

	gen "github.com/PlayerR9/MyGoLib/Utility/General"
	"github.com/markphelps/optional"
)

const (
	// Stack Implementation Types

	// linked_stack with no maximum size Implementation
	LINKED int = iota

	// ArrayStack with no maximum size Implementation
	ARRAY
)

type Stack[T any] struct {
	data           any
	implementation int
	size           int
	capacity       optional.Int
}

func NewStack[T any](implementation int, values ...T) Stack[T] {
	var stack Stack[T]

	switch implementation {
	case LINKED:
		stack.data = linked_stack[T]{
			top: nil,
		}
	case ARRAY:
		stack.data = make([]T, 0)
	default:
		panic("Invalid stack implementation type")
	}

	stack.implementation = implementation

	stack.size = 0
	stack.capacity = optional.Int{}

	for _, element := range values {
		stack.Push(element)
	}

	return stack
}

func NewLimitedStack[T any](implementation, capacity int, values ...T) Stack[T] {
	if capacity < 0 {
		panic("Cannot specify a negative capacity for a stack")
	}

	if len(values) > capacity {
		panic("Cannot specify more values than the capacity of a stack")
	}

	var stack Stack[T]

	switch implementation {
	case LINKED:
		stack.data = linked_stack[T]{
			top: nil,
		}
	case ARRAY:
		stack.data = make([]T, 0, capacity)
	default:
		panic("Invalid stack implementation type")
	}

	stack.implementation = implementation

	stack.capacity = optional.NewInt(capacity)
	stack.size = 0

	for _, element := range values {
		stack.Push(element)
	}

	return stack
}

func (stack *Stack[T]) Push(value T) {
	if stack.capacity.Present() && stack.size >= stack.capacity.MustGet() {
		panic(ErrFullStack{})
	}

	switch stack.implementation {
	case LINKED:
		new_node := stack_node[T]{
			value: gen.DeepCopy(value).(T),
			next:  stack.data.(linked_stack[T]).top,
		}

		tmp := stack.data.(linked_stack[T])
		tmp.top = &new_node
		stack.data = tmp
	case ARRAY:
		tmp := stack.data.([]T)
		tmp = append(tmp, gen.DeepCopy(value).(T))
		stack.data = tmp
	}

	stack.size++
}

func (stack *Stack[T]) Pop() T {
	if stack.size <= 0 {
		panic(ErrEmptyStack{})
	}

	var value T

	switch stack.implementation {
	case LINKED:
		value = stack.data.(linked_stack[T]).top.value

		tmp := stack.data.(linked_stack[T])
		tmp.top = tmp.top.next
		stack.data = tmp
	case ARRAY:
		tmp := stack.data.([]T)
		value = tmp[len(tmp)-1]

		tmp = tmp[:len(tmp)-1]
		stack.data = tmp
	}

	stack.size--

	return value
}

func (stack Stack[T]) Peek() T {
	if stack.size == 0 {
		panic(ErrEmptyStack{})
	}

	var value T

	switch stack.implementation {
	case LINKED:
		value = gen.DeepCopy(stack.data.(linked_stack[T]).top.value).(T)
	case ARRAY:
		tmp := stack.data.([]T)

		value = gen.DeepCopy(tmp[len(tmp)-1]).(T)
	}

	return value
}

func (stack Stack[T]) IsEmpty() bool {
	return stack.size == 0
}

func (stack Stack[T]) Size() int {
	return stack.size
}

func (stack *Stack[T]) ToSlice() []T {
	slice := make([]T, stack.size)

	switch stack.implementation {
	case LINKED:
		i := 0

		for node := stack.data.(linked_stack[T]).top; node != nil; node = node.next {
			slice[i] = gen.DeepCopy(node.value).(T)
			i++
		}
	case ARRAY:
		for i, element := range stack.data.([]T) {
			slice[i] = gen.DeepCopy(element).(T)
		}
	}

	slices.Reverse(slice)

	return slice
}

func (stack *Stack[T]) Clear() {
	stack.size = 0

	switch stack.implementation {
	case LINKED:
		stack.data = linked_stack[T]{
			top: nil,
		}
	case ARRAY:
		if stack.capacity.Present() {
			stack.data = make([]T, 0, stack.capacity.MustGet())
		} else {
			stack.data = make([]T, 0)
		}
	}
}

func (stack Stack[T]) IsFull() bool {
	return stack.capacity.Present() && stack.size == stack.capacity.MustGet()
}

func (stack *Stack[T]) String() string {
	if stack.size == 0 {
		return StackHead
	}

	var str strings.Builder

	switch stack.implementation {
	case LINKED:
		node := stack.data.(linked_stack[T]).top

		str.WriteString(fmt.Sprintf("%v", node.value))

		for node.next != nil {
			node = node.next
			str.WriteString(fmt.Sprintf("%v", node.value))
			str.WriteString(StackSep)
		}
	case ARRAY:
		tmp := stack.data.([]T)

		str.WriteString(fmt.Sprintf("%v", tmp[0]))

		for _, element := range tmp[1:] {
			str.WriteString(fmt.Sprintf("%v", element))
			str.WriteString(StackSep)
		}
	}

	str.WriteString(StackHead)

	return str.String()
}

// Stack Implementation Types

type stack_node[T any] struct {
	value T
	next  *stack_node[T]
}

type linked_stack[T any] struct {
	top *stack_node[T]
}

// Stack Errors

type ErrEmptyStack struct{}

func (e ErrEmptyStack) Error() string {
	return "Empty stack"
}

type ErrFullStack struct{}

func (e ErrFullStack) Error() string {
	return "Full stack"
}

var (
	StackHead string = " | →"
	StackSep  string = " | "
)
