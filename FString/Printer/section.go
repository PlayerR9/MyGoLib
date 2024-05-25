package Printer

import (
	"strings"
	"sync"
)

// sectionBuilder is a type that represents a section of a page.
type sectionBuilder struct {
	// builder is the string builder for the section.
	builder strings.Builder

	// words are the words in the section.
	words [][]string

	// lastLine is the last line of the section.
	lastLine int

	// mu is the mutex for the builder.
	mu sync.RWMutex
}

// newSectionBuilder creates a new section.
//
// Returns:
//   - *Section: The new section.
func newSectionBuilder() *sectionBuilder {
	return &sectionBuilder{
		words:    [][]string{{}},
		lastLine: 0,
	}
}

// GetWords is a function that returns the words of the section.
//
// Returns:
//   - [][]string: The words of the section.
func (s *sectionBuilder) GetWords() [][]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.words
}

// IsFirstOfLine is a function that returns true if the current position is the first
// position of a line.
//
// Returns:
//   - bool: True if the current position is the first position of a line.
func (s *sectionBuilder) IsFirstOfLine() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.words[s.lastLine]) != 0 {
		return false
	}

	return s.builder.Len() == 0
}

// Accept is a function that accepts the current word and
// creates a new line.
func (s *sectionBuilder) Accept() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.builder.Len() > 0 {
		s.words[s.lastLine] = append(s.words[s.lastLine], s.builder.String())
		s.builder.Reset()
	}

	s.words = append(s.words, []string{})
	s.lastLine++
}

// AddString adds a string to the section.
//
// Parameters:
//   - str: The string to add.
//
// Behaviors:
//   - If the string is empty, it is not added.
func (s *sectionBuilder) AddString(str string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if str == "" {
		return
	}

	s.builder.WriteString(str)
}

// AddRune adds a rune to the section.
//
// Parameters:
//   - r: The rune to add.
func (s *sectionBuilder) AddRune(r rune) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.builder.WriteRune(r)
}

// AcceptWord is a function that accepts the current in-progress word
// and resets the builder.
func (s *sectionBuilder) AcceptWord() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.builder.Len() == 0 {
		return
	}

	s.words[s.lastLine] = append(s.words[s.lastLine], s.builder.String())
	s.builder.Reset()
}
