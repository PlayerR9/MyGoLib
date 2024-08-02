package String

import (
	"strings"
	"unicode/utf8"
)

// String represents a unit of data in a draw table with a specific style.
type String struct {
	// content is the content of the unit.
	content string

	// length is the length of the content.
	length int
}

// Equals implements common.Objecter.
func (s *String) Equals(other *String) bool {
	if other == nil {
		return false
	}

	return s.content == other.content
}

// String returns the content of the unit as a string.
//
// Returns:
//   - string: The content of the unit.
func (s *String) String() string {
	return s.content
}

// Copy returns a copy of the unit.
//
// Returns:
//   - *String: A copy of the unit.
func (s *String) Copy() *String {
	return &String{
		content: s.content,
		length:  s.length,
	}
}

// Runes returns the content of the unit as a slice of runes.
//
// Parameters:
//   - width: The width of the table.
//   - height: The height of the table.
//
// Returns:
//   - [][]rune: The content of the unit as a slice of runes.
//   - error: An error if the content could not be converted to runes.
//
// Behaviors:
//   - Always assume that the width and height are greater than 0. No need to check for
//     this.
//   - Never errors.
func (s *String) Runes(width, height int) ([][]rune, error) {
	return [][]rune{[]rune(s.content)}, nil
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
func (s *String) GetContent() string {
	return s.content
}

// GetLength returns the length of the content.
//
// Returns:
//   - int: The length of the content.
func (s *String) GetLength() int {
	return s.length
}

// Fields splits a string into fields.
//
// Returns:
//   - []*String: The fields of the string.
func (s *String) Fields() []*String {
	fields := make([]*String, 0)
	var builder strings.Builder

	for _, r := range s.content {
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
func (s *String) AppendRune(r rune) {
	s.content += string(r)
	s.length++
}

// AppendString appends a string to the content of the unit.
//
// Parameters:
//   - otherS: The string to append.
func (s *String) AppendString(otherS string) {
	if otherS == "" {
		return
	}

	s.content += otherS
	s.length += utf8.RuneCountInString(otherS)
}

// PrependRune prepends a rune to the content of the unit.
//
// Parameters:
//   - r: The rune to prepend.
func (s *String) PrependRune(r rune) {
	s.content = string(r) + s.content
	s.length++
}

// PrependString prepends a string to the content of the unit.
//
// Parameters:
//   - otherS: The string to prepend.
func (s *String) PrependString(otherS string) {
	if otherS == "" {
		return
	}

	s.content = otherS + s.content
	s.length += utf8.RuneCountInString(otherS)
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
		return s.Copy()
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
