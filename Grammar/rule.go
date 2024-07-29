package Grammar

import (
	uc "github.com/PlayerR9/lib_units/common"
)

// Rule is a struct that represents a rule of type T.
type Rule[T TokenTyper] struct {
	// lhs is the left-hand side of the rule.
	lhs T

	// rhss is the right-hand side of the rule.
	rhss []T
}

// Iterator implements the uc.Iterable interface.
//
// Never returns nil.
func (r *Rule[T]) Iterator() uc.Iterater[T] {
	return uc.NewSimpleIterator(r.rhss)
}

// NewRule creates a new rule.
//
// Parameters:
//   - lhs: The left-hand side of the rule.
//   - rhss: The right-hand side of the rule.
//
// Returns:
//   - *Rule: The new rule.
//
// Returns nil if the rhss is empty.
func NewRule[T TokenTyper](lhs T, rhss []T) *Rule[T] {
	if len(rhss) == 0 {
		return nil
	}

	return &Rule[T]{
		lhs:  lhs,
		rhss: rhss,
	}
}
