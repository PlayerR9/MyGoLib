package Strings

import (
	"unicode/utf8"

	intf "github.com/PlayerR9/MyGoLib/Units/Common"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// TextSplit represents a split text with a maximum width and height.
type TextSplit struct {
	// lines is the lines of the split text.
	lines []*lineOfSplitter

	// maxWidth is the maximum length of a line.
	maxWidth int

	// maxHeight is the maximum number of lines.
	maxHeight int
}

// Copy is a method of intf.Copier that creates a shallow copy of the TextSplit.
//
// Returns:
//   - intf.Copier: A shallow copy of the TextSplit.
func (ts *TextSplit) Copy() intf.Copier {
	newTs := &TextSplit{
		maxWidth:  ts.maxWidth,
		lines:     make([]*lineOfSplitter, 0, ts.maxHeight),
		maxHeight: ts.maxHeight,
	}

	for _, line := range ts.lines {
		newTs.lines = append(newTs.lines, line.Copy().(*lineOfSplitter))
	}

	return newTs
}

// NewTextSplit creates a new TextSplit with the given maximum width and height.
//
// Parameters:
//   - maxWidth: The maximum length of a line.
//   - maxHeight: The maximum number of lines.
//
// Returns:
//   - *TextSplit: A pointer to the newly created TextSplit.
//   - error: An error of type *ers.ErrInvalidParameter if the maxWidth or
//     maxHeight is less than 0.
func NewTextSplit(maxWidth, maxHeight int) (*TextSplit, error) {
	if maxWidth < 0 {
		return nil, ers.NewErrInvalidParameter(
			"maxWidth",
			ers.NewErrGTE(0),
		)
	}

	if maxHeight < 0 {
		return nil, ers.NewErrInvalidParameter(
			"maxHeight",
			ers.NewErrGTE(0),
		)
	}

	return &TextSplit{
		maxWidth:  maxWidth,
		lines:     make([]*lineOfSplitter, 0, maxHeight),
		maxHeight: maxHeight,
	}, nil
}

// GetMaxWidth is a method that returns the maximum width of the TextSplit.
//
// Returns:
//   - int: The maximum width of the TextSplit.
func (ts *TextSplit) GetMaxHeight() int {
	return ts.maxHeight
}

// GetHeight is a method that returns the height of the TextSplit.
//
// Returns:
//   - int: The height of the TextSplit.
func (ts *TextSplit) GetHeight() int {
	return len(ts.lines)
}

// GetLines is a method that returns the lines of the TextSplit.
//
// Returns:
//   - []*SpltLine: The lines of the TextSplit.
func (ts *TextSplit) GetLines() []*lineOfSplitter {
	return ts.lines
}

// GetFirstLine is a method that returns the first line of the TextSplit.
//
// Returns:
//   - *SpltLine: The first line of the TextSplit.
//
// Behaviors:
//   - If the TextSplit is empty, the method returns nil.
func (ts *TextSplit) GetFirstLine() *lineOfSplitter {
	if len(ts.lines) == 0 {
		return nil
	}

	return ts.lines[0]
}

// GetFurthestRightEdge is a method that returns the number of characters in
// the longest line of the TextSplit.
//
// Returns:
//   - int: The number of characters in the longest line.
//   - bool: True if the TextSplit is not empty, and false otherwise.
func (ts *TextSplit) GetFurthestRightEdge() (int, bool) {
	if len(ts.lines) == 0 {
		return 0, false
	}

	max := ts.lines[0].Length()

	for _, line := range ts.lines[1:] {
		if line.Length() > max {
			max = line.len
		}
	}

	return max, true
}

// canInsertWordAt is a helper method that checks if a given word can be inserted
// into a specific line without exceeding the width of the TextSplit.
//
// Parameters:
//   - word: The word to check.
//   - lineIndex: The index of the line to check.
//
// Returns:
//   - bool: True if the word can be inserted into the line at lineIndex without
//     exceeding the width, and false otherwise.
func (ts *TextSplit) canInsertWordAt(word string, lineIndex int) bool {
	if lineIndex < 0 || lineIndex >= len(ts.lines) {
		return false
	}

	return ts.lines[lineIndex].len+utf8.RuneCountInString(word)+1 <= ts.maxWidth
}

/*
// insertWordAt is a helper method that inserts a given word into a specific line of the TextSplit.
//
// Parameters:
//   - word: The word to insert.
//   - lineIndex: The index of the line to insert the word into.
//
// Returns:
//   - error: An error of type *ers.ErrInvalidParameter if the lineIndex is out of bounds.
func (ts *TextSplit) insertWordAt(word string, lineIndex int) error {
	if lineIndex < 0 || lineIndex >= len(ts.lines) {
		return ers.NewErrInvalidParameter(
			"lineIndex",
			ers.NewErrOutOfBounds(lineIndex, 0, len(ts.lines)),
		)
	}

	// Check if adding the next word to the last line exceeds the width.
	// If it does, we shift the words of the last line to the left.
	for !ts.canInsertWordAt(word, lineIndex) && lineIndex >= 0 {
		firstWord := ts.lines[lineIndex].shiftLeft()
		ts.lines[lineIndex].InsertWord(word)
		word = firstWord

		lineIndex--
	}

	/*
		if !ts.CanInsertWord(word, lineIndex) {
			panic(fmt.Errorf("word %s cannot be inserted into line %d", word, lineIndex))
		}


	ts.lines[lineIndex].InsertWord(word)

	return nil
}
*/

// InsertWord is a method that attempts to insert a given word into
// the TextSplit.
//
// Parameters:
//   - word: The word to insert.
//
// Returns:
//   - bool: True if the word was successfully inserted, and false if the word is
//     too long to fit within the width of the TextSplit.
func (ts *TextSplit) InsertWord(word string) bool {
	if len(ts.lines) < ts.maxHeight {
		if utf8.RuneCountInString(word) > ts.maxWidth {
			return false
		}

		ts.lines = append(ts.lines, newLineOfSplitter(word))

		return true
	}

	lastLineIndex := ts.maxHeight - 1

	// Check if adding the next word to the last line exceeds the width.
	// If it does, we shift the words of the last line to the left.
	for !ts.canInsertWordAt(word, lastLineIndex) && lastLineIndex >= 0 {
		firstWord := ts.lines[lastLineIndex].shiftLeft()
		ts.lines[lastLineIndex].InsertWord(word)
		word = firstWord

		lastLineIndex--
	}

	ok := ts.canInsertWordAt(word, lastLineIndex)
	if !ok {
		return false
	}

	ts.lines[lastLineIndex].InsertWord(word)

	return true
}

// canShiftUp is an helper method that checks if the first word of a given line
// can be shifted up to the previous line without exceeding the width.
//
// Parameters:
//   - lineIndex: The index of the line to check.
//
// Returns:
//   - bool: True if the first word of the line at lineIndex can be shifted up to the
//     previous line without exceeding the width, and false otherwise.
func (ts *TextSplit) canShiftUp(lineIndex int) bool {
	return ts.canInsertWordAt(ts.lines[lineIndex].line[0], lineIndex-1)
}

// shiftUp is an helper method that shifts the first word of a given line up to
// the previous line.
//
// Parameters:
//   - lineIndex: The index of the line to shift up.
func (ts *TextSplit) shiftUp(lineIndex int) {
	ts.lines[lineIndex-1].InsertWord(ts.lines[lineIndex].shiftLeft())
}
