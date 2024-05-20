package SliceExt

// WeightFunc is a type that defines a function that assigns a weight to an element.
//
// Parameters:
//   - elem: The element to assign a weight to.
//
// Returns:
//   - float64: The weight of the element.
//   - bool: True if the weight is valid, otherwise false.
type WeightFunc[O any] func(elem O) (float64, bool)

// Weighter is an interface that represents a type that can assign weights to elements.
type Weighter[O any] interface {
	// GetData returns the data of the element.
	//
	// Returns:
	//   - O: The data of the element.
	GetData() O

	// GetWeight returns the weight of the element.
	//
	// Returns:
	//   - float64: The weight of the element.
	GetWeight() float64
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
//   - O: The data of the element.
func (we *WeightedElement[O]) GetData() O {
	return we.elem
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
func FilterByPositiveWeight[T Weighter[O], O any](S []T) []O {
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

	solution := make([]O, len(indices))
	for i, index := range indices {
		solution[i] = S[index].GetData()
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
func FilterByNegativeWeight[T Weighter[O], O any](S []T) []O {
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

	solution := make([]O, len(indices))
	for i, index := range indices {
		solution[i] = S[index].GetData()
	}

	return solution
}
