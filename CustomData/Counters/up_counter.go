// Package counter provides a set of interfaces and methods for managing counters.
package Counters

import (
	"errors"
	"fmt"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

// UpCounter represents a counter that increments upwards until it reaches an upper
// limit.
type UpCounter struct {
	// upperLimit is the maximum limit that the counter can reach.
	upperLimit int

	// currentCount is the current value of the counter.
	currentCount int

	// retreatCount is the number of times the counter has been retreated.
	retreatCount int
}

// NewUpCounter creates a new UpCounter with the specified upper limit.
// It panics with an error of type *ers.ErrInvalidParameter if the upper
// limit is less than 0.
//
// Parameters:
//
//   - upperLimit is the maximum limit that the counter can reach.
//
// Returns:
//
//   - The newly created UpCounter.
func NewUpCounter(upperLimit int) UpCounter {
	if upperLimit >= 0 {
		return UpCounter{upperLimit, 0, 0}
	}

	panic(ers.NewErrInvalidParameter("upperLimit").
		WithReason(fmt.Errorf("value (%d) must be positive", upperLimit)))
}

// IsDone checks if the UpCounter has reached its upper limit.
//
// Returns:
//
//   - true if the current count is greater than or equal to the upper limit.
//     false otherwise.
func (c *UpCounter) IsDone() bool {
	return c.currentCount >= c.upperLimit-c.retreatCount
}

// Advance increments the current count of the UpCounter by one.
// It panics with an error of type *ers.ErrCallFailed if the current
// count is already at or beyond the upper limit.
func (c *UpCounter) Advance() {
	if c.currentCount >= c.upperLimit-c.retreatCount {
		panic(ers.NewErrCallFailed("Advance", c.Advance).
			WithReason(errors.New("current count is already at or beyond the upper limit")),
		)
	}

	c.currentCount++
}

// Retreat increments the retreat count and, as a result, decrements the
// upper limit of the UpCounter by one.
// It panics with an error of type *ers.ErrCallFailed if the current
// count is already at or beyond the upper limit.
func (c *UpCounter) Retreat() {
	if c.currentCount >= c.upperLimit-c.retreatCount {
		panic(ers.NewErrCallFailed("Retreat", c.Retreat).WithReason(
			errors.New("current count is already at or beyond the upper limit"),
		))
	}

	c.retreatCount++
}

// GetRetreatCount returns the number of times the UpCounter has
// been retreated.
//
// Returns:
//
//   - An integer representing the retreat count.
func (c *UpCounter) GetRetreatCount() int {
	return c.retreatCount
}

// GetDistance calculates the distance between the current count and
// the upper limit of the UpCounter, that is, the number of times the
// counter can still be advanced before reaching the upper limit.
//
// Returns:
//
//   - An integer representing the distance.
func (c *UpCounter) GetDistance() int {
	return c.upperLimit - c.retreatCount - c.currentCount
}

// GetCurrentCount returns the current count of the UpCounter.
//
// Returns:
//
//   - An integer representing the current count.
func (c *UpCounter) GetCurrentCount() int {
	return c.currentCount
}

// GetInitialCount returns the upper limit of the UpCounter,
// which is the initial count.
//
// Returns:
//
//   - An integer representing the initial count.
func (c *UpCounter) GetInitialCount() int {
	return c.upperLimit
}

// String returns a string representation of the UpCounter.
// The string includes the upper limit, current count, retreat
// count, and whether the counter is done.
//
// Returns:
//
//   - A string representing the UpCounter.
func (c *UpCounter) String() string {
	return fmt.Sprintf("UpCounter[upperLimit=%d, currentCount=%d, retreatCount=%d, isDone=%t]",
		c.upperLimit, c.currentCount, c.retreatCount, c.IsDone())
}

// Reset resets the UpCounter to its initial state.
func (c *UpCounter) Reset() {
	c.currentCount = 0
	c.retreatCount = 0
}
