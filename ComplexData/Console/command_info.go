package Console

import (
	fsp "github.com/PlayerR9/MyGoLib/Formatting/FString"
)

// ConsoleFunc is a function type that represents a callback
// function for a console command.
//
// Parameters:
//   - flagMap: A map of string keys to any values representing
//     the arguments passed to the command.
//
// Returns:
//   - any: The result of the command. (if any)
//   - error: An error if the command fails.
type ConsoleFunc func(flagMap map[string]any) (any, error)

// CommandInfo represents a console command.
type CommandInfo struct {
	// description is the documentation of the command.
	description []string

	// args is a slice of string representing the arguments accepted by
	// the command. Order matters.
	args []string

	// fn is the function to call when the command is executed.
	fn ConsoleFunc
}

// FString generates a formatted string representation of a CommandInfo.
//
// Format:
//
//	<description>
//
//	Arguments: <arg 1> <arg 2> ...
//
// or:
//
//	<description>
//
//	Arguments: [No arguments available]
//
// Parameters:
//   - trav: The traversor to use for the CommandInfo.
//
// Returns:
//   - error: An error if the printing fails.
//
// Behaviors:
//   - <description> is printed according to the DocumentPrinter.
//   - If trav is nil, the function will do nothing.
func (inf *CommandInfo) FString(trav *fsp.Traversor) error {
	if trav == nil {
		return nil
	}

	/*

		docPrinter := fsd.NewDocumentPrinter(
			"Description",
			inf.description,
			"[No description available]",
		)
		err := docPrinter.FString(trav)
		if err != nil {
			return ue.NewErrWhile("printing the command description", err)
		}
	*/

	trav.EmptyLine()

	err := trav.AppendString("Arguments:")
	if err != nil {
		return err
	}

	if len(inf.args) == 0 {
		err := trav.AppendRune(' ')
		if err != nil {
			return err
		}

		err = trav.AppendString("[No arguments available]")
		if err != nil {
			return err
		}

		trav.AcceptLine()
	} else {
		trav.AcceptLine()

		err := trav.AddJoinedLine(" ", inf.args...)
		if err != nil {
			return err
		}
	}

	return nil
}

// NewCommandInfo is a function that creates a new command info.
//
// Parameters:
//   - description: The description of the command.
//   - fn: The function to call when the command is executed.
//   - args: A slice of string representing the arguments accepted by
//     the command. Order matters.
//
// Returns:
//   - *CommandInfo: The new command info.
func NewCommandInfo(description []string, fn ConsoleFunc, args []string) *CommandInfo {
	return &CommandInfo{
		description: description,
		fn:          fn,
		args:        args,
	}
}

// ParseArgs is a function that parses the arguments for a command.
//
// Parameters:
//   - args: A slice of strings representing the arguments passed to the command.
//
// Returns:
//   - map[string]any: A map of string keys to any values representing the arguments
//     parsed by the command.
//   - error: An error if the command fails.
//
// Errors:
//   - *ErrMissingArgument: If an argument is missing.
//   - *ErrArgumentNotRecognized: If an argument is not recognized.
func (inf *CommandInfo) ParseArgs(args []string) (map[string]any, error) {
	flagMap := make(map[string]any)

	for i, arg := range inf.args {
		if i >= len(args) {
			return nil, NewErrMissingArgument(arg)
		}

		flagMap[arg] = args[i]
	}

	if len(args) > len(inf.args) {
		return nil, NewErrArgumentNotRecognized(args[len(inf.args)])
	}

	return flagMap, nil
}
