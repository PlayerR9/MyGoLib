package Stacker

import (
	"fmt"

	ud "github.com/PlayerR9/MyGoLib/Units/Debugging"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// Clear clears the stack.
type ClearCmd[T any] struct {
	// data is a backup of the stack.
	data Stacker[T]
}

// Execute implements the Debugging.Commander interface.
//
// Never errors.
func (c *ClearCmd[T]) Execute(data Stacker[T]) error {
	// Save the data before clearing it
	c.data = data.Copy().(Stacker[T])

	data.Clear()

	return nil
}

// Undo implements the Debugging.Commander interface.
func (c *ClearCmd[T]) Undo(data Stacker[T]) error {
	data.Clear()

	// Restore the data
	values := c.data.Slice()

	for _, val := range values {
		ok := data.Push(val)
		if !ok {
			return fmt.Errorf("could not push value %v", val)
		}
	}

	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *ClearCmd[T]) Copy() uc.Copier {
	cCopy := &ClearCmd[T]{
		data: c.data.Copy().(Stacker[T]),
	}

	return cCopy
}

// NewClear is a function that creates a new ClearCmd.
//
// Returns:
//   - *ClearCmd: A pointer to the new ClearCmd.
func NewClear[T any]() *ClearCmd[T] {
	cmd := &ClearCmd[T]{}
	return cmd
}

// PushCmd is a command that pushes a value onto the stack.
type PushCmd[T any] struct {
	// value is the value to push onto the stack.
	value T
}

// Execute implements the Debugging.Commander interface.
func (c *PushCmd[T]) Execute(data Stacker[T]) error {
	ok := data.Push(c.value)
	if !ok {
		return fmt.Errorf("could not push value %v", c.value)
	}

	return nil
}

// Undo implements the Debugging.Commander interface.
func (c *PushCmd[T]) Undo(data Stacker[T]) error {
	_, ok := data.Pop()
	if !ok {
		return fmt.Errorf("could not pop value %v", c.value)
	}

	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *PushCmd[T]) Copy() uc.Copier {
	cCopy := &PushCmd[T]{
		value: uc.CopyOf(c.value).(T),
	}

	return cCopy
}

// Push pushes a value onto the stack.
//
// Parameters:
//   - value: The value to push onto the stack.
//
// Returns:
//   - *PushCmd: A pointer to the new PushCmd.
func NewPush[T any](value T) *PushCmd[T] {
	cmd := &PushCmd[T]{
		value: value,
	}
	return cmd
}

// PopCmd is a command that pops a value from the stack.
type PopCmd[T any] struct {
	// value is the value that was popped from the stack.
	value T
}

// Execute implements the Debugging.Commander interface.
func (c *PopCmd[T]) Execute(data Stacker[T]) error {
	val, ok := data.Pop()
	if !ok {
		return fmt.Errorf("could not pop value %v", c.value)
	}

	c.value = val

	return nil
}

// Undo implements the Debugging.Commander interface.
func (c *PopCmd[T]) Undo(data Stacker[T]) error {
	ok := data.Push(c.value)
	if !ok {
		return fmt.Errorf("could not push value %v", c.value)
	}

	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *PopCmd[T]) Copy() uc.Copier {
	cCopy := &PopCmd[T]{
		value: uc.CopyOf(c.value).(T),
	}

	return cCopy
}

// NewPop is a function that creates a new PopCmd.
//
// Returns:
//   - *PopCmd: A pointer to the new PopCmd.
func NewPop[T any]() *PopCmd[T] {
	cmd := &PopCmd[T]{}
	return cmd
}

// Value is a function that returns the value that was popped
// from the stack.
//
// Returns:
//   - T: The value that was popped from the stack.
//
// Must be called after the command has been executed.
func (c *PopCmd[T]) Value() T {
	return c.value
}

// CutNilValuesCmd is a command that removes nil values from the
// stack.
type CutNilValuesCmd[T any] struct {
	// data is a backup of the stack.
	data Stacker[T]
}

// Execute implements the Debugging.Commander interface.
//
// Never errors.
func (c *CutNilValuesCmd[T]) Execute(data Stacker[T]) error {
	c.data = data.Copy().(Stacker[T])

	data.CutNilValues()

	return nil
}

// Undo implements the Debugging.Commander interface.
func (c *CutNilValuesCmd[T]) Undo(data Stacker[T]) error {
	data.Clear()

	values := c.data.Slice()

	for _, val := range values {
		ok := data.Push(val)
		if !ok {
			return fmt.Errorf("could not push value %v", val)
		}
	}

	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *CutNilValuesCmd[T]) Copy() uc.Copier {
	cCopy := &CutNilValuesCmd[T]{
		data: c.data.Copy().(Stacker[T]),
	}

	return cCopy
}

// CutNilValues is a method that removes aCommon nil values from the stack.
//
// It also removes any empty or nil elements in the auxiliary stack.
func NewCutNilValues[T any]() *CutNilValuesCmd[T] {
	cmd := &CutNilValuesCmd[T]{}
	return cmd
}

// NewStackWithHistory creates a new stack that uses a specified stack as
// the main stack.
//
// Parameters:
//   - stack: The stack to use as the main stack.
//
// Returns:
//   - *ud.History[Stacker[T]]: A pointer to the new stack with history.
//
// Behaviors:
//   - If the stack parameter is nil, an ArrayStack is used as the main stack.
//   - In executions, history errors only if it would have errored normally.
func NewStackWithHistory[T any](stack Stacker[T]) *ud.History[Stacker[T]] {
	if stack == nil {
		stack = NewArrayStack[T]()
	}

	h := NewStackWithHistory(stack)

	return h
}
