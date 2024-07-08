package pkg

import (
	"strings"
	"sync"

	"github.com/PlayerR9/MyGoLib/Utility/object"
)

// sectionBuilder is a type that represents a section of a page.
type sectionBuilder struct {
	// buff is the string buff for the section.
	buff strings.Builder

	// lines are the lines in the section.
	lines [][]string

	// last_line is the last line of the section.
	last_line int

	// mu is the mutex for the builder.
	mu sync.Mutex
}

// Cleanup implements the object.Cleaner interface method.
func (sb *sectionBuilder) Cleanup() {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	for i := 0; i < len(sb.lines); i++ {
		line := object.CleanSliceOf(sb.lines[i])
		sb.lines[i] = line
		sb.lines[i] = nil
	}

	sb.lines = sb.lines[:0]

	sb.buff.Reset()
}

// new_section_builder creates a new section builder.
//
// Returns:
//   - *sectionBuilder: The new section builder.
func new_section_builder() *sectionBuilder {
	sb := &sectionBuilder{
		lines:     [][]string{{}},
		last_line: 0,
	}

	return sb
}

// remove_one is a function that removes the last character from the section.
//
// Returns:
//   - bool: True if a character was removed. False otherwise.
func (sb *sectionBuilder) remove_one() bool {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	size := sb.buff.Len()

	if size > 0 {
		str := sb.buff.String()
		str = str[:len(str)-1]

		sb.buff.Reset()
		sb.buff.WriteString(str)

		return true
	}

	for i := sb.last_line; i >= 0; i-- {
		words := sb.lines[i]

		for j := len(words) - 1; j >= 0; j-- {
			word_size := len(words[j])

			if word_size > 0 {
				words[j] = words[j][:word_size-1]
				return true
			}
		}
	}

	return false
}

// get_lines is a function that returns the words of the section.
//
// Returns:
//   - [][]string: The words of the section.
func (sb *sectionBuilder) get_lines() [][]string {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	return sb.lines
}

// is_first_of_line is a function that returns true if the current position is the first
// position of a line.
//
// Returns:
//   - bool: True if the current position is the first position of a line.
func (sb *sectionBuilder) is_first_of_line() bool {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	size := sb.buff.Len()
	if size > 0 {
		return false
	}

	size = len(sb.lines[sb.last_line])
	return size == 0
}

// accept is a function that accepts the current word and
// creates a new line.
//
// Parameters:
//   - right_delim: The right delimiter to use. If empty, it is not used.
func (sb *sectionBuilder) accept(right_delim string) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	size := sb.buff.Len()

	if size > 0 {
		if right_delim != "" {
			sb.buff.WriteString(right_delim)
		}

		str := sb.buff.String()
		sb.lines[sb.last_line] = append(sb.lines[sb.last_line], str)
		sb.buff.Reset()
	}

	sb.lines = append(sb.lines, []string{})
	sb.last_line++
}

// may_accept is a function that, like accept, accepts the current word and
// creates a new line. However, it only does so if the buffer is not empty.
//
// Parameters:
//   - right_delim: The delimiter to use. If empty, it is not used.
func (sb *sectionBuilder) may_accept(right_delim string) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	size := sb.buff.Len()
	if size == 0 {
		return
	}

	if right_delim != "" {
		sb.buff.WriteString(right_delim)
	}

	str := sb.buff.String()

	sb.lines[sb.last_line] = append(sb.lines[sb.last_line], str)
	sb.buff.Reset()

	sb.lines = append(sb.lines, []string{})
	sb.last_line++
}

// accept_word is a function that accepts the current in-progress word
// and resets the builder.
//
// Behaviors:
//   - If the buffer is empty, nothing happens.
func (sb *sectionBuilder) accept_word() {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	size := sb.buff.Len()
	if size == 0 {
		return
	}

	str := sb.buff.String()

	sb.lines[sb.last_line] = append(sb.lines[sb.last_line], str)
	sb.buff.Reset()
}

// write_rune adds a rune to the current, in-progress word.
//
// Parameters:
//   - r: The rune to write.
func (sb *sectionBuilder) write_rune(r rune) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	sb.buff.WriteRune(r)
}

// write_string adds a string to the current, in-progress word.
//
// Parameters:
//   - str: The string to write.
func (sb *sectionBuilder) write_string(str string) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	sb.buff.WriteString(str)
}
