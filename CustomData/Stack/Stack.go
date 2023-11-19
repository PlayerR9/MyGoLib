package Stack

import (
	"fmt"
)

type Item interface {
	Copy() Item
}

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

type Stack[T Item] interface {
	Push(value T) Stack[T]
	Pop() (T, Stack[T])
	Peek() T
	IsEmpty() bool
	Size() int
	ToSlice() []T
	Copy() Stack[T]
}

type node[T Item] struct {
	value T
	next  *node[T]
}

type LinkedStack[T Item] struct {
	top  *node[T]
	size int
}

func NewLinkedStack[T Item](elements ...T) LinkedStack[T] {
	stack := LinkedStack[T]{
		top:  nil,
		size: 0,
	}

	for _, element := range elements {
		stack = stack.Push(element).(LinkedStack[T])
	}

	return stack
}

func (stack LinkedStack[T]) Push(value T) Stack[T] {
	new_node := node[T]{
		value: value.Copy().(T),
		next:  stack.top,
	}

	stack.top = &new_node

	stack.size++

	return stack
}

func (stack LinkedStack[T]) Pop() (T, Stack[T]) {
	if stack.top == nil {
		panic(ErrEmptyStack{})
	}

	value := stack.top.value
	stack.top = stack.top.next

	stack.size--

	return value.Copy().(T), stack
}

func (stack LinkedStack[T]) Peek() T {
	if stack.top == nil {
		panic(ErrEmptyStack{})
	}

	return stack.top.value.Copy().(T)
}

func (stack LinkedStack[T]) IsEmpty() bool {
	return stack.top == nil
}

func (stack LinkedStack[T]) Size() int {
	return stack.size
}

func (stack LinkedStack[T]) ToSlice() []T {
	slice := make([]T, stack.size)
	i := 0

	for node := stack.top; node != nil; node = node.next {
		slice[i] = node.value.Copy().(T)
		i++
	}

	for i := 0; i < len(slice)/2; i++ {
		slice[i], slice[len(slice)-i-1] = slice[len(slice)-i-1], slice[i]
	}

	return slice
}

func (stack LinkedStack[T]) String() string {
	if stack.top == nil {
		return StackHead
	}

	var str string

	node := stack.top

	str += fmt.Sprintf("%v", node.value)

	for node.next != nil {
		node = node.next
		str = fmt.Sprintf("%v%s", node.value, StackSep) + str
	}

	return fmt.Sprintf("%s%s", str, StackHead)
}

func (stack LinkedStack[T]) Copy() Stack[T] {
	new_stack := LinkedStack[T]{
		top:  nil,
		size: stack.size,
	}

	if stack.top == nil {
		return new_stack
	}

	new_stack.top = &node[T]{
		value: stack.top.value.Copy().(T),
		next:  nil,
	}

	for n := stack.top.next; n != nil; n = n.next {
		new_node := node[T]{
			value: n.value.Copy().(T),
			next:  nil,
		}

		last_node := new_stack.top

		for last_node.next != nil {
			last_node = last_node.next
		}

		last_node.next = &new_node
	}

	return new_stack
}

type ArrayStack[T Item] struct {
	elements   []T
	hasMaxSize bool
}

func NewArrayStack[T Item](elements ...T) ArrayStack[T] {
	obj := ArrayStack[T]{
		elements:   make([]T, len(elements)),
		hasMaxSize: false,
	}

	for i, element := range elements {
		obj.elements[i] = element.Copy().(T)
	}

	return obj
}

func NewArrayStackWithMaxSize[T Item](max_size int, elements ...T) ArrayStack[T] {
	obj := ArrayStack[T]{
		elements:   make([]T, len(elements), max_size),
		hasMaxSize: true,
	}

	for i, element := range elements {
		obj.elements[i] = element.Copy().(T)
	}

	return obj
}

func (stack ArrayStack[T]) Push(value T) Stack[T] {
	if stack.hasMaxSize && len(stack.elements) == cap(stack.elements) {
		panic(ErrFullStack{})
	}

	stack.elements = append(stack.elements, value.Copy().(T))

	return stack
}

func (stack ArrayStack[T]) Pop() (T, Stack[T]) {
	if len(stack.elements) == 0 {
		panic(ErrEmptyStack{})
	}

	value := stack.elements[len(stack.elements)-1]

	stack.elements = stack.elements[:len(stack.elements)-1]

	return value.Copy().(T), stack
}

func (stack ArrayStack[T]) Peek() T {
	if len(stack.elements) == 0 {
		panic(ErrEmptyStack{})
	}

	return stack.elements[len(stack.elements)-1].Copy().(T)
}

func (stack ArrayStack[T]) IsEmpty() bool {
	return len(stack.elements) == 0
}

func (stack ArrayStack[T]) Size() int {
	return len(stack.elements)
}

func (stack ArrayStack[T]) ToSlice() []T {
	slice := make([]T, len(stack.elements))
	for i, element := range stack.elements {
		slice[i] = element.Copy().(T)
	}
	for i := 0; i < len(slice)/2; i++ {
		slice[i], slice[len(slice)-i-1] = slice[len(slice)-i-1], slice[i]
	}

	return slice
}

func (stack ArrayStack[T]) String() string {
	if len(stack.elements) == 0 {
		return StackHead
	}

	var str string

	str += fmt.Sprintf("%v", stack.elements[0])

	for _, element := range stack.elements[1:] {
		str += fmt.Sprintf("%v%s", element, StackSep)
	}

	return fmt.Sprintf("%s%s", str, StackHead)
}

func (stack ArrayStack[T]) Copy() Stack[T] {
	var obj ArrayStack[T]

	if stack.hasMaxSize {
		obj = ArrayStack[T]{
			elements:   make([]T, len(stack.elements), cap(stack.elements)),
			hasMaxSize: true,
		}
	} else {
		obj = ArrayStack[T]{
			elements:   make([]T, len(stack.elements)),
			hasMaxSize: false,
		}
	}
	for i, element := range stack.elements {
		obj.elements[i] = element.Copy().(T)
	}

	return obj
}
