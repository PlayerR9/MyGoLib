package CustomData

import (
	_ "github.com/markphelps/optional"

	itf "github.com/PlayerR9/MyGoLib/Interfaces"
)

//go:generate optional -type=Pair
type Pair[A any, B any] struct {
	First  A
	Second B
}

func (p *Pair[A, B]) Cleanup() {
	p.First = itf.Cleanup[A](p.First)
	p.Second = itf.Cleanup[B](p.Second)
}
