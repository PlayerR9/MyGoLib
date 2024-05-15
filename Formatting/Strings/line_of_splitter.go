package Strings

import (
	"strings"
	"unicode/utf8"

	intf "github.com/PlayerR9/MyGoLib/Units/Common"
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

// String is a method of fmt.Stringer that returns the string representation of the SpltLine.
//
// Returns:
//   - string: The resulting string.
func (sl *lineOfSplitter) String() string {
	return strings.Join(sl.line, " ")
}

// Copy is a method of intf.Copier that creates a shallow copy of the SpltLine.
//
// Returns:
//   - intf.Copier: A shallow copy of the SpltLine.
func (sl *lineOfSplitter) Copy() intf.Copier {
	newLine := make([]string, len(sl.line))
	copy(newLine, sl.line)

	return &lineOfSplitter{
		line: newLine,
		len:  sl.len,
	}
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
	splt := &lineOfSplitter{
		line: []string{word},
		len:  utf8.RuneCountInString(word),
	}

	return splt
}

// GetWords is a method of SpltLine that returns the words in the line.
//
// Returns:
//   - []string: The words in the line.
func (sl *lineOfSplitter) GetWords() []string {
	return sl.line
}

// Length is a method that returns the length of the line. (i.e., the
// number of characters in the line, including spaces between words.)
//
// Returns:
//   - int: The length of the line.
func (sl *lineOfSplitter) Length() int {
	return sl.len
}

// shiftLeft is an helper method of SpltLine that removes the first word of the line.
//
// Returns:
//   - string: The word that was removed.
func (sl *lineOfSplitter) shiftLeft() string {
	firstWord := sl.line[0]

	sl.line = sl.line[1:]
	sl.len -= (utf8.RuneCountInString(firstWord) + 1)

	return firstWord
}

// InsertWord is a method of SpltLine that adds a given word to the end of the line.
//
// If the word is an empty string, it is ignored.
//
// Parameters:
//   - word: The word to add to the line.
func (sl *lineOfSplitter) InsertWord(word string) {
	if word == "" {
		return
	}

	sl.line = append(sl.line, word)
	sl.len += (utf8.RuneCountInString(word) + 1)
}

// GetFirstWord is a method of SpltLine that returns the first word of the line.
//
// Returns:
//   - string: The first word of the line.
//   - bool: True if there is a first word, and false otherwise.
func (sl *lineOfSplitter) GetFirstWord() (string, bool) {
	if len(sl.line) == 0 {
		return "", false
	}

	return sl.line[0], true
}
