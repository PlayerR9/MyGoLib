package Go

import (
	"testing"
)

func TestParseFields(t *testing.T) {
	const (
		Input string = "a int"
	)

	fields, err := ParseFields(Input)
	if err != nil {
		t.Errorf("expected nil, got %s instead", err.Error())
	}

	if len(fields) != 1 {
		t.Errorf("expected 1, got %d instead", len(fields))
	}

	val, ok := fields["a"]
	if !ok {
		t.Errorf("expected true, got false instead")
	}

	if val != "int" {
		t.Errorf("expected 'int', got '%s' instead", val)
	}
}
