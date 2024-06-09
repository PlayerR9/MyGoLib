package slice

import uc "github.com/PlayerR9/MyGoLib/Units/common"

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

// WeightedElement is a type that represents an element with a weight.
type WeightedElement[O any] struct {
	// Elem is the element.
	elem O

	// Weight is the weight of the element.
	weight float64
}

// GetData returns the data of the element.
//
// Returns:
//   - *uc.Pair[O, error]: The data of the element.
//
// Behaviors:
//   - The second value of the pair is always nil.
func (we *WeightedElement[O]) GetData() *uc.Pair[O, error] {
	return uc.NewPair[O, error](we.elem, nil)
}

// GetWeight returns the weight of the element.
//
// Returns:
//   - float64: The weight of the element.
func (we *WeightedElement[O]) GetWeight() float64 {
	return we.weight
}

// NewWeightedElement creates a new WeightedElement with the given element and weight.
//
// Parameters:
//   - elem: The element.
//   - weight: The weight of the element.
//
// Returns:
//   - *WeightedElement: The new WeightedElement.
func NewWeightedElement[O any](elem O, weight float64) *WeightedElement[O] {
	return &WeightedElement[O]{
		elem:   elem,
		weight: weight,
	}
}
