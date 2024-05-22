package OrderedMap

import (
	uc "github.com/PlayerR9/MyGoLib/Units/Common"
)

// PredicateFilter is a type that defines a map filter function.
//
// Parameters:
//   - K: The type of the elements in the map.
//
// Returns:
//   - bool: True if the element satisfies the filter function, otherwise false.
type PredicateFilter[K comparable, V any] func(K, V) bool

// Intersect returns a PredicateFilter function that checks if an element
// satisfies all the PredicateFilter functions in funcs.
//
// Parameters:
//   - funcs: A map of PredicateFilter functions.
//
// Returns:
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     all the PredicateFilter functions in funcs.
//
// Behaviors:
//   - It returns false as soon as it finds a function in funcs that the element
//     does not satisfy.
//   - If no functions are provided, it returns a function that always returns true.
func Intersect[K comparable, V any](funcs ...PredicateFilter[K, V]) PredicateFilter[K, V] {
	if len(funcs) == 0 {
		return func(k K, v V) bool { return true }
	}

	return func(k K, v V) bool {
		for _, f := range funcs {
			if !f(k, v) {
				return false
			}
		}

		return true
	}
}

// ParallelIntersect returns a PredicateFilter function that checks if an element
// satisfies all the PredicateFilter functions in funcs.
// It runs the PredicateFilter functions in parallel and returns false as soon as
// it finds a function in funcs that the element does not satisfy.
//
// Parameters:
//   - funcs: A map of PredicateFilter functions.
//
// Returns:
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     all the PredicateFilter functions in funcs.
//
// Behaviors:
//   - It returns false as soon as it finds a function in funcs that the element
//     does not satisfy.
//   - If no functions are provided, it returns a function that always returns true.
func ParallelIntersect[K comparable, V any](funcs ...PredicateFilter[K, V]) PredicateFilter[K, V] {
	if len(funcs) == 0 {
		return func(k K, v V) bool { return true }
	}

	return func(k K, v V) bool {
		resultChan := make(chan bool, len(funcs))

		for _, f := range funcs {
			go func(f PredicateFilter[K, V]) {
				resultChan <- f(k, v)
			}(f)
		}

		for range funcs {
			if !<-resultChan {
				return false
			}
		}

		return true
	}
}

// Union returns a PredicateFilter function that checks if an element
// satisfies at least one of the PredicateFilter functions in funcs.
// It returns true as soon as it finds a function in funcs that the element
// satisfies.
//
// Parameters:
//   - funcs: A map of PredicateFilter functions.
//
// Returns:
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     at least one of the PredicateFilter functions in funcs.
//
// Behaviors:
//   - It returns true as soon as it finds a function in funcs that the element
//     satisfies.
//   - If no functions are provided, it returns a function that always returns false.
func Union[K comparable, V any](funcs ...PredicateFilter[K, V]) PredicateFilter[K, V] {
	if len(funcs) == 0 {
		return func(k K, v V) bool { return false }
	}

	return func(k K, v V) bool {
		for _, f := range funcs {
			if f(k, v) {
				return true
			}
		}

		return false
	}
}

// ParallelUnion returns a PredicateFilter function that checks if an element
// satisfies at least one of the PredicateFilter functions in funcs.
// It runs the PredicateFilter functions in parallel and returns true as soon as
// it finds a function in funcs that the element satisfies.
//
// Parameters:
//   - funcs: A map of PredicateFilter functions.
//
// Returns:
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     at least one of the PredicateFilter functions in funcs.
//
// Behaviors:
//   - It returns true as soon as it finds a function in funcs that the element
//     satisfies.
//   - If no functions are provided, it returns a function that always returns false.
func ParallelUnion[K comparable, V any](funcs ...PredicateFilter[K, V]) PredicateFilter[K, V] {
	if len(funcs) == 0 {
		return func(k K, v V) bool { return false }
	}

	return func(k K, v V) bool {
		resultChan := make(chan bool, len(funcs))

		for _, f := range funcs {
			go func(f PredicateFilter[K, V]) {
				resultChan <- f(k, v)
			}(f)
		}

		for range funcs {
			if <-resultChan {
				return true
			}
		}

		return false
	}
}

// MapFilter is a function that iterates over the map and applies the filter
// function to each element. The returned map contains the elements that satisfy
// the filter function.
//
// Parameters:
//   - S: map of elements.
//   - filter: function that takes an element and returns a bool.
//
// Returns:
//   - map[K]V: map of elements that satisfy the filter function.
//
// Behaviors:
//   - An element is said to satisfy the filter function if the function returns true
//     for that element.
//   - If S is empty, it returns a non-nil empty map.
//   - If filter is nil, it returns S as is
func MapFilter[K comparable, V any](S map[K]V, filter PredicateFilter[K, V]) map[K]V {
	if len(S) == 0 {
		return make(map[K]V)
	} else if filter == nil {
		return S
	}

	result := make(map[K]V)

	for k, v := range S {
		if filter(k, v) {
			result[k] = v
		}
	}

	return result
}

// FilterNilValues is a function that iterates over the map and
// removes elements whose value is nil.
//
// Parameters:
//   - S: map of elements.
//
// Returns:
//   - map[K]*V: map of elements that satisfy the filter function.
//
// Behaviors:
//   - If S is empty, it returns a non-nil empty map.
func FilterNilValues[K comparable, V any](S map[K]*V) map[K]*V {
	if len(S) == 0 {
		return make(map[K]*V)
	}

	result := make(map[K]*V)

	for k, v := range S {
		if v != nil {
			result[k] = v
		}
	}

	return result
}

