package runes

import (
	"unicode/utf8"
)

// CharStream is an interface for a character stream.
type CharStream interface {
	// IsDone checks whether the stream is done.
	//
	// Returns:
	//   - bool: True if the stream is done. False otherwise.
	IsDone() bool

	// Next returns the next character in the stream while advancing the position.
	//
	// Returns:
	//   - rune: The next character in the stream. utf8.RuneError if the stream is done.
	//   - bool: True if the stream has more characters. False otherwise.
	Next() (rune, bool)

	// Peek returns the next character in the stream without advancing the position.
	//
	// Returns:
	//   - rune: The next character in the stream. utf8.RuneError if the stream is done.
	//   - bool: True if the stream has more characters. False otherwise.
	Peek() (rune, bool)

	// Refuse undoes the last Next operation.
	//
	// Returns:
	//   - bool: True if the last Next operation was undone. False otherwise.
	Refuse() bool

	// RefuseMany undoes any Next operation since the last Accept operation.
	// This is useful when you want to abandon the current word and start a new one.
	RefuseMany()

	// Accept accepts the next character in the stream. This is useful for signifying
	// valid sequences that should not be undone.
	Accept()
}

// IsDone implements the CharStream interface.
func (s *Stream) IsDone() bool {
	return s.pos >= len(s.chars)
}

// Next implements the CharStream interface.
func (s *Stream) Next() (rune, bool) {
	if s.pos >= len(s.chars) {
		return utf8.RuneError, false
	}

	r := s.chars[s.pos]
	s.pos++

	return r, true
}

// Peek implements the CharStream interface.
func (s *Stream) Peek() (rune, bool) {
	if s.pos >= len(s.chars) {
		return utf8.RuneError, false
	}

	return s.chars[s.pos], true
}

// Refuse implements the CharStream interface.
func (s *Stream) Refuse() bool {
	if s.pos == 0 {
		return false
	}

	s.pos--
	return true
}

// RefuseMany implements the CharStream interface.
func (s *Stream) RefuseMany() {
	s.pos = s.last_accept
}

// Accept implements the CharStream interface.
func (s *Stream) Accept() {
	s.last_accept = s.pos
}

// Stream is a character stream.
type Stream struct {
	// chars is the character stream.
	chars []rune

	// pos is the current position in the stream.
	pos int

	// last_accept is the last accepted position in the stream.
	last_accept int
}

// NewStream creates a new stream with the given runes.
//
// Parameters:
//   - b: The rune slice of the stream.
//
// Returns:
//   - *Stream: The new stream. Never nil.
func NewStream(b []rune) *Stream {
	return &Stream{
		chars:       b,
		pos:         0,
		last_accept: 0,
	}
}
