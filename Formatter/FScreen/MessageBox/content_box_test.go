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
		if line[2] == ' ' {
			break
		}

		fmt.Println(string(line))
	}

	if string(testBox.table[0][0:5]) != "Hello" || string(testBox.table[0][6:11]) != "World" {
		t.Errorf("WriteLines did not correctly write short lines")
	}
}

func TestWriteLines_LongLine(t *testing.T) {
	testBox, _ := NewContentBox(18, 72)
	testBox.EnqueueContents([]string{"This is really a very long line that should be truncated and end with an ellipsis"}, StyleMap[NormalText])

	// DEBUG: Print the table
	for _, line := range testBox.table {
		if line[2] == ' ' {
			break
		}

		fmt.Println(string(line))
	}

	if string(testBox.table[0][0:17]) != "This is really..." {
		t.Errorf("WriteLines did not correctly truncate a long line")
	}
}

func TestWriteLines_ShiftUp(t *testing.T) {
	testBox, _ := NewContentBox(18, 6)

	contents := make([]string, 8)

	for i := range contents {
		contents[i] = "Line"
		testBox.EnqueueContents([]string{"Line"}, StyleMap[NormalText])

		if testBox.CanShiftUp() {
			testBox.ShiftUp()
		}
	}

	// DEBUG: Print the table
	for _, line := range testBox.table {
		if line[2] == ' ' {
			break
		}

		fmt.Println(string(line))
	}

	if string(testBox.table[4][0:4]) != "Line" {
		t.Errorf("WriteLines did not correctly shift the screen up")
	}
}
