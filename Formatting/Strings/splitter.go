package Strings

import (
	"math"
	"strings"
	"unicode/utf8"

	intf "github.com/PlayerR9/MyGoLib/Units/Common"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
	mext "github.com/PlayerR9/MyGoLib/Utility/MathExt"
	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

// SpltLine is a helper struct used in the SplitTextInEqualSizedLines function.
// It represents a line of text.
type SpltLine struct {
	// The Line field is a slice of strings, each representing a word in the line.
	Line []string

	// The Len field is an integer representing the total length of the line,
	// including spaces between words.
	Len int
}

// String is a method of fmt.Stringer that returns the string representation of the SpltLine.
//
// Returns:
//   - string: The resulting string.
func (sl *SpltLine) String() string {
	return strings.Join(sl.Line, " ")
}

// NewSpltLine creates a new SpltLine with the given initial word.
//
// Parameters:
//   - word: The initial word to add to the line.
//
// Returns:
//   - *SpltLine: A pointer to the newly created SpltLine.
func NewSpltLine(word string) *SpltLine {
	splt := &SpltLine{
		Line: []string{word},
		Len:  utf8.RuneCountInString(word),
	}

	return splt
}

// shiftLeft is an helper method of SpltLine that removes the first word of the line.
//
// Returns:
//   - string: The word that was removed.
func (sl *SpltLine) shiftLeft() string {
	firstWord := sl.Line[0]

	sl.Line = sl.Line[1:]
	sl.Len -= (utf8.RuneCountInString(firstWord) + 1)

	return firstWord
}

// InsertWord is a method of SpltLine that adds a given word to the end of the line.
//
// If the word is an empty string, it is ignored.
//
// Parameters:
//   - word: The word to add to the line.
func (sl *SpltLine) InsertWord(word string) {
	if word == "" {
		return
	}

	sl.Line = append(sl.Line, word)
	sl.Len += (utf8.RuneCountInString(word) + 1)
}

// Copy is a method of intf.Copier that creates a shallow copy of the SpltLine.
//
// Returns:
//   - intf.Copier: A shallow copy of the SpltLine.
func (sl *SpltLine) Copy() intf.Copier {
	newLine := make([]string, len(sl.Line))
	copy(newLine, sl.Line)

	return &SpltLine{
		Line: newLine,
		Len:  sl.Len,
	}
}

// TextSplit is a helper struct used in the SplitTextInEqualSizedLines function.
type TextSplit struct {
	// The Width represents the maximum length of a line.
	Width int

	// The Lines field is a slice of pointers to SpltLine structs, each representing
	// a line of text.
	Lines []*SpltLine
}

// NewTextSplit creates a new TextSplit with the given width.
//
// Parameters:
//   - width: The maximum length of a line.
//   - height: The maximum number of lines.
//
// Returns:
//   - *TextSplit: A pointer to the newly created TextSplit.
//   - error: An error of type *ers.ErrInvalidParameter if the width or height is less
//     than 0.
func NewTextSplit(width, height int) (*TextSplit, error) {
	if width < 0 {
		return nil, ers.NewErrInvalidParameter(
			"width",
			ers.NewErrGTE(0),
		)
	}

	if height < 0 {
		return nil, ers.NewErrInvalidParameter(
			"width",
			ers.NewErrGTE(0),
		)
	}

	return &TextSplit{
		Width: width,
		Lines: make([]*SpltLine, 0, height),
	}, nil
}

func (ts *TextSplit) GetFirstLine() *SpltLine {
	if len(ts.Lines) == 0 {
		return nil
	}

	return ts.Lines[0]
}

// GetFurthestRightEdge is a method that returns the length of the
// longest line in the TextSplit.
//
// Returns:
//   - int: The length of the longest line.
func (ts *TextSplit) GetFurthestRightEdge() int {
	if len(ts.Lines) == 0 {
		return ts.Width
	}

	max := ts.Lines[0].Len

	for _, line := range ts.Lines[1:] {
		if line.Len > max {
			max = line.Len
		}
	}

	return max
}

// CanInsertWord is a method that checks if a given word can be
// inserted into a specific line without exceeding the width of the TextSplit.
//
// Parameters:
//   - word: The word to check.
//   - lineIndex: The index of the line to check.
//
// Returns:
//   - bool: True if the word can be inserted into the line at lineIndex without
//     exceeding the width, and false otherwise.
func (ts *TextSplit) CanInsertWord(word string, lineIndex int) bool {
	if lineIndex < 0 || lineIndex >= len(ts.Lines) {
		return false
	}

	return ts.Lines[lineIndex].Len+utf8.RuneCountInString(word)+1 <= ts.Width
}

// InsertWordAt is a method that attempts to insert a given word
// into a specific line of the TextSplit.
//
// Parameters:
//   - word: The word to insert.
//   - lineIndex: The index of the line to insert the word into.
//
// Returns:
//   - error: An error of type *ers.ErrInvalidParameter if the lineIndex is out of bounds.
func (ts *TextSplit) InsertWordAt(word string, lineIndex int) error {
	if lineIndex < 0 || lineIndex >= len(ts.Lines) {
		return ers.NewErrInvalidParameter(
			"InsertWordAt",
			ers.NewErrOutOfBounds(lineIndex, 0, len(ts.Lines)),
		)
	}

	// Check if adding the next word to the last line exceeds the width.
	// If it does, we shift the words of the last line to the left.
	for !ts.CanInsertWord(word, lineIndex) && lineIndex >= 0 {
		firstWord := ts.Lines[lineIndex].shiftLeft()
		ts.Lines[lineIndex].InsertWord(word)
		word = firstWord

		lineIndex--
	}

	/*
		if !ts.CanInsertWord(word, lineIndex) {
			panic(fmt.Errorf("word %s cannot be inserted into line %d", word, lineIndex))
		}
	*/

	ts.Lines[lineIndex].InsertWord(word)

	return nil
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
	if len(ts.Lines) < cap(ts.Lines) {
		if utf8.RuneCountInString(word) > ts.Width {
			return false
		}

		ts.Lines = append(ts.Lines, NewSpltLine(word))

		return true
	}

	lastLineIndex := cap(ts.Lines) - 1

	// Check if adding the next word to the last line exceeds the width.
	// If it does, we shift the words of the last line to the left.
	for !ts.CanInsertWord(word, lastLineIndex) && lastLineIndex >= 0 {
		firstWord := ts.Lines[lastLineIndex].shiftLeft()
		ts.Lines[lastLineIndex].InsertWord(word)
		word = firstWord

		lastLineIndex--
	}

	/*
		ts.CanInsertWord(word, lastLineIndex)
	*/
	ts.Lines[lastLineIndex].InsertWord(word)

	return true
}

// Copy is a method of intf.Copier that creates a shallow copy of the TextSplit.
//
// Returns:
//   - intf.Copier: A shallow copy of the TextSplit.
func (ts *TextSplit) Copy() intf.Copier {
	newTs := &TextSplit{
		Width: ts.Width,
		Lines: make([]*SpltLine, 0, len(ts.Lines)),
	}

	for _, line := range ts.Lines {
		newTs.Lines = append(newTs.Lines, line.Copy().(*SpltLine))
	}

	return newTs
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
	return ts.CanInsertWord(ts.Lines[lineIndex].Line[0], lineIndex-1)
}

// shiftUp is an helper method that shifts the first word of a given line up to
// the previous line.
//
// Parameters:
//   - lineIndex: The index of the line to shift up.
func (ts *TextSplit) shiftUp(lineIndex int) {
	ts.Lines[lineIndex-1].InsertWord(ts.Lines[lineIndex].shiftLeft())
}

// CalculateNumberOfLines is a function that calculates the minimum number
// of lines needed to fit a given text within a specified width.
//
// Errors:
//   - *ers.ErrInvalidParameter: If the width is less than or equal to 0.
//   - *ErrLinesGreaterThanWords: If the calculated number of lines is greater
//     than the number of words in the text.
//
// Parameters:
//   - text: The slice of strings representing the text to calculate the number of
//     lines for.
//   - width: The width to fit the text within.
//
// Returns:
//
//   - int: The calculated number of lines needed to fit the text within the width.
//   - error: An error if it occurs during the calculation.
//
// The function calculates the total length of the text (Tl) and uses a mathematical
// formula to estimate the minimum number of lines needed to fit the text within the
// given width. The formula is explained in detail in the comments within the function.
func CalculateNumberOfLines(text []string, width int) (int, error) {
	if width <= 0 {
		return 0, ers.NewErrInvalidParameter(
			"width",
			ers.NewErrGT(0),
		)
	} else if len(text) == 0 {
		return 0, nil
	}

	// Euristic to calculate the least amount of splits needed
	// in order to fit the text within the width.
	//
	// Assuming:
	// 	- $n$ is the number of fields in the text using strings.Fields().
	//    - $\omega$ is the number of characters in a field (word).
	//    - $W$ is the total width. This considers only usable width, i.e.,
	//    width - padding.
	//    - $Tl$ = Total length of the text. Calculated by doing the
	// 	sum of the lengths of all the words plus the number of
	// 	spaces between them. $n - 1 + \Sum_{i=1}^n \omega_i$.
	//    - $x$ = Number of splits needed to fit the text within the width.
	//
	// Formula:
	//    $\frac{Tl - x}{x + 1} \leq W$
	//
	// Explanation:
	//
	// 	- $Tl - x$: For every split, the number of characters occupied by the
	//    text is reduced as the space between the splitted fields is removed.
	//    For example: "Hello World" has 11 characters. With one split, it becomes
	//    "Hello" and "World", which has 5 and 5 characters respectively. The
	//    total number of characters is 10, which is equal to 11 - 1.
	//		- $x + 1$: The number of lines is equal to the number of splits plus one as
	//    no splits (x = 0), gives us a single line (x + 1 = 1).
	// 	- $\frac{Tl - x}{x + 1}$: This divides the number of characters occupied by
	//    the title by the number of lines; giving us how many characters are
	//    occupied per line. $\leq W$ ensures that no line exceeds the
	//    width of the screen.
	//
	// Simplification:
	//    $$
	//    \begin{align}
	//    	\frac{Tl - x}{x + 1} &\leq W \\
	//    	Tl - x &\leq W(x + 1) \\
	//    	Tl - x &\leq Wx + W \\
	//    	Tl - W &\leq Wx + x \\
	//       Tl - W &\leq x(W + 1) \\
	//       \frac{Tl - W}{W + 1} &\leq x \\
	//       \lceil\frac{Tl - W}{W + 1}\rceil &\leq x
	//    \end{align}
	//    $$
	//   	Note: The ceil function is used as we cannot do non-integer splits.
	//		and, since we want $x$ to be greater or equal to the result of the
	//		division, we round up the result.
	//
	// Example: If we have the following text: "Hello World, this is a test",
	// with a width of 12 characters, we have:
	//    - $n = 6$
	//    - $W = 12$
	//    - $Tl = 27$
	//
	//	 	$\lceil\frac{27 - 12}{12 + 1}\rceil = \lceil\frac{15}{13}\rceil = 2$
	//
	// Solution:
	//		   *** Hello ***
	// 	*** World, this ***
	// 	 *** is a test ***

	var Tl float64

	for _, word := range text {
		// +1 for the space or any other separator
		Tl += float64(utf8.RuneCountInString(word) + 1)
	}
	Tl-- // Remove the last space or separator

	w := float64(width)

	numberOfLines := int(math.Ceil((Tl-w)/(w+1))) + 1

	if numberOfLines > len(text) {
		return 0, NewErrLinesGreaterThanWords(numberOfLines, len(text))
	} else {
		return numberOfLines, nil
	}
}

type TextSplitter struct {
	width  int
	height int

	candidates []*TextSplit
}

func NewTextSplitter(width int) (*TextSplitter, error) {
	if width <= 0 {
		return nil, ers.NewErrInvalidParameter(
			"width",
			ers.NewErrGT(0),
		)
	}

	tsr := &TextSplitter{
		width:      width,
		height:     -1,
		candidates: make([]*TextSplit, 0),
	}

	return tsr, nil
}

func (ts *TextSplitter) SetHeight(height int) error {
	if height < 1 {
		return ers.NewErrInvalidParameter(
			"height",
			ers.NewErrGT(0),
		)
	}

	ts.height = height

	return nil
}

// SplitInEqualSizedLines is a function that splits a given text into lines of
// equal width.
//
// Errors:
//   - *ers.ErrInvalidParameter: If the input text is empty or the width is less than
//     or equal to 0.
//   - *ErrLinesGreaterThanWords: If the number of lines needed to fit the text
//     within the width is greater than the number of words in the text.
//   - *ErrNoCandidateFound: If no candidate is found during the optimization process.
//
// Parameters:
//   - text: The slice of strings representing the text to split.
//
// Returns:
//   - *TextSplit: A pointer to the created TextSplit instance.
//   - error: An error of type *ErrEmptyText if the input text is empty, or an error
//     of type *ErrWidthTooSmall if the width is less than or equal to 0.
//
// The function calculates the minimum number of lines needed to fit the text within the
// width using the CalculateNumberOfLines function.
// Furthermore, it uses the Sum of Squared Mean (SQM) to find the optimal solution
// for splitting the text into lines of equal width.
//
// If maxHeight is not provided, the function calculates the number of lines needed to fit
// the text within the width using the CalculateNumberOfLines function.
func (tsr *TextSplitter) SplitInEqualSizedLines(text []string) error {
	if len(text) == 0 {
		return ers.NewErrInvalidParameter(
			"text",
			ers.NewErrEmptyString(),
		)
	}

	if tsr.height == -1 {
		var err error

		tsr.height, err = CalculateNumberOfLines(text, tsr.width)
		if err != nil {
			return err
		}
	}

	// We have to find the best way to split the text
	// such that all the lines are as close as possible to
	// the average number of characters per line.
	// Example: "Hello World, this is a test"

	// 1. Add each word to a different line.
	// Example:
	//		*** Hello ***
	//		*** World, ***
	//		 *** this ***
	// The rest is out of bounds. (This is not a problem as we will solve it later)

	// 2. If we still have words left, add them at the last line.
	// If the last line exceeds the width, move the first word of the last line
	// to the above line. Do this until all the words fit within the width.
	// Example:
	//		 *** Hello ***
	//		 *** World, ***
	//		*** this is a ***
	// if we were to add "test" to the last line, it would exceed the width.
	// So, we move "this" to the above line, and add "test" to the last line.
	//		   *** Hello ***
	//		*** World, this ***
	//		 *** is a test ***

	group := &TextSplit{
		Width: tsr.width,
		Lines: make([]*SpltLine, 0, tsr.height),
	}

	for _, word := range text {
		if !group.InsertWord(word) {
			return NewErrLinesGreaterThanWords(tsr.width, len(word))
		}
	}

	// 3. Now we have A solution to the problem, but not THE solution as
	// there may be other ways to split the text that are better than this one.
	// Here is an example where the solution is not optimal:
	//
	// Example: The text "Hi You They" has a Tl of 11 and, assuming
	// W is 8, we have:
	//		$\lceil\frac{11 - 8}{8 + 1}\rceil = \lceil\frac{3}{9}\rceil = 1$
	// This means that the text will be split into two lines:
	//		*** Hi ***
	//		*** You ***
	// With out algorithm, we add "They" to the last line but since it doesn't
	// exceed the width, we don't move any words to the above line.
	//		   *** Hi ***
	//		*** You They ***
	//
	// However, this is not the optimal solution as the following:
	//		*** Hi You ***
	//		 *** They ***
	// is better as the average number of characters per line is closer to the
	// average number of characters per line.

	// We can do this by using SQM (Sum of Squared Mean) as, the lower the SQM,
	// the better the solution.
	// In fact, the optimal solution has an SQM of 1, while our solution has an
	// SQM of 3.

	// 4. Now, we have to find the optimal solution. Because our solution prioritizes
	// the last line, we can do this only by moving words from the last line to the
	// above line; reducing the complexity of the problem.

	// 4.1. For each line that is not the first one, check if the first word of the
	// line can be moved to the above line without exceeding the width.
	// If yes, then it is a candidate for the optimal solution.

	candidates := []*TextSplit{group}

	for i := 0; i < len(candidates); i++ {
		for j := 1; j < tsr.height; j++ {
			// Check if the first word of the line can be moved to the above line.
			// If yes, then it is a candidate for the optimal solution.
			if !candidates[i].canShiftUp(j) {
				continue
			}

			// Copy the candidate as we don't want to modify the original one.
			candidateCopy := candidates[i].Copy().(*TextSplit)
			candidateCopy.shiftUp(j)
			candidates = append(candidates, candidateCopy)
		}
	}

	// 4.2. Calculate the SQM of each candidate and returns the ones with the lowest SQM.
	weights := slext.ApplyWeightFunc(candidates, func(candidate *TextSplit) (float64, bool) {
		values := make([]float64, 0, len(candidate.Lines))

		for _, line := range candidate.Lines {
			values = append(values, float64(line.Len))
		}

		sqm, err := mext.SQM(values)
		if err != nil {
			return 0, false
		}

		return sqm, true
	})
	if len(weights) != 0 {
		tsr.candidates = slext.FilterByNegativeWeight(weights)
	} else {
		tsr.candidates = make([]*TextSplit, 0)
	}

	tsr.height = -1

	return nil
}

// GetSolution returns the solution of the TextSplitter.
//
// Returns:
//   - *TextSplit: The solution of the TextSplitter.
//   - error: An error of type *ErrNoCandidateFound if no candidate is found.
func (tsr *TextSplitter) GetSolution() (*TextSplit, error) {
	if len(tsr.candidates) == 0 {
		return nil, NewErrNoCandidateFound()
	} else if len(tsr.candidates) == 1 {
		return tsr.candidates[0], nil
	}

	// If we have more than one candidate, we have to choose one
	// of them by following other criteria.
	//
	// (For now, we will just choose the first one.)
	// TODO: Choose the best candidate by following other criteria.

	return tsr.candidates[0], nil
}
