package Strings

/*
// GetMaxWidth is a method that returns the maximum width of the TextSplit.
//
// Returns:
//   - int: The maximum width of the TextSplit.
func (ts *TextSplit) GetMaxHeight() int {
	return ts.maxHeight
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
