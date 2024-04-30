// Package CnsPanel provides a structure and functions for handling
// console command flags.
package CnsPanel

// ConsoleBuilder represents a builder for a command line.
type ConsoleBuilder struct {
	// Name of the executable.
	execName string

	// Description of the executable.
	description [][]string

	// List of commands accepted by the console.
	commands []*CommandInfo
}

// SetExecutableName is a method of ConsoleBuilder that sets the
// executable name for a ConsoleBuilder.
//
// Parameters:
//   - name: The name of the executable.
func (b *ConsoleBuilder) SetExecutableName(name string) {
	b.execName = name
}

// AppendParagraph is a method of ConsoleBuilder that appends a
// paragraph to the description of a ConsoleBuilder.
//
// Parameters:
//   - contents: The contents of the paragraph to append.
func (b *ConsoleBuilder) AppendParagraph(contents ...string) {
	if b.description == nil {
		b.description = [][]string{contents}
	} else {
		b.description = append(b.description, contents)
	}
}

// AddCommand is a method of ConsoleBuilder that adds a new command
// to a ConsoleBuilder.
//
// It overwrites the command if it already exists.
//
// Parameters:
//   - command: The command to add.
func (b *ConsoleBuilder) AddCommand(command *CommandInfo) {
	if b.commands == nil {
		b.commands = []*CommandInfo{command}
	} else {
		b.commands = append(b.commands, command)
	}
}

// Build is a method of ConsoleBuilder that builds a ConsolePanel from a
// ConsoleBuilder.
//
// Returns:
//
//   - *ConsolePanel: A ConsolePanel built from the ConsoleBuilder.
func (b *ConsoleBuilder) Build() ConsolePanel {
	cm := ConsolePanel{
		executableName: b.execName,
		commands:       make([]*CommandInfo, 0),
	}

	if b.commands == nil {
		cm.description = make([][]string, 0)
	} else {
		// Remove duplicate commands prioritizing the last one added.
		seen := make(map[string]*CommandInfo)

		for _, command := range b.commands {
			seen[command.name] = command
		}

		for _, command := range seen {
			cm.commands = append(cm.commands, command)
		}

		cm.description = b.description
	}

	// Clear the ConsoleBuilder
	for i := range b.description {
		b.description[i] = nil
	}

	b.description = nil

	b.commands = nil

	return cm
}

// Reset is a method of ConsoleBuilder that resets a ConsoleBuilder.
func (b *ConsoleBuilder) Reset() {
	b.execName = ""

	for i := range b.description {
		b.description[i] = nil
	}

	b.description = nil
	b.commands = nil
}
