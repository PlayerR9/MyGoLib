package Queue

import (
	"log"
	"os"
)

// ERRORS

// ErrEmptyQueue is an error that is used when the queue is empty.
type ErrEmptyQueue struct{}

func (e ErrEmptyQueue) Error() string {
	return "Empty queue"
}

// ErrFullQueue is an error that is used when the queue is full.
type ErrFullQueue struct{}

func (e ErrFullQueue) Error() string {
	return "Full queue"
}

// GLOBAL VARIABLES

var (
	// DebugMode is a boolean that is used to enable or disable debug mode. When debug mode is enabled, the package will print debug messages.
	// **Note:** Debug mode is disabled by default.
	DebugMode bool        = false
	debugger  *log.Logger = log.New(os.Stdout, "[Queue] ", log.LstdFlags) // debugger is a logger that is used to print debug messages.
)

// node is a node in a linked list.
type node[T any] struct {
	value T        // The value of the node.
	next  *node[T] // The next node in the linked list.
}

// A Queue is a data structure that is used to store values in a first-in first-out (FIFO) manner. This means that the first value inserted into the queue is
// the first one to be removed.
type Queue[T any] interface {
	// Enqueue is a function that is used to enqueue a value.
	//
	// Parameters:
	//   - value: The value to enqueue.
	//
	// Returns:
	//   - Queue[T]: The queue.
	Enqueue(value T) Queue[T]

	// Dequeue is a function that is used to dequeue a value. Panics if the queue is empty.
	//
	// Returns:
	//   - T: The value dequeued.
	//   - Queue[T]: The queue.
	Dequeue() (T, Queue[T])

	// Peek is a function that is used to peek at the front value, without removing it. Panics if the queue is empty.
	//
	// Returns:
	//   - T: The front value.
	Peek() T

	// IsEmpty is a function that is used to check if the queue is empty.
	//
	// Returns:
	//   - bool: True if the queue is empty, false otherwise.
	IsEmpty() bool

	// Size is a function that is used to get the size of the queue.
	//
	// Returns:
	//   - int: The size of the queue.
	Size() int

	// ToSlice is a function that is used to get the queue as a slice. The first element of the slice is the front of the queue.
	//
	// Returns:
	//   - []T: The queue as a slice.
	//
	// Information:
	//   - The first element of the slice is the front of the queue.
	ToSlice() []T

	// ToString is a function that is used to get the queue as a string. The first element of the string is the front of the queue.
	//
	// Parameters:
	//   - f: The function to convert the values to strings.
	//   - sep: The separator to use between values.
	//
	// Returns:
	//   - string: The queue as a string.
	//
	// Information:
	//   - The first element of the string is the front of the queue.
	ToString(func(T) string, string) string

	// Copy is a function that is used to get a copy of the queue.
	//
	// Returns:
	//   - Queue[T]: The copy of the queue.
	Copy() Queue[T]
}

// LinkedQueue is a queue implemented using a linked list.
type LinkedQueue[T any] struct {
	front         *node[T]  // The front of the queue.
	back          *node[T]  // The back of the queue.
	size          int       // The size of the queue.
	copy_function func(T) T // The function to copy the values.
}

// NewLinkedQueue creates a new LinkedQueue.
//
// Parameters:
//   - copy_function: The function to copy the values.
//   - elements: The elements to initialize the queue with.
//
// Returns:
//   - LinkedQueue: A new LinkedQueue.
func NewLinkedQueue[T any](copy_function func(T) T, elements ...T) LinkedQueue[T] {
	// Create a new queue.
	queue := LinkedQueue[T]{
		front:         nil,
		back:          nil,
		size:          0,
		copy_function: copy_function,
	}

	// Enqueue the elements to the queue.
	for _, element := range elements {
		queue = queue.Enqueue(element).(LinkedQueue[T])
	}

	return queue
}

