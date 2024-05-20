package Slices

import (
	hlp "github.com/PlayerR9/MyGoLib/CustomData/Helpers"
	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/Common"
	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

// Accepter is an interface that represents an accepter.
type Accepter interface {
	// Accept returns true if the accepter accepts the element.
	//
	// Returns:
	//   - bool: True if the accepter accepts the element, false otherwise.
	Accept() bool
}

// Evaluate is the main function of the tree evaluator.
//
// Parameters:
//   - source: The source to evaluate.
//   - root: The root of the tree evaluator.
//
// Returns:
//   - error: An error if lexing fails.
//
// Errors:
//   - *ErrEmptyInput: The source is empty.
//   - *ers.ErrAt: An error occurred at a specific index.
//   - *ErrAllMatchesFailed: All matches failed.
func FrontierEvaluate[T Accepter](elem T, matcher uc.EvalManyFunc[T, T]) ([]T, error) {
	if matcher == nil {
		return nil, nil
	} else if elem.Accept() {
		return []T{elem}, nil
	}

	solutions := make([]hlp.SimpleHelper[T], 0)

	S := Stacker.NewArrayStack(elem)

	for {
		elem, err := S.Pop()
		if err != nil {
			break
		}

		nexts, err := matcher(elem)
		if err != nil {
			solutions = append(solutions, hlp.NewHResult(elem, err))

			continue
		}

		for _, next := range nexts {
			if next.Accept() {
				solutions = append(solutions, hlp.NewHResult(next, nil))

				continue
			} else {
				err := S.Push(next)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	solutions, ok := hlp.FilterSuccessOrFail(solutions)
	if !ok {
		// Determine the most likely error.
		// As of now, we will just return the first error.
		return nil, solutions[0].Second
	}

	// TODO: Fix once the MyGoLib is updated.

	result := make([]T, 0, len(solutions))

	for _, solution := range solutions {
		result = append(result, solution.First)
	}

	return result, nil
}

// DoWhile performs a do-while loop on a slice of elements.
//
// Parameters:
//   - todo: The elements to perform the do-while loop on.
//   - accept: The predicate filter to accept elements.
//   - f: The evaluation function to perform on the elements.
//
// Returns:
//   - []T: The elements that were accepted.
//
// Behaviors:
//   - If todo is empty, the function returns nil.
//   - If accept is nil, the function returns nil.
//   - If f is nil, the function returns the application of accept on todo.
//   - The function performs the do-while loop on the elements in todo.
func DoWhile[T any](todo []T, accept slext.PredicateFilter[T], f uc.EvalManyFunc[T, T]) []T {
	if len(todo) == 0 || accept == nil {
		return nil
	}

	done := make([]T, 0)

	for len(todo) > 0 {
		s1, s2 := slext.SFSeparate(todo, accept)
		if len(s1) > 0 {
			done = append(done, s1...)
		}

		if f == nil {
			return done
		}

		todo = todo[:0]

		for _, elem := range s2 {
			others, err := f(elem)
			if err != nil {
				continue
			}

			todo = append(todo, others...)
		}
	}

	return done
}
