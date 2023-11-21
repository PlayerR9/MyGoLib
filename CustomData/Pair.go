package CustomData

import _ "github.com/markphelps/optional"

//go:generate optional -type=Pair
type Pair[A any, B any] struct {
	First  A
	Second B
}
