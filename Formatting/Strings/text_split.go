package Strings

/*
// GetMaxWidth is a method that returns the maximum width of the TextSplit.
//
// Returns:
//   - int: The maximum width of the TextSplit.
func (ts *TextSplit) GetMaxHeight() int {
	return ts.maxHeight
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

// InsertWords is a method that attempts to insert multiple words into the TextSplit.
//
// Parameters:
//   - words: The words to insert.
//
// Returns:
//   - int: The index of the first word that could not be inserted, or -1 if all words were inserted.
func (ts *TextSplit) InsertWords(words ...string) int {
	for i, word := range words {
		if !ts.InsertWord(word) {
			return i
		}
	}

	return -1
}
*/
