package StringExt

import (
	"math"
	"strings"
	"unicode/utf8"

	"crypto/rand"
	"encoding/hex"
	"slices"

	ue "github.com/PlayerR9/MyGoLib/Units/errors"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
	mext "github.com/PlayerR9/MyGoLib/Utility/MathExt"
)

var (
	// calculateSplitRatio is a function that calculates the split ratio of a
	// TextSplit instance.
	calculateSplitRatio us.WeightFunc[*TextSplit]
)

func init() {
	calculateSplitRatio = func(candidate *TextSplit) (float64, bool) {
		values := make([]float64, 0, candidate.GetHeight())

		for _, line := range candidate.lines {
			values = append(values, float64(line.len))
		}

		sqm, ok := mext.SQM(values)

		return sqm, ok
	}
}

// ToUTF8Runes converts a string to a slice of runes.
//
// Parameters:
//   - s: The string to convert.
//
// Returns:
//   - runes: The slice of runes.
//   - error: An error of type *ErrInvalidUTF8Encoding if the string contains
//     invalid UTF-8 encoding.
//
// Behaviors:
//   - An empty string returns a nil slice with no errors.
//   - The function stops at the first invalid UTF-8 encoding; returning an
//     error and the runes found up to that point.
func ToUTF8Runes(s string) (runes []rune, err error) {
	if s == "" {
		return
	}

	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)

		if r == utf8.RuneError {
			err = NewErrInvalidUTF8Encoding()
			return
		}

		runes = append(runes, r)
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

	countStr := utf8.RuneCountInString(str)
	countSuffix := utf8.RuneCountInString(suffix)

	if countStr < countSuffix {
		return "", false
	}

	if countStr == countSuffix {
		return suffix, true
	} else {
		return str[:countStr-countSuffix] + suffix, true
	}
}

// FindContentIndexes searches for the positions of opening and closing
// tokens in a slice of strings.
//
// Parameters:
//   - openingToken: The string that marks the beginning of the content.
//   - closingToken: The string that marks the end of the content.
//   - contentTokens: The slice of strings in which to search for the tokens.
//
// Returns:
//   - result: An array of two integers representing the start and end indexes
//     of the content.
//   - err: Any error that occurred while searching for the tokens.
//
// Errors:
//   - *ue.ErrInvalidParameter: If the openingToken or closingToken is an
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
func FindContentIndexes(openingToken, closingToken string, contentTokens []string) (result [2]int, err error) {
	result[0] = -1
	result[1] = -1

	if openingToken == "" {
		err = ue.NewErrInvalidParameter("openingToken", ue.NewErrEmpty(openingToken))
		return
	} else if closingToken == "" {
		err = ue.NewErrInvalidParameter("closingToken", ue.NewErrEmpty(closingToken))
		return
	}

	openingTokenIndex := slices.Index(contentTokens, openingToken)
	if openingTokenIndex < 0 {
		err = NewErrTokenNotFound(openingToken, OpToken)
		return
	} else {
		result[0] = openingTokenIndex + 1
	}

	tokenBalance := 1
	closingTokenIndex := -1

	for i := result[0]; i < len(contentTokens); i++ {
		if contentTokens[i] == closingToken {
			tokenBalance--

			if tokenBalance == 0 {
				closingTokenIndex = i
				break
			}
		} else if contentTokens[i] == openingToken {
			tokenBalance++
		}
	}

	if closingTokenIndex != -1 {
		result[1] = closingTokenIndex + 1
		return
	}

	if tokenBalance < 0 {
		err = NewErrNeverOpened(openingToken, closingToken)
		return
	} else if tokenBalance != 1 || closingToken != "\n" {
		err = NewErrTokenNotFound(closingToken, ClToken)
		return
	}

	result[1] = len(contentTokens)
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
//   - *ue.ErrInvalidParameter: If the size is less than 1.
//   - Any error returned by the rand.Read function.
//
// Behaviors:
//   - The function uses the crypto/rand package to generate a random ID of the
//     specified size.
//   - The ID is returned as a hexadecimal string.
func GenerateID(size int) (string, error) {
	if size < 1 {
		return "", ue.NewErrInvalidParameter("size", ue.NewErrGT(0))
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
//   - error: An error if the width is less than 0.
//
// Behaviors:
//   - If the width is less than the length of the string, the string is
//     truncated to fit the width.
//   - If the width is greater than the length of the string, spaces are added
//     to the end of the string until the width is reached.
func FitString(s string, width int) (string, error) {
	if width < 0 {
		return "", ue.NewErrInvalidParameter("width", ue.NewErrGTE(0))
	}

	len := len([]rune(s))

	if width == 0 {
		return "", nil
	}

	if len == 0 {
		spacing := strings.Repeat(" ", width)

		return spacing, nil
	}

	if len == width {
		return s, nil
	}

	if len > width {
		return s[:width], nil
	}

	var builder strings.Builder
	spacing := strings.Repeat(" ", width-len)

	builder.WriteString(s)
	builder.WriteString(spacing)

	return builder.String(), nil
}

// CalculateNumberOfLines is a function that calculates the minimum number
// of lines needed to fit a given text within a specified width.
//
// Errors:
//   - *ue.ErrInvalidParameter: If the width is less than or equal to 0.
//   - *ErrLinesGreaterThanWords: If the calculated number of lines is greater
//     than the number of words in the text.
//
// Parameters:
//   - text: The slice of strings representing the text to calculate the number of
//     lines for.
//   - width: The width to fit the text within.
//
// Returns:
//
//   - int: The calculated number of lines needed to fit the text within the width.
//   - error: An error if it occurs during the calculation.
//
// The function calculates the total length of the text (Tl) and uses a mathematical
// formula to estimate the minimum number of lines needed to fit the text within the
// given width. The formula is explained in detail in the comments within the function.
//
// It also returns the calculated number of lines when it errors out
func CalculateNumberOfLines(text []string, width int) (int, error) {
	if width <= 0 {
		return 0, ue.NewErrInvalidParameter(
			"width",
			ue.NewErrGT(0),
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

	numberOfLines := int(math.Ceil((Tl-w)/(w+1))) + 1

	if numberOfLines > len(text) {
		return numberOfLines, NewErrLinesGreaterThanWords(numberOfLines, len(text))
	} else {
		return numberOfLines, nil
	}
}

// SplitInEqualSizedLines is a function that splits a given text into lines of
// equal width.
//
// Errors:
//   - *ue.ErrInvalidParameter: If the input text is empty or the width is less than
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
		return nil, ue.NewErrInvalidParameter("text", ue.NewErrEmpty(text))
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
	if err != nil {
		panic(err)
	}

	for _, word := range text {
		ok := group.InsertWord(word)

		if !ok {
			return nil, NewErrLinesGreaterThanWords(width, utf8.RuneCountInString(word))
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

	for i := 0; i < len(candidates); i++ {
		for j := 1; j < height; j++ {
			// Check if the first word of the line can be moved to the above line.
			// If yes, then it is a candidate for the optimal solution.
			ok := candidates[i].canShiftUp(j)
			if !ok {
				continue
			}

			// Copy the candidate as we don't want to modify the original one.
			candidateCopy := candidates[i].Copy().(*TextSplit)
			candidateCopy.shiftUp(j)
			candidates = append(candidates, candidateCopy)
		}
	}

	// 4.2. Calculate the SQM of each candidate and returns the ones with the lowest SQM.

	weights := us.ApplyWeightFunc(candidates, calculateSplitRatio)
	if len(weights) == 0 {
		return nil, NewErrNoCandidateFound()
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
