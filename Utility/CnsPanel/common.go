package CnsPanel

// CommandCallbackFunc is a function type that represents a callback
// function for a console command.
//
// Parameters:
//   - args: A map of string keys to any values representing the
//     arguments passed to the command.
//
// Returns:
//   - error: An error if the command fails.
type CommandCallbackFunc func(args map[string]any) error

// FlagCallbackFunc is a function type that represents a callback
// function for a console command flag.
//
// Parameters:
//   - args: A slice of strings representing the arguments passed to
//     the flag.
//
// Returns:
//   - any: Any value representing the result of the flag callback.
//   - error: An error if the flag fails.
type FlagCallbackFunc func(...string) (any, error)

// Arguments is a map of string keys to any values representing the
// arguments passed to a command.
//
// Key: The name of the argument.
// Value: The value of the argument.
type Arguments map[string]any
