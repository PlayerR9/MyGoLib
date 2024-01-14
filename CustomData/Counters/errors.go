package Counters

import "fmt"

type ErrInvalidParameter struct {
	parameter string
	reason    error
}

func (e *ErrInvalidParameter) Parameter(p string) *ErrInvalidParameter {
	e.parameter = p
	return e
}

func (e *ErrInvalidParameter) Reason(r error) *ErrInvalidParameter {
	e.reason = r
	return e
}

func (e *ErrInvalidParameter) Error() string {
	return fmt.Sprintf("invalid parameter %s: reason=%v", e.parameter, e.reason)
}

func (e *ErrInvalidParameter) Unwrap() error {
	return e.reason
}

type ErrCannotAdvanceCounter struct {
	counter Counter
	reason  error
}

func (e *ErrCannotAdvanceCounter) Counter(c Counter) *ErrCannotAdvanceCounter {
	e.counter = c
	return e
}

func (e *ErrCannotAdvanceCounter) Reason(r error) *ErrCannotAdvanceCounter {
	e.reason = r
	return e
}

func (e *ErrCannotAdvanceCounter) Error() string {
	return fmt.Sprintf("cannot advance %T: reason=%v", e.counter, e.reason)
}

func (e *ErrCannotAdvanceCounter) Unwrap() error {
	return e.reason
}

type ErrCannotRetreatCounter struct {
	counter Counter
	reason  error
}

func (e *ErrCannotRetreatCounter) Counter(c Counter) *ErrCannotRetreatCounter {
	e.counter = c
	return e
}

func (e *ErrCannotRetreatCounter) Reason(r error) *ErrCannotRetreatCounter {
	e.reason = r
	return e
}

func (e *ErrCannotRetreatCounter) Error() string {
	return fmt.Sprintf("cannot retreat %T: reason=%v", e.counter, e.reason)
}

func (e *ErrCannotRetreatCounter) Unwrap() error {
	return e.reason
}
