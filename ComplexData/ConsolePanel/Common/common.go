package Common

// ArgumentParserFunc is a function type that represents a function
// that parses a string argument.
//
// Parameters:
//   - string: The string to parse.
//
// Returns:
//   - any: The parsed value.
type ArgumentParserFunc func(string) (any, error)

// NoArgumentParser is a default argument parser function that returns
// the string as is.
//
// Parameters:
//   - string: The string to parse.
//
// Returns:
//   - any: The string as is.
//   - error: nil
func NoArgumentParser(s string) (any, error) {
	return s, nil
}

// FlagCallbackFunc is a function type that represents a callback
// function for a console command flag.
//
// Parameters:
//   - argMap: A map of string keys to any values representing the
//     arguments passed to the flag.
//
// Returns:
//   - ArgumentsMap: A map of string keys to any values representing
//     the parsed arguments.
//   - error: An error if the flag fails.
type FlagCallbackFunc func(argMap ArgumentsMap) (ArgumentsMap, error)

// NoFlagCallback is a default callback function for a console command flag
// when no callback is provided.
//
// Parameters:
//   - args: A slice of strings representing the arguments passed to
//     the flag.
//
// Returns:
//   - ArgumentsMap: A map of string keys to any values representing
//     the parsed arguments.
//   - error: nil
func NoFlagCallback(argMap ArgumentsMap) (ArgumentsMap, error) {
	return argMap, nil
}

// CommandCallbackFunc is a function type that represents a callback
// function for a console command.
//
// Parameters:
//   - args: A map of string keys to any values representing the
//     arguments passed to the command.
//
// Returns:
//   - error: An error if the command fails.
type CommandCallbackFunc func(args ArgumentsMap) error

// NoCommandCallback is a default callback function for a console command
// when no callback is provided.
//
// Parameters:
//   - args: A map of string keys to any values representing the
//     arguments passed to the command.
//
// Returns:
//   - error: nil
func NoCommandCallback(args ArgumentsMap) error {
	return nil
}

// ArgumentsMap is a map of string keys to any values representing the
// arguments passed to a command.
//
// Key: The name of the argument.
// Value: The value of the argument.
type ArgumentsMap map[string]any
