package ListLike

import (
	"fmt"
	"strings"
	"sync"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

// SafeQueue is a generic type in Go that represents a thread-safe queue data
// structure implemented using a linked list.
type SafeQueue[T any] struct {
	// front and back are pointers to the first and last nodes in the safe queue,
	// respectively.
	front, back *linkedNode[T]

	// frontMutex and backMutex are sync.RWMutexes, which are used to ensure that
	// concurrent reads and writes to the front and back nodes are thread-safe.
	frontMutex, backMutex sync.RWMutex

	// size is the current number of elements in the queue.
	size int
}

// NewSafeQueue is a function that creates and returns a new instance of a SafeQueue.
// It takes a variadic parameter of type T, which represents the initial values to be
// stored in the queue.
//
// If no initial values are provided, the function simply returns a new SafeQueue with
// all its fields set to their zero values.
//
// If initial values are provided, the function creates a new SafeQueue and initializes
// its size. It then creates a linked list of safeNodes from the initial values, with
// each node holding one value, and sets the front and back pointers of the queue.
// The new SafeQueue is then returned.
func NewSafeQueue[T any](values ...*T) *SafeQueue[T] {
	if len(values) == 0 {
		return new(SafeQueue[T])
	}

	queue := new(SafeQueue[T])
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

func (queue *SafeQueue[T]) Enqueue(value *T) {
	queue.backMutex.Lock()
	defer queue.backMutex.Unlock()

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

func (queue *SafeQueue[T]) Dequeue() *T {
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

func (queue *SafeQueue[T]) Peek() *T {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	if queue.front == nil {
		panic(ers.NewErrOperationFailed(
			"peek element", NewErrEmptyQueue(queue),
		))
	}

	return queue.front.value
}

func (queue *SafeQueue[T]) IsEmpty() bool {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	return queue.front == nil
}

func (queue *SafeQueue[T]) Size() int {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	queue.backMutex.RLock()
	defer queue.backMutex.RUnlock()

	return queue.size
}

func (queue *SafeQueue[T]) ToSlice() []*T {
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

func (queue *SafeQueue[T]) Clear() {
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

func (queue *SafeQueue[T]) IsFull() bool {
	return false
}

func (queue *SafeQueue[T]) String() string {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	queue.backMutex.RLock()
	defer queue.backMutex.RUnlock()

	var builder strings.Builder

	fmt.Fprintf(&builder, "SafeQueue[size=%d, values=[← ", queue.size)

	if !queue.IsEmpty() {
		fmt.Fprintf(&builder, "%v", queue.front.value)

		for node := queue.front.next; node != nil; node = node.next {
			fmt.Fprintf(&builder, " %v", node.value)
		}
	}

	fmt.Fprintf(&builder, "]]")

	return builder.String()
}
