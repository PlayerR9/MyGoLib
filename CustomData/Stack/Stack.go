package Stack

import (
	"log"
	"os"
)

// ERRORS

// ErrEmptyStack is an error that is used when the stack is empty.
type ErrEmptyStack struct{}

func (e ErrEmptyStack) Error() string {
	return "Empty stack"
}

// ErrFullStack is an error that is used when the stack is full.
type ErrFullStack struct{}

func (e ErrFullStack) Error() string {
	return "Full stack"
}

// GLOBAL VARIABLES

var (
	// DebugMode is a boolean that is used to enable or disable debug mode. When debug mode is enabled, the package will print debug messages.
	// **Note:** Debug mode is disabled by default.
	DebugMode bool = false

	debugger *log.Logger = log.New(os.Stdout, "[Stack] ", log.LstdFlags) // debugger is a logger that is used to print debug messages.
)

// A Stack is a data structure that is used to store values. It is a LIFO (Last In First Out) data structure, meaning that the last value inserted into the stack is the
// first one to be removed.
type Stack[T any] interface {
	// Push is a function that is used insert a value into the stack at the top.
	//
	// Parameters:
	//   - value: The value to insert into the stack.
	//
	// Returns:
	//   - Stack[T]: The stack.
	Push(value T) Stack[T]

	// Pop is a function that is used to remove a value from the stack at the top. Panics if the stack is empty.
	//
	// Returns:
	//   - T: The value removed from the stack.
	//   - Stack[T]: The stack.
	Pop() (T, Stack[T])

	// Peek is a function that is used to peek at the top value of the stack, without removing it. Panics if the stack is empty.
	//
	// Returns:
	//   - T: The top value of the stack.
	Peek() T

	// IsEmpty is a function that is used to check if the stack is empty.
	//
	// Returns:
	//   - bool: True if the stack is empty, false otherwise.
	IsEmpty() bool

	// Size is a function that is used to get the size of the stack.
	//
	// Returns:
	//   - int: The size of the stack.
	Size() int

	// ToSlice is a function that is used to get the stack as a slice. The first element of the slice is the top of the stack.
	//
	// Returns:
	//   - []T: The stack as a slice.
	//
	// Information:
	//   - The first element of the slice is the top of the stack.
	ToSlice() []T

	// ToString is a function that is used to get the stack as a string. The first element of the string is the bottom of the stack.
	//
	// Parameters:
	//   - f: The function to convert the values to strings.
	//   - sep: The separator to use between values.
	//
	// Returns:
	//   - string: The stack as a string.
	//
	// Information:
	//   - The first element of the string is the top of the stack.
	ToString(f func(T) string, sep string) string

	// Copy is a function that is used to copy the stack.
	//
	// Parameters:
	//   - f: The function to copy the values.
	//
	// Returns:
	//   - Stack[T]: The copy of the stack.
	Copy() Stack[T]
}

// node is a node in a linked list.
type node[T any] struct {
	value T        // The value of the node.
	next  *node[T] // The next node in the linked list.
}

// LinkedStack is a stack implemented using a linked list.
type LinkedStack[T any] struct {
	top           *node[T]  // The top of the stack.
	size          int       // The size of the stack.
	copy_function func(T) T // The function to copy the values.
}

// NewLinkedStack creates a new LinkedStack.
//
// Parameters:
//   - copy_function: The function to copy the values.
//   - elements: The elements to initialize the stack with.
//
// Returns:
//   - LinkedStack: A new LinkedStack.
func NewLinkedStack[T any](copy_function func(T) T, elements ...T) LinkedStack[T] {
	// Create the stack.
	stack := LinkedStack[T]{
		top:           nil,
		size:          0,
		copy_function: copy_function,
	}

	// Push the elements onto the stack.
	for _, element := range elements {
		stack = stack.Push(element).(LinkedStack[T])
	}

	return stack
}

