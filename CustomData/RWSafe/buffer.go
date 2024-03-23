package RWSafe

import (
	"fmt"
	"sync"

	Queue "github.com/PlayerR9/MyGoLib/ListLike/Queue"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

// Buffer is a thread-safe, generic data structure that allows multiple
// goroutines to produce and consume elements in a synchronized manner.
// It is implemented as a queue and uses channels to synchronize the
// goroutines.
// The Buffer should be initialized with the Init method before use.
type Buffer[T any] struct {
	// A pointer to the RWSafeQueue that stores the elements of the Buffer.
	q *Queue.SafeQueue[T]

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
//
//   - one that listens for incoming messages on the send channel and adds
//     them to the Buffer.
//   - another that sends messages from the Buffer
//     to the receive channel.
//
// This method should be called before using the Buffer.
//
// Parameters:
//   - bufferSize: The size of the buffer for the send and receive channels.
//     Must be a non-negative integer. If a negative integer is provided,
//     the method will panic with an *ers.InvalidParameterError.
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
func (b *Buffer[T]) Init(bufferSize int) error {
	if bufferSize < 0 {
		return ers.NewErrInvalidParameter("bufferSize").Wrap(
			fmt.Errorf("value (%d) cannot be negative", bufferSize),
		)
	}

	b.once.Do(func() {
		b.q = Queue.NewSafeQueue[T]()
		b.sendTo = make(chan T, bufferSize)
		b.receiveFrom = make(chan T, bufferSize)
		b.isClosed = false
		b.isNotEmptyOrClosed = sync.NewCond(new(sync.Mutex))

		b.wg.Add(2)

		go b.listenForIncomingMessages()
		go b.sendMessagesFromBuffer()
	})

	return nil
}

// GetSendChannel returns the send-only channel of the Buffer.
//
// This method is safe for concurrent use by multiple goroutines.
//
// Returns:
//   - chan<- T: The send-only channel of the Buffer.
func (b *Buffer[T]) GetSendChannel() chan<- T {
	return b.sendTo
}

// GetReceiveChannel returns the receive-only channel of the Buffer.
//
// This method is safe for concurrent use by multiple goroutines.
//
// Returns:
//   - <-chan T: The receive-only channel of the Buffer.
func (b *Buffer[T]) GetReceiveChannel() <-chan T {
	return b.receiveFrom
}

// listenForIncomingMessages is a method of the Buffer type that listens for
// incoming messages from the receiveChannel and enqueues them in the Buffer.
//
// It must be run in a separate goroutine to avoid blocking the main thread.
func (b *Buffer[T]) listenForIncomingMessages() {
	defer b.wg.Done()

	for msg := range b.receiveFrom {
		b.q.Enqueue(msg)

		b.isNotEmptyOrClosed.Signal()
	}

	b.isNotEmptyOrClosed.L.Lock()
	b.isClosed = true
	b.isNotEmptyOrClosed.L.Unlock()

	b.isNotEmptyOrClosed.Signal()
}

// sendMessagesFromBuffer is a method of the Buffer type that sends
// messages from the Buffer to the sendChannel.
//
// It must be run in a separate goroutine to avoid blocking the main thread.
func (b *Buffer[T]) sendMessagesFromBuffer() {
	defer b.wg.Done()

	for isClosed := false; !isClosed; {
		b.isNotEmptyOrClosed.L.Lock()
		for b.q.IsEmpty() && !b.isClosed {
			b.isNotEmptyOrClosed.Wait()
		}

		if b.isClosed {
			isClosed = true
		} else {
			msg, _ := b.q.Dequeue()
			b.sendTo <- msg
		}

		b.isNotEmptyOrClosed.L.Unlock()
	}

	for !b.q.IsEmpty() {
		msg, _ := b.q.Dequeue()
		b.sendTo <- msg
	}

	close(b.sendTo)
}

// CleanBuffer removes all elements from the Buffer, effectively resetting
// it to an empty state. Precalculated elements are kept as they are no longer
// in the buffer but in the channel. It locks the firstMutex to ensure
// thread-safety during the operation.
//
// This method is safe for concurrent use by multiple goroutines.
func (b *Buffer[T]) CleanBuffer() {
	b.q.Clear()
}

// Wait is a method of the Buffer type that waits for all goroutines
// launched by the Buffer to finish executing.
//
// This method is thread-safe and can be called from multiple goroutines.
func (b *Buffer[T]) Wait() {
	b.wg.Wait()
}
