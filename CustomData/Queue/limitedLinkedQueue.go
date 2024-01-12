package Queue

import (
	"fmt"
	"strings"
)

// LimitedLinkedQueue is a generic type in Go that represents a queue data structure with
// a limited capacity, implemented using a linked list.
type LimitedLinkedQueue[T any] struct {
	// front and back are pointers to the first and last nodes in the linked queue,
	// respectively.
	front, back *linkedNode[T]

	// size is the current number of elements in the queue. capacity is the maximum
	// number of elements the queue can hold.
	size, capacity int
}

func (queue *LimitedLinkedQueue[T]) Cleanup() {
	if queue.front != nil {
		queue.front.Cleanup()
		queue.front = nil
	}

	if queue.back != nil {
		queue.back.Cleanup()
		queue.back = nil
	}
}

// NewLimitedLinkedQueue is a function that creates and returns a new instance of a
// LimitedLinkedQueue.
// It takes an integer capacity, which represents the maximum number of elements the
// queue can hold, and a variadic parameter of type T,
// which represents the initial values to be stored in the queue.
//
// The function first checks if the provided capacity is negative. If it is, it returns
// an error of type ErrNegativeCapacity.
// It then checks if the number of initial values exceeds the provided capacity. If it
// does, it returns an error of type ErrTooManyValues.
//
// If the provided capacity and initial values are valid, the function creates a new
// LimitedLinkedQueue and initializes its size and capacity.
// It then creates a linked list of nodes from the initial values, with each node
// holding one value, and sets the front and back pointers of the queue.
// The new LimitedLinkedQueue is then returned.
func NewLimitedLinkedQueue[T any](capacity int, values ...T) (*LimitedLinkedQueue[T], error) {
	if capacity < 0 {
		return nil, &ErrNegativeCapacity{}
	} else if len(values) > capacity {
		return nil, &ErrTooManyValues{}
	}

	queue := new(LimitedLinkedQueue[T])
	queue.size = len(values)
	queue.capacity = capacity

	// First node
	node := linkedNode[T]{
		value: values[0],
	}

	queue.front = &node
	queue.back = &node

	// Subsequent nodes
	for _, element := range values[1:] {
		node = linkedNode[T]{
			value: element,
		}

		queue.back.next = &node
		queue.back = &node
	}

	return queue, nil
}

// Enqueue is a method of the LimitedLinkedQueue type. It is used to add an element to
// the end of the queue.
//
// The method takes a parameter, value, of a generic type T, which is the element to be
// added to the queue.
//
// Before adding the element, the method checks if the current size of the queue is equal
// to or greater than its capacity.
// If it is, it means the queue is full, and the method panics by throwing an ErrFullQueue
// error.
//
// If the queue is not full, the method creates a new linkedNode with the provided value.
// If the queue is currently empty (i.e., queue.back is nil),
// the new node is set as both the front and back of the queue. If the queue is not empty,
// the new node is added to the end of the queue by setting it
// as the next node of the current back node, and then updating the back pointer of the
// queue to the new node.
//
// Finally, the size of the queue is incremented by 1 to reflect the addition of the
// new element.
func (queue *LimitedLinkedQueue[T]) Enqueue(value T) {
	if queue.size >= queue.capacity {
		panic(&ErrFullQueue{})
	}

	queue_node := linkedNode[T]{
		value: value,
	}

	if queue.back == nil {
		queue.front = &queue_node
	} else {
		queue.back.next = &queue_node
	}

	queue.back = &queue_node

	queue.size++
}

func (queue *LimitedLinkedQueue[T]) Dequeue() (T, error) {
	if queue.front == nil {
		return *new(T), &ErrEmptyQueue{}
	}

	var value T

	value, queue.front = queue.front.value, queue.front.next
	if queue.front == nil {
		queue.back = nil
	}

	queue.size--

	return value, nil
}

func (queue *LimitedLinkedQueue[T]) MustDequeue() T {
	if queue.front == nil {
		panic(ErrEmptyQueue{})
	}

	var value T

	value, queue.front = queue.front.value, queue.front.next
	if queue.front == nil {
		queue.back = nil
	}

	queue.size--

	return value
}

func (queue *LimitedLinkedQueue[T]) Peek() (T, error) {
	if queue.front == nil {
		return *new(T), &ErrEmptyQueue{Peek}
	}

	return queue.front.value, nil
}

func (queue *LimitedLinkedQueue[T]) MustPeek() T {
	if queue.front == nil {
		panic(&ErrEmptyQueue{Peek})
	}

	return queue.front.value
}

func (queue *LimitedLinkedQueue[T]) IsEmpty() bool {
	return queue.front == nil
}

func (queue *LimitedLinkedQueue[T]) Size() int {
	return queue.size
}

func (queue *LimitedLinkedQueue[T]) ToSlice() []T {
	slice := make([]T, 0, queue.size)

	for queue_node := queue.front; queue_node != nil; queue_node = queue_node.next {
		slice = append(slice, queue_node.value)
	}

	return slice
}

func (queue *LimitedLinkedQueue[T]) Clear() {
	queue.front = nil
	queue.back = nil
	queue.size = 0
}

func (queue *LimitedLinkedQueue[T]) IsFull() bool {
	return queue.size >= queue.capacity
}

func (queue *LimitedLinkedQueue[T]) String() string {
	if queue.front == nil {
		return QueueHead
	}

	var builder strings.Builder

	builder.WriteString(QueueHead)
	builder.WriteString(fmt.Sprintf("%v", queue.front.value))

	for queue_node := queue.front.next; queue_node != nil; queue_node = queue_node.next {
		builder.WriteString(QueueSep)
		builder.WriteString(fmt.Sprintf("%v", queue_node.value))
	}

	return builder.String()
}
