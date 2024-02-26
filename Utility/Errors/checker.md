// Package Errors provides advanced error handling mechanisms.
package Errors

// RecoverFromPanic is a function that recovers from a panic and sets the
// provided error pointer to the recovered value. If the recovered value is
// an error, it sets the error pointer to that error. Otherwise, it creates
// a new error with the recovered value and sets the error pointer to that.
//
// This function should be called using the defer statement to recover from
// panics.
//
// Parameters:
//
//   - err: The error pointer to set.
//
// Example:
//
//	func MyFunction(n int) (result int, err error) {
//	    defer RecoverFromPanic(&err)
//	    // ...
//	    panic("something went wrong")
//	}
func RecoverFromPanic(err *error) {
	r := recover()
	if r == nil {
		return
	}

	if x, ok := r.(error); ok {
		if IsNoError(x) {
			*err = nil
		} else {
			*err = x
		}
	} else {
		*err = &ErrPanic{value: r}
	}
}

// PropagatePanic is a function that recovers from a panic and propagates it
// with the provided error. If the recovered value is an error, it panics
// with the provided error and the recovered error as its reason. If the
// recovered value is not an error, it panics with a new error that contains
// the recovered value.
//
// This function should be called using the defer statement to recover from
// panics.
//
// Parameters:
//
//   - err: The error to propagate with the recovered value.
func PropagatePanic[T interface{ Wrap(error) T }](err T) {
	r := recover()
	if r == nil {
		return
	}

	if x, ok := r.(error); ok {
		err.Wrap(x)
	} else {
		err.Wrap(&ErrPanic{value: r})
	}

	panic(err)
}

// ErrorOf is a generic function that accepts a function and a parameter
// of any type.
// It executes the provided function with the given parameter and returns
// any error that might occur during the execution of the function.
//
// The function uses a deferred function to recover from any panics that
// might occur during the execution of the provided function. If the
// recovered value is an error, it sets the returned error to that error.
// If the recovered value is not an error, it panics again with a new error
// that contains the recovered value.
//
// This function is useful for checking if a function can handle certain
// inputs without causing a panic.
//
// Parameters:
//
//   - f: The function to execute.
//   - param: The parameter to pass to the function.
//
// Returns:
//
//   - err: The error that occurred during the execution of the function, or nil if no error occurred.
//
// Example:
//
//		err := ErrorOf(func(n int) { panic("something went wrong") }, 42)
//		if err != nil {
//	    	fmt.Println(err)
//		}
func ErrorOf[T any](f func(T), param T) (err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		if x, ok := r.(error); ok {
			if IsNoError(x) {
				err = nil
			} else {
				err = x
			}
		} else {
			err = &ErrPanic{value: r}
		}
	}()

	f(param)

	return
}

// CheckFunc is a generic function that accepts a function and a parameter.
// It executes the provided function with the parameter and returns the result
// of the function if and only if the function does not return an error.
//
// Any error returned by the function triggers a panic. The panic can be
// recovered using the RecoverFromPanic function.
// However, it doesn't catch any panics that might occur during the execution
// of the function.
//
// Parameters:
//
//   - f: The function to execute. It should accept a value of type I and
//     return a value of type O and an error.
//   - param: The parameter to pass to the function.
//
// Returns:
//
//   - O: The result of the function if the function does not return an error.
//
// Example:
//
//	 CheckFunc(func(n int) (int, error) {
//		if n <= 0 {
//			return 0, NewErrInvalidParameter("n", fmt.Errorf("value (%d) must be positive", n))
//		}
//
//		return 60 / n, nil
//	 }, 42)
func CheckFunc[O any, I any](f func(I) (O, error), param I) O {
	res, err := f(param)
	if IsNoError(err) {
		return res
	}

	if x, ok := err.(*ErrPanic); ok {
		panic(x.value)
	} else {
		panic(err)
	}
}

// TryFunc is a function that accepts a function and, after having executed
// it, checks if the function returned an error. If the function returned an
// error, TryFunc triggers a panic with the error. The panic can be recovered
// using the RecoverFromPanic function.
//
// Parameters:
//
//   - f: The function to execute.
//
// Example:
//
//	 n := 42
//	 TryFunc(func() error {
//		if n <= 0 {
//			return NewErrInvalidParameter("n").
//	     	Wrap(fmt.Errorf("value (%d) must be positive", n))
//		}
//
//		return nil
//	})
func TryFunc(f func() error) {
	err := f()
	if IsNoError(err) {
		return
	}

	panic(err)
}
