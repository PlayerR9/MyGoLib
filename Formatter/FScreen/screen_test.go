package FScreen

import (
	"testing"
)

func TestWriteLines_ShortLines(t *testing.T) {
	Initialize("Test")
	Start()

	SendMessages(
		MessagePrint("Hello", "World"),
		MessageClose(),
	)

	Wait()

	if string(screen[1][2:7]) != "Hello" || string(screen[1][8:13]) != "World" {
		t.Errorf("WriteLines did not correctly write short lines")
	}
}

func TestWriteLines_LongLine(t *testing.T) {
	Initialize("Test")
	Start()

	SendMessages(
		MessagePrint("This is really a very long line that should be truncated and end with an ellipsis"),
		MessageClose(),
	)

	Wait()

	if string(screen[1][screenWidth-5:screenWidth-2]) != "..." {
		t.Errorf("WriteLines did not correctly truncate a long line")
	}
}

func TestWriteLines_ShiftUp(t *testing.T) {
	Initialize("Test")
	Start()

	contents := make([]string, screenHeight+1)

	for i := range contents {
		contents[i] = "Line"
		SendMessages(
			MessagePrint("Line"),
			MessageClose(),
		)
	}

	Wait()

	if string(screen[6][2:6]) != "Line" {
		t.Errorf("WriteLines did not correctly shift the screen up")
	}

	Start()

	SendMessages(
		MessageClearScreen(),
		MessagePrint(contents[0], contents[1:]...),
		MessageClose(),
	)

	Wait()

	if string(screen[1][52:56]) != "Line" {
		t.Errorf("WriteLines did not correctly shift the screen up")
	}
}
