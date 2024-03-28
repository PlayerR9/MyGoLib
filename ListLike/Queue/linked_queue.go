package ListLike

import (
	"fmt"
	"strings"

	itf "github.com/PlayerR9/MyGoLib/Interfaces"
	ll "github.com/PlayerR9/MyGoLib/ListLike"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
)

// LinkedQueue is a generic type that represents a queue data structure with
// or without a limited capacity, implemented using a linked list.
type LinkedQueue[T any] struct {
	// front and back are pointers to the first and last nodes in the linked queue,
	// respectively.
	front, back *linkedNode[T]

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

// Enqueue is a method of the LinkedQueue type. It is used to add an element to
// the end of the queue.
//
// Panics with an error of type *ErrCallFailed if the queue is full.
//
// Parameters:
//
//   - value: A pointer to a value of type T, which is the element to be added to the
//     queue.
func (queue *LinkedQueue[T]) Enqueue(value T) {
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
}

// Dequeue is a method of the LinkedQueue type. It is used to remove and return
// the element at the front of the queue.
//
// Panics with an error of type *ErrCallFailed if the queue is empty.
//
// Returns:
//
//   - T: The value of the element at the front of the queue.
func (queue *LinkedQueue[T]) Dequeue() (T, error) {
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

// MustDequeue is a method of the LinkedQueue type. It is used to remove and return
// the element at the front of the queue.
//
// Panics with an error of type *ll.ErrEmptyList if the queue is empty.
//
// Returns:
//
//   - T: The value of the element at the front of the queue.
func (queue *LinkedQueue[T]) MustDequeue() T {
	if queue.front == nil {
		panic(ll.NewErrEmptyList(queue))
	}

	toRemove := queue.front

	queue.front = queue.front.next
	if queue.front == nil {
		queue.back = nil
	}

	queue.size--
	toRemove.next = nil

	return toRemove.value
}

// Peek is a method of the LinkedQueue type. It is used to return the element at
// the front of the queue without removing it.
//
// Panics with an error of type *ErrCallFailed if the queue is empty.
//
// Returns:
//
//   - T: The value of the element at the front of the queue.
func (queue *LinkedQueue[T]) Peek() (T, error) {
	if queue.front == nil {
		return *new(T), ll.NewErrEmptyList(queue)
	}

	return queue.front.value, nil
}

// MustPeek is a method of the LinkedQueue type. It is used to return the element at
// the front of the queue without removing it.
//
// Panics with an error of type *ll.ErrEmptyList if the queue is empty.
//
// Returns:
//
//   - T: The value of the element at the front of the queue.
func (queue *LinkedQueue[T]) MustPeek() T {
	if queue.front == nil {
		panic(ll.NewErrEmptyList(queue))
	}

	return queue.front.value
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

	for queue_node := queue.front; queue_node != nil; queue_node = queue_node.next {
		builder.Append(queue_node.value)
	}

	return builder.Build()
}

// Clear is a method of the LinkedQueue type. It is used to remove all elements
// from the queue.
func (queue *LinkedQueue[T]) Clear() {
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

// String is a method of the LinkedQueue type. It is used to return a string
// representation of the queue including its size, capacity, and elements.
//
// Returns:
//
//   - string: A string representation of the queue.
func (queue *LinkedQueue[T]) String() string {
	var builder strings.Builder

	builder.WriteString("LinkedQueue[")

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

// CutNilValues is a method of the LinkedQueue type. It is used to remove all nil
// values from the queue.
func (queue *LinkedQueue[T]) CutNilValues() {
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

// Slice is a method of the LinkedQueue type. It is used to return a slice of the
// elements in the queue.
//
// Returns:
//
//   - []T: A slice of the elements in the queue.
func (queue *LinkedQueue[T]) Slice() []T {
	slice := make([]T, 0, queue.size)

	for queue_node := queue.front; queue_node != nil; queue_node = queue_node.next {
		slice = append(slice, queue_node.value)
	}

	return slice
}

// Copy is a method of the LinkedQueue type. It is used to create a shallow copy
// of the queue.
//
// Returns:
//
//   - itf.Copier: A copy of the queue.
func (queue *LinkedQueue[T]) Copy() itf.Copier {
	queueCopy := &LinkedQueue[T]{
		size: queue.size,
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
