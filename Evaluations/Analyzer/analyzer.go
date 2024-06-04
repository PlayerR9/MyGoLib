package Analyzer

import (
	"slices"
)

// FrequencyEvaluator is a structure that provides methods to evaluate the frequency of elements.
type FrequencyEvaluator[T comparable] struct {
	// elements is the map that stores the frequency of each element.
	elements map[T]int

	// keys is the slice that stores the keys of the elements.
	keys []T

	// total is the total number of elements.
	total float64
}

// NewFrequencyEvaluator creates a new evaluator to evaluate the frequency of elements.
//
// Returns:
//   - *FrequencyEvaluator[T]: The new evaluator.
func NewFrequencyEvaluator[T comparable]() *FrequencyEvaluator[T] {
	return &FrequencyEvaluator[T]{
		elements: make(map[T]int),
		keys:     make([]T, 0),
		total:    0,
	}
}

// AddElement adds an element to the evaluator.
//
// Parameters:
//   - element: The element to add.
func (fe *FrequencyEvaluator[T]) AddElement(element T) {
	_, ok := fe.elements[element]
	if !ok {
		fe.keys = append(fe.keys, element)
	}

	fe.elements[element]++
	fe.total++
}

// GetFrequency returns the frequency of each element.
//
// Returns:
//   - map[T]int: The frequency of each element.
func (fe *FrequencyEvaluator[T]) GetFrequency() map[T]int {
	return fe.elements
}

// GetInPercent returns the frequency of each element in percent.
//
// Returns:
//   - map[T]float64: The frequency of each element in percent.
func (fe *FrequencyEvaluator[T]) GetInPercent() map[T]float64 {
	prcnt := make(map[T]float64)

	for _, k := range fe.keys {
		prcnt[k] = float64(fe.elements[k]) / fe.total
	}

	return prcnt
}

// GetTotal returns the total number of elements.
//
// Returns:
//   - float64: The total number of elements.
func (fe *FrequencyEvaluator[T]) GetTotal() float64 {
	return fe.total
}

// GetKeys returns the keys of the elements.
//
// Returns:
//   - []T: The keys of the elements.
func (fe *FrequencyEvaluator[T]) GetKeys() []T {
	return fe.keys
}

// Reset resets the evaluator.
func (fe *FrequencyEvaluator[T]) Reset() {
	fe.elements = make(map[T]int)
	fe.keys = make([]T, 0)
	fe.total = 0
}

// Sort sorts the elements.
//
// Parameters:
//   - isAsc: A boolean value that indicates whether to sort in ascending order.
func (fe *FrequencyEvaluator[T]) Sort(isAsc bool) {
	var filterFunc func(a, b T) int

	if isAsc {
		filterFunc = func(a, b T) int {
			if fe.elements[a] < fe.elements[b] {
				return -1
			} else if fe.elements[a] > fe.elements[b] {
				return 1
			} else {
				return 0
			}
		}
	} else {
		filterFunc = func(a, b T) int {
			if fe.elements[a] > fe.elements[b] {
				return -1
			} else if fe.elements[a] < fe.elements[b] {
				return 1
			} else {
				return 0
			}
		}
	}

	slices.SortStableFunc(fe.keys, filterFunc)
}
