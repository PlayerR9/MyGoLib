package Queuer

import (
	"strconv"
	"strings"
	"sync"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
)

// LimitedSafeQueue is a generic type that represents a thread-safe queue data
// structure with or without a limited capacity, implemented using a linked list.
type LimitedSafeQueue[T any] struct {
	// front and back are pointers to the first and last nodes in the safe queue,
	// respectively.
	front, back *QueueSafeNode[T]

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

// Enqueue implements the Queuer interface.
func (queue *LimitedSafeQueue[T]) Enqueue(value T) bool {
	queue.backMutex.Lock()
	defer queue.backMutex.Unlock()

	if queue.size >= queue.capacity {
		return false
	}

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

	return true
}

// Dequeue implements the Queuer interface.
func (queue *LimitedSafeQueue[T]) Dequeue() (T, bool) {
	queue.frontMutex.Lock()
	defer queue.frontMutex.Unlock()

	if queue.front == nil {
		return *new(T), false
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

	return toRemove.Value, true
}

// Peek implements the Queuer interface.
func (queue *LimitedSafeQueue[T]) Peek() (T, bool) {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	if queue.front == nil {
		return *new(T), false
	}

	return queue.front.Value, true
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
func (queue *LimitedSafeQueue[T]) Capacity() int {
	return queue.capacity
}

// Iterator is a method of the LimitedSafeQueue type. It is used to return an iterator
// that can be used to iterate over the elements in the queue.
// However, the iterator does not share the queue's thread-safety.
//
// Returns:
//
//   - uc.Iterater[T]: An iterator that can be used to iterate over the elements
//     in the queue.
func (queue *LimitedSafeQueue[T]) Iterator() uc.Iterater[T] {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	queue.backMutex.RLock()
	defer queue.backMutex.RUnlock()

	var builder uc.Builder[T]

	for node := queue.front; node != nil; node = node.Next() {
		builder.Add(node.Value)
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

// GoString implements the fmt.GoStringer interface.
func (queue *LimitedSafeQueue[T]) GoString() string {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	queue.backMutex.RLock()
	defer queue.backMutex.RUnlock()

	values := make([]string, 0, queue.size)
	for node := queue.front; node != nil; node = node.Next() {
		values = append(values, uc.StringOf(node.Value))
	}

	var builder strings.Builder

	builder.WriteString("LimitedSafeQueue[capacity=")
	builder.WriteString(strconv.Itoa(queue.capacity))
	builder.WriteString(", size=")
	builder.WriteString(strconv.Itoa(queue.size))
	builder.WriteString(", values=[← ")
	builder.WriteString(strings.Join(values, ", "))
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

	for node := queue.front; node != nil; node = node.Next() {
		slice = append(slice, node.Value)
	}

	return slice
}

// Copy is a method of the LimitedSafeQueue type. It is used to create a shallow copy of
// the queue.
//
// Returns:
//
//   - uc.Copier: A copy of the queue.
func (queue *LimitedSafeQueue[T]) Copy() uc.Copier {
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
