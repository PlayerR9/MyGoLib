package Counters

import (
	"errors"
	"fmt"
)

type Counter interface {
	IsDone() bool
	Advance() error
	Retreat() error

	GetRetreatCount() int
	GetDistance() int
	GetCurrentCount() int
	GetInitialCount() int
}

type UpCounter struct {
	upperLimit, currentCount int
	retreatCount             int
}

func NewUpCounter(upperLimit int) (*UpCounter, error) {
	if upperLimit < 0 {
		return nil, new(ErrInvalidParameter).
			Parameter("upperLimit").
			Reason(fmt.Errorf("value (%d) must be positive", upperLimit))
	}

	return &UpCounter{upperLimit, 0, 0}, nil
}

func (c *UpCounter) IsDone() bool {
	return c.currentCount >= c.upperLimit
}

func (c *UpCounter) Advance() error {
	if c.currentCount >= c.upperLimit {
		return new(ErrCannotAdvanceCounter).
			Counter(c).
			Reason(errors.New("limit has already been reached"))
	}

	c.currentCount++
	return nil
}

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

func (c *UpCounter) GetRetreatCount() int {
	return c.retreatCount
}

func (c *UpCounter) GetDistance() int {
	return c.upperLimit - c.currentCount
}

func (c *UpCounter) GetCurrentCount() int {
	return c.currentCount
}

func (c *UpCounter) GetInitialCount() int {
	return c.upperLimit
}

type DownCounter struct {
	startingCount, currentCount int
	retreatCount                int
}

func NewDownCounter(startingCount int) (*DownCounter, error) {
	if startingCount < 0 {
		return nil, new(ErrInvalidParameter).
			Parameter("startingCount").
			Reason(fmt.Errorf("value (%d) must be positive", startingCount))
	}

	return &DownCounter{startingCount, startingCount, 0}, nil
}

func (c *DownCounter) IsDone() bool {
	return c.currentCount <= 0
}

func (c *DownCounter) Advance() error {
	if c.currentCount <= 0 {
		return new(ErrCannotAdvanceCounter).
			Counter(c).
			Reason(errors.New("current count is already at zero"))
	}

	c.currentCount--
	return nil
}

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

	c.startingCount--
	c.currentCount--
	c.retreatCount++

	return nil
}

func (c *DownCounter) GetRetreatCount() int {
	return c.retreatCount
}

func (c *DownCounter) GetDistance() int {
	return c.currentCount
}

func (c *DownCounter) GetCurrentCount() int {
	return c.currentCount
}

func (c *DownCounter) GetInitialCount() int {
	return c.startingCount
}
