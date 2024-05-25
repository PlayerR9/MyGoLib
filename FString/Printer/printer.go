package Printer

import (
	"unicode/utf8"

	fsd "github.com/PlayerR9/MyGoLib/FString/Document"
	fss "github.com/PlayerR9/MyGoLib/FString/Section"
	ue "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// Printer is a type that represents a formatted string.
type Printer struct {
	// pages are the pages of the document.
	pages []*pageBuilder

	// lastPage is the last page of the document.
	lastPage int

	// formatter is the formatter of the document.
	formatter *Formatter
}

// NewPrinter creates a new printer.
//
// Parameters:
//   - formatter: The formatter to use.
//
// Returns:
//   - *Printer: The new printer.
//
// Behaviors:
//   - If the formatter is nil, the function uses the formatter with nil values.
func NewPrinter(formatter *Formatter) *Printer {
	printer := &Printer{
		pages:    []*pageBuilder{NewPage()},
		lastPage: 0,
	}

	if formatter == nil {
		printer.formatter = NewFormatter(nil, nil, nil, nil)
	} else {
		printer.formatter = formatter
	}

	return printer
}

// ApplyFormat applies a format to a stringer.
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
func ApplyFormat[T FStringer](p *Printer, elem T) error {
	if p == nil {
		return ue.NewErrNilParameter("p")
	}

	trav := newTraversor(p.formatter, p)

	err := elem.FString(trav)
	if err != nil {
		return err
	}

	return nil
}

// ApplyFormatMany applies a format to a stringer.
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
func ApplyFormatMany[T FStringer](p *Printer, elems []T) error {
	if len(elems) == 0 {
		return nil
	}

	if p == nil {
		return ue.NewErrNilParameter("p")
	}

	for i, elem := range elems {
		err := elem.FString(newTraversor(p.formatter, p))
		if err != nil {
			return ue.NewErrAt(i+1, "FStringer element", err)
		}
	}

	return nil
}

// ApplyFunc applies a function to the printer. Useful for when you want to apply a function
// that does not implement the FStringer interface.
//
// Parameters:
//   - trav: The traversor to use.
//   - elem: The element to apply the function to.
//   - f: The function to apply.
//
// Returns:
//   - error: An error if the function fails.
//
// Errors:
//   - *ErrFinalization: If the finalization of the page fails.
//   - any error returned by the function.
func ApplyFunc[T any](trav *Traversor, elem T, f FStringFunc[T]) error {
	err := f(trav, elem)
	if err != nil {
		return err
	}

	return nil
}

// ApplyFuncMany applies a function to the printer. Useful for when you want to apply a function
// that does not implement the FStringer interface.
//
// Parameters:
//   - trav: The traversor to use.
//   - f: The function to apply.
//   - elems: The elements to apply the function to.
//
// Returns:
//   - error: An error if the function fails.
//
// Errors:
//   - *ErrFinalization: If the finalization of the page fails.
//   - *Errors.ErrAt: If an error occurs on a specific element.
//   - any error returned by the function.
func ApplyFuncMany[T any](trav *Traversor, f FStringFunc[T], elems []T) error {
	if len(elems) == 0 {
		return nil
	}

	for i, elem := range elems {
		err := f(trav, elem)
		if err != nil {
			return ue.NewErrAt(i+1, "element", err)
		}
	}

	return nil
}

// ApplyFormatFunc applies a format function to the printer.
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
func ApplyFormatFunc[T any](p *Printer, elem T, f FStringFunc[T]) error {
	if p == nil {
		return ue.NewErrNilParameter("p")
	}

	trav := newTraversor(p.formatter, p)

	err := f(trav, elem)
	if err != nil {
		return err
	}

	return nil
}

// ApplyFormatFuncMany applies a format function to the printer.
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
func ApplyFormatFuncMany[T any](p *Printer, f FStringFunc[T], elems []T) error {
	if len(elems) == 0 {
		return nil
	}

	if p == nil {
		return ue.NewErrNilParameter("p")
	}

	for i, elem := range elems {
		err := f(newTraversor(p.formatter, p), elem)
		if err != nil {
			return ue.NewErrAt(i+1, "element", err)
		}
	}

	return nil
}

// MakeDocument is a function that makes a document from the printer.
//
// Returns:
//   - *Document: The document.
func (p *Printer) MakeDocument() (*fsd.Document, error) {
	err := p.pages[p.lastPage].Finalize(new(fss.MultilineText))
	if err != nil {
		return nil, NewErrFinalization(err)
	}

	doc := fsd.NewDocument()

	return doc, nil

	// FIXME: This is a temporary fix.

	/*
		for _, page := range p.pages {
			iter := page.Iterator()

			for {
				section, err := iter.Consume()
				if err != nil {
					break
				}
			}
		}

		return doc, nil
	*/
}

// isFirstOfLine is a private function that returns true if the current position is the first
// position of a line.
//
// Returns:
//   - bool: True if the current position is the first position of a line.
func (p *Printer) isFirstOfLine() bool {
	return p.pages[p.lastPage].IsFirstOfLine()
}

// writeIndent is a private function that writes the indentation to the formatted string.
func (p *Printer) writeIndent(str string) {
	p.pages[p.lastPage].AddString(str)
}

// appendNBSP is a private function that appends a non-breaking space to the formatted
// string.
//
// Behaviors:
//   - This is equivalent to calling appendRune('\u00A0') but much more efficient.
func (p *Printer) appendNBSP() {
	// Non-breaking space
	p.pages[p.lastPage].AddRune(' ')
}

// appendRune is a private function that appends a rune to the formatted string.
//
// Parameters:
//   - replaceTab: The string to replace tabs with.
//   - char: The rune to append.
//
// Returns:
//   - error: An error if the rune could not be appended.
//
// Errors:
//   - *Errors.ErrInvalidRune: If the rune is invalid.
//   - *ErrFinalization: If the finalization of the page fails.
//   - *ErrInvalidPage: If the page is invalid.
func (p *Printer) appendRune(replaceTab string, char rune) error {
	if char == utf8.RuneError {
		err := p.pages[p.lastPage].Finalize(new(fss.MultilineText))
		if err != nil {
			return NewErrFinalization(err)
		}

		return ue.NewErrInvalidRune(nil)
	}

	switch char {
	case '\t':
		// Replace tabs with the replaceTab string
		p.pages[p.lastPage].AddString(replaceTab)
	case '\v':
		// Do nothing
	case '\r', '\n', '\u0085':
		// Line feed
		p.pages[p.lastPage].AcceptLine()
	case '\f':
		// Form feed
		err := p.pages[p.lastPage].Accept(new(fss.MultilineText))
		if err != nil {
			return NewErrInvalidPage(err)
		}
	case ' ':
		// Space
		p.pages[p.lastPage].AcceptWord()
	case '\u00A0':
		// Non-breaking space
		p.pages[p.lastPage].AddRune(' ')
	default:
		p.pages[p.lastPage].AddRune(char)
	}

	return nil
}

// appendString is a private function that appends a string to the formatted string.
//
// Parameters:
//   - replaceTab: The string to replace tabs with.
//   - str: The string to append.
//
// Returns:
//   - error: An error of type *Errors.ErrAt if an invalid rune is found in the string.
//
// Behaviors:
//   - If the string is empty, nothing is done.
func (p *Printer) appendString(replaceTab, str string) error {
	if str == "" {
		return nil
	}

	for j := 0; len(str) > 0; j++ {
		char, size := utf8.DecodeRuneInString(str)
		str = str[size:]

		if char == '\r' && size != 0 {
			nextRune, size := utf8.DecodeRuneInString(str)

			if nextRune == '\n' {
				str = str[size:]
			}
		}

		err := p.appendRune(replaceTab, char)
		if err != nil {
			return ue.NewErrAt(j, "character", err)
		}
	}

	return nil
}

// acceptWord is a private function that accepts the current word of the formatted string.
func (p *Printer) acceptWord() {
	p.pages[p.lastPage].AcceptWord()
}

// acceptLine is a private function that accepts the current line of the formatted string.
func (p *Printer) acceptLine() {
	p.pages[p.lastPage].AcceptLine()
}

// GetRaw returns the raw content of the printer.
//
// Returns:
//   - [][]fss.Sectioner: The raw content of the printer.
//   - error: An error if the finalization of the page fails.
//
// Errors:
//   - *ErrFinalization: If the finalization of the page fails.
func (p *Printer) GetRaw() ([][]fss.Sectioner, error) {
	err := p.pages[p.lastPage].Finalize(new(fss.MultilineText))
	if err != nil {
		return nil, NewErrFinalization(err)
	}

	var raw [][]fss.Sectioner

	for _, page := range p.pages {
		raw = append(raw, page.sections)
	}

	return raw, nil
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
