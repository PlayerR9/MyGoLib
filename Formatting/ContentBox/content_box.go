package Components

import (
	"fmt"

	cdd "github.com/PlayerR9/MyGoLib/ComplexData/Display"

	sx "github.com/PlayerR9/MyGoLib/Units/String"

	"github.com/gdamore/tcell"
)

const (
	// Hellip defines the string to be used as an ellipsis when the content
	// of the ContentBox is truncated.
	// It is set to "...", which is the standard representation of an ellipsis
	// in text.
	Hellip string = "..."

	// HellipLen defines the length of the Hellip string.
	// It is set to 3, which is the number of characters in the Hellip string.
	HellipLen int = len(Hellip)

	// Space defines the string to be used as a space when writing content
	// into the ContentBox.
	// It is set to " ", which is the standard representation of a space in
	// text.
	Space string = " "

	// FieldSpacing defines the number of spaces between each field (word)
	// when they are written into the ContentBox.
	// It is set to 1, meaning there will be one spaces between each field.
	FieldSpacing int = 1

	// IndentLevel defines the number of spaces used for indentation when
	// writing content into the ContentBox.
	// It is set to 2, meaning there will be two spaces at the start of each
	// new line of content.
	IndentLevel int = 3
)

// ContentBox represents a box that contains content.
type ContentBox struct {
	// lines is a two-dimensional slice of strings representing the content
	// of the box.
	lines [][]*sx.String

	// style is the style to be used when writing the content into the box.
	style tcell.Style
}

// Draw draws the content of the ContentBox into the specified draw table.
//
// Parameters:
//   - table - the draw table to draw the content into.
//   - x - the x coordinate to start drawing the content at.
//   - y - the y coordinate to start drawing the content at.
//
// Returns:
//   - error - an error if the content could not be drawn.
func (cb *ContentBox) Draw(table *cdd.DrawTable, x, y int) error {
	maxWidth := table.GetWidth() - x
	maxHeight := table.GetHeight() - y

	tss, err := cb.apply(maxWidth, maxHeight)
	if err != nil {
		return err
	}

	totalHeight := 0
	tableHeight := maxHeight

	for _, ts := range tss {
		totalHeight += ts.GetHeight()

		if totalHeight > tableHeight {
			break
		}
	}

	yCoord := 0

	for _, ts := range tss {
		currentHeight := ts.GetHeight()

		shouldExit := currentHeight+yCoord > tableHeight
		if shouldExit {
			return nil
		}

		cell := cdd.NewDtUnit(ts, cb.style)

		err := cell.Draw(table, x, y+yCoord)
		if err != nil {
			return fmt.Errorf("could not draw cell: %s", err.Error())
		}

		yCoord += currentHeight
	}

	return nil
}

// Draw draws the content of the ContentBox into the specified draw table.
//
// Parameters:
//   - table - the draw table to draw the content into.
//   - x - the x coordinate to start drawing the content at.
//   - y - the y coordinate to start drawing the content at.
//
// Returns:
//   - error - an error if the content could not be drawn.
func (cb *ContentBox) ForceDraw(table *cdd.DrawTable, x, y int) error {
	maxWidth := table.GetWidth() - x
	maxHeight := table.GetHeight() - y

	tss, err := cb.forceApply(maxWidth, maxHeight)
	if err != nil {
		return err
	}

	totalHeight := 0
	tableHeight := maxHeight

	for _, ts := range tss {
		totalHeight += ts.GetHeight()

		if totalHeight > tableHeight {
			break
		}
	}

	yCoord := 0

	for _, ts := range tss {
		currentHeight := ts.GetHeight()

		shouldExit := currentHeight+yCoord > tableHeight
		if shouldExit {
			return nil
		}

		cell := cdd.NewDtUnit(ts, cb.style)

		err := cell.ForceDraw(table, x, y+yCoord)
		if err != nil {
			return fmt.Errorf("could not draw cell: %s", err.Error())
		}

		yCoord += currentHeight
	}

	return nil
}

