package Slices

import (
	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/Common"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

// Laster is an interface for a stack of elements.
type Laster[T any] interface {
	// From creates a new Laster from a slice of elements.
	//
	// Parameters:
	//   - elems: The elements to add to the Laster.
	//
	// Returns:
	//   - Laster[T]: The new Laster.
	//   - error: An error if the elements could not be added.
	From(elems []T) (Laster[T], error)

	// GetLast gets the last element from the Laster.
	//
	// Returns:
	//   - T: The last element.
	//   - bool: If the last element was found.
	GetLast() (T, bool)

	// Append appends an element to the Laster.
	//
	// Parameters:
	//   - elem: The element to append.
	Append(elem T)

	uc.Copier
}

// NextsFunc is a function that gets the next elements.
//
// Parameters:
//   - wasDone: If the last element was done.
//   - last: The last element.
//
// Returns:
//   - []T: The next elements.
//   - error: An error if the next elements could not be found.
type NextsFunc[T any] func(wasDone bool, last T) ([]T, error)

// FilterNextsFunc is a function that filters the next elements.
//
// Parameters:
//   - wasDone: If the last element was done.
//   - nexts: The next elements.
//
// Returns:
//   - []T: The filtered next elements.
//   - error: An error if the next elements could not be filtered.
type FilterNextsFunc[T any] func(wasDone bool, nexts []T) ([]T, error)

// StackEvaluator evaluates a stack of elements.
type StackEvaluator[T any, E Laster[T]] struct {
	// eval is the evaluation function.
	eval uc.EvalOneFunc[T, bool]

	// nexts is the function to get the next elements.
	nexts NextsFunc[T]

	// filter is the function to filter the next elements.
	filter FilterNextsFunc[T]
}

// NewStackEvaluator creates a new StackEvaluator.
//
// Parameters:
//   - eval: The evaluation function.
//   - nexts: The function to get the next elements.
//
// Returns:
//   - *StackEvaluator[T, E]: The new StackEvaluator.
//   - error: An error of type *errors.ErrInvalidParameter if
//     the eval or nexts functions are nil.
//
// Behaviors:
//   - By default, the filter function will return the nexts if the
//     last element was not done. If the last element was done, it
//     will return nil.
func NewStackEvaluator[T any, E Laster[T]](eval uc.EvalOneFunc[T, bool], nexts NextsFunc[T]) (*StackEvaluator[T, E], error) {
	if eval == nil {
		return nil, ue.NewErrNilParameter("eval")
	} else if nexts == nil {
		return nil, ue.NewErrNilParameter("nexts")
	}

	return &StackEvaluator[T, E]{
		eval:  eval,
		nexts: nexts,
		filter: func(wasDone bool, nexts []T) ([]T, error) {
			if wasDone {
				return nil, nil
			}

			return nexts, nil
		},
	}, nil
}

// SetFilter sets the filter function.
//
// Parameters:
//   - filter: The filter function.
//
// Behaviors:
//   - If the filter function is nil, the normal filter function will be used.
func (se *StackEvaluator[T, E]) SetFilter(filter FilterNextsFunc[T]) {
	if filter == nil {
		return
	}

	se.filter = filter
}

// Evaluate evaluates the stack of elements from the given element.
//
// Parameters:
//   - elem: The element to start the evaluation.
//
// Returns:
//   - []E: The evaluated elements.
//   - error: An error if the elements could not be evaluated.
func (se *StackEvaluator[T, E]) Evaluate(elem T) ([]E, error) {
	var done []E

	first, err := (*new(E)).From([]T{elem})
	if err != nil {
		return nil, err
	}

	S := lls.NewLinkedStack(first.(E))

	for {
		top, err := S.Pop()
		if err != nil {
			break
		}

		last, ok := top.GetLast()
		if !ok {
			return nil, NewErrLastNotFound()
		}

		ok, err = se.eval(last)
		if err != nil {
			return nil, err
		}

		if ok {
			done = append(done, top)
		}

		nexts, err := se.nexts(ok, last)
		if err != nil {
			return nil, err
		}

		if len(nexts) == 0 {
			continue
		}

		nexts, err = se.filter(ok, nexts)
		if err != nil {
			return nil, err
		}

		for _, next := range nexts {
			topCopy := top.Copy().(E)

			topCopy.Append(next)

			err := S.Push(topCopy)
			if err != nil {
				return nil, err
			}
		}
	}

	return done, nil
}
