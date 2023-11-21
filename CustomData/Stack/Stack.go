package Stack

import (
	"fmt"
	"strings"

	gen "github.com/PlayerR9/MyGoLib/Utility/General"
	"github.com/markphelps/optional"
	"golang.org/x/exp/slices"
)

const (
	// Stack Implementation Types

	// LinkedStack Implementation
	LINKED int = iota

	// LinkedStack with a maximum size Implementation
	LINKED_SIZE

	// ArrayStack Implementation
	ARRAY

	// ArrayStack with a maximum size Implementation
	ARRAY_SIZE
)

type Stack[T any] struct {
	data     any
	size     int
	capacity optional.Int
	methods  *stack_methods[T]
}

func (stack *Stack[T]) Initialize(values ...T) {
	stack.size = 0
	stack.data = stack.methods.new()

	for _, element := range values {
		stack.methods.push(stack, element)
	}
}

func (stack *Stack[T]) Push(value T) {
	if stack.capacity.Present() && stack.size >= stack.capacity.MustGet() {
		panic(ErrFullStack{})
	}

	stack.methods.push(stack, value)
	stack.size++
}

func (stack *Stack[T]) Pop() T {
	if stack.size <= 0 {
		panic(ErrEmptyStack{})
	}

	stack.size--

	return stack.methods.pop(stack)
}

func (stack Stack[T]) Peek() T {
	if stack.size == 0 {
		panic(ErrEmptyStack{})
	}

	return stack.methods.peek(stack)
}

func (stack Stack[T]) IsEmpty() bool {
	return stack.size == 0
}

func (stack Stack[T]) Size() int {
	return stack.size
}

func (stack *Stack[T]) ToSlice() []T {
	slice := make([]T, stack.size)

	stack.methods.to_slice(*stack, slice)

	slices.Reverse[[]T, T](slice)

	return slice
}

func (stack *Stack[T]) Clear() {
	stack.size = 0
	stack.data = stack.methods.new()
}

func (stack Stack[T]) IsFull() bool {
	return stack.capacity.Present() && stack.size == stack.capacity.MustGet()
}

func (stack *Stack[T]) String() string {
	var str strings.Builder

	if stack.size != 0 {
		stack.methods.stringer(*stack, &str)
	}

	str.WriteString(StackHead)

	return str.String()
}

type stack_methods[T any] struct {
	new      func() any
	push     func(any, T)
	pop      func(any) T
	peek     func(any) T
	to_slice func(any, []T)
	stringer func(any, *strings.Builder)
}

func NewStack[T any](implementation int, capacity optional.Int) Stack[T] {
	var stack Stack[T]

	switch implementation {
	case LINKED, LINKED_SIZE:
		stack.data = linked_stack[T]{
			top: nil,
		}

		stack.size = 0

		if implementation == LINKED_SIZE {
			if !capacity.Present() {
				panic("Must specify capacity for a linked stack with a maximum size")
			}

			if capacity.MustGet() < 0 {
				panic("Cannot specify a negative capacity for a linked stack")
			}

			stack.capacity = capacity
		}

		if capacity.Present() && implementation != LINKED_SIZE {
			panic("Cannot specify capacity for a linked stack with no maximum size")
		}

		stack.methods = &stack_methods[T]{
			new: func() any {
				return linked_stack[T]{
					top: nil,
				}
			},

			push: func(data any, value T) {
				new_node := stack_node[T]{
					value: gen.DeepCopy(value).(T),
					next:  data.(linked_stack[T]).top,
				}

				tmp := data.(linked_stack[T])
				tmp.top = &new_node
				data = tmp
			},

			pop: func(data any) T {
				value := data.(linked_stack[T]).top.value

				tmp := data.(linked_stack[T])
				tmp.top = tmp.top.next
				data = tmp

				return value
			},

			peek: func(data any) T {
				return gen.DeepCopy(data.(linked_stack[T]).top.value).(T)
			},

			to_slice: func(data any, slice []T) {
				i := 0

				for node := data.(linked_stack[T]).top; node != nil; node = node.next {
					slice[i] = gen.DeepCopy(node.value).(T)
					i++
				}
			},

			stringer: func(data any, str *strings.Builder) {
				node := data.(linked_stack[T]).top

				str.WriteString(fmt.Sprintf("%v", node.value))

				for node.next != nil {
					node = node.next
					str.WriteString(fmt.Sprintf("%v", node.value))
					str.WriteString(StackSep)
				}
			},
		}
	case ARRAY, ARRAY_SIZE:
		stack.data = linked_stack[T]{
			top: nil,
		}

		stack.size = 0

		if capacity.Present() {
			if implementation == ARRAY_SIZE {
				stack.capacity = capacity
			} else {
				panic("Cannot specify capacity for an array stack with no maximum size")
			}
		} else if implementation == ARRAY_SIZE {
			panic("Must specify capacity for an array stack with a maximum size")
		}

		stack.methods = &stack_methods[T]{
			new: func() any {
				return make([]T, 0)
			},

			push: func(data any, value T) {
				tmp := data.([]T)
				tmp = append(tmp, gen.DeepCopy(value).(T))
				data = tmp
			},

			pop: func(data any) T {
				tmp := data.([]T)
				value := tmp[len(tmp)-1]

				tmp = tmp[:len(tmp)-1]
				data = tmp

				return value
			},

			peek: func(data any) T {
				tmp := data.([]T)

				return gen.DeepCopy(tmp[len(tmp)-1]).(T)
			},

			to_slice: func(data any, slice []T) {
				for i, element := range data.([]T) {
					slice[i] = gen.DeepCopy(element).(T)
				}
			},

			stringer: func(data any, str *strings.Builder) {
				tmp := data.([]T)

				str.WriteString(fmt.Sprintf("%v", tmp[0]))

				for _, element := range tmp[1:] {
					str.WriteString(fmt.Sprintf("%v", element))
					str.WriteString(StackSep)
				}
			},
		}
	default:
		panic(fmt.Sprintf("Invalid stack implementation type: %d", implementation))
	}

	return stack
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
