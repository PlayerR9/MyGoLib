package Builder

import (
	"unicode/utf8"

	tld "github.com/PlayerR9/MyGoLib/TextLike/Document"
	ue "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// FieldSplitter is a type that splits a string into fields.
type FieldSplitter struct {
	// replaceTab is the string to replace a tab character with.
	replaceTab string

	// content is the content of the FieldSplitter.
	content *DocumentBuilder

	// sections are the sections of the page.
	sections []*SectionBuilder
}

// NewFieldSplitter creates a new FieldSplitter instance.
//
// Parameters:
//   - replaceTab: The string to replace a tab character with.
//
// Returns:
//   - *FieldSplitter: The created FieldSplitter instance.
func NewFieldSplitter(replaceTab string) *FieldSplitter {
	return &FieldSplitter{
		replaceTab: replaceTab,
		content:    NewDocument(),
	}
}

// SplitSentenceIntoFields splits the string into fields, where each field is a
// substring separated by one or more whitespace charactue.
//
// Parameters:
//   - sentence: The string to split into fields.
//   - indentLevel: The number of spaces that a tab character is replaced with.
//
// Returns:
//   - [][]string: A two-dimensional slice of strings, where each inner slice
//     represents the fields of a line from the input string.
//   - error: An error of type *ue.ErrInvalidRuneAt if an invalid rune is found in
//     the sentence.
//
// Behaviors:
//   - Negative indentLevel values are converted to positive values.
//   - Empty sentences return a nil slice with no errors.
//   - The function handles the following whitespace characters: space, tab,
//     vertical tab, carriage return, line feed, and form feed.
//   - The function returns a partial result if an invalid rune is found where
//     the result are the fields found up to that point.
func (fs *FieldSplitter) Apply(sentence string) error {
	if sentence == "" {
		return nil
	}

	for j := 0; len(sentence) > 0; j++ {
		char, size := utf8.DecodeRuneInString(sentence)
		sentence = sentence[size:]

		if char == utf8.RuneError {
			fs.content.Finalize()

			return ue.NewErrAt(j, "character", ue.NewErrInvalidRune(nil))
		}

		switch char {
		case '\t':
			// Replace tabs with N spaces
			fs.content.AddString(fs.replaceTab)
		case '\v':
			// Do nothing
		case '\r':
			if size != 0 {
				nextRune, size := utf8.DecodeRuneInString(sentence)

				if nextRune == '\n' {
					sentence = sentence[size:]
				}
			}

			fallthrough
		case '\n', '\u0085':
			fs.content.AcceptSection()
		case '\f':
			fs.content.Accept()
		case ' ':
			fs.content.AcceptWord()
		case '\u00A0':
			fs.content.AddRune(' ')
		default:
			fs.content.AddRune(char)
		}
	}

	fs.content.Finalize()

	return nil
}

// GetDocument returns the content of the FieldSplitter as a Document.
//
// Returns:
//   - *Document: The content of the FieldSplitter.
func (fs *FieldSplitter) GetDocument() *DocumentBuilder {
	return fs.content
}

// Build is a function that builds the document.
//
// Returns:
//   - *tld.Document: The built document.
func (fs *FieldSplitter) Build() *tld.Document {
	doc := tld.NewDocument()

	for _, page := range fs.content.pages {
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
