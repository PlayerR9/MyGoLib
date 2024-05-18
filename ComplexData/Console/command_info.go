package Console

import "errors"

type ConsoleFunc func(flagMap map[string]any) (any, error)

type CommandInfo struct {
	description string

	args []string

	fn ConsoleFunc
}

func NewCommandInfo(description string, fn ConsoleFunc, args []string) *CommandInfo {
	return &CommandInfo{
		description: description,
		fn:          fn,
		args:        args,
	}
}

func (inf *CommandInfo) ParseArgs(args []string) (map[string]any, error) {
	if len(inf.args) > len(args) {
		return nil, errors.New("not enough arguments")
	} else if len(inf.args) < len(args) {
		return nil, errors.New("too many arguments")
	}

	flagMap := make(map[string]any)

	for i, arg := range inf.args {
		flagMap[arg] = args[i]
	}

	return flagMap, nil
}
