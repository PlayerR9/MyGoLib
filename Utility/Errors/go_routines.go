package Errors

import (
	"fmt"
	"sync"
)

type RecoverBuilder struct {
}

func (b *RecoverBuilder) Build() func() {
	return func() {
		r := recover()
		if r == nil {
			// success
			return
		}

		if x, ok := r.(error); ok {
			// error
		} else {
			// panic
		}
	}
}

func RecoverGeneral() {
	r := recover()
	if r == nil {
		// success
		// 1. Ignore: return //
		// 2. Handle: success()
		// 3. Default: default()
		return
	}

	if x, ok := r.(error); ok {
		// error
		// 1. Forwards: panic(x)
		// 2. Handle: error(x)
		// 3. Default: default() //
	} else {
		// panic:
		// 1. Forwards: panic(r) //
		// 2. Handle: err = &ErrPanic{value: r}
	}
}

// builder -> function

type Y struct {
	status error

	wg sync.WaitGroup

	whenSuccess func()
	whenDone    func(error) // or when error
}

func NewSimple(whenDone func(error)) *Y {
	y := &Y{}

	if whenDone == nil {
		y.whenDone = func(err error) {
			y.status = err
		}
	} else {
		y.whenDone = whenDone
	}

	return y
}

func NewComplex(whenError func(error), whenSuccess func()) *Y {
	y := &Y{}

	if whenError == nil {
		y.whenDone = func(err error) {
			y.status = err
		}
	} else {
		y.whenDone = whenError
	}

	if whenSuccess != nil {
		y.whenSuccess = whenSuccess
	} else {
		y.whenSuccess = func() {
			y.status = nil
		}
	}

	return y
}

func (y *Y) Run(z func(), routine func()) {
	y.wg.Add(1)

	go func() {
		defer y.wg.Done()

		defer func() {
			r := recover()
			if r == nil {
				if y.whenSuccess != nil {
					y.whenSuccess()
				} else {
					y.whenDone(nil)
				}

				return
			}

			// h.mu.Lock()
			// defer h.mu.Unlock()

			var err error

			if x, ok := r.(error); ok {
				err = x
			} else {
				err = &ErrPanic{value: r}
			}

			y.whenDone(err)
		}()

		routine()
	}()
}

type GRHandler interface {
	RunF(func())
}

type SimpleGRH struct {
	wg sync.WaitGroup

	status error
	mu     sync.RWMutex

	id string
}

type ComplexGRH struct {
	wg sync.WaitGroup

	status error
	mu     sync.RWMutex

	id          string
	whenSuccess func()
	whenError   func(error)
}

func (h *GRHandler) Wait() error {
	h.wg.Wait()

	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.status
}

func NewGRHandler(whenDone func(error)) *GRHandler {
	var h *GRHandler

	if whenDone == nil {
		h = &GRHandler{
			whenSuccess: func() {
				h.status = nil
			},
			whenError: func(err error) {
				h.status = err
			},
		}
	} else {
		h = &GRHandler{
			whenSuccess: func() {
				whenDone(nil)
			},
			whenError: whenDone,
		}
	}

	h.wg.Add(1)

	go func() {
		defer h.wg.Done()

		defer func() {
			r := recover()
			if r == nil {
				if whenSuccess != nil {
					whenSuccess()
				} else {
					whenDone(nil)
				}

				return
			}

			h.mu.Lock()
			defer h.mu.Unlock()

			var err error

			if x, ok := r.(error); ok {
				err = x
			} else {
				err = &ErrPanic{value: r}
			}

			if whenError != nil {
				whenError(err)
			} else {
				whenDone(err)
			}
		}()

		routine()
	}()

	return h

	/*

		h := &GRHandler{
			status: nil,
			id:     id,
		}

		/*
			id
			done
			panic
			err
			success
			routine


		var whenSuccess func()
		var whenError func(error)
		var whenPanic func(interface{})

	*/
}

func GoStartWithOptions(whenSuccess func(), whenError func(error), whenPanic func(any)) *GRHandler {
	h := &GRHandler{}

	if whenDone != nil {
		h.whenDone = whenDone
	} else {
		h.whenDone = func(err error) {
			h.status = err
		}
	}

	h.wg.Add(1)

	go func() {
		defer h.wg.Done()

		defer func() {
			r := recover()
			if r == nil {
				if whenSuccess != nil {
					whenSuccess()
				} else {
					whenDone(nil)
				}

				return
			}

			h.mu.Lock()
			defer h.mu.Unlock()

			if whenPanic == nil {
				switch x := r.(type) {
				case error:
					if whenError != nil {
						whenError(x)
					} else {
						whenDone(x)
					}
				default:
					whenDone(&ErrPanic{value: r})
				}
			} else {
				switch x := r.(type) {
				case error:
					if whenError != nil {
						whenError(x)
					} else {
						whenDone(x)
					}
				default:
					whenPanic(r)
				}
			}
		}()

		routine()
	}()

	return h

	/*

		h := &GRHandler{
			status: nil,
			id:     id,
		}

		/*
			id
			done
			panic
			err
			success
			routine


		var whenSuccess func()
		var whenError func(error)
		var whenPanic func(interface{})

	*/
}

func (h *GRHandler) Run(routine func()) *GRHandler {
	h := &GRHandler{}

	if whenDone != nil {
		h.whenDone = whenDone
	} else {
		h.whenDone = func(err error) {
			h.status = err
		}
	}

	h.wg.Add(1)

	go func() {
		defer h.wg.Done()

		defer func() {
			r := recover()
			if r == nil {
				if whenSuccess != nil {
					whenSuccess()
				} else {
					whenDone(nil)
				}

				return
			}

			h.mu.Lock()
			defer h.mu.Unlock()

			if whenPanic == nil {
				switch x := r.(type) {
				case error:
					if whenError != nil {
						whenError(x)
					} else {
						whenDone(x)
					}
				default:
					whenDone(&ErrPanic{value: r})
				}
			} else {
				switch x := r.(type) {
				case error:
					if whenError != nil {
						whenError(x)
					} else {
						whenDone(x)
					}
				default:
					whenPanic(r)
				}
			}
		}()

		routine()
	}()

	return h

	/*

		h := &GRHandler{
			status: nil,
			id:     id,
		}

		/*
			id
			done
			panic
			err
			success
			routine


		var whenSuccess func()
		var whenError func(error)
		var whenPanic func(interface{})

	*/
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
