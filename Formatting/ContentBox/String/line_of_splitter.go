package String

import intf "github.com/PlayerR9/MyGoLib/Units/Common"

// lineOfSplitter is a helper struct used in the SplitTextInEqualSizedLines function.
// It represents a line of text.
type lineOfSplitter struct {
	// The line field is a slice of strings, each representing a word in the line.
	line []*String

	// The len field is an integer representing the total length of the line,
	// including spaces between words.
	len int
}

// Copy is a method of intf.Copier that creates a shallow copy of the SpltLine.
//
// Returns:
//   - intf.Copier: A shallow copy of the SpltLine.
func (sl *lineOfSplitter) Copy() intf.Copier {
	newLine := make([]*String, len(sl.line))
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
func newLineOfSplitter(word *String) *lineOfSplitter {
	splt := &lineOfSplitter{
		line: []*String{word},
		len:  word.length,
	}

	return splt
}

// shiftLeft is an helper method of SpltLine that removes the first word of the line.
//
// Returns:
//   - string: The word that was removed.
func (sl *lineOfSplitter) shiftLeft() *String {
	firstWord := sl.line[0]

	sl.line = sl.line[1:]
	sl.len -= (firstWord.length + 1)

	return firstWord
}

// InsertWord is a method of SpltLine that adds a given word to the end of the line.
//
// If the word is an empty string, it is ignored.
//
// Parameters:
//   - word: The word to add to the line.
func (sl *lineOfSplitter) insertWord(word *String) {
	if word == nil {
		return
	}

	sl.line = append(sl.line, word)
	sl.len += (word.length + 1)
}
