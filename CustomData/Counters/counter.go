package Counters

import (
	"fmt"
)

type UpCounter struct {
	label         string
	elementsDone  int
	totalElements int
}

func NewUpCounter(label string, totalElements int) UpCounter {
	if label == "" {
		panic("label must not be empty")
	}

	if totalElements < 0 {
		panic(fmt.Sprintf("totalElements must be non-negative; got %d", totalElements))
	}

	return UpCounter{
		label:         label,
		elementsDone:  0,
		totalElements: totalElements,
	}
}

func (sc *UpCounter) Equal(other UpCounter) bool {
	return sc.label == other.label
}

func (sc *UpCounter) IsDone() bool {
	return sc.elementsDone >= sc.totalElements
}

func (sc *UpCounter) ContainsLabel(label string) bool {
	return sc.label == label
}

func (sc *UpCounter) String() string {
	return fmt.Sprintf("%s: %d of %d", sc.label, sc.elementsDone, sc.totalElements)
}

func (sc *UpCounter) Increment() {
	sc.elementsDone++
}

func (sc *UpCounter) Reduce() {
	sc.totalElements--
}
