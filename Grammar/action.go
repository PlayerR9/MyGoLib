package Grammar

import uc "github.com/PlayerR9/MyGoLib/Units/common"

// Actioner is an interface that defines the behavior of an action.
type Actioner interface {
}

// ShiftAction is the shift action.
type ShiftAction struct {
}

// NewShiftAction creates a new shift action.
//
// Returns:
//   - *ShiftAction: The new shift action. Never returns nil.
func NewShiftAction() *ShiftAction {
	return &ShiftAction{}
}

// ReduceAction is the reduce action.
type ReduceAction[T TokenTyper] struct {
	// rule is the rule to reduce.
	rule *Rule[T]
}

// NewReduceAction creates a new reduce action.
//
// Parameters:
//   - rule: The rule to reduce.
//
// Returns:
//   - *ReduceAction: The new reduce action.
//   - error: An error of type *common.ErrInvalidParameter if the rule is nil.
func NewReduceAction[T TokenTyper](rule *Rule[T]) (*ReduceAction[T], error) {
	if rule == nil {
		return nil, uc.NewErrNilParameter("rule")
	}

	return &ReduceAction[T]{
		rule: rule,
	}, nil
}

// AcceptAction is the accept action.
type AcceptAction[T TokenTyper] struct {
	// rule is the rule to accept.
	rule *Rule[T]
}

// NewAcceptAction creates a new accept action.
//
// Parameters:
//   - rule: The rule to accept.
//
// Returns:
//   - *AcceptAction: The new accept action.
//   - error: An error of type *common.ErrInvalidParameter if the rule is nil.
func NewAcceptAction[T TokenTyper](rule *Rule[T]) (*AcceptAction[T], error) {
	if rule == nil {
		return nil, uc.NewErrNilParameter("rule")
	}

	return &AcceptAction[T]{
		rule: rule,
	}, nil
}