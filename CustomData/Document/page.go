package Document

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
