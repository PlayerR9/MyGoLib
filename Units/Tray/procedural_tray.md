package Tray

import (
	"fmt"
)

// ProceduralTray is a struct that represents a tape.
type ProceduralTray[E Trayable[T], T any] struct {
	// source is the source of the procedural tray.
	source Trayer[E]

	// current is the current tape in the source.
	current Trayer[T]
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
func (st *ProceduralTray[E, T]) moveRightBy(arrow, n int) (int, int) {
	arrow, excess := st.current.moveRightBy(arrow, n)

	if excess > 0 {
		arrow = 0
		excess--

		tmp := st.source.Move(1)
		if tmp != 0 {
			return arrow, excess
		}

		elem, err := st.source.Read()
		if err != nil {
			panic(fmt.Errorf("could not read source: %w", err))
		}

		otherTray := elem.ToTray()

		otherTray.
			excess = otherTray.Move(excess)
	}

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
func (pt *ProceduralTray[E, T]) moveLeftBy(arrow, n int) (int, int) {
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
func (t *ProceduralTray[E, T]) Move(n int) int {
	if n == 0 {
		return 0
	}

	todo := n

	for {
		todo := t.current.Move(todo)
		if todo == 0 {
			return 0
		}

		var didMove bool

		if todo < 0 {
			didMove = t.source.Move(-1) == 0
		} else {
			didMove = t.source.Move(1) == 0
		}

		if !didMove {
			return todo
		}

		otherTray, err := t.source.Read()
		if err != nil {
			panic(fmt.Errorf("could not read source: %w", err))
		}

		var ok bool

		t.current, ok = otherTray.ToTray().(*SimpleTray[T])
		if !ok {
			panic(fmt.Errorf("could not convert to GeneralTray: %w", err))
		}

		if todo < 0 {
			todo++
			t.current.ArrowEnd()
		} else {
			todo--
			t.current.ArrowStart()
		}
	}
}

// Write implements the Trayer interface.
//
// However, because it writes on the dynamic tape, changes will not be saved once
// a new tape is generated.
func (t *ProceduralTray[E, T]) Write(elem T) error {
	if t.current == nil {
		otherTray, err := t.source.Read()
		if err != nil {
			return fmt.Errorf("could not read source: %w", err)
		}

		var ok bool

		t.current, ok = otherTray.ToTray().(*SimpleTray[T])
		if !ok {
			panic(fmt.Errorf("could not convert to GeneralTray: %w", err))
		}
	}

	return t.current.Write(elem)
}

// Read implements the Trayer interface.
func (t *ProceduralTray[E, T]) Read() (T, error) {
	if t.current == nil {
		otherTray, err := t.source.Read()
		if err != nil {
			return *new(T), fmt.Errorf("could not read source: %w", err)
		}

		var ok bool

		t.current, ok = otherTray.ToTray().(*SimpleTray[T])
		if !ok {
			panic(fmt.Errorf("could not convert to GeneralTray: %w", err))
		}
	}

	return t.current.Read()
}

// Delete implements the Trayer interface.
//
// However, because it deletes on the dynamic tape, changes will not be saved once
// a new tape is generated.
func (t *ProceduralTray[E, T]) Delete(n int) error {
	if t.current == nil {
		otherTray, err := t.source.Read()
		if err != nil {
			return fmt.Errorf("could not read source: %w", err)
		}

		var ok bool

		t.current, ok = otherTray.ToTray().(*SimpleTray[T])
		if !ok {
			panic(fmt.Errorf("could not convert to GeneralTray: %w", err))
		}
	}

	return t.current.Delete(n)
}

// Insert implements the Trayer interface.
//
// However, because it inserts on the dynamic tape, changes will not be saved once
// a new tape is generated.
func (t *ProceduralTray[E, T]) Insert(elems ...T) {
	if t.current == nil {
		otherTray, err := t.source.Read()
		if err != nil {
			panic(fmt.Errorf("could not read source: %w", err))
		}

		var ok bool

		t.current, ok = otherTray.ToTray().(*SimpleTray[T])
		if !ok {
			panic(fmt.Errorf("could not convert to GeneralTray: %w", err))
		}
	}

	t.current.Insert(elems...)
}

// ExtendTapeOnLeft implements the Trayer interface.
//
// However, because it extends the tape on the dynamic tape, changes will not be saved once
// a new tape is generated.
func (t *ProceduralTray[E, T]) ExtendTapeOnLeft(elems ...T) {
	if t.current == nil {
		otherTray, err := t.source.Read()
		if err != nil {
			panic(fmt.Errorf("could not read source: %w", err))
		}

		var ok bool

		t.current, ok = otherTray.ToTray().(*SimpleTray[T])
		if !ok {
			panic(fmt.Errorf("could not convert to GeneralTray: %w", err))
		}
	}

	t.current.ExtendTapeOnLeft(elems...)
}

// ExtendTapeOnRight implements the Trayer interface.
//
// However, because it extends the tape on the dynamic tape, changes will not be saved once
// a new tape is generated.
func (t *ProceduralTray[E, T]) ExtendTapeOnRight(elems ...T) {
	if t.current == nil {
		otherTray, err := t.source.Read()
		if err != nil {
			panic(fmt.Errorf("could not read source: %w", err))
		}

		var ok bool

		t.current, ok = otherTray.ToTray().(*SimpleTray[T])
		if !ok {
			panic(fmt.Errorf("could not convert to GeneralTray: %w", err))
		}
	}

	t.current.ExtendTapeOnRight(elems...)
}

// ArrowStart implements the Trayer interface.
func (t *ProceduralTray[E, T]) ArrowStart() {
	if t.current == nil {
		otherTray, err := t.source.Read()
		if err != nil {
			panic(fmt.Errorf("could not read source: %w", err))
		}

		var ok bool

		t.current, ok = otherTray.ToTray().(*SimpleTray[T])
		if !ok {
			panic(fmt.Errorf("could not convert to GeneralTray: %w", err))
		}
	}

	t.current.ArrowStart()
}

// ArrowEnd implements the Trayer interface.
func (t *ProceduralTray[E, T]) ArrowEnd() {
	if t.current == nil {
		otherTray, err := t.source.Read()
		if err != nil {
			panic(fmt.Errorf("could not read source: %w", err))
		}

		var ok bool

		t.current, ok = otherTray.ToTray().(*SimpleTray[T])
		if !ok {
			panic(fmt.Errorf("could not convert to GeneralTray: %w", err))
		}
	}

	t.current.ArrowEnd()
}

/*
// ShiftLeftOfArrow shifts the elements on the right of the
// arrow to the left by n positions.
//
// Parameters:
//   - n: The number of positions to shift the elements.
func (t *ProceduralTray[E, T]) ShiftLeftOfArrow(n int) {
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
func (t *ProceduralTray[E, T]) ShiftRightOfArrow(n int) {
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

// NewProceduralTray creates a new ProceduralTray with the given tape.
//
// Parameters:
//   - tape: The tape to use.
//
// Returns:
//   - *ProceduralTray: A pointer to the new ProceduralTray.
func NewProceduralTray[E Trayable[T], T any](tape []E) *ProceduralTray[E, T] {
	return &ProceduralTray[E, T]{
		source:  NewSimpleTray(tape),
		current: nil,
	}
}
