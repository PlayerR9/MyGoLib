package common

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Assert panics if the condition is false.
//
// Parameters:
//   - cond: The condition to check.
//   - msg: The message to show if the condition is false.
//
// The panic message is the string msg.
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
//
// The panic message is of error *ErrInvalidParameter.
func AssertParam(param string, cond bool, reason error) {
	if cond {
		return
	}

	err := NewErrInvalidParameter(param, reason)
	panic(err.Error())
}

// AssertIfZero panics if the element is zero.
//
// Parameters:
//   - elem: The element to check.
//   - msg: The message to show if the element is zero.
//
// The panic message is the string msg.
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
//
// The panic message is the equivalent of fmt.Sprintf(format, args).
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
//	    AssertErr(err, "SomeFunc(%s, %d)", param1, param2) // panic("In SomeFunc(param1, param2) = err")
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

// AssertOk panics if the condition is false.
//
// Parameters:
//   - ok: The condition to check.
//   - format: The format of the message to show if the condition is false.
//   - args: The arguments to format the message.
//
// The format should be the function name and the args should be the parameters.
//
// Example:
//
//	func MyFunc(param1 string, param2 int) {
//	    ok := SomeFunc(param1, param2)
//	    AssertOk(ok, "SomeFunc(%s, %d)", param1, param2) // panic("In SomeFunc(param1, param2) = false")
//	}
func AssertOk(ok bool, format string, args ...any) {
	if ok {
		return
	}

	var builder strings.Builder

	builder.WriteString("In ")
	fmt.Fprintf(&builder, format, args...)
	builder.WriteString(" = false")

	msg := builder.String()

	panic(msg)
}

// AssertNil panics if the element is nil but returns the element dereferenced
// if it is not nil.
//
// Parameters:
//   - elem: The element to check.
//   - param_name: The name of the parameter.
//
// Returns:
//   - T: The element if it is not nil.
//
// The panic message is the message "Parameter \"param_name\" must not be nil".
func AssertDerefNil[T any](elem *T, param_name string) T {
	if elem != nil {
		return *elem
	}

	values := []string{
		"Parameter",
		"(",
		strconv.Quote(param_name),
		")",
		"must not be nil",
	}

	panic(strings.Join(values, " "))
}

// AssertNil panics if the element is nil but returns the element dereferenced
// if it is not nil.
//
// Parameters:
//   - elem: The element to check.
//   - param_name: The name of the parameter.
//
// Returns:
//   - T: The element if it is not nil.
//
// The panic message is the message "Parameter \"param_name\" must not be nil".
func AssertNil[T any](elem *T, param_name string) {
	if elem != nil {
		return
	}

	values := []string{
		"Parameter",
		"(",
		strconv.Quote(param_name),
		")",
		"must not be nil",
	}

	panic(strings.Join(values, " "))
}
