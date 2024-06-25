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

// Slicer is an interface that provides a method to convert a data structure to a slice.
type Slicer[T any] interface {
	// Slice returns a slice containing all the elements in the data structure.
	//
	// Returns:
	//   - []T: A slice containing all the elements in the data structure.
	Slice() []T

	Iterable[T]
}

// SliceOf converts any type to a slice of elements of the same type.
//
// Parameters:
//   - elem: The element to convert to a slice.
//
// Returns:
//   - []T: The slice representation of the element.
//
// Behaviors:
//   - Nil elements are converted to nil slices.
//   - Slice elements are returned as is.
//   - Slicer elements have their Slice method called.
//   - Other elements are converted to slices containing a single element.
func SliceOf[T any](elem any) []T {
	if elem == nil {
		return nil
	}

	switch elem := elem.(type) {
	case []T:
		return elem
	case Slicer[T]:
		slice := elem.Slice()
		return slice
	default:
		return []T{elem.(T)}
	}
}

// Iterable is an interface that defines a method to get an iterator over a
// collection of elements of type T. It is implemented by data structures that
// can be iterated over.
type Iterable[T any] interface {
	// Iterator returns an iterator over the collection of elements.
	//
	// Returns:
	//   - Iterater[T]: An iterator over the collection of elements.
	Iterator() Iterater[T]
}

// IteratorOf converts any type to an iterator over elements of the same type.
//
// Parameters:
//   - elem: The element to convert to an iterator.
//
// Returns:
//   - Iterater[T]: The iterator over the element.
//
// Behaviors:
//   - IF elem is nil, an empty iterator is returned.
//   - IF elem -implements-> Iterater[T], the element is returned as is.
//   - IF elem -implements-> Iterable[T], the element's Iterator method is called.
//   - IF elem -implements-> []T, a new iterator over the slice is created.
//   - ELSE, a new iterator over a single-element collection is created.
func IteratorOf[T any](elem any) Iterater[T] {
	if elem == nil {
		var builder Builder[T]

		return builder.Build()
	}

	var iter Iterater[T]

	switch elem := elem.(type) {
	case Iterater[T]:
		iter = elem
	case Iterable[T]:
		iter = elem.Iterator()
	case []T:
		iter = &SimpleIterator[T]{
			values: &elem,
			index:  0,
		}
	default:
		iter = &SimpleIterator[T]{
			values: &[]T{elem.(T)},
			index:  0,
		}
	}

	return iter
}

// Iterater is an interface that defines methods for an iterator over a
// collection of elements of type T.
type Iterater[T any] interface {
	// Size returns the number of elements in the collection.
	//
	// Returns:
	//  - count: The number of elements in the collection.
	Size() (count int)

	// Consume advances the iterator to the next element in the
	// collection and returns the current element.
	//
	// Returns:
	//  - T: The current element in the collection.
	//  - error: An error if the iterator is exhausted or if an error occurred
	//    while consuming the element.
	Consume() (T, error)

	// Restart resets the iterator to the beginning of the
	// collection.
	Restart()
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