// NewContentBox creates a new ContentBox with the specified lines of content.
//
// Parameters:
//   - lines - a two-dimensional slice of strings representing the content of the box.
//   - style - the style to be used when writing the content into the box.
//
// Returns:
//   - *ContentBox - a pointer to the created ContentBox.
func NewContentBox(lines [][]*sx.String, style tcell.Style) *ContentBox {
	return &ContentBox{
		lines: lines,
		style: style,
	}
}

// processLine processes a line of text represented by a slice of fields.
// It calculates the number of lines the text would occupy if split into
// lines of a specified width. If the text cannot be split into lines of
// the specified width, it replaces the suffix of the text with a hellip
// and adds the resulting line to the TextSplitter. If the text can be split
// into more than one line, it creates a new line with the first field and
// as many subsequent fields as can fit into the line width, adding a hellip
// if necessary. If the text can be split into exactly one line, it splits
// the text into equal-sized lines and adds the first line to the TextSplitter.
//
// Parameters:
//   - isFirst - a boolean indicating whether the line is the first line of text.
//   - maxWidth - the maximum width of the line.
//   - ts - the TextSplitter to add the line to.
//   - words - a slice of fields representing the line of text.
//
// Returns:
//   - *sx.TextSplit - the updated TextSplitter.
//   - bool - a boolean indicating whether the text was truncated.
//   - error - an error if the text could not be processed.
func (cb *ContentBox) processLine(isFirst bool, maxWidth int, ts *sx.TextSplit, words []*sx.String, possibleNewLine bool) (*sx.TextSplit, bool, error) {
	if !isFirst {
		maxWidth -= IndentLevel
	}

	numberOfLines, err := sx.CalculateNumberOfLines(words, maxWidth)

	if err != nil {
		line := sx.Join(words, "").TrimEnd(maxWidth)

		ok := line.ReplaceSuffix(Hellip)
		if !ok {
			return nil, false, fmt.Errorf("suffix is bigger than maxWidth")
		}

		ok = ts.InsertWord(line)
		if !ok {
			panic("could not insert word")
		}

		return ts, true, nil
	}

	if numberOfLines > 1 {
		wordsProcessed := []*sx.String{words[0]}
		wpLen := words[0].GetLength()

		var nextField *sx.String

		for i, currentField := range words[1 : len(words)-1] {
			nextField = words[i+1]

			totalLen := wpLen + 2 + currentField.GetLength() +
				nextField.GetLength()

			if totalLen+HellipLen > maxWidth {
				currentField.AppendString(Hellip)

				wordsProcessed = append(wordsProcessed, currentField)
				wpLen += currentField.GetLength() + 1
				break
			}

			wordsProcessed = append(wordsProcessed, currentField)
			wpLen += currentField.GetLength() + 1
		}

		if wpLen+1+nextField.GetLength()+HellipLen <= maxWidth {
			wordsProcessed = append(wordsProcessed, nextField)
			wpLen += nextField.GetLength() + 1
		}

		firstNotInserted := ts.InsertWords(wordsProcessed)
		if firstNotInserted != -1 {
			panic(fmt.Sprintf("could not insert word %s", wordsProcessed[firstNotInserted]))
		}

		return ts, true, nil
	} else {
		halfTs, err := sx.SplitInEqualSizedLines(
			words, maxWidth, numberOfLines,
		)

		if err != nil {
			return nil, false, fmt.Errorf("could not split text: %s", err.Error())
		}

		wordsProcessed := halfTs.GetFirstLine()

		firstNotInserted := ts.InsertWords(wordsProcessed)
		if firstNotInserted != -1 {
			panic(fmt.Sprintf("could not insert word %s", wordsProcessed[firstNotInserted]))
		}

		return ts, false, nil
	}
}

