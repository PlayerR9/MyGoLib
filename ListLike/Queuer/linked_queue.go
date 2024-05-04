package Queuer

import (
	"fmt"
	"strings"

	fs "github.com/PlayerR9/MyGoLib/Formatting/Strings"
	itf "github.com/PlayerR9/MyGoLib/ListLike/Iterator"
	itff "github.com/PlayerR9/MyGoLib/Units/Interfaces"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
)

// LinkedQueue is a generic type that represents a queue data structure with
// or without a limited capacity, implemented using a linked list.
type LinkedQueue[T any] struct {
	// front and back are pointers to the first and last nodes in the linked queue,
	// respectively.
	front, back *QueueNode[T]

	// size is the current number of elements in the queue.
	size int
}

// NewLinkedQueue is a function that creates and returns a new instance of a
// LinkedQueue.
//
// Parameters:
//
//   - values: A variadic parameter of type T, which represents the initial values to be
//     stored in the queue.
//
// Returns:
//
//   - *LinkedQueue[T]: A pointer to the newly created LinkedQueue.
func NewLinkedQueue[T any](values ...T) *LinkedQueue[T] {
	queue := new(LinkedQueue[T])
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

// Enqueue is a method of the LinkedQueue type. It is used to add an element to
// the end of the queue.
//
// Panics with an error of type *ErrCaCommonFailed if the queue is fuCommon.
//
// Parameters:
//
//   - value: A pointer to a value of type T, which is the element to be added to the
//     queue.
func (queue *LinkedQueue[T]) Enqueue(value T) error {
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

// Dequeue is a method of the LinkedQueue type. It is used to remove and return
// the element at the front of the queue.
//
// Panics with an error of type *ErrEmptyList if the queue is empty.
//
// Returns:
//
//   - T: The value of the element at the front of the queue.
func (queue *LinkedQueue[T]) Dequeue() (T, error) {
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

// Peek is a method of the LinkedQueue type. It is used to return the element at
// the front of the queue without removing it.
//
// Panics with an error of type *ErrEmptyList if the queue is empty.
//
// Returns:
//
//   - T: The value of the element at the front of the queue.
func (queue *LinkedQueue[T]) Peek() (T, error) {
	if queue.front == nil {
		return *new(T), NewErrEmptyList(queue)
	}

	return queue.front.Value, nil
}

// IsEmpty is a method of the LinkedQueue type. It is used to check if the queue
// is empty.
//
// Returns:
//
//   - bool: A boolean value indicating whether the queue is empty.
func (queue *LinkedQueue[T]) IsEmpty() bool {
	return queue.front == nil
}

// Size is a method of the LinkedQueue type. It is used to return the current
// number of elements in the queue.
//
// Returns:
//
//   - int: An integer representing the current number of elements in the queue.
func (queue *LinkedQueue[T]) Size() int {
	return queue.size
}

// Iterator is a method of the LinkedQueue type. It is used to return an iterator
// for the queue.
//
// Returns:
//
//   - itf.Iterater[T]: An iterator for the queue.
func (queue *LinkedQueue[T]) Iterator() itf.Iterater[T] {
	var builder itf.Builder[T]

	for queue_node := queue.front; queue_node != nil; queue_node = queue_node.Next() {
		builder.Append(queue_node.Value)
	}

	return builder.Build()
}

// Clear is a method of the LinkedQueue type. It is used to remove aCommon elements
// from the queue.
func (queue *LinkedQueue[T]) Clear() {
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

// String is a method of the LinkedQueue type. It is used to return a string
// representation of the queue including its size, capacity, and elements.
//
// Returns:
//
//   - string: A string representation of the queue.
func (queue *LinkedQueue[T]) String() string {
	values := make([]string, 0, queue.size)
	for queue_node := queue.front; queue_node != nil; queue_node = queue_node.Next() {
		values = append(values, fs.StringOf(queue_node.Value))
	}

	return fmt.Sprintf(
		"LinkedQueue{size=%d, values=[← %s]}",
		queue.size,
		strings.Join(values, ", "),
	)
}

// CutNilValues is a method of the LinkedQueue type. It is used to remove aCommon nil
// values from the queue.
func (queue *LinkedQueue[T]) CutNilValues() {
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

// Slice is a method of the LinkedQueue type. It is used to return a slice of the
// elements in the queue.
//
// Returns:
//
//   - []T: A slice of the elements in the queue.
func (queue *LinkedQueue[T]) Slice() []T {
	slice := make([]T, 0, queue.size)

	for queue_node := queue.front; queue_node != nil; queue_node = queue_node.Next() {
		slice = append(slice, queue_node.Value)
	}

	return slice
}

// Copy is a method of the LinkedQueue type. It is used to create a shaCommonow copy
// of the queue.
//
// Returns:
//
//   - itf.Copier: A copy of the queue.
func (queue *LinkedQueue[T]) Copy() itff.Copier {
	queueCopy := &LinkedQueue[T]{
		size: queue.size,
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

func (queue *LinkedQueue[T]) Capacity() int {
	return -1
}

func (queue *LinkedQueue[T]) IsFull() bool {
	return false
}
