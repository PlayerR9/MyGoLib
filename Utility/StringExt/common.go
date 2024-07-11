package StringExt

import (
	"crypto/rand"
	"encoding/hex"
	"math"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
	mext "github.com/PlayerR9/MyGoLib/Utility/MathExt"
)

var (
	// calculate_split_ratio is a function that calculates the split ratio of a
	// TextSplit instance.
	//
	// Never false.
	calculate_split_ratio us.WeightFunc[*TextSplit]
)

func init() {
	calculate_split_ratio = func(candidate *TextSplit) (float64, bool) {
		uc.Assert(candidate != nil, "in calculateSplitRatio: candidate is nil")

		values := make([]float64, 0, candidate.GetHeight())

		for _, line := range candidate.lines {
			values = append(values, float64(line.len))
		}

		sqm, ok := mext.SQM(values)
		uc.Assert(ok, "in calculateSplitRatio: failed to calculate SQM as slice is empty")

		return sqm, true
	}
}

// ToUTF8Runes converts a string to a slice of runes.
//
// Parameters:
//   - s: The string to convert.
//
// Returns:
//   - runes: The slice of runes.
//   - error: An error of type *common.ErrAt if an invalid UTF-8 encoding is found.
//
// Behaviors:
//   - An empty string returns a nil slice with no errors.
//   - The function stops at the first invalid UTF-8 encoding; returning an
//     error and the runes found up to that point.
//   - The function converts '\r\n' to '\n'.
func ToUTF8Runes(s string) (runes []rune, err error) {
	if s == "" {
		return
	}

	for i := 0; len(s) > 0; i++ {
		r, size := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError {
			err = uc.NewErrAt(i+1, "character", NewErrInvalidUTF8Encoding())
			return
		}

		s = s[size:]

		if r != '\r' || len(s) == 0 {
			runes = append(runes, r)
			continue
		}

		other_r, size := utf8.DecodeRuneInString(s)
		if other_r == utf8.RuneError {
			err = uc.NewErrAt(i+2, "character", NewErrInvalidUTF8Encoding())
			return
		}

		if other_r == '\n' {
			runes = append(runes, '\n')
		} else {
			runes = append(runes, r, other_r)
		}

		s = s[size:]
	}

	return
}

// ReplaceSuffix replaces the end of the given string with the provided suffix.
//
// Parameters:
//   - str: The original string.
//   - suffix: The suffix to replace the end of the string.
//
// Returns:
//   - string: The resulting string after replacing the end with the suffix.
//   - bool: A boolean indicating if the operation was successful. (i.e., if the
//     suffix is shorter than the string).
//
// Behaviors:
//   - For quick error, use the *ErrLongerSuffix error type of this package.
//
// Examples:
//
//	const (
//		str    string = "hello world"
//		suffix string = "Bob"
//	)
//
//	result, err := ReplaceSuffix(str, suffix)
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(result) // Output: hello woBob
//	}
func ReplaceSuffix(str, suffix string) (string, bool) {
	if suffix == "" {
		return str, true
	}

	count_str := utf8.RuneCountInString(str)
	count_suffix := utf8.RuneCountInString(suffix)

	if count_str < count_suffix {
		return "", false
	}

	if count_str == count_suffix {
		return suffix, true
	}

	var builder strings.Builder

	builder.WriteString(str[:count_str-count_suffix])
	builder.WriteString(suffix)

	s := builder.String()

	return s, true
}

