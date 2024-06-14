package FSM

import (
	ut "github.com/PlayerR9/MyGoLib/Units/Tray"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// ActiveFSM is a struct that represents an active finite state machine.
type ActiveFSM[I any, S any, E uc.Enumer] struct {
	// currentState is the current state of the FSM.
	currentState S

	// Tray is the tray that the FSM uses to store data.
	Tray ut.Trayer[I]

	// m is a map that stores the values of the FSM.
	m map[E]any
}

// newActiveCmp is a constructor for ActiveFSM.
//
// Parameters:
//   - initState: The initial state of the FSM.
//   - tray: The tray that the FSM uses to store data.
//
// Returns:
//   - *ActiveFSM: A pointer to the newly created ActiveFSM.
func newActiveCmp[I any, S any, E uc.Enumer](initState S, tray ut.Trayer[I]) *ActiveFSM[I, S, E] {
	return &ActiveFSM[I, S, E]{
		currentState: initState,
		Tray:         tray,
		m:            make(map[E]any),
	}
}

// GetState returns the current state of the FSM.
//
// Returns:
//   - S: The current state of the FSM.
func (a *ActiveFSM[I, S, E]) GetState() S {
	return a.currentState
}

// changeState changes the current state of the FSM to newState.
//
// Parameters:
//   - newState: The new state of the FSM.
func (a *ActiveFSM[I, S, E]) changeState(newState S) {
	a.currentState = newState
	a.m = make(map[E]any)
}

// GetValue returns the value associated with key in the FSM.
//
// Parameters:
//   - key: The key of the value to retrieve.
//
// Returns:
//   - any: The value associated with key.
//   - bool: A boolean indicating whether the value was found.
func (a *ActiveFSM[I, S, E]) GetValue(key E) (any, bool) {
	val, ok := a.m[key]
	return val, ok
}
