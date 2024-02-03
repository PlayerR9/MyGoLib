package ListLike

import (
	"fmt"
	"strings"
	"sync"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

// safeNode is a generic type in Go that represents a node in a thread-safe linked
// list.
type safeNode[T any] struct {
	// value is the value stored in the node.
	value *T

	// next is a pointer to the next node in the linked list.
	next *safeNode[T]

	// mutex is a sync.RWMutex, which is used to ensure that concurrent reads and
	// writes to the node are thread-safe.
	mutex sync.RWMutex
}

// newSafeNode is a function that creates and returns a new instance of a safeNode.
// It takes a parameter, value, of a generic type T, which represents the value to
// be stored in the node.
// The function creates a new safeNode, initializes its value field with the provided
// value, and returns it.
// The next field of the node is left as nil, and the mutex field is left as its zero
// value.
func newSafeNode[T any](value *T) safeNode[T] {
	return safeNode[T]{
		value: value,
	}
}

// setNext is a method of the safeNode type. It is used to set the next node in the
// linked list.
//
// The method takes a parameter, next, which is a pointer to the safeNode that should
// be set as the next node.
//
// The method first acquires a lock on the node's mutex to ensure that no other
// goroutine can write to the node while it is being modified.
// It then sets the next field of the node to the provided next node. The lock is
// released when the method returns, either normally or due to a panic, thanks to the
// defer statement.
func (node *safeNode[T]) setNext(next *safeNode[T]) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	node.next = next
}

// getValue is a method of the safeNode type. It is used to get the value stored in
// the node.
//
// The method first acquires a read lock on the node's mutex to ensure that no other
// goroutine can write to the node while its value is being read.
// It then returns the value stored in the node. The read lock is released when the
// method returns, either normally or due to a panic, thanks to the defer statement.
func (node *safeNode[T]) getValue() *T {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	return node.value
}

// getNext is a method of the safeNode type. It is used to get the next node in the
// linked list.
//
// The method first acquires a read lock on the node's mutex to ensure that no other
// goroutine can write to the node while its next field is being read.
// It then returns the next node in the linked list. The read lock is released when
// the method returns, either normally or due to a panic, thanks to the defer statement.
func (node *safeNode[T]) getNext() *safeNode[T] {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	return node.next
}

// SafeQueue is a generic type in Go that represents a thread-safe queue data
// structure implemented using a linked list.
type SafeQueue[T any] struct {
	// front and back are pointers to the first and last nodes in the safe queue,
	// respectively.
	front, back *safeNode[T]

	// frontMutex and backMutex are sync.RWMutexes, which are used to ensure that
	// concurrent reads and writes to the front and back nodes are thread-safe.
	frontMutex, backMutex sync.RWMutex

	// size is the current number of elements in the queue.
	size int

	// sizeMutex is a sync.RWMutex, which is used to ensure that concurrent reads
	// and writes to the size field are thread-safe.
	sizeMutex sync.RWMutex
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
	queue.setSize(len(values))

	// First node
	node := newSafeNode(values[0])

	queue.setFront(&node)
	queue.setBack(&node)

	// Subsequent nodes
	for _, element := range values[1:] {
		node = newSafeNode(element)

		queue.back.setNext(&node)
		queue.setBack(&node)
	}

	return queue
}

func (queue *SafeQueue[T]) Enqueue(value *T) {
	node := newSafeNode(value)

	if queue.IsEmpty() {
		queue.setFront(&node)
	} else {
		queue.back.setNext(&node)
	}

	queue.setBack(&node)

	queue.setSize(queue.Size() + 1)
}

func (queue *SafeQueue[T]) Dequeue() *T {
	if queue.IsEmpty() {
		panic(ers.NewErrOperationFailed(
			"dequeue", NewErrEmptyQueue(queue),
		))
	}

	value := queue.front.getValue()
	queue.front = queue.front.getNext()

	if queue.IsEmpty() {
		queue.setBack(nil)
	}

	queue.setSize(queue.Size() - 1)

	return value
}

func (queue *SafeQueue[T]) Peek() *T {
	if queue.IsEmpty() {
		panic(ers.NewErrOperationFailed(
			"peek", NewErrEmptyQueue(queue),
		))
	}

	return queue.front.getValue()
}

func (queue *SafeQueue[T]) IsEmpty() bool {
	queue.frontMutex.RLock()
	defer queue.frontMutex.RUnlock()

	return queue.front == nil
}

func (queue *SafeQueue[T]) Size() int {
	queue.sizeMutex.RLock()
	defer queue.sizeMutex.RUnlock()

	return queue.size
}

func (queue *SafeQueue[T]) ToSlice() []*T {
	slice := make([]*T, 0, queue.Size())

	queue.frontMutex.RLock()
	node := queue.front
	queue.frontMutex.RUnlock()

	for node != nil {
		slice = append(slice, node.getValue())

		node = node.getNext()
	}

	return slice
}

func (queue *SafeQueue[T]) Clear() {
	queue.setFront(nil)
	queue.setBack(nil)
	queue.setSize(0)
}

func (queue *SafeQueue[T]) IsFull() bool {
	return false
}

func (queue *SafeQueue[T]) String() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "SafeQueue[size=%d, values=[← ", queue.Size())

	if !queue.IsEmpty() {
		fmt.Fprintf(&builder, "%v", queue.front.getValue())

		for node := queue.front.getNext(); node != nil; node = node.getNext() {
			fmt.Fprintf(&builder, " %v", node.getValue())
		}
	}

	fmt.Fprintf(&builder, "]]")

	return builder.String()
}

// setFront is a method of the SafeQueue type. It is used to set the front node
// of the queue.
//
// The method takes a parameter, node, which is a pointer to the safeNode that
// should be set as the front node.
//
// The method first acquires a lock on the queue's frontMutex to ensure that no
// other goroutine can write to the front node while it is being modified.
// It then sets the front field of the queue to the provided node. The lock is
// released when the method returns, either normally or due to a panic,
// thanks to the defer statement.
func (queue *SafeQueue[T]) setFront(node *safeNode[T]) {
	queue.frontMutex.Lock()
	defer queue.frontMutex.Unlock()

	queue.front = node
}

// setBack is a method of the SafeQueue type. It is used to set the back node of
// the queue.
//
// The method takes a parameter, node, which is a pointer to the safeNode that
// should be set as the back node.
//
// The method first acquires a lock on the queue's backMutex to ensure that no
// other goroutine can write to the back node while it is being modified.
// It then sets the back field of the queue to the provided node. The lock is
// released when the method returns, either normally or due to a panic,
// thanks to the defer statement.
func (queue *SafeQueue[T]) setBack(node *safeNode[T]) {
	queue.backMutex.Lock()
	defer queue.backMutex.Unlock()

	queue.back = node
}

// setSize is a method of the SafeQueue type. It is used to set the size of the
// queue.
//
// The method takes a parameter, size, which is an integer that represents the
// new size of the queue.
//
// The method first acquires a lock on the queue's sizeMutex to ensure that no
// other goroutine can write to the size field while it is being modified.
// It then sets the size field of the queue to the provided size. The lock is
// released when the method returns, either normally or due to a panic,
// thanks to the defer statement.
func (queue *SafeQueue[T]) setSize(size int) {
	queue.sizeMutex.Lock()
	defer queue.sizeMutex.Unlock()

	queue.size = size
}
