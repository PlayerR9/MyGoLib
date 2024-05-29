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
		printer := ffs.NewStdPrinter(ffs.NewFormatter())
		trav := printer.TraversorOf()

		err := trav.AddJoinedLine(" ", "Usage:", console.name, "<command>", "[flags]")
		if err != nil {
			return nil, err
		}

		mip := cdom.NewOrderedMapPrinter(
			"Here's a list of all the available commands:",
			console.commandMap,
			"Command",
			"[No commands available]",
		)

		err = mip.FString(trav)
		if err != nil {
			return nil, err
		}

		trav.Clean()

		doc := printer.GetPages()

		return ffs.Stringfy(doc), nil
	}

	doc, err := ffs.Sprintln(ffs.DefaultFormatter, "Display the help message.")
	if err != nil {
		return nil, err
	}

	return NewCommandInfo(
		HelpOpcode,
		ffs.Stringfy(doc),
		fn,
	), nil
}
