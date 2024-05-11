package DoWhile

import (
	uc "github.com/PlayerR9/MyGoLib/Units/Common"
	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

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
	if len(todo) == 0 {
		return nil
	} else if accept == nil {
		return nil
	} else if f == nil {
		s1, _ := slext.SFSeparate(todo, accept)

		return s1
	}

	done := make([]T, 0)

	for len(todo) > 0 {
		newElem := make([]T, 0)

		s1, s2 := slext.SFSeparate(todo, accept)
		if len(s1) > 0 {
			done = append(done, s1...)
		}

		if len(s2) == 0 {
			break
		}

		for _, elem := range s2 {
			others, err := f(elem)
			if err != nil {
				continue
			}

			newElem = append(newElem, others...)
		}

		todo = newElem
	}

	return done
}
