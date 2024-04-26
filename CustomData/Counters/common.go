// Package Counters provides a set of interfaces and methods for
// managing counters.
package Counters

import (
	"fmt"

	intf "github.com/PlayerR9/MyGoLib/Units/Interfaces"
)

// Counter interface defines the methods that a counter must implement.
// A counter is a tool that can be advanced or retreated, and it keeps
// track of its state.
type Counter interface {
	// IsDone checks if the counter has reached its limit.
	//
	// Returns:
	//   - bool: true if the counter has reached its limit, false otherwise.
	IsDone() bool

	// Advance advances the counter by one step.
	//
	// Returns:
	//   - error: An error if the counter is already at its maximum limit.
	Advance() error

	// Retreat retreats the counter by one step.
	//
	// Returns:
	//   - error: An error if the counter is already at its minimum limit.
	Retreat() error

	// GetRetreatCount returns the number of times the counter has been
	// retreated.
	//
	// Returns:
	//   - int: The number of times the counter has been retreated.
	GetRetreatCount() int

	// GetDistance returns the distance from the initial count to the current
	// count.
	//
	// Returns:
	//   - int: The distance from the initial count to the current count.
	GetDistance() int

	// GetCurrentCount returns the current count.
	//
	// Returns:
	//   - int: The current count.
	GetCurrentCount() int

	// GetInitialCount returns the initial count when the counter was first
	// created.
	//
	// Returns:
	//   - int: The initial count.
	GetInitialCount() int

	// Reset resets the counter to its initial state.
	Reset()

	fmt.Stringer

	intf.Copier
}
