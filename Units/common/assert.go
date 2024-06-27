package common

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
