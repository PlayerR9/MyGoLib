package Buffer

import (
	"sync"

	"github.com/PlayerR9/MyGoLib/ListLike/Queuer"

	rws "github.com/PlayerR9/MyGoLib/Safe/RWSafe"
	ers "github.com/PlayerR9/MyGoLib/Units/errors"

	uc "github.com/PlayerR9/MyGoLib/Units/Common"
)

// Locker is a thread-safe locker that allows multiple goroutines to wait for a condition.
type Locker[T uc.Enumer] struct {
	// preds is the map of predicates.
	preds *rws.SafeMap[T, bool]

	// obs is the map of observer predicates.
	obs *rws.SafeMap[T, func() bool]

	// cond is the condition variable.
	cond *sync.Cond
}

// NewLocker creates a new Locker.
//
// Use Locker.Set for observer boolean predicates.
//
// Parameters:
//   - keys: The keys to initialize the locker.
//
// Returns:
//   - *Locker[T]: A new Locker.
func NewLocker[T uc.Enumer](keys ...T) *Locker[T] {
	l := &Locker[T]{
		preds: rws.NewSafeMap[T, bool](),
		obs:   rws.NewSafeMap[T, func() bool](),
		cond:  sync.NewCond(&sync.Mutex{}),
	}

	for _, key := range keys {
		l.preds.Set(key, true)
	}

	return l
}

// Set sets the value of a predicate.
//
// Parameters:
//   - key: The key to set the value.
//   - value: The value to set.
func (l *Locker[T]) Set(key T, value func() bool) {
	l.cond.L.Lock()
	defer l.cond.L.Unlock()

	l.obs.Set(key, value)
}

// is checks if any of the predicates are true.
//
// Returns:
//   - map[T]bool: A map of the predicates and their values.
//   - bool: True if all predicates are true, false otherwise.
func (l *Locker[T]) is() (map[T]bool, bool) {
	iter := l.preds.Iterator()

	lockedMap := make(map[T]bool)

	for {
		val, err := iter.Consume()
		if err != nil {
			break
		}

		lockedMap[val.First] = val.Second
	}

	obsIter := l.obs.Iterator()

	for {
		val, err := obsIter.Consume()
		if err != nil {
			break
		}

		lockedMap[val.First] = val.Second()
	}

	for _, value := range lockedMap {
		if value {
			return lockedMap, false
		}
	}

	return lockedMap, true
}

// DoFunc is a function that executes a function while waiting for the condition to be false.
//
// Parameters:
//   - sm: The SafeMap to use.
//
// Returns:
//   - bool: True if the function should exit, false otherwise.
type DoFunc[T uc.Enumer] func(sm map[T]bool) bool

// Do executes a function while waiting for at least one of the conditions to be false.
//
// Parameters:
//   - f: The function to execute.
//
// Returns:
//   - bool: True if the function should exit, false otherwise.
func (l *Locker[T]) Do(f DoFunc[T]) bool {
	l.cond.L.Lock()
	defer l.cond.L.Unlock()

	var m map[T]bool
	var ok bool

	for {
		m, ok = l.is()
		if !ok {
			break
		}

		l.cond.Wait()
	}

	ok = f(m)

	return ok
}

// DoUntill executes a function while waiting for the condition to be false.
//
// The function will be executed until the condition returned by the function is true.
//
// Parameters:
//   - f: The function to execute.
func (l *Locker[T]) DoUntill(f DoFunc[T]) {
	shouldExit := false

	var m map[T]bool
	var ok bool

	for !shouldExit {
		l.cond.L.Lock()

		for {
			m, ok = l.is()
			if !ok {
				break
			}

			l.cond.Wait()
		}

		shouldExit = f(m)

		l.cond.L.Unlock()
	}
}

// Broadcast broadcasts the condition to all waiting goroutines.
//
// Parameters:
//   - key: The key to broadcast.
//   - value: The value to broadcast.
func (l *Locker[T]) Broadcast(key T, value bool) {
	l.cond.L.Lock()
	defer l.cond.L.Unlock()

	if value {
		l.preds.Set(key, func() bool { return true })
	} else {
		l.preds.Set(key, func() bool { return false })
	}

	l.cond.Broadcast()
}

// Signal signals the condition to a single waiting goroutine.
//
// Parameters:
//   - key: The key to signal.
//   - value: The value to signal.
func (l *Locker[T]) Signal(key T, value bool) {
	l.cond.L.Lock()
	defer l.cond.L.Unlock()

	if value {
		l.preds.Set(key, func() bool { return true })
	} else {
		l.preds.Set(key, func() bool { return false })
	}

	l.cond.Signal()
}

