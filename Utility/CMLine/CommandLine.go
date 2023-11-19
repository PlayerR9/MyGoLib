// git tag v0.1.26

package CMLine

import (
	"fmt"
)

type ConsoleCommandInfo struct {
	Name        string
	Description string
	Flags       []ConsoleFlagInfo
	Callback    func(args map[string]interface{}) (interface{}, error)
}

func (cci ConsoleCommandInfo) ToString(executable_name string) (str string) {
	str += fmt.Sprintf("Command: %s\nDescription:\n\t%s\n\nUsage: %s %s", cci.Name, cci.Description, executable_name, cci.Name)

	for _, flag := range cci.Flags {
		if !flag.Required {
			str += " ["
		} else {
			str += " "
		}

		str += flag.Name

		for _, arg := range flag.Args {
			str += fmt.Sprintf(" <%s>", arg)
		}

		if !flag.Required {
			str += "]"
		}
	}

	str += "\n\nFlags:"

	if len(cci.Flags) == 0 {
		str += " None"

		return
	}

	str += fmt.Sprintf("\n%s", cci.Flags[0].ToString(1))

	for _, flag := range cci.Flags[1:] {
		str += fmt.Sprintf("\n\n%s", flag.ToString(1))
	}

	return
}

type ConsoleFlagInfo struct {
	Name        string
	Args        []string
	Description string
	Required    bool
	Callback    func(...string) (interface{}, error)
}

func (cfi ConsoleFlagInfo) ToString(indent_level int) (str string) {
	var indentation string

	for i := 0; i < indent_level; i++ {
		indentation += "\t"
	}

	str += fmt.Sprintf("%s%s", indentation, cfi.Name)

	if len(cfi.Args) != 0 {
		for _, arg := range cfi.Args {
			str += " <" + arg + ">"
		}
	}

	str += fmt.Sprintf("\n\t%sDescription:\n\t\t%s%s", indentation, indentation, cfi.Description)

	if len(cfi.Args) != 0 {
		str += fmt.Sprintf("\n\t%sRequired: ", indentation)

		if cfi.Required {
			str += "Yes"
		} else {
			str += "No"
		}
	}

	return
}

func ParseCommandLine(args []string, commands []ConsoleCommandInfo) (string, interface{}, error) {
	found_index := -1

	for i, command := range commands {
		if command.Name == args[0] {
			if found_index != -1 {
				return "", nil, fmt.Errorf("command %s and %s share the same name", commands[found_index].Name, command.Name)
			}

			found_index = i
		}
	}

	if found_index == -1 {
		return "", nil, fmt.Errorf("command %s not found", args[0])
	}
	command := commands[found_index]

	flags, err := parse_console_flags(args[1:], command.Flags)
	if err != nil {
		return "", nil, fmt.Errorf("could not parse flags of command %s: %v", command.Name, err)
	}

	solution, err := command.Callback(flags)
	if err != nil {
		return "", nil, fmt.Errorf("could not execute command %s: %v", command.Name, err)
	}

	return command.Name, solution, nil
}

func parse_console_flags(args []string, flags []ConsoleFlagInfo) (map[string]interface{}, error) {
	results := make(map[string]interface{})
	var min int = 0
	var max int = 0

	for _, f := range flags {
		if f.Required {
			min += len(f.Args) + 1
		}

		max += len(f.Args) + 1
	}

	if len(args) < min {
		return nil, fmt.Errorf("not enough arguments; expected at least %d, got %d", min, len(args))
	} else if len(args) > max {
		return nil, fmt.Errorf("too many arguments; expected at most %d, got %d", max, len(args))
	}
	arg_index := 0

	for _, f := range flags {
		if arg_index >= len(args) {
			break
		}

		if f.Name != args[arg_index] {
			if f.Required {
				return nil, fmt.Errorf("required flag %s not present", f.Name)
			}

			continue
		}

		if len(f.Args)+arg_index >= len(args) {
			return nil, fmt.Errorf("flag %s present but not enough arguments specified", f.Name)
		}

		arg_index++

		args_tmp := make([]string, 0)

		for i := 0; i < len(f.Args); i++ {
			args_tmp = append(args_tmp, args[arg_index+i])
		}
		inf_tmp, err := f.Callback(args_tmp...)
		if err != nil {
			return nil, fmt.Errorf("could not parse flag %s: %v", f.Name, err)
		}
		results[f.Name] = inf_tmp

		arg_index += len(f.Args)
	}

	return results, nil
}
