package Helpers

import (
	uc "github.com/PlayerR9/MyGoLib/Units/Common"
	up "github.com/PlayerR9/MyGoLib/Units/Pair"
	slext "github.com/PlayerR9/MyGoLib/Units/Slices"
)

// Helper is an interface that represents a helper.
type Helperer[O any] interface {
	// GetData returns the data of the element.
	//
	// Returns:
	//   - *up.Pair[O, error]: The data of the element.
	GetData() *up.Pair[O, error]

	// GetWeight returns the weight of the element.
	//
	// Returns:
	//   - float64: The weight of the element.
	GetWeight() float64
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

// DoIfSuccess executes a function for each successful helper.
//
// Parameters:
//   - S: slice of helpers.
//   - f: the function to execute.
func DoIfSuccess[T Helperer[O], O any](S []T, f uc.DoFunc[O]) {
	if len(S) == 0 {
		return
	}

	for _, h := range S {
		p := h.GetData()

		if p.Second == nil {
			f(p.First)
		}
	}
}

// DoIfFailure executes a function for each failed helper.
//
// Parameters:
//   - S: slice of helpers.
//   - f: the function to execute.
func DoIfFailure[T Helperer[O], O any](S []T, f uc.DualDoFunc[O, error]) {
	if len(S) == 0 {
		return
	}

	for _, h := range S {
		p := h.GetData()

		if p.Second != nil {
			f(p.First, p.Second)
		}
	}
}

// ExtractResults extracts the results from the helpers. Unlike with the GetData
// method, this function returns only the results and not the pair of results and
// errors.
//
// Parameters:
//   - S: slice of helpers.
//
// Returns:
//   - []O: slice of results.
//
// Behaviors:
//   - The results are returned regardless of whether the helper is successful or not.
func ExtractResults[T Helperer[O], O any](S []T) []O {
	if len(S) == 0 {
		return nil
	}

	results := make([]O, 0, len(S))

	for _, h := range S {
		p := h.GetData()

		results = append(results, p.First)
	}

	return results
}

// WeightFunc is a type that defines a function that assigns a weight to an element.
//
// Parameters:
//   - elem: The element to assign a weight to.
//
// Returns:
//   - float64: The weight of the element.
//   - bool: True if the weight is valid, otherwise false.
type WeightFunc[O any] func(elem O) (float64, bool)

// ApplyWeightFunc is a function that iterates over the slice and applies the weight
// function to each element.
//
// Parameters:
//   - S: slice of elements.
//   - f: the weight function.
//
// Returns:
//   - []WeightResult[O]: slice of elements with their corresponding weights.
//
// Behaviors:
//   - If S is empty or f is nil, the function returns nil.
//   - If the weight function returns false, the element is not included in the result.
func ApplyWeightFunc[O any](S []O, f WeightFunc[O]) []*WeightedElement[O] {
	if len(S) == 0 || f == nil {
		return nil
	}

	trimmed := make([]*WeightedElement[O], 0)

	for _, e := range S {
		weight, ok := f(e)
		if !ok {
			continue
		}

		trimmed = append(trimmed, NewWeightedElement(e, weight))
	}

	return trimmed
}

// MaxSuccessOrFail returns the results with the maximum weight.
//
// Parameters:
//   - batch: The slice of results.
//
// Returns:
//   - []*up.Pair[O, error]: The results with the maximum weight.
//   - bool: True if the slice was filtered, false otherwise.
//
// Behaviors:
//   - If the slice is empty, the function returns a nil slice and true.
//   - The result can either be the sucessful results or the original slice.
//     Nonetheless, the maximum weight is always applied.
func MaxSuccessOrFail[T Helperer[O], O any](batch []T) ([]T, bool) {
	// 1. Remove nil elements.
	if len(batch) == 0 {
		return nil, true
	}

	// 2. Either get successful results or return the original slice.
	batch, ok := slext.SFSeparateEarly(batch, FilterIsNotSuccess[T, O])

	// 3. Get only the results with the maximum weight.
	solution := FilterByPositiveWeight(batch)

	return solution, ok
}

// MinSuccessOrFail returns the results with the minimum weight.
//
// Parameters:
//   - batch: The slice of results.
//
// Returns:
//   - []*up.Pair[O, error]: The results with the minimum weight.
//   - bool: True if the slice was filtered, false otherwise.
func MinSuccessOrFail[T Helperer[O], O any](batch []T) ([]T, bool) {
	if len(batch) == 0 {
		return nil, true
	}

	// 2. Either get successful results or return the original slice.
	batch, ok := slext.SFSeparateEarly(batch, FilterIsNotSuccess[T, O])

	// 3. Get only the results with the minimum weight.
	solution := FilterByNegativeWeight(batch)

	return solution, ok
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

	var allSolutions []*SimpleHelper[O]

	for _, h := range batch {
		res, err := f(h)

		helper := NewSimpleHelper(res, err)

		allSolutions = append(allSolutions, helper)
	}

	top := 0

	for i := 0; i < len(allSolutions); i++ {
		h := allSolutions[i]

		if h.reason == nil {
			allSolutions[top] = h
			top++
		}
	}

	if top != 0 {
		return allSolutions[:top], true
	} else {
		return allSolutions, false
	}
}

// EvaluateWeightHelpers is a function that evaluates a batch of helpers and returns
// the results.
//
// Parameters:
//   - batch: The slice of helpers.
//   - f: The evaluation function.
//   - wf: The weight function.
//   - isPositive: True if the weight is positive, false otherwise.
//
// Returns:
//   - []*WeightedHelper[O]: The results of the evaluation.
//   - bool: True if the slice was filtered, false otherwise.
//
// Behaviors:
//   - This function returns either the successful results or the original slice.
func EvaluateWeightHelpers[T any, O any](batch []T, f uc.EvalOneFunc[T, O], wf WeightFunc[T], isPositive bool) ([]*WeightedHelper[O], bool) {
	if len(batch) == 0 {
		return nil, true
	}

	var allSolutions []*WeightedHelper[O]

	for _, h := range batch {
		res, err := f(h)

		weight, ok := wf(h)
		if !ok {
			continue
		}

		h := NewWeightedHelper(res, err, weight)

		allSolutions = append(allSolutions, h)
	}

	top := 0

	for i := 0; i < len(allSolutions); i++ {
		h := allSolutions[i]

		if h.reason == nil {
			allSolutions[top] = h
			top++
		}
	}

	if top != 0 {
		allSolutions = allSolutions[:top]
	}

	if isPositive {
		allSolutions = FilterByPositiveWeight(allSolutions)
	} else {
		allSolutions = FilterByNegativeWeight(allSolutions)
	}

	return allSolutions, top > 0
}
