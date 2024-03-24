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
	ID   string
	Data string
	At   int
}

func (t *LeafToken) String() string {
	if t == nil {
		return "LeafToken[nil]"
	}

	return fmt.Sprintf("LeafToken[id=%s, data='%s', at=%d]",
		t.ID,
		t.Data,
		t.At,
	)
}

func (t *LeafToken) GetID() string {
	return t.ID
}

func (t *LeafToken) GetData() any {
	return t.Data
}

func (t *LeafToken) GetPos() int {
	return t.At
}

func NewLeafToken(id string, data string, at int) *LeafToken {
	return &LeafToken{
		id,
		data,
		at,
	}
}

type NonLeafToken struct {
	ID   string
	Data []Tokener
	At   int
}

func (t *NonLeafToken) String() string {
	if t == nil {
		return "NonLeafToken[nil]"
	}

	if len(t.Data) == 0 {
		return fmt.Sprintf("NonLeafToken[id=%s, data=[], at=%d]",
			t.ID,
			t.At,
		)
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "NonLeafToken[id=%s, data=[%v", t.ID, t.Data[0])

	for _, token := range t.Data[1:] {
		fmt.Fprintf(&builder, ", %v", token)
	}

	fmt.Fprintf(&builder, "], at=%d]", t.At)

	return builder.String()
}

func (t *NonLeafToken) GetID() string {
	return t.ID
}

func (t *NonLeafToken) GetData() any {
	return t.Data
}

func (t *NonLeafToken) GetPos() int {
	return t.At
}

func NewNonLeafToken(id string, at int, data ...Tokener) *NonLeafToken {
	return &NonLeafToken{
		ID:   id,
		At:   at,
		Data: data,
	}
}
