package Errors

import (
	"fmt"
	"sync"
)

type GRHandler struct {
	wg sync.WaitGroup

	status error
	mu     sync.RWMutex

	id       string
	whenDone func(error)
}

func (h *GRHandler) Wait() error {
	h.wg.Wait()

	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.status
}

type GRHandlerOption func(*GRHandler)

func WithIdentifier(id string) GRHandlerOption {
	return func(h *GRHandler) {
		h.id = id
	}
}

func WhenDone(whenDone func(error)) GRHandlerOption {
	return func(h *GRHandler) {
		h.whenDone = whenDone
	}
}

func WhenPanic(whenPanic func(interface{})) GRHandlerOption {
	return func(h *GRHandler) {
	}
}

func WhenError(whenError func(error)) GRHandlerOption {
	return func(h *GRHandler) {
	}

}

func WhenSuccess(whenSuccess func()) GRHandlerOption {
	return func(h *GRHandler) {
	}
}

func GoStart(id string, whenDone func(error), routine func()) *GRHandler {
	h := &GRHandler{
		status: nil,
		id:     id,
	}

	id
	done
	panic
	err
	success
	routine

	var whenDone func(error)

	var whenSuccess func()
	var whenError func(error)
	var whenPanic func(interface{})

	r := recover()
	if r == nil {
		whenSuccess() || whenDone()

		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	switch x := r.(type) {
	case error:
		h.whenError(x) || h.whenDone(x)
	default:
		h.whenPanic(r) || h.whenDone(&ErrPanic{value: r})
	}

	if whenDone == nil {
		h.whenDone = func(err error) {
			h.status = err
		}
	} else {
		h.whenDone = whenDone
	}

	h.wg.Add(1)

	go func() {
		defer h.wg.Done()

		defer func() {
			r := recover()
			if r == nil {
				whenSuccess()
				whenDone()

				return
			}

			h.mu.Lock()
			defer h.mu.Unlock()

			switch x := r.(type) {
			case error:
				h.whenError(x)
				h.whenDone(x)
			default:
				h.whenPanic(r)
				h.whenDone(&ErrPanic{value: r})
			}
		}()

		routine()
	}()

	return h
}

func (h *GRHandler) PanicIf() {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.status == nil {
		return
	}

	panic(fmt.Errorf("goroutine %s failed: %v", h.id, h.status))
}

func (h *GRHandler) Identifier() string {
	return h.id
}

func (h *GRHandler) Error() error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.status
}

type BatchBuilder struct {
	routines []func()
	ids      []string
	whenDone []func(error)
}

func (b *BatchBuilder) Add(identifier string, whenDone func(error), routine func()) {
	b.routines = append(b.routines, routine)

	b.ids = append(b.ids, identifier)
}

func (b *BatchBuilder) Build() []*GRHandler {
	if b.routines == nil {
		return nil
	}

	pairings := make([]*GRHandler, 0, len(b.routines))

	for i, routine := range b.routines {
		pairings = append(pairings, GoStart(b.ids[i], b.whenDone[i], routine))
	}

	return pairings
}

func WaitAll(batch []*GRHandler) {
	for _, pair := range batch {
		if pair == nil {
			continue
		}

		pair.Wait()
	}
}
