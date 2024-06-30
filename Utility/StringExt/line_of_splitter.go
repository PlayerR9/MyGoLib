package StringExt

import (
	"strings"
	"unicode/utf8"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// lineOfSplitter is a helper struct used in the SplitTextInEqualSizedLines function.
// It represents a line of text.
type lineOfSplitter struct {
	// The line field is a slice of strings, each representing a word in the line.
	line []string

	// The len field is an integer representing the total length of the line,
	// including spaces between words.
	len int
}

// Copy implements the common.Copier interface.
func (sl *lineOfSplitter) Copy() uc.Copier {
	newLine := make([]string, len(sl.line))
	copy(newLine, sl.line)

	losCopy := &lineOfSplitter{
		line: newLine,
		len:  sl.len,
	}

	return losCopy
}

// GetRunes is a method of SpltLine that returns the runes of the line.
//
// Always returns a slice of runes with one line.
//
// Returns:
//   - [][]rune: A slice of runes representing the words in the line.
//
// Behaviors:
//   - It is always a slice of runes with one line.
func (sl *lineOfSplitter) GetRunes() [][]rune {
	if len(sl.line) == 0 {
		return [][]rune{{}}
	}

	str := strings.Join(sl.line, " ")

	return [][]rune{[]rune(str)}
}

// newLineOfSplitter is a helper function that creates a new line of
// splitter with the given word.
//
// Parameters:
//   - word: The initial word to add to the line.
//
// Returns:
//   - *lineOfSplitter: A pointer to the newly created line of splitter.
func newLineOfSplitter(word string) *lineOfSplitter {
	len := utf8.RuneCountInString(word)

	splt := &lineOfSplitter{
		line: []string{word},
		len:  len,
	}

	return splt
}

// shiftLeft is an helper method of SpltLine that removes the first word of the line.
//
// Returns:
//   - string: The word that was removed.
func (sl *lineOfSplitter) shiftLeft() string {
	firstWord := sl.line[0]

	sl.line = sl.line[1:]
	sl.len -= utf8.RuneCountInString(firstWord)
	sl.len-- // Remove the extra space

	return firstWord
}

// InsertWord is a method of SpltLine that adds a given word to the end of the line.
//
// Parameters:
//   - word: The word to add to the line.
//
// Behaviors:
//   - If the word is an empty string, it is ignored.
func (sl *lineOfSplitter) insertWord(word string) {
	if word == "" {
		return
	}

	sl.line = append(sl.line, word)
	sl.len += utf8.RuneCountInString(word)
	sl.len++ // Add the extra space
}
