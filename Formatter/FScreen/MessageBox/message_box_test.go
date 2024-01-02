package MessageBox

import (
	"testing"
)

func TestWriteLines_ShortLines(t *testing.T) {
	testBox := NewMessageBox(80, 20)
	go testBox.Run()

	testBox.SendMessages(
		NewTextMessage(NormalText,
			"Hello",
			"World",
		),
	)

	testBox.Pause()

	if string(testBox.table[1][2:7]) != "Hello" || string(testBox.table[1][8:13]) != "World" {
		t.Errorf("WriteLines did not correctly write short lines")
	}
}

func TestWriteLines_LongLine(t *testing.T) {
	testBox := NewMessageBox(10, 10)
	go testBox.Run()

	testBox.SendMessages(
		NewTextMessage(NormalText,
			"This is really a very long line that should be truncated and end with an ellipsis",
		),
	)

	testBox.Fini()

	if string(testBox.table[1][testBox.height-5:testBox.width-2]) != "..." {
		t.Errorf("WriteLines did not correctly truncate a long line")
	}
}

func TestWriteLines_ShiftUp(t *testing.T) {
	testBox := NewMessageBox(10, 10)
	go testBox.Run()

	contents := make([]string, testBox.height+1)

	for i := range contents {
		contents[i] = "Line"
		testBox.SendMessages(
			NewTextMessage(NormalText, "Line"),
		)
	}

	testBox.Fini()

	if string(testBox.table[6][2:6]) != "Line" {
		t.Errorf("WriteLines did not correctly shift the screen up")
	}

	testBox = NewMessageBox(10, 10)
	go testBox.Run()

	testBox.SendMessages(
		NewTextMessage(NormalText,
			contents...,
		),
	)

	testBox.Fini()

	if string(testBox.table[1][52:56]) != "Line" {
		t.Errorf("WriteLines did not correctly shift the screen up")
	}
}