// SFSeparate is a function that iterates over the map and applies the filter
// function to each element. The returned slices contain the elements that
// satisfy and do not satisfy the filter function.
//
// Parameters:
//   - S: map of elements.
//   - filter: function that takes an element and returns a bool.
//
// Returns:
//   - map[K]V: map of elements that satisfy the filter function.
//   - map[K]V: map of elements that do not satisfy the filter function.
//
// Behavior:
//   - If S is empty, the function returns two empty slices.
func SFSeparate[K comparable, V any](S map[K]V, filter PredicateFilter[K, V]) (map[K]V, map[K]V) {
	success := make(map[K]V)
	failed := make(map[K]V)

	for k, v := range S {
		if filter(k, v) {
			success[k] = v
		} else {
			failed[k] = v
		}
	}

	return success, failed
}

// SFSeparateEarly is a variant of SFSeparate that returns all successful elements.
// If there are none, it returns the original map and false.
//
// Parameters:
//   - S: map of elements.
//   - filter: function that takes an element and returns a bool.
//
// Returns:
//   - map[K]V: map of elements that satisfy the filter function or the original map.
//   - bool: true if there are successful elements, otherwise false.
//
// Behavior:
//   - If S is empty, the function returns an empty map and true.
func SFSeparateEarly[K comparable, V any](S map[K]V, filter PredicateFilter[K, V]) (map[K]V, bool) {
	if len(S) == 0 {
		return map[K]V{}, true
	}

	success := make(map[K]V, 0)
	for k, v := range S {
		if filter(k, v) {
			success[k] = v
		}
	}

	if len(success) == 0 {
		return S, false
	}

	return success, true
}

// RemoveDuplicates removes duplicate elements from the map.
//
// Parameters:
//   - S: map of elements.
//
// Returns:
//   - map[K]V: map of elements with duplicates removed.
//
// Behavior:
//   - The function preserves the order of the elements in the map.
//   - If there are multiple duplicates of an element, only the first
//     occurrence is kept.
//   - If there are less than two elements in the map, the function
//     returns the original map.
func RemoveDuplicates[K comparable, V comparable](S map[K]V) map[K]V {
	if len(S) < 2 {
		return S
	}

	seen := make(map[K]bool)
	unique := make(map[K]V, 0)

	for k, V := range S {
		if _, ok := seen[k]; !ok {
			seen[k] = true
			unique[k] = V
		}
	}

	return unique
}

// RemoveDuplicatesFunc removes duplicate elements from the map.
//
// Parameters:
//   - S: map of elements.
//   - equals: function that takes two elements and returns a bool.
//
// Returns:
//   - map[K]V: map of elements with duplicates removed.
//
// Behavior:
//   - The function preserves the order of the elements in the map.
//   - If there are multiple duplicates of an element, only the first
//     occurrence is kept.
//   - If there are less than two elements in the map, the function
//     returns the original map.
func RemoveDuplicatesFunc[K comparable, V uc.Objecter](S map[K]V) map[K]V {
	if len(S) < 2 {
		return S
	}

	indexMap := make(map[K]V)
	unique := make(map[K]V, 0)

	for k, v := range S {
		found := false

		for _, value := range indexMap {
			if value.Equals(v) {
				found = true
				break
			}
		}

		if !found {
			unique[k] = v

			indexMap[k] = v
		}
	}

	return unique
}

// KeyOfDuplicate returns the index of the first duplicate element in the map.
//
// Parameters:
//   - S: map of elements.
//
// Returns:
//   - K: key of the first duplicate element or a zero value if there are no duplicates.
//   - bool: true if there are duplicates, otherwise false.
func KeyOfDuplicate[K, V comparable](S map[K]V) (K, bool) {
	if len(S) < 2 {
		return *new(K), false
	}

	seen := make(map[V]bool)

	for k, v := range S {
		if _, ok := seen[v]; ok {
			return k, true
		}

		seen[v] = true
	}

	return *new(K), false
}

// KeyOfDuplicateFunc returns the index of the first duplicate element in the map.
//
// Parameters:
//   - S: map of elements.
//
// Returns:
//   - K: key of the first duplicate element or a zero value if there are no duplicates.
//   - bool: true if there are duplicates, otherwise false.
func KeyOfDuplicateFunc[K comparable, V uc.Objecter](S map[K]V) (K, bool) {
	if len(S) < 2 {
		return *new(K), false
	}

	indexMap := make(map[K]V)

	for k, v := range S {
		found := false

		for _, value := range indexMap {
			if value.Equals(v) {
				found = true
				break
			}
		}

		if found {
			return k, true
		}

		indexMap[k] = v
	}

	return *new(K), false
}

// PurgeKeysNotIn removes elements from the map whose keys are not in the validKeys slice.
//
// Parameters:
//   - M: map of elements.
//   - validKeys: slice of keys to keep.
//
// Returns:
//   - map[K]V: map of elements with keys not in validKeys removed.
//
// Behavior:
//   - If validKeys is empty, the function returns an empty map.
//   - If M is empty, the function returns an empty map.
func PurgeKeysNotIn[K comparable, V any](M map[K]V, validKeys []K) map[K]V {
	if len(validKeys) == 0 {
		return make(map[K]V)
	} else if len(M) == 0 {
		return make(map[K]V)
	}

	for _, key := range validKeys {
		_, ok := M[key]
		if !ok {
			delete(M, key)
		}

		if len(M) == 0 {
			break
		}
	}

	return M
}
