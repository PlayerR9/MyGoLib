package Grammar

import (
	"fmt"
	"slices"
	"strings"
)

const (
	// LeftToRight is the direction of a production from left to right.
	LeftToRight   string = "->"
	StartSymbolID string = "source"
	EndSymbolID   string = "EOF"
)

// Grammar represents a context-free grammar.
//
// A context-free grammar is a set of productions, each of which
// consists of a non-terminal symbol and a sequence of symbols.
//
// The non-terminal symbol is the left-hand side of the production,
// and the sequence of symbols is the right-hand side of the production.
//
// The grammar also contains a set of symbols, which are the
// non-terminal and terminal symbols in the grammar.
type Grammar struct {
	// productions is a slice of productions in the grammar.
	productions []Productioner

	// skipProductions is a slice of booleans that indicate whether to
	// skip a production. If skipProductions[i] is true, then
	// productions[i] is skipped.
	skipProductions []bool

	// symbols is a slice of symbols in the grammar.
	symbols []string
}

// String is a method of Grammar that returns a string representation
// of a Grammar.
//
// Returns:
//
//   - string: A string representation of a Grammar.
func (g *Grammar) String() string {
	if g == nil {
		return "Grammar[nil]"
	}

	if len(g.productions) == 0 {
		return "Grammar[prouctions=[], symbols=[]]"
	}

	var builder strings.Builder

	if g.skipProductions[0] {
		fmt.Fprintf(&builder, "Grammar[productions=[%v (skip)", g.productions[0])
	} else {
		fmt.Fprintf(&builder, "Grammar[productions=[%v", g.productions[0])
	}

	for i, production := range g.productions[1:] {
		if g.skipProductions[i+1] {
			fmt.Fprintf(&builder, ", %v (skip)", production)
		} else {
			fmt.Fprintf(&builder, ", %v", production)
		}
	}

	fmt.Fprintf(&builder, "], symbols=[%v", g.symbols[0])

	for _, symbol := range g.symbols[1:] {
		fmt.Fprintf(&builder, ", %v", symbol)
	}

	builder.WriteString("]]")

	return builder.String()
}

// GetProductions is a method of Grammar that returns a copy of the
// productions in the Grammar.
//
// Returns:
//
//   - []*Production: A copy of the productions in the Grammar.
func (g *Grammar) GetProductions() []Productioner {
	productions := make([]Productioner, len(g.productions))
	copy(productions, g.productions)

	return productions
}

// GetSymbols is a method of Grammar that returns a copy of the symbols
// in the Grammar.
//
// Returns:
//
//   - []string: A copy of the symbols in the Grammar.
func (g *Grammar) GetSymbols() []string {
	symbols := make([]string, len(g.symbols))
	copy(symbols, g.symbols)

	return symbols
}

func (g *Grammar) LhsToSkip() []string {
	if len(g.productions) == 0 {
		return make([]string, 0)
	}

	skip := make([]string, 0)

	for i, production := range g.productions {
		if !g.skipProductions[i] {
			continue
		}

		lhs := production.GetLhs()

		if !slices.Contains(skip, lhs) {
			skip = append(skip, lhs)
		}
	}

	return skip
}

type MatchedResult struct {
	Matched   Tokener
	RuleIndex int
}

func NewMatchResult(matched Tokener, ruleIndex int) MatchedResult {
	return MatchedResult{Matched: matched, RuleIndex: ruleIndex}
}

func (g *Grammar) Match(at int, b any) []MatchedResult {
	matches := make([]MatchedResult, 0)

	for i, p := range g.productions {
		matched := p.Match(at, b)
		if matched != nil {
			matches = append(matches, NewMatchResult(matched, i))
		}
	}

	return matches
}