// createTextSplitter takes a two-dimensional slice of strings
// representing a list of fields and a width, and creates a
// TextSplitter that splits the fields into lines of the specified
// width. It processes the first line of fields separately from
// the other lines. If an error occurs while processing a line,
// it returns an error with a message indicating the line number
// and the original error.
//
// The function returns a pointer to the created TextSplitter
// and an error. If no errors occur during the creation of the
// TextSplitter, the error is nil.
func (cb *ContentBox) createTextSplitter(lines [][]*sx.String, maxWidth, maxHeight int) (*sx.TextSplit, error) {
	ts, err := sx.NewTextSplit(maxWidth-IndentLevel, maxHeight)
	if err != nil {
		return nil, fmt.Errorf("could not create TextSplitter: %s", err.Error())
	}

	possibleNewLine := false

	ts, possibleNewLine, err = cb.processLine(true, maxWidth, ts, lines[0], true)
	if err != nil {
		return nil, err
	}

	for _, line := range lines[1:] {
		if possibleNewLine {
			for len(line) > 0 {
				ok := ts.InsertWord(line[0])
				if !ok {
					break
				}

				line = line[1:]
			}
		}

		if len(line) == 0 {
			continue
		}

		ts, possibleNewLine, err = cb.processLine(false, maxWidth, ts, line, possibleNewLine)
		if err != nil {
			return nil, fmt.Errorf("could not process line: %s", err.Error())
		}
	}

	return ts, nil
}

// apply takes a maximum width and height, and applies the content of the ContentBox
// to the specified width and height. It splits the content into lines of the specified
// width, and optimizes the text if possible.
//
// Parameters:
//   - maxWidth - the maximum width of the content.
//   - maxHeight - the maximum height of the content.
//
// Returns:
//   - []*sx.TextSplit - a slice of TextSplit objects representing the optimized content.
//   - error - an error if the content could not be applied.
func (cb *ContentBox) apply(maxWidth, maxHeight int) ([]*sx.TextSplit, error) {
	finalTs := make([]*sx.TextSplit, 0, len(cb.lines))

	for _, line := range cb.lines {
		sentences := [][]*sx.String{line}

		ts, err := cb.createTextSplitter(sentences, maxWidth, maxHeight)
		if err != nil {
			return nil, err
		}

		// If it is possible to optimize the text, optimize it.
		// Otherwise, the unoptimized text is also fine.
		optimizedTs, err := sx.SplitInEqualSizedLines(ts.GetLines(), maxWidth, -1)
		if err != nil {
			finalTs = append(finalTs, ts)
		} else {
			finalTs = append(finalTs, optimizedTs)
		}
	}

	return finalTs, nil
}

// apply takes a maximum width and height, and applies the content of the ContentBox
// to the specified width and height. It splits the content into lines of the specified
// width, and optimizes the text if possible.
//
// Parameters:
//   - maxWidth - the maximum width of the content.
//   - maxHeight - the maximum height of the content.
//
// Returns:
//   - []*sx.TextSplit - a slice of TextSplit objects representing the optimized content.
//   - error - an error if the content could not be applied.
func (cb *ContentBox) forceApply(maxWidth, maxHeight int) ([]*sx.TextSplit, error) {
	finalTs := make([]*sx.TextSplit, 0, len(cb.lines))

	for _, line := range cb.lines {
		sentences := [][]*sx.String{line}

		ts, err := cb.createTextSplitter(sentences, maxWidth, maxHeight)
		if err != nil {
			return nil, err
		}

		// If it is possible to optimize the text, optimize it.
		// Otherwise, the unoptimized text is also fine.
		optimizedTs, err := sx.SplitInEqualSizedLines(ts.GetLines(), maxWidth, -1)
		if err != nil {
			finalTs = append(finalTs, ts)
		} else {
			finalTs = append(finalTs, optimizedTs)
		}
	}

	return finalTs, nil
}
