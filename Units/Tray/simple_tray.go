package Tray

import (
	"slices"

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

// GetLeftDistance implements the Trayer interface.
func (t *SimpleTray[T]) GetLeftDistance() int {
	return t.arrow
}

// GetRightDistance implements the Trayer interface.
func (t *SimpleTray[T]) GetRightDistance() int {
	return t.size - t.arrow - 1
}

// moveRightBy is a helper function that moves the arrow to the right by n positions.
//
// Parameters:
//   - arrow: The current position of the arrow.
//   - n: The number of positions to move the arrow.
//
// Returns:
//   - int: The updated position of the arrow.
//   - int: The number of positions that the arrow has moved beyond the end of the tape.
func (st *SimpleTray[T]) moveRightBy(arrow, n int) (int, int) {
	if st.size == 0 {
		return arrow, n
	}

	arrow += n
	excess := arrow - (st.size - 1)

	if excess > 0 {
		arrow = st.size - 1
	} else {
		excess = 0
	}

	return arrow, excess
}

// moveLeftBy is a helper function that moves the arrow to the left by n positions.
//
// Parameters:
//   - n: The number of positions to move the arrow.
//
// Returns:
//   - int: The number of positions that the arrow has moved beyond the start of the tape.
//
// Assumptions:
//   - n is non-negative.
func (st *SimpleTray[T]) moveLeftBy(arrow, n int) (int, int) {
	if st.size == 0 {
		return arrow, n
	}

	arrow -= n
	excess := -arrow

	if arrow < 0 {
		arrow = 0
	} else {
		excess = 0
	}

	return arrow, excess
}

// Move implements the Trayer interface.
func (t *SimpleTray[T]) Move(n int) int {
	if n == 0 {
		return 0
	}

	var arrow, excess int

	if n < 0 {
		arrow, excess = t.moveLeftBy(t.arrow, -n)
		excess = -excess
	} else {
		arrow, excess = t.moveRightBy(t.arrow, n)
	}

	t.arrow = arrow

	return excess
}

// Write implements the Trayer interface.
func (t *SimpleTray[T]) Write(elem T) error {
	if len(t.tape) == 0 {
		return ers.NewErrEmpty(t)
	}

	t.tape[t.arrow] = elem

	return nil
}

// Read implements the Trayer interface.
func (t *SimpleTray[T]) Read() (T, error) {
	if len(t.tape) == 0 {
		return *new(T), ers.NewErrEmpty(t)
	}

	return t.tape[t.arrow], nil
}

// ReadMany implements the Trayer interface.
func (t *SimpleTray[T]) ReadMany(n int) []T {
	if n == 0 || t.size == 0 {
		return nil
	}

	var left, right int

	if n < 0 {
		left, _ = t.moveLeftBy(t.arrow, -n)
		right = t.arrow + 1
	} else {
		left = t.arrow
		right, _ = t.moveRightBy(t.arrow, n)
	}

	return t.tape[left:right]
}

// Delete implements the Trayer interface.
func (st *SimpleTray[T]) Delete(n int) int {
	if n == 0 {
		return 0
	} else if st.size == 0 {
		return n
	}

	var left, right, excess int

	if n < 0 {
		left, excess = st.moveLeftBy(st.arrow, -n)
		right = st.arrow + 1

		st.arrow = left
	} else {
		left = st.arrow
		right, excess = st.moveRightBy(st.arrow, n)
	}

	st.tape = slices.Delete(st.tape, left, right)
	st.size = len(st.tape)

	if st.arrow >= st.size {
		if st.size == 0 {
			st.arrow = 0
		} else {
			st.arrow = st.size - 1
		}
	}

	return excess
}

// ExtendTapeOnLeft implements the Trayer interface.
func (t *SimpleTray[T]) ExtendTapeOnLeft(elems ...T) {
	if len(elems) == 0 {
		return
	}

	if t.size == 0 {
		t.tape = elems
		t.arrow = len(elems) - 1
		t.size = len(elems)
	} else {
		t.tape = slices.Insert(t.tape, t.arrow, elems...)
		t.arrow += len(elems)
		t.size = len(t.tape)
	}
}

// ExtendTapeOnRight implements the Trayer interface.
func (t *SimpleTray[T]) ExtendTapeOnRight(elems ...T) {
	if len(elems) == 0 {
		return
	}

	if t.size == 0 {
		t.tape = elems
		t.arrow = 0
		t.size = len(elems)
	} else {
		t.tape = slices.Insert(t.tape, t.arrow+1, elems...)
		t.size = len(t.tape)
	}
}

// ArrowStart moves the arrow to the start of the tape.
func (t *SimpleTray[T]) ArrowStart() {
	t.arrow = 0
}

// ArrowEnd moves the arrow to the end of the tape.
func (t *SimpleTray[T]) ArrowEnd() {
	if len(t.tape) == 0 {
		t.arrow = 0
	} else {
		t.arrow = len(t.tape) - 1
	}
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
		size:  len(tape),
	}
}
