package Debugger

// PrintMessager is an interface that defines a message to print.
type PrintMessager interface {
}

// PrintfMessage is a struct that represents a message to print formatted text.
type PrintfMessage struct {
	// format is the format string.
	format string

	// v is the values to print.
	v []any
}

// NewPrintfMessage is a function that creates a new PrintfMessage.
//
// Parameters:
//   - format: The format string.
//   - v: The values to print.
//
// Returns:
//   - *PrintfMessage: The new PrintfMessage.
func NewPrintfMessage(format string, v []any) *PrintfMessage {
	return &PrintfMessage{
		format: format,
		v:      v,
	}
}

// PrintMessage is a struct that represents a message to print.
type PrintMessage struct {
	// v is the values to print.
	v []any
}

// NewPrintMessage is a function that creates a new PrintMessage.
//
// Parameters:
//   - v: The values to print.
//
// Returns:
//   - *PrintMessage: The new PrintMessage.
func NewPrintMessage(v []any) *PrintMessage {
	return &PrintMessage{
		v: v,
	}
}

// PrintlnMessage is a struct that represents a message to print a line.
type PrintlnMessage struct {
	// v is the values to print.
	v []any
}

// NewPrintlnMessage is a function that creates a new PrintlnMessage.
//
// Parameters:
//   - v: The values to print.
//
// Returns:
//   - *PrintlnMessage: The new PrintlnMessage.
func NewPrintlnMessage(v []any) *PrintlnMessage {
	return &PrintlnMessage{
		v: v,
	}
}
