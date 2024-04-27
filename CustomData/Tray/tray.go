package Tray

import (
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

type Tray[T any] struct {
	tape  []T
	arrow int
}

func NewTray[T any](tape []T) *Tray[T] {
	return &Tray[T]{
		tape:  tape,
		arrow: 0,
	}
}

func (t *Tray[T]) MoveLeft(n int) bool {
	if n == 0 {
		return true
	}

	t.arrow -= n

	if t.arrow < 0 {
		t.arrow = 0
		return false
	} else if t.arrow >= len(t.tape) {
		t.arrow = len(t.tape) - 1
		return false
	}

	return true
}

func (t *Tray[T]) MoveRight(n int) bool {
	if n == 0 {
		return true
	}

	t.arrow += n

	if t.arrow < 0 {
		t.arrow = 0
		return false
	} else if t.arrow >= len(t.tape) {
		t.arrow = len(t.tape) - 1
		return false
	}

	return true
}

func (t *Tray[T]) Write(elem T) error {
	if len(t.tape) == 0 {
		return NewErrEmptyTray()
	}

	t.tape[t.arrow] = elem

	return nil
}

func (t *Tray[T]) Read() (T, error) {
	if len(t.tape) == 0 {
		return *new(T), NewErrEmptyTray()
	}

	return t.tape[t.arrow], nil
}

func (t *Tray[T]) Delete(n int) error {
	if n < 0 {
		return ers.NewErrInvalidParameter(
			"n",
			ers.NewErrGTE(0),
		)
	} else if n == 0 {
		return nil
	} else if len(t.tape) == 0 {
		return NewErrEmptyTray()
	} else if n >= len(t.tape) {
		t.tape = []T{}
		t.arrow = 0
		return nil
	}

	if t.arrow == 0 {
		t.tape = t.tape[n:]
	} else if t.arrow == len(t.tape)-1 {
		t.arrow -= n
		t.tape = t.tape[:t.arrow+1]
	} else {
		t.tape = append(t.tape[:t.arrow], t.tape[t.arrow+n:]...)
	}

	return nil
}

func (t *Tray[T]) Insert(elems ...T) {
	if len(t.tape) == 0 {
		t.tape = elems
	} else {
		t.tape = append(t.tape[:t.arrow], append(elems, t.tape[t.arrow:]...)...)
	}
}

func (t *Tray[T]) ExtendTapeOnLeft(elems ...T) {
	if len(t.tape) == 0 {
		t.tape = elems
		t.arrow = 0
	} else {
		t.tape = append(elems, t.tape...)
		t.arrow += len(elems)
	}
}

func (t *Tray[T]) ExtendTapeOnRight(elems ...T) {
	if len(t.tape) == 0 {
		t.tape = elems
		t.arrow = 0
	} else {
		t.tape = append(t.tape, elems...)
	}
}

func (t *Tray[T]) ArrowStart() {
	t.arrow = 0
}

func (t *Tray[T]) ArrowEnd() {
	t.arrow = len(t.tape) - 1
}

func (t *Tray[T]) ShiftLeftOfArrow(n int) {
	if t.arrow == 0 || len(t.tape) == 0 || n == 0 {
		return
	}

	// FIXME: This is not done yet.

	for i := 0; i < n; i++ {
		// Everything on the right of the arrow is shifted to the left one position;
		// replacing the element at the arrow position.
		t.tape = append(t.tape[:t.arrow], t.tape[t.arrow+1:]...)
	}
}

func (t *Tray[T]) ShiftRightOfArrow(n int) {
	if t.arrow == len(t.tape) || len(t.tape) == 0 || n == 0 {
		return
	}

	// FIXME: This is not done yet.

	for i := 0; i < n; i++ {
		// Everything on the left of the arrow is shifted to the right one position;
		// replacing the element at the arrow position.
		t.tape = append(t.tape[:t.arrow], t.tape[t.arrow+1:]...)
		t.arrow--

	}
}
