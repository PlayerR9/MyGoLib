package maps

// SeenMap is a map that keeps track of which keys have been seen.
type SeenMap[T comparable] struct {
	// table is the underlying map.
	table map[T]bool
}

// NewSeenMap creates a new SeenMap.
//
// Returns:
//   - *SeenMap[T]: A pointer to the new SeenMap. Never returns nil.
func NewSeenMap[T comparable]() *SeenMap[T] {
	return &SeenMap[T]{
		table: make(map[T]bool),
	}
}

// See marks the key as seen.
//
// Parameters:
//   - key: The key to mark as seen.
func (s *SeenMap[T]) See(key T) {
	s.table[key] = true
}

// IsSeen returns true if the key has been seen.
//
// Parameters:
//   - key: The key to check.
//
// Returns:
//   - bool: True if the key has been seen, false otherwise.
func (s *SeenMap[T]) IsSeen(key T) bool {
	v, ok := s.table[key]
	return ok && v
}

// FilterSeen returns the elements that have not been seen.
//
// Parameters:
//   - elems: The elements to filter.
//
// Returns:
//   - []T: The elements that have not been seen.
func (s *SeenMap[T]) FilterSeen(elems []T) []T {
	var top int

	for i := 0; i < len(elems); i++ {
		elem := elems[i]

		if !s.IsSeen(elem) {
			elems[top] = elems[i]
			top++
		}
	}

	return elems[:top]
}
