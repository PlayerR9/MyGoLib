package Queue

type ErrEmptyQueue struct{}

func (e ErrEmptyQueue) Error() string {
	return "Empty queue"
}

// node is a node in a linked list.
type node[T any] struct {
	value T        // The value of the node.
	next  *node[T] // The next node in the linked list.
}

// A queue is a data structure that is used to store values in a first-in first-out (FIFO) manner. This means that the first value inserted into the queue is
// the first one to be removed.
type Queue[T any] interface {
	// Enqueue is a function that is used to enqueue a value.
	//
	// Parameters:
	//   - value: The value to enqueue.
	Enqueue(value T)

	// Dequeue is a function that is used to dequeue a value. Panics if the queue is empty.
	//
	// Returns:
	//   - T: The value dequeued.
	Dequeue() T

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

	// ToSlice is a function that is used to get the queue as a slice.
	//
	// Returns:
	//   - []T: The queue as a slice.
	//
	// Information:
	//   - The first element of the slice is the front of the queue.
	ToSlice() []T

	// ToString is a function that is used to get the queue as a string.
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
}

// LinkedQueue is a queue implemented using a linked list.
type LinkedQueue[T any] struct {
	front *node[T] // The front of the queue.
	back  *node[T] // The back of the queue.
	size  int      // The size of the queue.
}

// NewLinkedQueue creates a new LinkedQueue.
//
// Parameters:
//   - elements: The elements to initialize the queue with.
//
// Returns:
//   - LinkedQueue: A new LinkedQueue.
func NewLinkedQueue[T any](elements []T) LinkedQueue[T] {
	queue := LinkedQueue[T]{
		front: nil,
		back:  nil,
		size:  0,
	}

	for _, element := range elements {
		queue.Enqueue(element)
	}

	return queue
}

// Enqueue is a function that is used enqueue a value.
//
// Parameters:
//   - value: The value to enqueue.
func (queue *LinkedQueue[T]) Enqueue(value T) {
	node := &node[T]{
		value: value,
		next:  nil,
	}

	if queue.back == nil {
		queue.front = node
	} else {
		queue.back.next = node
	}

	queue.back = node

	queue.size++
}

// Dequeue is a function that is used to dequeue a value. Panics if the queue is empty.
//
// Returns:
//   - T: The value dequeued.
func (queue *LinkedQueue[T]) Dequeue() T {
	if queue.front == nil {
		panic(ErrEmptyQueue{}.Error())
	}

	value := queue.front.value
	queue.front = queue.front.next

	if queue.front == nil {
		queue.back = nil
	}

	queue.size--

	return value
}

// Peek is a function that is used to peek at the front value, without removing it. Panics if the queue is empty.
//
// Returns:
//   - T: The front value.
func (queue LinkedQueue[T]) Peek() T {
	if queue.front == nil {
		panic(ErrEmptyQueue{}.Error())
	}

	return queue.front.value
}

// IsEmpty is a function that is used to check if the queue is empty.
//
// Returns:
//   - bool: True if the queue is empty, false otherwise.
func (queue LinkedQueue[T]) IsEmpty() bool {
	return queue.front == nil
}

// Size is a function that is used to get the size of the queue.
//
// Returns:
//   - int: The size of the queue.
func (queue LinkedQueue[T]) Size() int {
	return queue.size
}

// ToSlice is a function that is used to get the queue as a slice.
//
// Returns:
//   - []T: The queue as a slice.
//
// Information:
//   - The first element of the slice is the front of the queue.
func (queue LinkedQueue[T]) ToSlice() []T {
	slice := make([]T, queue.Size())

	i := 0
	for node := queue.front; node != nil; node = node.next {
		slice[i] = node.value
		i++
	}

	return slice
}

// ToString is a function that is used to get the queue as a string.
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
func (queue LinkedQueue[T]) ToString(f func(T) string, sep string) string {
	str := ""

	for i, node := 0, queue.front; node != nil; i, node = i+1, node.next {
		str += f(node.value)

		if i < queue.Size()-1 {
			str += sep
		}
	}

	return str
}

// ArrayQueue is a queue implemented using an array.
type ArrayQueue[T any] struct {
	elements []T // The elements of the queue.
	front    int // The index of the front of the queue.
	back     int // The index of the back of the queue.
}

// NewArrayQueue creates a new ArrayQueue.
//
// Parameters:
//   - elements: The elements to initialize the queue with.
//
// Returns:
//   - ArrayQueue: A new ArrayQueue.
func NewArrayQueue[T any](elements []T) ArrayQueue[T] {
	queue := ArrayQueue[T]{
		elements: make([]T, len(elements)),
		front:    0,
		back:     len(elements) - 1,
	}

	copy(queue.elements, elements)

	return queue
}

// Enqueue is a function that is used enqueue a value.
//
// Parameters:
//   - value: The value to enqueue.
func (queue *ArrayQueue[T]) Enqueue(value T) {
	queue.back = (queue.back + 1) % len(queue.elements)
	queue.elements[queue.back] = value
}

// Dequeue is a function that is used to dequeue a value. Panics if the queue is empty.
//
// Returns:
//   - T: The value dequeued.
func (queue *ArrayQueue[T]) Dequeue() T {
	if queue.front == queue.back {
		panic(ErrEmptyQueue{}.Error())
	}

	value := queue.elements[queue.front]
	queue.front = (queue.front + 1) % len(queue.elements)

	return value
}

// Peek is a function that is used to peek at the front value, without removing it. Panics if the queue is empty.
//
// Returns:
//   - T: The front value.
func (queue ArrayQueue[T]) Peek() T {
	if queue.front == queue.back {
		panic(ErrEmptyQueue{}.Error())
	}

	return queue.elements[queue.front]
}

// IsEmpty is a function that is used to check if the queue is empty.
//
// Returns:
//   - bool: True if the queue is empty, false otherwise.
func (queue ArrayQueue[T]) IsEmpty() bool {
	return queue.front == queue.back
}

// Size is a function that is used to get the size of the queue.
//
// Returns:
//   - int: The size of the queue.
func (queue ArrayQueue[T]) Size() int {
	return (queue.back+1)%len(queue.elements) - queue.front
}

// ToSlice is a function that is used to get the queue as a slice.
//
// Returns:
//   - []T: The queue as a slice.
//
// Information:
//   - The first element of the slice is the front of the queue.
func (queue ArrayQueue[T]) ToSlice() []T {
	slice := make([]T, queue.Size())

	for i, j := queue.front, 0; i != queue.back; i, j = (i+1)%len(queue.elements), j+1 {
		slice[j] = queue.elements[i]
	}

	slice[queue.Size()-1] = queue.elements[queue.back]

	return slice
}

// ToString is a function that is used to get the queue as a string.
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
func (queue ArrayQueue[T]) ToString(f func(T) string, sep string) string {
	str := ""

	for i, j := queue.front, 0; i != queue.back; i, j = (i+1)%len(queue.elements), j+1 {
		str += f(queue.elements[i]) + sep
	}

	str += f(queue.elements[queue.back])

	return str
}
