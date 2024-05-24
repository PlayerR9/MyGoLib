package Operations

import (
	"fmt"
	"strings"

	tlt "github.com/PlayerR9/MyGoLib/TreeLike/Traversor"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
	uc "github.com/PlayerR9/MyGoLib/Units/Common"
)

type InfPrinter struct {
	indentLevel int
}

// Equals implements Common.Objecter.
func (i *InfPrinter) Equals(other uc.Objecter) bool {
	panic("unimplemented")
}

// String implements Common.Objecter.
func (i *InfPrinter) String() string {
	panic("unimplemented")
}

func (i *InfPrinter) Copy() uc.Objecter {
	return &InfPrinter{
		indentLevel: i.indentLevel,
	}
}

func NewInfPrinter() *InfPrinter {
	return &InfPrinter{
		indentLevel: 0,
	}
}

func (i *InfPrinter) IncIndent() {
	i.indentLevel++
}

func PrintTree[T any](tree *tr.Tree[T]) []string {
	var lines []string
	var builder strings.Builder

	f := func(elem T, obj uc.Objecter) (bool, error) {
		inf, ok := obj.(*InfPrinter)
		if !ok {
			return false, fmt.Errorf("invalid objecter type: %T", obj)
		}

		builder.WriteString(strings.Repeat("| ", inf.indentLevel))
		builder.WriteString(uc.StringOf(elem))
		builder.WriteString("\n")

		inf.IncIndent()

		return true, nil
	}

	err := tlt.DFS(tree, NewInfPrinter(), f)
	if err != nil {
		panic(err)
	}

	return lines
}
