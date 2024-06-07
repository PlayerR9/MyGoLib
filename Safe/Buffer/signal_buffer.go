package Buffer

import (
	"sync"

	rws "github.com/PlayerR9/MyGoLib/Safe/RWSafe"
)

// SignalBuffer is a thread-safe buffer that can be used to signal changes
// to a channel. It is safe for concurrent use by multiple goroutines.
type SignalBuffer struct {
	// The value contained in the SignalBuffer.
	value *rws.Subject[bool]

	// A receive-only channel of type T. Messages from the SignalBuffer
	// are sent to this channel.
	receiveFrom chan bool

	// A WaitGroup to wait for all goroutines to finish.
	wg sync.WaitGroup

	// closed is a channel that is closed when the SignalBuffer is closed.
	closed chan struct{}
}

// NewSignalBuffer creates a new SignalBuffer instance.
//
// Returns:
//   - *SignalBuffer: A pointer to the new SignalBuffer.
func NewSignalBuffer() *SignalBuffer {
	return &SignalBuffer{
		value:       rws.NewSubject(false),
		receiveFrom: make(chan bool),
		closed:      make(chan struct{}),
	}

}

// Start starts the SignalBuffer by launching one goroutine that sends messages
// from the SignalBuffer to the send channel.
func (b *SignalBuffer) Start() {
	b.wg.Add(1)
	go b.sendMessagesFromSignalBuffer()
}

// GetSignalChan returns the send-only channel of the SignalBuffer.
//
// This method is safe for concurrent use by multiple goroutines.
//
// Returns:
//   - <-chan bool: The receive-only channel of the SignalBuffer.
func (b *SignalBuffer) GetSignalChan() <-chan bool {
	return b.receiveFrom
}

// SignalChange signals that the SignalBuffer has changed and that the send
// goroutine should send a message to the send channel.
//
// This method is safe for concurrent use by multiple goroutines.
func (b *SignalBuffer) SignalChange() {
	b.value.Set(true)

	select {
	case b.receiveFrom <- true:
	default:
	}
}

// sendMessagesFromSignalBuffer is a method of the SignalBuffer type that sends
// messages from the SignalBuffer to the sendChannel.
//
// It must be run in a separate goroutine to avoid blocking the main thread.
func (b *SignalBuffer) sendMessagesFromSignalBuffer() {
	defer b.wg.Done()

	for {
		select {
		case <-b.closed:
			return
		default:
			if b.value.Get() {
				b.receiveFrom <- true
				b.value.Set(false)
			}
		}
	}
}

// CleanSignalBuffer removes all elements from the SignalBuffer, effectively resetting
// it to an empty state. Precalculated elements are kept as they are no longer
// in the buffer but in the channel. It locks the firstMutex to ensure
// thread-safety during the operation.
//
// This method is safe for concurrent use by multiple goroutines.
func (b *SignalBuffer) CleanSignalBuffer() {
	b.value.Set(false)
}

// Wait is a method of the SignalBuffer type that waits for all goroutines
// launched by the SignalBuffer to finish executing.
//
// This method is thread-safe and can be called from multiple goroutines.
func (b *SignalBuffer) Wait() {
	b.wg.Wait()
}

// Close is a method of the SignalBuffer type that closes the SignalBuffer.
func (b *SignalBuffer) Close() {
	close(b.closed)

	b.wg.Wait()

	close(b.receiveFrom)
}
