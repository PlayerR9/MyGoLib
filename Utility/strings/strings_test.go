package strings

import (
	"testing"
)

func TestFindContentIndexes(t *testing.T) {
	const (
		OpToken string = "("
		ClToken string = ")"
	)

	var (
		ContentTokens []string = []string{
			"(", "(", "a", "+", "b", ")", "*", "c", ")", "+", "d",
		}
	)

	indices, err := FindContentIndexes(OpToken, ClToken, ContentTokens)
	if err != nil {
		t.Errorf("expected no error, got %s instead", err.Error())
	}

	if indices[0] != 1 {
		t.Errorf("expected 1, got %d instead", indices[0])
	}

	if indices[1] != 9 {
		t.Errorf("expected 9, got %d instead", indices[1])
	}
}

func TestOrString(t *testing.T) {
	TestValues := []string{"a", "b", "c "}

	str := OrString(TestValues, false, false)
	if str != "a, b, or c" {
		t.Errorf("OrString(%q) = %q; want %q", TestValues, str, "a, b, or c")
	}
}
