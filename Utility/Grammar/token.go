package Grammar

import (
	"fmt"
	"strings"
	"unicode"
)

func IsTerminal(identifier string) bool {
	firstLetter := []rune(identifier)[0]

	return unicode.IsLetter(firstLetter) && unicode.IsUpper(firstLetter)
}

type Tokener interface {
	fmt.Stringer

	GetID() string
	GetData() any
	GetPos() int
}

type LeafToken struct {
	id   string
	data string
	at   int
}

func (t *LeafToken) String() string {
	if t == nil {
		return "LeafToken[nil]"
	}

	return fmt.Sprintf("LeafToken[id=%s, data='%s', at=%d]",
		t.id,
		t.data,
		t.at,
	)
}

func (t *LeafToken) GetID() string {
	return t.id
}

func (t *LeafToken) GetData() any {
	return t.data
}

func (t *LeafToken) GetPos() int {
	return t.at
}

func NewLeafToken(id string, data string, at int) *LeafToken {
	return &LeafToken{
		id,
		data,
		at,
	}
}

type NonLeafToken struct {
	id   string
	data []Tokener
	at   int
}

func (t *NonLeafToken) String() string {
	if t == nil {
		return "NonLeafToken[nil]"
	}

	if len(t.data) == 0 {
		return fmt.Sprintf("NonLeafToken[id=%s, data=[], at=%d]",
			t.id,
			t.at,
		)
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "NonLeafToken[id=%s, data=[%v", t.id, t.data[0])

	for _, token := range t.data[1:] {
		fmt.Fprintf(&builder, ", %v", token)
	}

	fmt.Fprintf(&builder, "], at=%d]", t.at)

	return builder.String()
}

func (t *NonLeafToken) GetID() string {
	return t.id
}

func (t *NonLeafToken) GetData() any {
	return t.data
}

func (t *NonLeafToken) GetPos() int {
	return t.at
}

func NewNonLeafToken(id string, at int, data ...Tokener) *NonLeafToken {
	return &NonLeafToken{
		id:   id,
		at:   at,
		data: data,
	}
}
