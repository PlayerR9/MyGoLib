package Document

import (
	tls "github.com/PlayerR9/MyGoLib/TextLike/Sections"
)

// Page is a type that represents a page of a document.
type Page struct {
	// sections are the sections of the page.
	sections []tls.Sectioner
}
