package Helpers

import up "github.com/PlayerR9/MyGoLib/Units/Pair"

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
//   - *up.Pair[O, error]: The data of the element.
//
// Behaviors:
//   - The second value of the pair is always nil.
func (we *WeightedElement[O]) GetData() *up.Pair[O, error] {
	return up.NewPair[O, error](we.elem, nil)
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
