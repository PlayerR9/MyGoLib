package common

import (
	"fmt"
	"reflect"
)

// Assert panics if the condition is false.
//
// Parameters:
//   - cond: The condition to check.
//   - msg: The message to show if the condition is false.
func Assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

// AssertParam panics if the condition is false with a parameter name.
//
// Parameters:
//   - param: The name of the parameter.
//   - cond: The condition to check.
//   - reason: The reason why the parameter is invalid.
func AssertParam(param string, cond bool, reason error) {
	if cond {
		return
	}

	err := NewErrInvalidParameter(param, reason)
	panic(err)
}

// AssertIfZero panics if the element is zero.
//
// Parameters:
//   - elem: The element to check.
//   - msg: The message to show if the element is zero.
func AssertIfZero(elem any, msg string) {
	value := reflect.ValueOf(elem)
	ok := value.IsZero()
	if ok {
		panic(msg)
	}
}

// AssertF panics if the condition is false.
//
// Parameters:
//   - cond: The condition to check.
//   - format: The format of the message to show if the condition is false.
//   - args: The arguments to format the message.
func AssertF(cond bool, format string, args ...any) {
	if cond {
		return
	}

	msg := fmt.Sprintf(format, args...)
	panic(msg)
}
