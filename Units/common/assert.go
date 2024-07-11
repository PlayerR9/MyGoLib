package common

import (
	"fmt"
	"reflect"
	"strings"
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

// AssertErr panics if the error is not nil.
//
// Parameters:
//   - err: The error to check.
//   - format: The format of the message to show if the error is not nil.
//   - args: The arguments to format the message.
//
// The format should be the function name and the args should be the parameters.
//
// Example:
//
//	func MyFunc(param1 string, param2 int) {
//	    res, err := SomeFunc(param1, param2)
//	    AssertErr(err, "SomeFunc(%s, %d)", param1, param2)
//	}
func AssertErr(err error, format string, args ...any) {
	if err == nil {
		return
	}

	var builder strings.Builder

	builder.WriteString("In ")
	fmt.Fprintf(&builder, format, args...)
	builder.WriteString(" = ")
	builder.WriteString(err.Error())

	msg := builder.String()

	panic(msg)
}
