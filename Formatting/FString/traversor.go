package FString

import (
	"fmt"
	"strings"
	"sync"
	"unicode/utf8"

	pkg "github.com/PlayerR9/MyGoLib/Formatting/FString/pkg"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	utstr "github.com/PlayerR9/MyGoLib/Utility/strings"
)

var (
	// NBSP is the non-breaking space rune.
	NBSP rune
)

func init() {
	NBSP = pkg.NBSP
}

/////////////////////////////////////////////////

// Traversor is a type that represents a traversor for a formatted string.
type Traversor struct {
	// config is the configuration of the traversor.
	config *FormatConfig

	// indentation is the string that is used for indentation
	// on the left side of the traversor.
	indentation string

	// hasIndent is a flag that indicates if the traversor has indentation.
	hasIndent bool

	// leftConfig is the configuration for the left symbol of the traversor.
	leftConfig *DelimiterConfig

	// rightDelim is the right delimiter of the traversor.
	rightDelim string

	// source is the buffer of the traversor.
	source *pkg.Buffer

	// mu is the mutex of the traversor.
	mu *sync.Mutex
}

// Cleanup implements the Cleaner interface.
func (trav *Traversor) Clean() {
	trav.source = nil
}

// newTraversor creates a new traversor.
//
// Parameters:
//   - config: The configuration of the traversor.
//   - source: The source of the traversor.
//
// Returns:
//   - *Traversor: The new traversor.
func newTraversor(config *FormatConfig, source *pkg.Buffer) *Traversor {
	trav := &Traversor{
		config:      config,
		source:      source, // shared source
		hasIndent:   false,
		indentation: "",
		leftConfig:  nil,
		rightDelim:  "",
	}

	indentConfig := config.indentation
	if indentConfig != nil {
		trav.indentation = indentConfig.GetIndentation()
		trav.hasIndent = true
	}

	leftConfig := config.delimiterLeft
	if leftConfig != nil {
		trav.leftConfig = leftConfig
	}

	rightConfig := config.delimiterRight
	if rightConfig != nil {
		trav.rightDelim = rightConfig.str
	}

	return trav
}

// writeIndent writes the indentation string to the traversor if
// the traversor has indentation and the traversor is at the first
// of the line.
func (trav *Traversor) writeIndent() {
	ok := trav.source.IsFirstOfLine()
	if !ok {
		return
	}

	if trav.hasIndent {
		trav.source.ForceWriteString(trav.indentation)
	}

	if trav.leftConfig != nil {
		trav.source.ForceWriteString(trav.leftConfig.str)
	}
}

// writeRune appends a rune to the current, in-progress line of the traversor.
//
// Parameters:
//   - r: The rune to append.
//
// Returns:
//   - error: An error if the rune could not be appended.
func (trav *Traversor) writeRune(r rune) error {
	trav.writeIndent()

	switch r {
	case utf8.RuneError:
		return uc.NewErrInvalidRune(nil)
	case NBSP:
		trav.source.WriteRune(r)
	default:
		trav.source.Write(r)
	}

	return nil
}

// writeString appends a string to the current, in-progress line of the traversor.
//
// Parameters:
//   - str: The string to append.
//
// Returns:
//   - int: The index of the first invalid UTF-8 character. -1 if there is no error.
func (trav *Traversor) writeString(str string) int {
	trav.writeIndent()

	if str == "" {
		return -1
	}

	chars, idx := utstr.ToUtf8(str)
	if idx != -1 {
		return idx
	}

	for _, r := range chars {
		trav.source.Write(r)
	}

	return -1
}

// writeLine writes a line to the traversor. If there is any in-progress line,
// then the line is appended to the line before accepting it. Otherwise, a new line
// with the line is added to the source.
//
// Parameters:
//   - line: The line to write.
//
// Returns:
//   - error: An error if the line could not be written.
//
// Behaviors:
//   - If line is empty, then an empty line is added to the source.
func (trav *Traversor) writeLine(line string) error {
	trav.source.AcceptLine(trav.rightDelim) // Accept the current line if any.

	trav.writeIndent()

	if line == "" {
		trav.source.WriteEmptyLine(trav.rightDelim)
	} else {
		chars, idx := utstr.ToUtf8(line)
		if idx != -1 {
			return fmt.Errorf("invalid UTF-8 encoding at index %d", idx)
		}

		for _, r := range chars {
			trav.source.Write(r)
		}
	}

	trav.source.AcceptLine(trav.rightDelim) // Accept the line.

	return nil
}