// FindContentIndexes searches for the positions of opening and closing
// tokens in a slice of strings.
//
// Parameters:
//   - op_token: The string that marks the beginning of the content.
//   - cl_token: The string that marks the end of the content.
//   - tokens: The slice of strings in which to search for the tokens.
//
// Returns:
//   - result: An array of two integers representing the start and end indexes
//     of the content.
//   - err: Any error that occurred while searching for the tokens.
//
// Errors:
//   - *uc.ErrInvalidParameter: If the openingToken or closingToken is an
//     empty string.
//   - *ErrTokenNotFound: If the opening or closing token is not found in the
//     content.
//   - *ErrNeverOpened: If the closing token is found without any
//     corresponding opening token.
//
// Behaviors:
//   - The first index of the content is inclusive, while the second index is
//     exclusive.
//   - This function returns a partial result when errors occur. ([-1, -1] if
//     errors occur before finding the opening token, [index, 0] if the opening
//     token is found but the closing token is not found.
func FindContentIndexes(op_token, cl_token string, tokens []string) (result [2]int, err error) {
	result[0] = -1
	result[1] = -1

	if op_token == "" {
		err = uc.NewErrInvalidParameter("openingToken", uc.NewErrEmpty(op_token))
		return
	} else if cl_token == "" {
		err = uc.NewErrInvalidParameter("closingToken", uc.NewErrEmpty(cl_token))
		return
	}

	op_tok_idx := slices.Index(tokens, op_token)
	if op_tok_idx < 0 {
		err = NewErrTokenNotFound(op_token, OpToken)
		return
	} else {
		result[0] = op_tok_idx + 1
	}

	balance := 1
	cl_tok_idx := -1

	for i := result[0]; i < len(tokens) && cl_tok_idx == -1; i++ {
		curr_tok := tokens[i]

		if curr_tok == cl_token {
			balance--

			if balance == 0 {
				cl_tok_idx = i
			}
		} else if curr_tok == op_token {
			balance++
		}
	}

	if cl_tok_idx != -1 {
		result[1] = cl_tok_idx + 1
		return
	}

	if balance < 0 {
		err = NewErrNeverOpened(op_token, cl_token)
		return
	} else if balance != 1 || cl_token != "\n" {
		err = NewErrTokenNotFound(cl_token, ClToken)
		return
	}

	result[1] = len(tokens)
	return
}

// GenerateID generates a random ID of the specified size (in bytes).
//
// Parameters:
//   - size: The size of the ID to generate (in bytes).
//
// Returns:
//   - string: The generated ID.
//   - error: An error if the ID cannot be generated.
//
// Errors:
//   - *uc.ErrInvalidParameter: If the size is less than 1.
//   - Any error returned by the rand.Read function.
//
// Behaviors:
//   - The function uses the crypto/rand package to generate a random ID of the
//     specified size.
//   - The ID is returned as a hexadecimal string.
func GenerateID(size int) (string, error) {
	if size < 1 {
		return "", uc.NewErrInvalidParameter("size", uc.NewErrGT(0))
	}

	b := make([]byte, size) // 128 bits

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	encoded := hex.EncodeToString(b)

	return encoded, nil
}

// FitString fits a string to the specified width by adding spaces to the end
// of the string until the width is reached.
//
// Parameters:
//   - width: The width to fit the string to.
//
// Returns:
//   - string: The string with spaces added to the end to fit the width.
//
// Behaviors:
//   - If the width is less than 0, it is set to 0.
//   - If the width is less than the length of the string, the string is
//     truncated to fit the width.
//   - If the width is greater than the length of the string, spaces are added
//     to the end of the string until the width is reached.
func FitString(s string, width int) string {
	if width < 0 {
		width = 0
	}

	len := utf8.RuneCountInString(s)

	if width == 0 {
		return ""
	}

	switch len {
	case width:
		// Do nothing
	case 0:
		s = strings.Repeat(" ", width)
	default:
		if len < width {
			var builder strings.Builder
			spacing := strings.Repeat(" ", width-len)

			builder.WriteString(s)
			builder.WriteString(spacing)

			s = builder.String()
		} else {
			s = s[:width]
		}
	}

	return s
}

