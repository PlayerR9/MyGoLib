package FScreen

import "fmt"

type ScreenCounter struct {
	Label         string
	ElementsDone  int
	TotalElements int
}

func MakeCounter(label string, totalElements int) ScreenCounter {
	return ScreenCounter{
		Label:         label,
		ElementsDone:  0,
		TotalElements: totalElements,
	}
}

func (sc ScreenCounter) Equal(other ScreenCounter) bool {
	return sc.Label == other.Label
}

func (sc ScreenCounter) IsDone() bool {
	return sc.ElementsDone >= sc.TotalElements
}

func (sc ScreenCounter) HasLabel(label string) bool {
	return sc.Label == label
}

func (sc ScreenCounter) String() string {
	return fmt.Sprintf("%s: %d of %d", sc.Label, sc.ElementsDone, sc.TotalElements)
}

func (sc *ScreenCounter) Increment() {
	sc.ElementsDone++
}

func (sc *ScreenCounter) Reduce() {
	sc.TotalElements--
}
