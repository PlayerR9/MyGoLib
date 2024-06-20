package errors

// ErrOrSol is a struct that holds a list of errors and a list of solutions.
type ErrOrSol[T any] struct {
	// errorList is a list of errors.
	errorList []error

	// solutionList is a list of solutions.
	solutionList []T

	// level is the level of the error or solution.
	level int

	// ignoreErr is a flag that indicates if the error should be ignored.
	ignoreErr bool
}

// AddErr adds an error to the list of errors if the level is greater or equal
// to the current level.
//
// Parameters:
//   - err: The error to add.
//   - level: The level of the error.
//
// Behaviors:
//   - If an error has been added with a level greater than the current level,
//     the error list is reset and the new level is updated.
//   - If the error is nil, the ignoreErr flag is set to true and the error list is reset.
func (e *ErrOrSol[T]) AddErr(err error, level int) {
	if level < e.level || e.ignoreErr {
		// Do nothing.
		return
	}

	if err == nil {
		e.ignoreErr = true
		e.errorList = nil
	} else {
		if level == e.level {
			e.errorList = append(e.errorList, err)
		} else {
			e.errorList = []error{err}
			e.level = level
		}
	}
}

// AddSol adds a solution to the list of solutions if the level is greater or equal
// to the current level.
//
// Parameters:
//   - sol: The solution to add.
//   - level: The level of the solution.
//
// Behaviors:
//   - If a solution has been added with a level greater than the current level,
//     the solution list is reset and the new level is updated.
//   - This function sets the ignoreErr flag to true and resets the error list.
func (e *ErrOrSol[T]) AddSol(sol T, level int) {
	if level < e.level {
		// Do nothing.
		return
	}

	if e.level == level {
		e.solutionList = append(e.solutionList, sol)
	} else {
		e.solutionList = []T{sol}
		e.level = level
	}

	if !e.ignoreErr {
		e.ignoreErr = true
		e.errorList = nil
	}
}

// AddAny adds an element to the list of errors or solutions if the level is greater or equal
// to the current level.
//
// Parameters:
//   - elem: The element to add.
//   - level: The level of the element.
//
// Behaviors:
//   - If an error has been added with a level greater than the current level,
//     the error list is reset and the new level is updated.
//   - If a solution has been added with a level greater than the current level,
//     the solution list is reset and the new level is updated.
func (e *ErrOrSol[T]) AddAny(elem any, level int) {
	if level < e.level {
		// Do nothing.
		return
	}

	switch elem := elem.(type) {
	case error:
		if e.ignoreErr {
			// Do nothing.
			return
		}

		if elem == nil {
			e.ignoreErr = true
			e.errorList = nil
		} else {
			if level == e.level {
				e.errorList = append(e.errorList, elem)
			} else {
				e.errorList = []error{elem}
				e.level = level
			}
		}
	case T:
		if e.level == level {
			e.solutionList = append(e.solutionList, elem)
		} else {
			e.solutionList = []T{elem}
			e.level = level
		}

		if !e.ignoreErr {
			e.ignoreErr = true
			e.errorList = nil
		}
	}
}

// HasError checks if errors are not ignored and if the error list is not empty.
//
// Returns:
//   - bool: True if errors are not ignored and the error list is not empty, otherwise false.
func (e *ErrOrSol[T]) HasError() bool {
	return !e.ignoreErr && len(e.errorList) > 0
}

// GetErrors returns the list of errors.
//
// Returns:
//   - []error: The list of errors.
func (e *ErrOrSol[T]) GetErrors() []error {
	return e.errorList
}

// GetSolutions returns the list of solutions.
//
// Returns:
//   - []T: The list of solutions.
func (e *ErrOrSol[T]) GetSolutions() []T {
	return e.solutionList
}
