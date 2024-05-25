package Document

import (
	"testing"
)

func TestFString(t *testing.T) {
	doc := NewDocument("Hello", "World")

	// FIXME: This test is not complete.
	t.Fatalf("Expected no error, but got %v", doc)

	/*
		expected := make([]string, 0)

		for _, line := range trav.GetLines() {
			expected = append(expected, line.GetLines()...)
		}

		var builder Builder

		builder.SetSeparator(nil)

		actual := builder.Build().Apply(doc.Tmp())

		if len(actual) != len(expected) {
			// DEBUG: Show the expected and actual values.
			fmt.Println("Expected:", expected)

			fmt.Println("Actual:", actual)

			t.Errorf("Expected: %v, Actual: %v", expected, actual)
		}

		for i := range actual {
			if actual[i] != expected[i] {
				t.Errorf("Expected: %v, Actual: %v", expected, actual)
			}
		}
	*/
}