func (stack LinkedStack[T]) Push(value T) Stack[T] {
	// Create a new node.
	new_node := node[T]{
		value: stack.copy_function(value),
		next:  stack.top,
	}

	stack.top = &new_node // Set the new node as the top of the stack.

	stack.size++ // Increment the size of the stack.

	return stack
}

func (stack LinkedStack[T]) Pop() (T, Stack[T]) {
	if stack.top == nil {
		// If the stack is empty, panic.
		if DebugMode {
			debugger.Panic(ErrEmptyStack{}.Error())
		} else {
			panic(ErrEmptyStack{}.Error())
		}
	}

	value := stack.top.value   // Get the value of the top of the stack.
	stack.top = stack.top.next // Set the next node as the top of the stack.

	stack.size-- // Decrement the size of the stack.

	return stack.copy_function(value), stack
}

func (stack LinkedStack[T]) Peek() T {
	if stack.top == nil {
		// If the stack is empty, panic.
		if DebugMode {
			debugger.Panic(ErrEmptyStack{}.Error())
		} else {
			panic(ErrEmptyStack{}.Error())
		}
	}

	return stack.copy_function(stack.top.value) // Get the value of the top of the stack.
}

func (stack LinkedStack[T]) IsEmpty() bool {
	return stack.top == nil // If the top of the stack is nil, the stack is empty.
}

func (stack LinkedStack[T]) Size() int {
	return stack.size
}

func (stack LinkedStack[T]) ToSlice() []T {
	slice := make([]T, stack.size) // Create a slice with the size of the stack.

	// Add the values to the slice.
	i := 0
	for node := stack.top; node != nil; node = node.next {
		slice[i] = stack.copy_function(node.value)
		i++
	}

	// Reverse the slice.
	for i := 0; i < len(slice)/2; i++ {
		slice[i], slice[len(slice)-i-1] = slice[len(slice)-i-1], slice[i]
	}

	return slice
}

func (stack LinkedStack[T]) ToString(f func(T) string, sep string) string {
	str := ""

	for node := stack.top; node != nil; node = node.next {
		str += f(node.value)

		if node.next != nil {
			str += sep
		}
	}

	return str + " | →"
}

func (stack LinkedStack[T]) Copy() Stack[T] {
	new_stack := LinkedStack[T]{
		top:           nil,
		size:          stack.size,
		copy_function: stack.copy_function,
	}

	if stack.top == nil {
		return new_stack
	}

	new_stack.top = &node[T]{
		value: new_stack.copy_function(stack.top.value),
		next:  nil,
	}

	for n := stack.top.next; n != nil; n = n.next {
		// Create a new node.
		new_node := node[T]{
			value: new_stack.copy_function(n.value),
			next:  nil,
		}

		// Find the last node.
		last_node := new_stack.top

		for last_node.next != nil {
			last_node = last_node.next
		}

		// Add the new node to the end.
		last_node.next = &new_node
	}

	return new_stack
}

// ArrayStack is a stack implemented using an array. It panics if the stack is full.
type ArrayStack[T any] struct {
	elements      []T       // The elements of the stack.
	hasMaxSize    bool      // Whether or not the stack has a maximum size.
	copy_function func(T) T // The function to copy the values.
}

// NewArrayStack creates a new ArrayStack without a maximum size.
//
// Parameters:
//   - copy_function: The function to copy the values.
//   - elements: The elements to initialize the stack with.
//
// Returns:
//   - ArrayStack: A new ArrayStack.
func NewArrayStack[T any](copy_function func(T) T, elements ...T) ArrayStack[T] {
	// Create the stack.
	obj := ArrayStack[T]{
		elements:      make([]T, len(elements)),
		hasMaxSize:    false,
		copy_function: copy_function,
	}

	// Push the elements onto the stack.
	for i, element := range elements {
		obj.elements[i] = obj.copy_function(element)
	}

	return obj
}

