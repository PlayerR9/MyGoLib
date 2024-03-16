package Grammar

import (
	"fmt"
	"slices"
	"strings"

	itf "github.com/PlayerR9/MyGoLib/Interfaces"
	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

// ProductionDirection represents the direction of a production.
type ProductionDirection bool

const (
	// LeftToRight is the direction of a production from left to right.
	LeftToRight ProductionDirection = false

	// RightToLeft is the direction of a production from right to left.
	RightToLeft ProductionDirection = true
)

// String is a method of ProductionDirection that returns a string
// representation of a ProductionDirection. It returns "->" if the
// direction is LeftToRight, and "<-" if the direction is RightToLeft.
//
// Returns:
//
//   - string: A string representation of a ProductionDirection.
func (d ProductionDirection) String() string {
	if d == LeftToRight {
		return "->"
	} else {
		return "<-"
	}
}

// Production represents a production in a grammar.
type Production struct {
	// Left-hand side of the production.
	lhs string

	// Right-hand side of the production.
	rhs []string

	// Direction of the production.
	direction ProductionDirection
}

// String is a method of Production that returns a string representation
// of a Production.
//
// Returns:
//
//   - string: A string representation of a Production.
func (p *Production) String() string {
	if p == nil {
		return "Production[nil]"
	}

	if len(p.rhs) == 0 {
		return fmt.Sprintf("Production[%s %v]", p.lhs, p.direction)
	}

	if p.direction == LeftToRight {
		return fmt.Sprintf("Production[%s %v %s]", p.lhs, p.direction, strings.Join(p.rhs, " "))
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "Production[%s %v %s", p.lhs, p.direction, p.rhs[len(p.rhs)-1])

	for i := len(p.rhs) - 2; i >= 0; i-- {
		fmt.Fprintf(&builder, " %s", p.rhs[i])
	}

	builder.WriteString("]")

	return builder.String()
}

// Iterator is a method of Production that returns an iterator for the
// production that iterates over the right-hand side of the production.
//
// Returns:
//
//   - itf.Iterater[string]: An iterator for the production.
func (p *Production) Iterator() itf.Iterater[string] {
	if p.direction == LeftToRight {
		return itf.IteratorFromSlice(p.rhs)
	}

	var builder itf.Builder[string]

	for i := len(p.rhs) - 1; i >= 0; i-- {
		builder.Append(p.rhs[i])
	}

	return builder.Build()
}

// NewProduction is a function that returns a new Production with the
// given left-hand side and right-hand side.
//
// Parameters:
//
//   - lhs: The left-hand side of the production.
//   - rhs: The right-hand side of the production.
//
// Returns:
//
//   - *Production: A new Production with the given left-hand side and
//     right-hand side.
func NewProduction(lhs string, rhs ...string) *Production {
	return &Production{lhs: lhs, rhs: rhs}
}

// IsEqual is a method of Production that returns whether the production
// is equal to another production. Two productions are equal if their
// left-hand sides are equal and their right-hand sides are equal.
//
// Parameters:
//
//   - other: The other production to compare to.
//
// Returns:
//
//   - bool: Whether the production is equal to the other production.
func (p *Production) IsEqual(other *Production) bool {
	if other == nil {
		return false
	} else if p.lhs != other.lhs {
		return false
	} else if len(p.rhs) != len(other.rhs) {
		return false
	}

	for i, symbol := range p.rhs {
		if symbol != other.rhs[i] {
			return false
		}
	}

	return true
}

// GetSymbols is a method of Production that returns a slice of symbols
// in the production. The slice contains the left-hand side of the
// production and the right-hand side of the production, with no
// duplicates.
//
// Returns:
//
//   - []string: A slice of symbols in the production.
func (p *Production) GetSymbols() []string {
	symbols := make([]string, len(p.rhs)+1)
	copy(symbols, p.rhs)

	symbols[len(symbols)-1] = p.lhs

	// Remove duplicates
	for i := 0; i < len(symbols); {
		index := slices.Index(symbols[i+1:], symbols[i])

		if index != -1 {
			symbols = append(symbols[:index], symbols[index+1:]...)
		} else {
			i++
		}
	}

	return symbols
}

// SetDirection is a method of Production that sets the direction of the
// production.
//
// Parameters:
//
//   - direction: The direction of the production.
func (p *Production) SetDirection(direction ProductionDirection) {
	p.direction = direction
}

// Size is a method of Production that returns the number of symbols in
// the right-hand side of the production.
//
// Returns:
//
//   - int: The number of symbols in the right-hand side of the
//     production.
func (p *Production) Size() int {
	return len(p.rhs)
}

// GetRhsAt is a method of Production that returns the symbol at the
// given index in the right-hand side of the production.
//
// Parameters:
//
//   - index: The index of the symbol to get.
//
// Returns:
//
//   - string: The symbol at the given index in the right-hand side of
//     the production.
//   - error: An error of type *ErrInvalidParameter if the index is
//     invalid.
func (p *Production) GetRhsAt(index int) (string, error) {
	if index < 0 || index >= len(p.rhs) {
		return "", ers.NewErrInvalidParameter("index").
			Wrap(ers.NewErrOutOfBound(index, 0, len(p.rhs)))
	}

	return p.rhs[index], nil
}

// GetLHS is a method of Production that returns the left-hand side of
// the production.
//
// Returns:
//
//   - string: The left-hand side of the production.
func (p *Production) GetLHS() string {
	return p.lhs
}

// IsLeftToRight is a method of Production that returns whether the
// production is left-to-right.
//
// Returns:
//
//   - bool: Whether the production is left-to-right.
func (p *Production) IsLeftToRight() bool {
	return p.direction == LeftToRight
}
