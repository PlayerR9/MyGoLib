package runes

import "testing"

func TestDuplicates(t *testing.T) {
	wm := NewWordMatcher()

	err := wm.AddWord("hello")
	if err != nil {
		t.Errorf("error adding word: %s", err.Error())
	}

	err = wm.AddWord("hello1")
	if err != nil {
		t.Errorf("error adding word: %s", err.Error())
	}

	err = wm.AddWord("hello")
	if err == nil {
		t.Errorf("expected error adding duplicate word")
	}
}

func TestMatch(t *testing.T) {
	test_words := []string{"foo", "bar", "baz", "foobar"}

	wm := NewWordMatcher()

	for _, word := range test_words {
		err := wm.AddWord(word)
		if err != nil {
			t.Errorf("error adding word: %s", err.Error())
		}
	}

	is := NewStream([]rune("foobar"))

	word, err := wm.Match(is)
	if err != nil {
		t.Errorf("error matching word: %s", err.Error())
	}

	if word != "foobar" {
		t.Errorf("expected word to be 'foobar', got '%s'", word)
	}
}

func TestMiddleMatch(t *testing.T) {
	test_words := []string{"foo", "bar", "baz", "foobar"}

	wm := NewWordMatcher()

	for _, word := range test_words {
		err := wm.AddWord(word)
		if err != nil {
			t.Errorf("error adding word: %s", err.Error())
		}
	}

	is := NewStream([]rune("fooba"))

	word, err := wm.Match(is)
	if err != nil {
		t.Errorf("error matching word: %s", err.Error())
	}

	if word != "foo" {
		t.Errorf("expected word to be 'foo', got '%s'", word)
	}
}
