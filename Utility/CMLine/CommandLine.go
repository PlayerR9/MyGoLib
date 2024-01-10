package CMLine

import (
	"fmt"
	"strings"
)

// ConsoleCommandInfo represents a console command.
// It contains the name and description of the command, as well as a list of
// flags that the command accepts.
// The Callback function is invoked when the command is executed, with the
// provided arguments. It returns a result and an error, if any.
type ConsoleCommandInfo struct {
	// Name is the name of the console command.
	Name string

	// Description provides a brief explanation of what the command does.
	Description string

	// Flags is a slice of ConsoleFlagInfo that represents the flags accepted by
	// the command.
	Flags []ConsoleFlagInfo

	// Callback is a function that is invoked when the command is executed.
	// It takes a map of string keys to any type values as arguments, and returns
	// a result of any type and an error, if any.
	Callback func(args map[string]any) (any, error)
}

// FString formats the ConsoleCommandInfo into a string that can be displayed
// to the user.
// It includes the command name, description, usage, and flags.
// The executableName parameter is used to show the command usage.
// If the command has no flags, "None" is displayed.
// Each flag is formatted using its FString method.
// The method returns the formatted string.
func (cci ConsoleCommandInfo) FString(executableName string) string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "Command: %s\n", cci.Name)
	fmt.Fprintf(&builder, "Description:\n\t%s\n\n", cci.Description)
	fmt.Fprintf(&builder, "Usage: %s %s", executableName, cci.Name)

	for _, flag := range cci.Flags {
		if flag.Required {
			fmt.Fprintf(&builder, " %s", flag.Name)
		} else {
			fmt.Fprintf(&builder, " [%s]", flag.Name)
		}

		for _, arg := range flag.Args {
			fmt.Fprintf(&builder, " <%s>", arg)
		}
	}

	fmt.Fprint(&builder, "\n\nFlags:")

	if len(cci.Flags) == 0 {
		fmt.Fprint(&builder, " None")
	} else {
		for i, flag := range cci.Flags {
			if i > 0 {
				fmt.Fprint(&builder, "\n\n")
			}
			fmt.Fprintf(&builder, "%s", flag.FString(1))
		}
	}

	return builder.String()
}

// ConsoleFlagInfo represents a flag for a console command.
// It contains the name, arguments, and description of the flag,
// as well as a boolean indicating whether the flag is required.
// The Callback function is invoked when the flag is used,
// with the provided arguments. It returns a result and an error, if any.
type ConsoleFlagInfo struct {
	// Name is the name of the flag.
	Name string

	// Args is a slice of strings that represents the arguments accepted by the flag.
	Args []string

	// Description provides a brief explanation of what the flag does.
	Description string

	// Required is a boolean that indicates whether the flag is required.
	Required bool

	// Callback is a function that is invoked when the flag is used.
	// It takes a variadic number of string arguments,
	// and returns a result of any type and an error, if any.
	Callback func(...string) (any, error)
}

// FString formats the ConsoleFlagInfo into a string that can be displayed to the user.
// It includes the flag name, description, and whether it's required.
// The indentLevel parameter is used to control the indentation of the output.
// If the flag has arguments, they are included in the output.
// The method returns the formatted string.
func (cfi ConsoleFlagInfo) FString(indentLevel int) string {
	var builder strings.Builder

	indentation := strings.Repeat("\t", indentLevel)

	fmt.Fprintf(&builder, "%s%s", indentation, cfi.Name)

	for _, arg := range cfi.Args {
		fmt.Fprintf(&builder, " <%s>", arg)
	}

	fmt.Fprintf(&builder, "\n\t%sDescription:\n\t\t%s%s", indentation, indentation, cfi.Description)

	if len(cfi.Args) != 0 {
		fmt.Fprintf(&builder, "\n\t%sRequired: ", indentation)

		if cfi.Required {
			fmt.Fprint(&builder, "Yes")
		} else {
			fmt.Fprint(&builder, "No")
		}
	}

	return builder.String()
}

// ParseCommandLine parses the command line arguments and executes the
// appropriate command.
// It takes a slice of arguments and a slice of ConsoleCommandInfo.
// The function first checks if the command exists in the provided commands.
// If the command is not found, an error is returned.
// If the command is found, the function checks if there are any other
// commands with the same name.
// If there are, an error is returned.
// The function then parses the flags for the command.
// If the flags cannot be parsed, an error is returned.
// Finally, the function executes the command's callback function and returns
// the command's name, the result of the callback, and any error that occurred.
func ParseCommandLine(args []string, commands []ConsoleCommandInfo) (string, any, error) {
	commandMap := make(map[string]ConsoleCommandInfo)
	for _, command := range commands {
		if _, exists := commandMap[command.Name]; exists {
			return "", nil, fmt.Errorf("duplicate command name: %s", command.Name)
		}
		commandMap[command.Name] = command
	}

	commandName := args[0]
	command, exists := commandMap[commandName]
	if !exists {
		return "", nil, fmt.Errorf("command %s not found", commandName)
	}

	flags, err := parseConsoleFlags(args[1:], command.Flags)
	if err != nil {
		return "", nil, fmt.Errorf("could not parse flags of command %s: %v", command.Name, err)
	}

	solution, err := command.Callback(flags)
	if err != nil {
		return "", nil, fmt.Errorf("could not execute command %s: %v", command.Name, err)
	}

	return command.Name, solution, nil
}

// parseConsoleFlags parses the console flags from the command line arguments.
// It takes a slice of arguments and a slice of ConsoleFlagInfo.
// The function first checks if the number of arguments is within the expected range.
// If not, an error is returned.
// The function then iterates over the flags.
// If a required flag is not present, an error is returned.
// If a flag is present but not enough arguments are specified for it, an error is returned.
// The function then calls the flag's callback function with the arguments and stores the result.
// If the callback function returns an error, it is returned by the function.
// Finally, the function returns a map of flag names to results and any error that occurred.
func parseConsoleFlags(args []string, flags []ConsoleFlagInfo) (map[string]any, error) {
	flagMap := make(map[string]ConsoleFlagInfo)
	for _, flag := range flags {
		flagMap[flag.Name] = flag
	}

	results := make(map[string]any)
	argIndex := 0

	for argIndex < len(args) {
		flagName := args[argIndex]
		flag, exists := flagMap[flagName]
		if !exists {
			return nil, fmt.Errorf("unknown flag: %s", flagName)
		}

		if len(flag.Args)+argIndex >= len(args) {
			return nil, fmt.Errorf("flag %s present but not enough arguments specified", flag.Name)
		}

		argIndex++

		argsTmp := make([]string, len(flag.Args))
		for i := 0; i < len(flag.Args); i++ {
			argsTmp[i] = args[argIndex+i]
		}

		infTmp, err := flag.Callback(argsTmp...)
		if err != nil {
			return nil, fmt.Errorf("could not parse flag %s: %v", flag.Name, err)
		}
		results[flag.Name] = infTmp

		argIndex += len(flag.Args)
	}

	return results, nil
}
