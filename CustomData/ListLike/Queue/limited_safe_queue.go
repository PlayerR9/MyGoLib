package ListLike

import (
	"fmt"
	"strings"
	"sync"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

// LimitedSafeQueue is a generic type in Go that represents a thread-safe queue data
// structure implemented using a linked list.
type LimitedSafeQueue[T any] struct {
	// front and back are pointers to the first and last nodes in the safe queue,
	// respectively.
	front, back *linkedNode[T]

	// frontMutex and backMutex are sync.RWMutexes, which are used to ensure that
	// concurrent reads and writes to the front and back nodes are thread-safe.
	frontMutex, backMutex sync.RWMutex

	// size is the current number of elements in the queue.
	size int

	// capacity is the maximum number of elements that the queue can hold.
	capacity int
}

// NewLimitedSafeQueue is a function that creates and returns a new instance of a LimitedSafeQueue.
// It takes a variadic parameter of type T, which represents the initial values to be
// stored in the queue.
//
// If no initial values are provided, the function simply returns a new LimitedSafeQueue with
// all its fields set to their zero values.
//
// If initial values are provided, the function creates a new LimitedSafeQueue and initializes
// its size. It then creates a linked list of safeNodes from the initial values, with
// each node holding one value, and sets the front and back pointers of the queue.
// The new LimitedSafeQueue is then returned.
func NewLimitedSafeQueue[T any](capacity int, values ...*T) *LimitedSafeQueue[T] {
	if capacity < 0 {
		panic(ers.NewErrInvalidParameter(
			"capacity", fmt.Errorf("negative capacity (%d) is not allowed", capacity),
		))
	} else if len(values) > capacity {
		panic(ers.NewErrInvalidParameter(
			"values", fmt.Errorf("number of values (%d) exceeds the provided capacity (%d)",
				len(values), capacity),
		))
	}

	if len(values) == 0 {
		return new(LimitedSafeQueue[T])
	}

	queue := new(LimitedSafeQueue[T])
	queue.size = len(values)

	// First node
	node := &linkedNode[T]{value: values[0]}

	queue.front = node
	queue.back = node

	// Subsequent nodes
	for _, element := range values[1:] {
		node = &linkedNode[T]{value: element}

		queue.back.next = node
		queue.back = node
	}

	return queue
}

func (queue *LimitedSafeQueue[T]) Enqueue(value *T) {
	queue.backMutex.Lock()
	defer queue.backMutex.Unlock()

	if queue.size == queue.capacity {
		panic(ers.NewErrOperationFailed(
			"enqueue element", NewErrFullQueue(queue),
		))
	}

	node := &linkedNode[T]{value: value}

	if queue.back == nil {
		queue.frontMutex.Lock()
		queue.front = node
		queue.frontMutex.Unlock()
	} else {
		queue.back.next = node
	}

	queue.back = node
	queue.size++
}

func (queue *LimitedSafeQueue[T]) Dequeue() *T {
	queue.frontMutex.Lock()
	defer queue.frontMutex.Unlock()

	if queue.front == nil {
		panic(ers.NewErrOperationFailed(
			"dequeue element", NewErrEmptyQueue(queue),
		))
	}

	value := queue.front.value

	if queue.front.next == nil {
		queue.front = nil

		queue.backMutex.Lock()
		queue.back = nil
		queue.backMutex.Unlock()
	} else {
		queue.front = queue.front.next
	}

	queue.size--

	return value
}

func (queue *LimitedSafeQueue[T]) Peek() *T {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	if queue.front == nil {
		panic(ers.NewErrOperationFailed(
			"peek element", NewErrEmptyQueue(queue),
		))
	}

	return queue.front.value
}

func (queue *LimitedSafeQueue[T]) IsEmpty() bool {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	return queue.front == nil
}

func (queue *LimitedSafeQueue[T]) Size() int {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	queue.backMutex.RLock()
	defer queue.backMutex.RUnlock()

	return queue.size
}

func (queue *LimitedSafeQueue[T]) ToSlice() []*T {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	queue.backMutex.RLock()
	defer queue.backMutex.RUnlock()

	slice := make([]*T, 0, queue.size)

	for node := queue.front; node != nil; node = node.next {
		slice = append(slice, node.value)
	}

	return slice
}

func (queue *LimitedSafeQueue[T]) Clear() {
	queue.frontMutex.Lock()
	defer queue.frontMutex.Unlock()

	if queue.front == nil {
		return // nothing to clear
	}

	queue.backMutex.Lock()
	defer queue.backMutex.Unlock()

	// 1. First node
	prevNode := queue.front
	prevNode.value = nil

	// 2. Subsequent nodes
	for node := queue.front.next; node != nil; node = node.next {
		node.value = nil
		prevNode.next = nil
	}

	// 3. Clear the queue fields
	queue.front = nil
	queue.back = nil
	queue.size = 0
}

func (queue *LimitedSafeQueue[T]) IsFull() bool {
	queue.backMutex.RLock()
	defer queue.backMutex.RUnlock()

	return queue.size >= queue.capacity
}

func (queue *LimitedSafeQueue[T]) String() string {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	queue.backMutex.RLock()
	defer queue.backMutex.RUnlock()

	var builder strings.Builder

	fmt.Fprintf(&builder, "LimitedSafeQueue[size=%d, capacity=%d, values=[← ", queue.capacity, queue.size)

	if !queue.IsEmpty() {
		fmt.Fprintf(&builder, "%v", queue.front.value)

		for node := queue.front.next; node != nil; node = node.next {
			fmt.Fprintf(&builder, " %v", node.value)
		}
	}

	fmt.Fprintf(&builder, "]]")

	return builder.String()
}
