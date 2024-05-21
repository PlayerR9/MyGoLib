package Pair

// RemoveNilPairs removes all nil pairs from the given slice of pairs.
//
// Parameters:
//   - pairs: The slice of pairs.
//
// Returns:
//   - []*Pair[A, B]: The slice of pairs without nil pairs.
//
// Behaviors:
//   - If the slice is empty, the function returns nil.
//   - This has the side effect of modifying the original slice when at
//     least one nil pair is found. BEWARE!
func RemoveNilPairs[A any, B any](pairs []*Pair[A, B]) []*Pair[A, B] {
	top := 0

	for i := 0; i < len(pairs); i++ {
		if pairs[i] != nil {
			pairs[top] = pairs[i]
			top++
		}
	}

	if top == 0 {
		return nil
	}

	return pairs[:top]
}

// ExtractFirsts extracts all the first elements from the given slice of pairs.
//
// Parameters:
//   - pairs: The slice of pairs.
//
// Returns:
//   - []A: The slice of first elements.
//
// Behaviors:
//   - If the slice is empty, the function returns nil.
//   - If the slice contains only nil pairs, the function returns nil.
//   - This has the side effect of modifying the original slice when at
//     least one nil pair is found. BEWARE!
func ExtractFirsts[A any, B any](pairs []*Pair[A, B]) []A {
	pairs = RemoveNilPairs(pairs)
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
//
// Behaviors:
//   - If the slice is empty, the function returns nil.
//   - If the slice contains only nil pairs, the function returns nil.
//   - This has the side effect of modifying the original slice when at
//     least one nil pair is found. BEWARE!
func ExtractSeconds[A any, B any](pairs []*Pair[A, B]) []B {
	pairs = RemoveNilPairs(pairs)
	if len(pairs) == 0 {
		return nil
	}

	seconds := make([]B, 0, len(pairs))

	for _, pair := range pairs {
		seconds = append(seconds, pair.Second)
	}

	return seconds
}