func (queue LinkedQueue[T]) Enqueue(value T) Queue[T] {
	// Create a new node.
	node := &node[T]{
		value: queue.copy_function(value),
		next:  nil,
	}

	if queue.back == nil {
		queue.front = node // If the queue is empty, set the front to the new node.
	} else {
		queue.back.next = node // Otherwise, set the next of the back to the new node.
	}

	queue.back = node // Set the back to the new node.

	queue.size++ // Increment the size.

	return queue
}

func (queue LinkedQueue[T]) Dequeue() (T, Queue[T]) {
	if queue.front == nil {
		// If the queue is empty, panic.
		if DebugMode {
			debugger.Panic(ErrEmptyQueue{}.Error())
		} else {
			panic(ErrEmptyQueue{}.Error())
		}
	}

	value := queue.front.value     // Get the value at the front of the queue.
	queue.front = queue.front.next // Set the front to the next of the front.

	if queue.front == nil {
		// If the front is nil, set the back to nil.
		queue.back = nil
	}

	queue.size-- // Decrement the size.

	return queue.copy_function(value), queue
}

func (queue LinkedQueue[T]) Peek() T {
	if queue.front == nil {
		// If the queue is empty, panic.
		if DebugMode {
			debugger.Panic(ErrEmptyQueue{}.Error())
		} else {
			panic(ErrEmptyQueue{}.Error())
		}
	}

	return queue.copy_function(queue.front.value) // Get the value at the front of the queue.
}

func (queue LinkedQueue[T]) IsEmpty() bool {
	return queue.front == nil // If the front is nil, the queue is empty.
}

func (queue LinkedQueue[T]) Size() int {
	return queue.size
}

func (queue LinkedQueue[T]) ToSlice() []T {
	slice := make([]T, queue.Size()) // Create a new slice with the size of the queue.

	// Iterate over the queue and add the values to the slice.
	i := 0
	for node := queue.front; node != nil; node = node.next {
		slice[i] = queue.copy_function(node.value)
		i++
	}

	return slice
}

func (queue LinkedQueue[T]) ToString(f func(T) string, sep string) string {
	if queue.front == nil {
		// If the queue is empty, return an empty string.
		return ""
	}

	str := f(queue.front.value) // Add the first value without the separator.

	// Add the rest of the values with the separator.
	for node := queue.front.next; node != nil; node = node.next {
		str += sep + f(node.value)
	}

	return str
}

func (queue LinkedQueue[T]) Copy() Queue[T] {
	// Create a new queue.
	obj := LinkedQueue[T]{
		front:         nil,
		back:          nil,
		size:          0,
		copy_function: queue.copy_function,
	}

	// Iterate over the queue and add the values to the new queue.
	for node := queue.front; node != nil; node = node.next {
		obj.Enqueue(queue.copy_function(node.value))
	}

	return obj
}

// ArrayQueue is a queue implemented using an array. If maximum size is reached, no more elements can be added.
type ArrayQueue[T any] struct {
	elements      []T       // The elements of the queue.
	hasMaxSize    bool      // Whether the queue has a maximum size.
	copy_function func(T) T // The function to copy the values.
}

// NewArrayQueue creates a new ArrayQueue with no maximum size.
//
// Parameters:
//   - copy_function: The function to copy the values.
//   - elements: The elements to initialize the queue with.
//
// Returns:
//   - ArrayQueue: A new ArrayQueue.
func NewArrayQueue[T any](copy_function func(T) T, elements ...T) ArrayQueue[T] {
	// Create a new queue.
	queue := ArrayQueue[T]{
		elements:      make([]T, len(elements)),
		hasMaxSize:    false,
		copy_function: copy_function,
	}

	// Copy the elements to the queue.
	for i, element := range elements {
		queue.elements[i] = queue.copy_function(element)
	}

	return queue
}