// AppendRune appends a rune to the half-line of the traversor.
//
// Parameters:
//   - r: The rune to append.
//
// Returns:
//   - error: An error if the rune could not be appended.
//
// common:
//   - *common.ErrInvalidRune: If the rune is invalid.
//
// Behaviors:
//   - If the half-line is nil, then a new half-line is created.
func (trav *Traversor) AppendRune(r rune) error {
	if trav.source == nil {
		return nil
	}

	err := trav.writeRune(r)
	if err != nil {
		return err
	}

	return nil
}

// AppendString appends a string to the half-line of the traversor.
//
// Parameters:
//   - str: The string to append.
//
// Returns:
//   - error: An error if the string is not valid UTF-8.
//
// Behaviors:
//   - IF str is empty: nothing is done.
func (trav *Traversor) AppendString(str string) error {
	if trav.source == nil {
		return nil
	}

	idx := trav.writeString(str)
	if idx != -1 {
		return fmt.Errorf("invalid UTF-8 encoding at index %d", idx)
	}

	return nil
}

// AppendStrings appends multiple strings to the half-line of the traversor.
//
// Parameters:
//   - strs: The strings to append.
//
// Returns:
//   - error: An error of type *common.ErrAt if there is an error appending a string.
//
// Behaviors:
//   - This is equivalent to calling AppendString for each string in strs but more efficient.
func (trav *Traversor) AppendStrings(strs []string) error {
	if trav.source == nil || len(strs) == 0 {
		return nil
	}

	for i, str := range strs {
		idx := trav.writeString(str)
		if idx != -1 {
			return uc.NewErrAt(i, "string", fmt.Errorf("invalid UTF-8 encoding at index %d", idx))
		}
	}

	return nil
}

// AppendJoinedString appends a joined string to the half-line of the traversor.
//
// Parameters:
//   - sep: The separator to use.
//   - fields: The fields to join.
//
// Returns:
//   - error: An error if some field or the separator is not valid UTF-8 encoding.
//
// Behaviors:
//   - This is equivalent to calling AppendString(strings.Join(fields, sep)).
func (trav *Traversor) AppendJoinedString(sep string, fields ...string) error {
	if trav.source == nil || len(fields) == 0 {
		return nil
	}

	str := strings.Join(fields, sep)

	idx := trav.writeString(str)
	if idx != -1 {
		return fmt.Errorf("invalid UTF-8 encoding at index %d", idx)
	}

	return nil
}

// AcceptWord is a function that, if there is any in-progress word, then said word is added
// to the source.
func (trav *Traversor) AcceptWord() {
	if trav.source == nil {
		return
	}

	trav.source.AcceptWord()
}

// AcceptLine is a function that accepts the current line of the traversor.
//
// Behaviors:
//   - This also accepts the current word if any.
func (trav *Traversor) AcceptLine() {
	if trav.source == nil {
		return
	}

	trav.source.AcceptLine(trav.rightDelim)
}

// AddLine adds a line to the traversor. If there is any in-progress line, then the line is
// appended to the line before accepting it. Otherwise, a new line with the line is added to
// the source.
//
// Parameters:
//   - line: The line to add.
//
// Returns:
//   - error: An error of type *common.ErrAt if there is an error adding the line.
//
// Behaviors:
//   - If line is empty, then an empty line is added to the source.
func (trav *Traversor) AddLine(line string) error {
	if trav.source == nil {
		return nil
	}

	err := trav.writeLine(line)
	if err != nil {
		return err
	}

	return nil
}

// AddLines adds multiple lines to the traversor in a more efficient way than
// adding each line individually.
//
// Parameters:
//   - lines: The lines to add.
//
// Returns:
//   - error: An error of type *common.ErrAt if there is an error adding a line.
//
// Behaviors:
//   - If there are no lines, then nothing is done.
func (trav *Traversor) AddLines(lines []string) error {
	if trav.source == nil || len(lines) == 0 {
		return nil
	}

	for i, line := range lines {
		err := trav.writeLine(line)
		if err != nil {
			return uc.NewErrAt(i, "line", err)
		}
	}

	return nil
}

