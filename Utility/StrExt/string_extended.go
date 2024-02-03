package StrExt

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"strings"
	"unicode/utf8"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
	"github.com/markphelps/optional"
)

// ReplaceSuffix replaces the end of the given string with the provided suffix.
// The function checks the lengths of the string and the suffix to determine the appropriate action:
//   - If the length of the suffix is greater than the length of the string, the function returns an empty string and an ErrSuffixTooLong error.
//   - If the length of the suffix is equal to the length of the string, the function returns the suffix.
//   - If the length of the suffix is zero, the function returns the original string.
//   - Otherwise, the function replaces the end of the string with the suffix and returns the result.
//
// Parameters:
//
//   - str: The original string.
//   - suffix: The suffix to replace the end of the string.
//
// Returns:
//
//   - The modified string, or an error if the suffix is too long.
func ReplaceSuffix(str, suffix string) string {
	countStr := utf8.RuneCountInString(str)
	countSuffix := utf8.RuneCountInString(suffix)

	if countStr < countSuffix {
		panic(ers.NewErrOperationFailed(
			"replace suffix", fmt.Errorf("suffix (%s) is longer than the string (%s)", suffix, str),
		))
	}

	if countStr == countSuffix {
		return suffix
	}

	if countSuffix == 0 {
		return str
	}

	return str[:countStr-countSuffix] + suffix
}

// FindContentIndexes searches for the positions of opening and closing tokens
// in a slice of strings.
// It returns the start and end indexes of the content between the tokens, and
// an error if any.
//
// Parameters:
//
//   - openingToken: The string that marks the beginning of the content.
//   - closingToken: The string that marks the end of the content.
//   - contentTokens: The slice of strings in which to search for the tokens.
//
// Returns:
//   - The start index of the content (inclusive).
//   - The end index of the content (exclusive).
//   - An error if the opening or closing tokens are empty, or if they are not
//     found in the correct order in the contentTokens.
//
// If the openingToken is found but the closingToken is not, the function will
// return an error.
// If the closingToken is found before the openingToken, the function will
// return an error.
// If the closingToken is a newline ("\n") and it is not found, the function will
// return the length of the contentTokens as the end index.
//
// Errors returned can be of type ErrOpeningTokenEmpty, ErrClosingTokenEmpty,
// ErrOpeningTokenNotFound, or a generic error with a message about the closing token.
func FindContentIndexes(openingToken, closingToken string, contentTokens []string) (int, int, error) {
	if openingToken == "" {
		return 0, 0, ers.NewErrInvalidParameter(
			"openingToken", errors.New("opening token cannot be empty"),
		)
	}

	if closingToken == "" {
		return 0, 0, ers.NewErrInvalidParameter(
			"closingToken", errors.New("closing token cannot be empty"),
		)
	}

	openingTokenIndex := slices.Index(contentTokens, openingToken)
	if openingTokenIndex < 0 {
		return 0, 0, ers.NewErrOperationFailed(
			"find content indexes", &ErrOpeningTokenNotFound{openingToken},
		)
	}

	tokenStartIndex := openingTokenIndex + 1

	tokenBalance := 1
	closingTokenIndex := slices.IndexFunc(contentTokens[tokenStartIndex:], func(token string) bool {
		if token == closingToken {
			tokenBalance--
		} else if token == openingToken {
			tokenBalance++
		}

		return tokenBalance == 0
	})

	if closingTokenIndex != -1 {
		return tokenStartIndex, tokenStartIndex + closingTokenIndex + 1, nil
	}

	if tokenBalance < 0 {
		return 0, 0, &ErrNeverOpened{openingToken, closingToken}
	} else if tokenBalance == 1 && closingToken == "\n" {
		return tokenStartIndex, len(contentTokens), nil
	}

	return 0, 0, &ErrClosingTokenNotFound{closingToken}
}

