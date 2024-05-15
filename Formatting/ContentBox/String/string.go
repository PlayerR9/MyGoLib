package String

import (
	"fmt"
	"strings"
	"unicode/utf8"

	cdd "github.com/PlayerR9/MyGoLib/ComplexData/Display"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
	"github.com/gdamore/tcell"

	com "github.com/PlayerR9/MyGoLib/Units/Common"
)

// String represents a unit of data in a draw table with a specific style.
type String struct {
	// content is the content of the unit.
	content string

	// style is the style of the unit.
	style tcell.Style

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
		style:   dtu.style,
		length:  dtu.length,
	}
}

// Draw is a method of cdd.TableDrawer that draws the unit to the table at the given x and y
// coordinates.
//
// Parameters:
//   - table: The table to draw the unit to.
//   - x: The x coordinate to draw the unit at.
//   - y: The y coordinate to draw the unit at.
//
// Returns:
//   - error: An error of type *ers.ErrInvalidParameter if the table is nil.
//
// Behaviors:
//   - Out of bounds values are ignored.
func (dtu *String) Draw(table *cdd.DrawTable, x, y int) error {
	if table == nil {
		return ers.NewErrNilParameter("table")
	}

	if err := table.IsXInBounds(x); err != nil {
		return ers.NewErrInvalidParameter("x", err)
	} else if err := table.IsYInBounds(y); err != nil {
		return ers.NewErrInvalidParameter("y", err)
	}

	exceeding := x + dtu.length - table.GetWidth()

	if exceeding > 0 {
		return fmt.Errorf("content exceeds table width by %d", exceeding)
	}

	table.WriteLineAt(x, y, dtu.content, dtu.style, true)

	return nil
}

// Draw is a method of cdd.TableDrawer that draws the unit to the table at the given x and y
// coordinates.
//
// Parameters:
//   - table: The table to draw the unit to.
//   - x: The x coordinate to draw the unit at.
//   - y: The y coordinate to draw the unit at.
//
// Returns:
//   - error: An error of type *ers.ErrInvalidParameter if the table is nil.
//
// Behaviors:
//   - Out of bounds values are ignored.
func (dtu *String) ForceDraw(table *cdd.DrawTable, x, y int) error {
	if table == nil {
		return ers.NewErrNilParameter("table")
	}

	table.WriteLineAt(x, y, dtu.content, dtu.style, true)

	return nil
}

// NewDtUnit creates a new DtUnit with the given content and style.
//
// Parameters:
//   - content: The content of the unit.
//   - style: The style of the unit.
//
// Returns:
//   - *DtUnit: The new DtUnit.
func NewString(content string, style tcell.Style) *String {
	return &String{
		content: content,
		style:   style,
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

// GetStyle returns the style of the unit.
//
// Returns:
//   - tcell.Style: The style of the unit.
func (dtu *String) GetStyle() tcell.Style {
	return dtu.style
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
			fields = append(fields, NewString(builder.String(), dtu.style))
			builder.Reset()
		} else {
			builder.WriteRune(r)
		}
	}

	if builder.Len() > 0 {
		fields = append(fields, NewString(builder.String(), dtu.style))
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
		return &String{
			content: s.content,
			style:   s.style,
			length:  s.length,
		}
	}

	if limit <= 0 {
		return &String{
			content: "",
			style:   s.style,
			length:  0,
		}
	}

	return &String{
		content: s.content[:limit],
		style:   s.style,
		length:  limit,
	}
}
