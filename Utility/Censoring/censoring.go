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

// WithLabel sets the censorLabel of the Context to the given label.
// If the given label is an empty string, the censorLabel is set to DefaultCensorLabel.
// It returns the Context itself to allow for method chaining, following the builder
// pattern.
func (ctx *Context) WithLabel(label string) *Context {
	label = strings.TrimSpace(label)
	if label != "" {
		ctx.censorLabel = label
	} else {
		ctx.censorLabel = DefaultCensorLabel
	}

	return ctx
}

// WithMode sets the notCensored value of the Context to the negation of the given mode.
// It returns the Context itself to allow for method chaining, following the builder
// pattern.
func (ctx *Context) WithMode(mode CensorValue) *Context {
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
// If the provided label is not an empty string, it replaces the current label
// of the Builder.
// The label is trimmed of any leading or trailing white space before being set.
func WithLabel(label string) BuilderOption {
	return func(b *Builder) {
		label = strings.TrimSpace(label)
		if label != "" {
			b.label = label
		}
	}
}

// WithMode returns a BuilderOption that sets the censorship mode of a Builder.
// If the provided mode is true, the Builder will not censor values.
func WithMode(mode CensorValue) BuilderOption {
	return func(b *Builder) {
		b.isNotCensored = !mode
	}
}

// WithSeparator returns a BuilderOption that sets the separator of a Builder.
// The separator is used to separate values in the output string.
func WithSeparator(sep rune) BuilderOption {
	return func(b *Builder) {
		b.sep = sep
	}
}

// WithFilters returns a BuilderOption that sets the filters of a Builder.
// Filters are functions that modify the output string of the Builder.
func WithFilters(filters ...FilterFunc) BuilderOption {
	return func(b *Builder) {
		b.filters = filters
	}
}

// WithValues returns a BuilderOption that sets the values of a Builder.
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

// Make creates a new Builder with the given BuilderOptions.
// It initializes the Builder's Context to the given Context, sets its
// label to DefaultCensorLabel, and initializes its filters and values
// to empty slices.
// Each provided BuilderOption is then applied to the Builder in the order
// they were provided.
// After applying the BuilderOptions, it splits each value in the Builder's
// values slice by the Builder's separator, and replaces the original values
// with the resulting fields.
// It returns a pointer to the newly created Builder.
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

// String is a method of the Builder type that returns a string
// representation of the Builder.
// If the Builder is not censored, it joins the values with the
// separator and returns the resulting string.
// If the Builder is censored, it creates a new slice of censored
// values, checks each value against the filters, and replaces
// unacceptable values with the censor label. It then joins the
// censored values with the separator and returns the resulting string.
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
	f(strings.Join(b.values, string(b.sep)))
}