// GetOrdinalSuffix returns the ordinal suffix for a given integer.
//
// Parameters:
//   - number: The integer for which to get the ordinal suffix.
//
// The function returns a string that represents the number with its ordinal suffix.
//
// For example, for the number 1, the function returns "1st"; for the number 2, it returns "2nd"; and so on.
// For numbers ending in 11, 12, or 13, the function returns the number with the suffix "th" (e.g., "11th", "12th", "13th").
// For negative numbers, the function also returns the number with the suffix "th".
//
// If the last digit of the number is 0 or greater than 3 (and the last two digits are not 11, 12, or 13), the function returns the number with the suffix "th".
// If the last digit of the number is 1, 2, or 3 (and the last two digits are not 11, 12, or 13), the function returns the number with the corresponding ordinal suffix ("st", "nd", or "rd").
func GetOrdinalSuffix(number int) string {
	if number < 0 {
		return fmt.Sprintf("%dth", number)
	}

	lastTwoDigits := number % 100
	lastDigit := lastTwoDigits % 10

	if lastTwoDigits >= 11 && lastTwoDigits <= 13 {
		return fmt.Sprintf("%dth", number)
	}

	if lastDigit == 0 || lastDigit > 3 {
		return fmt.Sprintf("%dth", number)
	}

	switch lastDigit {
	case 1:
		return fmt.Sprintf("%dst", number)
	case 2:
		return fmt.Sprintf("%dnd", number)
	case 3:
		return fmt.Sprintf("%drd", number)
	}

	return ""
}

// SpltLine is a helper struct used in the SplitTextInEqualSizedLines function.
// It represents a line of text.
type SpltLine struct {
	// The Line field is a slice of strings, each representing a word in the line.
	Line []string

	// The Len field is an integer representing the total length of the line,
	// including spaces between words.
	Len int
}

// NewSpltLine creates a new SpltLine instance.
// The initial word is passed as an argument.
// The length of the line is calculated based on the initial word.
// It returns a pointer to the created SpltLine instance.
func NewSpltLine(word string) *SpltLine {
	splt := new(SpltLine)

	splt.Line = []string{word}
	splt.Len = utf8.RuneCountInString(word)

	return splt
}

// shiftLeft is a method on the SpltLine struct.
// It removes the first word from the Line slice and decreases the Len field by
// the length of the word plus one (for the space).
// The method returns the word that was removed.
func (sl *SpltLine) shiftLeft() string {
	firstWord := sl.Line[0]

	sl.Line = sl.Line[1:]
	sl.Len -= (utf8.RuneCountInString(firstWord) + 1)

	return firstWord
}

// InsertWord is a method on the SpltLine struct.
// It adds a given word to the end of the Line slice and increases the Len field
// by the length of the word plus one (for the space).
// The method does not return any value.
func (sl *SpltLine) InsertWord(word string) {
	if word == "" {
		return
	}

	sl.Line = append(sl.Line, word)
	sl.Len += (utf8.RuneCountInString(word) + 1)
}

// deepCopy is a method on the SpltLine struct.
// It creates a deep copy of the SpltLine and returns a pointer to it.
// The method creates a new slice for the Line field and copies all words
// from the original Line slice to the new one.
// The Len field is copied directly.
func (sl *SpltLine) deepCopy() *SpltLine {
	newLine := make([]string, len(sl.Line))
	copy(newLine, sl.Line)

	return &SpltLine{
		Line: newLine,
		Len:  sl.Len,
	}
}

// String is a method on the SpltLine struct.
// It converts the SpltLine to a string.
// The method joins all words in the Line slice with a space and returns
// the resulting string.
func (sl SpltLine) String() string {
	return strings.Join(sl.Line, " ")
}

// TextSplitter is a helper struct used in the SplitTextInEqualSizedLines function.
// It holds the width of the lines and a slice of pointers to SpltLine structs.
type TextSplitter struct {
	// The Width represents the maximum length of a line.
	Width int

	// The Lines field is a slice of pointers to SpltLine structs, each representing
	// a line of text.
	Lines []*SpltLine
}

// GetFurthestRightEdge is a method on the TextSplitter struct.
// It iterates over all lines in the TextSplitter and finds the length
// of the longest line.
// If no lines exist, it returns the width of the TextSplitter.
// Otherwise, it returns the length of the longest line.
func (ts *TextSplitter) GetFurthestRightEdge() int {
	max := -1

	for _, line := range ts.Lines {
		if max == -1 || line.Len > max {
			max = line.Len
		}
	}

	if max == -1 {
		return ts.Width
	} else {
		return max
	}
}

