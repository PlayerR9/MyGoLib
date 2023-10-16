// git tag v0.1.20

package CMLine

import (
	"fmt"
)

// ConsoleCommandInfo is a struct that contains information about a console command.
//
// Fields:
//   - Name: The name of the command.
//   - Description: The description of the command.
//   - Flags: The flags of the command.
//   - Callback: A function that is called when the command is parsed.
type ConsoleCommandInfo struct {
	// The name of the command.
	Name string

	// The description of the command.
	Description string

	// The flags of the command.
	Flags []ConsoleFlagInfo

	// A function that is called when the command is parsed.
	Callback func(args map[string]interface{}) (interface{}, error)
}

// ToString returns a string representation of the ConsoleCommandInfo struct.
//
// Returns:
//   - string: A string representation of the ConsoleCommandInfo struct.
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

// ConsoleFlagInfo is a struct that contains information about a console flag.
//
// Fields:
//   - Name: The name of the flag.
//   - Args: The argument name of the flag.
//   - Description: The description of the flag.
//   - Required: Whether or not the flag is required.
//   - Callback: A function that is called when the flag is parsed.
type ConsoleFlagInfo struct {
	// The name of the flag.
	Name string

	// The argument name of the flag.
	Args []string

	// The description of the flag.
	Description string

	// Whether or not the flag is required.
	Required bool

	// A function that is called when the flag is parsed.
	Callback func(...string) (interface{}, error)
}

// ToString returns a string representation of the ConsoleFlagInfo struct.
//
// Returns:
//   - string: A string representation of the ConsoleFlagInfo struct.
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

// HelpToString returns a string representation of the help of a command.
//
// Parameters:
//   - executable_name: The name of the executable.
//   - flags: The flags of the command.
//
// Returns:
//   - string: A string representation of the help of a command.
func HelpToString(executable_name string, commands []ConsoleCommandInfo) (str string) {
	if len(commands) == 0 {
		panic("no commands specified")
	}

	str += fmt.Sprintf("** HELP PAGE OF %s **\n\n", executable_name)

	for _, command := range commands {
		str += fmt.Sprintf("\n%s\n", command.ToString(executable_name))
	}

	return
}

// ParseCommandLine parses the command line arguments.
//
// Parameters:
//   - args: The command line arguments.
//   - commands: The commands to parse.
//
// Returns:
//   - string: The name of the command.
//   - interface{}: The solution of the command.
//   - error: An error if one occurred.
func ParseCommandLine(args []string, commands []ConsoleCommandInfo) (string, interface{}, error) {
	// Check if the command is present
	if len(args) == 0 {
		return "", nil, fmt.Errorf("no command specified")
	}

	// Find the command
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

	// Parse flags
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

// parse_console_flags parses the flags of a command.
//
// Parameters:
//   - args: The arguments of the command.
//   - flags: The flags of the command.
//
// Returns:
//   - map[string]interface{}: A map of the flags and their values.
//   - error: An error if one occurred.
func parse_console_flags(args []string, flags []ConsoleFlagInfo) (map[string]interface{}, error) {
	results := make(map[string]interface{})

	// Check if enough arguments are present
	var min int = 1
	var max int = 1

	for _, f := range flags {
		if f.Required {
			min += len(f.Args) + 1
		}

		max += len(f.Args) + 1
	}

	if len(args) < min {
		return nil, fmt.Errorf("not enough arguments; expected at least %d, got %d", min, len(args))
	} else if len(args) >= max {
		return nil, fmt.Errorf("too many arguments; expected at most %d, got %d", max, len(args))
	}

	// Parse flags
	arg_index := 1

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

		// Call callback function for flag
		inf_tmp, err := f.Callback(args_tmp...)
		if err != nil {
			return nil, fmt.Errorf("could not parse flag %s: %v", f.Name, err)
		}

		// Set result
		results[f.Name] = inf_tmp

		arg_index += len(f.Args)
	}

	return results, nil
}
