package Tray

import (
	ers "github.com/PlayerR9/MyGoLib/Units/errors"
)

// SimpleTray is a struct that represents a tape.
type SimpleTray[T any] struct {
	// tape is a slice of elements on the tape.
	tape []T

	// arrow is the position of the arrow on the tape.
	arrow int

	// size is the size of the tape.
	size int
}

// Move implements the Trayer interface.
func (t *SimpleTray[T]) Move(n int) int {
	if n == 0 {
		return 0
	}

	var diff int

	if n < 0 {
		diff = t.arrow + n

		if diff < 0 {
			t.arrow = 0
		} else {
			t.arrow += n
			diff = 0
		}
	} else {
		diff = (t.arrow + n) - (t.size - 1)

		if diff > 0 {
			t.arrow = t.size - 1
		} else {
			t.arrow += n
			diff = 0
		}
	}

	return diff
}

// Write writes the given element to the tape at the arrow position.
//
// Parameters:
//   - elem: The element to write.
//
// Returns:
//   - error: An error of type *ers.Empty[*SimpleTray] if the tape is empty.
func (t *SimpleTray[T]) Write(elem T) error {
	if len(t.tape) == 0 {
		return ers.NewErrEmpty(t)
	}

	t.tape[t.arrow] = elem

	return nil
}

// Read reads the element at the arrow position.
//
// Returns:
//   - T: The element at the arrow position.
//   - error: An error of type *ers.Empty[*SimpleTray] if the tape is empty.
func (t *SimpleTray[T]) Read() (T, error) {
	if len(t.tape) == 0 {
		return *new(T), ers.NewErrEmpty(t)
	}

	return t.tape[t.arrow], nil
}

// ReadLeft reads the element to the left of the arrow position.
//
// Returns:
//   - T: The element to the left of the arrow position.
//   - error: An error if elements to the left of the arrow position cannot be found.
//
// Errors:
//   - *ers.ErrEmpty[*GeneralTray]: If the tape is empty.
//   - *ers.ErrInvalidParameter: If n is less than 0.
func (t *SimpleTray[T]) Delete(n int) error {
	if n < 0 {
		return ers.NewErrInvalidParameter(
			"n",
			ers.NewErrGTE(0),
		)
	} else if n == 0 {
		return nil
	} else if len(t.tape) == 0 {
		return ers.NewErrEmpty(t)
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

// Insert inserts the given elements to the tape at the arrow position.
//
// Parameters:
//   - elems: The elements to insert.
func (t *SimpleTray[T]) Insert(elems ...T) {
	if len(t.tape) == 0 {
		t.tape = elems
	} else {
		t.tape = append(t.tape[:t.arrow], append(elems, t.tape[t.arrow:]...)...)
	}
}

// ExtendTapeOnLeft extends the tape on the left with the given elements.
//
// Parameters:
//   - elems: The elements to add.
func (t *SimpleTray[T]) ExtendTapeOnLeft(elems ...T) {
	if len(t.tape) == 0 {
		t.tape = elems
		t.arrow = 0
	} else {
		t.tape = append(elems, t.tape...)
		t.arrow += len(elems)
	}
}

// ExtendTapeOnRight extends the tape on the right with the given elements.
//
// Parameters:
//   - elems: The elements to add.
func (t *SimpleTray[T]) ExtendTapeOnRight(elems ...T) {
	if len(t.tape) == 0 {
		t.tape = elems
		t.arrow = 0
	} else {
		t.tape = append(t.tape, elems...)
	}
}

// ArrowStart moves the arrow to the start of the tape.
func (t *SimpleTray[T]) ArrowStart() {
	t.arrow = 0
}

// ArrowEnd moves the arrow to the end of the tape.
func (t *SimpleTray[T]) ArrowEnd() {
	t.arrow = len(t.tape) - 1
}

/*

// ShiftLeftOfArrow shifts the elements on the right of the
// arrow to the left by n positions.
//
// Parameters:
//   - n: The number of positions to shift the elements.
func (t *GeneralTray[T]) ShiftLeftOfArrow(n int) {
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

// ShiftRightOfArrow shifts the elements on the left of the
// arrow to the right by n positions.
//
// Parameters:
//   - n: The number of positions to shift the elements.
func (t *GeneralTray[T]) ShiftRightOfArrow(n int) {
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
*/

// NewSimpleTray creates a new GeneralTray with the given tape.
//
// Parameters:
//   - tape: The tape to use.
//
// Returns:
//   - *GeneralTray: A pointer to the new GeneralTray.
func NewSimpleTray[T any](tape []T) *SimpleTray[T] {
	return &SimpleTray[T]{
		tape:  tape,
		arrow: 0,
	}
}
