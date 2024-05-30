package CmdLineParser

import (
	"strconv"
	"strings"
)

// ErrOpcodeHelp represents an error where the opcode help is reserved.
type ErrOpcodeHelp struct{}

// Error returns the error message: "opcode \"help\" is reserved".
//
// Returns:
//   - string: The error message.
func (e *ErrOpcodeHelp) Error() string {
	return "opcode \"help\" is reserved"
}

// NewErrOpcodeHelp creates a new ErrOpcodeHelp.
//
// Returns:
//   - *ErrOpcodeHelp: A pointer to the new ErrOpcodeHelp.
func NewErrOpcodeHelp() *ErrOpcodeHelp {
	return &ErrOpcodeHelp{}
}

// ErrCommandNotFound represents an error where a command is not found.
type ErrCommandNotFound struct {
	// The command that was not found.
	Command string
}

// Error returns the error message: "command <command> not found".
//
// Returns:
//   - string: The error message.
func (e *ErrCommandNotFound) Error() string {
	var builder strings.Builder

	builder.WriteString("command ")
	builder.WriteString(strconv.Quote(e.Command))
	builder.WriteString(" not found")

	return builder.String()
}

// NewErrCommandNotFound creates a new ErrCommandNotFound with the provided command.
//
// Parameters:
//   - command: The command that was not found.
//
// Returns:
//   - *ErrCommandNotFound: A pointer to the new ErrCommandNotFound.
func NewErrCommandNotFound(command string) *ErrCommandNotFound {
	return &ErrCommandNotFound{Command: command}
}
