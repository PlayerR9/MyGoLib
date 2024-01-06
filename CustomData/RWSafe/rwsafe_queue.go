package RWSafe

import (
	"sync"
)

// safeElement represents a thread-safe element in a RWSafeQueue.
// It is a generic type that can hold any type T that satisfies the empty
// interface.
// Each safeElement has a value of type T, a pointer to the next safeElement
// in the RWSafeQueue, and a RWMutex for concurrent read/write protection.
type safeElement[T any] struct {
	// The value that this element holds. It is of a generic type T.
	value T

	// A pointer to the next safeElement in the RWSafeQueue. It is nil if this
	// element is the last one in the RWSafeQueue.
	next *safeElement[T]

	// A RWMutex to protect concurrent reads/writes to this element.
	mutex sync.RWMutex
}

// trimFrom removes all elements from the current element onwards in
// the linked list.
// This method is recursive, calling trimFrom on the next element until
// it reaches the end of the list.
// After the recursion, it sets the next element of the current element
// to nil, effectively detaching all subsequent elements.
//
// This method is primarily used for cleaning up the RWSafeQueue and for
// testing purposes.
func (se *safeElement[T]) trimFrom() {
	if se.getNext() != nil {
		se.getNext().trimFrom()
		se.setNext(nil)
	}
}

// getValue safely retrieves the value stored in the safeElement.
// It uses a read lock to ensure thread safety during the operation.
// The read lock is acquired at the start of the function and is
// automatically released when the function returns, regardless of
// how the function exits.
func (se *safeElement[T]) getValue() T {
	se.mutex.RLock()
	defer se.mutex.RUnlock()

	return se.value
}

// setNext safely sets the next element in the safeElement list.
// It acquires a write lock before setting the next element to ensure
// thread safety.
// The write lock is released after the next element is set.
//
// Parameters:
//   - next: The safeElement that should be set as the next element in
//     the list.
func (se *safeElement[T]) setNext(next *safeElement[T]) {
	se.mutex.Lock()
	defer se.mutex.Unlock()

	se.next = next
}

// getNext safely retrieves the next safeElement in the linked list.
// It acquires a read lock before accessing the next element to ensure
// thread safety.
// The read lock is automatically released when the function returns,
// regardless of how the function exits.
//
// Returns:
//   - The next safeElement in the linked list. If this is the last
//     element, it returns nil.
func (se *safeElement[T]) getNext() *safeElement[T] {
	se.mutex.RLock()
	defer se.mutex.RUnlock()

	return se.next
}

// RWSafeQueue is a thread-safe queue that allows multiple
// goroutines to produce and consume elements in a synchronized manner.
// It is implemented as a linked list with pointers to the first and
// last elements.
type RWSafeQueue[T any] struct {
	// A pointer to the first safeElement in the RWSafeQueue. It is nil
	// if the RWSafeQueue is empty.
	first *safeElement[T]

	// A RWMutex to protect concurrent reads/writes to the first element.
	firstMutex sync.RWMutex

	// A pointer to the last safeElement in the RWSafeQueue.
	last *safeElement[T]

	// A RWMutex to protect concurrent reads/writes to the last element.
	lastMutex sync.RWMutex
}

// NewRWSafeQueue initializes a new RWSafeQueue with a given set of
// initial values.
//
// Parameters:
//
//   - values: Zero or more initial values to add to the queue. The
//     values are of type T.
//
// Returns:
//
//   - A pointer to the newly created RWSafeQueue.
func NewRWSafeQueue[T any](values ...T) *RWSafeQueue[T] {
	queue := new(RWSafeQueue[T])

	for _, val := range values {
		queue.Enqueue(val)
	}

	return queue
}

// Enqueue adds a new value to the end of the RWSafeQueue.
//
// Parameters:
//
//   - value: The value to add to the queue. The value is of type T.
func (b *RWSafeQueue[T]) Enqueue(value T) {
	next := safeElement[T]{
		value: value,
	}

	if b.IsEmpty() {
		b.firstMutex.Lock()
		b.first = &next
		b.firstMutex.Unlock()
	} else {
		b.last.setNext(&next)
	}

	b.lastMutex.Lock()
	b.last = &next
	b.lastMutex.Unlock()
}

// Dequeue removes the first value from the RWSafeQueue and returns it.
//
// The method returns the value that was removed from the queue.
//
// If the queue is empty, the method panics with an appropriate message.
func (b *RWSafeQueue[T]) Dequeue() T {
	if b.IsEmpty() {
		panic("Cannot dequeue from an empty queue")
	}

	value := b.first.getValue()

	b.firstMutex.Lock()
	b.first = b.first.getNext()
	b.firstMutex.Unlock()

	if b.IsEmpty() {
		b.lastMutex.Lock()
		b.last = nil
		b.lastMutex.Unlock()
	}

	return value
}

// IsEmpty checks if the RWSafeQueue is empty.
//
// The method returns a boolean indicating whether the queue is empty.
// It returns true if the queue is empty, and false otherwise.
func (b *RWSafeQueue[T]) IsEmpty() bool {
	b.firstMutex.RLock()
	defer b.firstMutex.RUnlock()

	return b.first == nil
}

// Clear removes all values from the RWSafeQueue.
func (b *RWSafeQueue[T]) Clear() {
	if b.IsEmpty() {
		return
	}

	b.firstMutex.Lock()

	if b.first.next != nil {
		b.first.next.trimFrom()
	}
	b.first = nil

	b.firstMutex.Unlock()
}
