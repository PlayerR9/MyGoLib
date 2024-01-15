// Package counter provides a set of interfaces and methods for managing counters.
package Counters

import (
	"errors"
	"fmt"
	"strings"
)

// Counter interface defines the methods that a counter must implement.
// A counter is a tool that can be advanced or retreated, and it keeps track of
// its state.
type Counter interface {
	// IsDone checks if the counter has reached its limit.
	IsDone() bool

	// Advance advances the counter by one step. It returns an error if the counter
	// is already at its limit.
	Advance() error

	// Retreat retreats the counter by one step. It returns an error if the counter
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
}

// UpCounter represents a counter that increments upwards until it reaches an upper
// limit.
// It keeps track of the current count and the number of times it has been retreated.
type UpCounter struct {
	// upperLimit is the maximum limit that the counter can reach.
	upperLimit int

	// currentCount is the current value of the counter.
	currentCount int

	// retreatCount is the number of times the counter has been retreated.
	retreatCount int
}

// NewUpCounter creates a new UpCounter with the specified upper limit.
// It returns an error if the upper limit is less than 0.
//
// The function takes the following parameter:
//
//   - upperLimit is the maximum limit that the counter can reach.
//
// The function returns the following:
//
//   - A pointer to the newly created UpCounter.
//   - An error if the upper limit is less than 0.
func NewUpCounter(upperLimit int) (*UpCounter, error) {
	if upperLimit < 0 {
		return nil, new(ErrInvalidParameter).
			Parameter("upperLimit").
			Reason(fmt.Errorf("value (%d) must be positive", upperLimit))
	}

	return &UpCounter{upperLimit, 0, 0}, nil
}

// IsDone checks if the UpCounter has reached its upper limit.
// It returns true if the current count is greater than or equal to the
// upper limit, false otherwise.
func (c *UpCounter) IsDone() bool {
	return c.currentCount >= c.upperLimit
}

// Advance increments the current count of the UpCounter by one.
// It returns an error if the current count is already at or beyond the
// upper limit.
//
// The function returns the following:
//
//   - An error if the current count is already at or beyond the upper limit.
//   - nil if the current count was successfully incremented.
func (c *UpCounter) Advance() error {
	if c.currentCount >= c.upperLimit {
		return new(ErrCannotAdvanceCounter).
			Counter(c).
			Reason(errors.New("limit has already been reached"))
	}

	c.currentCount++
	return nil
}

// Retreat decrements the current count of the UpCounter by one.
// It also increments the retreat count.
// It returns an error if the current count is already at or below zero.
//
// The function returns the following:
//
//   - An error if the current count is already at or below zero.
//   - nil if the current count was successfully decremented.
func (c *UpCounter) Retreat() error {
	if c.currentCount <= 0 {
		return new(ErrCannotRetreatCounter).
			Counter(c).
			Reason(errors.New("limit cannot be lower than the current count"))
	}

	c.currentCount--
	c.retreatCount++
	return nil
}

// GetRetreatCount returns the number of times the UpCounter has
// been retreated.
// It returns an integer representing the retreat count.
func (c *UpCounter) GetRetreatCount() int {
	return c.retreatCount
}

// GetDistance calculates the distance between the current count and
// the upper limit of the UpCounter.
// It returns an integer representing the distance.
func (c *UpCounter) GetDistance() int {
	return c.upperLimit - c.currentCount
}

// GetCurrentCount returns the current count of the UpCounter.
// It returns an integer representing the current count.
func (c *UpCounter) GetCurrentCount() int {
	return c.currentCount
}

// GetInitialCount returns the upper limit of the UpCounter,
// which is the initial count.
// It returns an integer representing the initial count.
func (c *UpCounter) GetInitialCount() int {
	return c.upperLimit
}

// String returns a string representation of the UpCounter.
// The string includes the upper limit, current count, retreat
// count, and whether the counter is done.
//
// The function returns the following:
//
//   - A string representing the UpCounter.
func (c *UpCounter) String() string {
	var builder strings.Builder

	builder.WriteString("UpCounter{")
	builder.WriteString(fmt.Sprintf("upperLimit: %d, ", c.upperLimit))
	builder.WriteString(fmt.Sprintf("currentCount: %d, ", c.currentCount))
	builder.WriteString(fmt.Sprintf("retreatCount: %d, ", c.retreatCount))

	if c.IsDone() {
		builder.WriteString("isDone: true")
	} else {
		builder.WriteString("isDone: false")
	}

	builder.WriteRune('}')

	return builder.String()
}

