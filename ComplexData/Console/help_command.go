package Console

import (
	cdom "github.com/PlayerR9/MyGoLib/CustomData/OrderedMap"
	fsd "github.com/PlayerR9/MyGoLib/FString/Document"
	fsp "github.com/PlayerR9/MyGoLib/FString/Printer"
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

		form := fsp.NewFormatter(
			fsp.NewIndentConfig("   ", 0),
			nil,
			nil,
			nil,
		)

		printer := fsp.NewPrinter(form)

		err := fsp.ApplyFormat(printer, mip)
		if err != nil {
			return nil, ue.NewErrWhile("printing the command list", err)
		}

		doc, err := printer.MakeDocument()
		if err != nil {
			return nil, ue.NewErrWhile("making the document", err)
		}

		return doc, nil
	}

	return NewCommandInfo(
		fsd.NewDocument("Display the help message."),
		fn,
		[]string{},
	), nil
}
