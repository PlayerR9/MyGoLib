package Queuer

import (
	"strconv"
	"strings"
	"sync"

	rws "github.com/PlayerR9/MyGoLib/Safe/RWSafe"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
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
	mu sync.RWMutex

	// size is the size that observers observe.
	size *rws.Subject[int]
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
		return &SafeQueue[T]{
			size: rws.NewSubject(0),
		}
	}

	queue := &SafeQueue[T]{
		size: rws.NewSubject(len(values)),
	}

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
//
// Always returns true.
func (queue *SafeQueue[T]) Enqueue(value T) bool {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	node := NewQueueSafeNode(value)

	if queue.back == nil {
		queue.front = node
	} else {
		queue.back.SetNext(node)
	}

	queue.back = node

	queue.size.ModifyState(func(size int) int {
		return size + 1
	})

	return true
}

// Dequeue implements the Queuer interface.
func (queue *SafeQueue[T]) Dequeue() (T, bool) {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	if queue.front == nil {
		return *new(T), false
	}

	toRemove := queue.front

	if queue.front.Next() == nil {
		queue.front = nil
		queue.back = nil
	} else {
		queue.front = queue.front.Next()
	}

	queue.size.ModifyState(func(size int) int {
		return size - 1
	})

	return toRemove.Value, true
}

// Peek implements the Queuer interface.
func (queue *SafeQueue[T]) Peek() (T, bool) {
	queue.mu.RLock()
	defer queue.mu.RUnlock()

	if queue.front == nil {
		return *new(T), false
	}

	return queue.front.Value, true
}

// IsEmpty is a method of the SafeQueue type. It is used to check if the queue is
// empty.
//
// Returns:
//
//   - bool: A boolean value that is true if the queue is empty, and false otherwise.
func (queue *SafeQueue[T]) IsEmpty() bool {
	queue.mu.RLock()
	defer queue.mu.RUnlock()

	return queue.front == nil
}

// Size is a method of the SafeQueue type. It is used to return the number of
// elements in the queue.
//
// Returns:
//
//   - int: An integer that represents the number of elements in the queue.
func (queue *SafeQueue[T]) Size() int {
	queue.mu.RLock()
	defer queue.mu.RUnlock()

	return queue.size.Get()
}

// Iterator is a method of the SafeQueue type. It is used to return an iterator
// that can be used to iterate over the elements in the queue.
// However, the iterator does not share the queue's thread-safety.
//
// Returns:
//
//   - uc.Iterater[T]: An iterator that can be used to iterate over the elements
//     in the queue.
func (queue *SafeQueue[T]) Iterator() uc.Iterater[T] {
	queue.mu.RLock()
	defer queue.mu.RUnlock()

	var builder uc.Builder[T]

	for node := queue.front; node != nil; node = node.Next() {
		builder.Add(node.Value)
	}

	return builder.Build()
}

// Clear is a method of the SafeQueue type. It is used to remove all elements
// from the queue, making it empty.
func (queue *SafeQueue[T]) Clear() {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	if queue.front == nil {
		return // Queue is already empty
	}

	queue.front = nil
	queue.back = nil

	queue.size.Set(0)
}

// GoString implements the fmt.GoStringer interface.
func (queue *SafeQueue[T]) GoString() string {
	queue.mu.RLock()
	defer queue.mu.RUnlock()

	size := queue.size.Get()

	values := make([]string, 0, size)
	for node := queue.front; node != nil; node = node.Next() {
		values = append(values, uc.StringOf(node.Value))
	}

	var builder strings.Builder

	builder.WriteString("SafeQueue{size=")
	builder.WriteString(strconv.Itoa(size))
	builder.WriteString(", values=[← ")
	builder.WriteString(strings.Join(values, ", "))
	builder.WriteString("]}")

	return builder.String()
}

// CutNilValues is a method of the SafeQueue type. It is used to remove all nil
// values from the queue.
func (queue *SafeQueue[T]) CutNilValues() {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	if queue.front == nil {
		return // Queue is empty
	}

	if gen.IsNil(queue.front.Value) && queue.front == queue.back {
		// Single node
		queue.front = nil
		queue.back = nil

		queue.size.Set(0)

		return
	}

	var toDelete *QueueSafeNode[T] = nil

	// 1. First node
	if gen.IsNil(queue.front.Value) {
		toDelete = queue.front

		queue.front = queue.front.Next()

		toDelete.SetNext(nil)

		queue.size.ModifyState(func(size int) int {
			return size - 1
		})

		if queue.front == nil {
			queue.back = nil
		}
	}

	prev := queue.front

	// 2. Subsequent nodes (except last)
	for node := queue.front.Next(); node.Next() != nil; node = node.Next() {
		if !gen.IsNil(node.Value) {
			prev = node
		} else {
			prev.SetNext(node.Next())

			queue.size.ModifyState(func(size int) int {
				return size - 1
			})

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

		queue.size.ModifyState(func(size int) int {
			return size - 1
		})
	}
}

// Slice is a method of the SafeQueue type. It is used to return a slice of the
// elements in the queue.
//
// Returns:
//
//   - []T: A slice of the elements in the queue.
func (queue *SafeQueue[T]) Slice() []T {
	queue.mu.RLock()
	defer queue.mu.RUnlock()

	slice := make([]T, 0, queue.size.Get())

	for node := queue.front; node != nil; node = node.Next() {
		slice = append(slice, node.Value)
	}

	return slice
}

// Copy is a method of the SafeQueue type. It is used to create a shallow copy of
// the queue.
//
// Returns:
//   - uc.Copier: A copy of the queue.
//
// Behaviors:
//   - Does not copy the observers.
func (queue *SafeQueue[T]) Copy() uc.Copier {
	queue.mu.RLock()
	defer queue.mu.RUnlock()

	queueCopy := &SafeQueue[T]{
		size: rws.NewSubject(queue.size.Get()),
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

// Capacity is a method of the SafeQueue type. It is used to return the maximum
// number of elements that the queue can store.
//
// Returns:
//   - int: -1
func (queue *SafeQueue[T]) Capacity() int {
	return -1
}

// IsFull is a method of the SafeQueue type. It is used to check if the queue is
// full.
//
// Returns:
//   - bool: false
func (queue *SafeQueue[T]) IsFull() bool {
	return false
}

// SetIsEmptyObserver is a method of the SafeQueue type. It is used to set an
// observer that will be notified when the queue becomes empty or non-empty.
func (queue *SafeQueue[T]) ObserveSize(action func(int)) {
	queue.size.Attach(NewIsEmptyObserver(action))
}

type IsEmptyObserver struct {
	action func(int)
}

func (o *IsEmptyObserver) Notify(size int) {
	o.action(size)
}

func NewIsEmptyObserver(action func(int)) *IsEmptyObserver {
	if action == nil {
		return nil
	}

	return &IsEmptyObserver{
		action: action,
	}
}
