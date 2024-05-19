package Console

import (
	"errors"
	"fmt"
	"strings"

	cdm "github.com/PlayerR9/MyGoLib/CustomData/OrderedMap"
)

const (
	HelpOpcode string = "help"
)

type Console struct {
	commandMap *cdm.OrderedMap[string, *CommandInfo]
}

func NewConsole() *Console {
	console := &Console{
		commandMap: cdm.NewOrderedMap[string, *CommandInfo](),
	}

	helpCommand := func(flagMap map[string]any) (any, error) {
		lines := []string{
			"Here are the commands you can use:",
		}

		var builder strings.Builder

		iter := console.commandMap.Iterator()

		for {
			entry, err := iter.Consume()
			if err != nil {
				break
			}

			builder.WriteString(entry.First)
			builder.WriteRune(' ')
			builder.WriteRune(':')
			builder.WriteRune(' ')
			builder.WriteString(entry.Second.description)

			lines = append(lines, builder.String())
			builder.Reset()
		}

		return strings.Join(lines, "\n"), nil
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

	cInfo, ok := c.commandMap.GetEntry(args[1])
	if !ok {
		return nil, fmt.Errorf("%q is not a valid command", args[1])
	}

	flagMap, err := cInfo.ParseArgs(args[2:])
	if err != nil {
		return nil, err
	}

	return NewParsedCommand(args[1], flagMap, cInfo.fn), nil
}
