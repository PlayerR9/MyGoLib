package Document

import (
	"fmt"
	"testing"

	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
)

func TestFString(t *testing.T) {
	doc := NewDocument("Hello", "World")

	expected := doc.FString(0)

	var builder ffs.Builder

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
}
