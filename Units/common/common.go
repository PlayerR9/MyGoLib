package common

// ExtractFirsts extracts all the first elements from the given slice of pairs.
//
// Parameters:
//   - pairs: The slice of pairs.
//
// Returns:
//   - []A: The slice of first elements.
func ExtractFirsts[A any, B any](pairs []Pair[A, B]) []A {
	if len(pairs) == 0 {
		return nil
	}

	firsts := make([]A, 0, len(pairs))

	for _, pair := range pairs {
		firsts = append(firsts, pair.First)
	}

	return firsts
}

// ExtractSeconds extracts all the second elements from the given slice of pairs.
//
// Parameters:
//   - pairs: The slice of pairs.
//
// Returns:
//   - []B: The slice of second elements.
func ExtractSeconds[A any, B any](pairs []Pair[A, B]) []B {
	if len(pairs) == 0 {
		return nil
	}

	seconds := make([]B, 0, len(pairs))

	for _, pair := range pairs {
		seconds = append(seconds, pair.Second)
	}

	return seconds
}
