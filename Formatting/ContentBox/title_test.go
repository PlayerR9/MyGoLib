package ContentBox

import (
	"testing"
)

func TestTitle(t *testing.T) {
	const (
		Title        string = "Test Title"
		ExpectedLine string = " *** Test Title *** "
	)

	title := NewTitleBox(Title)

	lines, err := title.DrawTitle(20)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
	}

	if len(lines) != 1 {
		t.Errorf("Expected 1 line, but got %d", len(lines))
	}

	if lines[0] != ExpectedLine {
		t.Errorf("Expected line to be '%s', but got '%s'", ExpectedLine, lines[0])
	}
}

func TestMiddleSplit(t *testing.T) {
	type titleTest struct {
		title         string
		width         int
		height        int
		expectedLines []string
	}

	tests := []titleTest{
		{
			title:  "This is a very long title",
			width:  13,
			height: 5,
			expectedLines: []string{
				"*** This *** ",
				"*** is a *** ",
				"*** very *** ",
				"*** long *** ",
				"*** title ***",
			},
		},
		{
			title:  "Hello world, this is a test",
			width:  19,
			height: 3,
			expectedLines: []string{
				"   *** Hello ***   ",
				"*** world, this ***",
				" *** is a test *** ",
			},
		},
		{
			title:  "Hi You They",
			width:  14,
			height: 2,
			expectedLines: []string{
				"*** Hi You ***",
				" *** They *** ",
			},
		},
	}

	for i, test := range tests {
		title := NewTitleBox(test.title)

		lines, err := title.DrawTitle(test.width)
		if err != nil {
			t.Errorf("At test %d, expected no error, but got %s", i, err.Error())
		}

		if len(lines) != len(test.expectedLines) {
			t.Errorf("At test %d, expected %d lines, but got %d", i, len(test.expectedLines), len(lines))
		}

		for j := 0; j < len(lines); j++ {
			if lines[j] != test.expectedLines[j] {
				t.Errorf("At test %d, expected line %d to be '%s', but got '%s'", i, j, test.expectedLines[j], lines[j])
			}
		}
	}
}

func TestTitleTruncation(t *testing.T) {
	const (
		Title        string = "Thisisaverylongtitle"
		ExpectedLine string = "*** Th... ***"
	)

	title := NewTitleBox(Title)

	lines, err := title.DrawTitle(13)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
	}

	if len(lines) != 1 {
		t.Errorf("Expected 1 line, but got %d", len(lines))
	}

	if lines[0] != ExpectedLine {
		t.Errorf("Expected line to be '%s', but got '%s'", ExpectedLine, lines[0])
	}
}
