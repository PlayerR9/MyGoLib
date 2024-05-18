package Console

import (
	"errors"
	"fmt"
	"strings"

	cdm "github.com/PlayerR9/MyGoLib/CustomData/SortedMap"
)

const (
	HelpOpcode string = "help"
)

type Console struct {
	commandMap *cdm.SortedMap[string, *CommandInfo]
}

func NewConsole() *Console {
	console := &Console{
		commandMap: cdm.NewSortedMap[string, *CommandInfo](),
	}

	helpCommand := func(flagMap map[string]any) (any, error) {
		var builder strings.Builder

		builder.WriteString("Here are the commands you can use:\n")

		for _, entry := range console.commandMap.GetEntries() {
			builder.WriteString(fmt.Sprintf("%s : %s\n", entry.First, entry.Second.description))
		}

		return builder.String(), nil
	}

	console.commandMap.AddEntry(
		HelpOpcode,
		NewCommandInfo(
			"Display the help message",
			helpCommand,
			[]string{},
		),
	)

	return console
}

func (c *Console) AddCommand(name string, info *CommandInfo) {
	if name == HelpOpcode {
		return
	}

	c.commandMap.AddEntry(name, info)
}

func (c *Console) ParseArgs(args []string) (*ParsedCommand, error) {
	if len(args) < 2 {
		return nil, errors.New("invalid command")
	}

	cInfo, err := c.commandMap.GetEntry(args[1])
	if err != nil {
		return nil, fmt.Errorf("%q is not a valid command", args[1])
	}

	flagMap, err := cInfo.ParseArgs(args[2:])
	if err != nil {
		return nil, err
	}

	return NewParsedCommand(args[1], flagMap, cInfo.fn), nil
}
