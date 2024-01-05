package MessageBox

import (
	"fmt"
	"testing"
)

func TestWriteLines_ShortLines(t *testing.T) {
	testBox, _ := NewContentBox(76, 18)
	testBox.EnqueueContents([]string{"Hello", "World"}, StyleMap[NormalText])

	// DEBUG: Print the table
	for _, line := range testBox.table {
		fmt.Println(string(line))
	}

	if string(testBox.table[0][0:5]) != "Hello" || string(testBox.table[0][6:11]) != "World" {
		t.Errorf("WriteLines did not correctly write short lines")
	}
}

func TestWriteLines_LongLine(t *testing.T) {
	testBox := new(MessageBox)
	sendTo, _ := testBox.Init(80, 20)

	sendTo <- NewTextMessage(NormalText,
		"This is really a very long line that should be truncated and end with an ellipsis",
	)

	testBox.Fini()

	if string(testBox.content.table[1][testBox.content.height-5:testBox.content.width-2]) != "..." {
		t.Errorf("WriteLines did not correctly truncate a long line")
	}
}

func TestWriteLines_ShiftUp(t *testing.T) {
	testBox := new(MessageBox)
	sendTo, _ := testBox.Init(80, 20)

	contents := make([]string, testBox.content.height+1)

	for i := range contents {
		contents[i] = "Line"
		sendTo <- NewTextMessage(NormalText, "Line")
	}

	testBox.Fini()

	if string(testBox.content.table[6][2:6]) != "Line" {
		t.Errorf("WriteLines did not correctly shift the screen up")
	}

	sendTo, _ = testBox.Init(80, 20)

	sendTo <- NewTextMessage(NormalText,
		contents...,
	)

	testBox.Fini()

	if string(testBox.content.table[1][52:56]) != "Line" {
		t.Errorf("WriteLines did not correctly shift the screen up")
	}
}
