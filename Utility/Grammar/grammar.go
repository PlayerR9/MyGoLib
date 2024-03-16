package Grammar

import (
	"fmt"
	"slices"
	"strings"
)

// GrammarBuilder represents a builder for a grammar.
//
// The default direction of the productions is LeftToRight.
type GrammarBuilder struct {
	// Slice of productions to add to the grammar.
	productions []*Production

	// Direction of the productions.
	direction ProductionDirection
}

// String is a method of GrammarBuilder that returns a string
// representation of a GrammarBuilder.
//
// Returns:
//
//   - string: A string representation of a GrammarBuilder.
func (b *GrammarBuilder) String() string {
	if b.productions == nil {
		return "GrammarBuilder[nil]"
	}

	if len(b.productions) == 0 {
		return "GrammarBuilder[total=0, productions=[]]"
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "GrammarBuilder[total=%d, productions=[%v", len(b.productions), b.productions[0])

	for _, production := range b.productions[1:] {
		fmt.Fprintf(&builder, ", %v", production)
	}

	builder.WriteString("]]")

	return builder.String()
}

// AddProduction is a method of GrammarBuilder that adds a production to
// the GrammarBuilder.
//
// Parameters:
//
//   - p: The production to add to the GrammarBuilder.
func (b *GrammarBuilder) AddProduction(p *Production) {
	if p == nil {
		return
	}

	if b.productions == nil {
		b.productions = []*Production{p}
	} else {
		b.productions = append(b.productions, p)
	}
}

// SetDirection is a method of GrammarBuilder that sets the direction of
// the productions in the GrammarBuilder.
//
// Parameters:
//
//   - direction: The direction of the productions.
func (b *GrammarBuilder) SetDirection(direction ProductionDirection) {
	b.direction = direction
}

// Build is a method of GrammarBuilder that builds a Grammar from the
// GrammarBuilder.
//
// Returns:
//
//   - *Grammar: A Grammar built from the GrammarBuilder.
func (b *GrammarBuilder) Build() *Grammar {
	if b.productions == nil {
		b.direction = LeftToRight

		return &Grammar{
			productions: make([]*Production, 0),
			symbols:     make([]string, 0),
		}
	}

	grammar := Grammar{
		symbols: make([]string, 0),
	}

	// 1. Remove duplicates
	for i := 0; i < len(b.productions); {
		index := slices.IndexFunc(b.productions[i+1:], func(p *Production) bool {
			return p.IsEqual(b.productions[i])
		})

		if index != -1 {
			b.productions = append(b.productions[:index], b.productions[index+1:]...)
		} else {
			i++
		}
	}

	// 2. Make sure all productions follow the same direction
	for _, p := range b.productions {
		p.SetDirection(b.direction)
	}

	// 3. Add productions to grammar
	grammar.productions = make([]*Production, len(b.productions))
	copy(grammar.productions, b.productions)

	// 3. Add symbols to grammar
	for _, p := range b.productions {
		for _, symbol := range p.GetSymbols() {
			if !slices.Contains(grammar.symbols, symbol) {
				grammar.symbols = append(grammar.symbols, symbol)
			}
		}
	}

	// 4. Clear builder
	for i := range b.productions {
		b.productions[i] = nil
	}

	b.productions = nil
	b.direction = LeftToRight

	return &grammar
}

// Reset is a method of GrammarBuilder that resets a GrammarBuilder.
func (b *GrammarBuilder) Reset() {
	for i := range b.productions {
		b.productions[i] = nil
	}

	b.productions = nil

	b.direction = LeftToRight
}

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
	productions []*Production

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

	fmt.Fprintf(&builder, "Grammar[productions=[%v", g.productions[0])

	for _, production := range g.productions[1:] {
		fmt.Fprintf(&builder, ", %v", production)
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
func (g *Grammar) GetProductions() []*Production {
	productions := make([]*Production, len(g.productions))
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
