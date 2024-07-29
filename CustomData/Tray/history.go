package Tray

import (
	ud "github.com/PlayerR9/MyGoLib/Units/Debugging"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// NewSimpleTrayWithHistory creates a new SimpleTray with the given tape and history.
//
// Parameters:
//   - tape: The tape to use.
//
// Returns:
//   - *Debugging.History: A pointer to the new SimpleTray.
func NewSimpleTrayWithHistory[T any](tape []T) *ud.History[*SimpleTray[T]] {
	st := NewSimpleTray(tape)

	h := ud.NewHistory(st)

	return h
}

// MoveCmd is a command that moves the arrow on the tape.
type MoveCmd[T any] struct {
	// n is the number of positions to move the arrow.
	n int

	// excess is the number of positions that the arrow has moved beyond
	// the end of the tape.
	excess int

	// prevArrow is the previous position of the arrow.
	prevArrow int
}

// Execute implements the Debugging.Commander interface.
//
// Never errors.
func (mc *MoveCmd[T]) Execute(data *SimpleTray[T]) error {
	mc.prevArrow = data.arrow

	excess := data.Move(mc.n)

	mc.excess = excess

	return nil
}

// Undo implements the Debugging.Commander interface.
//
// Never errors.
func (mc *MoveCmd[T]) Undo(data *SimpleTray[T]) error {
	data.arrow = mc.prevArrow

	return nil
}

// Copy implements the Debugging.Commander interface.
func (mc *MoveCmd[T]) Copy() uc.Copier {
	mcCopy := &MoveCmd[T]{
		n:         mc.n,
		excess:    mc.excess,
		prevArrow: mc.prevArrow,
	}

	return mcCopy
}

// NewMoveCmd creates a new MoveCmd.
//
// Parameters:
//   - n: The number of positions to move the arrow.
//
// Returns:
//   - *MoveCmd: A pointer to the new MoveCmd.
func NewMoveCmd[T any](n int) *MoveCmd[T] {
	cmd := &MoveCmd[T]{
		n: n,
	}

	return cmd
}

// GetExcess returns the number of positions that the arrow has moved beyond
// the end of the tape.
//
// Call this after the command has been executed.
//
// Returns:
//   - int: The number of positions that the arrow has moved beyond the end of the tape.
func (cm *MoveCmd[T]) GetExcess() int {
	return cm.excess
}

// WriteCmd is a command that writes an element to the tape.
type WriteCmd[T any] struct {
	// elem is the element to write.
	elem T

	// prevElem is the previous element at the arrow position.
	prevElem T
}

// Execute implements the Debugging.Commander interface.
func (wc *WriteCmd[T]) Execute(data *SimpleTray[T]) error {
	prev, err := data.Read()
	if err != nil {
		return err
	}

	wc.prevElem = prev

	err = data.Write(wc.elem)
	if err != nil {
		return err
	}

	return nil
}

// Undo implements the Debugging.Commander interface.
func (wc *WriteCmd[T]) Undo(data *SimpleTray[T]) error {
	err := data.Write(wc.prevElem)
	if err != nil {
		return err
	}

	return nil
}

// Copy implements the Debugging.Commander interface.
func (wc *WriteCmd[T]) Copy() uc.Copier {
	wcCopy := &WriteCmd[T]{
		elem:     wc.elem,
		prevElem: wc.prevElem,
	}

	return wcCopy
}

// NewWriteCmd creates a new WriteCmd.
//
// Parameters:
//   - elem: The element to write.
//
// Returns:
//   - *WriteCmd: A pointer to the new WriteCmd.
func NewWriteCmd[T any](elem T) *WriteCmd[T] {
	cmd := &WriteCmd[T]{
		elem: elem,
	}

	return cmd
}

// DeleteCmd is a command that deletes an element from the tape.
type DeleteCmd[T any] struct {
	// n is the number of elements to delete.
	n int

	// excess is the number of elements that have been deleted beyond
	// the end of the tape.
	excess int

	// prevTape is the previous tape.
	prevElems []T

	// prevArrow is the previous position of the arrow.
	prevArrow int
}

// Execute implements the Debugging.Commander interface.
//
// Never errors.
func (dc *DeleteCmd[T]) Execute(data *SimpleTray[T]) error {
	prevTape := make([]T, len(data.tape))
	copy(prevTape, data.tape)

	dc.prevArrow = data.arrow

	excess := data.Delete(dc.n)

	dc.prevElems = prevTape
	dc.excess = excess

	return nil
}

// Undo implements the Debugging.Commander interface.
//
// Never errors.
func (dc *DeleteCmd[T]) Undo(data *SimpleTray[T]) error {
	data.tape = dc.prevElems
	data.size = len(dc.prevElems)

	data.arrow = dc.prevArrow

	return nil
}

// Copy implements the Debugging.Commander interface.
func (dc *DeleteCmd[T]) Copy() uc.Copier {
	dcCopy := &DeleteCmd[T]{
		n:         dc.n,
		excess:    dc.excess,
		prevElems: dc.prevElems,
		prevArrow: dc.prevArrow,
	}

	return dcCopy
}

// NewDeleteCmd creates a new DeleteCmd.
//
// Parameters:
//   - n: The number of elements to delete.
//
// Returns:
//   - *DeleteCmd: A pointer to the new DeleteCmd.
func NewDeleteCmd[T any](n int) *DeleteCmd[T] {
	cmd := &DeleteCmd[T]{
		n: n,
	}

	return cmd
}

// GetExcess returns the number of elements that have been deleted beyond
// the end of the tape.
//
// Call this after the command has been executed.
//
// Returns:
//   - int: The number of elements that have been deleted beyond the end of the tape.
func (dc *DeleteCmd[T]) GetExcess() int {
	return dc.excess
}

// ExtendLeftCmd is a command that extends the tape on the left.
type ExtendLeftCmd[T any] struct {
	// elems are the elements to extend the tape with.
	elems []T

	// prevTape is the previous tape.
	prevTape []T

	// prevArrow is the previous position of the arrow.
	prevArrow int
}

// Execute implements the Debugging.Commander interface.
//
// Never errors.
func (elc *ExtendLeftCmd[T]) Execute(data *SimpleTray[T]) error {
	prevTape := make([]T, len(data.tape))
	copy(prevTape, data.tape)

	elc.prevArrow = data.arrow

	data.ExtendTapeOnLeft(elc.elems...)

	elc.prevTape = prevTape

	return nil
}

// Undo implements the Debugging.Commander interface.
//
// Never errors.
func (elc *ExtendLeftCmd[T]) Undo(data *SimpleTray[T]) error {
	data.tape = elc.prevTape
	data.size = len(elc.prevTape)

	data.arrow = elc.prevArrow

	return nil
}

// Copy implements the Debugging.Commander interface.
func (elc *ExtendLeftCmd[T]) Copy() uc.Copier {
	elcCopy := &ExtendLeftCmd[T]{
		elems:     elc.elems,
		prevTape:  elc.prevTape,
		prevArrow: elc.prevArrow,
	}

	return elcCopy
}

// NewExtendLeftCmd creates a new ExtendLeftCmd.
//
// Parameters:
//   - elems: The elements to extend the tape with.
//
// Returns:
//   - *ExtendLeftCmd: A pointer to the new ExtendLeftCmd.
func NewExtendLeftCmd[T any](elems ...T) *ExtendLeftCmd[T] {
	cmd := &ExtendLeftCmd[T]{
		elems: elems,
	}

	return cmd
}

// ExtendRightCmd is a command that extends the tape on the right.
type ExtendRightCmd[T any] struct {
	// elems are the elements to extend the tape with.
	elems []T

	// prevTape is the previous tape.
	prevTape []T

	// prevArrow is the previous position of the arrow.
	prevArrow int
}

// Execute implements the Debugging.Commander interface.
//
// Never errors.
func (erc *ExtendRightCmd[T]) Execute(data *SimpleTray[T]) error {
	prevTape := make([]T, len(data.tape))
	copy(prevTape, data.tape)

	erc.prevArrow = data.arrow

	data.ExtendTapeOnRight(erc.elems...)

	erc.prevTape = prevTape

	return nil
}

// Undo implements the Debugging.Commander interface.
//
// Never errors.
func (erc *ExtendRightCmd[T]) Undo(data *SimpleTray[T]) error {
	data.tape = erc.prevTape
	data.size = len(erc.prevTape)

	data.arrow = erc.prevArrow

	return nil
}

// Copy implements the Debugging.Commander interface.
func (erc *ExtendRightCmd[T]) Copy() uc.Copier {
	ercCopy := &ExtendRightCmd[T]{
		elems:     erc.elems,
		prevTape:  erc.prevTape,
		prevArrow: erc.prevArrow,
	}

	return ercCopy
}

// NewExtendRightCmd creates a new ExtendRightCmd.
//
// Parameters:
//   - elems: The elements to extend the tape with.
//
// Returns:
//   - *ExtendRightCmd: A pointer to the new ExtendRightCmd.
func NewExtendRightCmd[T any](elems ...T) *ExtendRightCmd[T] {
	cmd := &ExtendRightCmd[T]{
		elems: elems,
	}

	return cmd
}

// ArrowStartCmd is a command that moves the arrow to the start of the tape.
type ArrowStartCmd[T any] struct {
	// prevArrow is the previous position of the arrow.
	prevArrow int
}

// Execute implements the Debugging.Commander interface.
//
// Never errors.
func (asc *ArrowStartCmd[T]) Execute(data *SimpleTray[T]) error {
	asc.prevArrow = data.arrow

	data.ArrowStart()

	return nil
}

// Undo implements the Debugging.Commander interface.
//
// Never errors.
func (asc *ArrowStartCmd[T]) Undo(data *SimpleTray[T]) error {
	data.arrow = asc.prevArrow

	return nil
}

// Copy implements the Debugging.Commander interface.
func (asc *ArrowStartCmd[T]) Copy() uc.Copier {
	ascCopy := &ArrowStartCmd[T]{
		prevArrow: asc.prevArrow,
	}

	return ascCopy
}

// NewArrowStartCmd creates a new ArrowStartCmd.
//
// Returns:
//   - *ArrowStartCmd: A pointer to the new ArrowStartCmd.
func NewArrowStartCmd[T any]() *ArrowStartCmd[T] {
	cmd := &ArrowStartCmd[T]{}

	return cmd
}

// ArrowEndCmd is a command that moves the arrow to the end of the tape.
type ArrowEndCmd[T any] struct {
	// prevArrow is the previous position of the arrow.
	prevArrow int
}

// Execute implements the Debugging.Commander interface.
//
// Never errors.
func (aec *ArrowEndCmd[T]) Execute(data *SimpleTray[T]) error {
	aec.prevArrow = data.arrow

	data.ArrowEnd()

	return nil
}

// Undo implements the Debugging.Commander interface.
//
// Never errors.
func (aec *ArrowEndCmd[T]) Undo(data *SimpleTray[T]) error {
	data.arrow = aec.prevArrow

	return nil
}

// Copy implements the Debugging.Commander interface.
func (aec *ArrowEndCmd[T]) Copy() uc.Copier {
	aecCopy := &ArrowEndCmd[T]{
		prevArrow: aec.prevArrow,
	}

	return aecCopy
}

// NewArrowEndCmd creates a new ArrowEndCmd.
//
// Returns:
//   - *ArrowEndCmd: A pointer to the new ArrowEndCmd.
func NewArrowEndCmd[T any]() *ArrowEndCmd[T] {
	cmd := &ArrowEndCmd[T]{}

	return cmd
}
