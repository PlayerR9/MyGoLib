package slice

import uc "github.com/PlayerR9/MyGoLib/Units/common"

// PredicateFilter is a type that defines a slice filter function.
//
// Parameters:
//   - T: The type of the elements in the slice.
//
// Returns:
//   - bool: True if the element satisfies the filter function, otherwise false.
type PredicateFilter[T any] func(T) bool

// Intersect returns a PredicateFilter function that checks if an element
// satisfies all the PredicateFilter functions in funcs.
//
// Parameters:
//   - funcs: A slice of PredicateFilter functions.
//
// Returns:
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     all the PredicateFilter functions in funcs.
//
// Behavior:
//   - If no filter functions are provided, then all elements are considered to satisfy
//     the filter function.
//   - It returns false as soon as it finds a function in funcs that the element
//     does not satisfy.
func Intersect[T any](funcs ...PredicateFilter[T]) PredicateFilter[T] {
	if len(funcs) == 0 {
		return func(elem T) bool { return true }
	}

	return func(elem T) bool {
		for _, f := range funcs {
			if !f(elem) {
				return false
			}
		}

		return true
	}
}

// ParallelIntersect returns a PredicateFilter function that checks if an element
// satisfies all the PredicateFilter functions in funcs concurrently.
//
// Parameters:
//   - funcs: A slice of PredicateFilter functions.
//
// Returns:
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     all the PredicateFilter functions in funcs.
//
// Behavior:
//   - If no filter functions are provided, then all elements are considered to satisfy
//     the filter function.
//   - It returns false as soon as it finds a function in funcs that the element
//     does not satisfy.
func ParallelIntersect[T any](funcs ...PredicateFilter[T]) PredicateFilter[T] {
	if len(funcs) == 0 {
		return func(elem T) bool { return true }
	}

	return func(elem T) bool {
		resultChan := make(chan bool, len(funcs))

		for _, f := range funcs {
			go func(f PredicateFilter[T]) {
				resultChan <- f(elem)
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
//
// Parameters:
//   - funcs: A slice of PredicateFilter functions.
//
// Returns:
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     at least one of the PredicateFilter functions in funcs.
//
// Behavior:
//   - If no filter functions are provided, then no elements are considered to satisfy
//     the filter function.
//   - It returns true as soon as it finds a function in funcs that the element
//     satisfies.
func Union[T any](funcs ...PredicateFilter[T]) PredicateFilter[T] {
	if len(funcs) == 0 {
		return func(elem T) bool { return false }
	}

	return func(elem T) bool {
		for _, f := range funcs {
			if f(elem) {
				return true
			}
		}

		return false
	}
}

// ParallelUnion returns a PredicateFilter function that checks if an element
// satisfies at least one of the PredicateFilter functions in funcs concurrently.
//
// Parameters:
//   - funcs: A slice of PredicateFilter functions.
//
// Returns:
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     at least one of the PredicateFilter functions in funcs.
//
// Behavior:
//   - If no filter functions are provided, then no elements are considered to satisfy
//     the filter function.
//   - It returns true as soon as it finds a function in funcs that the element
//     satisfies.
func ParallelUnion[T any](funcs ...PredicateFilter[T]) PredicateFilter[T] {
	if len(funcs) == 0 {
		return func(elem T) bool { return false }
	}

	return func(elem T) bool {
		resultChan := make(chan bool, len(funcs))

		for _, f := range funcs {
			go func(f PredicateFilter[T]) {
				resultChan <- f(elem)
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

// SliceFilter is a function that iterates over the slice and applies the filter
// function to each element.
//
// Parameters:
//   - S: slice of elements.
//   - filter: function that takes an element and returns a bool.
//
// Returns:
//   - []T: slice of elements that satisfy the filter function.
//
// Behavior:
//   - If S is empty, the function returns a nil slice.
//   - If S has only one element and it satisfies the filter function, the function
//     returns a slice with that element. Otherwise, it returns a nil slice.
//   - An element is said to satisfy the filter function if the function returns true
//     when applied to the element.
//   - If the filter function is nil, the function returns the original slice.
func SliceFilter[T any](S []T, filter PredicateFilter[T]) []T {
	if len(S) == 0 {
		return nil
	} else if filter == nil {
		return S
	}

	top := 0

	for i := 0; i < len(S); i++ {
		if filter(S[i]) {
			S[top] = S[i]
			top++
		}
	}

	return S[:top]
}

// FilterNilValues is a function that iterates over the slice and removes the
// nil elements.
//
// Parameters:
//   - S: slice of elements.
//
// Returns:
//   - []*T: slice of elements that satisfy the filter function.
//
// Behavior:
//   - If S is empty, the function returns a nil slice.
func FilterNilValues[T any](S []*T) []*T {
	if len(S) == 0 {
		return nil
	}

	top := 0

	for i := 0; i < len(S); i++ {
		if S[i] != nil {
			S[top] = S[i]
			top++
		}
	}

	return S[:top]
}

// FilterNilPredicates is a function that iterates over the slice and removes the
// nil predicate functions.
//
// Parameters:
//   - S: slice of predicate functions.
//
// Returns:
//   - []PredicateFilter: slice of predicate functions that are not nil.
//
// Behavior:
//   - If S is empty, the function returns a nil slice.
func FilterNilPredicates[T any](S []PredicateFilter[T]) []PredicateFilter[T] {
	if len(S) == 0 {
		return nil
	}

	top := 0

	for i := 0; i < len(S); i++ {
		if S[i] != nil {
			S[top] = S[i]
			top++
		}
	}

	return S[:top]
}

// SFSeparate is a function that iterates over the slice and applies the filter
// function to each element. The returned slices contain the elements that
// satisfy and do not satisfy the filter function.
//
// Parameters:
//   - S: slice of elements.
//   - filter: function that takes an element and returns a bool.
//
// Returns:
//   - []T: slice of elements that satisfy the filter function.
//   - []T: slice of elements that do not satisfy the filter function.
//
// Behavior:
//   - If S is empty, the function returns two empty slices.
func SFSeparate[T any](S []T, filter PredicateFilter[T]) ([]T, []T) {
	if len(S) == 0 {
		return []T{}, []T{}
	}

	var failed []T

	top := 0

	for i := 0; i < len(S); i++ {
		if filter(S[i]) {
			S[top] = S[i]
			top++
		} else {
			failed = append(failed, S[i])
		}
	}

	return S[:top], failed
}

// SFSeparateEarly is a variant of SFSeparate that returns all successful elements.
// If there are none, it returns the original slice and false.
//
// Parameters:
//   - S: slice of elements.
//   - filter: function that takes an element and returns a bool.
//
// Returns:
//   - []T: slice of elements that satisfy the filter function or the original slice.
//   - bool: true if there are successful elements, otherwise false.
//
// Behavior:
//   - If S is empty, the function returns an empty slice and true.
func SFSeparateEarly[T any](S []T, filter PredicateFilter[T]) ([]T, bool) {
	if len(S) == 0 {
		return []T{}, true
	}

	top := 0

	for i := 0; i < len(S); i++ {
		if filter(S[i]) {
			S[top] = S[i]
			top++
		}
	}

	if top == 0 {
		return S, false
	} else {
		return S[:top], true
	}
}

// FilterIsSuccess filters any helper that is not successful.
//
// Parameters:
//   - h: The helper to filter.
//
// Returns:
//   - bool: True if the helper is successful, false otherwise.
//
// Behaviors:
//   - It assumes that the h is not nil.
func FilterIsSuccess[T Helperer[O], O any](h T) bool {
	return h.GetData().Second == nil
}

// FilterByPositiveWeight is a function that iterates over weight results and
// returns the elements with the maximum weight.
//
// Parameters:
//   - S: slice of weight results.
//
// Returns:
//   - []T: slice of elements with the maximum weight.
//
// Behaviors:
//   - If S is empty, the function returns a nil slice.
//   - If multiple elements have the same maximum weight, they are all returned.
//   - If S contains only one element, that element is returned.
func FilterByPositiveWeight[T Helperer[O], O any](S []T) []T {
	if len(S) == 0 {
		return nil
	}

	maxWeight := S[0].GetWeight()
	indices := []int{0}

	for i, e := range S[1:] {
		currentWeight := e.GetWeight()

		if currentWeight > maxWeight {
			maxWeight = currentWeight
			indices = []int{i + 1}
		} else if currentWeight == maxWeight {
			indices = append(indices, i+1)
		}
	}

	solution := make([]T, len(indices))

	for i, index := range indices {
		solution[i] = S[index]
	}

	return solution
}

// FilterByNegativeWeight is a function that iterates over weight results and
// returns the elements with the minimum weight.
//
// Parameters:
//   - S: slice of weight results.
//
// Returns:
//   - []T: slice of elements with the minimum weight.
//
// Behaviors:
//   - If S is empty, the function returns a nil slice.
//   - If multiple elements have the same minimum weight, they are all returned.
//   - If S contains only one element, that element is returned.
func FilterByNegativeWeight[T Helperer[O], O any](S []T) []T {
	if len(S) == 0 {
		return nil
	}

	minWeight := S[0].GetWeight()
	indices := []int{0}

	for i, e := range S[1:] {
		currentWeight := e.GetWeight()

		if currentWeight < minWeight {
			minWeight = currentWeight
			indices = []int{i + 1}
		} else if currentWeight == minWeight {
			indices = append(indices, i+1)
		}
	}

	solution := make([]T, len(indices))
	for i, index := range indices {
		solution[i] = S[index]
	}

	return solution
}

// SuccessOrFail returns the results with the maximum weight.
//
// Parameters:
//   - batch: The slice of results.
//   - useMax: True if the maximum weight should be used, false otherwise.
//
// Returns:
//   - []*uc.Pair[O, error]: The results with the maximum weight.
//   - bool: True if the slice was filtered, false otherwise.
//
// Behaviors:
//   - If the slice is empty, the function returns a nil slice and true.
//   - The result can either be the sucessful results or the original slice.
//     Nonetheless, the maximum weight is always applied.
func SuccessOrFail[T Helperer[O], O any](batch []T, useMax bool) ([]T, bool) {
	// 1. Remove nil elements.
	if len(batch) == 0 {
		return nil, true
	}

	success, fail := SFSeparate(batch, FilterIsSuccess[T, O])

	var target, solution []T

	if len(success) == 0 {
		target = fail
	} else {
		target = success
	}

	if useMax {
		solution = FilterByPositiveWeight(target)
	} else {
		solution = FilterByNegativeWeight(target)
	}

	return solution, len(success) > 0
}

// EvaluateSimpleHelpers is a function that evaluates a batch of helpers and returns
// the results.
//
// Parameters:
//   - batch: The slice of helpers.
//   - f: The evaluation function.
//
// Returns:
//   - []*SimpleHelper[O]: The results of the evaluation.
//   - bool: True if the slice was filtered, false otherwise.
//
// Behaviors:
//   - This function returns either the successful results or the original slice.
func EvaluateSimpleHelpers[T any, O any](batch []T, f uc.EvalOneFunc[T, O]) ([]*SimpleHelper[O], bool) {
	if len(batch) == 0 {
		return nil, true
	}

	solutions := make([]*SimpleHelper[O], 0, len(batch))

	for _, h := range batch {
		res, err := f(h)

		helper := NewSimpleHelper(res, err)
		solutions = append(solutions, helper)
	}

	success, fail := SFSeparate(solutions, FilterIsSuccess)

	var result []*SimpleHelper[O]

	if len(success) == 0 {
		result = fail
	} else {
		result = success
	}

	return result, len(success) > 0
}

// EvaluateWeightHelpers is a function that evaluates a batch of helpers and returns
// the results.
//
// Parameters:
//   - batch: The slice of helpers.
//   - f: The evaluation function.
//   - wf: The weight function.
//   - useMax: True if the maximum weight should be used, false otherwise.
//
// Returns:
//   - []*WeightedHelper[O]: The results of the evaluation.
//   - bool: True if the slice was filtered, false otherwise.
//
// Behaviors:
//   - This function returns either the successful results or the original slice.
func EvaluateWeightHelpers[T any, O any](batch []T, f uc.EvalOneFunc[T, O], wf WeightFunc[T], useMax bool) ([]*WeightedHelper[O], bool) {
	if len(batch) == 0 {
		return nil, true
	}

	solutions := make([]*WeightedHelper[O], 0, len(batch))

	for _, h := range batch {
		res, err := f(h)

		weight, ok := wf(h)
		if !ok {
			continue
		}

		h := NewWeightedHelper(res, err, weight)
		solutions = append(solutions, h)
	}

	success, fail := SFSeparate(solutions, FilterIsSuccess)

	var target, result []*WeightedHelper[O]

	if len(success) == 0 {
		target = fail
	} else {
		target = success
	}

	if useMax {
		result = FilterByPositiveWeight(target)
	} else {
		result = FilterByNegativeWeight(target)
	}
	return result, len(success) > 0
}

// RemoveEmpty is a function that removes the empty elements from a slice.
//
// Parameters:
//   - elems: The slice of elements.
//
// Returns:
//   - []T: The slice of elements without the empty elements.
func RemoveEmpty[T comparable](elems []T) []T {
	top := 0

	for i := 0; i < len(elems); i++ {
		if elems[i] != *new(T) {
			elems[top] = elems[i]
			top++
		}
	}

	return elems[:top]
}
