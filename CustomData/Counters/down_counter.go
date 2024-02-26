// Package counter provides a set of interfaces and methods for managing counters.
package Counters

import (
	"errors"
	"fmt"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

// DownCounter represents a counter that decrements downwards until it reaches zero.
type DownCounter struct {
	// startingCount is the initial value of the counter.
	startingCount int

	// currentCount is the current value of the counter.
	currentCount int

	// retreatCount is the number of times the counter has been retreated.
	retreatCount int
}

// NewDownCounter creates a new DownCounter with the specified starting count.
// It panics with an error of type *ers.ErrInvalidParameter if the starting count
// is less than 0.
//
// Parameters:
//
//   - startingCount is the initial value of the counter.
//
// Returns:
//
//   - The newly created DownCounter.
func NewDownCounter(startingCount int) (DownCounter, error) {
	if startingCount < 0 {
		return DownCounter{}, ers.NewErrInvalidParameter("startingCount").
			Wrap(fmt.Errorf("value (%d) must be positive", startingCount))
	}

	return DownCounter{startingCount, startingCount, 0}, nil
}

// IsDone checks if the DownCounter has reached zero.
//
// Returns:
//
//   - true if the current count is less than or equal to zero.
//     false otherwise.
func (c *DownCounter) IsDone() bool {
	return c.currentCount-c.retreatCount <= 0
}

// Advance decrements the current count of the DownCounter by one.
// It panics with an error of type *ers.ErrCallFailed if the current
// count is already at or below zero.
func (c *DownCounter) Advance() error {
	if c.currentCount-c.retreatCount <= 0 {
		return errors.New("current count is already at or below zero")
	}

	c.currentCount--

	return nil
}

// Retreat increments the retrat count of the DownCounter by one and, as a result,
// decrements the current count by one.
// It panics with an error of type *ers.ErrCallFailed if the starting
// count or current count is already at or below zero.
func (c *DownCounter) Retreat() error {
	if c.currentCount-c.retreatCount <= 0 {
		return errors.New("current count is already at or below zero")
	}

	c.retreatCount++

	return nil
}

// GetRetreatCount returns the number of times the DownCounter
// has been retreated.
//
// Returns:
//
//   - An integer representing the number of times the counter has been retreated.
func (c *DownCounter) GetRetreatCount() int {
	return c.retreatCount
}

// GetDistance returns the current count of the DownCounter.
// This is equivalent to the distance from zero, as the DownCounter
// decrements towards zero.
//
// Returns:
//
//   - An integer representing the current count.
func (c *DownCounter) GetDistance() int {
	return c.currentCount - c.retreatCount
}

// GetCurrentCount returns the current count of the DownCounter.
//
// Returns:
//
//   - An integer representing the current count.
func (c *DownCounter) GetCurrentCount() int {
	return c.currentCount
}

// GetInitialCount returns the starting count of the DownCounter,
// which is the initial count.
//
// Returns:
//
//   - An integer representing the initial count.
func (c *DownCounter) GetInitialCount() int {
	return c.startingCount
}

// String returns a string representation of the DownCounter.
// The string includes the starting count, current count, retreat
// count, and whether the counter is done.
//
// Returns:
//
//   - A string representing the DownCounter.
func (c *DownCounter) String() string {
	return fmt.Sprintf("DownCounter[startingCount=%d, currentCount=%d, retreatCount=%d, isDone=%t]",
		c.startingCount, c.currentCount, c.retreatCount, c.IsDone())
}

// Reset resets the DownCounter to its initial state.
func (c *DownCounter) Reset() {
	c.currentCount = c.startingCount
	c.retreatCount = 0
}
