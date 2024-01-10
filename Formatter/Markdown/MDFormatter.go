package Markdown

import (
	"fmt"
)

type TableAllignmentType int

const (
	TableStyleCenter TableAllignmentType = iota
	TableStyleLeft
	TableStyleRight
)

func (a TableAllignmentType) String() string {
	return [...]string{
		":---:",
		":---",
		"---:",
	}[a]
}

var (
	tokPipe      string = " | "
	tokNewLine   string = "\n"
	PedanticMode bool   = true
)

type Header struct {
	Text  string
	Style TableAllignmentType
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

	if !PedanticMode {
		for i, row := range t.rows {
			if len(row) < len(t.headers) {
				panic(fmt.Sprintf("row %d has less elements than the number of headers", i))
			}
		}
	} else {
		for _, row := range t.rows {
			if len(row) >= len(t.headers) {
				continue
			}

			for len(row) < len(t.headers) {
				row = append(row, "")
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
	var text []string

	appendRow := func(row []string) {
		text = append(text, tokPipe)
		for _, element := range row {
			text = append(text, element, tokPipe)
		}
		text = append(text, tokNewLine)
	}

	headerTexts := make([]string, len(t.headers))
	headerStyles := make([]string, len(t.headers))
	for i, header := range t.headers {
		headerTexts[i] = header.Text
		headerStyles[i] = header.Style.String()
	}

	appendRow(headerTexts)
	appendRow(headerStyles)

	for _, row := range t.rows {
		appendRow(row)
	}

	return text
}
