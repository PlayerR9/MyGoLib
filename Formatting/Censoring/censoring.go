// Package censoring provides utilities for filtering strings.
package Censoring

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
//   - Context: A new Context with default values.
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
//   - label: The label to be set.
//
// Returns:
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
//   - mode: The mode to be set.
//
// Returns:
//   - *Context: A pointer to the Context itself.
func (ctx *Context) WithMode(mode CensorValue) *Context {
	ctx.notCensored = !mode

	return ctx
}
