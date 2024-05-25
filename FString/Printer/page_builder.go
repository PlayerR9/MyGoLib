package Printer

import (
	tls "github.com/PlayerR9/MyGoLib/FString/Section"
)

// pageBuilder is a type that represents a page of a document.
type pageBuilder struct {
	// sections are the sections of the page.
	sections []tls.Sectioner

	// buffer is the buffer of the page.
	buffer *sectionBuilder
}

// NewPage creates a new page.
//
// Returns:
//   - *Page: The new page.
func NewPage() *pageBuilder {
	return &pageBuilder{
		sections: make([]tls.Sectioner, 0),
		buffer:   nil,
	}
}

// AcceptLine is a function that accepts the current in-progress line
// (if any).
func (p *pageBuilder) AcceptLine() {
	if p.buffer == nil {
		return
	}

	p.buffer.Accept()
}

// IsFirstOfLine is a function that checks if the current buffer is the first of the line.
//
// Returns:
//   - bool: True if the current buffer is the first of the line.
func (p *pageBuilder) IsFirstOfLine() bool {
	if p.buffer == nil {
		return true
	}

	return p.buffer.IsFirstOfLine()
}

// Accept is a function that accepts the current in-progress buffer
// by converting it to the specified section type. Lastly, the section
// is added to the page.
//
// Parameters:
//   - sectionType: The section type to convert the buffer to.
//
// Returns:
//   - error: An error if the buffer could not be converted to the section type.
//
// Behaviors:
//   - Even when the buffer is empty, the section is still added to the page.
//     To avoid this, use the Finalize function.
func (p *pageBuilder) Accept(sectionType tls.Sectioner) error {
	if p.buffer != nil {
		p.buffer.AcceptWord()
	} else {
		p.buffer = newSectionBuilder()
	}

	// Convert the buffer to a section
	err := sectionType.FromTextBlock(p.buffer.GetWords())
	if err != nil {
		return err
	}

	// Add the section to the page
	p.sections = append(p.sections, sectionType)

	// Reset the buffer
	p.buffer = nil

	return nil
}

// AddString adds a string to the page.
//
// Parameters:
//   - str: The string to add.
//
// Behaviors:
//   - If the string is empty, it is not added.
func (p *pageBuilder) AddString(str string) {
	if str == "" {
		return
	}

	if p.buffer == nil {
		p.buffer = newSectionBuilder()
	}

	p.buffer.AddString(str)
}

// AddRune adds a rune to the page.
//
// Parameters:
//   - r: The rune to add.
func (p *pageBuilder) AddRune(r rune) {
	if p.buffer == nil {
		p.buffer = newSectionBuilder()
	}

	p.buffer.AddRune(r)
}

// AcceptWord is a function that accepts the current in-progress word
// if any.
func (p *pageBuilder) AcceptWord() {
	if p.buffer == nil {
		return
	}

	p.buffer.AcceptWord()
}

// Finalize is a function that finalizes the page by converting the buffer
// to the specified section type and adding it to the page.
//
// Parameters:
//   - sectionType: The section type to convert the buffer to.
//
// Returns:
//   - error: An error if the buffer could not be converted to the section type.
//
// Behaviors:
//   - If the buffer is empty, the function does nothing.
func (p *pageBuilder) Finalize(sectionType tls.Sectioner) error {
	if p.buffer == nil {
		return nil
	}

	p.buffer.AcceptWord()

	// Convert the buffer to a section
	err := sectionType.FromTextBlock(p.buffer.GetWords())
	if err != nil {
		return err
	}

	// Add the section to the page
	p.sections = append(p.sections, sectionType)

	// Reset the buffer
	p.buffer = nil

	return nil
}

//////////////////////////////////////////////////////////////

/*
// Iterator is a function that returns an iterator for the page.
//
// Returns:
//   - ui.Iterater[*SectionBuilder]: The iterator for the page.
func (p *PageBuilder) Iterator() ui.Iterater[*SectionBuilder] {
	return ui.NewGenericIterator(p.sections)
}
*/

/*
import (
	"strings"

	Tr "github.com/PlayerR9/MyGoLib/CustomData/Tray"
)

type Builder struct {
	tabStop  int
	allChars []rune
}

func (b *Builder) AddChar(char rune) {
	switch char {
	case '\u0000':
		// null : Ignore this character
	case '\a':
		// Bell : Ignore this character
	case '\b':
		// backspace : Remove the last character
	case '\t':
		// Tab : Add spaces until the next tab stop
	case '\n':
		// line feed : Add a new line or move to the left edge and down
	case '\v':
		// vertical tab : Add vertical tabulation
	case '\f':
		// form feed : Go to the next page
	case '\r':
		// carriage return : Move to the start of the line (alone)
		// or move to the start of the line and down (with line feed)
	case '\u001A':
		// Control-Z : End of file for Windows text-mode file i/o
	case '\u001B':
		// escape : Introduce an escape sequence (next character)
	default:
	}
}

func removeStartingChar(tray *Tr.Tray[rune], char rune) {
	tray.ArrowStart()

	for {
		c, err := tray.Read()
		if err != nil || c != char {
			return
		}

		err = tray.Delete(1)
		if err != nil {
			return
		}
	}
}

func removeTrailingChar(tray *Tr.Tray[rune], char rune) {
	tray.ArrowEnd()

	for {
		c, err := tray.Read()
		if err != nil || c != char {
			return
		}

		err = tray.Delete(1)
		if err != nil {
			return
		}
	}
}

func (b *Builder) dealWithBackspace(tray *Tr.Tray[rune]) {
	// Backspace : Replace the previous character with the next one in the sequence.
	// However, if there is no previous or next character, do nothing.

	// 1. Remove starting and trailing backspaces
	removeStartingChar(tray, '\b')
	removeTrailingChar(tray, '\b')

	tray.ArrowStart()
	if !tray.MoveRight(1) {
		return
	}

	for {
		c, err := tray.Read()
		if err != nil {
			return
		}

		if c == '\b' {
			if !tray.MoveLeft(1) {
				return
			}

			tray.ShiftRightOfArrow(1)
			tray.ShiftRightOfArrow(1)
		} else {
			if !tray.MoveRight(1) {
				return
			}
		}
	}
}

func (b *Builder) dealWithLineFeed() [][]rune {
	// Line feed : Move to the left edge and down
	result := make([][]rune, 0)

	currentLine := make([]rune, 0)

	for _, char := range b.allChars {
		if char != '\n' {
			currentLine = append(currentLine, char)
		} else {
			result = append(result, currentLine)
			currentLine = make([]rune, 0)
		}
	}

	if len(currentLine) > 0 {
		result = append(result, currentLine)
	}

	return result
}

func (b *Builder) Build(char rune) string {
	// Backspace : Replace the previous character with the next one in the sequence.
	// However, if there is no previous or next character, do nothing.

	var builder strings.Builder

	for i := 0; i < len(b.allChars); {
		switch char {
		case '\u0000':
			// null : Ignore this character
		case '\a':
			// Bell : Ignore this character
		case '\t':
			// Tab : Add spaces until the next tab stop
		case '\n':
			// line feed : Add a new line or move to the left edge and down
		case '\v':
			// vertical tab : Add vertical tabulation
		case '\f':
			// form feed : Go to the next page
		case '\r':
			// carriage return : Move to the start of the line (alone)
			// or move to the start of the line and down (with line feed)
		case '\u001A':
			// Control-Z : End of file for Windows text-mode file i/o
		case '\u001B':
			// escape : Introduce an escape sequence (next character)
		default:
		}
	}

	return builder.String()
}
*/
