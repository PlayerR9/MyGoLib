package Counters

import (
	"strconv"
	"strings"

	intf "github.com/PlayerR9/MyGoLib/Units/common"
	ers "github.com/PlayerR9/MyGoLib/Units/errors"
)

// UpCounter represents a counter that increments upwards until it
// reaches an upper limit.
type UpCounter struct {
	// upperLimit is the maximum limit that the counter can reach.
	upperLimit int

	// currentCount is the current value of the counter.
	currentCount int

	// retreatCount is the number of times the counter has been retreated.
	retreatCount int
}

// Equals implements common.Objecter.
func (c *UpCounter) Equals(other intf.Equaler) bool {
	if other == nil {
		return false
	}

	otherC, ok := other.(*UpCounter)
	if !ok {
		return false
	}

	return c.upperLimit == otherC.upperLimit &&
		c.currentCount == otherC.currentCount &&
		c.retreatCount == otherC.retreatCount
}

// IsDone checks if the UpCounter has reached its upper limit.
//
// Returns:
//   - bool: true if the counter has reached its upper limit, false
//     otherwise.
func (c *UpCounter) IsDone() bool {
	return c.currentCount >= c.upperLimit-c.retreatCount
}

// Advance increments the current count of the UpCounter by one.
//
// Returns:
//   - error: An error of type *ErrCurrentCountAboveUpperLimit if the current
//     count is already at or beyond the upper limit.
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
//   - error: An error of type *ErrCurrentCountAboveUpperLimit if the current
//     count is already at or beyond the upper limit.
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
//   - int: The number of times the counter has been retreated.
func (c *UpCounter) GetRetreatCount() int {
	return c.retreatCount
}

// GetDistance calculates the distance between the current count and
// the upper limit of the UpCounter, that is, the number of times the
// counter can still be advanced before reaching the upper limit.
//
// Returns:
//   - int: The distance between the current count and the upper limit.
func (c *UpCounter) GetDistance() int {
	return c.upperLimit - c.retreatCount - c.currentCount
}

// GetCurrentCount returns the current count of the UpCounter.
//
// Returns:
//   - int: The current count of the UpCounter.
func (c *UpCounter) GetCurrentCount() int {
	return c.currentCount
}

// GetInitialCount returns the upper limit of the UpCounter,
// which is the initial count.
//
// Returns:
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
//   - string: A string representation of the UpCounter.
func (c *UpCounter) String() string {
	var builder strings.Builder

	builder.WriteString("UpCounter[upperLimit=")
	builder.WriteString(strconv.Itoa(c.upperLimit))
	builder.WriteString(", currentCount=")
	builder.WriteString(strconv.Itoa(c.currentCount))
	builder.WriteString(", retreatCount=")
	builder.WriteString(strconv.Itoa(c.retreatCount))
	builder.WriteString(", isDone=")
	builder.WriteString(strconv.FormatBool(c.IsDone()))
	builder.WriteString("]")

	return builder.String()
}

// Reset resets the UpCounter to its initial state.
func (c *UpCounter) Reset() {
	c.currentCount = 0
	c.retreatCount = 0
}

// Copy creates a shallow copy of the UpCounter.
//
// Returns:
//   - intf.Copier: A shallow copy of the UpCounter.
func (c *UpCounter) Copy() intf.Copier {
	return &UpCounter{
		upperLimit:   c.upperLimit,
		currentCount: c.currentCount,
		retreatCount: c.retreatCount,
	}
}

// NewUpCounter creates a new UpCounter with the specified upper limit.
//
// Parameters:
//   - upperLimit: The maximum limit that the counter can reach.
//
// Returns:
//   - *UpCounter: A pointer to the newly created UpCounter.
//   - error: An error of type *ers.ErrInvalidParameter if the upper limit is
//     less than zero.
func NewUpCounter(upperLimit int) (*UpCounter, error) {
	if upperLimit < 0 {
		return nil, ers.NewErrInvalidParameter(
			"upperLimit",
			ers.NewErrGTE(0),
		)
	}

	return &UpCounter{upperLimit, 0, 0}, nil
}
