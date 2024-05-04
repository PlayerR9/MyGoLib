package Queuer

import (
	"fmt"
	"strings"

	fs "github.com/PlayerR9/MyGoLib/Formatting/Strings"
	itf "github.com/PlayerR9/MyGoLib/ListLike/Iterator"
	itff "github.com/PlayerR9/MyGoLib/Units/Interfaces"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
)

// LimitedLinkedQueue is a generic type that represents a queue data structure with
// or without a limited capacity, implemented using a linked list.
type LimitedLinkedQueue[T any] struct {
	// front and back are pointers to the first and last nodes in the linked queue,
	// respectively.
	front, back *QueueNode[T]

	// size is the current number of elements in the queue.
	size int

	// capacity is the maximum number of elements the queue can hold.
	capacity int
}

// NewLimitedLinkedQueue is a function that creates and returns a new instance of a
// LimitedLinkedQueue.
//
// Parameters:
//
//   - values: A variadic parameter of type T, which represents the initial values to be
//     stored in the queue.
//
// Returns:
//
//   - *LimitedLinkedQueue[T]: A pointer to the newly created LimitedLinkedQueue.
func NewLimitedLinkedQueue[T any](values ...T) *LimitedLinkedQueue[T] {
	queue := new(LimitedLinkedQueue[T])
	queue.size = len(values)

	if len(values) == 0 {
		return queue
	}

	// First node
	node := NewQueueNode(values[0])

	queue.front = node
	queue.back = node

	// Subsequent nodes
	for _, element := range values[1:] {
		node = NewQueueNode(element)

		queue.back.SetNext(node)
		queue.back = node
	}

	return queue
}

// Enqueue is a method of the LimitedLinkedQueue type. It is used to add an element to
// the end of the queue.
//
// Panics with an error of type *ErrCallFailed if the queue is full.
//
// Parameters:
//
//   - value: A pointer to a value of type T, which is the element to be added to the
//     queue.
func (queue *LimitedLinkedQueue[T]) Enqueue(value T) error {
	if queue.size >= queue.capacity {
		return NewErrFullList(queue)
	}

	queue_node := NewQueueNode(value)

	if queue.back == nil {
		queue.front = queue_node
	} else {
		queue.back.SetNext(queue_node)
	}

	queue.back = queue_node

	queue.size++

	return nil
}

// Dequeue is a method of the LimitedLinkedQueue type. It is used to remove and return
// the element at the front of the queue.
//
// Panics with an error of type *ErrCallFailed if the queue is empty.
//
// Returns:
//
//   - T: The value of the element at the front of the queue.
func (queue *LimitedLinkedQueue[T]) Dequeue() (T, error) {
	if queue.front == nil {
		return *new(T), NewErrEmptyList(queue)
	}

	toRemove := queue.front

	queue.front = queue.front.Next()
	if queue.front == nil {
		queue.back = nil
	}

	queue.size--
	toRemove.SetNext(nil)

	return toRemove.Value, nil
}

// Peek is a method of the LimitedLinkedQueue type. It is used to return the element at
// the front of the queue without removing it.
//
// Panics with an error of type *ErrCallFailed if the queue is empty.
//
// Returns:
//
//   - T: The value of the element at the front of the queue.
func (queue *LimitedLinkedQueue[T]) Peek() (T, error) {
	if queue.front == nil {
		return *new(T), NewErrEmptyList(queue)
	}

	return queue.front.Value, nil
}

// IsEmpty is a method of the LimitedLinkedQueue type. It is used to check if the queue
// is empty.
//
// Returns:
//
//   - bool: A boolean value indicating whether the queue is empty.
func (queue *LimitedLinkedQueue[T]) IsEmpty() bool {
	return queue.front == nil
}

// Size is a method of the LimitedLinkedQueue type. It is used to return the current
// number of elements in the queue.
//
// Returns:
//
//   - int: An integer representing the current number of elements in the queue.
func (queue *LimitedLinkedQueue[T]) Size() int {
	return queue.size
}

func (queue *LimitedLinkedQueue[T]) Capacity() int {
	return queue.capacity
}

