package ContentBox

import (
	"testing"
)

func TestWriteLines_ShortLines(t *testing.T) {
	mlt := NewMultiLineText("Hello", "World")

	lines, err := mlt.Draw(18, 2)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
	}

	if len(lines) != 1 {
		t.Fatalf("Expected 1 lines, but got %d", len(lines))
	}

	if lines[0] != "Hello World" {
		t.Errorf("Expected first line to be 'Hello World', but got '%s'", lines[0])
	}
}

func TestWriteLines_LongLine(t *testing.T) {
	mlt := NewMultiLineText(
		"This is really a very long line that should be truncated and end with an ellipsis",
	)

	lines, err := mlt.Draw(18, 1)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
	}

	if len(lines) != 1 {
		t.Fatalf("Expected 1 lines, but got %d", len(lines))
	}

	if lines[0] != "This is really... " {
		t.Errorf("Expected first line to be 'This is really... ', but got '%s'", lines[0])
	}
}
