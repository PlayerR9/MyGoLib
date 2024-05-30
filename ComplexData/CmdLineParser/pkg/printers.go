package pkg

import (
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
)

type DescriptionPrinter struct {
	lines []string
}

func (dp *DescriptionPrinter) FString(trav *ffs.Traversor) error {
	if trav == nil {
		return nil
	}

	err := trav.AddLines(dp.lines)
	if err != nil {
		return err
	}

	return nil
}

func NewDescriptionPrinter(lines []string) *DescriptionPrinter {
	return &DescriptionPrinter{
		lines: lines,
	}
}
