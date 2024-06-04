package FString

import (
	"strings"
	"sync"
	"unicode/utf8"

	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

/////////////////////////////////////////////////

// checkString is a private function that checks a string for invalid runes.
//
// Parameters:
//   - str: The string to check.
//
// Returns:
//   - []rune: The runes of the string.
//   - error: An error of type *Errors.ErrAt if an invalid rune is found in the string.
func checkString(str string) ([]rune, error) {
	var runes []rune

	for j := 0; len(str) > 0; j++ {
		char, size := utf8.DecodeRuneInString(str)
		str = str[size:]

		if char == utf8.RuneError {
			return nil, ue.NewErrAt(j, "character", ue.NewErrInvalidRune(nil))
		}

		if char == '\r' && size != 0 {
			nextRune, size := utf8.DecodeRuneInString(str)

			if nextRune == '\n' {
				str = str[size:]
			}
		}

		runes = append(runes, char)
	}

	return runes, nil
}

// sectionBuilder is a type that represents a section of a page.
type sectionBuilder struct {
	// buff is the string buff for the section.
	buff strings.Builder

	// lines are the lines in the section.
	lines [][]string

	// lastLine is the last line of the section.
	lastLine int

	// mu is the mutex for the builder.
	mu sync.RWMutex
}

// Cleanup implements the Cleanup interface method.
func (sb *sectionBuilder) Clean() {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	for i := 0; i < len(sb.lines); i++ {
		sb.lines[i] = nil
	}

	sb.lines = nil

	sb.buff.Reset()
}

// newSectionBuilder creates a new section.
//
// Returns:
//   - *Section: The new section.
func newSectionBuilder() *sectionBuilder {
	return &sectionBuilder{
		lines:    [][]string{{}},
		lastLine: 0,
	}
}

// removeOne is a function that removes the last character from the section.
//
// Returns:
//   - bool: True if a character was removed. False otherwise.
func (sb *sectionBuilder) removeOne() bool {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	if sb.buff.Len() > 0 {
		str := sb.buff.String()
		str = str[:len(str)-1]

		sb.buff.Reset()
		sb.buff.WriteString(str)

		return true
	}

	for i := sb.lastLine; i >= 0; i-- {
		words := sb.lines[i]

		for j := len(words) - 1; j >= 0; j-- {
			if len(words[j]) > 0 {
				words[j] = words[j][:len(words[j])-1]
				return true
			}
		}
	}

	return false
}

// getLines is a function that returns the words of the section.
//
// Returns:
//   - [][]string: The words of the section.
func (sb *sectionBuilder) getLines() [][]string {
	sb.mu.RLock()
	defer sb.mu.RUnlock()

	return sb.lines
}

// isFirstOfLine is a function that returns true if the current position is the first
// position of a line.
//
// Returns:
//   - bool: True if the current position is the first position of a line.
func (sb *sectionBuilder) isFirstOfLine() bool {
	sb.mu.RLock()
	defer sb.mu.RUnlock()

	return sb.buff.Len() == 0 && len(sb.lines[sb.lastLine]) == 0
}

// accept is a function that accepts the current word and
// creates a new line.
func (sb *sectionBuilder) accept(delim string) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	if sb.buff.Len() > 0 {
		if delim != "" {
			sb.buff.WriteString(delim)
		}

		sb.lines[sb.lastLine] = append(sb.lines[sb.lastLine], sb.buff.String())
		sb.buff.Reset()
	}

	sb.lines = append(sb.lines, []string{})
	sb.lastLine++
}

// accept is a function that accepts the current word and
// creates a new line.
func (sb *sectionBuilder) mayAccept(delim string) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	if sb.buff.Len() == 0 {
		return
	}

	if delim != "" {
		sb.buff.WriteString(delim)
	}

	sb.lines[sb.lastLine] = append(sb.lines[sb.lastLine], sb.buff.String())
	sb.buff.Reset()

	sb.lines = append(sb.lines, []string{})
	sb.lastLine++
}

// acceptWord is a function that accepts the current in-progress word
// and resets the builder.
func (sb *sectionBuilder) acceptWord() {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	if sb.buff.Len() == 0 {
		return
	}

	sb.lines[sb.lastLine] = append(sb.lines[sb.lastLine], sb.buff.String())
	sb.buff.Reset()
}

// writeRune adds a rune to the current, in-progress word.
//
// Parameters:
//   - r: The rune to write.
func (sb *sectionBuilder) writeRune(r rune) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	sb.buff.WriteRune(r)
}

// writeString adds a string to the current, in-progress word.
//
// Parameters:
//   - str: The string to write.
func (sb *sectionBuilder) writeString(str string) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	sb.buff.WriteString(str)
}

// buffer is a type that represents a buffer of a document.
type buffer struct {
	// pages are the pages of the buffer.
	pages [][]*sectionBuilder

	// buff is the in-progress section of the buffer.
	buff *sectionBuilder

	// lastPage is the last page of the buffer.
	lastPage int
}

// Cleanup implements the Cleanup interface method.
func (b *buffer) Clean() {
	// pages are the pages of the buffer.
	for i := 0; i < len(b.pages); i++ {
		for j := 0; j < len(b.pages[i]); j++ {
			b.pages[i][j].Clean()
			b.pages[i][j] = nil
		}

		b.pages[i] = nil
	}

	b.pages = nil

	b.buff.Clean()
	b.buff = nil
}

// newBuffer creates a new buffer.
//
// Returns:
//   - *buffer: The new buffer.
func newBuffer() *buffer {
	return &buffer{
		pages:    [][]*sectionBuilder{{}},
		buff:     nil,
		lastPage: 0,
	}
}

