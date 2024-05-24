package Builder

import (
	"strings"

	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
)

// SectionBuilder is a type that represents a section of a page.
type SectionBuilder struct {
	// builder is the string builder for the section.
	builder strings.Builder

	// words are the words in the section.
	words []string
}

// Accept is a function that accepts the current word and
// creates a new word.
func (s *SectionBuilder) Accept() {
	if s.builder.Len() == 0 {
		return
	}

	s.words = append(s.words, s.builder.String())
	s.builder.Reset()
}

// AddString adds a string to the section.
//
// Parameters:
//   - str: The string to add.
//
// Behaviors:
//   - If the string is empty, it is not added.
func (s *SectionBuilder) AddString(str string) {
	if str == "" {
		return
	}

	s.builder.WriteString(str)
}

// AddRune adds a rune to the section.
//
// Parameters:
//   - r: The rune to add.
func (s *SectionBuilder) AddRune(r rune) {
	s.builder.WriteRune(r)
}

// AcceptWord is a function that accepts the current in-progress word
// and creates a new word.
func (s *SectionBuilder) AcceptWord() {
	s.Accept()
}

// Finalize is a function that finalizes the section.
func (s *SectionBuilder) Finalize() {
	if s.builder.Len() == 0 {
		return
	}

	s.words = append(s.words, s.builder.String())
	s.builder.Reset()
}

// FString is a function that prints the section.
//
// Parameters:
//   - trav: The traversor to use for printing.
//
// Returns:
//   - error: An error if the printing fails.
func (s *SectionBuilder) FString(trav *ffs.Traversor) error {
	if trav == nil {
		return nil
	}

	for _, word := range s.words {
		if err := trav.AppendString(word); err != nil {
			return err
		}
	}

	return nil
}

// NewSection creates a new section.
//
// Returns:
//   - *Section: The new section.
func NewSection() *SectionBuilder {
	return &SectionBuilder{
		words: make([]string, 0),
	}
}
