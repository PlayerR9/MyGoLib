package Queuer

import (
	"strconv"
	"strings"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
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

// Enqueue implements the Queuer interface.
func (queue *LimitedLinkedQueue[T]) Enqueue(value T) bool {
	if queue.size >= queue.capacity {
		return false
	}

	queue_node := NewQueueNode(value)

	if queue.back == nil {
		queue.front = queue_node
	} else {
		queue.back.SetNext(queue_node)
	}

	queue.back = queue_node

	queue.size++

	return true
}

// Dequeue implements the Queuer interface.
func (queue *LimitedLinkedQueue[T]) Dequeue() (T, bool) {
	if queue.front == nil {
		return *new(T), false
	}

	toRemove := queue.front

	queue.front = queue.front.Next()
	if queue.front == nil {
		queue.back = nil
	}

	queue.size--
	toRemove.SetNext(nil)

	return toRemove.Value, true
}

// Peek implements the Queuer interface.
func (queue *LimitedLinkedQueue[T]) Peek() (T, bool) {
	if queue.front == nil {
		return *new(T), false
	}

	return queue.front.Value, true
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
//   - uc.Iterater[T]: An iterator for the queue.
func (queue *LimitedLinkedQueue[T]) Iterator() uc.Iterater[T] {
	var builder uc.Builder[T]

	for queue_node := queue.front; queue_node != nil; queue_node = queue_node.Next() {
		builder.Add(queue_node.Value)
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

// GoString implements the fmt.GoStringer interface.
func (queue *LimitedLinkedQueue[T]) GoString() string {
	values := make([]string, 0, queue.size)
	for queue_node := queue.front; queue_node != nil; queue_node = queue_node.Next() {
		values = append(values, uc.StringOf(queue_node.Value))
	}

	var builder strings.Builder

	builder.WriteString("LimitedLinkedQueue[capacity=")
	builder.WriteString(strconv.Itoa(queue.capacity))
	builder.WriteString(", size=")
	builder.WriteString(strconv.Itoa(queue.size))
	builder.WriteString(", values=[← ")
	builder.WriteString(strings.Join(values, ", "))
	builder.WriteString("]]")

	return builder.String()
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
//   - uc.Copier: A copy of the queue.
func (queue *LimitedLinkedQueue[T]) Copy() uc.Copier {
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
