package ListLike

import (
	"fmt"
	"strings"

	itf "github.com/PlayerR9/MyGoLib/CustomData/Iterators"
	ll "github.com/PlayerR9/MyGoLib/ListLike"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
	itff "github.com/PlayerR9/MyGoLibUnits/Interfaces"
)

// LimitedLinkedQueue is a generic type that represents a queue data structure with
// or without a limited capacity, implemented using a linked list.
type LimitedLinkedQueue[T any] struct {
	// front and back are pointers to the first and last nodes in the linked queue,
	// respectively.
	front, back *linkedNode[T]

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
	node := &linkedNode[T]{
		value: values[0],
	}

	queue.front = node
	queue.back = node

	// Subsequent nodes
	for _, element := range values[1:] {
		node = &linkedNode[T]{
			value: element,
		}

		queue.back.next = node
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
		return ll.NewErrFullList(queue)
	}

	queue_node := &linkedNode[T]{
		value: value,
	}

	if queue.back == nil {
		queue.front = queue_node
	} else {
		queue.back.next = queue_node
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
		return *new(T), ll.NewErrEmptyList(queue)
	}

	toRemove := queue.front

	queue.front = queue.front.next
	if queue.front == nil {
		queue.back = nil
	}

	queue.size--
	toRemove.next = nil

	return toRemove.value, nil
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
		return *new(T), ll.NewErrEmptyList(queue)
	}

	return queue.front.value, nil
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

func (queue *LimitedLinkedQueue[T]) Capacity() (int, bool) {
	return queue.capacity, true
}

// Iterator is a method of the LimitedLinkedQueue type. It is used to return an iterator
// for the queue.
//
// Returns:
//
//   - itf.Iterater[T]: An iterator for the queue.
func (queue *LimitedLinkedQueue[T]) Iterator() itf.Iterater[T] {
	var builder itf.Builder[T]

	for queue_node := queue.front; queue_node != nil; queue_node = queue_node.next {
		builder.Append(queue_node.value)
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
	for node := queue.front.next; node != nil; node = node.next {
		prev = node
		prev.next = nil
	}

	prev.next = nil

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
	var builder strings.Builder

	fmt.Fprintf(&builder, "LimitedLinkedQueue[capacity=%d, ", queue.capacity)

	if queue.size == 0 {
		builder.WriteString("size=0, values=[← ]]")

		return builder.String()
	}

	fmt.Fprintf(&builder, "size=%d, values=[← %v", queue.size, queue.front.value)

	for queue_node := queue.front.next; queue_node != nil; queue_node = queue_node.next {
		fmt.Fprintf(&builder, ", %v", queue_node.value)
	}

	builder.WriteString("]]")

	return builder.String()
}

// CutNilValues is a method of the LimitedLinkedQueue type. It is used to remove all nil
// values from the queue.
func (queue *LimitedLinkedQueue[T]) CutNilValues() {
	if queue.front == nil {
		return // Queue is empty
	}

	if gen.IsNil(queue.front.value) && queue.front == queue.back {
		// Single node
		queue.front = nil
		queue.back = nil
		queue.size = 0

		return
	}

	var toDelete *linkedNode[T] = nil

	// 1. First node
	if gen.IsNil(queue.front.value) {
		toDelete = queue.front

		queue.front = queue.front.next

		toDelete.next = nil
		queue.size--
	}

	prev := queue.front

	// 2. Subsequent nodes (except last)
	for node := queue.front.next; node.next != nil; node = node.next {
		if !gen.IsNil(node.value) {
			prev = node
		} else {
			prev.next = node.next
			queue.size--

			if toDelete != nil {
				toDelete.next = nil
			}

			toDelete = node
		}
	}

	if toDelete != nil {
		toDelete.next = nil
	}

	// 3. Last node
	if gen.IsNil(queue.back.value) {
		prev.next = nil
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

	for queue_node := queue.front; queue_node != nil; queue_node = queue_node.next {
		slice = append(slice, queue_node.value)
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
	node := &linkedNode[T]{
		value: queue.front.value,
	}

	queueCopy.front = node
	queueCopy.back = node

	// Subsequent nodes
	for queue_node := queue.front.next; queue_node != nil; queue_node = queue_node.next {
		node = &linkedNode[T]{
			value: queue_node.value,
		}

		queueCopy.back.next = node
		queueCopy.back = node
	}

	return queueCopy
}
