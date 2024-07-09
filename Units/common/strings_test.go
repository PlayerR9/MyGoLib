package common

import (
	"testing"
)

func TestOrString(t *testing.T) {
	TestValues := []string{"a", "b", "c "}

	str := OrString(TestValues, false)
	if str != "a, b, or c" {
		t.Errorf("OrString(%q) = %q; want %q", TestValues, str, "a, b, or c")
	}
}
