// Package censoring provides utilities for filtering strings.
package Censoring

import (
	"fmt"
	"slices"
	"strings"
)

// FilterFunc is a function type that defines the filtering criteria for a string.
// It takes a string as input and returns a boolean indicating whether the string
// is acceptable or not.
type FilterFunc func(string) bool

// DefaultCensorLabel is a constant that defines the default label used to replace
// unacceptable strings when they are filtered. Its default value is "[***]".
const (
	DefaultCensorLabel string = "[***]"
)

// CensorValue is a type that represents whether a string has been censored or not.
// It is used in conjunction with the FilterFunc type to determine the censoring
// status of a string.
type CensorValue bool

// Censored and NotCensored are constants of type CensorValue that represent the two
// possible censoring statuses.
// Censored represents a string that has been censored, while NotCensored represents
// a string that has not been censored.
const (
	Censored    CensorValue = true
	NotCensored CensorValue = false
)

// Context is a type that encapsulates the context for a censoring operation.
// It contains a censorLabel, which is the label used to replace unacceptable strings,
// and a notCensored value, which indicates whether the string is censored or not.
type Context struct {
	// The label used to replace unacceptable strings
	censorLabel string

	// true if the string is not censored
	notCensored CensorValue
}

func (ctx *Context) Cleanup() {}

// NewContext creates a new Context with default values.
// The default censorLabel is DefaultCensorLabel and the default notCensored value
// is false.
// It returns a pointer to the newly created Context.
func NewContext() *Context {
	return &Context{
		censorLabel: DefaultCensorLabel,
		notCensored: false,
	}
}

// SetLabel sets the censorLabel of the Context to the given label.
// If the given label is an empty string, the censorLabel is set to DefaultCensorLabel.
// It returns the Context itself to allow for method chaining, following the builder
// pattern.
func (ctx *Context) SetLabel(label string) *Context {
	if label != "" {
		ctx.censorLabel = label
	} else {
		ctx.censorLabel = DefaultCensorLabel
	}

	return ctx
}

// CensorMode sets the notCensored value of the Context to the negation of the given mode.
// It returns the Context itself to allow for method chaining, following the builder
// pattern.
func (ctx *Context) CensorMode(mode CensorValue) *Context {
	ctx.notCensored = !mode

	return ctx
}

// Builder is a type that provides a fluent interface for constructing a
// censoring operation.
// It contains a slice of FilterFuncs, a separator string, a slice of values
// to be censored,
// a censor label, a Context, and a CensorValue indicating whether the
// string is not censored.
type Builder struct {
	// The FilterFuncs used to determine whether each value should be censored
	filters []FilterFunc

	// The string used to join the values together
	sep string

	// The strings to be censored
	values []string

	// The string used to replace censored values
	label string

	// The Context for the censoring operation
	ctx *Context

	// true if the string is not censored
	isNotCensored CensorValue
}

func (b *Builder) Cleanup() {
	if b.filters != nil {
		for i := 0; i < len(b.filters); i++ {
			b.filters[i] = nil
		}

		b.filters = nil
	}

	b.values = nil

	if b.ctx != nil {
		b.ctx.Cleanup()
		b.ctx = nil
	}
}

// Make creates a new Builder with the given filters, separator, and values.
// It initializes the Builder's Context to the given Context, sets its
// isNotCensored value to false,
// and converts the given values to strings.
// If the given separator is not an empty string, it splits each value by the
// separator and appends the resulting fields to the values.
// It returns a pointer to the newly created Builder.
func (ctx *Context) Make(filters []FilterFunc, sep string, values ...any) *Builder {
	builder := new(Builder)
	builder.ctx = ctx
	builder.isNotCensored = false

	stringValues := make([]string, 0, len(values))

	for _, value := range values {
		switch x := value.(type) {
		case *Builder:
			if x == nil {
				stringValues = append(stringValues, "")
				continue
			}

			// Partial application to get the
			// uncensored string representation
			var str string

			x.Apply(func(s string) { str = s })

			stringValues = append(stringValues, str)
		case Builder:
			// Partial application to get the
			// uncensored string representation
			var str string

			x.Apply(func(s string) { str = s })

			stringValues = append(stringValues, str)
		default:
			stringValues = append(stringValues, fmt.Sprintf("%v", value))
		}
	}

	if sep != "" {
		var fields []string

		for i := 0; i < len(stringValues); i += len(fields) {
			fields = strings.Split(stringValues[i], sep)
			stringValues = append(stringValues[:i], append(fields, stringValues[i+1:]...)...)
		}
	}

	builder.filters = filters
	builder.sep = sep
	builder.label = ""
	builder.values = stringValues

	return builder
}

// SetLabel sets the label of the Builder to the given label.
// If the given label is an empty string, it sets the label to the censorLabel
// of the Builder's Context.
// It returns the Builder itself to allow for method chaining, following the
// builder pattern.
func (b *Builder) SetLabel(label string) *Builder {
	if label != "" {
		b.label = label
	} else {
		b.label = b.ctx.censorLabel
	}

	return b
}

// CensorMode sets the isNotCensored value of the Builder to the negation
// of the given mode.
// It returns the Builder itself to allow for method chaining, following
// the builder pattern.
func (b *Builder) CensorMode(mode CensorValue) *Builder {
	b.isNotCensored = !mode

	return b
}

// String is a method of the Builder type that returns a string
// representation of the Builder.
// If the Builder is not censored, it joins the values with the
// separator and returns the resulting string.
// If the Builder is censored, it creates a new slice of censored
// values, checks each value against the filters, and replaces
// unacceptable values with the censor label. It then joins the
// censored values with the separator and returns the resulting string.
func (b Builder) String() string {
	if !b.IsCensored() {
		return strings.Join(b.values, b.sep)
	}

	censoredValues := make([]string, 0, len(b.values))

	for _, word := range b.values {
		if slices.ContainsFunc(b.filters, func(filter FilterFunc) bool {
			return filter(strings.ToLower(word))
		}) {
			if b.label != "" {
				censoredValues = append(censoredValues, b.label)
			} else {
				censoredValues = append(censoredValues, b.ctx.censorLabel)
			}
		} else {
			censoredValues = append(censoredValues, word)
		}
	}

	return strings.Join(censoredValues, b.sep)
}

// IsCensored is a method of the Builder type that returns a CensorValue
// indicating whether the Builder is censored.
// If the Builder's Context is nil, or if either the Context's notCensored
// value or the Builder's isNotCensored value is false, it returns Censored.
// Otherwise, it returns NotCensored.
func (b *Builder) IsCensored() CensorValue {
	if b.ctx == nil || !b.ctx.notCensored || !b.isNotCensored {
		return Censored
	}

	return NotCensored
}

// Apply is a method of the Builder type that applies a given function to
// the non-censored string representation of the Builder.
// It joins the Builder's values with the separator to create the non-censored
// string, and then passes this string to the given function.
// This method allows you to perform operations on the non-censored string,
// regardless of the current censor level.
func (b *Builder) Apply(f func(s string)) {
	f(strings.Join(b.values, b.sep))
}
