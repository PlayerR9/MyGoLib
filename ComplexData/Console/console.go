package Console

import (
	"errors"
	"fmt"

	cdom "github.com/PlayerR9/MyGoLib/CustomData/OrderedMap"
)

const (
	// HelpOpcode is the opcode for the help command.
	HelpOpcode string = "help"
)

// CommandInfo represents a console command.
type Console struct {
	// name is the name of the console.
	name string

	// commandMap is a map of command opcodes to CommandInfo.
	commandMap *cdom.OrderedMap[string, *CommandInfo]
}

// NewConsole is a function that creates a new console.
//
// Parameters:
//   - name: The name of the console.
//
// Returns:
//   - *Console: The new console.
func NewConsole(name string) *Console {
	console := &Console{
		name:       name,
		commandMap: cdom.NewOrderedMap[string, *CommandInfo](),
	}

	helpCommand, err := MakeHelpCommand(console)
	if err != nil {
		panic(err)
	}

	console.commandMap.AddEntry(
		HelpOpcode,
		helpCommand,
	)

	return console
}

// AddCommand adds a command to the console.
//
// Parameters:
//   - info: The information about the command.
func (c *Console) AddCommand(info *CommandInfo) {
	if info == nil {
		return
	}

	name := info.GetName()
	if name == HelpOpcode {
		return
	}

	c.commandMap.AddEntry(name, info)
}

// ParseArgs parses the arguments for a command.
//
// Parameters:
//   - args: A slice of strings representing the arguments passed to the command.
//
// Returns:
//   - *ParsedCommand: The parsed command.
//   - error: An error if the command fails.
func (c *Console) ParseArgs(args []string) (*ParsedCommand, error) {
	if len(args) == 0 {
		return nil, errors.New("missing command name")
	}

	commandName := args[0]

	cInfo, ok := c.commandMap.GetEntry(commandName)
	if !ok {
		return nil, fmt.Errorf("command %q does not exist", commandName)
	}

	flagMap, err := cInfo.ParseArgs(args[1:])
	if err != nil {
		return nil, err
	}

	return NewParsedCommand(commandName, flagMap, cInfo.fn), nil
}
