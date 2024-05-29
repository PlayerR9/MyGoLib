package Console

import (
	cdom "github.com/PlayerR9/MyGoLib/CustomData/OrderedMap"
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	ue "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// MakeHelpCommand is a function that creates a help command for a console.
//
// The help command, when run, returns as the first value a slice of
// []*ContentBox.MultilineText representing the list of available commands
// according to FString formatting.
//
// Parameters:
//   - flagMap: A map of string keys to any values representing
//     the arguments passed to the command.
//
// Returns:
//   - *CommandInfo: The help command.
//   - error: An error of type *Errors.ErrInvalidParameter if the console is nil.
func MakeHelpCommand(console *Console) (*CommandInfo, error) {
	if console == nil {
		return nil, ue.NewErrNilParameter("console")
	}

	fn := func(flagMap map[string]any) (any, error) {
		mip := cdom.NewOrderedMapPrinter(
			"Here's a list of all the available commands:",
			console.commandMap,
			"Command",
			"[No commands available]",
		)

		doc, err := ffs.SprintFString(ffs.NewFormatter(), mip)
		if err != nil {
			return nil, ue.NewErrWhile("printing the command list", err)
		}

		return ffs.Stringfy(doc), nil
	}

	doc, err := ffs.Sprintln(ffs.DefaultFormatter, "Display the help message.")
	if err != nil {
		return nil, err
	}

	return NewCommandInfo(
		ffs.Stringfy(doc),
		fn,
		[]string{},
	), nil
}
