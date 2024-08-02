package Counters

import (
	"strconv"
	"strings"
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

// GoString implements the fmt.GoStringer interface.
func (c UpCounter) GoString() string {
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

// NewUpCounter creates a new UpCounter with the specified upper limit.
//
// Parameters:
//   - upperLimit: The maximum limit that the counter can reach.
//
// Returns:
//   - UpCounter: A new UpCounter with the specified upper limit.
//
// If the upperLimit is less than 0, it is set to 0.
func NewUpCounter(upperLimit int) UpCounter {
	if upperLimit < 0 {
		upperLimit = 0
	}

	return UpCounter{upperLimit, 0, 0}
}

// IsDone checks if the UpCounter has reached its upper limit.
//
// Returns:
//   - bool: true if the counter has reached its upper limit, false
//     otherwise.
func (c UpCounter) IsDone() bool {
	return c.currentCount >= c.upperLimit-c.retreatCount
}

// Advance increments the current count of the UpCounter by one.
//
// Returns:
//   - UpCounter: A new UpCounter with the incremented current count.
//   - bool: true if the counter has not reached its upper limit, false
//     otherwise.
func (c UpCounter) Advance() (UpCounter, bool) {
	if c.currentCount >= c.upperLimit-c.retreatCount {
		return c, false
	}

	c.currentCount++

	return c, true
}

// Retreat increments the retreat count and, as a result, decrements the
// upper limit of the UpCounter by one.
//
// Returns:
//   - UpCounter: A new UpCounter with the decremented upper limit.
//   - bool: true if the counter has not reached its upper limit, false
//     otherwise.
func (c UpCounter) Retreat() (UpCounter, bool) {
	if c.currentCount >= c.upperLimit-c.retreatCount {
		return c, false
	}

	c.retreatCount++

	return c, true
}

// RetreatCount returns the number of times the UpCounter has
// been retreated.
//
// Returns:
//   - int: The number of times the counter has been retreated.
func (c UpCounter) RetreatCount() int {
	return c.retreatCount
}

// Distance calculates the distance between the current count and
// the upper limit of the UpCounter, that is, the number of times the
// counter can still be advanced before reaching the upper limit.
//
// Returns:
//   - int: The distance between the current count and the upper limit.
func (c UpCounter) Distance() int {
	return c.upperLimit - c.retreatCount - c.currentCount
}

// CurrentCount returns the current count of the UpCounter.
//
// Returns:
//   - int: The current count of the UpCounter.
func (c UpCounter) CurrentCount() int {
	return c.currentCount
}

// InitialCount returns the upper limit of the UpCounter,
// which is the initial count.
//
// Returns:
//   - int: The upper limit of the UpCounter.
func (c UpCounter) InitialCount() int {
	return c.upperLimit
}

// Reset resets the UpCounter to its initial state.
func (c UpCounter) Reset() UpCounter {
	c.currentCount = 0
	c.retreatCount = 0

	return c
}
