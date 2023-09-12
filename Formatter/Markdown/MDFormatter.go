package Markdown

import "fmt"

const (
	STYLE_CENTER = iota
	STYLE_LEFT
	STYLE_RIGHT
)

const (
	ALIGN_LEFT   = ":---"
	ALIGN_CENTER = ":---:"
	ALIGN_RIGHT  = "---:"

	PIPE = " | "

	NEW_LINE = "\n"
)

var Pedantic bool = false // If false, the formatter will try to fix errors in the table.

type Header struct {
	text  string
	style int
}

func NewHeader(text string, style int) Header {
	if style < STYLE_CENTER || style > STYLE_RIGHT {
		panic(fmt.Sprintf("invalid style: %d", style))
	}

	return Header{text, style}
}

type Table struct {
	headers []Header
	rows    [][]string
}

func NewTable() Table {
	return Table{
		headers: make([]Header, 0),
		rows:    make([][]string, 0),
	}
}

func NewTableWithHeaders(headers ...Header) Table {
	return Table{
		headers: headers,
		rows:    make([][]string, 0),
	}
}

func (t *Table) AppendHeader(header Header) {
	t.headers = append(t.headers, header)

	if !Pedantic {
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

func (t *Table) AddRow(elements ...string) error {
	if len(elements) != len(t.headers) {
		return fmt.Errorf("number of elements does not match number of headers")
	}

	t.rows = append(t.rows, elements)

	return nil
}

func (t Table) ToText() []string {
	text := make([]string, 0)

	// Headers
	text = append(text, PIPE)

	for _, header := range t.headers {
		text = append(text, header.text, PIPE)
	}

	text = append(text, NEW_LINE)

	// Alignments
	text = append(text, PIPE)

	for _, header := range t.headers {
		switch header.style {
		case STYLE_CENTER:
			text = append(text, ALIGN_CENTER)
		case STYLE_LEFT:
			text = append(text, ALIGN_LEFT)
		case STYLE_RIGHT:
			text = append(text, ALIGN_RIGHT)
		}

		text = append(text, PIPE)
	}

	text = append(text, NEW_LINE)

	// Rows
	for _, row := range t.rows {
		text = append(text, PIPE)

		for _, element := range row {
			text = append(text, element, PIPE)
		}

		text = append(text, NEW_LINE)
	}

	return text
}
