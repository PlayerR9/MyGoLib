package Markdown

import (
	"fmt"
	"log"
	"os"
)

const (
	// Styles
	STYLE_CENTER = iota
	STYLE_LEFT
	STYLE_RIGHT
)

var (
	// Tokens
	tok_align_center string = ":---:"
	tok_align_left   string = ":---"
	tok_align_right  string = "---:"
	tok_pipe         string = " | "
	tok_new_line     string = "\n"

	// Global variables
	pedantic  bool        = true                                                // Pedantic mode is enabled by default
	DebugMode bool        = false                                               // Debug mode is disabled by default
	debugger  *log.Logger = log.New(os.Stdout, "[MDFormatter] ", log.LstdFlags) // Debugger
)

// GetPedanticMode returns the current pedantic mode. When pedantic mode is enabled, the formatter will panic when rows and columns do not match,
// otherwise it will fill the missing elements with empty strings. **Note:** Pedantic mode is enabled by default.
//
// Returns:
//   - bool: The current pedantic mode.
func GetPedanticMode() bool {
	return pedantic
}

// SetPedanticMode sets the pedantic mode. If pedantic mode was already enabled, if the debug mode is enabled, it will print a warning.
func SetPedantic() {
	if pedantic {
		if DebugMode {
			debugger.Printf("pedantic mode is already enabled")
		}
	} else {
		pedantic = true

		if DebugMode {
			debugger.Printf("pedantic mode enabled")
		}
	}
}

// Header is a struct that contains the text and the style of a header. The style can be one of the following: STYLE_CENTER, STYLE_LEFT or STYLE_RIGHT.
type Header struct {
	text  string // Text of the header
	style int    // Style of the header
}

// NewHeader creates a new header with the given text and style. If the style is not one of the following: STYLE_CENTER, STYLE_LEFT or STYLE_RIGHT, it will panic.
//
// Parameters:
//   - text: Text of the header.
//   - style: Style of the header.
//
// Returns:
//   - Header: The new header.
func NewHeader(text string, style int) Header {
	if style < STYLE_CENTER || style > STYLE_RIGHT {
		if DebugMode {
			debugger.Panicf("invalid style: %d", style)
		} else {
			panic(fmt.Sprintf("invalid style: %d", style))
		}
	}

	return Header{text, style}
}

// Table is a struct that contains the headers and the rows of a table.
type Table struct {
	headers []Header   // Headers of the table
	rows    [][]string // Rows of the table
}

// NewTable creates a new table with no headers and no rows. Useful when you want to create a table and add the headers and rows later.
// This clashes with the pedantic mode, so if you want to use this function, you should disable the pedantic mode. For creating a table with headers,
// use NewTableWithHeaders instead.
//
// Returns:
//   - Table: The new table.
func NewTable() Table {
	return Table{
		headers: make([]Header, 0),
		rows:    make([][]string, 0),
	}
}

// NewTableWithHeaders creates a new table with the given headers and no rows. Useful when you want to create a table and add the rows later. This
// function does not clash with the pedantic mode, so you can use it even if the pedantic mode is enabled.
//
// Parameters:
//   - headers: Headers of the table.
//
// Returns:
//   - Table: The new table.
func NewTableWithHeaders(headers ...Header) Table {
	return Table{
		headers: headers,
		rows:    make([][]string, 0),
	}
}

// AppendHeader appends a header to the table. If the pedantic mode is enabled, it will panic if a row has less elements than the number of headers.
// Otherwise, it will fill the missing elements with empty strings.
//
// Parameters:
//   - header: Header to append.
func (t *Table) AppendHeader(header Header) {
	t.headers = append(t.headers, header)

	if !pedantic {
		for i, row := range t.rows {
			if len(row) < len(t.headers) {
				panic(fmt.Sprintf("row %d has less elements than the number of headers", i))
			}
		}
	} else {
		for _, row := range t.rows {
			if len(row) < len(t.headers) {
				for len(row) < len(t.headers) {
					row = append(row, "")
				}
			}
		}
	}
}

func (t *Table) AddRow(elements []string) error {
	if len(elements) != len(t.headers) {
		return fmt.Errorf("number of elements does not match number of headers")
	}

	t.rows = append(t.rows, elements)

	return nil
}

func (t Table) ToText() []string {
	text := make([]string, 0)

	// Headers
	text = append(text, tok_pipe)

	for _, header := range t.headers {
		text = append(text, header.text, tok_pipe)
	}

	text = append(text, tok_new_line)

	// Alignments
	text = append(text, tok_pipe)

	for _, header := range t.headers {
		switch header.style {
		case STYLE_CENTER:
			text = append(text, tok_align_center)
		case STYLE_LEFT:
			text = append(text, tok_align_left)
		case STYLE_RIGHT:
			text = append(text, tok_align_right)
		}

		text = append(text, tok_pipe)
	}

	text = append(text, tok_new_line)

	// Rows
	for _, row := range t.rows {
		text = append(text, tok_pipe)

		for _, element := range row {
			text = append(text, element, tok_pipe)
		}

		text = append(text, tok_new_line)
	}

	return text
}
