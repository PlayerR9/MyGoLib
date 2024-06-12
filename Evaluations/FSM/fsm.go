package FSM

import (
	"slices"

	ut "github.com/PlayerR9/MyGoLib/Units/Tray"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

type ActiveFSM[I any, S any, E uc.Enumer] struct {
	currentState S
	Tray         ut.Trayer[I]
	m            map[E]any
}

func newActiveCmp[I any, S any, E uc.Enumer](initState S, tray ut.Trayer[I]) *ActiveFSM[I, S, E] {
	return &ActiveFSM[I, S, E]{
		currentState: initState,
		Tray:         tray,
		m:            make(map[E]any),
	}
}

func (a *ActiveFSM[I, S, E]) GetState() S {
	return a.currentState
}

func (a *ActiveFSM[I, S, E]) changeState(newState S) {
	a.currentState = newState
	a.m = make(map[E]any)
}

func (a *ActiveFSM[I, S, E]) GetValue(key E) (any, bool) {
	val, ok := a.m[key]
	return val, ok
}

type FsmBuilder[I any, S any, R any, E uc.Enumer] struct {
	InitFn      InitFunc[I, S]
	ShouldEndFn EndCond[I, S, E]

	GetResFn EvalFunc[I, S, R, E]
	NextFn   TransFunc[I, S, E]

	detsBefore map[E]DetFunc[I, S, E]
	orderDets  []E
}

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

type FSM[I any, S any, R any, E uc.Enumer] struct {
	InitFn      InitFunc[I, S]
	ShouldEndFn EndCond[I, S, E]

	GetResFn EvalFunc[I, S, R, E]
	NextFn   TransFunc[I, S, E]

	detsBefore map[E]DetFunc[I, S, E]
	orderDets  []E
}
