package Markdown

import (
	"fmt"
)

const (
	STYLE_CENTER = iota
	STYLE_LEFT
	STYLE_RIGHT
)

var (
	tok_align_center string = ":---:"
	tok_align_left   string = ":---"
	tok_align_right  string = "---:"
	tok_pipe         string = " | "
	tok_new_line     string = "\n"
	Pedantic         bool   = true
)

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

func (t *Table) AddRow(elements []string) error {
	if len(elements) != len(t.headers) {
		return fmt.Errorf("number of elements does not match number of headers")
	}

	t.rows = append(t.rows, elements)

	return nil
}

func (t Table) ToText() []string {
	text := make([]string, 0)
	text = append(text, tok_pipe)

	for _, header := range t.headers {
		text = append(text, header.text, tok_pipe)
	}

	text = append(text, tok_new_line)
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
	for _, row := range t.rows {
		text = append(text, tok_pipe)

		for _, element := range row {
			text = append(text, element, tok_pipe)
		}

		text = append(text, tok_new_line)
	}

	return text
}