// CanInsertWord is a method on the TextSplitter struct.
// It checks if a given word can be inserted into a specific line without
// exceeding the width of the TextSplitter.
// The method takes a string (word) and an integer (lineIndex) as input.
// It returns true if the word can be inserted into the line at lineIndex
// without exceeding the width, and false otherwise.
func (ts *TextSplitter) CanInsertWord(word string, lineIndex int) bool {
	if lineIndex < 0 || lineIndex >= len(ts.Lines) {
		return false
	}

	return ts.Lines[lineIndex].Len+utf8.RuneCountInString(word)+1 <= ts.Width
}

func (ts *TextSplitter) InsertWordAt(word string, lineIndex int) {
	if lineIndex < 0 || lineIndex >= len(ts.Lines) {
		panic(ers.NewErrInvalidParameter(
			"lineIndex", ers.NewErrOutOfBound(0, len(ts.Lines)-1, lineIndex).
				WithLowerBound(true).WithUpperBound(false),
		))
	}

	// Check if adding the next word to the last line exceeds the width.
	// If it does, we shift the words of the last line to the left.
	for !ts.CanInsertWord(word, lineIndex) && lineIndex >= 0 {
		firstWord := ts.Lines[lineIndex].shiftLeft()
		ts.Lines[lineIndex].InsertWord(word)
		word = firstWord

		lineIndex--
	}

	// FIXME: Check why this is needed
	ts.CanInsertWord(word, lineIndex)
	ts.Lines[lineIndex].InsertWord(word)
}

