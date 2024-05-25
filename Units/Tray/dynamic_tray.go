package Tray

import (
	"fmt"
)

// DynamicTray is a struct that represents a tape.
type DynamicTray[E any, T any] struct {
	// source is the source of the procedural tray.
	source Trayer[E]

	// current is the current tape in the source.
	current *SimpleTray[T]

	// transition is the transition function of the dynamic tray.
	transition func(E) *SimpleTray[T]
}

// Move implements the Trayer interface.
func (t *DynamicTray[E, T]) Move(n int) int {
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

		t.current = t.transition(otherTray)

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
func (t *DynamicTray[E, T]) Write(elem T) error {
	if t.current == nil {
		otherTray, err := t.source.Read()
		if err != nil {
			return fmt.Errorf("could not read source: %w", err)
		}

		t.current = t.transition(otherTray)
	}

	return t.current.Write(elem)
}

// Read reads the element at the arrow position.
//
// Returns:
//   - T: The element at the arrow position.
//   - error: An error of type *ers.Empty[*DynamicTray] if the tape is empty.
func (t *DynamicTray[E, T]) Read() (T, error) {
	if t.current == nil {
		otherTray, err := t.source.Read()
		if err != nil {
			return *new(T), fmt.Errorf("could not read source: %w", err)
		}

		t.current = t.transition(otherTray)
	}

	return t.current.Read()
}

// Delete implements the Trayer interface.
//
// However, because it deletes on the dynamic tape, changes will not be saved once
// a new tape is generated.
func (t *DynamicTray[E, T]) Delete(n int) error {
	if t.current == nil {
		otherTray, err := t.source.Read()
		if err != nil {
			return fmt.Errorf("could not read source: %w", err)
		}

		t.current = t.transition(otherTray)
	}

	return t.current.Delete(n)
}

// Insert implements the Trayer interface.
//
// However, because it inserts on the dynamic tape, changes will not be saved once
// a new tape is generated.
func (t *DynamicTray[E, T]) Insert(elems ...T) {
	if t.current == nil {
		otherTray, err := t.source.Read()
		if err != nil {
			panic(fmt.Errorf("could not read source: %w", err))
		}

		t.current = t.transition(otherTray)
	}

	t.current.Insert(elems...)
}

// ExtendTapeOnLeft implements the Trayer interface.
//
// However, because it extends the tape on the dynamic tape, changes will not be saved once
// a new tape is generated.
func (t *DynamicTray[E, T]) ExtendTapeOnLeft(elems ...T) {
	if t.current == nil {
		otherTray, err := t.source.Read()
		if err != nil {
			panic(fmt.Errorf("could not read source: %w", err))
		}

		t.current = t.transition(otherTray)
	}

	t.current.ExtendTapeOnLeft(elems...)
}

// ExtendTapeOnRight implements the Trayer interface.
//
// However, because it extends the tape on the dynamic tape, changes will not be saved once
// a new tape is generated.
func (t *DynamicTray[E, T]) ExtendTapeOnRight(elems ...T) {
	if t.current == nil {
		otherTray, err := t.source.Read()
		if err != nil {
			panic(fmt.Errorf("could not read source: %w", err))
		}

		t.current = t.transition(otherTray)
	}

	t.current.ExtendTapeOnRight(elems...)
}

// ArrowStart implements the Trayer interface.
func (t *DynamicTray[E, T]) ArrowStart() {
	if t.current == nil {
		otherTray, err := t.source.Read()
		if err != nil {
			panic(fmt.Errorf("could not read source: %w", err))
		}

		t.current = t.transition(otherTray)
	}

	t.current.ArrowStart()
}

// ArrowEnd implements the Trayer interface.
func (t *DynamicTray[E, T]) ArrowEnd() {
	if t.current == nil {
		otherTray, err := t.source.Read()
		if err != nil {
			panic(fmt.Errorf("could not read source: %w", err))
		}

		t.current = t.transition(otherTray)
	}

	t.current.ArrowEnd()
}

/*
// ShiftLeftOfArrow shifts the elements on the right of the
// arrow to the left by n positions.
//
// Parameters:
//   - n: The number of positions to shift the elements.
func (t *DynamicTray[E, T]) ShiftLeftOfArrow(n int) {
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
func (t *DynamicTray[E, T]) ShiftRightOfArrow(n int) {
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

// NewDynamicTray creates a new DynamicTray with the given tape.
//
// Parameters:
//   - tape: The tape to use.
//   - f: The transition function.
//
// Returns:
//   - *DynamicTray: A pointer to the new DynamicTray.
func NewDynamicTray[E, T any](tape []E, f func(E) *SimpleTray[T]) *DynamicTray[E, T] {
	return &DynamicTray[E, T]{
		source:     NewSimpleTray(tape),
		current:    nil,
		transition: f,
	}
}
