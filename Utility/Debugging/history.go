package Debugging

import (
	uc "github.com/PlayerR9/lib_units/common"
)

// Commander is an interface that represents a command that can be
// executed and undone.
type Commander[T any] interface {
	// Execute executes the command.
	//
	// Parameters:
	//   - data: The data to execute the command on.
	//
	// Returns:
	//   - error: An error if the execution fails.
	Execute(data T) error

	// Undo undoes the command.
	//
	// Parameters:
	//   - data: The data to undo the command on.
	//
	// Returns:
	//   - error: An error if the undo fails.
	Undo(data T) error

	uc.Copier
}

// Command represents a generic command that can be executed and undone.
type Command[T any] struct {
	// execute represents the function to execute the command.
	execute func(data T) error

	// undo represents the function to undo the command.
	undo func(data T) error
}

// Execute implements the Commander interface.
func (c *Command[T]) Execute(data T) error {
	err := c.execute(data)
	if err != nil {
		return err
	}

	return nil
}

// Undo implements the Commander interface.
func (c *Command[T]) Undo(data T) error {
	err := c.undo(data)
	if err != nil {
		return err
	}

	return nil
}

// Copy implements the Commander interface.
func (c *Command[T]) Copy() uc.Copier {
	cCopy := &Command[T]{
		execute: c.execute,
		undo:    c.undo,
	}

	return cCopy
}

// NewCommand creates a new command with the given execute and undo functions.
//
// Parameters:
//   - execute: The function to execute the command.
//   - undo: The function to undo the command.
//
// Returns:
//   - Commander: The new command.
//
// Behaviors:
//   - If either the execute or undo functions are nil, nil is returned.
func NewCommand[T any](execute, undo func(data T) error) Commander[T] {
	if execute == nil || undo == nil {
		return nil
	}

	cmd := &Command[T]{
		execute: execute,
		undo:    undo,
	}

	return cmd
}

// History represents a history of commands that can be executed and undone.
type History[T any] struct {
	// data represents the data that the commands are executed on.
	data T

	// commands represents the commands that have been executed.
	commands []Commander[T]
}

// Copy implements the uc.Copier interface.
func (h *History[T]) Copy() uc.Copier {
	hCopy := &History[T]{
		data:     uc.CopyOf(h.data).(T),
		commands: make([]Commander[T], len(h.commands)),
	}

	for i, cmd := range h.commands {
		cmdCopy := cmd.Copy().(Commander[T])
		hCopy.commands[i] = cmdCopy
	}

	return hCopy
}

// NewHistory creates a new history with the given data.
//
// Parameters:
//   - data: The data to create the history with.
//
// Returns:
//   - *History: The new history.
func NewHistory[T any](data T) *History[T] {
	h := &History[T]{
		data:     data,
		commands: make([]Commander[T], 0),
	}

	return h
}

// ExecuteCommand executes a command on the history.
//
// Parameters:
//   - cmd: The command to execute.
//
// Returns:
//   - error: An error if the execution fails.
//
// Behaviors:
//   - If the command is nil, no action is taken.
func (h *History[T]) ExecuteCommand(cmd Commander[T]) error {
	if cmd == nil {
		return nil
	}

	err := cmd.Execute(h.data)
	h.commands = append(h.commands, cmd)

	if err != nil {
		return err
	}

	return nil
}

// UndoLastCommand undoes the last command executed on the history.
//
// Returns:
//   - error: An error if the undo fails.
//
// Behaviors:
//   - If there are no commands to undo, no action is taken.
func (h *History[T]) UndoLastCommand() error {
	if len(h.commands) == 0 {
		return nil
	}

	lc := h.commands[len(h.commands)-1]
	err := lc.Undo(h.data)
	h.commands = h.commands[:len(h.commands)-1]

	if err != nil {
		return err
	}

	return nil
}

// ReadData reads the data from the history.
//
// Parameters:
//   - f: The function to read the data.
//
// Behaviors:
//   - The function is called with the data from the history.
func (h *History[T]) ReadData(f func(data T)) {
	f(h.data)
}

// Accept accepts the history, clearing the commands.
//
// WARNING: Because the commands are cleared, they cannot be undone after
// accepting the history. Thus, be sure to use this method only when
// you are certain that you will never need to undo the commands.
func (h *History[T]) Accept() {
	h.commands = h.commands[:0]
}

// Reject rejects the history, undoing all commands.
//
// Returns:
//   - error: An error if the undo fails.
func (h *History[T]) Reject() error {
	for len(h.commands) > 0 {
		err := h.UndoLastCommand()
		if err != nil {
			return err
		}
	}

	return nil
}

// GetData returns the data from the history.
func (h *History[T]) GetData() T {
	return h.data
}
