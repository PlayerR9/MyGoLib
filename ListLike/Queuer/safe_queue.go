package Queuer

import (
	"fmt"
	"strings"
	"sync"

	itf "github.com/PlayerR9/MyGoLib/ListLike/Iterator"
	uc "github.com/PlayerR9/MyGoLib/Units/Common"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
)

// SafeQueue is a generic type that represents a thread-safe queue data
// structure with or without a limited capacity, implemented using a linked list.
type SafeQueue[T any] struct {
	// front and back are pointers to the first and last nodes in the safe queue,
	// respectively.
	front, back *QueueSafeNode[T]

	// frontMutex and backMutex are sync.RWMutexes, which are used to ensure that
	// concurrent reads and writes to the front and back nodes are thread-safe.
	frontMutex, backMutex sync.RWMutex

	// size is the current number of elements in the queue.
	size int
}

// NewSafeQueue is a function that creates and returns a new instance of a
// SafeQueue.
//
// Parameters:
//
//   - values: A variadic parameter of type T, which represents the initial values to be
//     stored in the queue.
//
// Return:
//
//   - *SafeQueue[T]: A pointer to the newly created SafeQueue.
func NewSafeQueue[T any](values ...T) *SafeQueue[T] {
	if len(values) == 0 {
		return new(SafeQueue[T])
	}

	queue := new(SafeQueue[T])
	queue.size = len(values)

	// First node
	node := NewQueueSafeNode(values[0])

	queue.front = node
	queue.back = node

	// Subsequent nodes
	for _, element := range values[1:] {
		node = NewQueueSafeNode(element)

		queue.back.SetNext(node)
		queue.back = node
	}

	return queue
}

// Enqueue is a method of the SafeQueue type. It is used to add an element to the
// back of the queue.
//
// Panics with an error of type *ErrCallFailed if the queue is fu
//
// Parameters:
//
//   - value: The value of type T to be added to the queue.
func (queue *SafeQueue[T]) Enqueue(value T) error {
	queue.backMutex.Lock()
	defer queue.backMutex.Unlock()

	node := NewQueueSafeNode(value)

	if queue.back == nil {
		queue.frontMutex.Lock()
		queue.front = node
		queue.frontMutex.Unlock()
	} else {
		queue.back.SetNext(node)
	}

	queue.back = node
	queue.size++

	return nil
}

// Dequeue is a method of the SafeQueue type. It is used to remove and return the
// element at the front of the queue.
//
// Panics with an error of type *ErrCallFailed if the queue is empty.
//
// Returns:
//
//   - T: The value of the element at the front of the queue.
func (queue *SafeQueue[T]) Dequeue() (T, error) {
	queue.frontMutex.Lock()
	defer queue.frontMutex.Unlock()

	if queue.front == nil {
		return *new(T), NewErrEmptyList(queue)
	}

	toRemove := queue.front

	if queue.front.Next() == nil {
		queue.front = nil

		queue.backMutex.Lock()
		queue.back = nil
		queue.backMutex.Unlock()
	} else {
		queue.front = queue.front.Next()
	}

	queue.size--
	toRemove.SetNext(nil)

	return toRemove.Value, nil
}

// Peek is a method of the SafeQueue type. It is used to return the element at the
// front of the queue without removing it.
//
// Panics with an error of type *ErrCallFailed if the queue is empty.
//
// Returns:
//
//   - T: The value of the element at the front of the queue.
func (queue *SafeQueue[T]) Peek() (T, error) {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	if queue.front == nil {
		return *new(T), NewErrEmptyList(queue)
	}

	return queue.front.Value, nil
}

// IsEmpty is a method of the SafeQueue type. It is used to check if the queue is
// empty.
//
// Returns:
//
//   - bool: A boolean value that is true if the queue is empty, and false otherwise.
func (queue *SafeQueue[T]) IsEmpty() bool {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	return queue.front == nil
}

// Size is a method of the SafeQueue type. It is used to return the number of
// elements in the queue.
//
// Returns:
//
//   - int: An integer that represents the number of elements in the queue.
func (queue *SafeQueue[T]) Size() int {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	queue.backMutex.RLock()
	defer queue.backMutex.RUnlock()

	return queue.size
}

// Iterator is a method of the SafeQueue type. It is used to return an iterator
// that can be used to iterate over the elements in the queue.
// However, the iterator does not share the queue's thread-safety.
//
// Returns:
//
//   - itf.Iterater[T]: An iterator that can be used to iterate over the elements
//     in the queue.
func (queue *SafeQueue[T]) Iterator() itf.Iterater[T] {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	queue.backMutex.RLock()
	defer queue.backMutex.RUnlock()

	var builder itf.Builder[T]

	for node := queue.front; node != nil; node = node.Next() {
		builder.Append(node.Value)
	}

	return builder.Build()
}

// Clear is a method of the SafeQueue type. It is used to remove all elements
// from the queue, making it empty.
func (queue *SafeQueue[T]) Clear() {
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

// String is a method of the SafeQueue type. It returns a string representation of
// the queue, including its size, capacity, and the elements it contains.
func (queue *SafeQueue[T]) String() string {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	queue.backMutex.RLock()
	defer queue.backMutex.RUnlock()

	values := make([]string, 0, queue.size)
	for node := queue.front; node != nil; node = node.Next() {
		values = append(values, uc.StringOf(node.Value))
	}

	return fmt.Sprintf(
		"SafeQueue{size=%d, values=[← %s]}",
		queue.size,
		strings.Join(values, ", "),
	)
}

// CutNilValues is a method of the SafeQueue type. It is used to remove all nil
// values from the queue.
func (queue *SafeQueue[T]) CutNilValues() {
	queue.frontMutex.Lock()
	defer queue.frontMutex.Unlock()

	queue.backMutex.Lock()
	defer queue.backMutex.Unlock()

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

	var toDelete *QueueSafeNode[T] = nil

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
		queue.back = prev
		queue.back.SetNext(nil)
		queue.size--
	}
}

// Slice is a method of the SafeQueue type. It is used to return a slice of the
// elements in the queue.
//
// Returns:
//
//   - []T: A slice of the elements in the queue.
func (queue *SafeQueue[T]) Slice() []T {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	queue.backMutex.RLock()
	defer queue.backMutex.RUnlock()

	slice := make([]T, 0, queue.size)

	for node := queue.front; node != nil; node = node.Next() {
		slice = append(slice, node.Value)
	}

	return slice
}

// Copy is a method of the SafeQueue type. It is used to create a shallow copy of
// the queue.
//
// Returns:
//
//   - itf.Copier: A copy of the queue.
func (queue *SafeQueue[T]) Copy() uc.Copier {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	queue.backMutex.RLock()
	defer queue.backMutex.RUnlock()

	queueCopy := &SafeQueue[T]{
		size: queue.size,
	}

	if queue.front == nil {
		return queueCopy
	}

	// First node
	node := NewQueueSafeNode(queue.front.Value)

	queueCopy.front = node
	queueCopy.back = node

	// Subsequent nodes
	for qNode := queue.front.Next(); qNode != nil; qNode = qNode.Next() {
		node = NewQueueSafeNode(qNode.Value)

		queueCopy.back.SetNext(node)
		queueCopy.back = node
	}

	return queueCopy
}

func (queue *SafeQueue[T]) Capacity() int {
	return -1
}

func (queue *SafeQueue[T]) IsFull() bool {
	return false
}
