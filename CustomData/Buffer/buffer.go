package Buffer

import "sync"

type safeElement[T interface{}] struct {
	value T
	next  *safeElement[T]

	mutex sync.RWMutex
}

func (se *safeElement[T]) trimLast() {
	if se.getNext().getNext() == nil {
		se.setNext(nil)
	} else {
		se.getNext().trimLast()
	}
}

func (se *safeElement[T]) getValue() T {
	se.mutex.RLock()
	defer se.mutex.RUnlock()

	return se.value
}

func (se *safeElement[T]) setNext(next *safeElement[T]) {
	se.mutex.Lock()
	se.next = next
	se.mutex.Unlock()
}

func (se *safeElement[T]) getNext() *safeElement[T] {
	se.mutex.RLock()
	defer se.mutex.RUnlock()

	return se.next
}

func (b *Buffer[T]) IsEmpty() bool {
	b.firstMutex.RLock()
	defer b.firstMutex.RUnlock()

	return b.first == nil
}

type Buffer[T interface{}] struct {
	first      *safeElement[T]
	firstMutex sync.RWMutex

	last      *safeElement[T]
	lastMutex sync.RWMutex

	channel    chan T
	wg         sync.WaitGroup
	once       sync.Once
	isNotEmpty *sync.Cond
	isClosed   bool
}

func NewBuffer[T interface{}]() Buffer[T] {
	return Buffer[T]{
		first:      nil,
		firstMutex: sync.RWMutex{},
		last:       nil,
		lastMutex:  sync.RWMutex{},
		channel:    make(chan T),
		wg:         sync.WaitGroup{},
		once:       sync.Once{},
		isNotEmpty: sync.NewCond(new(sync.Mutex)),
	}
}

func (b *Buffer[T]) setFirst(element *safeElement[T]) {
	b.firstMutex.Lock()
	defer b.firstMutex.Unlock()

	b.first = element
}

func (b *Buffer[T]) setLast(element *safeElement[T]) {
	b.lastMutex.Lock()
	defer b.lastMutex.Unlock()

	b.last = element
}

func (b *Buffer[T]) Run() {
	b.once.Do(func() {
		b.wg.Add(1)

		for msg := range b.channel {
			next := &safeElement[T]{
				value: msg,
			}

			if b.IsEmpty() {
				b.setFirst(next)
			} else {
				b.lastMutex.Lock()
				b.last.setNext(next)
				b.lastMutex.Unlock()
			}

			b.setLast(next)

			b.isNotEmpty.Signal()
		}
	})
}

func (b *Buffer[T]) Get() (T, bool) {
	b.isNotEmpty.L.Lock()
	defer b.isNotEmpty.L.Unlock()

	for b.IsEmpty() && !b.isClosed {
		b.isNotEmpty.Wait()
	}

	if b.IsEmpty() && b.isClosed {
		b.wg.Done()
		return *new(T), false
	}

	value := b.first.getValue()

	if b.first.getNext() == nil {
		b.setFirst(nil)
		b.setLast(nil)
	} else {
		b.firstMutex.Lock()
		b.first.setNext(b.first.getNext())
		b.firstMutex.Unlock()
	}

	return value, true
}

// DEBUG TOOL: Allows to pause the message box to examine its contents
// Do not use this function outside of debugging
func (b *Buffer[T]) Pause() {
	close(b.channel)

	b.isNotEmpty.L.Lock()
	b.isClosed = true
	b.isNotEmpty.L.Unlock()

	b.wg.Wait()
}

func (b *Buffer[T]) Fini() {
	close(b.channel)

	b.isNotEmpty.L.Lock()
	b.isClosed = true
	b.isNotEmpty.L.Unlock()

	b.wg.Wait()

	// Close and free resources
	b.first = nil
	b.last = nil
	b.channel = nil
}

func (b *Buffer[T]) AbruptFini() {
	close(b.channel)

	// Close everything
	b.isNotEmpty.L.Lock()

	b.firstMutex.Lock()
	b.first.trimLast()
	b.firstMutex.Unlock()

	b.isClosed = true

	b.isNotEmpty.L.Unlock()

	b.wg.Wait()

	// Close and free resources
	b.first = nil
	b.last = nil
	b.channel = nil
}

func (b *Buffer[T]) SendMessages(message T, optionalMessages ...T) {
	b.channel <- message

	for _, message := range optionalMessages {
		b.channel <- message
	}
}
