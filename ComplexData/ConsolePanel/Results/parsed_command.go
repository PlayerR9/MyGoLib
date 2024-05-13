// Package CnsPanel provides a structure and functions for handling
// console command flags.
package Results

import (
	com "github.com/PlayerR9/MyGoLib/ComplexData/ConsolePanel/Common"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// ParsedCommand represents a parsed console command.
type ParsedCommand struct {
	// args is the arguments passed to the command.
	args com.ArgumentsMap

	// callback is the function to call when the command is used.
	callback com.CommandCallbackFunc
}

// NewParsedCommand creates a new ParsedCommand with the provided name, arguments,
// and callback function.
//
// Parameters:
//   - args: A map of string keys to any values representing the arguments passed
//     to the command.
//   - callbackFunc: The function to call when the command is used.
//
// Returns:
//   - *ParsedCommand: A pointer to the new ParsedCommand.
//   - error: An error of type *ers.ErrInvalidParameter if the callbackFunc is nil.
//
// Behaviors:
//   - If callbackFunc is nil, NoCommandCallback is used.
func NewParsedCommand(args com.ArgumentsMap, callbackFunc com.CommandCallbackFunc) (*ParsedCommand, error) {
	if callbackFunc == nil {
		return nil, ers.NewErrNilParameter("callbackFunc")
	}

	return &ParsedCommand{
		args:     args,
		callback: callbackFunc,
	}, nil
}

// Execute executes the callback function for the parsed command.
//
// Returns:
//   - error: An error if the callback function fails.
func (pc *ParsedCommand) Execute() error {
	return pc.callback(pc.args)
}
