package Components

import (
	"errors"
	"fmt"
	"strings"

	cdd "github.com/PlayerR9/MyGoLib/ComplexData/Display"

	sx "github.com/PlayerR9/MyGoLib/Units/String"

	"github.com/gdamore/tcell"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

const (
	Asterisks string = "***"

	AsterisksLen int = len(Asterisks)

	TitleMinWidth int = 2 * (AsterisksLen + 1)
)

// Title represents the header of a process or application.
// It contains information about the title, current process, counters, message buffer,
// channels for receiving messages and errors, synchronization primitives, and the width
// of the header.
type Title struct {
	title    string
	subtitle string // empty string if no subtitle
	style    tcell.Style
}

// Draw draws the title to the draw table.
//
// Parameters:
//   - table: The draw table.
//   - x: The x coordinate to draw the title at.
//   - y: The y coordinate to draw the title at.
//
// Returns:
//   - error: An error if the title could not be drawn.
func (t *Title) Draw(table *cdd.DrawTable, x, y int) error {
	if table == nil {
		return ers.NewErrNilParameter("table")
	}

	// 1. Generate the full title
	var fullTitle *sx.String

	if t.subtitle == "" {
		fullTitle = sx.NewString(t.title)
	} else {
		var builder strings.Builder

		builder.WriteString(t.title)
		builder.WriteRune(' ')
		builder.WriteRune('-')
		builder.WriteRune(' ')
		builder.WriteString(t.subtitle)

		fullTitle = sx.NewString(builder.String())
	}

	// 2. Generate the lines
	width := table.GetWidth()

	lines, err := t.tryToFitLines(width, x, fullTitle)
	if err != nil {
		return fmt.Errorf("failed to fit lines: %s", err.Error())
	}

	// 3. Check if the lines fit in the draw table
	if len(lines) > table.GetHeight() {
		return errors.New("lines do not fit in draw table")
	}

	// 4. Write the lines with centered alignment
	for i := 0; i < len(lines); i++ {
		startPos := (width - lines[i].GetLength()) / 2

		cell := cdd.NewDtUnit(lines[i], t.style)

		err := cell.Draw(table, startPos, i)
		if err != nil {
			return fmt.Errorf("failed to draw line %d: %s", i, err.Error())
		}
	}

	return nil
}

// ForceDraw draws the title to the draw table.
//
// Parameters:
//   - table: The draw table.
//   - x: The x coordinate to draw the title at.
//   - y: The y coordinate to draw the title at.
//
// Returns:
//   - error: An error if the title could not be drawn.
func (t *Title) ForceDraw(table *cdd.DrawTable, x, y int) error {
	if table == nil {
		return ers.NewErrNilParameter("table")
	}

	// 1. Generate the full title
	var fullTitle *sx.String

	if t.subtitle == "" {
		fullTitle = sx.NewString(t.title)
	} else {
		var builder strings.Builder

		builder.WriteString(t.title)
		builder.WriteRune(' ')
		builder.WriteRune('-')
		builder.WriteRune(' ')
		builder.WriteString(t.subtitle)

		fullTitle = sx.NewString(builder.String())
	}

	// 2. Generate the lines
	width := table.GetWidth()

	lines, err := forceGenerateLines(fullTitle, width, x)
	if err != nil {
		return fmt.Errorf("failed to generate lines: %s", err.Error())
	}

	// 3. Write the lines with centered alignment
	for i := 0; i < len(lines); i++ {
		startPos := (width - lines[i].GetLength()) / 2

		cell := cdd.NewDtUnit(lines[i], t.style)

		err := cell.ForceDraw(table, startPos, i)
		if err != nil {
			return fmt.Errorf("failed to draw line %d: %s", i, err.Error())
		}
	}

	return nil
}

// NewTitle creates a new Title with the given title and a style.
//
// Parameters:
//   - title: The title of the new Title.
//   - style: The style of the new Title.
//
// Returns:
//   - *Title: The new Title.
func NewTitle(title string, style tcell.Style) *Title {
	return &Title{
		title:    title,
		subtitle: "",
		style:    style,
	}
}

// SetSubtitle sets the subtitle of the Title.
//
// Parameters:
//   - subtitle: The new subtitle.
//
// Behaviors:
//   - If the subtitle is an empty string, the subtitle is removed.
func (t *Title) SetSubtitle(subtitle string) {
	t.subtitle = subtitle
}

// generateLines is a helper method that generates the lines of the title.
//
// Parameters:
//   - fullTitle: The full title.
//   - width: The width of the lines.
//
// Returns:
//   - []*sx.String: The lines of the title.
//   - error: An error if the full title could not be split in lines.
func generateLines(fullTitle *sx.String, width int, x int) ([]*sx.String, error) {
	contents := fullTitle.Fields()

	numberOfLines, err := sx.CalculateNumberOfLines(contents, width-x-TitleMinWidth)
	if err != nil {
		return nil, fmt.Errorf("could not calculate number of lines: %s", err.Error())
	}

	ts, err := sx.SplitInEqualSizedLines(contents, width-x-TitleMinWidth, numberOfLines)
	if err != nil {
		return nil, fmt.Errorf("could not split text in equal sized lines: %s", err.Error())
	}

	lines := ts.GetLines()

	for _, line := range lines {
		line.PrependRune(' ')
		line.PrependString(Asterisks)
		line.AppendRune(' ')
		line.AppendString(Asterisks)
	}

	return lines, nil
}

// tryToFitLines is a helper method that tries to fit the full title in the draw table.
//
// Parameters:
//   - table: The draw table.
//   - fullTitle: The full title.
//
// Returns:
//   - []*sx.String: The lines of the title.
//   - error: An error if the full title could not be split in lines.
func (t *Title) tryToFitLines(width int, x int, fullTitle *sx.String) ([]*sx.String, error) {
	lines, err := generateLines(fullTitle, width, x)
	if err == nil {
		return lines, nil
	}

	fullTitle = fullTitle.TrimEnd(width - x - TitleMinWidth)

	ok := fullTitle.ReplaceSuffix(Hellip)
	if !ok {
		return nil, errors.New("hellip is longer than the full title")
	}

	fullTitle.PrependRune(' ')
	fullTitle.PrependString(Asterisks)
	fullTitle.AppendRune(' ')
	fullTitle.AppendString(Asterisks)

	return []*sx.String{fullTitle}, nil
}

// forceGenerateLines is a helper method that generates the lines of the title.
//
// Parameters:
//   - fullTitle: The full title.
//   - width: The width of the lines.
//
// Returns:
//   - []*sx.String: The lines of the title.
//   - error: An error if the full title could not be split in lines.
func forceGenerateLines(fullTitle *sx.String, width int, x int) ([]*sx.String, error) {
	contents := fullTitle.Fields()

	numberOfLines, err := sx.CalculateNumberOfLines(contents, width-x-TitleMinWidth)
	if err != nil && !ers.As[*sx.ErrLinesGreaterThanWords](err) {
		return nil, fmt.Errorf("could not calculate number of lines: %s", err.Error())
	}

	ts, err := sx.SplitInEqualSizedLines(contents, width-x-TitleMinWidth, numberOfLines)
	if err != nil {
		return nil, fmt.Errorf("could not split text in equal sized lines: %s", err.Error())
	}

	lines := ts.GetLines()

	for _, line := range lines {
		line.PrependRune(' ')
		line.PrependString(Asterisks)
		line.AppendRune(' ')
		line.AppendString(Asterisks)
	}

	return lines, nil
}

// GetTitle returns the title of the Title.
//
// Returns:
//   - string: The title of the Title.
func (t *Title) GetTitle() string {
	return t.title
}

// GetSubtitle returns the subtitle of the Title.
//
// Returns:
//   - string: The subtitle of the Title.
func (t *Title) GetSubtitle() string {
	return t.subtitle
}
