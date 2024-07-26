package runes

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// WordMatcher is the word matcher.
type WordMatcher struct {
	// words is the list of words.
	words [][]rune
}

// NewWordMatcher returns a new WordMatcher.
//
// Returns:
//   - *WordMatcher: The new WordMatcher. Never nil.
func NewWordMatcher() *WordMatcher {
	return &WordMatcher{
		words: make([][]rune, 0),
	}
}

// AddWord adds a word to the matcher. It ignores empty or duplicated words.
//
// Parameters:
//   - word: The word to add.
//
// Returns:
//   - error: An error if the word is invalid.
//
// Errors:
//   - *common.ErrAt: When the word is not a valid UTF-8 string.
func (wm *WordMatcher) AddWord(word string) error {
	if word == "" {
		return nil
	}

	chars, err := StringToUtf8(word)
	if err != nil {
		return err
	}

	uc.Assert(len(chars) > 0, "chars is empty")

	var indices []int

	for i, w := range wm.words {
		uc.Assert(len(w) > 0, "word is empty")

		if len(w) == len(chars) && w[0] == chars[0] {
			indices = append(indices, i)
		}
	}

	for i := 1; i < len(chars) && len(indices) > 0; i++ {
		filter := func(idx int) bool {
			return wm.words[idx][i] == chars[i]
		}

		indices = us.SliceFilter(indices, filter)
	}

	if len(indices) == 0 {
		wm.words = append(wm.words, chars)
	}

	return nil
}

// Match matches the input stream.
//
// Parameters:
//   - is: The input stream to match.
//
// Returns:
//   - string: The matched word.
//   - error: An error if the stream could not be matched.
//
// Errors:
//   - *common.ErrAt: If the input stream is not a valid UTF-8 stream.
func (wm *WordMatcher) Match(is CharStream) (string, error) {
	if is == nil {
		return "", uc.NewErrNilParameter("is")
	}

	indices := make([]int, 0, len(wm.words))

	for i := range wm.words {
		indices = append(indices, i)
	}

	m := &matcher{
		indices: indices,
		pos:     0,
		wm:      wm,
	}

	var next_count int

	for {
		char, ok := is.Peek() // Assume anything
		if !ok {
			word, err := m.get_sol()

			size := utf8.RuneCountInString(word)
			for i := size; i < next_count; i++ {
				ok := is.Refuse()
				uc.AssertOk(ok, "is.Refuse()")
			}

			if err != nil {
				return "", err
			}

			return word, nil
		}

		ok = m.match(char)
		if ok {
			word, err := m.get_sol()

			size := utf8.RuneCountInString(word)
			for i := size; i < next_count; i++ {
				ok := is.Refuse()
				uc.AssertOk(ok, "is.Refuse()")
			}

			if err != nil {
				return "", err
			}

			return word, nil
		}

		// We can continue to match...
		is.Next()
		next_count++
	}
}

// get_word_at returns the word at the given index.
//
// Parameters:
//   - idx: The index of the word to return.
//
// Returns:
//   - []rune: The word.
//   - bool: True if the word was found, false otherwise.
func (wm *WordMatcher) get_word_at(idx int) ([]rune, bool) {
	if idx < 0 || idx >= len(wm.words) {
		return nil, false
	}

	return wm.words[idx], true
}

// matcher is a helper struct for controlling the matching process.
type matcher struct {
	// indices is the list of indices of the words that were matched.
	indices []int

	// pos is the current position in the stream.
	pos int

	// word is the current word being matched.
	word strings.Builder

	// wm is the WordMatcher that is being used.
	wm *WordMatcher
}

// match is a helper function that matches the next character in the stream.
//
// Returns:
//   - bool: True if the matching process cannot continue. False otherwise.
func (m *matcher) match(char rune) bool {
	uc.Assert(len(m.indices) > 0, "m.indices is empty")

	var top int

	for i := 0; i < len(m.indices); i++ {
		idx := m.indices[i]

		word, ok := m.wm.get_word_at(idx)
		uc.AssertOk(ok, "wm.get_word_at(%d)", idx)

		if m.pos >= len(word) || word[m.pos] == char {
			m.indices[top] = idx
			top++
		}
	}

	if top == 0 {
		return true
	}

	m.pos++
	m.indices = m.indices[:top]
	m.word.WriteRune(char)

	return false
}

// get_sol is a helper function that returns the solution.
//
// Returns:
//   - string: The solution.
//   - error: An error if the solution could not be found.
func (m *matcher) get_sol() (string, error) {
	var indices []int

	if len(m.indices) > 0 {
		// Try with the words that were not discarded yet.
		for _, idx := range m.indices {
			word, ok := m.wm.get_word_at(idx)
			uc.AssertOk(ok, "wm.get_word_at(%d)", idx)

			if len(word) <= m.pos {
				indices = append(indices, idx)
			}
		}
	}

	if len(indices) == 0 {
		if m.word.Len() == 0 {
			return "", errors.New("no matches found")
		} else {
			return "", fmt.Errorf("no matches found for %q", m.word.String())
		}
	}

	// Find the longest matching word.

	var max_len int
	var max_idx []int

	for _, idx := range indices {
		word, ok := m.wm.get_word_at(idx)
		uc.AssertOk(ok, "wm.get_word_at(%d)", idx)

		if len(max_idx) == 0 || len(word) > max_len {
			max_len = len(word)
			max_idx = []int{idx}
		} else if len(word) == max_len {
			max_idx = append(max_idx, idx)
		}
	}

	if len(max_idx) > 1 {
		values := make([]string, 0, len(max_idx))

		for _, idx := range max_idx {
			word, ok := m.wm.get_word_at(idx)
			uc.AssertOk(ok, "wm.get_word_at(%d)", idx)

			values = append(values, string(word))
		}

		return "", fmt.Errorf("multiple matches found for %q: %s", m.word.String(), strings.Join(values, ", "))
	}

	word, ok := m.wm.get_word_at(max_idx[0])
	uc.AssertOk(ok, "wm.get_word_at(%d)", max_idx[0])

	return string(word), nil
}

// MultiMatcher kinda works like WordMatcher but, unlike WordMatcher, it only matches a specific set of characters.
//
// Parameters:
//   - chars: The characters to match.
//   - stream: The CharStream to use.
//
// Returns:
//   - string: The matched string.
//   - error: An error if the matching process failed.
//
// Errors:
//   - *common.ErrInvalidParameter: If the input stream is nil or the input characters are empty.
func MultiMatcher(chars []rune, stream CharStream) (string, error) {
	if stream == nil {
		return "", uc.NewErrNilParameter("stream")
	} else if len(chars) == 0 {
		return "", uc.NewErrInvalidParameter("chars", uc.NewErrEmpty(chars))
	}

	var builder strings.Builder
	var size int

	for _, c := range chars {
		char, ok := stream.Peek()
		if !ok {
			for size > 0 {
				ok := stream.Refuse()
				uc.Assert(ok, "stream.Refuse()")
			}

			return "", fmt.Errorf("expected '%c', got nothing instead", c)
		}

		builder.WriteRune(char)

		if char != c {
			for size > 0 {
				ok := stream.Refuse()
				uc.Assert(ok, "stream.Refuse()")
			}

			return "", fmt.Errorf("expected '%c', got '%c' instead", c, char)
		}

		stream.Next() // Consume the peeked char
		size++
	}

	return builder.String(), nil
}
