package FSM

import (
	ut "github.com/PlayerR9/MyGoLib/Units/Tray"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// DetFunc is a function that determines the value of an element of the FSM.
//
// Parameters:
//   - fsm: The active FSM.
//
// Returns:
//   - any: The value of the element.
//   - error: An error if the function fails.
type DetFunc[I any, S any, E uc.Enumer] func(fsm *ActiveFSM[I, S, E]) (any, error)

// TransFunc is a function that transitions the FSM to the next state.
//
// Parameters:
//   - fsm: The active FSM.
//
// Returns:
//   - S: The next state of the FSM.
//   - error: An error if the function fails.
type TransFunc[I any, S any, E uc.Enumer] func(fsm *ActiveFSM[I, S, E]) (S, error)

// EvalFunc is a function that evaluates the FSM.
//
// Parameters:
//   - fsm: The active FSM.
//
// Returns:
//   - any: The result of the evaluation.
//   - error: An error if the function fails.
type EvalFunc[I any, S any, R any, E uc.Enumer] func(fsm *ActiveFSM[I, S, E]) (R, error)

// EndCond is a function that determines whether the FSM should end.
//
// Parameters:
//   - fsm: The active FSM.
//
// Returns:
//   - bool: A boolean indicating whether the FSM should end.
type EndCond[I any, S any, E uc.Enumer] func(fsm *ActiveFSM[I, S, E]) bool

// InitFunc is a function that initializes the FSM.
//
// Parameters:
//   - tray: The tray that the FSM uses to store data.
//
// Returns:
//   - S: The initial state of the FSM.
//   - error: An error if the function fails.
type InitFunc[I any, S any] func(tray ut.Trayer[I]) (S, error)
