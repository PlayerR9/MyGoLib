package FString

import (
	"strings"

	ue "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// Printer is a type that represents a formatted string.
type Printer struct {
	// buffer is the buffer of the document.
	buff *buffer

	// formatter is the formatter of the document.
	formatter FormatConfig
}

// NewPrinter creates a new printer.
//
// Parameters:
//   - form: The formatter to use.
//
// Returns:
//   - *Printer: The new printer.
//
// Behaviors:
//   - If the formatter is nil, the function uses the formatter with nil values.
func NewPrinter(form FormatConfig) *Printer {
	return &Printer{
		buff:      newBuffer(),
		formatter: form,
	}
}

// NewPrinterFromConfig creates a new printer from a configuration.
//
// Parameters:
//   - opts: The configuration to use.
//
// Returns:
//   - *Printer: The new printer.
//
// Behaviors:
//   - If the configuration is nil, the function uses the default configuration.
//   - Panics if an invalid configuration type is given (i.e., not IndentConfig, DelimiterConfig,
//     or SeparatorConfig).
func NewPrinterFromConfig(opts ...Configer) *Printer {
	return &Printer{
		buff:      newBuffer(),
		formatter: NewFormatter(opts...),
	}
}

// GetTraversor returns a traversor for the printer.
//
// Returns:
//   - *Traversor: The traversor for the printer.
func (p *Printer) GetTraversor() *Traversor {
	return newTraversor(p.formatter, p.buff)
}

// Apply applies a format to a stringer.
//
// Parameters:
//   - p: The printer to use.
//   - elem: The element to format.
//
// Returns:
//   - error: An error if the formatting fails.
//
// Errors:
//   - *ErrInvalidParameter: If the printer is nil.
//   - *ErrFinalization: If the finalization of the page fails.
//   - any error returned by the element's FString method.
//
// Behaviors:
//   - If the formatter is nil, the function uses the nil formatter.
//   - If the element is nil, the function does nothing.
func Apply[T FStringer](p *Printer, elem T) error {
	if p == nil {
		return ue.NewErrNilParameter("p")
	}

	trav := newTraversor(p.formatter, p.buff)

	err := elem.FString(trav)
	if err != nil {
		return err
	}

	return nil
}

// ApplyMany applies a format to a stringer.
//
// Parameters:
//   - p: The printer to use.
//   - elems: The elements to format.
//
// Returns:
//   - error: An error if the formatting fails.
//
// Errors:
//   - *ErrInvalidParameter: If the printer is nil.
//   - *ErrFinalization: If the finalization of the page fails.
//   - *Errors.ErrAt: If an error occurs on a specific element.
//
// Behaviors:
//   - If the formatter is nil, the function uses the nil formatter.
//   - If an element is nil, the function skips the element.
//   - If all elements are nil, the function does nothing.
func ApplyMany[T FStringer](p *Printer, elems []T) error {
	if len(elems) == 0 {
		return nil
	}

	if p == nil {
		return ue.NewErrNilParameter("p")
	}

	for i, elem := range elems {
		err := elem.FString(newTraversor(p.formatter, p.buff))
		if err != nil {
			return ue.NewErrAt(i+1, "FStringer element", err)
		}
	}

	return nil
}

// ApplyFunc applies a format function to the printer.
//
// Parameters:
//   - p: The printer to use.
//   - elem: The element to apply the function to.
//   - f: The function to apply.
//
// Returns:
//   - error: An error if the function fails.
//
// Errors:
//   - *ErrInvalidParameter: If the printer is nil.
//   - any error returned by the function.
func ApplyFunc[T any](p *Printer, elem T, f FStringFunc[T]) error {
	if p == nil {
		return ue.NewErrNilParameter("p")
	}

	trav := newTraversor(p.formatter, p.buff)

	err := f(trav, elem)
	if err != nil {
		return err
	}

	return nil
}

// ApplyFuncMany applies a format function to the printer.
//
// Parameters:
//   - p: The printer to use.
//   - f: The function to apply.
//   - elems: The elements to apply the function to.
//
// Returns:
//   - error: An error if the function fails.
//
// Errors:
//   - *ErrInvalidParameter: If the printer is nil.
//   - *Errors.ErrAt: If an error occurs on a specific element.
//   - any error returned by the function.
func ApplyFuncMany[T any](p *Printer, f FStringFunc[T], elems []T) error {
	if len(elems) == 0 {
		return nil
	}

	if p == nil {
		return ue.NewErrNilParameter("p")
	}

	for i, elem := range elems {
		err := f(newTraversor(p.formatter, p.buff), elem)
		if err != nil {
			return ue.NewErrAt(i+1, "element", err)
		}
	}

	return nil
}

// GetPages returns the pages of the printer.
//
// Returns:
//   - [][][][]string: The pages of the printer.
func (p *Printer) GetPages() [][][][]string {
	p.buff.finalize()

	pages := p.buff.pages

	// Reset the buffer
	p.buff = newBuffer()

	allStrings := make([][][][]string, 0, len(pages))

	for _, page := range pages {
		sectionLines := make([][][]string, 0)

		for _, section := range page {
			sectionLines = append(sectionLines, section.getLines())
		}

		allStrings = append(allStrings, sectionLines)
	}

	return allStrings
}

// Cleanup implements the Cleaner interface.
func (p *Printer) Cleanup() {
	p.buff.Cleanup()

	p.buff = nil
}

// PrintFString prints a formatted string.
//
// Parameters:
//   - form: The formatter to use.
//   - elem: The element to print.
//
// Returns:
//   - [][][][]string: The pages of the formatted string.
func PrintFString[T FStringer](form FormatConfig, elem T) ([][][][]string, error) {
	p := NewPrinter(form)

	err := Apply(p, elem)
	if err != nil {
		return nil, err
	}

	return p.GetPages(), nil
}

// Printc prints a character.
//
// Parameters:
//   - form: The formatter to use.
//   - c: The character to print.
//
// Returns:
//   - [][][][]string: The pages of the formatted character.
//   - error: An error if the printing fails.
func Printc(form FormatConfig, c rune) ([][][][]string, error) {
	p := NewPrinter(form)

	trav := newTraversor(form, p.buff)

	err := trav.AppendRune(c)
	if err != nil {
		return nil, err
	}

	return p.GetPages(), nil
}

// Print prints a string.
//
// Parameters:
//   - form: The formatter to use.
//   - strs: The strings to print.
//
// Returns:
//   - [][][][]string: The pages of the formatted strings.
//   - error: An error if the printing fails.
func Print(form FormatConfig, strs ...string) ([][][][]string, error) {
	p := NewPrinter(form)

	trav := newTraversor(form, p.buff)

	var err error

	switch len(strs) {
	case 0:
		// Do nothing
	case 1:
		err = trav.AppendString(strs[0])
	default:
		err = trav.AppendStrings(strs)
	}

	if err != nil {
		return nil, err
	}

	// apply

	return p.GetPages(), nil
}

// Printj prints a joined string.
//
// Parameters:
//   - form: The formatter to use.
//   - sep: The separator to use.
//   - strs: The strings to join.
//
// Returns:
//   - [][][][]string: The pages of the formatted strings.
//   - error: An error if the printing fails.
func Printj(form FormatConfig, sep string, strs ...string) ([][][][]string, error) {
	p := NewPrinter(form)

	trav := newTraversor(form, p.buff)

	err := trav.AppendJoinedString(sep, strs...)
	if err != nil {
		return nil, err
	}

	return p.GetPages(), nil
}

// Fprint prints a formatted string.
//
// Parameters:
//   - form: The formatter to use.
//   - a: The elements to print.
//
// Returns:
//   - [][][][]string: The pages of the formatted strings.
//   - error: An error if the printing fails.
func Fprint(form FormatConfig, a ...interface{}) ([][][][]string, error) {
	p := NewPrinter(form)

	trav := newTraversor(form, p.buff)

	err := trav.Print()
	if err != nil {
		return nil, err
	}

	return p.GetPages(), nil
}

// Fprintf prints a formatted string.
//
// Parameters:
//   - form: The formatter to use.
//   - format: The format string.
//   - a: The elements to print.
//
// Returns:
//   - [][][][]string: The pages of the formatted strings.
//   - error: An error if the printing fails.
func Fprintf(form FormatConfig, format string, a ...interface{}) ([][][][]string, error) {
	p := NewPrinter(form)

	trav := newTraversor(form, p.buff)

	err := trav.Printf(format, a...)
	if err != nil {
		return nil, err
	}

	return p.GetPages(), nil
}

// Println prints a string with a newline.
//
// Parameters:
//   - form: The formatter to use.
//   - lines: The lines to print.
//
// Returns:
//   - [][][][]string: The pages of the formatted strings.
//   - error: An error if the printing fails.
func Println(form FormatConfig, lines ...string) ([][][][]string, error) {
	p := NewPrinter(form)

	trav := newTraversor(form, p.buff)

	var err error

	switch len(lines) {
	case 0:
		trav.EmptyLine()
	case 1:
		err = trav.AddLine(lines[0])
	default:
		err = trav.AddLines(lines)
	}

	if err != nil {
		return nil, err
	}

	return p.GetPages(), nil
}

// Printjln prints a joined string with a newline.
//
// Parameters:
//   - form: The formatter to use.
//   - sep: The separator to use.
//   - lines: The lines to join.
//
// Returns:
//   - [][][][]string: The pages of the formatted strings.
//   - error: An error if the printing fails.
func Printjln(form FormatConfig, sep string, lines ...string) ([][][][]string, error) {
	p := NewPrinter(form)

	trav := newTraversor(form, p.buff)

	err := trav.AddJoinedLine(sep, lines...)
	if err != nil {
		return nil, err
	}

	return p.GetPages(), nil
}

// Fprintln prints a formatted string with a newline.
//
// Parameters:
//   - form: The formatter to use.
//   - a: The elements to print.
//
// Returns:
//   - [][][][]string: The pages of the formatted strings.
//   - error: An error if the printing fails.
func Fprintln(form FormatConfig, a ...interface{}) ([][][][]string, error) {
	p := NewPrinter(form)

	trav := newTraversor(form, p.buff)

	err := trav.Println(a...)
	if err != nil {
		return nil, err
	}

	return p.GetPages(), nil
}

// Stringify converts a formatted string to a string.
//
// Parameters:
//   - doc: The formatted string.
//
// Returns:
//   - [][]string: The stringified formatted string.
func Stringfy(doc [][][][]string) []string {
	var pages []string

	for _, page := range doc {
		var sections []string

		for _, section := range page {
			var lines []string

			for _, line := range section {
				lines = append(lines, strings.Join(line, " "))
			}

			sections = append(sections, strings.Join(lines, "\n"))
		}

		pages = append(pages, strings.Join(sections, "\n"))
	}

	return pages
}

//////////////////////////////////////////////////////////////
/*
const (
	// Hellip is the ellipsis character.
	Hellip string = "..."

	// HellipLen is the length of the ellipsis character.
	HellipLen int = len(Hellip)

	// MarginLeft is the left margin of the content box.
	MarginLeft int = 1
)


// addLine is a private function that adds a line to the formatted string.
//
// Parameters:
//   - mlt: The line to add.
func (p *printerSource) addLine(mlt *cb.MultiLineText) {
	if mlt == nil {
		return
	}

	p.lines = append(p.lines, mlt)
}

// GetLines returns the lines of the formatted string.
//
// Returns:
//   - []*MultiLineText: The lines of the formatted string.
func (p *printerSource) GetLines() []*cb.MultiLineText {
	return p.lines
}

/*
func (p *printerSource) Boxed(width, height int) ([]string, error) {
	p.fix()

	all_fields := p.getAllFields()

	fss := make([]*printerSource, 0, len(all_fields))

	for _, fields := range all_fields {
		p := &printerSource{
			lines: fields,
		}

		fss = append(fss, p)
	}

	lines := make([]string, 0)

	for _, p := range fss {
		ts, err := p.generateContentBox(width, height)
		if err != nil {
			return nil, err
		}

		leftLimit, ok := ts.GetFurthestRightEdge()
		if !ok {
			panic("could not get furthest right edge")
		}

		for _, line := range ts.GetLines() {
			fitted, err := sext.FitString(line.String(), leftLimit)
			if err != nil {
				return nil, err
			}

			lines = append(lines, fitted)
		}
	}

	return lines, nil
}


func (p *printerSource) fix() {
	// 1. Fix newline boundaries
	newLines := make([]string, 0)

	for _, line := range p.lines {
		newFields := strings.Split(line, "\n")

		newLines = append(newLines, newFields...)
	}

	p.lines = newLines
}

// Must call Fix() before calling this function.
func (p *printerSource) getAllFields() [][]string {
	// TO DO: Handle special WHITESPACE characters

	fieldList := make([][]string, 0)

	for _, content := range p.lines {
		fields := strings.Fields(content)

		if len(fields) != 0 {
			fieldList = append(fieldList, fields)
		}
	}

	return fieldList
}
*/

/*
// GetDocument returns the content of the FieldSplitter as a Document.
//
// Returns:
//   - *Document: The content of the FieldSplitter.
func (p *FieldSplitter) GetDocument() *FieldSplitter {
	return p.content
}


// Build is a function that builds the document.
//
// Returns:
//   - *tld.Document: The built document.
func (p *FieldSplitter) Build() *tld.Document {
	doc := tld.NewDocument()

	for _, page := range p.content.pages {
		iter := page.Iterator()

		for {
			section, err := iter.Consume()
			if err != nil {
				break
			}
		}
	}

	return doc
}
*/
