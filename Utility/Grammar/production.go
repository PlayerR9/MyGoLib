package Grammar

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	itf "github.com/PlayerR9/MyGoLib/Interfaces"
	Stack "github.com/PlayerR9/MyGoLib/ListLike/Stack"
	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

type Productioner interface {
	fmt.Stringer
	Equals(Productioner) bool
	GetLhs() string
	GetSymbols() []string
	Match(int, any) Tokener

	itf.Copier
}

// Production represents a production in a grammar.
type Production struct {
	// Left-hand side of the production.
	lhs string

	// Right-hand side of the production.
	rhs []string
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
		return fmt.Sprintf("Production[%s %s]", p.lhs, LeftToRight)
	} else {
		return fmt.Sprintf("Production[%s %s %s]", p.lhs, LeftToRight, strings.Join(p.rhs, " "))
	}
}

// Equals is a method of Production that returns whether the production
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
func (p *Production) Equals(other Productioner) bool {
	if p == nil || other == nil || other.GetLhs() != p.lhs {
		return false
	}

	val, ok := other.(*Production)
	if !ok || len(val.rhs) != len(p.rhs) {
		return false
	}

	for i, symbol := range p.rhs {
		if symbol != val.rhs[i] {
			return false
		}
	}

	return true
}

// GetLhs is a method of Production that returns the left-hand side of
// the production.
//
// Returns:
//
//   - string: The left-hand side of the production.
func (p *Production) GetLhs() string {
	return p.lhs
}

// Iterator is a method of Production that returns an iterator for the
// production that iterates over the right-hand side of the production.
//
// Returns:
//
//   - itf.Iterater[string]: An iterator for the production.
func (p *Production) Iterator() itf.Iterater[string] {
	return itf.IteratorFromSlice(p.rhs)
}

// ReverseIterator is a method of Production that returns a reverse
// iterator for the production that iterates over the right-hand side of
// the production in reverse.
//
// Returns:
//
//   - itf.Iterater[string]: A reverse iterator for the production.
func (p *Production) ReverseIterator() itf.Iterater[string] {
	slice := make([]string, len(p.rhs))
	copy(slice, p.rhs)
	slices.Reverse(slice)

	return itf.IteratorFromSlice(slice)
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

// b must be of Stack.Stacker[Tokener]
func (p *Production) Match(at int, b any) Tokener {
	switch val := b.(type) {
	case Stack.Stacker[Tokener]:
		popped := Stack.NewArrayStack[Tokener]()

		for i := len(p.rhs) - 1; i >= 0; i-- {
			top, err := val.Peek()
			if err != nil {
				// Push back the popped symbols.
				for !popped.IsEmpty() {
					elem, _ := popped.Pop()
					val.Push(elem)
				}

				return nil
			}

			popped.Push(top)
			val.Pop()
		}

		slice := popped.Slice()
		slices.Reverse(slice)

		// Push back the popped symbols.
		for !popped.IsEmpty() {
			elem, _ := popped.Pop()
			val.Push(elem)
		}

		return NewNonLeafToken(p.lhs, at, slice...)
	default:
		return nil
	}
}

func (p *Production) Copy() itf.Copier {
	rhsCopy := make([]string, len(p.rhs))
	copy(rhsCopy, p.rhs)

	return &Production{
		lhs: p.lhs,
		rhs: rhsCopy,
	}
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

func (p *Production) IndexOfRhs(rhs string) int {
	return slices.Index(p.rhs, rhs)
}

type RegProduction struct {
	lhs string
	rhs string
	rxp *regexp.Regexp
}

func (r *RegProduction) String() string {
	if r == nil {
		return "RegProduction[nil]"
	}

	if r.rxp == nil {
		return fmt.Sprintf("RegProduction[lhs=%s, rhs=%s, rxp=N/A]", r.lhs, r.rhs)
	}

	return fmt.Sprintf("RegProduction[lhs=%s, rhs=%s, rxp=%v]", r.lhs, r.rhs, r.rxp)
}

func (p *RegProduction) Equals(other Productioner) bool {
	if p == nil || other == nil || other.GetLhs() != p.lhs {
		return false
	}

	val, ok := other.(*RegProduction)
	return ok && val.rhs == p.rhs
}

func (p *RegProduction) GetLhs() string {
	return p.lhs
}

func (p *RegProduction) GetSymbols() []string {
	return []string{p.lhs}
}

// return nil if no match
func (p *RegProduction) Match(at int, b any) Tokener {
	val, ok := b.([]byte)
	if !ok {
		return nil
	}

	data := p.rxp.Find(val)
	if data == nil {
		return nil
	}

	return NewLeafToken(p.lhs, string(data), at)
}

func (p *RegProduction) Copy() itf.Copier {
	return &RegProduction{
		lhs: p.lhs,
		rhs: p.rhs,
		rxp: p.rxp,
	}
}

func NewRegProduction(lhs string, regex string) *RegProduction {
	return &RegProduction{
		lhs: lhs,
		rhs: "^" + regex,
	}
}
