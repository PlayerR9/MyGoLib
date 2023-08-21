package Stack

type ErrEmptyStack struct{}

func (e ErrEmptyStack) Error() string {
	return "Empty stack"
}

// A stack is a data structure that is used to store values. It is a LIFO (Last In First Out) data structure, meaning that the last value inserted into the stack is the
// first one to be removed.
type Stack[T any] interface {
	// push is a function that is used insert a value into the stack at the top.
	//
	// Parameters:
	//   - value: The value to insert into the stack.
	//
	// Returns:
	//   - Stack[T]: The stack.
	Push(value T) Stack[T]

	// pop is a function that is used to remove a value from the stack at the top. Panics if the stack is empty.
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

	// ToSlice is a function that is used to get the stack as a slice.
	//
	// Returns:
	//   - []T: The stack as a slice.
	//
	// Information:
	//   - The first element of the slice is the top of the stack.
	ToSlice() []T

	// ToString is a function that is used to get the stack as a string.
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
	Copy(f func(T) T) Stack[T]
}

// node is a node in a linked list.
type node[T any] struct {
	value T        // The value of the node.
	next  *node[T] // The next node in the linked list.
}

// LinkedStack is a stack implemented using a linked list.
type LinkedStack[T any] struct {
	top  *node[T] // The top of the stack.
	size int      // The size of the stack.
}

// NewLinkedStack creates a new LinkedStack.
//
// Parameters:
//   - elements: The elements to initialize the stack with.
//
// Returns:
//   - LinkedStack: A new LinkedStack.
func NewLinkedStack[T any](elements ...T) LinkedStack[T] {
	stack := LinkedStack[T]{
		top:  nil,
		size: 0,
	}

	for _, element := range elements {
		stack = stack.Push(element).(LinkedStack[T])
	}

	return stack
}

// Push is a function that is used insert a value into the stack at the top.
//
// Parameters:
//   - value: The value to insert into the stack.
func (stack LinkedStack[T]) Push(value T) Stack[T] {
	new_node := node[T]{
		value: value,
		next:  stack.top,
	}

	stack.top = &new_node

	stack.size++

	return stack
}

// Pop is a function that is used to remove a value from the stack at the top. Panics if the stack is empty.
//
// Returns:
//   - T: The value removed from the stack.
func (stack LinkedStack[T]) Pop() (T, Stack[T]) {
	if stack.top == nil {
		panic(ErrEmptyStack{}.Error())
	}

	value := stack.top.value
	stack.top = stack.top.next

	stack.size--

	return value, stack
}

// Peek is a function that is used to peek at the top value of the stack, without removing it. Panics if the stack is empty.
//
// Returns:
//   - T: The top value of the stack.
func (stack LinkedStack[T]) Peek() T {
	if stack.top == nil {
		panic(ErrEmptyStack{}.Error())
	}

	return stack.top.value
}

// IsEmpty is a function that is used to check if the stack is empty.
//
// Returns:
//   - bool: True if the stack is empty, false otherwise.
func (stack LinkedStack[T]) IsEmpty() bool {
	return stack.top == nil
}

// Size is a function that is used to get the size of the stack.
//
// Returns:
//   - int: The size of the stack.
func (stack LinkedStack[T]) Size() int {
	return stack.size
}

// ToSlice is a function that is used to get the stack as a slice.
//
// Returns:
//   - []T: The stack as a slice.
//
// Information:
//   - The first element of the slice is the top of the stack.
func (stack LinkedStack[T]) ToSlice() []T {
	slice := make([]T, 0)

	for node := stack.top; node != nil; node = node.next {
		slice = append(slice, node.value)
	}

	return slice
}

// ToString is a function that is used to get the stack as a string.
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
func (stack LinkedStack[T]) ToString(f func(T) string, sep string) string {
	str := ""

	for node := stack.top; node != nil; node = node.next {
		str += f(node.value)

		if node.next != nil {
			str += sep
		}
	}

	return str
}

func (stack LinkedStack[T]) Copy(f func(T) T) Stack[T] {
	new_stack := LinkedStack[T]{
		top:  nil,
		size: stack.size,
	}

	if stack.top == nil {
		return new_stack
	}

	new_stack.top = &node[T]{
		value: f(stack.top.value),
		next:  nil,
	}

	for n := stack.top.next; n != nil; n = n.next {
		// Create a new node.
		new_node := node[T]{
			value: f(n.value),
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

// ArrayStack is a stack implemented using an array.
type ArrayStack[T any] struct {
	elements []T // The elements of the stack.
}

// NewArrayStack creates a new ArrayStack.
//
// Parameters:
//   - elements: The elements to initialize the stack with.
//
// Returns:
//   - ArrayStack: A new ArrayStack.
func NewArrayStack[T any](elements ...T) ArrayStack[T] {
	if elements == nil {
		elements = make([]T, 0)
	}

	return ArrayStack[T]{
		elements: elements,
	}
}

// Push is a function that is used insert a value into the stack at the top.
//
// Parameters:
//   - value: The value to insert into the stack.
func (stack ArrayStack[T]) Push(value T) Stack[T] {
	stack.elements = append(stack.elements, value)

	return stack
}

// Pop is a function that is used to remove a value from the stack at the top. Panics if the stack is empty.
//
// Returns:
//   - T: The value removed from the stack.
func (stack ArrayStack[T]) Pop() (T, Stack[T]) {
	if len(stack.elements) == 0 {
		panic(ErrEmptyStack{}.Error())
	}

	value := stack.elements[len(stack.elements)-1]
	stack.elements = stack.elements[:len(stack.elements)-1]

	return value, stack
}

// Peek is a function that is used to peek at the top value of the stack, without removing it. Panics if the stack is empty.
//
// Returns:
//   - T: The top value of the stack.
func (stack ArrayStack[T]) Peek() T {
	if len(stack.elements) == 0 {
		panic(ErrEmptyStack{}.Error())
	}

	return stack.elements[len(stack.elements)-1]
}

// IsEmpty is a function that is used to check if the stack is empty.
//
// Returns:
//   - bool: True if the stack is empty, false otherwise.
func (stack ArrayStack[T]) IsEmpty() bool {
	return len(stack.elements) == 0
}

// Size is a function that is used to get the size of the stack.
//
// Returns:
//   - int: The size of the stack.
func (stack ArrayStack[T]) Size() int {
	return len(stack.elements)
}

// ToSlice is a function that is used to get the stack as a slice.
//
// Returns:
//   - []T: The stack as a slice.
//
// Information:
//   - The first element of the slice is the top of the stack.
func (stack ArrayStack[T]) ToSlice() []T {
	slice := make([]T, len(stack.elements))
	copy(slice, stack.elements)

	return slice
}

// ToString is a function that is used to get the stack as a string.
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
func (stack ArrayStack[T]) ToString(f func(T) string, sep string) string {
	str := ""

	for i, element := range stack.elements {
		str += f(element)

		if i < len(stack.elements)-1 {
			str += sep
		}
	}

	return str
}

func (stack ArrayStack[T]) Copy(f func(T) T) Stack[T] {
	new_stack := ArrayStack[T]{
		elements: make([]T, len(stack.elements)),
	}

	for i, element := range stack.elements {
		new_stack.elements[i] = f(element)
	}

	return new_stack
}
