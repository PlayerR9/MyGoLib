package ConsolePanel

// FlagParseResult represents the result of parsing a flag.
type FlagParseResult struct {
	// Parsed arguments.
	Args map[string]any

	// Index of the last unsuccessful parse argument.
	Index int
}

// NewFlagParseResult creates a new FlagParseResult with the given
// arguments, index, and ignorable boolean.
//
// Parameters:
//   - args: The arguments to parse.
//   - index: The index of the last unsuccessful parse argument.
//
// Returns:
//   - *FlagParseResult: A pointer to the new FlagParseResult.
func NewFlagParseResult(args map[string]any, index int) *FlagParseResult {
	return &FlagParseResult{
		Args:  args,
		Index: index,
	}
}
