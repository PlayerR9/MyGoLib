package RWSafe

import (
	"sync"

	q "github.com/PlayerR9/MyGoLib/CustomData/Queue"
)

// Buffer is a thread-safe, generic data structure that allows multiple
// goroutines to produce and consume elements in a synchronized manner.
// It is implemented as a queue and uses channels to synchronize the
// goroutines.
// The Buffer should be initialized with the Init method before use.
type Buffer[T any] struct {
	// A pointer to the RWSafeQueue that stores the elements of the Buffer.
	queue *q.SafeQueue[T]

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

func (b *Buffer[T]) Cleanup() {
	b.wg.Wait()

	if b.queue != nil {
		b.queue.Cleanup()
		b.queue = nil
	}

	b.sendTo = nil
	b.receiveFrom = nil
	b.isNotEmptyOrClosed = nil
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
		b.queue = q.NewSafeQueue[T]()
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
		b.queue.Enqueue(msg)

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
		for b.queue.IsEmpty() && !b.isClosed {
			b.isNotEmptyOrClosed.Wait()
		}

		if b.isClosed {
			isClosed = true
		} else {
			b.sendTo <- b.queue.MustDequeue()
		}

		b.isNotEmptyOrClosed.L.Unlock()
	}

	for !b.queue.IsEmpty() {
		b.sendTo <- b.queue.MustDequeue()
	}

	close(b.sendTo)
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
	b.queue.Clear()
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
