// Package censoring provides utilities for filtering strings.
package Censoring

import (
	"fmt"
	"slices"
	"strings"
)

// FilterFunc is a function type that defines the filtering criteria for a string,
// based on which the string is either accepted or rejected.
//
// Parameters:
//
//   - string: The string to be filtered.
//
// Returns:
//
//   - bool: A boolean value that is true if the string is rejected, and false otherwise.
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

const (
	// Censored represents a string that has been censored
	Censored CensorValue = true

	// NotCensored represents a string that has not been censored
	NotCensored CensorValue = false
)

// Context is a type that encapsulates the context for a censoring operation.
type Context struct {
	// The label used to replace unacceptable strings
	censorLabel string

	// true if the string is not censored
	notCensored CensorValue
}

// NewContext creates a new Context with default values.
// The default censorLabel is DefaultCensorLabel and the context is censored.
//
// Returns:
//
//   - *Context: A pointer to the newly created Context.
func NewContext() *Context {
	return &Context{
		censorLabel: DefaultCensorLabel,
		notCensored: false,
	}
}

// WithLabel sets the censorLabel of the Context to the given label.
// Empty strings will be replaced with DefaultCensorLabel.
//
// Parameters:
//
//   - label: The label to be set.
//
// Returns:
//
//   - *Context: A pointer to the Context itself.
func (ctx *Context) WithLabel(label string) *Context {
	if label != "" {
		ctx.censorLabel = label
	} else {
		ctx.censorLabel = DefaultCensorLabel
	}

	return ctx
}

// WithMode sets the notCensored value of the Context to the given mode.
//
// Parameters:
//
//   - mode: The mode to be set.
//
// Returns:
//
//   - *Context: A pointer to the Context itself.
func (ctx *Context) WithMode(mode CensorValue) *Context {
	ctx.notCensored = !mode

	return ctx
}

// Builder is a type that provides a fluent interface for constructing a
// censoring operation.
type Builder struct {
	// The FilterFuncs used to determine whether each value should be censored
	filters []FilterFunc

	// The string used to join the values together
	sep rune

	// The strings to be censored
	values []string

	// The string used to replace censored values
	label string

	// The Context for the censoring operation
	ctx *Context

	// true if the string is not censored
	isNotCensored CensorValue
}

// BuilderOption is a function type that modifies the properties of a Builder.
type BuilderOption func(*Builder)

// WithLabel returns a BuilderOption that sets the label of a Builder.
// Empty strings will be ignored and the label will not be set.
//
// Parameters:
//
//   - label: The label to be set.
//
// Returns:
//
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
//
//   - mode: The mode to be set.
//
// Returns:
//
//   - BuilderOption: A function that modifies the censorship mode of a Builder.
func WithMode(mode CensorValue) BuilderOption {
	return func(b *Builder) {
		b.isNotCensored = !mode
	}
}

// WithSeparator returns a BuilderOption that sets the separator of a Builder.
//
// Parameters:
//
//   - sep: The separator to be set.
//
// Returns:
//
//   - BuilderOption: A function that modifies the separator of a Builder.
func WithSeparator(sep rune) BuilderOption {
	return func(b *Builder) {
		b.sep = sep
	}
}

// WithFilters returns a BuilderOption that sets the filters of a Builder.
//
// Parameters:
//
//   - filters: The filters to be set.
//
// Returns:
//
//   - BuilderOption: A function that modifies the filters of a Builder.
func WithFilters(filters ...FilterFunc) BuilderOption {
	return func(b *Builder) {
		b.filters = filters
	}
}

// WithValues returns a BuilderOption that sets the values of a Builder.
//
// Parameters:
//
//   - values: The values to be set.
//
// Returns:
//
//   - BuilderOption: A function that modifies the values of a Builder.
//
// Important:
//
// Values can be of any type and are converted to strings before being set.
// If a value is a Builder, its output string is used as the value.
func WithValues(values ...any) BuilderOption {
	return func(b *Builder) {
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

		b.values = stringValues
	}
}

// Make is a method of Context that creates a new Builder with the given BuilderOptions.
//
// Parameters:
//
//   - options: The BuilderOptions to be applied to the Builder.
//
// Returns:
//
//   - *Builder: A pointer to the newly created Builder.
//
// Important:
//
// The BuilderOptions are applied in the order they are provided. By default, the
// Builder's label is set to DefaultCensorLabel and separator is set to a space
// character.
func (ctx *Context) Make(options ...BuilderOption) *Builder {
	builder := &Builder{
		ctx:     ctx,
		label:   DefaultCensorLabel,
		filters: make([]FilterFunc, 0),
		values:  make([]string, 0),
		sep:     ' ',
	}

	for _, option := range options {
		option(builder)
	}

	var fields []string

	for i := 0; i < len(builder.values); i += len(fields) {
		fields = strings.FieldsFunc(builder.values[i], func(r rune) bool {
			return r == rune(builder.sep)
		})

		builder.values = append(builder.values[:i], append(fields, builder.values[i+1:]...)...)
	}

	return builder
}

// String is a method of the Builder type that returns a the censored or uncensored
// content of the Builder as a string.
// The value returned depends on the censoring status of the Builder.
//
// Returns:
//
//   - string: The censored or uncensored content of the Builder.
func (b *Builder) String() string {
	if !b.IsCensored() {
		return strings.Join(b.values, string(b.sep))
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

	return strings.Join(censoredValues, string(b.sep))
}

// IsCensored is a method of the Builder type that returns a CensorValue
// indicating whether the Builder is censored.
//
// Returns:
//
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
//
//   - f: The function to be applied to the non-censored string.
func (b *Builder) Apply(f func(s string)) {
	f(strings.Join(b.values, string(b.sep)))
}
