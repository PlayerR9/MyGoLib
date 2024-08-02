package StringExt

import (
	"strings"
	"unicode/utf8"

	uc "github.com/PlayerR9/lib_units/common"
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

// Copy is a method that creates a copy of the TextSplit.
//
// Returns:
//   - *TextSplit: A copy of the TextSplit.
func (ts *TextSplit) Copy() *TextSplit {
	lines_copy := make([]*lineOfSplitter, 0, len(ts.lines))

	for _, line := range ts.lines {
		lines_copy = append(lines_copy, line.Copy())
	}

	tsCopy := &TextSplit{
		maxWidth:  ts.maxWidth,
		lines:     lines_copy,
		maxHeight: ts.maxHeight,
	}

	return tsCopy
}

// GetRunes is a method of TextSplit that returns the runes of the TextSplit.
//
// Returns:
//   - [][]rune: A slice of runes representing the words in the TextSplit.
//
// Behaviors:
//   - It is always a slice of runes with one line.
func (ts *TextSplit) GetRunes() [][]rune {
	if len(ts.lines) == 0 {
		return [][]rune{{}}
	}

	runeTable := make([][]rune, 0, len(ts.lines))

	for _, line := range ts.lines {
		row := line.GetRunes()
		runeTable = append(runeTable, row[0])
	}

	return runeTable
}

// NewTextSplit creates a new TextSplit with the given maximum width and height.
//
// Parameters:
//   - maxWidth: The maximum length of a line.
//   - maxHeight: The maximum number of lines.
//
// Returns:
//   - *TextSplit: A pointer to the newly created TextSplit.
//   - error: An error of type *uc.ErrInvalidParameter if the maxWidth or
//     maxHeight is less than 0.
func NewTextSplit(maxWidth, maxHeight int) (*TextSplit, error) {
	if maxWidth < 0 {
		return nil, uc.NewErrInvalidParameter(
			"maxWidth",
			uc.NewErrGTE(0),
		)
	}

	if maxHeight < 0 {
		return nil, uc.NewErrInvalidParameter(
			"maxHeight",
			uc.NewErrGTE(0),
		)
	}

	ts := &TextSplit{
		maxWidth:  maxWidth,
		lines:     make([]*lineOfSplitter, 0, maxHeight),
		maxHeight: maxHeight,
	}

	return ts, nil
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

	wordLen := utf8.RuneCountInString(word)
	totalLen := ts.lines[lineIndex].len + wordLen

	return totalLen+1 <= ts.maxWidth
}

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
		wordLen := utf8.RuneCountInString(word)

		if wordLen > ts.maxWidth {
			return false
		}

		los := newLineOfSplitter(word)

		ts.lines = append(ts.lines, los)

		return true
	}

	lastLineIndex := ts.maxHeight - 1

	// Check if adding the next word to the last line exceeds the width.
	// If it does, we shift the words of the last line to the left.
	ok := ts.canInsertWordAt(word, lastLineIndex)

	for !ok && lastLineIndex >= 0 {
		lastLine := ts.lines[lastLineIndex]

		firstWord := lastLine.shiftLeft()
		lastLine.insertWord(word)

		word = firstWord

		lastLineIndex--
	}

	ok = ts.canInsertWordAt(word, lastLineIndex)
	if !ok {
		return false
	}

	lastLine := ts.lines[lastLineIndex]
	lastLine.insertWord(word)

	return true
}

// InsertWords is a method that attempts to insert multiple words into the TextSplit.
//
// Parameters:
//   - words: The words to insert.
//
// Returns:
//   - int: The index of the first word that could not be inserted, or -1 if all words were inserted.
func (ts *TextSplit) InsertWords(words []string) int {
	for i, word := range words {
		ok := ts.InsertWord(word)
		if !ok {
			return i
		}
	}

	return -1
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
	ok := ts.canInsertWordAt(ts.lines[lineIndex].line[0], lineIndex-1)

	return ok
}

// shiftUp is an helper method that shifts the first word of a given line up to
// the previous line.
//
// Parameters:
//   - lineIndex: The index of the line to shift up.
func (ts *TextSplit) shiftUp(lineIndex int) {
	lastLine := ts.lines[lineIndex]
	firstWord := lastLine.shiftLeft()

	secondLastLine := ts.lines[lineIndex-1]
	secondLastLine.insertWord(firstWord)
}

// GetHeight is a method that returns the height of the TextSplit.
//
// Returns:
//   - int: The height of the TextSplit.
func (ts *TextSplit) GetHeight() int {
	height := len(ts.lines)
	return height
}

// GetLines is a method that returns the lines of the TextSplit.
//
// Returns:
//   - []*SpltLine: The lines of the TextSplit.
func (ts *TextSplit) GetLines() []string {
	if len(ts.lines) == 0 {
		return nil
	}

	lines := make([]string, 0, len(ts.lines))

	for _, line := range ts.lines {
		str := strings.Join(line.line, " ")
		lines = append(lines, str)
	}

	return lines
}

// GetFirstLine is a method that returns the first line of the TextSplit.
//
// Returns:
//   - []string: The first line of the TextSplit, or nil if the TextSplit is empty.
//
// Behaviors:
//   - If the TextSplit is empty, the method returns nil.
func (ts *TextSplit) GetFirstLine() []string {
	if len(ts.lines) == 0 {
		return nil
	}

	return ts.lines[0].line
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

	max := ts.lines[0].len

	for _, line := range ts.lines[1:] {
		if line.len > max {
			max = line.len
		}
	}

	return max, true
}
