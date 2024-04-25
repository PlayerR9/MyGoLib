package Buffer

import (
	"fmt"
	"sync"

	sll "github.com/PlayerR9/MyGoLib/ListLike/Queue"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// Buffer is a thread-safe, generic data structure that allows multiple
// goroutines to produce and consume elements in a synchronized manner.
// It is implemented as a queue and uses channels to synchronize the
// goroutines.
// The Buffer should be initialized with the Init method before use.
type Buffer[T any] struct {
	// A pointer to the RWSafeQueue that stores the elements of the Buffer.
	q *sll.SafeQueue[T]

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

// NewBuffer creates a new Buffer instance.
//
// Parameters:
//   - bufferSize: The size of the buffer for the send and receive channels.
//     Must be a non-negative integer. If a negative integer is provided,
//     the method will panic with an *ers.InvalidParameterError.
//
// Returns:
//   - *Buffer: A pointer to the newly created Buffer instance.
//   - error: An error of type *ers.InvalidParameterError if
//     the buffer size is negative.
//
// Information: To close the buffer, just close the send-only channel.
// Once that is done, a cascade of events will happen:
//   - The goroutine that listens for incoming messages will stop listening
//     and exit.
//   - The goroutine that sends messages from the Buffer to the receive
//     channel will stop sending messages once the Buffer is empty, and then exit.
//   - The Buffer will be cleaned up.
//
// Of course, a Close method is also provided to manually close the Buffer but
// it is not necessary to call it if the send-only channel is closed.
func NewBuffer[T any](bufferSize int) (*Buffer[T], error) {
	if bufferSize < 0 {
		return nil, ers.NewErrInvalidParameter(
			"bufferSize",
			fmt.Errorf("value (%d) cannot be negative", bufferSize),
		)
	}

	b := &Buffer[T]{
		q:                  sll.NewSafeQueue[T](),
		sendTo:             make(chan T, bufferSize),
		receiveFrom:        make(chan T, bufferSize),
		isClosed:           false,
		isNotEmptyOrClosed: sync.NewCond(new(sync.Mutex)),
	}

	return b, nil
}

// Start is a method of the Buffer type that starts the Buffer by launching
// the goroutines that listen for incoming messages and send messages from
// the Buffer to the send channel.
func (b *Buffer[T]) Start() {
	b.once.Do(func() {
		b.wg.Add(2)

		go b.listenForIncomingMessages()
		go b.sendMessagesFromBuffer()
	})
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
			msg, err := b.q.Dequeue()
			if err == nil {
				b.sendTo <- msg
			}
		}

		b.isNotEmptyOrClosed.L.Unlock()
	}

	for !b.q.IsEmpty() {
		msg, err := b.q.Dequeue()
		if err == nil {
			b.sendTo <- msg
		}
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

// Close is a method of the Buffer type that closes the Buffer.
func (b *Buffer[T]) Close() {
	close(b.receiveFrom)
}
