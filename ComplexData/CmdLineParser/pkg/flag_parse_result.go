package pkg

// FlagParseResult represents the result of parsing a flag.
type FlagParseResult struct {
	// Parsed arguments.
	Args map[string]any
}

// NewFlagParseResult creates a new FlagParseResult with the given
// arguments, index, and ignorable boolean.
//
// Parameters:
//   - args: The arguments to parse.
//
// Returns:
//   - *FlagParseResult: A pointer to the new FlagParseResult.
func NewFlagParseResult(args map[string]any) *FlagParseResult {
	return &FlagParseResult{
		Args: args,
	}
}
