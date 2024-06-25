package CmdLineParser

import (
	pkg "github.com/PlayerR9/MyGoLib/ComplexData/CmdLineParser/pkg"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

type ArgBuilder struct {
	headers []string
	fns     []pkg.ArgumentParserFunc
}

func NewArgBuilder() *ArgBuilder {
	return &ArgBuilder{
		headers: make([]string, 0),
		fns:     make([]pkg.ArgumentParserFunc, 0),
	}
}

// NewArgument creates a new argument.
//
// Parameters:
//   - name: The name of the argument.
//   - argumentParserFunc: The function that parses the argument.
//
// Returns:
//   - *Argument: A pointer to the newly created argument.
//
// Behaviors:
//   - If argumentParserFunc is nil, the default NoArgumentParser is used.
func (b *ArgBuilder) SetArg(name string, parserFunc pkg.ArgumentParserFunc) *ArgBuilder {
	b.headers = append(b.headers, name)
	b.fns = append(b.fns, parserFunc)

	return b
}

func (b *ArgBuilder) Build() ([]*pkg.ArgInfo, error) {
	var arguments []*pkg.ArgInfo

	for i, header := range b.headers {
		arg, err := pkg.NewArgument(header, b.fns[i])
		if err != nil {
			return nil, uc.NewErrAt(i+1, "argument", err)
		}

		arguments = append(arguments, arg)

		index := us.FindEquals(arguments, arg)
		if index != -1 {
			arguments[index] = arg
		} else {
			arguments = append(arguments, arg)
		}
	}

	b.Reset()

	// Format:
	// - "x": only x
	// - "x-y": from x to y (both inclusive)
	// - "x-": from x to +inf
	// - "-y": from 0 to y
	// - "-": any

	return arguments, nil
}

func (b *ArgBuilder) Reset() {
	for i := range b.headers {
		b.fns[i] = nil
	}

	b.headers = b.headers[:0]
	b.fns = b.fns[:0]
}

type CmdBuilder struct {
	names        []string
	descriptions [][]string
	callbacks    []pkg.CommandCallbackFunc
	flagBuilders []*FlagBuilder
}

func NewCmdBuilder() *CmdBuilder {
	return &CmdBuilder{
		names:        make([]string, 0),
		descriptions: make([][]string, 0),
		callbacks:    make([]pkg.CommandCallbackFunc, 0),
		flagBuilders: make([]*FlagBuilder, 0),
	}
}

func (b *CmdBuilder) SetCmd(name string, description []string, callback pkg.CommandCallbackFunc, flagBuilder *FlagBuilder) *CmdBuilder {
	b.names = append(b.names, name)
	b.descriptions = append(b.descriptions, description)
	b.callbacks = append(b.callbacks, callback)
	b.flagBuilders = append(b.flagBuilders, flagBuilder)

	return b
}

func (b *CmdBuilder) Build() ([]*pkg.CommandInfo, error) {
	var commands []*pkg.CommandInfo

	for i, name := range b.names {
		if name == HelpOpcode {
			return nil, uc.NewErrAt(i+1, "command", NewErrOpcodeHelp())
		}

		flagInfos, err := b.flagBuilders[i].Build()
		if err != nil {
			return nil, uc.NewErrAt(i+1, "command", err)
		}

		newCommand, err := pkg.NewCommandInfo(name, b.descriptions[i], b.callbacks[i], flagInfos)
		if err != nil {
			return nil, uc.NewErrAt(i+1, "command", err)
		}

		index := us.FindEquals(commands, newCommand)
		if index != -1 {
			commands[index] = newCommand
		} else {
			commands = append(commands, newCommand)
		}
	}

	b.Reset()

	return commands, nil
}

func (b *CmdBuilder) Reset() {
	b.names = b.names[:0]
	b.descriptions = b.descriptions[:0]

	for i := range b.callbacks {
		b.callbacks[i] = nil

		if b.flagBuilders[i] != nil {
			b.flagBuilders[i].Reset()
			b.flagBuilders[i] = nil
		}
	}

	b.callbacks = b.callbacks[:0]
	b.flagBuilders = b.flagBuilders[:0]
}

type FlagBuilder struct {
	names        []string
	requireds    []bool
	descriptions [][]string
	callbacks    []pkg.FlagCallbackFunc
	argBuilders  []*ArgBuilder
}

func NewFlagBuilder() *FlagBuilder {
	return &FlagBuilder{
		names:        make([]string, 0),
		requireds:    make([]bool, 0),
		descriptions: make([][]string, 0),
		callbacks:    make([]pkg.FlagCallbackFunc, 0),
		argBuilders:  make([]*ArgBuilder, 0),
	}
}

// NewFlagInfo creates a new FlagInfo with the given name and
// arguments.
//
// Parameters:
//   - name: The name of the flag.
//   - isRequired: A boolean indicating whether the flag is required.
//   - callback: The function that parses the flag arguments.
//   - args: A slice of strings representing the arguments accepted by
//     the flag.
//
// Returns:
//   - *FlagInfo: A pointer to the new FlagInfo.
//
// Behaviors:
//   - Any nil arguments are filtered out.
//   - If 'callback' is nil, a default callback is used that returns nil without error.
func (b *FlagBuilder) SetFlag(name string, isRequired bool, description []string, callback pkg.FlagCallbackFunc, argBuilder *ArgBuilder) *FlagBuilder {
	b.names = append(b.names, name)
	b.requireds = append(b.requireds, isRequired)
	b.descriptions = append(b.descriptions, description)
	b.callbacks = append(b.callbacks, callback)
	b.argBuilders = append(b.argBuilders, argBuilder)

	return b
}

func (b *FlagBuilder) Build() ([]*pkg.FlagInfo, error) {
	var flagList []*pkg.FlagInfo

	for i, flag := range b.names {
		var argInfos []*pkg.ArgInfo

		if b.argBuilders[i] == nil {
			argInfos = make([]*pkg.ArgInfo, 0)
		} else {
			var err error

			argInfos, err = b.argBuilders[i].Build()
			if err != nil {
				return nil, uc.NewErrAt(i+1, "flag", err)
			}
		}

		newFlag, err := pkg.NewFlagInfo(
			flag,
			b.descriptions[i],
			b.requireds[i],
			b.callbacks[i],
			argInfos,
		)
		if err != nil {
			return nil, uc.NewErrAt(i+1, "flag", err)
		}

		index := us.FindEquals(flagList, newFlag)
		if index != -1 {
			flagList[index] = newFlag
		} else {
			flagList = append(flagList, newFlag)
		}
	}

	b.Reset()

	return flagList, nil
}

func (b *FlagBuilder) Reset() {
	b.names = b.names[:0]
	b.requireds = b.requireds[:0]
	b.descriptions = b.descriptions[:0]

	for i := range b.callbacks {
		b.callbacks[i] = nil

		if b.argBuilders[i] != nil {
			b.argBuilders[i].Reset()
			b.argBuilders[i] = nil
		}
	}

	b.callbacks = b.callbacks[:0]
	b.argBuilders = b.argBuilders[:0]
}
