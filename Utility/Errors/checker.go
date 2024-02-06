// Package Errors provides advanced error handling mechanisms.
package Errors

import (
	"fmt"
)

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

	if recErr, ok := r.(error); ok {
		*err = recErr
	} else {
		*err = fmt.Errorf("panic: %v", r)
	}
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
		if r != nil {
			if recErr, ok := r.(error); ok {
				err = recErr
			} else {
				panic(fmt.Errorf("panic: %v", r))
			}
		}
	}()

	f(param)

	return
}

// builder is a generic struct that provides a fluent interface for
// building error handling chains and validating values.
type builder[T any] struct {
	// target is the value that the builder is operating on.
	target T
}

// Done is a method of the builder[T] that returns the target value of the
// builder. This method is used to complete the error handling chain and
// return the final value (if needed).
//
// Returns:
//
//   - T: The target value of the builder.
func (b *builder[T]) Done() T {
	return b.target
}

// On is a constructor function that creates a new builder[T] with the
// provided target value.
//
// Parameters:
//
//   - target: The target value to operate on.
//
// Returns:
//
//   - *builder[T]: A new builder[T] with the provided target value.
func On[T any](target T) *builder[T] {
	return &builder[T]{
		target: target,
	}
}

// Check is a method that validates a condition on the target value of the
// builder.
// If the condition is not met, it triggers a panic with the error returned
// by the condition.
//
// This method is designed to be used in conjunction with RecoverFromPanic.
// This combination allows a function to handle errors without the need for
// explicit error checks at every step.
//
// Parameters:
//
//   - conditions: A list of conditions to check on the target value of the
//     builder. Each condition should return an error if the condition is
//     not met.
//
// Returns:
//
//   - *builder[T]: The builder with the target value.
//
// Example:
//
//	 func MyFunction(n int) (result int, err error) {
//		 defer RecoverFromPanic(&err)
//
//		 On(n).Check(
//		 	InvalidParameter("n").Positive(),
//		 	InvalidParameter("n").Not(0),
//		 )
//
//		 return 60 / n, nil
//	 }
//
// In this example, if n is less than or equal to 0, the Check method will
// trigger a panic with an ErrInvalidParameter error. The RecoverFromPanic
// function will then capture the panic and assign the error to the err
// variable.
//
// The Check method allows you to define custom conditions for your types,
// making it a flexible tool for error handling.
func (b *builder[T]) Check(conditions ...ConditionCheck) *builder[T] {
	for _, condition := range conditions {
		err := condition(b.target)
		if err != nil {
			panic(err)
		}
	}

	return b
}

// CheckFunc is a generic function that accepts a function and a builder.
// It executes the provided function with the target value of the builder and
// returns a new builder with the result of the function as its target value.
//
// If the provided function returns an error, CheckFunc triggers a panic with
// the error. The panic can be recovered using the RecoverFromPanic function.
//
// This function is useful for chaining operations on a value while handling
// errors in a consistent manner.
//
// Parameters:
//
//   - f: The function to execute. It should accept a value of type I and
//     return a value of type O and an error.
//   - b: The builder whose target value should be passed to the function.
//
// Returns:
//
//   - *builder[O]: A new builder with the result of the function as its
//     target value.
//
// Example:
//
//	 CheckFunc(func(n int) (int, error) {
//		if n <= 0 {
//			return 0, NewErrInvalidParameter("n", fmt.Errorf("value (%d) must be positive", n))
//		}
//
//		return 60 / n, nil
//	 }, On(42))
func CheckFunc[I any, O any](f func(I) (O, error), b *builder[I]) *builder[O] {
	res, err := f(b.target)
	if err != nil {
		panic(err)
	}

	return &builder[O]{
		target: res,
	}
}

// Using is a method of builder[T] that allows the user to override the error
// returned by the builder with a custom error. If any of the provided conditions
// are not met, it triggers a panic with the provided custom error.
//
// This method is designed to be used in conjunction with RecoverFromPanic.
// This combination allows a function to handle errors without the need for
// explicit error checks at every step.
//
// Parameters:
//
//   - err: The error to use as the override.
//   - conditions: A list of functions that check conditions on the target value
//     of the builder. Each function should return an error if the condition is not
//     met.
//
// Returns:
//
//   - *builder[T]: The builder with the target value.
//
// Example:
//
//	 On(42).Using(fmt.Errorf("custom error"), func(target int) error {
//		if target <= 0 {
//			return NewErrInvalidParameter("n", fmt.Errorf("value (%d) must be positive", target))
//			}
//		return nil
//	 })
//
// In this example, the call to On will return the custom error provided to the
// Using method if the condition is not met.
//
// The Using method allows you to override the error returned by the builder
// with a custom error, making it a flexible tool for error handling.
func (b *builder[T]) Using(err error, conditions ...ConditionCheck) *builder[T] {
	for _, condition := range conditions {
		if condition(b.target) != nil {
			panic(err)
		}
	}

	return b
}
