// Package ListLike provides a Queuer interface that defines methods for a queue data
// structure.
package ListLike

import (
	"errors"
	"fmt"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
	"github.com/markphelps/optional"
)

// Queuer is an interface that defines methods for a queue data structure.
type Queuer[T any] interface {
	// The Enqueue method adds a value of type T to the end of the queue.
	// If the queue is full, it will panic.
	Enqueue(value *T)

	// The Dequeue method is a convenience method that dequeues an element from the
	// queue and returns it.
	// If the queue is empty, it will panic.
	Dequeue() *T

	// Peek is a method that returns the value at the front of the queue without
	// removing it.
	// If the queue is empty, it will panic.
	Peek() *T

	// The IsEmpty method checks if the queue is empty and returns a boolean value
	// indicating whether it is empty or not.
	IsEmpty() bool

	// The Size method returns the number of elements currently in the queue.
	Size() int

	// The Capacity method returns the maximum number of elements that the queue can hold.
	Capacity() optional.Int

	// The ToSlice method returns a slice containing all the elements in the queue.
	ToSlice() []*T

	// The Clear method is used to remove all elements from the queue, making it
	// empty.
	Clear()

	// The IsFull method checks if the queue is full, meaning it has reached its
	// maximum capacity and cannot accept any more elements.
	IsFull() bool

	// The String method returns a string representation of the queue.
	fmt.Stringer

	// The CutNilValues method is used to remove all nil values from the queue.
	CutNilValues()
}

// linkedNode represents a node in a linked list.
type linkedNode[T any] struct {
	// value is the value stored in the node.
	value *T

	// next is a pointer to the next linkedNode in the list.
	next *linkedNode[T]
}

// QueueIterator is a generic type that represents an iterator for a queue.
// It is used to iterate over the elements in a queue.
type QueueIterator[T any] struct {
	// values is a slice of type T that represents the elements stored in the queue.
	values []*T

	// currentIndex is an integer that keeps track of the current index position of
	// the iterator in the queue.
	currentIndex int
}

// NewQueueIterator is a function that creates and returns a new QueueIterator
// object for a given queue.
//
// Parameters:
//
//   - queue: A queue of type Queuer[T] that the iterator will be created for.
//
// Returns:
//
//   - *QueueIterator[T]: A pointer to a new QueueIterator object.
func NewQueueIterator[T any](queue Queuer[T]) *QueueIterator[T] {
	return &QueueIterator[T]{
		values:       queue.ToSlice(),
		currentIndex: 0,
	}
}

// GetNext is a method of the QueueIterator type. It is used to move the iterator
// to the next element in the queue and return the value of that element.
//
// Panics with an error of type *ErrOperationFailed if the iterator is out of bounds.
//
// Returns:
//
//   - *T: A pointer to the next element in the queue.
func (iterator *QueueIterator[T]) GetNext() *T {
	if len(iterator.values) > iterator.currentIndex {
		value := iterator.values[iterator.currentIndex]
		iterator.currentIndex++

		return value
	}

	panic(ers.NewErrOperationFailed(
		"get next element", errors.New("iterator is out of bounds")),
	)
}

// HasNext is a method of the QueueIterator type. It is used to check if there are
// more elements in the queue that the iterator is pointing to.
//
// Returns:
//
//   - bool: A boolean value indicating whether there are more elements in the queue.
func (iterator *QueueIterator[T]) HasNext() bool {
	return iterator.currentIndex < len(iterator.values)
}
