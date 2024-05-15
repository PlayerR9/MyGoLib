package Components

import (
	"testing"

	cdd "github.com/PlayerR9/MyGoLib/ComplexData/Display"
)

func TestWriteLines_ShortLines(t *testing.T) {
	mlt := NewMultiLineText()

	err := mlt.AppendSentence("Hello World")
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	table, err := cdd.NewDrawTable(18, 2)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	err = mlt.Draw(table)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	lines := table.GetLines()

	if lines[0] != "Hello World       " {
		t.Errorf("Expected first line to be 'Hello World       ', but got '%s'", lines[0])
	}
}

func TestWriteLines_LongLine(t *testing.T) {
	mlt := NewMultiLineText()

	err := mlt.AppendSentence(
		"This is really a very long line that should be truncated and end with an ellipsis",
	)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	table, err := cdd.NewDrawTable(18, 1)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	err = mlt.Draw(table)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	lines := table.GetLines()

	if lines[0] != "This is really... " {
		t.Fatalf("Expected first line to be 'This is really... ', but got '%s'", lines[0])
	}
}
