package Buffer

import (
	"sync"
)

// safeElement represents a thread-safe element in a buffer.
// It is a generic type that can hold any type T that satisfies the empty
// interface.
// Each safeElement has a value of type T, a pointer to the next safeElement
// in the buffer,
// and a RWMutex for concurrent read/write protection.
type safeElement[T interface{}] struct {
	// The value that this element holds. It is of a generic type T.
	value T

	// A pointer to the next safeElement in the buffer. It is nil if this
	// element is the last one in the buffer.
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
// This method is primarily used for cleaning up the buffer and for
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
//   - The next safeElement in the linked list. If this is the last element,
//     it returns nil.
func (se *safeElement[T]) getNext() *safeElement[T] {
	se.mutex.RLock()
	defer se.mutex.RUnlock()

	return se.next
}

// Buffer is a thread-safe, generic data structure that allows multiple
// goroutines to produce and consume elements in a synchronized manner.
// It is implemented as a linked list with pointers to the first and last elements,
// and uses channels to synchronize the goroutines.
// The Buffer should be initialized with the Init method before use.
type Buffer[T interface{}] struct {
	// A pointer to the first safeElement in the Buffer. It is nil if the Buffer is empty.
	first *safeElement[T]

	// A RWMutex to protect concurrent reads/writes to the first element.
	firstMutex sync.RWMutex

	// A pointer to the last safeElement in the Buffer. It is nil if the Buffer is empty.
	last *safeElement[T]

	// A RWMutex to protect concurrent reads/writes to the last element.
	lastMutex sync.RWMutex

	// A send-only channel of type T. Messages from the Buffer are sent to this channel.
	sendTo chan T

	// A receive-only channel of type T. Messages sent to this channel are added to the Buffer.
	receiveFrom chan T

	// A WaitGroup to wait for all goroutines to finish.
	wg sync.WaitGroup

	// A Once to ensure that the close and run operations on the channels are done only once..
	once sync.Once

	// A condition variable to signal when the Buffer is not empty or closed.
	isNotEmptyOrClosed *sync.Cond

	// A boolean indicating whether the Buffer is closed.
	isClosed bool
}

// Init initializes a Buffer instance.
// It ensures the initialization is done only once, even if called
// multiple times.
// It creates two channels and a condition variable, and starts two
// goroutines:
// one that listens for incoming messages on the send channel and adds
// them to the Buffer, and another that sends messages from the Buffer
// to the receive channel.
// This method should be called before using the Buffer.
//
// Parameters:
//   - bufferSize: The size of the buffer for the send and receive channels.
//     Must be a non-negative integer. If a negative integer is provided,
//     the method will panic.
//
// Returns:
//   - A send-only channel of type T. Messages sent to this channel are
//     added to the Buffer.
//   - A receive-only channel of type T. Messages from the Buffer are sent
//     to this channel.
//
// Information: To close the buffer, just close the send-only channel.
// Once that is done, a cascade of events will happen:
//   - The goroutine that listens for incoming messages will stop listening
//     and exit.
//   - The goroutine that sends messages from the Buffer to the receive
//     channel will stop sending messages once the Buffer is empty, and then exit.
//   - The Buffer will be cleaned up.
func (b *Buffer[T]) Init(bufferSize int) (chan<- T, <-chan T) {
	if bufferSize < 0 {
		panic("bufferSize cannot be negative")
	}

	b.once.Do(func() {
		b.sendTo = make(chan T, bufferSize)
		b.receiveFrom = make(chan T, bufferSize)
		b.isClosed = false
		b.isNotEmptyOrClosed = sync.NewCond(new(sync.Mutex))

		b.wg.Add(2)

		go b.listenForIncomingMessages()
		go b.sendMessagesFromBuffer()
	})

	return b.receiveFrom, b.sendTo
}

// listenForIncomingMessages is a method of the Buffer type.
// It runs in a separate goroutine and listens for incoming messages
// from the receiveChannel.
// Each incoming message is enqueued in the Buffer.
// If the receiveChannel is closed, the method sets the Buffer's isClosed
// field to true and exits the loop.
// After exiting the loop, it signals the isNotEmptyOrClosed condition
// variable one more time to ensure that any waiting goroutines proceed.
//
// This method is not thread-safe and should only be called from within
// the Buffer's goroutine.
//
// Note: This method assumes that the setFirst, setLast, and isEmpty methods
// of the Buffer type are thread-safe.
func (b *Buffer[T]) listenForIncomingMessages() {
	defer b.wg.Done()

	for msg := range b.receiveFrom {
		// Enqueue the message in the buffer
		next := safeElement[T]{
			value: msg,
		}

		if b.isEmpty() {
			b.setFirst(&next)
		} else {
			b.last.setNext(&next)
		}

		b.setLast(&next)

		b.isNotEmptyOrClosed.Signal()
	}

	b.isNotEmptyOrClosed.L.Lock()
	b.isClosed = true
	b.isNotEmptyOrClosed.L.Unlock()

	b.isNotEmptyOrClosed.Signal()
}

// sendMessagesFromBuffer is a method of the Buffer type.
// It runs in a separate goroutine and sends messages from the
// Buffer to the sendChannel.
// The method waits until the Buffer is not empty or is closed
// before sending each message.
// If the Buffer is closed, the method stops sending messages
// and exits the loop.
// After exiting the loop, it sends any remaining messages in
// the Buffer to the sendChannel, and then closes the sendChannel.
//
// This method is not thread-safe and should only be called from
// within the Buffer's goroutine.
//
// Note: This method assumes that the dequeue and isEmpty
// methods of the Buffer type are thread-safe.
func (b *Buffer[T]) sendMessagesFromBuffer() {
	defer b.wg.Done()

	for isClosed := false; !isClosed; {
		b.isNotEmptyOrClosed.L.Lock()
		for b.isEmpty() && !b.isClosed {
			b.isNotEmptyOrClosed.Wait()
		}

		if b.isClosed {
			isClosed = true
		} else {
			b.sendTo <- b.dequeue()
		}

		b.isNotEmptyOrClosed.L.Unlock()
	}

	for !b.isEmpty() {
		b.sendTo <- b.dequeue()
	}

	close(b.sendTo)
}

// dequeue is a method of the Buffer type.
// It removes and returns the value at the front of the Buffer.
// If the Buffer is empty, it returns the zero value of type T.
// After removing the value, if the Buffer becomes empty, it
// sets both the first and last pointers to nil.
// Otherwise, it updates the first pointer to point to the next
// element in the Buffer.
//
// This method is not thread-safe and should only be called from
// within the Buffer's goroutine.
//
// Note: This method assumes that the getValue, getNext, setFirst,
// and setLast methods of the Buffer type are thread-safe.
func (b *Buffer[T]) dequeue() T {
	value := b.first.getValue()

	b.setFirst(b.first.getNext())

	if b.isEmpty() {
		b.setLast(nil)
	}

	return value
}

// isEmpty safely checks if the Buffer is empty.
// It acquires a read lock before accessing the first element to ensure
// thread safety.
// The read lock is automatically released when the function returns,
// regardless of how the function exits.
//
// Returns:
//   - true if the Buffer is empty (i.e., the first element is nil), and
//     false otherwise.
func (b *Buffer[T]) isEmpty() bool {
	b.firstMutex.RLock()
	defer b.firstMutex.RUnlock()

	return b.first == nil
}

// setFirst safely sets the first element of the Buffer.
// It acquires a write lock before setting the first element to ensure
// thread safety.
// The write lock is released after the first element is set.
//
// Parameters:
//   - element: The safeElement that should be set as the first element
//     in the Buffer.
func (b *Buffer[T]) setFirst(element *safeElement[T]) {
	b.firstMutex.Lock()
	defer b.firstMutex.Unlock()

	b.first = element
}

// setLast safely sets the last element of the Buffer.
// It acquires a write lock before setting the last element to ensure
// thread safety.
// The write lock is released after the last element is set.
//
// Parameters:
//   - element: The safeElement that should be set as the last element
//     in the Buffer.
func (b *Buffer[T]) setLast(element *safeElement[T]) {
	b.lastMutex.Lock()
	defer b.lastMutex.Unlock()

	b.last = element
}

// CleanBuffer removes all elements from the Buffer, effectively resetting
// it to an empty state. Precalculated elements are kept as they are no longer
// in the buffer but in the channel. It locks the firstMutex to ensure
// thread-safety during the operation.
//
// This method is safe for concurrent use by multiple goroutines.
//
// Note: This method assumes that the trimFrom method of the safeElement
// type is thread-safe.
//
// If the Buffer is already empty, this method does nothing.
func (b *Buffer[T]) CleanBuffer() {
	if b.isEmpty() {
		return
	}

	b.firstMutex.Lock()

	if b.first.next != nil {
		b.first.next.trimFrom()
	}
	b.first = nil

	b.firstMutex.Unlock()
}

// Cleanup is a method of the Buffer type.
// It resets the Buffer by clearing its internal state,
// setting the first and last pointers,
// and the receiveChannel and sendChannel to nil.
// This effectively frees all resources used by the Buffer.
//
// This method is not thread-safe and should only be called when
// no other goroutines are accessing the Buffer. It waits until
// all goroutines launched by the Buffer have finished executing
// before proceeding.
func (b *Buffer[T]) Cleanup() {
	b.wg.Wait()

	b.first = nil
	b.last = nil
	b.sendTo = nil
	b.receiveFrom = nil
}

// Wait is a method of the Buffer type.
// It blocks the calling goroutine until all goroutines launched by
// the Buffer have finished executing.
// This is achieved by waiting for the internal WaitGroup to be done.
//
// This method is thread-safe and can be called from multiple goroutines.
func (b *Buffer[T]) Wait() {
	b.wg.Wait()
}
