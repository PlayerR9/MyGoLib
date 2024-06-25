package OrderedMap

import (
	"slices"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	uts "github.com/PlayerR9/MyGoLib/Utility/Sorting"
)

// MapKeyer is an interface that represents a key of a map.
type MapKeyer interface {
	uc.Copier
}

// LessMap is a generic data structure that represents a sorted map.
type LessMap[K MapKeyer, V any] struct {
	keys   []K
	values []V

	sf uts.SortFunc[K]
}

// Copy creates a deep copy of the sorted map.
//
// Returns:
//   - uc.Copier: A deep copy of the sorted map.
func (s *LessMap[K, V]) Copy() uc.Copier {
	keysCopy := make([]K, 0, len(s.keys))
	for _, key := range s.keys {
		keysCopy = append(keysCopy, key.Copy().(K))
	}

	valuesCopy := make([]V, 0, len(s.values))
	for _, value := range s.values {
		valuesCopy = append(valuesCopy, uc.CopyOf(value).(V))
	}

	return &LessMap[K, V]{
		keys:   keysCopy,
		values: valuesCopy,
		sf:     s.sf,
	}
}

// NewLessMap creates a new sorted map.
//
// Parameters:
//   - sf: The sort function to use for sorting the keys.
//
// Returns:
//   - *LessMap[K, V]: A pointer to the newly created sorted map.
//
// Behaviors:
//   - If the sort function is nil, the function returns nil.
func NewLessMap[K MapKeyer, V any](sf uts.SortFunc[K]) *LessMap[K, V] {
	if sf == nil {
		return nil
	}

	return &LessMap[K, V]{
		keys:   make([]K, 0),
		values: make([]V, 0),
		sf:     sf,
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
func (s *LessMap[K, V]) AddEntry(key K, value V) {
	pos, ok := slices.BinarySearchFunc(s.keys, key, s.sf)
	if ok {
		s.values[pos] = value
	} else {
		s.keys = slices.Insert(s.keys, pos, key)
		s.values = slices.Insert(s.values, pos, value)
	}
}

// GetEntry gets the value of the entry with the provided key.
//
// Parameters:
//   - key: The key of the entry.
//
// Returns:
//   - V: The value of the entry.
//   - bool: A boolean indicating if the key exists in the sorted map.
//
// Errors:
//   - One can use the error *ErrKeyNotFound from the package
//     when the key does not exist.
func (s *LessMap[K, V]) GetEntry(key K) (V, bool) {
	pos, ok := slices.BinarySearchFunc(s.keys, key, s.sf)
	if !ok {
		return *new(V), false
	}

	return s.values[pos], true
}

// Size returns the number of entries in the sorted map.
//
// Returns:
//   - int: The number of entries in the sorted map.
func (s *LessMap[K, V]) Size() int {
	return len(s.keys)
}

// Values returns the values of the entries in the sorted map.
//
// Returns:
//   - []V: The values of the entries in the sorted map.
func (s *LessMap[K, V]) Values() []V {
	return s.values
}

// Keys returns the keys of the entries in the sorted map.
//
// Returns:
//   - []K: The keys of the entries in the sorted map.
func (s *LessMap[K, V]) Keys() []K {
	return s.keys
}

// GetEntries returns the entries in the sorted map.
//
// Returns:
//   - []*uc.Pair[K, V]: The entries in the sorted map.
//
// Behaviors:
//   - The entries are returned in the order of the keys.
//   - There are no nil pairs in the returned slice.
//   - Prefer using Iterator() method for iterating over the entries
//     instead of this method.
func (s *LessMap[K, V]) GetEntries() []uc.Pair[K, V] {
	entries := make([]uc.Pair[K, V], 0, len(s.keys))

	for i, key := range s.keys {
		p := uc.NewPair(key, s.values[i])

		entries = append(entries, p)
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
func (s *LessMap[K, V]) Delete(key K) {
	pos, ok := slices.BinarySearchFunc(s.keys, key, s.sf)
	if !ok {
		return
	}

	s.keys = slices.Delete(s.keys, pos, pos+1)
	s.values = slices.Delete(s.values, pos, pos+1)
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
func (s *LessMap[K, V]) ModifyValueFunc(key K, f ModifyValueFunc[V]) error {
	pos, ok := slices.BinarySearchFunc(s.keys, key, s.sf)
	if !ok {
		return NewErrKeyNotFound()
	}

	newValue, err := f(s.values[pos])
	if err != nil {
		return err
	}

	s.values[pos] = newValue

	return nil
}

// Iterator returns an iterator for the sorted map.
//
// Returns:
//   - uc.Iterater[*uc.Pair[K, V]]: An iterator for the sorted map.
//
// Behaviors:
//   - The iterator returns the entries in the order of the keys as pairs.
func (s *LessMap[K, V]) Iterator() uc.Iterater[uc.Pair[K, V]] {
	var builder uc.Builder[uc.Pair[K, V]]

	for i, key := range s.keys {
		p := uc.NewPair(key, s.values[i])

		builder.Add(p)
	}

	return builder.Build()
}

// DoFunc performs a function on each entry in the sorted map.
//
// Parameters:
//   - f: The function to perform on each entry.
//
// Returns:
//   - error: An error if the function fails.
//
// Behaviors:
//   - The function 'f' is called for each entry in the sorted map.
//   - If the function 'f' returns an error, the iteration stops.
func (s *LessMap[K, V]) DoFunc(f func(K, V) error) error {
	for i, key := range s.keys {
		err := f(key, s.values[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// GetAt returns the value at the provided index.
//
// Parameters:
//   - index: The index of the value to retrieve.
//
// Returns:
//   - V: The value at the provided index.
//
// Errors:
//   - *uc.ErrInvalidParameter: The index is out of bounds.
func (s *LessMap[K, V]) GetAt(index int) (V, error) {
	if index < 0 || index >= len(s.keys) {
		return *new(V), uc.NewErrInvalidParameter(
			"index",
			uc.NewErrOutOfBounds(index, 0, len(s.keys)),
		)
	}

	return s.values[index], nil
}
