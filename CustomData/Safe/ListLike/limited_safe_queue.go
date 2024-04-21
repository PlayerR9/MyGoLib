package ListLike

import (
	"fmt"
	"strings"
	"sync"

	itf "github.com/PlayerR9/MyGoLib/CustomData/Iterators"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
	itff "github.com/PlayerR9/MyGoLibUnits/Interfaces"
)

// LimitedSafeQueue is a generic type that represents a thread-safe queue data
// structure with or without a limited capacity, implemented using a linked list.
type LimitedSafeQueue[T any] struct {
	// front and back are pointers to the first and last nodes in the safe queue,
	// respectively.
	front, back *queueLinkedNode[T]

	// frontMutex and backMutex are sync.RWMutexes, which are used to ensure that
	// concurrent reads and writes to the front and back nodes are thread-safe.
	frontMutex, backMutex sync.RWMutex

	// size is the current number of elements in the queue.
	size int

	// capacity is the maximum number of elements that the queue can hold.
	capacity int
}

// NewLimitedSafeQueue is a function that creates and returns a new instance of a
// LimitedSafeQueue.
//
// Parameters:
//
//   - values: A variadic parameter of type T, which represents the initial values to be
//     stored in the queue.
//
// Return:
//
//   - *LimitedSafeQueue[T]: A pointer to the newly created LimitedSafeQueue.
func NewLimitedSafeQueue[T any](values ...T) *LimitedSafeQueue[T] {
	if len(values) == 0 {
		return new(LimitedSafeQueue[T])
	}

	queue := new(LimitedSafeQueue[T])
	queue.size = len(values)

	// First node
	node := &queueLinkedNode[T]{value: values[0]}

	queue.front = node
	queue.back = node

	// Subsequent nodes
	for _, element := range values[1:] {
		node = &queueLinkedNode[T]{value: element}

		queue.back.next = node
		queue.back = node
	}

	return queue
}

// Enqueue is a method of the LimitedSafeQueue type. It is used to add an element to the
// back of the queue.
//
// Panics with an error of type *ErrCallFailed if the queue is fu
//
// Parameters:
//
//   - value: The value of type T to be added to the queue.
func (queue *LimitedSafeQueue[T]) Enqueue(value T) error {
	queue.backMutex.Lock()
	defer queue.backMutex.Unlock()

	if queue.size >= queue.capacity {
		return NewErrFullList(queue)
	}

	node := &queueLinkedNode[T]{value: value}

	if queue.back == nil {
		queue.frontMutex.Lock()
		queue.front = node
		queue.frontMutex.Unlock()
	} else {
		queue.back.next = node
	}

	queue.back = node
	queue.size++

	return nil
}

// Dequeue is a method of the LimitedSafeQueue type. It is used to remove and return the
// element at the front of the queue.
//
// Panics with an error of type *ErrCallFailed if the queue is empty.
//
// Returns:
//
//   - T: The value of the element at the front of the queue.
func (queue *LimitedSafeQueue[T]) Dequeue() (T, error) {
	queue.frontMutex.Lock()
	defer queue.frontMutex.Unlock()

	if queue.front == nil {
		return *new(T), NewErrEmptyList(queue)
	}

	toRemove := queue.front

	if queue.front.next == nil {
		queue.front = nil

		queue.backMutex.Lock()
		queue.back = nil
		queue.backMutex.Unlock()
	} else {
		queue.front = queue.front.next
	}

	queue.size--
	toRemove.next = nil

	return toRemove.value, nil
}

// Peek is a method of the LimitedSafeQueue type. It is used to return the element at the
// front of the queue without removing it.
//
// Panics with an error of type *ErrCallFailed if the queue is empty.
//
// Returns:
//
//   - T: The value of the element at the front of the queue.
func (queue *LimitedSafeQueue[T]) Peek() (T, error) {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	if queue.front == nil {
		return *new(T), NewErrEmptyList(queue)
	}

	return queue.front.value, nil
}

// IsEmpty is a method of the LimitedSafeQueue type. It is used to check if the queue is
// empty.
//
// Returns:
//
//   - bool: A boolean value that is true if the queue is empty, and false otherwise.
func (queue *LimitedSafeQueue[T]) IsEmpty() bool {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	return queue.front == nil
}

// Size is a method of the LimitedSafeQueue type. It is used to return the number of
// elements in the queue.
//
// Returns:
//
//   - int: An integer that represents the number of elements in the queue.
func (queue *LimitedSafeQueue[T]) Size() int {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	queue.backMutex.RLock()
	defer queue.backMutex.RUnlock()

	return queue.size
}