// NewArrayQueueWithMaxSize creates a new ArrayQueue with a maximum size. Panics if maxSize is negative.
//
// Parameters:
//   - maxSize: The maximum size of the queue.
//   - elements: The elements to initialize the queue with.
//
// Returns:
//   - ArrayQueue: A new ArrayQueue.
func NewArrayQueueWithMaxSize[T any](maxSize int, elements ...T) ArrayQueue[T] {
	if maxSize < 0 {
		// If maxSize is negative, panic.
		if DebugMode {
			debugger.Panic("maxSize cannot be negative")
		} else {
			panic("maxSize cannot be negative")
		}
	}

	if len(elements) > maxSize {
		// If the number of elements is greater than the size, panic.
		if DebugMode {
			debugger.Panic("maxSize is smaller than the number of elements")
		} else {
			panic("maxSize is smaller than the number of elements")
		}
	}

	// Create a new queue.
	queue := ArrayQueue[T]{
		elements:   make([]T, len(elements), maxSize),
		hasMaxSize: true,
	}

	// Copy the elements to the queue.
	for i, element := range elements {
		queue.elements[i] = queue.copy_function(element)
	}

	return queue
}

func (queue ArrayQueue[T]) Enqueue(value T) Queue[T] {
	if queue.hasMaxSize && len(queue.elements) == cap(queue.elements) {
		// If the queue has a maximum size and the queue is full, panic.
		if DebugMode {
			debugger.Panic(ErrFullQueue{}.Error())
		} else {
			panic(ErrFullQueue{}.Error())
		}
	}

	// Append the value to the queue.
	queue.elements = append(queue.elements, queue.copy_function(value))

	return queue
}

func (queue ArrayQueue[T]) Dequeue() (T, Queue[T]) {
	if len(queue.elements) == 0 {
		// If the queue is empty, panic.
		if DebugMode {
			debugger.Panic(ErrEmptyQueue{}.Error())
		} else {
			panic(ErrEmptyQueue{}.Error())
		}
	}

	value := queue.elements[0]          // Get the value at the front of the queue.
	queue.elements = queue.elements[1:] // Remove the value from the queue.

	return queue.copy_function(value), queue
}

func (queue ArrayQueue[T]) Peek() T {
	if len(queue.elements) == 0 {
		// If the queue is empty, panic.
		if DebugMode {
			debugger.Panic(ErrEmptyQueue{}.Error())
		} else {
			panic(ErrEmptyQueue{}.Error())
		}
	}

	return queue.copy_function(queue.elements[0]) // Get the value at the front of the queue.
}

func (queue ArrayQueue[T]) IsEmpty() bool {
	return len(queue.elements) == 0 // If the queue is empty, the queue is empty.
}

func (queue ArrayQueue[T]) Size() int {
	return len(queue.elements)
}

func (queue ArrayQueue[T]) ToSlice() []T {
	slice := make([]T, len(queue.elements)) // Create a new slice with the size of the queue.

	for i, element := range queue.elements {
		slice[i] = queue.copy_function(element)
	}

	return slice
}

func (queue ArrayQueue[T]) ToString(f func(T) string, sep string) string {
	var str string

	if len(queue.elements) != 0 {
		str = f(queue.elements[0]) // Add the first element without the separator.

		// Add the rest of the elements with the separator.
		for _, element := range queue.elements[1:] {
			str += sep + f(element)
		}
	} else {
		str = ""
	}

	return "← | " + str
}

func (queue ArrayQueue[T]) Copy() Queue[T] {
	// Create a new queue.
	var obj ArrayQueue[T]

	if queue.hasMaxSize {
		obj = ArrayQueue[T]{
			elements:      make([]T, 0, cap(queue.elements)),
			hasMaxSize:    true,
			copy_function: queue.copy_function,
		}
	} else {
		obj = ArrayQueue[T]{
			elements:      make([]T, 0),
			hasMaxSize:    false,
			copy_function: queue.copy_function,
		}
	}

	// Iterate over the queue and add the values to the new queue.
	for _, element := range queue.elements {
		obj.elements = append(obj.elements, obj.copy_function(element))
	}

	return obj
}
