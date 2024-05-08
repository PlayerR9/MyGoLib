package StringExt

import (
	"testing"
)

func TestReplaceSuffix(t *testing.T) {
	const (
		Str    string = "hello world"
		Suffix string = "Bob"
	)

	result, err := ReplaceSuffix(Str, Suffix)
	if err != nil {
		t.Errorf("expected no error, got %s instead", err.Error())
	}

	if result != "hello woBob" {
		t.Errorf("expected 'hello woBob', got %s instead", result)
	}
}

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

func TestSplitSentenceIntoFields(t *testing.T) {
	const (
		Sentence string = "\tHello, \vworld!\nThis is a test.\r\n"
		Indent   int    = 3
	)

	lines, err := AdvancedFieldsSplitter(Sentence, Indent)
	if err != nil {
		t.Errorf("expected no error, got %s instead", err.Error())
	}

	if len(lines) != 2 {
		t.Errorf("expected 2, got %d instead", len(lines))
	}

	if len(lines[0]) != 2 {
		t.Errorf("expected 2, got %d instead", len(lines[0]))
	}

	if len(lines[1]) != 4 {
		t.Errorf("expected 4, got %d instead", len(lines[1]))
	}

	if lines[0][0] != "   Hello," {
		t.Errorf("expected '   Hello,', got %s instead", lines[0][0])
	}

	if lines[0][1] != " world!" {
		t.Errorf("expected ' world!', got %s instead", lines[0][1])
	}

	if lines[1][0] != "This" {
		t.Errorf("expected 'This', got %s instead", lines[1][0])
	}

	if lines[1][1] != "is" {
		t.Errorf("expected 'is', got %s instead", lines[1][1])
	}

	if lines[1][2] != "a" {
		t.Errorf("expected 'a', got %s instead", lines[1][2])
	}

	if lines[1][3] != "test." {
		t.Errorf("expected 'test.', got %s instead", lines[1][3])
	}
}
