package Grammar

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// GrammarBuilder represents a builder for a grammar.
//
// The default direction of the productions is LeftToRight.
type GrammarBuilder struct {
	// Slice of productions to add to the grammar.
	productions []Productioner

	// Slice of booleans that indicate whether to skip a production.
	// If skipProductions[i] is true, then productionss[i] is skipped.
	skipProductions []bool
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
		return "GrammarBuilder[total=0, productions=[], skipProductions=[]]"
	}

	var builder strings.Builder

	if b.skipProductions[0] {
		fmt.Fprintf(&builder, "GrammarBuilder[total=%d, productions=[%v (skip)", len(b.productions), b.productions[0])
	} else {
		fmt.Fprintf(&builder, "GrammarBuilder[total=%d, productions=[%v", len(b.productions), b.productions[0])
	}

	for i, production := range b.productions[1:] {
		if b.skipProductions[i+1] {
			fmt.Fprintf(&builder, ", %v (skip)", production)
		} else {
			fmt.Fprintf(&builder, ", %v", production)
		}
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
func (b *GrammarBuilder) AddProduction(p Productioner, skip bool) {
	if p == nil {
		return
	}

	if b.productions == nil {
		b.productions = []Productioner{p}
		b.skipProductions = []bool{skip}
	} else {
		b.productions = append(b.productions, p)
		b.skipProductions = append(b.skipProductions, skip)
	}
}

// Build is a method of GrammarBuilder that builds a Grammar from the
// GrammarBuilder.
//
// Returns:
//
//   - *Grammar: A Grammar built from the GrammarBuilder.
func (b *GrammarBuilder) Build() (*Grammar, error) {
	if b.productions == nil {
		return &Grammar{
			productions:     make([]Productioner, 0),
			symbols:         make([]string, 0),
			skipProductions: make([]bool, 0),
		}, nil
	}

	grammar := Grammar{
		symbols: make([]string, 0),
	}

	// 1. Remove duplicates
	for i := 0; i < len(b.productions); {
		index := slices.IndexFunc(b.productions[i+1:], func(p Productioner) bool {
			return p.Equals(b.productions[i])
		})

		if index != -1 {
			b.productions = append(b.productions[:index], b.productions[index+1:]...)
		} else {
			i++
		}
	}

	// 2. Add productions to grammar
	grammar.productions = make([]Productioner, len(b.productions))
	copy(grammar.productions, b.productions)

	grammar.skipProductions = make([]bool, len(b.skipProductions))
	copy(grammar.skipProductions, b.skipProductions)

	// 3. Add symbols to grammar
	for _, p := range b.productions {
		for _, symbol := range p.GetSymbols() {
			if !slices.Contains(grammar.symbols, symbol) {
				grammar.symbols = append(grammar.symbols, symbol)
			}
		}
	}

	// 4. Compile all regular expressions
	var err error

	for i, p := range grammar.productions {
		val, ok := p.(*RegProduction)
		if !ok {
			continue
		}

		val.rxp, err = regexp.Compile(val.rhs)
		if err != nil {
			return nil, err
		}

		grammar.productions[i] = val
	}

	// 4. Clear builder
	for i := range b.productions {
		b.productions[i] = nil
	}

	b.productions = nil

	return &grammar, nil
}

// Reset is a method of GrammarBuilder that resets a GrammarBuilder.
func (b *GrammarBuilder) Reset() {
	for i := range b.productions {
		b.productions[i] = nil
	}

	b.productions = nil
}
