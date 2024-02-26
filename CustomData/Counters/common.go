// Package counter provides a set of interfaces and methods for managing counters.
package Counters

import (
	"fmt"
)

// Counter interface defines the methods that a counter must implement.
// A counter is a tool that can be advanced or retreated, and it keeps track of
// its state.
type Counter interface {
	// IsDone checks if the counter has reached its limit.
	IsDone() bool

	// Advance advances the counter by one step. It panics if the counter
	// is already at its limit.
	Advance() error

	// Retreat retreats the counter by one step. It panics if the counter
	// is already at its minimum limit.
	Retreat() error

	// GetRetreatCount returns the number of times the counter has been retreated.
	GetRetreatCount() int

	// GetDistance returns the distance from the initial count to the current count.
	GetDistance() int

	// GetCurrentCount returns the current count.
	GetCurrentCount() int

	// GetInitialCount returns the initial count when the counter was first created.
	GetInitialCount() int

	// String returns a string representation of the counter.
	fmt.Stringer

	// DeepCopy creates a deep copy of the counter.
	DeepCopy() Counter

	// Reset resets the counter to its initial state.
	Reset()
}
