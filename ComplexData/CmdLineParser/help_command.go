package CmdLineParser

import (
	"fmt"
	"strconv"

	pkg "github.com/PlayerR9/MyGoLib/ComplexData/CmdLineParser/pkg"
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

func generateHelpDoc(opcode string) ([]string, error) {
	doc, err := ffs.Sprintln(
		ffs.DefaultFormatter,
		"Displays help information for the",
		strconv.Quote(opcode),
		"command.",
	)
	if err != nil {
		return nil, fmt.Errorf("error generating help doc: %w", err)
	}

	return ffs.Stringfy(doc), nil
}

func generateHelpCallback(cpnl *CmdLineParser) pkg.CommandCallbackFunc {
	displayCommand := func(command *pkg.CommandInfo) ([]string, error) {
		doc, err := ffs.SprintFString(ffs.DefaultFormatter, command)
		if err != nil {
			return nil, fmt.Errorf("error printing command %q: %w", command.GetOpcode(), err)
		}

		return ffs.Stringfy(doc), nil
	}

	return func(args map[string]map[string][]any) (any, error) {
		if len(args) == 0 {
			doc, err := ffs.SprintFString(ffs.DefaultFormatter, cpnl)
			if err != nil {
				return nil, fmt.Errorf("error printing console panel: %w", err)
			}

			return ffs.Stringfy(doc), nil
		} else {
			var pages [][]string

			for opcode := range args {
				command, ok := cpnl.GetCommand(opcode)
				if !ok {
					return nil, NewErrCommandNotFound(opcode)
				}

				page, err := displayCommand(command)
				if err != nil {
					return nil, err
				}

				pages = append(pages, page)
			}

			return pages, nil
		}
	}
}

func generateHelpCommand(console *CmdLineParser) (*pkg.CommandInfo, error) {
	// Add the help command
	doc, err := ffs.Sprintln(ffs.DefaultFormatter, "Displays help information for the console.")
	if err != nil {
		return nil, fmt.Errorf("error generating help doc: %w", err)
	}

	flagDoc, err := generateHelpDoc(HelpOpcode)
	if err != nil {
		return nil, err
	}

	flag, err := pkg.NewFlagInfo(
		HelpOpcode,
		flagDoc,
		false,
		nil,
		make([]*pkg.ArgInfo, 0),
	)
	if err != nil {
		return nil, err
	}

	flagList := []*pkg.FlagInfo{flag}

	for i, command := range console.commandList {
		doc, err := generateHelpDoc(command.GetOpcode())
		if err != nil {
			panic(err)
		}

		flag, err := pkg.NewFlagInfo(
			command.GetOpcode(),
			doc,
			false,
			nil,
			make([]*pkg.ArgInfo, 0),
		)
		if err != nil {
			return nil, ue.NewErrAt(i+1, "flag", err)
		}

		flagList = append(flagList, flag)
	}

	command, err := pkg.NewCommandInfo(
		HelpOpcode,
		ffs.Stringfy(doc),
		generateHelpCallback(console),
		flagList,
	)
	if err != nil {
		return nil, err
	}

	return command, nil
}
