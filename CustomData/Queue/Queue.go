package Queue

// QueueHead and QueueSep are constants used in the String() method of the Queuer interface
// to format the string representation of a queue.
const (
	// QueueHead is a string constant that represents the start of the queue. It is used to
	// indicate where elements are removed from the queue.
	// The value of QueueHead is "← | ", which visually indicates the direction of element
	// removal.
	QueueHead string = "← | "

	// QueueSep is a string constant that is used as a separator between elements in the string
	// representation of the queue.
	// The value of QueueSep is " | ", which provides a clear visual separation between individual
	// elements in the queue.
	QueueSep string = " | "
)

// Package queue provides a Queuer interface that defines methods for a queue data structure.
//
// Queuer is an interface that defines methods for a queue data structure.
// It includes methods to add and remove elements, check if the queue is empty or full,
// get the size of the queue, convert the queue to a slice, clear the queue, and get a string
// representation of the queue.
type Queuer[T any] interface {
	// The Enqueue method adds a value of type T to the end of the queue.
	Enqueue(value T)

	// The Dequeue method removes and returns the element at the front of the queue.
	// If the queue is empty, it returns an error.
	Dequeue() (T, error)

	// The MustDequeue method is a convenience method that dequeues an element from the queue
	// and returns it.
	// If the queue is empty, it will panic.
	MustDequeue() T

	// The Peek method returns the element at the front of the queue without removing it.
	Peek() (T, error)

	// MustPeek is a method that returns the value at the front of the queue without removing
	// it.
	// If the queue is empty, it will panic.
	MustPeek() T

	// The IsEmpty method checks if the queue is empty and returns a boolean value indicating
	// whether it is empty or not.
	IsEmpty() bool

	// The Size method returns the number of elements currently in the queue.
	Size() int

	// The ToSlice method returns a slice containing all the elements in the queue.
	ToSlice() []T

	// The Clear method is used to remove all elements from the queue, making it empty.
	Clear()

	// The IsFull method checks if the queue is full, meaning it has reached its maximum
	// capacity and cannot accept any more elements.
	IsFull() bool

	// The String method returns a string representation of the Queuer.
	String() string
}

// linkedNode represents a node in a linked list. It holds a value of a generic type and a
// reference to the next node in the list.
//
// The value field is of a generic type T, which can be any type such as int, string, or a
// custom type.
// It represents the value stored in the node.
//
// The next field is a pointer to the next linkedNode in the list. This allows for traversal
// through the linked list by pointing to the subsequent node in the sequence.
type linkedNode[T any] struct {
	value T
	next  *linkedNode[T]
}

// QueueIterator is a generic type in Go that represents an iterator for a queue.
//
// The values field is a slice of type T, which represents the elements stored in the queue.
//
// The currentIndex field is an integer that keeps track of the current index position of the
// iterator in the queue.
// It is used to iterate over the elements in the queue.
type QueueIterator[T any] struct {
	values       []T
	currentIndex int
}

// NewQueueIterator is a function that creates and returns a new QueueIterator object for a
// given queue.
// It takes a queue of type Queuer[T] as an argument, where T can be any type.
//
// The function uses the ToSlice method of the queue to get a slice of its values, and
// initializes the currentIndex to -1, indicating that the iterator is at the start of the
// queue.
//
// The returned QueueIterator can be used to iterate over the elements in the queue.
func NewQueueIterator[T any](queue Queuer[T]) *QueueIterator[T] {
	return &QueueIterator[T]{
		values:       queue.ToSlice(),
		currentIndex: 0,
	}
}

// GetNext is a method of the QueueIterator type. It is used to move the iterator to the next
// element in the queue and return the value of that element.
// If the iterator is at the end of the queue, the method panics by throwing an
// ErrOutOfBoundsIterator error.
//
// This method is typically used in a loop to iterate over all the elements in a queue.
func (iterator *QueueIterator[T]) GetNext() T {
	if len(iterator.values) <= iterator.currentIndex {
		panic(&ErrOutOfBoundsIterator{})
	}

	value := iterator.values[iterator.currentIndex]
	iterator.currentIndex++

	return value
}

// HasNext is a method of the QueueIterator type. It returns true if there are more elements
// that the iterator is pointing to, and false otherwise.
//
// This method is typically used in conjunction with the GetNext method to iterate over and
// access all the elements in a queue.
func (iterator *QueueIterator[T]) HasNext() bool {
	return iterator.currentIndex < len(iterator.values)
}
