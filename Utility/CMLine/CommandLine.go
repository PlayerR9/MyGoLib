package CMLine

import (
	"errors"
	"fmt"
	"strings"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

// ConsoleFlagInfo represents a flag for a console command.
// It contains the name, arguments, and description of the flag,
// as well as a boolean indicating whether the flag is required.
// The Callback function is invoked when the flag is used,
// with the provided arguments. It returns a result and an error, if any.
type ConsoleFlagInfo struct {
	// Name is the name of the flag.
	name string

	// Args is a slice of strings that represents the arguments accepted by the flag.
	args []string

	// Description provides a brief explanation of what the flag does.
	description []string

	// Required is a boolean that indicates whether the flag is required.
	required bool

	// Callback is a function that is invoked when the flag is used.
	// It takes a variadic number of string arguments,
	// and returns a result of any type and an error, if any.
	callback func(...string) (any, error)
}

// FlagInfoOption is a function type that modifies ConsoleFlagInfo.
type FlagInfoOption func(*ConsoleFlagInfo)

// WithArgs is a FlagInfoOption that sets the arguments for a ConsoleFlagInfo.
// It trims the space from each argument and ignores empty arguments.
//
// Parameters:
//
//   - args: The arguments to set.
//
// Returns:
//
//   - A FlagInfoOption that sets the arguments for a ConsoleFlagInfo.
func WithArgs(args ...string) FlagInfoOption {
	return func(flag *ConsoleFlagInfo) {
		flag.args = make([]string, 0, len(args))

		for _, arg := range args {
			arg = strings.TrimSpace(arg)
			if arg != "" {
				flag.args = append(flag.args, arg)
			}
		}

		flag.args = flag.args[:len(flag.args)]
	}
}

// WithDescription is a FlagInfoOption that sets the description for a ConsoleFlagInfo.
// It splits each line of the description by newline characters.
//
// Parameters:
//
//   - description: The description to set.
//
// Returns:
//
//   - A FlagInfoOption that sets the description for a ConsoleFlagInfo.
func WithDescription(description ...string) FlagInfoOption {
	return func(flag *ConsoleFlagInfo) {
		for _, line := range description {
			fields := strings.Split(line, "\n")

			flag.description = append(flag.description, fields...)
		}
	}
}

// WithRequired is a FlagInfoOption that sets whether a ConsoleFlagInfo is required.
//
// Parameters:
//
//   - required: Whether the flag is required.
//
// Returns:
//
//   - A FlagInfoOption that sets whether a ConsoleFlagInfo is required.
func WithRequired(required bool) FlagInfoOption {
	return func(flag *ConsoleFlagInfo) {
		flag.required = required
	}
}

// ConsoleCommandInfo represents a console command.
// It contains the name and description of the command, as well as a list of
// flags that the command accepts.
// The Callback function is invoked when the command is executed, with the
// provided arguments. It returns a result and an error, if any.
type ConsoleCommandInfo struct {
	// Name is the name of the console command.
	name string

	// Description provides a brief explanation of what the command does.
	description []string

	// Flags is a slice of ConsoleFlagInfo that represents the flags accepted by
	// the command.
	flags []*ConsoleFlagInfo

	// Callback is a function that is invoked when the command is executed.
	// It takes a map of string keys to any type values as arguments, and returns
	// a result of any type and an error, if any.
	callback func(args map[string]any) (any, error)

	// errReason is the reason why the command could not be created.
	errReason error
}

// CommandInfoOption is a function type that modifies ConsoleCommandInfo.
type CommandInfoOption func(*ConsoleCommandInfo)

// WithFlag is a CommandInfoOption that adds a new flag to a ConsoleCommandInfo.
//
// Parameters:
//
//   - name: The name of the flag.
//   - callback: The function to call when the flag is used.
//   - options: The options to apply to the flag.
//
// Returns:
//
//   - A CommandInfoOption that adds a new flag to a ConsoleCommandInfo.
func WithFlag(name string, callback func(...string) (any, error), options ...FlagInfoOption) CommandInfoOption {
	return func(command *ConsoleCommandInfo) {
		newFlag := &ConsoleFlagInfo{
			name:        name,
			args:        make([]string, 0),
			description: make([]string, 0),
			required:    false,
			callback:    callback,
		}

		name = strings.TrimSpace(name)
		if name == "" {
			command.errReason = fmt.Errorf("could not create flag: %v",
				ers.NewErrInvalidParameter("name").WithReason(errors.New("flag name cannot be empty")))
			return
		}

		if callback == nil {
			command.errReason = fmt.Errorf("could not create flag %s: %v", name,
				ers.NewErrInvalidParameter("callback").WithReason(errors.New("flag callback cannot be nil")))
			return
		}

		for _, option := range options {
			option(newFlag)
		}

		command.flags = append(command.flags, newFlag)
	}
}

// WithCallback is a CommandInfoOption that sets the callback for a ConsoleCommandInfo.
//
// Parameters:
//
//   - callback: The function to call when the command is used.
//
// Returns:
//
//   - A CommandInfoOption that sets the callback for a ConsoleCommandInfo.
func WithCallback(callback func(map[string]any) (any, error)) CommandInfoOption {
	return func(command *ConsoleCommandInfo) {
		if callback == nil {
			command.errReason = ers.NewErrInvalidParameter("callback").
				WithReason(errors.New("callback cannot be nil"))

			return
		}

		command.callback = callback
	}
}

// CMLine represents a command line interface.
// It contains the name of the executable and a map of commands that the
// interface accepts.
type CMLine struct {
	// executableName is the name of the executable.
	executableName string

	// commands is a map of string keys to ConsoleCommandInfo values that
	// represents the commands accepted by the interface.
	commands map[string]*ConsoleCommandInfo

	// errReason is the reason why the CMLine could not be created.
	errReason error
}

// CMLineOption is a function type that modifies CMLine.
type CMLineOption func(*CMLine)

// WithExecutableName is a CMLineOption that sets the executable name for a CMLine.
// Note: It trims the space from the name.
//
// Parameters:
//
//   - name: The name of the executable.
//
// Returns:
//
//   - A CMLineOption that sets the executable name for a CMLine.
func WithExecutableName(name string) CMLineOption {
	return func(cm *CMLine) {
		name = strings.TrimSpace(name)
		if name == "" {
			cm.errReason = ers.NewErrInvalidParameter("name").
				WithReason(errors.New("executable name cannot be empty"))

			return
		}

		cm.executableName = name
	}
}

// WithCommand is a CMLineOption that adds a new command to a CMLine.
//
// Note:
//
//   - commands with the same name are overwritten. Last one wins.
//   - It trims the space from the name.
//
// Parameters:
//
//   - name: The name of the command.
//   - options: The options to apply to the command.
//
// Returns:
//
//   - A CMLineOption that adds a new command to a CMLine.
func WithCommand(name string, options ...CommandInfoOption) CMLineOption {
	return func(cm *CMLine) {
		name = strings.TrimSpace(name)
		if name == "" {
			cm.errReason = ers.NewErrInvalidParameter("name").
				WithReason(errors.New("name cannot be empty"))

			return
		}

		newCommand := &ConsoleCommandInfo{
			name:        name,
			description: make([]string, 0),
			flags:       make([]*ConsoleFlagInfo, 0),
			callback:    nil,
			errReason:   nil,
		}

		for i, option := range options {
			option(newCommand)

			if newCommand.errReason != nil {
				cm.errReason = fmt.Errorf("invalid option %d for command %s: %v", i, name, newCommand.errReason)

				return
			}
		}

		cm.commands[name] = newCommand
	}
}

// NewCMLine creates a new CMLine with the given options.
//
// Parameters:
//
//   - options: The options to apply to the CMLine.
//
// Returns:
//
//   - A pointer to the created CMLine.
func NewCMLine(options ...CMLineOption) *CMLine {
	cml := &CMLine{
		commands:  make(map[string]*ConsoleCommandInfo),
		errReason: nil,
	}

	for _, option := range options {
		option(cml)
	}

	// Check if any command errors occurred
	for key, command := range cml.commands {
		cml.errReason = fmt.Errorf("could not create command %s: %v", key, command.errReason)

		return cml
	}

	return cml
}

// ParseCommandLine parses the command line arguments and executes the
// appropriate command.
//
// Parameters:
//
//   - args: The command line arguments.
//
// Returns:
//
//   - The name of the command that was executed.
//   - The result of the command's callback function.
//   - Any error that occurred.
//
// The function first checks if the command exists in the provided commands.
// If the command is not found, an error is returned.
// The function then parses the flags for the command.
// If the flags cannot be parsed, an error is returned.
// Finally, the function executes the command's callback function and returns
// the command's name, the result of the callback, and any error that occurred.
// If there was a building error, it returns the error reason.
func (cml *CMLine) ParseCommandLine(args []string) (string, any, error) {
	// Check if there was a building error
	if cml.errReason != nil {
		return "", nil, cml.errReason
	}

	// Check if any arguments were provided
	if len(args) == 0 {
		return "", nil, ers.NewErrInvalidParameter("args").
			WithReason(errors.New("no arguments provided"))
	}

	// Get the command from the command map
	command, exists := cml.commands[args[0]]
	if !exists {
		return "", nil, fmt.Errorf("command '%s' not found", args[0])
	}

	// Create a map to store the flags
	commandFlags := make(map[string]any)

	if len(args) > 1 {
		// Parse the flags if provided
		var err error

		commandFlags, err = parseConsoleFlags(args[1:], command.flags)
		if err != nil {
			return "", nil, err
		}
	}

	// Check if the command has a callback function
	if command.callback == nil {
		return command.name, nil, nil
	}

	// Call the callback function with the flags
	commandSolution, err := command.callback(commandFlags)
	if err != nil {
		return "", nil, fmt.Errorf("failed to execute command '%s': reason=%v", command.name, err)
	}

	return command.name, commandSolution, nil
}

// parseConsoleFlags parses the console flags from the command line arguments.
//
// Parameters:
//
//   - args: The command line arguments.
//   - flags: The flags to parse.
//
// Returns:
//
//   - A map of flag names to results.
//   - Any error that occurred.
//
// The function first checks if the number of arguments is within the expected range.
// If not, an error is returned.
// The function then iterates over the flags.
// If a required flag is not present, an error is returned.
// If a flag is present but not enough arguments are specified for it, an error is returned.
// The function then calls the flag's callback function with the arguments and stores the result.
// If the callback function returns an error, it is returned by the function.
// Finally, the function returns a map of flag names to results and any error that occurred.
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
func parseConsoleFlags(args []string, flags []*ConsoleFlagInfo) (map[string]any, error) {
	// Create a map to store the console flags for quick lookup
	consoleFlagMap := make(map[string]*ConsoleFlagInfo)
	for _, consoleFlag := range flags {
		consoleFlagMap[consoleFlag.name] = consoleFlag
	}

	// Create a map to store the parsed results
	parsedResults := make(map[string]any)
	currentArgIndex := 0

	for currentArgIndex < len(args) {
		// Get the console flag name from the current argument
		consoleFlagName := args[currentArgIndex]

		// Check if the console flag exists in the map
		consoleFlag, exists := consoleFlagMap[consoleFlagName]
		if !exists {
			return nil, fmt.Errorf("unknown flag '%s' provided", consoleFlagName)
		}

		// Check if there are enough arguments for the console flag
		if len(consoleFlag.args)+currentArgIndex >= len(args) {
			return nil, fmt.Errorf("flag '%s' requires more arguments", consoleFlag.name)
		}

		// Move to the next argument
		currentArgIndex++

		// Create a temporary slice to store the arguments for the console flag
		tempArgs := make([]string, len(args[currentArgIndex:]))
		copy(tempArgs, args[currentArgIndex:])

		// Call the callback function for the console flag with the arguments
		parsedFlag, err := consoleFlag.callback(tempArgs...)
		if err != nil {
			return nil, fmt.Errorf("failed to parse flag '%s': reason=%v", consoleFlag.name, err)
		}

		// Store the result of the callback function in the parsed results map
		// amd move to the next argument
		parsedResults[consoleFlag.name] = parsedFlag
		currentArgIndex += len(consoleFlag.args)
	}

	return parsedResults, nil
}

// FString formats the CMLine into a string that can be displayed to the user.
//
// Returns:
//
//   - The formatted string.
func (cml *CMLine) FString() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "Usage: %s <command> [flags]\n", cml.executableName)

	if len(cml.commands) == 0 {
		fmt.Fprint(&builder, "Commands: None\n")
	} else {
		fmt.Fprint(&builder, "Commands:\n")

		for _, command := range cml.commands {
			fmt.Fprintf(&builder, "%s\n", command.FString(1))
		}
	}

	return builder.String()
}