// CalculateNumberOfLines is a function that calculates the minimum number
// of lines needed to fit a given text within a specified width.
//
// Parameters:
//   - text: The slice of strings representing the text to calculate the number of
//     lines for.
//   - width: The width to fit the text within.
//
// Returns:
//   - int: The calculated number of lines needed to fit the text within the width.
//   - error: An error if it occurs during the calculation.
//
// Errors:
//   - *uc.ErrInvalidParameter: If the width is less than or equal to 0.
//   - *ErrLinesGreaterThanWords: If the calculated number of lines is greater
//     than the number of words in the text.
//
// The function calculates the total length of the text (Tl) and uses a mathematical
// formula to estimate the minimum number of lines needed to fit the text within the
// given width. The formula is explained in detail in the comments within the function.
//
// It also returns the calculated number of lines when it errors out
func CalculateNumberOfLines(text []string, width int) (int, error) {
	if width <= 0 {
		return 0, uc.NewErrInvalidParameter(
			"width",
			uc.NewErrGT(0),
		)
	} else if len(text) == 0 {
		return 0, nil
	}

	// Euristic to calculate the least amount of splits needed
	// in order to fit the text within the width.
	//
	// Assuming:
	// 	- $n$ is the number of fields in the text using strings.Fields().
	//    - $\omega$ is the number of characters in a field (word).
	//    - $W$ is the total width. This considers only usable width, i.e.,
	//    width - padding.
	//    - $Tl$ = Total length of the text. Calculated by doing the
	// 	sum of the lengths of all the words plus the number of
	// 	spaces between them. $n - 1 + \Sum_{i=1}^n \omega_i$.
	//    - $x$ = Number of splits needed to fit the text within the width.
	//
	// Formula:
	//    $\frac{Tl - x}{x + 1} \leq W$
	//
	// Explanation:
	//
	// 	- $Tl - x$: For every split, the number of characters occupied by the
	//    text is reduced as the space between the splitted fields is removed.
	//    For example: "Hello World" has 11 charactue. With one split, it becomes
	//    "Hello" and "World", which has 5 and 5 characters respectively. The
	//    total number of characters is 10, which is equal to 11 - 1.
	//		- $x + 1$: The number of lines is equal to the number of splits plus one as
	//    no splits (x = 0), gives us a single line (x + 1 = 1).
	// 	- $\frac{Tl - x}{x + 1}$: This divides the number of characters occupied by
	//    the title by the number of lines; giving us how many characters are
	//    occupied per line. $\leq W$ ensures that no line exceeds the
	//    width of the screen.
	//
	// Simplification:
	//    $$
	//    \begin{align}
	//    	\frac{Tl - x}{x + 1} &\leq W \\
	//    	Tl - x &\leq W(x + 1) \\
	//    	Tl - x &\leq Wx + W \\
	//    	Tl - W &\leq Wx + x \\
	//       Tl - W &\leq x(W + 1) \\
	//       \frac{Tl - W}{W + 1} &\leq x \\
	//       \lceil\frac{Tl - W}{W + 1}\rceil &\leq x
	//    \end{align}
	//    $$
	//   	Note: The ceil function is used as we cannot do non-integer splits.
	//		and, since we want $x$ to be greater or equal to the result of the
	//		division, we round up the result.
	//
	// Example: If we have the following text: "Hello World, this is a test",
	// with a width of 12 characters, we have:
	//    - $n = 6$
	//    - $W = 12$
	//    - $Tl = 27$
	//
	//	 	$\lceil\frac{27 - 12}{12 + 1}\rceil = \lceil\frac{15}{13}\rceil = 2$
	//
	// Solution:
	//		   *** Hello ***
	// 	*** World, this ***
	// 	 *** is a test ***

	var Tl float64

	for _, word := range text {
		// +1 for the space or any other separator
		Tl += float64(utf8.RuneCountInString(word)) + 1
	}
	Tl-- // Remove the last space or separator

	w := float64(width)

	line_count := int(math.Ceil((Tl-w)/(w+1))) + 1

	if line_count > len(text) {
		err := NewErrLinesGreaterThanWords(line_count, len(text))
		return line_count, err
	}

	return line_count, nil
}