// getPages returns the pages of the StdPrinter.
//
// Returns:
//   - [][][][]string: The pages of the StdPrinter.
func (b *buffer) getPages() [][][][]string {
	b.finalize()

	pages := b.pages

	allStrings := make([][][][]string, 0, len(pages))

	for _, page := range pages {
		sectionLines := make([][][]string, 0)

		for _, section := range page {
			sectionLines = append(sectionLines, section.getLines())
		}

		allStrings = append(allStrings, sectionLines)
	}

	return allStrings
}

// isFirstOfLine is a private function that returns true if the current position is the first
// position of a line.
//
// Returns:
//   - bool: True if the current position is the first position of a line.
func (b *buffer) isFirstOfLine() bool {
	return b.buff == nil || b.buff.isFirstOfLine()
}

// forceWriteString is a private function that writes the indentation to the formatted string.
func (b *buffer) forceWriteString(str string) {
	if str == "" {
		return
	}

	if b.buff == nil {
		b.buff = newSectionBuilder()
	}

	b.buff.writeString(str)
}

// Accept is a function that accepts the current in-progress buffer
// by converting it to the specified section type. Lastly, the section
// is added to the page.
//
// Parameters:
//   - sectionType: The section type to convert the buffer to.
//
// Behaviors:
//   - Even when the buffer is empty, the section is still added to the page.
//     To avoid this, use the Finalize function.
func (b *buffer) accept() {
	if b.buff != nil {
		b.buff.acceptWord()
	}

	b.pages[b.lastPage] = append(b.pages[b.lastPage], b.buff)

	b.buff = nil
}

// write is a private function that appends a rune to the buffer
// while dealing with special characters.
//
// Parameters:
//   - char: The rune to append.
func (b *buffer) write(char rune) {
	switch char {
	case '\t':
		// Tab : Add spaces until the next tab stop
		if b.buff == nil {
			b.buff = newSectionBuilder()
		} else {
			b.buff.acceptWord()
		}

		b.buff.writeRune(char)

		b.buff.acceptWord()
	case '\v':
		// vertical tab : Add vertical tabulation

		// Do nothing
	case '\r', '\n', '\u0085':
		// carriage return : Move to the start of the line (alone)
		// or move to the start of the line and down (with line feed)
		// line feed : Add a new line or move to the left edge and down

		b.accept()
	case '\f':
		// form feed : Go to the next page
		b.accept()

		b.lastPage++
		b.pages = append(b.pages, []*sectionBuilder{})
	case ' ':
		// Space
		if b.buff != nil {
			b.buff.acceptWord()
		}
	case '\u0000', '\a':
		// null : Ignore this character
		// Bell : Ignore this character
	case '\b':
		// backspace : Remove the last character
		if b.buff != nil {
			ok := b.buff.removeOne()
			if ok {
				return
			}
		}

		for i := b.lastPage; i >= 0; i-- {
			sections := b.pages[i]

			for j := len(sections) - 1; j >= 0; j-- {
				section := sections[j]

				ok := section.removeOne()
				if ok {
					return
				}
			}
		}
	case '\u001A':
		// Control-Z : End of file for Windows text-mode file i/o
		b.finalize()
	case '\u001B':
		// escape : Introduce an escape sequence (next character)
		// Do nothing
	default:
		// NBSP : Non-breaking space
		// any other normal character
		if b.buff == nil {
			b.buff = newSectionBuilder()
		}

		if char == NBSP {
			// Non-breaking space
			b.buff.writeRune(' ')
		} else {
			b.buff.writeRune(char)
		}
	}
}

// writeRune is a private function that appends a rune to the buffer
// without checking for special characters.
//
// Parameters:
//   - r: The rune to append.
func (b *buffer) writeRune(r rune) {
	if b.buff == nil {
		b.buff = newSectionBuilder()
	}

	b.buff.writeRune(r)
}

// writeRunes is a private function that writes runes to the formatted string.
//
// Parameters:
//   - runes: The runes to write.
func (b *buffer) writeRunes(runes []rune) {
	for _, r := range runes {
		b.write(r)
	}
}

// writeBytes is a private function that writes bytes to the formatted string.
//
// Parameters:
//   - b: The bytes to write.
//
// Returns:
//   - int: The number of bytes written.
func (b *buffer) writeBytes(data []byte) int {
	if len(data) == 0 {
		return 0
	}

	var count int

	for count = 0; len(data) > 0; count++ {
		r, size := utf8.DecodeRune(data)

		b.write(r)

		data = data[size:]
	}

	return count
}

// acceptWord is a private function that accepts the current word of the formatted string.
func (b *buffer) acceptWord() {
	if b.buff != nil {
		b.buff.acceptWord()
	}
}

// acceptLine is a private function that accepts the current line of the formatted string.
func (b *buffer) acceptLine(delim string) {
	if b.buff != nil {
		b.buff.mayAccept(delim)
	}
}

// writeEmptyLine is a private function that accepts the current line
// regardless of the whether the line is empty or not.
func (b *buffer) writeEmptyLine(delim string) {
	if b.buff == nil {
		b.buff = newSectionBuilder()
	}

	b.buff.accept(delim)
}

// finalize is a private function that finalizes the buffer.
func (b *buffer) finalize() {
	if b.buff == nil {
		return
	}

	b.buff.acceptWord()

	b.pages[b.lastPage] = append(b.pages[b.lastPage], b.buff)

	b.buff = nil
}