// FString formats the ConsoleFlagInfo into a string that can be displayed to the user.
//
// Parameters:
//
//   - indentLevel: The indentation level of the output.
//
// Returns:
//
//   - The formatted string.
func (cfi *ConsoleFlagInfo) FString(indentLevel int) string {
	var builder strings.Builder

	indentation := strings.Repeat("\t", indentLevel)

	// Add the flag name
	fmt.Fprintf(&builder, "%s%s", indentation, cfi.name)

	// Add the arguments
	for _, arg := range cfi.args {
		fmt.Fprintf(&builder, " <%s>", arg)
	}

	// Add the description
	fmt.Fprintf(&builder, "\n%sDescription:\n%s\t%s", indentation, indentation, cfi.description)

	// Add the required information
	if len(cfi.args) != 0 {
		fmt.Fprintf(&builder, "\n%sRequired: ", indentation)

		if cfi.required {
			fmt.Fprint(&builder, "Yes")
		} else {
			fmt.Fprint(&builder, "No")
		}
	}

	return builder.String()
}

// FString formats the ConsoleCommandInfo into a string that can be displayed
// to the user.
//
// Parameters:
//
//   - indentLevel: The indentation level of the output.
//
// Returns:
//
//   - The formatted string.
func (cci *ConsoleCommandInfo) FString(indentLevel int) string {
	indentation := strings.Repeat("\t", indentLevel)
	var builder strings.Builder

	// Add the command name
	fmt.Fprintf(&builder, "%sCommand: %s\n", indentation, cci.name)

	// Add the command description
	if len(cci.description) == 0 {
		fmt.Fprintf(&builder, "%sDescription: [No description provided]\n", indentation)
	} else {
		fmt.Fprintf(&builder, "%sDescription:\n", indentation)
		for _, line := range cci.description {
			fmt.Fprintf(&builder, "%s\t%s\n", indentation, line)
		}
	}

	// Add the usage information for each flag
	for _, flag := range cci.flags {
		fmt.Fprintf(&builder, "%sUsage: %s", indentation, cci.name)

		if flag.required {
			fmt.Fprintf(&builder, " %s", flag.name)
		} else {
			fmt.Fprintf(&builder, " [%s]", flag.name)
		}

		for _, arg := range flag.args {
			fmt.Fprintf(&builder, " <%s>", arg)
		}

		builder.WriteString("\n")
	}

	// Add the flag information
	if len(cci.flags) == 0 {
		fmt.Fprintf(&builder, "%sFlags: None\n", indentation)
	} else {
		fmt.Fprintf(&builder, "%sFlags:\n", indentation)

		for _, flag := range cci.flags {
			fmt.Fprintf(&builder, "%s\n", flag.FString(indentLevel+1))
		}
	}

	return builder.String()
}
