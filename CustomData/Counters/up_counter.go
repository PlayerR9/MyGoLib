// Package Counters provides a set of interfaces and methods for managing counters.
package Counters

import (
	"fmt"

	ers "github.com/PlayerR9/MyGoLibUnits/Errors"
	itf "github.com/PlayerR9/MyGoLibUnits/Interfaces"
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
//
// Parameters:
//
//   - upperLimit: The maximum limit that the counter can reach.
//
// Returns:
//
//   - *UpCounter: A pointer to the newly created UpCounter.
//   - error: An error of type *ers.ErrInvalidParameter if the upper limit
//     is less than 0.
func NewUpCounter(upperLimit int) (*UpCounter, error) {
	if upperLimit < 0 {
		return nil, ers.NewErrInvalidParameter(
			"upperLimit",
			fmt.Errorf("value (%d) must be positive", upperLimit),
		)
	}

	return &UpCounter{upperLimit, 0, 0}, nil
}

// IsDone checks if the UpCounter has reached its upper limit.
//
// Returns:
//
//   - bool: true if the counter has reached its upper limit, false otherwise.
func (c *UpCounter) IsDone() bool {
	return c.currentCount >= c.upperLimit-c.retreatCount
}

// Advance increments the current count of the UpCounter by one.
//
// Returns:
//
//   - error: An error of type *ErrCurrentCountAboveUpperLimit if the current count
//     is already at or beyond the upper limit.
func (c *UpCounter) Advance() error {
	if c.currentCount >= c.upperLimit-c.retreatCount {
		return NewErrCurrentCountAboveUpperLimit()
	}

	c.currentCount++

	return nil
}

// Retreat increments the retreat count and, as a result, decrements the
// upper limit of the UpCounter by one.
//
// Returns:
//
//   - error: An error of type *ErrCurrentCountAboveUpperLimit if the current count
//     is already at or beyond the upper limit.
func (c *UpCounter) Retreat() error {
	if c.currentCount >= c.upperLimit-c.retreatCount {
		return NewErrCurrentCountAboveUpperLimit()
	}

	c.retreatCount++

	return nil
}

// GetRetreatCount returns the number of times the UpCounter has
// been retreated.
//
// Returns:
//
//   - int: The number of times the counter has been retreated.
func (c *UpCounter) GetRetreatCount() int {
	return c.retreatCount
}

// GetDistance calculates the distance between the current count and
// the upper limit of the UpCounter, that is, the number of times the
// counter can still be advanced before reaching the upper limit.
//
// Returns:
//
//   - int: The distance between the current count and the upper limit.
func (c *UpCounter) GetDistance() int {
	return c.upperLimit - c.retreatCount - c.currentCount
}

// GetCurrentCount returns the current count of the UpCounter.
//
// Returns:
//
//   - int: The current count of the UpCounter.
func (c *UpCounter) GetCurrentCount() int {
	return c.currentCount
}

// GetInitialCount returns the upper limit of the UpCounter,
// which is the initial count.
//
// Returns:
//
//   - int: The upper limit of the UpCounter.
func (c *UpCounter) GetInitialCount() int {
	return c.upperLimit
}

// String returns a string representation of the UpCounter.
// The string includes the upper limit, current count, retreat
// count, and whether the counter is done.
//
// This should be used for debugging and logging purposes.
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

// Copy creates a shallow copy of the UpCounter.
//
// Returns:
//
//   - itf.Copier: A shallow copy of the UpCounter.
func (c *UpCounter) Copy() itf.Copier {
	return &UpCounter{
		upperLimit:   c.upperLimit,
		currentCount: c.currentCount,
		retreatCount: c.retreatCount,
	}
}
