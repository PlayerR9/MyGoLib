package String

import (
	"strings"
	"unicode/utf8"

	com "github.com/PlayerR9/MyGoLib/Units/Common"
)

// String represents a unit of data in a draw table with a specific style.
type String struct {
	// content is the content of the unit.
	content string

	// length is the length of the content.
	length int
}

// String returns the content of the unit as a string.
//
// Returns:
//   - string: The content of the unit.
func (dtu *String) String() string {
	return dtu.content
}

// Copy returns a copy of the unit.
//
// Returns:
//   - com.Copier: A copy of the unit.
func (dtu *String) Copy() com.Copier {
	return &String{
		content: dtu.content,
		length:  dtu.length,
	}
}

// GetRunes returns the content of the unit as a slice of runes.
//
// Returns:
//   - [][]rune: The content of the unit as a slice of runes.
//
// Behaviors:
//   - The content is returned as a slice of runes with one line.
func (s *String) GetRunes() [][]rune {
	runes := make([][]rune, 1)
	runes[0] = []rune(s.content)

	return runes

}

// NewDtUnit creates a new DtUnit with the given content and style.
//
// Parameters:
//   - content: The content of the unit.
//
// Returns:
//   - *DtUnit: The new DtUnit.
func NewString(content string) *String {
	return &String{
		content: content,
		length:  utf8.RuneCountInString(content),
	}
}

// GetContent returns the content of the unit.
//
// Returns:
//   - string: The content of the unit.
func (dtu *String) GetContent() string {
	return dtu.content
}

// GetLength returns the length of the content.
//
// Returns:
//   - int: The length of the content.
func (dtu *String) GetLength() int {
	return dtu.length
}

// Fields splits a string into fields.
//
// Returns:
//   - []*String: The fields of the string.
func (dtu *String) Fields() []*String {
	fields := make([]*String, 0)
	var builder strings.Builder

	for _, r := range dtu.content {
		if r == ' ' {
			fields = append(fields, NewString(builder.String()))
			builder.Reset()
		} else {
			builder.WriteRune(r)
		}
	}

	if builder.Len() > 0 {
		fields = append(fields, NewString(builder.String()))
	}

	return fields
}

// AppendRune appends a rune to the content of the unit.
//
// Parameters:
//   - r: The rune to append.
func (dtu *String) AppendRune(r rune) {
	dtu.content += string(r)
	dtu.length++
}

// AppendString appends a string to the content of the unit.
//
// Parameters:
//   - s: The string to append.
func (dtu *String) AppendString(s string) {
	if s == "" {
		return
	}

	dtu.content += s
	dtu.length += utf8.RuneCountInString(s)
}

// PrependRune prepends a rune to the content of the unit.
//
// Parameters:
//   - r: The rune to prepend.
func (dtu *String) PrependRune(r rune) {
	dtu.content = string(r) + dtu.content
	dtu.length++
}

// PrependString prepends a string to the content of the unit.
//
// Parameters:
//   - s: The string to prepend.
func (dtu *String) PrependString(s string) {
	if s == "" {
		return
	}

	dtu.content = s + dtu.content
	dtu.length += utf8.RuneCountInString(s)
}

// ReplaceSuffix replaces the end of the string with the given suffix.
// It fails if the suffix is longer than the string.
//
// Parameters:
//   - suffix: The suffix to replace the end of the string.
//
// Returns:
//   - bool: True if the suffix was replaced, and false otherwise.
func (s *String) ReplaceSuffix(suffix string) bool {
	if suffix == "" {
		return true
	}

	countSuffix := utf8.RuneCountInString(suffix)

	if s.length < countSuffix {
		return false
	}

	if s.length == countSuffix {
		s.content = suffix
	} else {
		s.content = s.content[:s.length-countSuffix] + suffix
	}

	return true
}

// TrimEnd trims the end of the string to the given limit.
//
// Parameters:
//   - limit: The limit to trim the string to.
//
// Returns:
//   - *String: The trimmed string.
//
// Behaviors:
//   - If the limit is less than or equal to 0, the string is set to an empty string.
//   - If the limit is greater than the length of the string, the string is unchanged.
func (s *String) TrimEnd(limit int) *String {
	if s.content == "" || limit >= s.length {
		return s.Copy().(*String)
	}

	if limit <= 0 {
		return &String{
			content: "",
			length:  0,
		}
	}

	return &String{
		content: s.content[:limit],
		length:  limit,
	}
}
