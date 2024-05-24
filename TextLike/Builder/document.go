package Builder

import ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"

// DocumentBuilder is a generic data structure that represents a document.
type DocumentBuilder struct {
	// pages are the pages of the document.
	pages []*PageBuilder

	// lastPage is the last page of the document.
	lastPage int
}

// Accept is a function that accepts the current in-progress page
// and creates a new page.
func (d *DocumentBuilder) Accept() {
	if d.lastPage != -1 {
		d.pages[d.lastPage].Accept()
	}

	d.pages = append(d.pages, NewPage())
	d.lastPage++
}

// AddString adds a string to the document.
//
// Parameters:
//   - str: The string to add.
//
// Behaviors:
//   - If the string is empty, it is not added.
func (d *DocumentBuilder) AddString(str string) {
	if str == "" {
		return
	}

	if d.lastPage == -1 {
		d.pages = append(d.pages, NewPage())
		d.lastPage = 0
	}

	d.pages[d.lastPage].AddString(str)
}

// AddRune adds a rune to the document.
//
// Parameters:
//   - r: The rune to add.
func (d *DocumentBuilder) AddRune(r rune) {
	if d.lastPage == -1 {
		d.pages = append(d.pages, NewPage())
		d.lastPage = 0
	}

	d.pages[d.lastPage].AddRune(r)
}

// AcceptWord is a function that accepts the current in-progress word
// and creates a new word.
func (d *DocumentBuilder) AcceptWord() {
	if d.lastPage != -1 {
		d.pages[d.lastPage].AcceptWord()
	}
}

// AcceptSection is a function that accepts the current in-progress section
// and creates a new section.
func (d *DocumentBuilder) AcceptSection() {
	if d.lastPage != -1 {
		d.pages[d.lastPage].Accept()
	}
}

// Finalize is a function that finalizes the document.
func (d *DocumentBuilder) Finalize() {
	if d.lastPage == -1 {
		return
	}

	d.pages[d.lastPage].Finalize()
}

// FString is a function that prints the document.
//
// Parameters:
//   - trav: The traversor to use for printing.
//
// Returns:
//   - error: An error if the printing fails.
func (d *DocumentBuilder) FString(trav *ffs.Traversor) error {
	for _, page := range d.pages {
		err := page.FString(trav)
		if err != nil {
			return err
		}
	}

	return nil
}

// NewDocument creates a new document.
//
// Parameters:
//   - pages: The pages to add to the document.
//
// Returns:
//   - *Document: A pointer to the newly created document.
func NewDocument() *DocumentBuilder {
	return &DocumentBuilder{
		pages:    make([]*PageBuilder, 0),
		lastPage: -1,
	}
}

/*

// String returns the string representation of the document.
//
// Returns:
//   - string: The string representation of the document.
func (d *Document) String() string {
	return strings.Join(d.lines, "\n")
}

func (d *Document) Tmp() []string {
	return d.lines
}

// FString returns the formatted string representation of the document.
//
// Parameters:
//   - indentLevel: The level of indentation.
//
// Returns:
//   - []string: The formatted string representation of the document.
func (d *Document) FString(trav *ffs.Traversor) error {
	err := trav.AddLines(d.lines)
	if err != nil {
		return err
	}

	return nil
}

// NewDocument creates a new document.
//
// Parameters:
//   - sentences: The sentences to add to the document.
//
// Returns:
//   - *Document: A pointer to the newly created document.
//
// Behaviors:
//   - The sentences are separated by a space and on the same line.
func NewDocument(sentences ...string) *Document {
	d := &Document{
		lines: make([]string, 0),
	}

	d.AddLine(sentences...)

	return d
}

// AddLine adds sentences to the document separated by a space.
// The line is split by the newline character.
//
// Parameters:
//   - line: The line to add.
//
// Returns:
//   - *Document: A pointer to the document. This allows for chaining.
//
// Example:
//   - AddLine("Hello,", "world!")
//   - AddLine("This is a sentence.")
func (d *Document) AddLine(sentences ...string) *Document {
	if len(sentences) == 0 {
		return d
	}

	var builder strings.Builder

	if strings.HasSuffix(sentences[0], "\n") {
		builder.WriteString(strings.TrimSuffix(sentences[0], "\n"))

		d.lines = append(d.lines, strings.Split(builder.String(), "\n")...)
	} else {
		builder.WriteString(sentences[0])
	}

	for _, sentence := range sentences[1:] {
		if strings.HasSuffix(sentence, "\n") {
			builder.WriteRune(' ')
			builder.WriteString(strings.TrimSuffix(sentence, "\n"))

			d.lines = append(d.lines, strings.Split(builder.String(), "\n")...)
		} else {
			builder.WriteRune(' ')
			builder.WriteString(sentence)
		}
	}

	if builder.Len() != 0 {
		d.lines = append(d.lines, strings.Split(builder.String(), "\n")...)
	}

	return d
}

*/