// DownCounter represents a counter that decrements downwards until it reaches zero.
// It keeps track of the starting count, the current count, and the number of times
// it has been retreated.
type DownCounter struct {
	// startingCount is the initial value of the counter.
	startingCount int

	// currentCount is the current value of the counter.
	currentCount int

	// retreatCount is the number of times the counter has been retreated.
	retreatCount int
}

// NewDownCounter creates a new DownCounter with the specified starting count.
// It returns an error if the starting count is less than 0.
//
// The function takes the following parameter:
//
//   - startingCount is the initial value of the counter.
//
// The function returns the following:
//
//   - A pointer to the newly created DownCounter.
//   - An error if the starting count is less than 0.
func NewDownCounter(startingCount int) (*DownCounter, error) {
	if startingCount < 0 {
		return nil, new(ErrInvalidParameter).
			Parameter("startingCount").
			Reason(fmt.Errorf("value (%d) must be positive", startingCount))
	}

	return &DownCounter{startingCount, startingCount, 0}, nil
}

// IsDone checks if the DownCounter has reached zero.
// It returns true if the current count is less than or equal to zero,
// false otherwise.
func (c *DownCounter) IsDone() bool {
	return c.currentCount <= 0
}

// Advance decrements the current count of the DownCounter by one.
// It returns an error if the current count is already at or below zero.
//
// The function returns the following:
//
//   - An error if the current count is already at or below zero.
//   - nil if the current count was successfully decremented.
func (c *DownCounter) Advance() error {
	if c.currentCount <= 0 {
		return new(ErrCannotAdvanceCounter).
			Counter(c).
			Reason(errors.New("current count is already at zero"))
	}

	c.currentCount--
	return nil
}

// Retreat increments the starting count and current count of the
// DownCounter by one, and increments the retreat count.
// It returns an error if the starting count or current count is
// already at or below zero.
//
// The function returns the following:
//
//   - An error if the starting count or current count is already at
//     or below zero.
//   - nil if the starting count, current count, and retreat count were
//     successfully incremented.
func (c *DownCounter) Retreat() error {
	if c.startingCount <= 0 {
		return new(ErrCannotRetreatCounter).
			Counter(c).
			Reason(errors.New("starting count cannot be lower than zero"))
	}

	if c.currentCount <= 0 {
		return new(ErrCannotRetreatCounter).
			Counter(c).
			Reason(errors.New("current count is already at zero"))
	}

	c.startingCount++
	c.currentCount++
	c.retreatCount++

	return nil
}

// GetRetreatCount returns the number of times the DownCounter
// has been retreated.
// It returns an integer representing the retreat count.
func (c *DownCounter) GetRetreatCount() int {
	return c.retreatCount
}

// GetDistance returns the current count of the DownCounter.
// This is equivalent to the distance from zero, as the DownCounter
// decrements towards zero.
// It returns an integer representing the current count.
func (c *DownCounter) GetDistance() int {
	return c.currentCount
}

// GetCurrentCount returns the current count of the DownCounter.
// It returns an integer representing the current count.
func (c *DownCounter) GetCurrentCount() int {
	return c.currentCount
}

// GetInitialCount returns the starting count of the DownCounter,
// which is the initial count.
// It returns an integer representing the initial count.
func (c *DownCounter) GetInitialCount() int {
	return c.startingCount
}

// String returns a string representation of the DownCounter.
// The string includes the starting count, current count, retreat
// count, and whether the counter is done.
//
// The function returns the following:
//
//   - A string representing the DownCounter.
func (c *DownCounter) String() string {
	var builder strings.Builder

	builder.WriteString("DownCounter{")
	builder.WriteString(fmt.Sprintf("startingCount: %d, ", c.startingCount))
	builder.WriteString(fmt.Sprintf("currentCount: %d, ", c.currentCount))
	builder.WriteString(fmt.Sprintf("retreatCount: %d, ", c.retreatCount))

	if c.IsDone() {
		builder.WriteString("isDone: true")
	} else {
		builder.WriteString("isDone: false")
	}

	builder.WriteRune('}')

	return builder.String()
}
