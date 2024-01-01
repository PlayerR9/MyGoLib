package FScreen

import (
	"strings"

	"github.com/gdamore/tcell"
)

const (
	IndentLevel int  = 2
	Separator   rune = '='
	Space       rune = ' '
)

type MessageBox struct {
	table          [][]rune
	styles         []tcell.Style
	firstEmptyLine int

	width, height int

	isShifting bool
}

// WARNING: This function doesn't do any checks
func (mb *MessageBox) DrawLine(x, y int) {
	if y >= mb.firstEmptyLine {
		return
	}

	for i, cell := range mb.table[y] {
		screen.SetContent(x+i, y, cell, nil, mb.styles[y])
	}
}

// Returns the first empty index after the last character
// WARNING: This function doesn't do any checks
func (mb *MessageBox) WriteAt(x, y int, word string) int {
	row := mb.table[y]

	for _, char := range word {
		row[x] = char
		x++
	}

	return x
}

// WARNING: This function doesn't do any checks
func (mb *MessageBox) CanWriteAt(x int, word string) bool {
	return x+len(word) <= mb.width
}

// Returns the first empty index after the last character
// WARNING: This function doesn't do any checks
func (mb *MessageBox) WriteFieldsAt(x, y int, fields []string) (int, int) {
	if len(fields) == 1 {
		if mb.CanWriteAt(x, fields[0]) {
			return mb.WriteAt(x, y, fields[0]), y
		}

		x = mb.WriteAt(x, y, fields[0][:mb.width-x-HellipLen])
		return mb.WriteAt(x, y, Hellip), y
	} else if !mb.CanWriteAt(x, fields[0]) {
		x = mb.WriteAt(x, y, fields[0][:mb.width-x-HellipLen])
		x = mb.WriteAt(x, y, Hellip)

		return IndentLevel, y + 1
	}

	x = mb.WriteAt(x, y, fields[0])

	index := -1
	for i, field := range fields[1:] {
		if !mb.CanWriteAt(x+2, field) {
			index = i
			break
		}

		x = mb.WriteAt(x+2, y, field)
	}

	if index != -1 {
		mb.WriteAt(x, y, Hellip)

		return IndentLevel, y + 1
	}

	return x, y
}

// -1 = no fields, >= 0 index of the last field that fits
// WARNING: This function doesn't do any checks
func (mb *MessageBox) CanWriteFieldsAt(x int, fields []string) int {
	if !mb.CanWriteAt(x, fields[0]) {
		return -1
	}

	x += len(fields[0])

	for i, field := range fields[1:] {
		if !mb.CanWriteAt(x+2, field) {
			return i
		}

		x += len(field) + 2
	}

	return len(fields) - 1
}

// WARNING: This function doesn't do any checks
func (mb *MessageBox) EnqueueContents(contents []string, style tcell.Style, isIndented bool) []string {
	fields := strings.Fields(contents[0])
	if len(fields) == 0 {
		return nil
	}

	x := 0
	if isIndented {
		x = IndentLevel
	}

	var y int
	x, y = mb.WriteFieldsAt(x, mb.firstEmptyLine, fields)
	mb.styles[y] = style
	if y != mb.firstEmptyLine {
		mb.firstEmptyLine = y
		return contents[1:]
	}

	// Try to see if the next content fits on the same line
	for index, content := range contents[1:] {
		fields = strings.Fields(content)
		if len(fields) == 0 {
			continue // Skip empty lines
		}

		lastValidField := mb.CanWriteFieldsAt(x, fields)
		if lastValidField == -1 {
			mb.firstEmptyLine++
			return contents[index:]
		}

		for _, field := range fields[:lastValidField+1] {
			x = mb.WriteAt(x+2, y, field)
		}

		if lastValidField != len(fields)-1 {
			mb.firstEmptyLine++
			return append([]string{strings.Join(fields[lastValidField+1:], " ")}, contents[index+1:]...)
		}
	}

	return nil
}

// WARNING: This function doesn't do any checks
func (mb *MessageBox) EnqueueSeparator() {
	for i := 0; i < mb.width; i++ {
		mb.table[mb.firstEmptyLine][i] = Separator
	}

	mb.firstEmptyLine++
}

// WARNING: This function doesn't do any checks
func (mb *MessageBox) EnqueueLineBreak() {
	for i := 0; i < mb.width; i++ {
		mb.table[mb.firstEmptyLine][i] = Space
	}

	mb.firstEmptyLine++
}

func (mb *MessageBox) CanShiftUp() bool {
	if mb.isShifting {
		if mb.firstEmptyLine == 0 {
			mb.isShifting = false
			return false
		}
	} else if mb.firstEmptyLine < mb.height/2 {
		return false
	}

	return true
}

func (mb *MessageBox) ShiftUp() {
	size := 1

	for size < mb.firstEmptyLine && mb.table[mb.firstEmptyLine-size][0] == Space {
		size++
	}

	for i := 0; i < mb.height-size; i++ {
		mb.table[i] = mb.table[i+size]
		mb.styles[i] = mb.styles[i+size]
	}

	mb.firstEmptyLine -= size
	mb.isShifting = true
}
