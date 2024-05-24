package Document

// Document is a type that represents a document.
type Document struct {
	// pages are the pages of the document.
	pages []*Page
}

// NewDocument creates a new document.
//
// Returns:
//   - *Document: The new document.
func NewDocument() *Document {
	return &Document{
		pages: make([]*Page, 0),
	}
}
