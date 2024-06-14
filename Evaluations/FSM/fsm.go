package FSM

import (
	"fmt"

	ut "github.com/PlayerR9/MyGoLib/Units/Tray"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

// FSM is a struct that represents a finite state machine.
type FSM[I any, S any, R any, E uc.Enumer] struct {
	// InitFn is the function that initializes the FSM.
	InitFn InitFunc[I, S]

	// ShouldEndFn is the function that determines whether the FSM should end.
	ShouldEndFn EndCond[I, S, E]

	// GetResFn is the function that retrieves the result of the FSM.
	GetResFn EvalFunc[I, S, R, E]

	// NextFn is the function that transitions the FSM to the next state.
	NextFn TransFunc[I, S, E]

	// detsBefore is a map that stores the functions that are executed before the
	// FSM transitions to the next state.
	detsBefore map[E]DetFunc[I, S, E]

	// orderDets is a slice that stores the order in which the functions are executed.
	orderDets []E
}

// Run runs the FSM.
//
// Parameters:
//   - inputStream: The input stream for the FSM.
//
// Returns:
//   - []R: A slice of results from the FSM.
//   - error: An error if the function fails.
func (fsm *FSM[I, S, R, E]) Run(inputStream ut.Trayable[I]) ([]R, error) {
	if inputStream == nil {
		return nil, ue.NewErrNilParameter("inputStream")
	}

	var solution []R

	stream := inputStream.ToTray()
	stream.ArrowStart()

	initState, err := fsm.InitFn(stream)
	if err != nil {
		return solution, fmt.Errorf("error initializing: %w", err)
	}

	active := newActiveCmp[I, S, E](initState, stream)

	// End condition: Check if the FSM has reached the end.
	for {
		ok := fsm.ShouldEndFn(active)
		if ok {
			break
		}

		// Action: Determine all the elements of the FSM.
		for _, elem := range fsm.orderDets {
			fn, ok := fsm.detsBefore[elem]
			if !ok {
				return solution, fmt.Errorf("no function for element %s", elem.String())
			}

			sol, err := fn(active)
			if err != nil {
				return solution, fmt.Errorf("error determining %s: %w", elem.String(), err)
			}

			active.m[elem] = sol
		}

		// Transition: Get the element that will determine the next state.
		res, err := fsm.GetResFn(active)
		if err != nil {
			return solution, fmt.Errorf("error evaluating: %w", err)
		}

		solution = append(solution, res)

		// Transition: Change the state.
		nextState, err := fsm.NextFn(active)
		if err != nil {
			return solution, fmt.Errorf("error transitioning: %w", err)
		}

		active.changeState(nextState)
	}

	return solution, nil
}
