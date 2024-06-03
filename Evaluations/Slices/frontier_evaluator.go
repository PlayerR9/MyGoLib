package Slices

import (
	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/Common"
	hlp "github.com/PlayerR9/MyGoLib/Utility/Helpers"

	up "github.com/PlayerR9/MyGoLib/Units/Pair"
)

// Accepter is an interface that represents an accepter.
type Accepter interface {
	// Accept returns true if the accepter accepts the element.
	//
	// Returns:
	//   - bool: True if the accepter accepts the element, false otherwise.
	Accept() bool
}

// FrontierEvaluator is a type that represents a frontier evaluator.
type FrontierEvaluator[T Accepter] struct {
	// matcher is the matcher.
	matcher uc.EvalManyFunc[T, T]

	// solutions is the list of solutions.
	solutions []*hlp.WeightedHelper[T]
}

// NewFrontierEvaluator creates a new frontier evaluator.
//
// Parameters:
//   - matcher: The matcher.
//
// Returns:
//   - *FrontierEvaluator: The new frontier evaluator.
//
// Behaviors:
//   - If matcher is nil, then the frontier evaluator will return nil for any evaluation.
func NewFrontierEvaluator[T Accepter](matcher uc.EvalManyFunc[T, T]) *FrontierEvaluator[T] {
	fe := &FrontierEvaluator[T]{
		matcher:   matcher,
		solutions: make([]*hlp.WeightedHelper[T], 0),
	}

	return fe
}

// Evaluate evaluates the frontier evaluator given an element.
//
// Parameters:
//   - elem: The element to evaluate.
//
// Behaviors:
//   - If the element is accepted, the solutions will be set to the element.
//   - If the element is not accepted, the solutions will be set to the results of the matcher.
//   - If the matcher returns an error, the solutions will be set to the error.
//   - The evaluations assume that, the more the element is elaborated, the more the weight increases.
//     Thus, it is assumed to be the most likely solution as it is the most elaborated. Euristic: Depth.
func (fe *FrontierEvaluator[T]) Evaluate(elem T) {
	if fe.matcher == nil {
		fe.solutions = nil
		return
	}

	if elem.Accept() {
		fe.solutions = []*hlp.WeightedHelper[T]{hlp.NewWeightedHelper(elem, nil, 0.0)}
		return
	}

	fe.solutions = make([]*hlp.WeightedHelper[T], 0)

	S := Stacker.NewArrayStack(up.NewPair(elem, 0.0))

	for {
		p, err := S.Pop()
		if err != nil {
			break
		}

		nexts, err := fe.matcher(p.First)
		if err != nil {
			fe.solutions = append(fe.solutions, hlp.NewWeightedHelper(p.First, err, p.Second))
			continue
		}

		newPairs := make([]*up.Pair[T, float64], len(nexts))

		for _, next := range nexts {
			newPairs = append(newPairs, up.NewPair(next, p.Second+1.0))
		}

		for _, pair := range newPairs {
			if pair.First.Accept() {
				fe.solutions = []*hlp.WeightedHelper[T]{hlp.NewWeightedHelper(pair.First, nil, pair.Second)}
			} else {
				err := S.Push(pair)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

// GetResults gets the results of the frontier evaluator.
//
// Returns:
//   - []T: The results of the frontier evaluator.
//   - error: An error if the frontier evaluator failed.
//
// Behaviors:
//   - If the solutions are empty, the function returns nil.
//   - If the solutions contain errors, the function returns the first error.
//   - Otherwise, the function returns the solutions.
func (fe *FrontierEvaluator[T]) GetResults() ([]T, error) {
	if len(fe.solutions) == 0 {
		return nil, nil
	}

	results, ok := hlp.SuccessOrFail(fe.solutions, true)
	if !ok {
		// Determine the most likely error.
		// As of now, we will just return the first error.
		return hlp.ExtractResults(results), results[0].GetData().Second
	}

	return hlp.ExtractResults(results), nil
}
