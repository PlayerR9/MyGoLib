package FSM

import (
	"slices"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

// FsmBuilder is a struct that represents a builder for a finite state machine.
type FsmBuilder[I any, S any, R any, E uc.Enumer] struct {
	// InitFn is the function that initializes the FSM.
	InitFn InitFunc[I, S]

	// ShouldEndFn is the function that determines whether the FSM should end.
	ShouldEndFn EndCond[I, S, E]

	// GetResFn is the function that retrieves the result of the FSM.
	GetResFn EvalFunc[I, S, R, E]

	// NextFn is the function that transitions the FSM to the next state.
	NextFn TransFunc[I, S, E]

	// detsBefore is a map that stores the functions that are executed before the
	detsBefore map[E]DetFunc[I, S, E]

	// orderDets is a slice that stores the order in which the functions are executed.
	orderDets []E
}

// AddDetFn adds a function that is executed before the FSM transitions to the next state.
//
// Parameters:
//   - elem: The element that the function determines the value of.
//   - fn: The function that determines the value of the element.
func (b *FsmBuilder[I, S, R, E]) AddDetFn(elem E, fn DetFunc[I, S, E]) {
	if b.detsBefore == nil {
		b.detsBefore = make(map[E]DetFunc[I, S, E])
	}

	_, ok := b.detsBefore[elem]
	if ok {
		index := slices.Index(b.orderDets, elem)
		b.orderDets = slices.Delete(b.orderDets, index, index+1)
	}

	b.detsBefore[elem] = fn
	b.orderDets = append(b.orderDets, elem)
}

// Build creates a new FSM from the builder.
//
// Returns:
//   - *FSM: A pointer to the newly created FSM.
//   - error: An error if the function fails.
func (b *FsmBuilder[I, S, R, E]) Build() (*FSM[I, S, R, E], error) {
	if b.InitFn == nil {
		return nil, ue.NewErrNilParameter("InitFn")
	}

	if b.ShouldEndFn == nil {
		return nil, ue.NewErrNilParameter("ShouldEndFn")
	}

	if b.GetResFn == nil {
		return nil, ue.NewErrNilParameter("GetResFn")
	}

	if b.NextFn == nil {
		return nil, ue.NewErrNilParameter("NextFn")
	}

	alias := &FSM[I, S, R, E]{
		InitFn:      b.InitFn,
		ShouldEndFn: b.ShouldEndFn,
		GetResFn:    b.GetResFn,
		NextFn:      b.NextFn,
		detsBefore:  b.detsBefore,
		orderDets:   b.orderDets,
	}

	return alias, nil
}
