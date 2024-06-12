package FSM

import (
	"fmt"

	ut "github.com/PlayerR9/MyGoLib/Units/Tray"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

type DetFunc[I any, S any, E uc.Enumer] func(*ActiveFSM[I, S, E]) (any, error)

type TransFunc[I any, S any, E uc.Enumer] func(*ActiveFSM[I, S, E]) (S, error)

type EvalFunc[I any, S any, R any, E uc.Enumer] func(*ActiveFSM[I, S, E]) (R, error)

type EndCond[I any, S any, E uc.Enumer] func(*ActiveFSM[I, S, E]) bool

type InitFunc[I any, S any] func(ut.Trayer[I]) (S, error)

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