// AddJoinedLine adds a joined line to the traversor. This is a more efficient way to do
// the same as AddLine(strings.Join(fields, sep)).
//
// Parameters:
//   - sep: The separator to use.
//   - fields: The fields to join.
//
// Returns:
//   - error: An error of type *common.ErrInvalidRuneAt if there is an invalid rune
//     in the line.
//
// Behaviors:
//   - If fields is empty, then nothing is done.
func (trav *Traversor) AddJoinedLine(sep string, fields ...string) error {
	if trav.source == nil || len(fields) == 0 {
		return nil
	}

	str := strings.Join(fields, sep)

	err := trav.writeLine(str)
	if err != nil {
		return err
	}

	return nil
}

// EmptyLine adds an empty line to the traversor. This is a more efficient way to do
// the same as AddLine("") or AddLines([]string{""}).
//
// Behaviors:
//   - If the half-line is not empty, then the half-line is added to the source
//     (half-line is reset) and an empty line is added to the source.
func (trav *Traversor) EmptyLine() {
	if trav.source == nil {
		return
	}

	trav.source.AcceptLine(trav.rightDelim) // Accept the current line if any.

	trav.writeIndent()

	trav.source.ForceAcceptLine(trav.rightDelim) // Accept the line.
}

// Write implements the io.Writer interface for the traversor.
func (trav *Traversor) Write(p []byte) (int, error) {
	if trav.source == nil {
		return 0, nil
	}

	n, err := trav.source.WriteBytes(p)
	if err != nil {
		return n, err
	}

	return len(p), nil
}

// Print is a function that writes to the traversor using the fmt.Fprint function.
//
// Parameters:
//   - a: The arguments to write.
func (trav *Traversor) Print(a ...interface{}) error {
	if trav.source == nil {
		return nil
	}

	_, err := fmt.Fprint(trav, a...)
	return err
}

// Printf is a function that writes to the traversor using the fmt.Fprintf function.
//
// Parameters:
//   - format: The format string.
//   - a: The arguments to write.
func (trav *Traversor) Printf(format string, a ...interface{}) error {
	if trav.source == nil {
		return nil
	}

	_, err := fmt.Fprintf(trav, format, a...)
	return err
}

// Println is a function that writes to the traversor using the fmt.Fprintln function.
//
// Parameters:
//   - a: The arguments to write.
func (trav *Traversor) Println(a ...interface{}) error {
	if trav.source == nil {
		return nil
	}

	_, err := fmt.Fprintln(trav, a...)
	return err
}

// GetConfig is a method that returns a copy of the configuration of the traversor.
//
// Parameters:
//   - options: The options to apply to the configuration.
//
// Returns:
//   - FormatConfig: A copy of the configuration of the traversor.
func (trav *Traversor) GetConfig(options ...ConfigOption) *FormatConfig {
	configCopy := trav.config.Copy().(*FormatConfig)

	for _, option := range options {
		option(configCopy)
	}

	return configCopy
}

// Lock locks the traversor. Be aware of deadlocks.
func (trav *Traversor) Lock() {
	trav.mu.Lock()
}

// Unlock unlocks the traversor. Be aware of deadlocks.
func (trav *Traversor) Unlock() {
	trav.mu.Unlock()
}

//////////////////////////////////////////////////////////////

/*
// GetIndent returns the indentation string of the traversor.
//
// Returns:
//   - string: The indentation string of the traversor.
func (trav *Traversor) GetIndent() string {
	if trav.indent == nil {
		return ""
	} else {
		return trav.indentStr
	}
}

// ApplyIndent applies the indentation configuration to a specified string.
//
// Parameters:
//   - str: The string to apply the indentation to.
//
// Returns:
//   - string: The string with the indentation applied.
func (trav *Traversor) ApplyIndent(isFirstLine bool, str string) string {
	if trav.indent == nil || !trav.source.isFirstOfLine() {
		return str
	}

	var builder strings.Builder

	builder.WriteString(trav.indentStr)
	builder.WriteString(str)

	return builder.String()
}
*/

/*
// AddMultiline adds a multiline to the traversor. But first, it accepts any in-progress
// half-line.
//
// Parameters:
//   - mlt: The multiline to add.
//
// Behaviors:
//   - If the multiline is nil, then nothing is done.
func (trav *Traversor) AddMultiline(mlt *cb.MultiLineText) {
	if mlt == nil {
		return
	}

	trav.AcceptHalfLine()
	trav.source.addLine(mlt)
}
*/
