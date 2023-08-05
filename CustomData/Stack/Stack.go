package Stack

/* node is a node in a linked list. */
type node[T any] struct {
	value T        // The value of the node.
	next  *node[T] // The next node in the linked list.
}

/* StackInterface is an interface for a stack. */
type Stack[T any] interface {
	Push(T)                                 // Pushes a value onto the stack
	Pop() T                                 // Pops a value off the stack. Panics if the stack is empty.
	Peek() T                                // Peeks at the top value of the stack. Panics if the stack is empty.
	IsEmpty() bool                          // Returns true if the stack is empty, false otherwise.
	Size() int                              // Returns the size of the stack.
	ToSlice() []T                           // Returns the stack as a slice.
	ToString(func(T) string, string) string // Returns the stack as a string.
}

/* LinkedStack is a stack implemented using a linked list. */
type LinkedStack[T any] struct {
	top *node[T] // The top of the stack.
}

/* NewLinkedStack creates a new LinkedStack.
Parameters:
	elements: The elements to initialize the stack with.
Returns:
	A new LinkedStack.
Complexity:
	O(n), where n is the number of elements.
*/
func NewLinkedStack[T any](elements []T) LinkedStack[T] {
	stack := LinkedStack[T]{
		top: nil,
	}

	for _, element := range elements {
		stack.Push(element)
	}

	return stack
}

/* Push pushes a value onto the stack.
Parameters:
	value: The value to push onto the stack.
Complexity:
	O(1).
*/
func (stack *LinkedStack[T]) Push(value T) {
	new_node := node[T]{
		value: value,
		next:  stack.top,
	}

	stack.top = &new_node
}

/* Pop pops a value off the stack. Panics if the stack is empty.
Returns:
	The value popped off the stack.
Complexity:
	O(1).
*/
func (stack *LinkedStack[T]) Pop() T {
	if stack.top == nil {
		panic("cannot pop from empty stack")
	}

	value := stack.top.value
	stack.top = stack.top.next

	return value
}

/* Peek peeks at the top value of the stack. Panics if the stack is empty.
Returns:
	The top value of the stack.
Complexity:
	O(1).
*/
func (stack LinkedStack[T]) Peek() T {
	if stack.top == nil {
		panic("cannot peek empty stack")
	}

	return stack.top.value
}

/* IsEmpty returns true if the stack is empty, false otherwise.
Returns:
	true if the stack is empty, false otherwise.
Complexity:
	O(1).
*/
func (stack LinkedStack[T]) IsEmpty() bool {
	return stack.top == nil
}

/* Size returns the size of the stack.
Returns:
	The size of the stack.
Complexity:
	O(n), where n is the number of elements.
*/
func (stack LinkedStack[T]) Size() int {
	size := 0

	for node := stack.top; node != nil; node = node.next {
		size++
	}

	return size
}

/* ToSlice returns the stack as a slice.
Returns:
	The stack as a slice.
Complexity:
	O(n), where n is the number of elements.
*/
func (stack LinkedStack[T]) ToSlice() []T {
	slice := make([]T, 0)

	for node := stack.top; node != nil; node = node.next {
		slice = append(slice, node.value)
	}

	return slice
}

/* ToString returns the stack as a string.
Parameters:
	f: The function to convert the values to strings.
	sep: The separator to use between values.
Returns:
	The stack as a string.
Complexity:
	O(n), where n is the number of elements.
*/
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

/* ArrayStack is a stack implemented using an array. */
type ArrayStack[T any] struct {
	elements []T // The elements of the stack.
}

/* NewArrayStack creates a new ArrayStack.
Parameters:
	elements: The elements to initialize the stack with.
Returns:
	A new ArrayStack.
Complexity:
	O(n), where n is the number of elements.
*/
func NewArrayStack[T any](elements []T) ArrayStack[T] {
	if elements == nil {
		elements = make([]T, 0)
	}

	return ArrayStack[T]{
		elements: elements,
	}
}

/* Push pushes a value onto the stack.
Parameters:
	value: The value to push onto the stack.
Complexity:
	O(1) amortized.
*/
func (stack *ArrayStack[T]) Push(value T) {
	stack.elements = append(stack.elements, value)
}

/* Pop pops a value off the stack. Panics if the stack is empty.
Returns:
	The value popped off the stack.
Complexity:
	O(1) amortized.
*/
func (stack *ArrayStack[T]) Pop() T {
	if len(stack.elements) == 0 {
		panic("cannot pop from empty stack")
	}

	value := stack.elements[len(stack.elements)-1]
	stack.elements = stack.elements[:len(stack.elements)-1]

	return value
}

/* Peek peeks at the top value of the stack. Panics if the stack is empty.
Returns:
	The top value of the stack.
Complexity:
	O(1).
*/
func (stack ArrayStack[T]) Peek() T {
	if len(stack.elements) == 0 {
		panic("cannot peek empty stack")
	}

	return stack.elements[len(stack.elements)-1]
}

/* IsEmpty returns true if the stack is empty, false otherwise.
Returns:
	true if the stack is empty, false otherwise.
Complexity:
	O(1).
*/
func (stack ArrayStack[T]) IsEmpty() bool {
	return len(stack.elements) == 0
}

/* Size returns the size of the stack.
Returns:
	The size of the stack.
Complexity:
	O(1).
*/
func (stack ArrayStack[T]) Size() int {
	return len(stack.elements)
}

/* ToSlice returns the stack as a slice.
Returns:
	The stack as a slice.
Complexity:
	O(n), where n is the number of elements.
*/
func (stack ArrayStack[T]) ToSlice() []T {
	slice := make([]T, len(stack.elements))
	copy(slice, stack.elements)

	return slice
}

/* ToString returns the stack as a string.
Parameters:
	f: The function to convert the values to strings.
	sep: The separator to use between values.
Returns:
	The stack as a string.
Complexity:
	O(n), where n is the number of elements.
*/
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