// Iterator is a method of the LimitedLinkedQueue type. It is used to return an iterator
// for the queue.
//
// Returns:
//
//   - itf.Iterater[T]: An iterator for the queue.
func (queue *LimitedLinkedQueue[T]) Iterator() itf.Iterater[T] {
	var builder itf.Builder[T]

	for queue_node := queue.front; queue_node != nil; queue_node = queue_node.Next() {
		builder.Append(queue_node.Value)
	}

	return builder.Build()
}

// Clear is a method of the LimitedLinkedQueue type. It is used to remove all elements
// from the queue.
func (queue *LimitedLinkedQueue[T]) Clear() {
	if queue.front == nil {
		return // Queue is already empty
	}

	// 1. First node
	prev := queue.front

	// 2. Subsequent nodes
	for node := queue.front.Next(); node != nil; node = node.Next() {
		prev = node
		prev.SetNext(nil)
	}

	prev.SetNext(nil)

	// 3. Reset queue fields
	queue.front = nil
	queue.back = nil
	queue.size = 0
}

// IsFull is a method of the LimitedLinkedQueue type. It is used to check if the queue is
// full.
//
// Returns:
//
//   - isFull: A boolean value indicating whether the queue is full.
func (queue *LimitedLinkedQueue[T]) IsFull() bool {
	return queue.size >= queue.capacity
}

// String is a method of the LimitedLinkedQueue type. It is used to return a string
// representation of the queue including its size, capacity, and elements.
//
// Returns:
//
//   - string: A string representation of the queue.
func (queue *LimitedLinkedQueue[T]) String() string {
	values := make([]string, 0, queue.size)
	for queue_node := queue.front; queue_node != nil; queue_node = queue_node.Next() {
		values = append(values, fs.StringOf(queue_node.Value))
	}

	return fmt.Sprintf(
		"LimitedLinkedQueue[capacity=%d, size=%d, values=[← %s]]",
		queue.capacity,
		queue.size,
		strings.Join(values, ", "),
	)
}

// CutNilValues is a method of the LimitedLinkedQueue type. It is used to remove all nil
// values from the queue.
func (queue *LimitedLinkedQueue[T]) CutNilValues() {
	if queue.front == nil {
		return // Queue is empty
	}

	if gen.IsNil(queue.front.Value) && queue.front == queue.back {
		// Single node
		queue.front = nil
		queue.back = nil
		queue.size = 0

		return
	}

	var toDelete *QueueNode[T] = nil

	// 1. First node
	if gen.IsNil(queue.front.Value) {
		toDelete = queue.front

		queue.front = queue.front.Next()

		toDelete.SetNext(nil)
		queue.size--
	}

	prev := queue.front

	// 2. Subsequent nodes (except last)
	for node := queue.front.Next(); node.Next() != nil; node = node.Next() {
		if !gen.IsNil(node.Value) {
			prev = node
		} else {
			prev.SetNext(node.Next())
			queue.size--

			if toDelete != nil {
				toDelete.SetNext(nil)
			}

			toDelete = node
		}
	}

	if toDelete != nil {
		toDelete.SetNext(nil)
	}

	// 3. Last node
	if gen.IsNil(queue.back.Value) {
		prev.SetNext(nil)
		queue.back = prev
		queue.size--
	}
}

// Slice is a method of the LimitedLinkedQueue type. It is used to return a slice of the
// elements in the queue.
//
// Returns:
//
//   - []T: A slice of the elements in the queue.
func (queue *LimitedLinkedQueue[T]) Slice() []T {
	slice := make([]T, 0, queue.size)

	for queue_node := queue.front; queue_node != nil; queue_node = queue_node.Next() {
		slice = append(slice, queue_node.Value)
	}

	return slice
}

// Copy is a method of the LimitedLinkedQueue type. It is used to create a shallow copy
// of the queue.
//
// Returns:
//
//   - itf.Copier: A copy of the queue.
func (queue *LimitedLinkedQueue[T]) Copy() itff.Copier {
	queueCopy := &LimitedLinkedQueue[T]{
		size:     queue.size,
		capacity: queue.capacity,
	}

	if queue.size == 0 {
		return queueCopy
	}

	// First node
	node := NewQueueNode(queue.front.Value)

	queueCopy.front = node
	queueCopy.back = node

	// Subsequent nodes
	for queue_node := queue.front.Next(); queue_node != nil; queue_node = queue_node.Next() {
		node = NewQueueNode(queue_node.Value)

		queueCopy.back.SetNext(node)
		queueCopy.back = node
	}

	return queueCopy
}
