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

type TokenStatus int

const (
	TkComplete TokenStatus = iota
	TkIncomplete
	TkError
)

func (s TokenStatus) String() string {
	return [...]string{
		"complete",
		"incomplete",
		"error",
	}[s]
}

type Tokener interface {
	fmt.Stringer

	GetID() string
	GetData() any
	GetPos() int
	GetStatus() TokenStatus
	SetStatus(TokenStatus)
}

type LeafToken struct {
	id     string
	data   string
	at     int
	status TokenStatus
}

func (t *LeafToken) String() string {
	if t == nil {
		return "LeafToken[nil]"
	}

	return fmt.Sprintf("LeafToken[id=%s, data='%s', at=%d, status=%v]",
		t.id,
		t.data,
		t.at,
		t.status,
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

func (t *LeafToken) GetStatus() TokenStatus {
	return t.status
}

func (t *LeafToken) SetStatus(status TokenStatus) {
	t.status = status
}

func NewLeafToken(id string, data string, at int) *LeafToken {
	return &LeafToken{
		id,
		data,
		at,
		TkIncomplete,
	}
}

type NonLeafToken struct {
	id        string
	data      []Tokener
	at        int
	status    TokenStatus
	lookahead string
}

func (t *NonLeafToken) String() string {
	if t == nil {
		return "NonLeafToken[nil]"
	}

	if len(t.data) == 0 {
		return fmt.Sprintf("NonLeafToken[id=%s, data=[], at=%d, status=%v]",
			t.id,
			t.at,
			t.status,
		)
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "NonLeafToken[id=%s, data=[%v", t.id, t.data[0])

	for _, token := range t.data[1:] {
		fmt.Fprintf(&builder, ", %v", token)
	}

	fmt.Fprintf(&builder, "], at=%d, status=%v]",
		t.at,
		t.status,
	)

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

func (t *NonLeafToken) GetStatus() TokenStatus {
	return t.status
}

func (t *NonLeafToken) SetStatus(status TokenStatus) {
	t.status = status
}

func NewNonLeafToken(id string, at int, lookahead string, data ...Tokener) *NonLeafToken {
	return &NonLeafToken{
		id:        id,
		at:        at,
		data:      data,
		status:    TkIncomplete,
		lookahead: lookahead,
	}
}
