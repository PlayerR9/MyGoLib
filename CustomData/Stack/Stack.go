package Stack

import (
	"fmt"
	"slices"
	"strings"

	gen "github.com/PlayerR9/MyGoLib/Utility/General"
	"github.com/markphelps/optional"
)

type StackImplementationType int

const (
	LinkedStack StackImplementationType = iota
	ArrayStack
)

func (implementation StackImplementationType) String() string {
	return [...]string{
		"Stack implemented with a linked list",
		"Stack implemented with an array",
	}[implementation]
}

type Stack[T any] struct {
	data           any
	implementation StackImplementationType
	size           int
	capacity       optional.Int
}

func NewStack[T any](implementation StackImplementationType, values ...T) Stack[T] {
	var stack Stack[T]

	switch implementation {
	case LinkedStack:
		stack.data = linked_stack[T]{
			top: nil,
		}
	case ArrayStack:
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

func NewLimitedStack[T any](implementation StackImplementationType, capacity int, values ...T) Stack[T] {
	if capacity < 0 {
		panic("Cannot specify a negative capacity for a stack")
	}

	if len(values) > capacity {
		panic("Cannot specify more values than the capacity of a stack")
	}

	var stack Stack[T]

	switch implementation {
	case LinkedStack:
		stack.data = linked_stack[T]{
			top: nil,
		}
	case ArrayStack:
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
	case LinkedStack:
		new_node := stack_node[T]{
			value: gen.DeepCopy(value).(T),
			next:  stack.data.(linked_stack[T]).top,
		}

		tmp := stack.data.(linked_stack[T])
		tmp.top = &new_node
		stack.data = tmp
	case ArrayStack:
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
	case LinkedStack:
		value = stack.data.(linked_stack[T]).top.value

		tmp := stack.data.(linked_stack[T])
		tmp.top = tmp.top.next
		stack.data = tmp
	case ArrayStack:
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
	case LinkedStack:
		value = gen.DeepCopy(stack.data.(linked_stack[T]).top.value).(T)
	case ArrayStack:
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
	case LinkedStack:
		i := 0

		for node := stack.data.(linked_stack[T]).top; node != nil; node = node.next {
			slice[i] = gen.DeepCopy(node.value).(T)
			i++
		}
	case ArrayStack:
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
	case LinkedStack:
		stack.data = linked_stack[T]{
			top: nil,
		}
	case ArrayStack:
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
	case LinkedStack:
		node := stack.data.(linked_stack[T]).top

		str.WriteString(fmt.Sprintf("%v", node.value))

		for node.next != nil {
			node = node.next
			str.WriteString(fmt.Sprintf("%v", node.value))
			str.WriteString(StackSep)
		}
	case ArrayStack:
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

var (
	StackHead string = " | →"
	StackSep  string = " | "
)