// SplitInEqualSizedLines is a function that splits a given text into lines of
// equal width.
//
// Errors:
//   - *uc.ErrInvalidParameter: If the input text is empty or the width is less than
//     or equal to 0.
//   - *ErrLinesGreaterThanWords: If the number of lines needed to fit the text
//     within the width is greater than the number of words in the text.
//   - *ErrNoCandidateFound: If no candidate is found during the optimization process.
//
// Parameters:
//   - text: The slice of strings representing the text to split.
//
// Returns:
//   - *TextSplit: A pointer to the created TextSplit instance.
//   - error: An error of type *ErrEmptyText if the input text is empty, or an error
//     of type *ErrWidthTooSmall if the width is less than or equal to 0.
//
// The function calculates the minimum number of lines needed to fit the text within the
// width using the CalculateNumberOfLines function.
// Furthermore, it uses the Sum of Squared Mean (SQM) to find the optimal solution
// for splitting the text into lines of equal width.
//
// If maxHeight is not provided, the function calculates the number of lines needed to fit
// the text within the width using the CalculateNumberOfLines function.
func SplitInEqualSizedLines(text []string, width, height int) (*TextSplit, error) {
	if len(text) == 0 {
		err := uc.NewErrInvalidParameter("text", uc.NewErrEmpty(text))
		return nil, err
	}

	if height == -1 {
		var err error

		height, err = CalculateNumberOfLines(text, width)
		if err != nil {
			return nil, err
		}
	}

	// We have to find the best way to split the text
	// such that all the lines are as close as possible to
	// the average number of characters per line.
	// Example: "Hello World, this is a test"

	// 1. Add each word to a different line.
	// Example:
	//		*** Hello ***
	//		*** World, ***
	//		 *** this ***
	// The rest is out of bounds. (This is not a problem as we will solve it later)

	// 2. If we still have words left, add them at the last line.
	// If the last line exceeds the width, move the first word of the last line
	// to the above line. Do this until all the words fit within the width.
	// Example:
	//		 *** Hello ***
	//		 *** World, ***
	//		*** this is a ***
	// if we were to add "test" to the last line, it would exceed the width.
	// So, we move "this" to the above line, and add "test" to the last line.
	//		   *** Hello ***
	//		*** World, this ***
	//		 *** is a test ***

	group, err := NewTextSplit(width, height)
	uc.AssertErr(err, "NewTextSplit(%d, %d)", width, height)

	for _, word := range text {
		ok := group.InsertWord(word)

		if !ok {
			err := NewErrLinesGreaterThanWords(width, utf8.RuneCountInString(word))
			return nil, err
		}
	}

	// 3. Now we have A solution to the problem, but not THE solution as
	// there may be other ways to split the text that are better than this one.
	// Here is an example where the solution is not optimal:
	//
	// Example: The text "Hi You They" has a Tl of 11 and, assuming
	// W is 8, we have:
	//		$\lceil\frac{11 - 8}{8 + 1}\rceil = \lceil\frac{3}{9}\rceil = 1$
	// This means that the text will be split into two lines:
	//		*** Hi ***
	//		*** You ***
	// With out algorithm, we add "They" to the last line but since it doesn't
	// exceed the width, we don't move any words to the above line.
	//		   *** Hi ***
	//		*** You They ***
	//
	// However, this is not the optimal solution as the following:
	//		*** Hi You ***
	//		 *** They ***
	// is better as the average number of characters per line is closer to the
	// average number of characters per line.

	// We can do this by using SQM (Sum of Squared Mean) as, the lower the SQM,
	// the better the solution.
	// In fact, the optimal solution has an SQM of 1, while our solution has an
	// SQM of 3.

	// 4. Now, we have to find the optimal solution. Because our solution prioritizes
	// the last line, we can do this only by moving words from the last line to the
	// above line; reducing the complexity of the problem.

	// 4.1. For each line that is not the first one, check if the first word of the
	// line can be moved to the above line without exceeding the width.
	// If yes, then it is a candidate for the optimal solution.

	candidates := []*TextSplit{group}

	for _, candidate := range candidates {
		for j := 1; j < height; j++ {
			// Check if the first word of the line can be moved to the above line.
			// If yes, then it is a candidate for the optimal solution.
			ok := candidate.canShiftUp(j)
			if !ok {
				continue
			}

			// Copy the candidate as we don't want to modify the original one.
			candidateCopy := candidate.Copy().(*TextSplit)
			candidateCopy.shiftUp(j)
			candidates = append(candidates, candidateCopy)
		}
	}

	// 4.2. Calculate the SQM of each candidate and returns the ones with the lowest SQM.

	weights := us.ApplyWeightFunc(candidates, calculate_split_ratio)
	if len(weights) == 0 {
		err := NewErrNoCandidateFound()
		return nil, err
	}

	// 4.3. Return the candidates with the lowest SQM.
	weights = us.FilterByNegativeWeight(weights)
	candidates = us.ExtractResults(weights)

	// If we have more than one candidate, we have to choose one
	// of them by following other criteria.
	//
	// (For now, we will just choose the first one.)
	// TODO: Choose the best candidate by following other criteria.

	return candidates[0], nil
}

// LevenshteinTable is a table of words for the Levenshtein distance.
type LavenshteinTable struct {
	// words is the list of words.
	word_list [][]rune

	// word_length_list is the list of word lengths.
	word_length_list []int
}

// NewLevenshteinTable creates a new Levenshtein table
// with the given words.
//
// Parameters:
//   - words: The words to add to the table.
//
// Returns:
//   - *LevenshteinTable: The new Levenshtein table.
func NewLevenshteinTable(words ...string) (*LavenshteinTable, error) {
	lt := &LavenshteinTable{}

	for i, word := range words {
		ok := lt.AddWord(word)
		if !ok {
			return nil, uc.NewErrWhileAt("adding words", i+1, "word", NewErrInvalidUTF8Encoding())
		}
	}

	return lt, nil
}

// AddWord adds a word to the table.
//
// Parameters:
//   - word: The word to add.
//
// Returns:
//   - bool: Whether the word was added successfully.
//
// If this fails, it means that the word could not be converted to UTF-8 runes.
func (lt *LavenshteinTable) AddWord(word string) bool {
	chars, err := ToUTF8Runes(word)
	if err != nil {
		return false
	}

	lt.word_list = append(lt.word_list, chars)
	lt.word_length_list = append(lt.word_length_list, len(chars))

	return true
}