// NewArrayStackWithMaxSize creates a new ArrayStack with a maximum size.
//
// Parameters:
//   - copy_function: The function to copy the values.
//   - max_size: The maximum size of the stack.
//   - elements: The elements to initialize the stack with.
//
// Returns:
//   - ArrayStack: A new ArrayStack.
func NewArrayStackWithMaxSize[T any](copy_function func(T) T, max_size int, elements ...T) ArrayStack[T] {
	// Create the stack.
	obj := ArrayStack[T]{
		elements:      make([]T, len(elements), max_size),
		hasMaxSize:    true,
		copy_function: copy_function,
	}

	// Push the elements onto the stack.
	for i, element := range elements {
		obj.elements[i] = obj.copy_function(element)
	}

	return obj
}

func (stack ArrayStack[T]) Push(value T) Stack[T] {
	if stack.hasMaxSize && len(stack.elements) == cap(stack.elements) {
		// If the stack is full, panic.
		if DebugMode {
			debugger.Panic(ErrFullStack{}.Error())
		} else {
			panic(ErrFullStack{}.Error())
		}
	}

	// Add the value to the stack.
	stack.elements = append(stack.elements, stack.copy_function(value))

	return stack
}

func (stack ArrayStack[T]) Pop() (T, Stack[T]) {
	if len(stack.elements) == 0 {
		// If the stack is empty, panic.
		if DebugMode {
			debugger.Panic(ErrEmptyStack{}.Error())
		} else {
			panic(ErrEmptyStack{}.Error())
		}
	}

	value := stack.elements[len(stack.elements)-1] // Get the value of the top of the stack.

	stack.elements = stack.elements[:len(stack.elements)-1] // Remove the value from the stack.

	return stack.copy_function(value), stack
}

func (stack ArrayStack[T]) Peek() T {
	if len(stack.elements) == 0 {
		// If the stack is empty, panic.
		if DebugMode {
			debugger.Panic(ErrEmptyStack{}.Error())
		} else {
			panic(ErrEmptyStack{}.Error())
		}
	}

	return stack.copy_function(stack.elements[len(stack.elements)-1]) // Get the value of the top of the stack.
}

func (stack ArrayStack[T]) IsEmpty() bool {
	return len(stack.elements) == 0 // If the stack is empty, the length of the elements is 0.
}

func (stack ArrayStack[T]) Size() int {
	return len(stack.elements)
}

func (stack ArrayStack[T]) ToSlice() []T {
	slice := make([]T, len(stack.elements)) // Create a slice with the size of the stack.

	// Add the values to the slice.
	for i, element := range stack.elements {
		slice[i] = stack.copy_function(element)
	}

	// Reverse the slice.
	for i := 0; i < len(slice)/2; i++ {
		slice[i], slice[len(slice)-i-1] = slice[len(slice)-i-1], slice[i]
	}

	return slice
}

func (stack ArrayStack[T]) ToString(f func(T) string, sep string) string {
	var str string

	if len(stack.elements) != 0 {
		str := f(stack.elements[0]) // Get the value of the top of the stack.

		// Add the values to the string.
		for _, element := range stack.elements[1:] {
			str += f(element) + sep
		}
	} else {
		str = ""
	}

	return str + " | →"
}

func (stack ArrayStack[T]) Copy() Stack[T] {
	// Create the stack.
	var obj ArrayStack[T]

	if stack.hasMaxSize {
		obj = ArrayStack[T]{
			elements:      make([]T, len(stack.elements), cap(stack.elements)),
			hasMaxSize:    true,
			copy_function: stack.copy_function,
		}
	} else {
		obj = ArrayStack[T]{
			elements:      make([]T, len(stack.elements)),
			hasMaxSize:    false,
			copy_function: stack.copy_function,
		}
	}

	// Copy the values.
	for i, element := range stack.elements {
		obj.elements[i] = stack.copy_function(element)
	}

	return obj
}
