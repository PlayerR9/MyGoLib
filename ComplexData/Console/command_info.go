package Console

import (
	cdom "github.com/PlayerR9/MyGoLib/CustomData/OrderedMap"
	fsp "github.com/PlayerR9/MyGoLib/Formatting/FString"

	ue "github.com/PlayerR9/MyGoLib/Units/Errors"
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
	// name is the name of the command.
	name string

	// description is the documentation of the command.
	description []string

	// args is a slice of string representing the arguments accepted by
	// the command. Order matters.
	args *cdom.OrderedMap[string, *Flag]

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

	if inf.args.Size() == 0 {
		err = trav.AppendString("[No arguments available]")
		if err != nil {
			return err
		}
	} else {
		trav.AcceptLine()

		iter := inf.args.Iterator()

		for count := 0; ; count++ {
			item, err := iter.Consume()
			if err != nil {
				break
			}

			err = item.Second.FString(trav)
			if err != nil {
				return ue.NewErrWhileAt("printing", count+1, "argument", err)
			}
		}
	}

	trav.AcceptLine()

	return nil
}

// NewCommandInfo is a function that creates a new command info.
//
// Parameters:
//   - name: The name of the command.
//   - description: The description of the command.
//   - fn: The function to call when the command is executed.
//   - args: A slice of Flag representing the arguments accepted by the command.
//
// Returns:
//   - *CommandInfo: The new command info.
func NewCommandInfo(name string, description []string, fn ConsoleFunc, args ...*Flag) *CommandInfo {
	ci := &CommandInfo{
		name:        name,
		description: description,
		fn:          fn,
		args:        cdom.NewOrderedMap[string, *Flag](),
	}

	for _, arg := range args {
		ci.args.AddEntry(arg.name, arg)
	}

	return ci
}

// GetName returns the name of the command.
//
// Returns:
//   - string: The name of the command.
func (inf *CommandInfo) GetName() string {
	return inf.name
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

	iter := inf.args.Iterator()

	for count := 0; ; count++ {
		item, err := iter.Consume()
		if err != nil {
			break
		}

		if count >= len(args) {
			return nil, NewErrMissingArgument(item.First)
		}

		sol, err := item.Second.parsingFunc(args[count])
		if err != nil {
			return nil, ue.NewErrWhileAt("parsing", count+1, "argument", err)
		}

		flagMap[item.First] = sol
	}

	size := iter.Size()

	if len(args) > size {
		return nil, NewErrArgumentNotRecognized(args[size])
	}

	return flagMap, nil
}
