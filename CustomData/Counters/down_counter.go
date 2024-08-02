package Counters

import (
	"strconv"
	"strings"
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

// GoString implements the fmt.GoStringer interface.
func (c DownCounter) GoString() string {
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

// NewDownCounter creates a new DownCounter with the specified starting count.
//
// Parameters:
//   - startingCount: The initial value of the counter.
//
// Returns:
//   - DownCounter: A new DownCounter with the specified starting count.
//
// If the startingCount is less than 0, it is set to 0.
func NewDownCounter(startingCount int) DownCounter {
	if startingCount < 0 {
		startingCount = 0
	}

	return DownCounter{startingCount, startingCount, 0}
}

// IsDone checks if the DownCounter has reached zero.
//
// Returns:
//   - bool: true if the counter has reached zero, false otherwise.
func (c DownCounter) IsDone() bool {
	return c.currentCount-c.retreatCount <= 0
}

// Advance decrements the current count of the DownCounter by one.
//
// Returns:
//   - DownCounter: A new DownCounter with the decremented current count.
//   - bool: true if the counter has not reached zero, false otherwise.
func (c DownCounter) Advance() (DownCounter, bool) {
	if c.currentCount-c.retreatCount <= 0 {
		return c, false
	}

	c.currentCount--

	return c, true
}

// Retreat increments the retrat count of the DownCounter by one and, as a
// result, decrements the current count by one.
//
// Returns:
//   - DownCounter: A new DownCounter with the decremented current count.
//   - bool: true if the counter has not reached zero, false otherwise.
func (c DownCounter) Retreat() (DownCounter, bool) {
	if c.currentCount-c.retreatCount <= 0 {
		return c, false
	}

	c.retreatCount++

	return c, true
}

// RetreatCount returns the number of times the DownCounter
// has been retreated.
//
// Returns:
//   - int: The number of times the DownCounter has been retreated.
func (c DownCounter) RetreatCount() int {
	return c.retreatCount
}

// Distance returns the current count of the DownCounter.
// This is equivalent to the distance from zero, as the DownCounter
// decrements towards zero.
//
// Returns:
//   - int: The current count of the DownCounter.
func (c DownCounter) Distance() int {
	return c.currentCount - c.retreatCount
}

// CurrentCount returns the current count of the DownCounter.
//
// Returns:
//   - int: The current count of the DownCounter.
func (c DownCounter) CurrentCount() int {
	return c.currentCount
}

// InitialCount returns the starting count of the DownCounter,
// which is the initial count.
//
// Returns:
//   - int: The starting count of the DownCounter.
func (c DownCounter) InitialCount() int {
	return c.startingCount
}

// Reset resets the DownCounter to its initial state.
func (c DownCounter) Reset() DownCounter {
	c.currentCount = c.startingCount
	c.retreatCount = 0

	return c
}
