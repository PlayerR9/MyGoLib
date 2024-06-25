package Counters

import (
	"strconv"
	"strings"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// DownCounter represents a counter that decrements downwards until it
// reaches zero.
type DownCounter struct {
	// startingCount is the initial value of the counter.
	startingCount int

	// currentCount is the current value of the counter.
	currentCount int

	// retreatCount is the number of times the counter has been retreated.
	retreatCount int
}

// Equals implements common.Objecter.
func (c *DownCounter) Equals(other uc.Equaler) bool {
	if other == nil {
		return false
	}

	otherC, ok := other.(*DownCounter)
	if !ok {
		return false
	}

	return c.startingCount == otherC.startingCount &&
		c.currentCount == otherC.currentCount &&
		c.retreatCount == otherC.retreatCount
}

// IsDone checks if the DownCounter has reached zero.
//
// Returns:
//   - bool: true if the counter has reached zero, false otherwise.
func (c *DownCounter) IsDone() bool {
	return c.currentCount-c.retreatCount <= 0
}

// Advance decrements the current count of the DownCounter by one.
//
// Returns:
//   - error: An error of type *ErrCurrentCountBelowZero if the current
//     count is already at or below zero.
func (c *DownCounter) Advance() error {
	if c.currentCount-c.retreatCount <= 0 {
		return NewErrCurrentCountBelowZero()
	}

	c.currentCount--

	return nil
}

// Retreat increments the retrat count of the DownCounter by one and, as a
// result, decrements the current count by one.
//
// Returns:
//   - error: An error of type *ErrCurrentCountBelowZero if the current count
//     is already at or below zero.
func (c *DownCounter) Retreat() error {
	if c.currentCount-c.retreatCount <= 0 {
		return NewErrCurrentCountBelowZero()
	}

	c.retreatCount++

	return nil
}

// GetRetreatCount returns the number of times the DownCounter
// has been retreated.
//
// Returns:
//   - int: The number of times the DownCounter has been retreated.
func (c *DownCounter) GetRetreatCount() int {
	return c.retreatCount
}

// GetDistance returns the current count of the DownCounter.
// This is equivalent to the distance from zero, as the DownCounter
// decrements towards zero.
//
// Returns:
//   - int: The current count of the DownCounter.
func (c *DownCounter) GetDistance() int {
	return c.currentCount - c.retreatCount
}

// GetCurrentCount returns the current count of the DownCounter.
//
// Returns:
//   - int: The current count of the DownCounter.
func (c *DownCounter) GetCurrentCount() int {
	return c.currentCount
}

// GetInitialCount returns the starting count of the DownCounter,
// which is the initial count.
//
// Returns:
//   - int: The starting count of the DownCounter.
func (c *DownCounter) GetInitialCount() int {
	return c.startingCount
}

// String returns a string representation of the DownCounter.
// The string includes the starting count, current count, retreat
// count, and whether the counter is done.
//
// This is used for debugging and logging purposes.
//
// Returns:
//   - string: A string representation of the DownCounter.
func (c *DownCounter) String() string {
	var builder strings.Builder

	builder.WriteString("DownCounter[startingCount=")
	builder.WriteString(strconv.Itoa(c.startingCount))
	builder.WriteString(", currentCount=")
	builder.WriteString(strconv.Itoa(c.currentCount))
	builder.WriteString(", retreatCount=")
	builder.WriteString(strconv.Itoa(c.retreatCount))
	builder.WriteString(", isDone=")
	builder.WriteString(strconv.FormatBool(c.IsDone()))
	builder.WriteRune(']')

	return builder.String()
}

// Reset resets the DownCounter to its initial state.
func (c *DownCounter) Reset() {
	c.currentCount = c.startingCount
	c.retreatCount = 0
}

// Copy creates a shallow copy of the DownCounter.
//
// Returns:
//   - uc.Copier: A shallow copy of the DownCounter.
func (c *DownCounter) Copy() uc.Copier {
	return &DownCounter{c.startingCount, c.currentCount, c.retreatCount}
}

// NewDownCounter creates a new DownCounter with the specified starting count.
//
// Parameters:
//   - startingCount: The initial value of the counter.
//
// Returns:
//   - *DownCounter: A pointer to the new DownCounter.
//   - error: An error of type *uc.ErrInvalidParameter if the starting count is
//     less than zero.
func NewDownCounter(startingCount int) (*DownCounter, error) {
	if startingCount < 0 {
		return nil, uc.NewErrInvalidParameter(
			"startingCount",
			uc.NewErrGTE(0),
		)
	}

	return &DownCounter{startingCount, startingCount, 0}, nil
}
