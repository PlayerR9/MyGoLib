package common

import "errors"

// ExtractFirsts extracts all the first elements from the given slice of pairs.
//
// Parameters:
//   - pairs: The slice of pairs.
//
// Returns:
//   - []A: The slice of first elements.
func ExtractFirsts[A any, B any](pairs []Pair[A, B]) []A {
	if len(pairs) == 0 {
		return nil
	}

	firsts := make([]A, 0, len(pairs))

	for _, pair := range pairs {
		firsts = append(firsts, pair.First)
	}

	return firsts
}

// ExtractSeconds extracts all the second elements from the given slice of pairs.
//
// Parameters:
//   - pairs: The slice of pairs.
//
// Returns:
//   - []B: The slice of second elements.
func ExtractSeconds[A any, B any](pairs []Pair[A, B]) []B {
	if len(pairs) == 0 {
		return nil
	}

	seconds := make([]B, 0, len(pairs))

	for _, pair := range pairs {
		seconds = append(seconds, pair.Second)
	}

	return seconds
}

// Unwrapper is an interface that defines a method to unwrap an error.
type Unwrapper interface {
	// Unwrap returns the error that this error wraps.
	//
	// Returns:
	//   - error: The error that this error wraps.
	Unwrap() error

	// ChangeReason changes the reason of the error.
	//
	// Parameters:
	//   - reason: The new reason of the error.
	ChangeReason(reason error)

	error
}

// Is is function that checks if an error is of type T.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - bool: true if the error is of type T, false otherwise (including if the error is nil).
func Is[T error](err error) bool {
	if err == nil {
		return false
	}

	var target T

	ok := errors.As(err, &target)
	return ok
}

// IsNoError checks if an error is a no error error or if it is nil.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - bool: True if the error is a no error error or if it is nil, otherwise false.
func IsNoError(err error) bool {
	if err == nil {
		return true
	}

	var errNoError *ErrNoError

	ok := errors.As(err, &errNoError)
	return ok
}

// IsErrIgnorable checks if an error is an *ErrIgnorable or *ErrInvalidParameter error.
// If the error is nil, the function returns false.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - bool: True if the error is an *ErrIgnorable or *ErrInvalidParameter error,
//     otherwise false.
func IsErrIgnorable(err error) bool {
	if err == nil {
		return false
	}

	var ignorable *ErrIgnorable

	ok := errors.As(err, &ignorable)
	if ok {
		return true
	}

	var invalid *ErrInvalidParameter

	ok = errors.As(err, &invalid)
	return ok
}

// LimitErrorMsg limits the error message to a certain number of unwraps.
// It returns the top level error for allowing to print the error message
// with the limit of unwraps applied.
//
// If the error is nil or the limit is less than 0, the function does nothing.
//
// Parameters:
//   - err: The error to limit.
//   - limit: The limit of unwraps.
//
// Returns:
//   - error: The top level error with the limit of unwraps applied.
func LimitErrorMsg(err error, limit int) error {
	if err == nil || limit < 0 {
		return err
	}

	currErr := err

	for i := 0; i < limit; i++ {
		target, ok := currErr.(Unwrapper)
		if !ok {
			return err
		}

		reason := target.Unwrap()
		if reason == nil {
			return err
		}

		currErr = reason
	}

	// Limit reached
	target, ok := currErr.(Unwrapper)
	if !ok {
		return err
	}

	target.ChangeReason(nil)

	return err
}
