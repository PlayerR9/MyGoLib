package ContentBox

import (
	"fmt"
	"strings"
	"unicode/utf8"

	cdd "github.com/PlayerR9/MyGoLib/ComplexData/Display"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"

	fs "github.com/PlayerR9/MyGoLib/Formatting/Strings"
	sext "github.com/PlayerR9/MyGoLib/Utility/StringExt"

	"github.com/gdamore/tcell"
	"github.com/markphelps/optional"
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
	Space rune = ' '

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

type ContentBox struct {
	lines    [][]string
	maxWidth int
	sep      optional.Rune
}

type ContentBoxOption func(*ContentBox)

func WithSeparator(separator rune) ContentBoxOption {
	return func(cb *ContentBox) {
		cb.sep = optional.NewRune(separator)
	}
}

func WithLines(lines [][]string) ContentBoxOption {
	return func(cb *ContentBox) {
		if lines == nil {
			cb.lines = make([][]string, 0)
		} else {
			cb.lines = lines
		}
	}
}

func NewContentBox(maxWidth, maxHeight int, options ...ContentBoxOption) ([]*fs.TextSplit, error) {
	if maxWidth <= 0 {
		return nil, ers.NewErrInvalidParameter(
			"maxWidth",
			ers.NewErrGT(0),
		)
	}

	cb := &ContentBox{
		lines:    make([][]string, 0),
		maxWidth: maxWidth,
	}

	for _, option := range options {
		option(cb)
	}

	finalTs := make([]*fs.TextSplit, 0, len(cb.lines))

	for _, line := range cb.lines {
		var sentences [][]string

		if !cb.sep.Present() {
			sentences = [][]string{line}
		} else {
			s := cb.sep.MustGet()

			for _, field := range line {
				sentences = append(sentences, strings.FieldsFunc(field, func(r rune) bool {
					return r == s
				}))
			}
		}

		ts, err := cb.createTextSplitter(sentences, maxHeight)
		if err != nil {
			return nil, err
		}

		// If it is possible to optimize the text, optimize it.
		// Otherwise, the unoptimized text is also fine.
		optimizedTs, err := cb.finalOptimization(ts)
		if err != nil {
			finalTs = append(finalTs, ts)
		} else {
			finalTs = append(finalTs, optimizedTs)
		}
	}

	return finalTs, nil
}

// processFirstLine processes the first line of text represented by
// a slice of fields.
// It calculates the number of lines the text would occupy if split
// into lines of a specified width. If the text cannot be split into
// lines of the specified width, it replaces the suffix of the text
// with a hellip and adds the resulting line to the TextSplitter. If
// the text can be split into more than one line, it creates a new
// line with the first field and as many subsequent fields as can fit
// into the line width, adding a hellip if necessary. If the text can
// be split into exactly one line, it splits the text into equal-sized
// lines and adds the first line to the TextSplitter.
//
// The function returns the updated TextSplitter, a boolean indicating
// whether the text was truncated, and an error. If an error occurs
// while replacing the suffix or splitting the text, the error is
// returned. Otherwise, the error is nil.
func (cb *ContentBox) processFirstLine(words []string, ts *fs.TextSplit) (*fs.TextSplit, bool, error) {
	numberOfLines, err := fs.CalculateNumberOfLines(words, cb.maxWidth)
	if err != nil {
		var line string

		if !cb.sep.Present() {
			line = strings.Join(words, "")
		} else {
			line = strings.Join(words, string(cb.sep.MustGet()))
		}

		line, err = sext.ReplaceSuffix(line[:cb.maxWidth], Hellip)
		if err != nil {
			return nil, false, fmt.Errorf("could not replace suffix: %v", err)
		}

		ok := ts.InsertWord(line)
		if !ok {
			panic("could not insert word")
		}

		return ts, true, nil
	}

	if numberOfLines > 1 {
		wordsProcessed := []string{words[0]}
		wpLen := utf8.RuneCountInString(words[0])

		var nextField string

		for i, currentField := range words[1 : len(words)-1] {
			nextField = words[i+1]

			totalLen := wpLen + 2 + utf8.RuneCountInString(currentField) +
				utf8.RuneCountInString(nextField)

			if totalLen+HellipLen > cb.maxWidth {
				wordsProcessed = append(wordsProcessed, currentField+Hellip)
				wpLen += utf8.RuneCountInString(currentField) + HellipLen + 1
				break
			}

			wordsProcessed = append(wordsProcessed, currentField)
			wpLen += utf8.RuneCountInString(currentField) + 1
		}

		if wpLen+1+utf8.RuneCountInString(nextField)+HellipLen <= cb.maxWidth {
			wordsProcessed = append(wordsProcessed, nextField)
			wpLen += utf8.RuneCountInString(nextField) + 1
		}

		for _, word := range wordsProcessed {
			ok := ts.InsertWord(word)
			if !ok {
				panic("could not insert word")
			}
		}
	} else {
		halfTs, err := fs.SplitInEqualSizedLines(
			words, cb.maxWidth, numberOfLines,
		)
		if err != nil {
			return nil, false, fmt.Errorf("could not split text: %v", err)
		}

		for _, words := range halfTs.GetFirstLine().GetWords() {
			ok := ts.InsertWord(words)
			if !ok {
				panic("could not insert word")
			}
		}
	}

	return ts, false, nil
}

// processOtherLines processes the other lines of text (i.e., all lines
// except the first) represented by a slice of fields. It calculates
// the number of lines the text would occupy if split into lines of a
// specified width. If the text cannot be split into lines of the
// specified width, it replaces the suffix of the text with a hellip and
// adds the resulting line to the TextSplitter. If the text can be split
// into more than one line, it creates a new line with the first field
// and as many subsequent fields as can fit into the line width, adding
// a hellip if necessary. If the text can be split into exactly one line,
// it splits the text into equal-sized lines and adds the first line to
// the TextSplitter.
//
// The function returns the updated TextSplitter, a boolean indicating
// whether the text was truncated, and an error. If an error occurs while
// replacing the suffix or splitting the text, the error is returned.
// Otherwise, the error is nil.
func (cb *ContentBox) processOtherLines(fields []string, ts *fs.TextSplit, possibleNewLine bool) (*fs.TextSplit, bool, error) {
	if possibleNewLine {
		for len(fields) > 0 {
			ok := ts.InsertWord(fields[0])
			if !ok {
				break
			}

			fields = fields[1:]
		}
	}

	if len(fields) == 0 {
		return ts, possibleNewLine, nil
	}

	numberOfLines, err := fs.CalculateNumberOfLines(fields, cb.maxWidth-IndentLevel)
	if err != nil {
		var line string

		if !cb.sep.Present() {
			line = strings.Join(fields, "")
		} else {
			line = strings.Join(fields, string(cb.sep.MustGet()))
		}

		line, err = sext.ReplaceSuffix(line[:cb.maxWidth-IndentLevel], Hellip)
		if err != nil {
			return nil, false, fmt.Errorf("could not replace suffix: %v", err)
		}

		ok := ts.InsertWord(line)
		if !ok {
			panic("could not insert word")
		}
	} else if numberOfLines > 1 {
		splt := sext.NewSpltLine(fields[0])

		for _, field := range fields[1:] {
			if splt.Len+1+utf8.RuneCountInString(field)+HellipLen+IndentLevel > cb.maxWidth {
				splt.InsertWord(field + Hellip)
				break
			}

			splt.InsertWord(field)
		}

		ts.Lines = append(ts.Lines, splt)

		return ts, true, nil
	}

	halfTs, err := fs.SplitInEqualSizedLines(
		fields, cb.maxWidth-IndentLevel, numberOfLines,
	)
	if err != nil {
		return nil, false, fmt.Errorf("could not split text: %v", err)
	}

	ok := ts.InsertWord(halfTs.GetFirstLine().String())
	if !ok {
		panic("could not insert word")
	}

	return ts, false, nil
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
func (cb *ContentBox) createTextSplitter(lines [][]string, maxHeight int) (*fs.TextSplit, error) {
	ts, err := fs.NewTextSplit(cb.maxWidth-IndentLevel, maxHeight)
	if err != nil {
		return nil, fmt.Errorf("could not create TextSplitter: %v", err)
	}

	possibleNewLine := false

	ts, possibleNewLine, err = cb.processFirstLine(lines[0], ts)
	if err != nil {
		return nil, err
	}

	for _, line := range lines[1:] {
		ts, possibleNewLine, err = cb.processOtherLines(line, ts, possibleNewLine)
		if err != nil {
			return nil, err
		}
	}

	return ts, nil
}

// finalOptimization takes a TextSplitter and a width, and optimizes the
// lines of text in the TextSplitter by splitting them into equal-sized
// lines of the specified width.
// It first converts the lines of the TextSplitter into a slice of strings,
// and then calls the SplitTextInEqualSizedLines function with this slice,
// the specified width, and -1 as the arguments.
//
// The function returns a new TextSplitter resulting from the optimization
// and an error.
// If an error occurs during the optimization, the error is returned.
// Otherwise, the error is nil.
func (cb *ContentBox) finalOptimization(ts *fs.TextSplit) (*fs.TextSplit, error) {
	text := make([]string, 0)

	for _, line := range ts.GetLines() {
		text = append(text, line.String())
	}

	return fs.SplitInEqualSizedLines(text, cb.maxWidth, -1)
}

// writeLines takes a TextSplitter, a tcell.Style, and a WriteOnlyDrawTable,
// and writes the lines of text from the TextSplitter into the drawTable
// with the specified style.
// It first calculates the rightmost limit of the text, and then for each
// line of the drawTable, it writes the corresponding line of text from the
// TextSplitter at the beginning of the line, and fills the rest of the line
// with spaces. The first line of the drawTable is written at the beginning,
// while the rest of the lines are indented by a constant IndentLevel.
func writeLines(ts *fs.TextSplit, style tcell.Style, table *cdd.DrawTable, yCoord int) error {
	height := ts.GetHeight()

	if height > table.GetHeight() {
		return ers.NewErrInvalidParameter(
			"ts.Lines",
			fmt.Errorf("length (%d) exceeds height (%d)", height, table.GetHeight()),
		)
	}

	rightMostLimit, ok := ts.GetFurthestRightEdge()
	if !ok {
		panic("could not get furthest right edge")
	}

	var line, emptyLine string

	// First line
	line = ts.GetFirstLine().String()
	emptyLine = strings.Repeat(string(Space), rightMostLimit-len(line))

	table.WriteLineAt(0, yCoord, line, style, true)
	table.WriteLineAt(len(line), yCoord, emptyLine, style, true)

	// Rest of the lines
	tsLines := ts.GetLines()

	for i := 1; i < len(tsLines); i++ {
		line = tsLines[i].String()
		emptyLine = strings.Repeat(string(Space), rightMostLimit-len(line))

		table.WriteLineAt(IndentLevel, i+yCoord, line, style, true)
		table.WriteLineAt(IndentLevel+len(line), i+yCoord, emptyLine, style, true)
	}

	return nil
}

type MakeContentBoxFunc func() ([]*fs.TextSplit, error)

func ContentBoxWrite(table *cdd.DrawTable, style tcell.Style, make MakeContentBoxFunc) error {
	tss, err := make()
	if err != nil {
		return fmt.Errorf("could not fit lines: %v", err)
	}

	totalHeight := 0
	tableHeight := table.GetHeight()

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

		if err := writeLines(ts, style, table, yCoord); err != nil {
			return fmt.Errorf("could not write lines: %v", err)
		}

		yCoord += currentHeight
	}

	return nil
}