// GetClosest gets the closest word to a target.
//
// Parameters:
//   - target: The target.
//
// Returns:
//   - string: The closest word.
//   - error: The error if any occurs.
//
// Errors:
//   - *common.ErrInvalidParameter: If the target is empty.
//   - *ErrNoClosestWordFound: If no closest word is found.
func (lt *LavenshteinTable) GetClosest(target []rune) (string, error) {
	if len(target) == 0 {
		return "", uc.NewErrInvalidParameter("target", uc.NewErrEmpty(target))
	}

	target_len := len(target)

	closest_idx := -1
	var min int

	for i, word := range lt.word_list {
		d := levenshtein_distance(target, target_len, word, lt.word_length_list[i])

		if closest_idx == -1 || d < min {
			min = d
			closest_idx = i
		}
	}

	if closest_idx == -1 {
		return "", NewErrNoClosestWordFound()
	}

	word := lt.word_list[closest_idx]

	return string(word), nil
}

// levenshteinDistance calculates the Levenshtein distance between two strings.
//
// Parameters:
//   - target: The target.
//   - target_len: The target length.
//   - other: The other.
//   - other_len: The other length.
//
// Returns:
//   - int: The Levenshtein distance.
func levenshtein_distance(target []rune, target_len int, other []rune, other_len int) int {
	matrix := make([][]int, 0, target_len+1)

	for i := 0; i <= target_len; i++ {
		row := make([]int, other_len+1)

		matrix = append(matrix, row)
	}

	// Initialize the matrix
	for i := 0; i <= target_len; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= other_len; j++ {
		matrix[0][j] = j
	}

	// Compute the distances
	for i := 1; i <= target_len; i++ {
		for j := 1; j <= other_len; j++ {
			if target[i-1] == other[j-1] {
				matrix[i][j] = matrix[i-1][j-1] // No operation needed
			} else {
				deletion := matrix[i-1][j] + 1
				insertion := matrix[i][j-1] + 1
				substitution := matrix[i-1][j-1] + 1

				min_first := uc.Min(deletion, insertion)
				min_second := uc.Min(min_first, substitution)
				matrix[i][j] = min_second
			}
		}
	}

	d := matrix[target_len][other_len]

	return d
}

// LastInstanceOfWS finds the last instance of whitespace in the characters.
//
// Parameters:
//   - chars: The characters.
//   - from_idx: The starting index. (inclusive)
//   - to_idx: The ending index. (exclusive)
//
// Returns:
//   - int: The index of the last whitespace character. -1 if not found.
//
// Behaviors:
//   - If from_idx < 0, from_idx is set to 0.
//   - If to_idx >= len(chars), to_idx is set to len(chars) - 1.
//   - If from_idx > to_idx, from_idx and to_idx are swapped.
func LastInstanceOfWS(chars []rune, from_idx, to_idx int) int {
	if len(chars) == 0 {
		return -1
	}

	if from_idx < 0 {
		from_idx = 0
	}

	if to_idx >= len(chars) {
		to_idx = len(chars)
	}

	if from_idx > to_idx {
		from_idx, to_idx = to_idx, from_idx
	}

	for i := to_idx - 1; i >= from_idx; i-- {
		ok := unicode.IsSpace(chars[i])
		if ok {
			return i
		}
	}

	return -1
}

// FirstInstanceOfWS finds the first instance of whitespace in the characters.
//
// Parameters:
//   - chars: The characters.
//   - from_idx: The starting index. (inclusive)
//   - to_idx: The ending index. (exclusive)
//
// Returns:
//   - int: The index of the first whitespace character. -1 if not found.
//
// Behaviors:
//   - If from_idx < 0, from_idx is set to 0.
//   - If to_idx >= len(chars), to_idx is set to len(chars) - 1.
//   - If from_idx > to_idx, from_idx and to_idx are swapped.
//
// FIXME: Remove this function once MyGoLib is updated.
func FirstInstanceOfWS(chars []rune, from_idx, to_idx int) int {
	if len(chars) == 0 {
		return -1
	}

	if from_idx < 0 {
		from_idx = 0
	}

	if to_idx >= len(chars) {
		to_idx = len(chars)
	}

	if from_idx > to_idx {
		from_idx, to_idx = to_idx, from_idx
	}

	for i := from_idx; i < to_idx; i++ {
		ok := unicode.IsSpace(chars[i])
		if ok {
			return i
		}
	}

	return -1
}