// Capacity is a method of the LimitedSafeQueue type. It is used to return the maximum
// number of elements the queue can hold.
//
// Returns:
//
//   - optional.Int: An optional integer that represents the maximum number of
//     elements the queue can hold.
func (queue *LimitedSafeQueue[T]) Capacity() (int, bool) {
	return queue.capacity, true
}

// Iterator is a method of the LimitedSafeQueue type. It is used to return an iterator
// that can be used to iterate over the elements in the queue.
// However, the iterator does not share the queue's thread-safety.
//
// Returns:
//
//   - itf.Iterater[T]: An iterator that can be used to iterate over the elements
//     in the queue.
func (queue *LimitedSafeQueue[T]) Iterator() itf.Iterater[T] {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	queue.backMutex.RLock()
	defer queue.backMutex.RUnlock()

	var builder itf.Builder[T]

	for node := queue.front; node != nil; node = node.next {
		builder.Append(node.value)
	}

	return builder.Build()
}

// Clear is a method of the LimitedSafeQueue type. It is used to remove all elements
// from the queue, making it empty.
func (queue *LimitedSafeQueue[T]) Clear() {
	queue.frontMutex.Lock()
	defer queue.frontMutex.Unlock()

	queue.backMutex.Lock()
	defer queue.backMutex.Unlock()

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

// IsFull is a method of the LimitedSafeQueue type. It is used to check if the queue is
// full, meaning it has reached its maximum capacity and cannot accept any more
// elements.
//
// Returns:
//
//   - isFull: A boolean value that is true if the queue is full, and false otherwise.
func (queue *LimitedSafeQueue[T]) IsFull() (isFull bool) {
	queue.backMutex.RLock()
	defer queue.backMutex.RUnlock()

	return queue.size >= queue.capacity
}

// String is a method of the LimitedSafeQueue type. It returns a string representation of
// the queue, including its size, capacity, and the elements it contains.
func (queue *LimitedSafeQueue[T]) String() string {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	queue.backMutex.RLock()
	defer queue.backMutex.RUnlock()

	var builder strings.Builder

	fmt.Fprintf(&builder, "LimitedSafeQueue[capacity=%d, ", queue.capacity)

	if queue.size == 0 {
		builder.WriteString("size=0, values=[← ]]")

		return builder.String()
	}

	fmt.Fprintf(&builder, "size=%d, values=[← %v", queue.size, queue.front.value)

	for node := queue.front.next; node != nil; node = node.next {
		fmt.Fprintf(&builder, ", %v", node.value)
	}

	builder.WriteString("]]")

	return builder.String()
}

// CutNilValues is a method of the LimitedSafeQueue type. It is used to remove all nil
// values from the queue.
func (queue *LimitedSafeQueue[T]) CutNilValues() {
	queue.frontMutex.Lock()
	defer queue.frontMutex.Unlock()

	queue.backMutex.Lock()
	defer queue.backMutex.Unlock()

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

	var toDelete *queueLinkedNode[T] = nil

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
		queue.back = prev
		queue.back.next = nil
		queue.size--
	}
}

// Slice is a method of the LimitedSafeQueue type. It is used to return a slice of the
// elements in the queue.
//
// Returns:
//
//   - []T: A slice of the elements in the queue.
func (queue *LimitedSafeQueue[T]) Slice() []T {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	queue.backMutex.RLock()
	defer queue.backMutex.RUnlock()

	slice := make([]T, 0, queue.size)

	for node := queue.front; node != nil; node = node.next {
		slice = append(slice, node.value)
	}

	return slice
}

// Copy is a method of the LimitedSafeQueue type. It is used to create a shallow copy of
// the queue.
//
// Returns:
//
//   - itf.Copier: A copy of the queue.
func (queue *LimitedSafeQueue[T]) Copy() itff.Copier {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	queue.backMutex.RLock()
	defer queue.backMutex.RUnlock()

	queueCopy := &LimitedSafeQueue[T]{
		size: queue.size,
	}

	if queue.front == nil {
		return queueCopy
	}

	// First node
	node := &queueLinkedNode[T]{value: queue.front.value}

	queueCopy.front = node
	queueCopy.back = node

	// Subsequent nodes
	for qNode := queue.front.next; qNode != nil; qNode = qNode.next {
		node = &queueLinkedNode[T]{value: qNode.value}

		queueCopy.back.next = node
		queueCopy.back = node
	}

	return queueCopy
}
