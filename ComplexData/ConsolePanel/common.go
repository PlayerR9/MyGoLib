package ConsolePanel

import (
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
)

const (
	// HelpOpcode is the opcode for the help command.
	HelpOpcode string = "help"
)

// BoolFS returns "Yes" if val is true, "No" otherwise.
//
// Parameters:
//   - val: The boolean value.
//
// Returns:
//   - string: "Yes" if val is true, "No" otherwise.
//   - error: nil
func BoolFString(val bool) (string, error) {
	var res string

	if val {
		res = "Yes"
	} else {
		res = "No"
	}

	return res, nil
}

type descriptionPrinter struct {
	lines []string
}

func (dp *descriptionPrinter) FString(trav *ffs.Traversor) error {
	if trav == nil {
		return nil
	}

	err := trav.AddLines(dp.lines)
	if err != nil {
		return err
	}

	return nil
}
