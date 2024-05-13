package SortedMap

import (
	"slices"

	cdp "github.com/PlayerR9/MyGoLib/CustomData/Pair"

	uc "github.com/PlayerR9/MyGoLib/Units/Common"
)

// SortedMap is a generic data structure that represents a sorted map.
type SortedMap[K comparable, V any] struct {
	// mapping is the map of keys to values.
	mapping map[K]V

	// keys is the sorted list of keys.
	keys []K
}

// Copy creates a deep copy of the sorted map.
//
// Returns:
//   - uc.Copier: A deep copy of the sorted map.
func (s *SortedMap[K, V]) Copy() uc.Copier {
	sCopy := &SortedMap[K, V]{
		mapping: make(map[K]V),
		keys:    make([]K, len(s.keys)),
	}

	for _, key := range s.keys {
		sCopy.mapping[key] = s.mapping[key]
	}

	copy(sCopy.keys, s.keys)

	return sCopy
}

// NewSortedMap creates a new sorted map.
//
// Returns:
//   - *SortedMap[K, V]: A pointer to the newly created sorted map.
func NewSortedMap[K comparable, V any]() *SortedMap[K, V] {
	return &SortedMap[K, V]{
		mapping: make(map[K]V),
		keys:    make([]K, 0),
	}
}

// AddEntry adds an entry to the sorted map.
//
// Parameters:
//   - key: The key of the entry.
//   - value: The value of the entry.
//
// Behaviors:
//   - If the key already exists, the value is updated.
func (s *SortedMap[K, V]) AddEntry(key K, value V) {
	_, ok := s.mapping[key]
	if !ok {
		s.keys = append(s.keys, key)
	}

	s.mapping[key] = value
}

// GetEntry gets the value of the entry with the provided key.
//
// Parameters:
//   - key: The key of the entry.
//
// Returns:
//   - V: The value of the entry.
//   - error: An error of type *ErrKeyNotFound if the key does not exist.
func (s *SortedMap[K, V]) GetEntry(key K) (V, error) {
	value, ok := s.mapping[key]

	if !ok {
		return *new(V), NewErrKeyNotFound(key)
	}

	return value, nil
}

// Size returns the number of entries in the sorted map.
//
// Returns:
//   - int: The number of entries in the sorted map.
func (s *SortedMap[K, V]) Size() int {
	return len(s.keys)
}

// Values returns the values of the entries in the sorted map.
//
// Returns:
//   - []V: The values of the entries in the sorted map.
func (s *SortedMap[K, V]) Values() []V {
	values := make([]V, 0, len(s.keys))

	for _, key := range s.keys {
		values = append(values, s.mapping[key])
	}

	return values
}

// Keys returns the keys of the entries in the sorted map.
//
// Returns:
//   - []K: The keys of the entries in the sorted map.
func (s *SortedMap[K, V]) Keys() []K {
	keys := make([]K, len(s.keys))
	copy(keys, s.keys)

	return keys
}

// GetEntries returns the entries in the sorted map.
//
// Returns:
//   - []*cdp.Pair[K, V]: The entries in the sorted map.
//
// Behaviors:
//   - The entries are returned in the order of the keys.
//   - There are no nil pairs in the returned slice.
func (s *SortedMap[K, V]) GetEntries() []*cdp.Pair[K, V] {
	entries := make([]*cdp.Pair[K, V], 0, len(s.keys))

	for _, key := range s.keys {
		entries = append(entries, cdp.NewPair(key, s.mapping[key]))
	}

	return entries
}

// Delete deletes the entry with the provided key from the sorted map.
//
// Parameters:
//   - key: The key of the entry to delete.
//
// Behaviors:
//   - If the key does not exist, nothing happens.
func (s *SortedMap[K, V]) Delete(key K) {
	_, ok := s.mapping[key]
	if !ok {
		return
	}

	delete(s.mapping, key)

	index := slices.Index(s.keys, key)
	if index == -1 {
		return
	}

	s.keys = slices.Delete(s.keys, index, index+1)
}

// ModifyValueFunc is a method that modifies a value of the sorted map.
//
// Parameters:
//   - key: The key of the value to modify.
//   - f: The function that modifies the value.
//
// Returns:
//   - error: An error if the change fails.
//
// Errors:
//   - *ErrKeyNotFound: The key does not exist in the sorted map.
//   - Any error returned by the function 'f'.
func (s *SortedMap[K, V]) ModifyValueFunc(key K, f ModifyValueFunc[V]) error {
	oldValue, ok := s.mapping[key]
	if !ok {
		return NewErrKeyNotFound(key)
	}

	newValue, err := f(oldValue)
	if err != nil {
		return err
	}

	s.mapping[key] = newValue

	return nil
}