// Get returns the value of a predicate.
//
// Parameters:
//   - key: The key to get the value.
//
// Returns:
//   - bool: The value of the predicate or false if the key does not exist.
func (l *Locker[T]) Get(key T) bool {
	val, ok := l.preds.Get(key)
	if !ok {
		return false
	}

	return val()
}

type BufferCondition int

const (
	IsEmpty BufferCondition = iota
	IsNotClosed
)

func (bc BufferCondition) String() string {
	return [...]string{"IsEmpty", "IsNotClosed"}[bc]
}

// Buffer is a thread-safe, generic data structure that allows multiple
// goroutines to produce and consume elements in a synchronized manner.
// It is implemented as a queue and uses channels to synchronize the
// goroutines.
// The Buffer should be initialized with the Init method before use.
type Buffer[T any] struct {
	// A pointer to the RWSafeQueue that stores the elements of the Buffer.
	q *Queuer.SafeQueue[T]

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

	ineoc *Locker[BufferCondition]
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
			ers.NewErrGTE(0),
		)
	}

	b := &Buffer[T]{
		q:                  Queuer.NewSafeQueue[T](),
		sendTo:             make(chan T, bufferSize),
		receiveFrom:        make(chan T, bufferSize),
		ineoc:              NewLocker(IsNotClosed),
		isNotEmptyOrClosed: sync.NewCond(new(sync.Mutex)),
	}

	b.ineoc.Set(IsEmpty, func() bool { return b.q.IsEmpty() })

	return b, nil
}

// Start is a method of the Buffer type that starts the Buffer by launching
// the goroutines that listen for incoming messages and send messages from
// the Buffer to the send channel.
func (b *Buffer[T]) Start() {
	b.once.Do(func() {
		b.wg.Add(2)

		go b.sendMessagesFromBuffer()
		go b.listenForIncomingMessages()

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

	for msg := range b.sendTo {
		b.q.Enqueue(msg)
		b.isNotEmptyOrClosed.Broadcast()
	}

	b.ineoc.Broadcast(IsNotClosed, false)
}

// sendMessagesFromBuffer is a method of the Buffer type that sends
// messages from the Buffer to the sendChannel.
//
// It must be run in a separate goroutine to avoid blocking the main thread.
func (b *Buffer[T]) sendMessagesFromBuffer() {
	defer b.wg.Done()

	for {
		b.isNotEmptyOrClosed.L.Lock()
		for !b.isClosed.Get() && b.q.IsEmpty() {
			b.isNotEmptyOrClosed.Wait()
		}

		msg, err := b.q.Dequeue()

		if err == nil {
			b.receiveFrom <- msg
		}

		if b.isClosed.Get() {
			b.isNotEmptyOrClosed.L.Unlock()
			break
		}

		b.isNotEmptyOrClosed.L.Unlock()
	}

	for {
		msg, err := b.q.Dequeue()
		if err != nil {
			break
		}

		b.receiveFrom <- msg
	}

	close(b.receiveFrom)
}

// CleanBuffer removes all elements from the Buffer, effectively resetting
// it to an empty state. Precalculated elements are kept as they are no longer
// in the buffer but in the channel. It locks the firstMutex to ensure
// thread-safety during the operation.
//
// This method is safe for concurrent use by multiple goroutines.
func (b *Buffer[T]) CleanBuffer() {
	b.q.Clear()

	b.isNotEmptyOrClosed.Broadcast()
}

// Wait is a method of the Buffer type that waits for all goroutines
// launched by the Buffer to finish executing.
//
// This method is thread-safe and can be called from multiple goroutines.
func (b *Buffer[T]) Wait() {
	b.wg.Wait()
}

// Close is a method of the Buffer type that closes the Buffer
// and waits for all goroutines to finish executing.
func (b *Buffer[T]) Close() {
	close(b.sendTo)

	b.wg.Wait()
}

// Send is a method of the Buffer type that sends a message to the send channel.
//
// Parameters:
//   - msg: The message to send.
//
// Behaviors:
//   - If the send channel is nil, the method will return immediately.
func (b *Buffer[T]) Send(msg T) {
	if b.sendTo == nil {
		return
	}

	b.sendTo <- msg
}

// Receive is a method of the Buffer type that receives a message from the receive channel.
//
// Returns:
//   - T: The message received from the receive channel.
//   - bool: A boolean indicating if the message was received successfully.
//
// Behaviors:
//   - If the receive channel is nil, the method will return a zero value and false.
//   - This method will block until a message is received from the receive channel.
func (b *Buffer[T]) Receive() (T, bool) {
	if b.receiveFrom == nil {
		return *new(T), false
	}

	msg, ok := <-b.receiveFrom
	if !ok {
		return msg, false
	}

	return msg, true
}
