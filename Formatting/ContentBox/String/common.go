package String

import (
	"math"
	"strings"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
	mext "github.com/PlayerR9/MyGoLib/Utility/MathExt"
	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

// CalculateNumberOfLines is a function that calculates the minimum number
// of lines needed to fit a given text within a specified width.
//
// Errors:
//   - *ers.ErrInvalidParameter: If the width is less than or equal to 0.
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
func CalculateNumberOfLines(text []*String, width int) (int, error) {
	if width <= 0 {
		return 0, ers.NewErrInvalidParameter(
			"width",
			ers.NewErrGT(0),
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
	//    For example: "Hello World" has 11 characters. With one split, it becomes
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
		Tl += float64(word.length + 1)
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
//   - *ers.ErrInvalidParameter: If the input text is empty or the width is less than
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
func SplitInEqualSizedLines(text []*String, width, height int) (*TextSplit, error) {
	if err := ers.NewErrEmpty(text).ErrorIf(); err != nil {
		return nil, ers.NewErrInvalidParameter("text", err)
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
		if !group.InsertWord(word) {
			return nil, NewErrLinesGreaterThanWords(width, word.length)
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
			if !candidates[i].canShiftUp(j) {
				continue
			}

			// Copy the candidate as we don't want to modify the original one.
			candidateCopy := candidates[i].Copy().(*TextSplit)
			candidateCopy.shiftUp(j)
			candidates = append(candidates, candidateCopy)
		}
	}

	// 4.2. Calculate the SQM of each candidate and returns the ones with the lowest SQM.
	weights := slext.ApplyWeightFunc(candidates, func(candidate *TextSplit) (float64, bool) {
		values := make([]float64, 0, candidate.GetHeight())

		for _, line := range candidate.lines {
			values = append(values, float64(line.len))
		}

		sqm, err := mext.SQM(values)
		if err != nil {
			return 0, false
		}

		return sqm, true
	})
	if len(weights) == 0 {
		return nil, NewErrNoCandidateFound()
	}

	// 4.3. Return the candidates with the lowest SQM.
	candidates = slext.FilterByNegativeWeight(weights)

	// If we have more than one candidate, we have to choose one
	// of them by following other criteria.
	//
	// (For now, we will just choose the first one.)
	// TODO: Choose the best candidate by following other criteria.

	return candidates[0], nil
}

// Join joins a slice of strings with a separator.
//
// Parameters:
//   - elems: The slice of strings to join.
//   - sep: The separator to join the strings with.
//
// Returns:
//   - *String: The joined string.
//
// Behaviors:
//   - If the slice of strings is empty, a nil string is returned.
//   - The style of the first string is used for the joined string.
//   - The separator is not added to the end of the joined string.
func Join(elems []*String, sep string) *String {
	switch len(elems) {
	case 0:
		return nil
	case 1:
		return elems[0].Copy().(*String)
	default:
		var builder strings.Builder

		builder.WriteString(elems[0].content)

		for _, elem := range elems[1:] {
			builder.WriteString(sep)
			builder.WriteString(elem.content)
		}

		return NewString(builder.String())
	}
}

// FieldsFunc splits a string into fields using a function to determine the separator.
//
// Parameters:
//   - s: The string to split.
//   - sep: The function to determine the separator.
//
// Returns:
//   - []*String: The fields of the string.
func FieldsFunc(s *String, sep string) []*String {
	if sep == "" {
		return []*String{s.Copy().(*String)}
	} else {
		fields := make([]*String, 0)
		var builder strings.Builder

		runes := []rune(sep)
		counter := 0

		for _, r := range s.content {
			if r != runes[counter] {
				counter = 0
				builder.WriteRune(r)

				continue
			}

			counter++

			if counter == len(runes) {
				fields = append(fields, NewString(builder.String()))
				builder.Reset()
				counter = 0
			}
		}

		if builder.Len() > 0 {
			fields = append(fields, NewString(builder.String()))
		}

		return fields
	}
}

// Repeat repeats a string count times.
//
// Parameters:
//   - s: The string to repeat.
//   - count: The number of times to repeat the string.
//
// Returns:
//   - *String: The repeated string.
//
// Behaviors:
//   - If the count is less than or equal to 0, nil is returned.
func Repeat(s *String, count int) *String {
	if count <= 0 {
		return nil
	}

	var builder strings.Builder

	for i := 0; i < count; i++ {
		builder.WriteString(s.content)
	}

	return NewString(builder.String())
}