// InsertWord is a method on the TextSplitter struct.
// It attempts to insert a given word into the TextSplitter.
// If the length of the word is greater than the width of the TextSplitter,
// it returns an ErrWordTooLong error.
// If the word can be inserted without exceeding the width, it is added to
// the Lines slice and the method returns nil.
// If adding the word to the last line would exceed the width, the words of
// the last line are shifted to the left until the word can be inserted or
// there are no more lines.
func (ts *TextSplitter) InsertWord(word string) bool {
	if len(ts.Lines) < cap(ts.Lines) {
		if utf8.RuneCountInString(word) > ts.Width {
			return false
		}

		ts.Lines = append(ts.Lines, &SpltLine{
			Line: []string{word},
			Len:  utf8.RuneCountInString(word),
		})

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

	// FIXME: Check why this is needed
	ts.CanInsertWord(word, lastLineIndex)
	ts.Lines[lastLineIndex].InsertWord(word)

	return true
}

// deepCopy is a method on the TextSplitter struct.
// It creates a deep copy of the TextSplitter and returns a pointer to it.
// The method iterates over all lines in the TextSplitter and adds a deep
// copy of each line to the Lines slice of the new TextSplitter.
func (ts *TextSplitter) deepCopy() *TextSplitter {
	newTs := &TextSplitter{
		Width: ts.Width,
		Lines: make([]*SpltLine, 0, len(ts.Lines)),
	}

	for _, line := range ts.Lines {
		newTs.Lines = append(newTs.Lines, line.deepCopy())
	}

	return newTs
}

// canShiftUp is a method on the TextSplitter struct.
// It checks if the first word of a given line can be shifted up to
// the previous line without exceeding the width.
// The method takes an integer (lineIndex) as input.
// It returns true if the first word of the line at lineIndex can be
// shifted up to the previous line without exceeding the width, and
// false otherwise.
func (ts *TextSplitter) canShiftUp(lineIndex int) bool {
	return ts.CanInsertWord(ts.Lines[lineIndex].Line[0], lineIndex-1)
}

// shiftUp is a method on the TextSplitter struct.
// It shifts the first word of a given line up to the previous line.
// The method takes an integer (lineIndex) as input.
// It does not return any value.
// The method assumes that it is safe to shift the word up, i.e., the
// canShiftUp method should be called before calling this method.
func (ts *TextSplitter) shiftUp(lineIndex int) {
	ts.Lines[lineIndex-1].InsertWord(ts.Lines[lineIndex].shiftLeft())
}

// CalculateNumberOfLines is a function that calculates the minimum number
// of lines needed
// to fit a given text within a specified width.
// It takes a slice of strings (text) and an integer (width) as input.
// The function returns the calculated number of lines and an error.
//
// If the input text is empty, the function returns 0 and an ErrEmptyText
// error.
// If the width is less than or equal to 0, the function returns 0 and an
// ErrWidthTooSmall error.
//
// The function calculates the total length of the text (Tl) and uses a
// mathematical formula to estimate the minimum number of lines needed to
// fit the text within the given width.
// The formula is explained in detail in the comments within the function.
//
// If the calculated number of lines is greater than the number of words in
// the text, the function returns 0 and an ErrWidthTooSmall error, indicating
// that the width is too small.
//
// Otherwise, the function returns the calculated number of lines and no error.
func CalculateNumberOfLines(text []string, width int) (int, bool) {
	if len(text) == 0 {
		panic(ers.NewErrInvalidParameter(
			"text", errors.New("text cannot be empty"),
		))
	}
	if width <= 0 {
		panic(ers.NewErrInvalidParameter(
			"width", fmt.Errorf("negative or zero width (%d) is not allowed", width),
		))
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

	var Tl int

	for _, word := range text {
		Tl += utf8.RuneCountInString(word) + 1 // +1 for the space or any other separator
	}
	Tl-- // Remove the last space or separator

	numberOfLines := int(math.Ceil(float64(Tl-width)/float64(width+1))) + 1

	return numberOfLines, numberOfLines > len(text)
}

// SplitTextInEqualSizedLines is a function that splits a given text into
// lines of equal width.
// It takes a slice of strings (text) and two integers (width, height) as
// input.
// The function returns a pointer to a TextSplitter struct and an error.
//
// If the input text is empty, the function returns nil and an ErrEmptyText
// error.
// If the width is less than or equal to 0, the function returns
// nil and an ErrWidthTooSmall error.
//
// If the height is less than 1, the function calculates the height itself
// using the CalculateNumberOfLines function.
// If CalculateNumberOfLines returns an error, SplitTextInEqualSizedLines
// returns nil and the same error.
//
// The function then creates a new TextSplitter struct and starts inserting
// words into it.
// If a word cannot be inserted, the function returns an error.
//
// After all words have been inserted, the function calculates the Sum of
// Squared Mean (SQM) for the current solution and starts looking for an
// optimal solution by moving words from the last line to the above line.
//
// The function keeps track of all candidates for the optimal solution and
// their corresponding SQM.
// The candidate with the lowest SQM is considered the optimal solution.
//
// If there are multiple candidates with the same lowest SQM, the function
// simply returns the first one.
//
// The function finally returns the first line of the optimal solution and no error.
func SplitTextInEqualSizedLines(text []string, width int, maxHeight optional.Int) (ts *TextSplitter, err error) {
	defer ers.RecoverFromPanic(&err)

	ers.Check(len(text), ers.NewErrInvalidParameter(
		"text", errors.New("text cannot be empty"),
	))
	ers.Check(width <= 0, ers.NewErrInvalidParameter(
		"width", fmt.Errorf("negative or zero width (%d) is not allowed", width),
	))

	height := maxHeight.OrElse(ers.CheckPred[int, int](width, func(width int) (int, bool) {
		return CalculateNumberOfLines(text, width)
	}, nil))
	ers.Check(height < 1, ers.NewErrInvalidParameter(
		"height", fmt.Errorf("negative or zero height (%d) is not allowed", height),
	))

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

	group := &TextSplitter{
		Width: width,
		Lines: make([]*SpltLine, 0, height),
	}

	index := slices.IndexFunc(text, group.InsertWord)
	ers.Check(index != -1, ers.NewErrOperationFailed(
		"insert word", fmt.Errorf("word at index %d (%s) is too long", index, text[index]),
	))

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

	// 4.2. Calculate the SQM of each candidate and returns the ones with the lowest SQM.
	calculateSQMOf := func(ts *TextSplitter) float64 {
		var average float64

		for _, line := range ts.Lines {
			average += float64(line.Len)
		}

		average /= float64(len(ts.Lines))

		var sqm float64

		for _, line := range ts.Lines {
			sqm += math.Pow(float64(line.Len)-average, 2)
		}

		return math.Sqrt(sqm / float64(len(ts.Lines)))
	}

	candidateList := make([]*TextSplitter, 0)
	newCandidates := []*TextSplitter{group}
	minSQM := math.MaxFloat64

	for len(newCandidates) > 0 {
		candidatesToCheck := make([]*TextSplitter, len(newCandidates))
		copy(candidatesToCheck, newCandidates)

		newCandidates = make([]*TextSplitter, 0) // Clear the slice as we don't need it anymore.

		for _, candidate := range candidatesToCheck {
			if SQM := calculateSQMOf(candidate); SQM < minSQM {
				minSQM = SQM
				candidateList = []*TextSplitter{candidate}
			} else if SQM == minSQM {
				candidateList = append(candidateList, candidate)
			}

			for i := 1; i < height; i++ {
				// Check if the first word of the line can be moved to the above line.
				// If yes, then it is a candidate for the optimal solution.
				if !candidate.canShiftUp(i) {
					continue
				}

				// Copy the candidate as we don't want to modify the original one.
				candidateCopy := candidate.deepCopy()
				candidateCopy.shiftUp(i)
				newCandidates = append(newCandidates, candidateCopy)
			}
		}
	}

	// 5. If we still have more than one candidate, we have to choose one
	// of them by following other criteria.
	// However, for now, we will just choose the first one.
	//
	// TODO: Choose the best candidate by following other criteria.

	// 6. Return the first line of the candidate.
	return candidateList[0], nil
}

// SplitSentenceIntoFields takes a string and an indent level as input.
// It splits the string into fields, where each field is a substring separated by one or more whitespace characters.
// The function also handles special characters such as tabs, vertical tabs, carriage returns, line feeds, and form feeds.
// The indent level determines the number of spaces that a tab character is replaced with.
//
// The function returns a two-dimensional slice of strings.
// Each inner slice represents the fields of a line from the input string.
// Lines that do not contain any fields (i.e., they are empty or consist only of whitespace) are not included in the output.
//
// If the input string is empty, the function returns an empty two-dimensional slice.
// If the indent level is negative, the function returns an ErrInvalidParameter error.
// If the input string contains an invalid rune, the function returns an error.
func SplitSentenceIntoFields(sentence string, indentLevel int) [][]string {
	if sentence == "" {
		return [][]string{}
	}

	ers.Check(indentLevel < 0, ers.NewErrInvalidParameter(
		"indentLevel", errors.New("indent level cannot be negative"),
	))

	lines := make([][]string, 0)
	words := make([]string, 0)

	var builder strings.Builder

	for j := 0; len(sentence) > 0; j++ {
		char, size := utf8.DecodeRuneInString(sentence)
		sentence = sentence[size:]

		ers.Check(char == utf8.RuneError, ers.NewErrOperationFailed(
			"rune at index", fmt.Errorf("rune at index %d is invalid", j),
		))

		switch char {
		case '\t':
			builder.WriteString(strings.Repeat(" ", indentLevel)) // 3 spaces
		case '\v':
			// Do nothing
		case '\r':
			if utf8.RuneCountInString(sentence) > 0 {
				nextRune, size := utf8.DecodeRuneInString(sentence)

				if nextRune == '\n' {
					sentence = sentence[size:]
				}
			}

			fallthrough
		case '\n', '\u0085':
			if builder.Len() != 0 {
				words = append(words, builder.String())
				builder.Reset()
			}

			lines = append(lines, words)
			words = make([]string, 0)
		case ' ':
			if builder.Len() != 0 {
				words = append(words, builder.String())
				builder.Reset()
			}
		case '\u00A0':
			builder.WriteRune(' ')
		case '\f':
			if builder.Len() != 0 {
				words = append(words, builder.String())
				builder.Reset()
			}

			lines = append(lines, words)
			words = make([]string, 0)

			lines = append(lines, []string{string(char)})
		default:
			builder.WriteRune(char)
		}
	}

	if builder.Len() != 0 {
		words = append(words, builder.String())
	}

	if len(words) > 0 {
		lines = append(lines, words)
	}

	return lines
}
