package OrderedMap

import (
	"slices"
	"strings"

	uc "github.com/PlayerR9/lib_units/common"
	lup "github.com/PlayerR9/lib_units/pair"
	lustr "github.com/PlayerR9/lib_units/strings"
)

// OrderedMap is a generic data structure that represents a sorted map.
type OrderedMap[K comparable, V any] struct {
	// mapping is the map of keys to values.
	mapping map[K]V

	// keys is the sorted list of keys.
	keys []K
}

/* // Equals is a method that checks if two OrderedMaps are equal.
//
// Returns:
//   - bool: true if the maps are equal, false otherwise.
func (s *OrderedMap[K, V]) Equals(other *OrderedMap[K, V]) bool {
	if other == nil {
		return false
	}

	if len(s.keys) != len(other.keys) {
		return false
	}

	for key, value := range s.mapping {
		val, ok := other.mapping[key]
		if !ok {
			return false
		}

		if !uc.EqualOf(value, val) {
			return false
		}
	}

	return true
} */

// String implements common.Objecter.
func (s *OrderedMap[K, V]) String() string {
	if len(s.keys) == 0 {
		return "{}"
	} else if len(s.keys) == 1 {
		var builder strings.Builder

		builder.WriteRune('{')
		builder.WriteString(lustr.GoStringOf(s.keys[0]))
		builder.WriteString(" : ")
		builder.WriteString(lustr.GoStringOf(s.mapping[s.keys[0]]))
		builder.WriteRune('}')

		return builder.String()
	}

	var builder strings.Builder

	builder.WriteRune('{')
	builder.WriteString(lustr.GoStringOf(s.keys[0]))
	builder.WriteString(" : ")
	builder.WriteString(lustr.GoStringOf(s.mapping[s.keys[0]]))

	for _, key := range s.keys[1:] {
		builder.WriteString(", ")
		builder.WriteString(lustr.GoStringOf(key))
		builder.WriteString(" : ")
		builder.WriteString(lustr.GoStringOf(s.mapping[key]))
	}

	builder.WriteRune('}')

	return builder.String()
}

// Copy creates a shallow copy of the sorted map.
//
// Returns:
//   - *SortedMap[K, V]: A shallow copy of the sorted map.
func (s *OrderedMap[K, V]) Copy() *OrderedMap[K, V] {
	sCopy := &OrderedMap[K, V]{
		mapping: make(map[K]V),
		keys:    make([]K, len(s.keys)),
	}

	for _, key := range s.keys {
		sCopy.mapping[key] = s.mapping[key]
	}

	copy(sCopy.keys, s.keys)

	return sCopy
}

// NewOrderedMap creates a new sorted map.
//
// Returns:
//   - *SortedMap[K, V]: A pointer to the newly created sorted map.
func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
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
func (s *OrderedMap[K, V]) AddEntry(key K, value V) {
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
//   - bool: A boolean indicating if the key exists in the sorted map.
//
// Errors:
//   - One can use the error *ErrKeyNotFound from the package
//     when the key does not exist.
func (s *OrderedMap[K, V]) GetEntry(key K) (V, bool) {
	value, ok := s.mapping[key]

	return value, ok
}

// Size returns the number of entries in the sorted map.
//
// Returns:
//   - int: The number of entries in the sorted map.
func (s *OrderedMap[K, V]) Size() int {
	return len(s.keys)
}

// Values returns the values of the entries in the sorted map.
//
// Returns:
//   - []V: The values of the entries in the sorted map.
func (s *OrderedMap[K, V]) Values() []V {
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
func (s *OrderedMap[K, V]) Keys() []K {
	keys := make([]K, len(s.keys))
	copy(keys, s.keys)

	return keys
}

// GetEntries returns the entries in the sorted map.
//
// Returns:
//   - []*lup.Pair[K, V]: The entries in the sorted map.
//
// Behaviors:
//   - The entries are returned in the order of the keys.
//   - There are no nil pairs in the returned slice.
//   - Prefer using Iterator() method for iterating over the entries
//     instead of this method.
func (s *OrderedMap[K, V]) GetEntries() []lup.Pair[K, V] {
	entries := make([]lup.Pair[K, V], 0, len(s.keys))

	for _, key := range s.keys {
		entries = append(entries, lup.NewPair(key, s.mapping[key]))
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
func (s *OrderedMap[K, V]) Delete(key K) {
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
func (s *OrderedMap[K, V]) ModifyValueFunc(key K, f ModifyValueFunc[V]) error {
	oldValue, ok := s.mapping[key]
	if !ok {
		return NewErrKeyNotFound()
	}

	newValue, err := f(oldValue)
	if err != nil {
		return err
	}

	s.mapping[key] = newValue

	return nil
}

// SortKeys sorts the keys of the sorted map.
//
// Parameters:
//   - less: The function that defines the sorting order.
//
// Behaviors:
//   - The keys are sorted in place using the slice.SortFunc function.
//   - The function 'less' should return < 0 if the first key is less than the second key.
//   - The function 'less' should return > 0 if the first key is greater than the second key.
//   - The function 'less' should return 0 if the first key is equal to the second key.
func (s *OrderedMap[K, V]) SortKeys(less func(K, K) int) {
	slices.SortFunc(s.keys, less)
}

// Iterator returns an iterator for the sorted map.
//
// Returns:
//   - uc.Iterater[*lup.Pair[K, V]]: An iterator for the sorted map.
//
// Behaviors:
//   - The iterator returns the entries in the order of the keys as pairs.
func (s *OrderedMap[K, V]) Iterator() uc.Iterater[lup.Pair[K, V]] {
	var builder uc.Builder[lup.Pair[K, V]]

	for _, key := range s.keys {
		builder.Add(lup.NewPair(key, s.mapping[key]))
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
func (s *OrderedMap[K, V]) DoFunc(f func(K, V) error) error {
	for _, key := range s.keys {
		if err := f(key, s.mapping[key]); err != nil {
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
func (s *OrderedMap[K, V]) GetAt(index int) (V, error) {
	if index < 0 || index >= len(s.keys) {
		return *new(V), uc.NewErrInvalidParameter(
			"index",
			uc.NewErrOutOfBounds(index, 0, len(s.keys)),
		)
	}

	return s.mapping[s.keys[index]], nil
}
