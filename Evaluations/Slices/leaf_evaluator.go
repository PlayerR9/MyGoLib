package Slices

import (
	ui "github.com/PlayerR9/MyGoLib/Units/Iterators"
	up "github.com/PlayerR9/MyGoLib/Units/Pair"
)

// LeafEvaluater is an interface that represents a leaf evaluater.
type LeafEvaluater[A, T, E, R any] interface {
	// Init is a function type that initializes the memory and the first branch.
	//
	// Parameters:
	//   - elems: The elements.
	//
	// Returns:
	//   - T: The first branch.
	//   - error: An error if the initialization fails.
	Init(elems []A) (T, error)

	// Core is a function type that performs the core evaluation.
	//
	// Parameters:
	//   - index: The index of the loop element.
	//   - lpe: The loop element.
	//
	// Returns:
	//   - *up.Pair[R, error]: The result and an error.
	//   - error: An error if the evaluation fails (reserved for panic-level of critical errors).
	Core(index int, lpe E) (*up.Pair[R, error], error)

	// Next is a function type that performs the next evaluation.
	//
	// Parameters:
	//   - pair: The result and an error.
	//   - branch: The current branch.
	//
	// Returns:
	//   - []T: The new branches.
	//   - error: An error if the evaluation fails (reserved for panic-level of critical errors).
	Next(pair *up.Pair[R, error], branch T) ([]T, error)

	ui.Iterable[E]
}

// LeafEvaluable is an interface that represents a leaf evaluable.
type LeafEvaluable[A, T, E, R any] interface {
	// Evaluator is a function type that returns the leaf evaluator.
	//
	// Returns:
	//   - LeafEvaluater[A, T, E, R]: The leaf evaluator.
	Evaluator() LeafEvaluater[A, T, E, R]
}

// Evaluate performs a leaf evaluation with a loop.
//
// Parameters:
//   - le: The evaluable element.
//   - args: The arguments.
//
// Returns:
//   - []T: The branches.
//   - error: An error if the evaluation fails (reserved for panic-level of critical errors).
//
// Behaviors:
//   - The function performs a leaf evaluation with a loop.
//   - The function returns the branches.
//   - If le is nil, the function returns nil.
func Evaluate[A, T, E, R any](elem LeafEvaluable[A, T, E, R], args []A) ([]T, error) {
	if elem == nil {
		return nil, nil
	}

	ev := elem.Evaluator()
	if ev == nil {
		return nil, nil
	}

	firstBranch, err := ev.Init(args)
	if err != nil {
		return nil, err
	}

	branches := []T{firstBranch}

	index := 0
	iter := ev.Iterator()

	for {
		value, err := iter.Consume()
		if err != nil {
			break
		}

		pair, err := ev.Core(index, value)
		if err != nil {
			return branches, err
		}

		var newBranches []T

		for _, branch := range branches {
			tmp, err := ev.Next(pair, branch)
			if err != nil {
				return branches, err
			}

			if len(tmp) > 0 {
				newBranches = append(newBranches, tmp...)
			}
		}

		if len(newBranches) == 0 {
			return newBranches, nil
		}

		branches = newBranches
		index++
	}

	return branches, nil
}
