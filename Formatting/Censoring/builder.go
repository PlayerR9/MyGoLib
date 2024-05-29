// Package censoring provides utilities for filtering strings.
package Censoring

import (
	"path/filepath"
	"slices"
	"strings"

	uc "github.com/PlayerR9/MyGoLib/Units/Common"
	slext "github.com/PlayerR9/MyGoLib/Units/Slices"
)

// Builder is a type that provides a fluent interface for constructing a
// censoring operation.
type Builder struct {
	// The FilterFuncs used to determine whether each value should be censored
	filters []slext.PredicateFilter[string]

	// The strings to be censored
	values []string

	// The string used to replace censored values
	label string

	// The Context for the censoring operation
	ctx *Context

	// true if the string is not censored
	isNotCensored CensorValue
}

// String is a method of the Builder type that returns a the censored or uncensored
// content of the Builder as a string.
// The value returned depends on the censoring status of the Builder.
//
// Returns:
//   - string: The censored or uncensored content of the Builder.
func (b *Builder) String() string {
	if !b.IsCensored() {
		return filepath.Join(b.values...)
	}

	censoredValues := make([]string, 0, len(b.values))

	for _, word := range b.values {
		filterLower := func(filter slext.PredicateFilter[string]) bool {
			return filter(strings.ToLower(word))
		}

		if slices.ContainsFunc(b.filters, filterLower) {
			if b.label != "" {
				censoredValues = append(censoredValues, b.label)
			} else {
				censoredValues = append(censoredValues, b.ctx.censorLabel)
			}
		} else {
			censoredValues = append(censoredValues, word)
		}
	}

	return filepath.Join(censoredValues...)
}

// BuilderOption is a function type that modifies the properties of a Builder.
type BuilderOption func(*Builder)

// WithLabel returns a BuilderOption that sets the label of a Builder.
// Empty strings will be ignored and the label will not be set.
//
// Parameters:
//   - label: The label to be set.
//
// Returns:
//   - BuilderOption: A function that modifies the label of a Builder.
func WithLabel(label string) BuilderOption {
	return func(b *Builder) {
		if label != "" {
			b.label = label
		}
	}
}

// WithMode returns a BuilderOption that sets the censorship mode of a Builder.
//
// Parameters:
//   - mode: The mode to be set.
//
// Returns:
//   - BuilderOption: A function that modifies the censorship mode of a Builder.
func WithMode(mode CensorValue) BuilderOption {
	return func(b *Builder) {
		b.isNotCensored = !mode
	}
}

// WithFilters returns a BuilderOption that sets the filters of a Builder.
//
// Parameters:
//   - filters: The filters to be set.
//
// Returns:
//   - BuilderOption: A function that modifies the filters of a Builder.
func WithFilters(filters ...slext.PredicateFilter[string]) BuilderOption {
	return func(b *Builder) {
		b.filters = filters
	}
}

// WithValues returns a BuilderOption that sets the values of a Builder.
//
// Parameters:
//   - values: The values to be set.
//
// Returns:
//   - BuilderOption: A function that modifies the values of a Builder.
//
// Important:
//   - Values can be of any type and are converted to strings before being set.
//     If a value is a Builder, its output string is used as the value.
func WithValues(values ...any) BuilderOption {
	return func(b *Builder) {
		stringValues := make([]string, 0, len(values))

		for _, value := range values {
			var str string

			switch x := value.(type) {
			case *Builder:
				if x == nil {
					stringValues = append(stringValues, "")
					continue
				}

				// Partial application to get the
				// uncensored string representation
				x.Apply(func(s string) { str = s })
			default:
				str = uc.StringOf(value)
			}

			stringValues = append(stringValues, str)
		}

		b.values = stringValues
	}
}

// Make is a method of Context that creates a new Builder with the given BuilderOptions.
//
// Parameters:
//   - options: The BuilderOptions to be applied to the Builder.
//
// Returns:
//   - *Builder: The new Builder with the given options.
//
// Important:
//   - The BuilderOptions are applied in the order they are provided. By default, the
//     Builder's label is set to DefaultCensorLabel and separator is set to a space
//     character.
func (ctx *Context) Make(options ...BuilderOption) *Builder {
	builder := &Builder{
		ctx:     ctx,
		label:   DefaultCensorLabel,
		filters: make([]slext.PredicateFilter[string], 0),
		values:  make([]string, 0),
	}

	for _, option := range options {
		option(builder)
	}

	return builder
}

// IsCensored is a method of the Builder type that returns a CensorValue
// indicating whether the Builder is censored.
//
// Returns:
//   - CensorValue: Censored if the Builder is censored, and NotCensored otherwise.
func (b *Builder) IsCensored() CensorValue {
	if b.ctx == nil || !b.ctx.notCensored || !b.isNotCensored {
		return Censored
	}

	return NotCensored
}

// Apply is a method of the Builder type that applies a given function to
// the non-censored string representation of the Builder.
// However, because this allows you to perform operations on the non-censored
// string, regardless of the current censor level, it is important to use it
// with caution and only when necessary.
//
// Parameters:
//   - f: The function to be applied to the non-censored string.
func (b *Builder) Apply(f func(s string)) {
	path := filepath.Join(b.values...)

	f(path)
}
