package FScreen

import (
	"fmt"
	"sync"
)

type ScreenCounter struct {
	label         string
	elementsDone  int
	totalElements int
	mutex         sync.RWMutex
}

// WARNING: It assumes everything is right
func NewScreenCounter(label string, totalElements int) ScreenCounter {
	return ScreenCounter{
		label:         label,
		elementsDone:  0,
		totalElements: totalElements,
	}
}

func (sc *ScreenCounter) Equal(other *ScreenCounter) bool {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()

	return other.ContainsLabel(sc.label)
}

func (sc *ScreenCounter) IsDone() bool {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()

	return sc.elementsDone >= sc.totalElements
}

func (sc *ScreenCounter) ContainsLabel(label string) bool {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()

	return sc.label == label
}

func (sc *ScreenCounter) String() string {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()

	return fmt.Sprintf("%s: %d of %d", sc.label, sc.elementsDone, sc.totalElements)
}

func (sc *ScreenCounter) Increment() {
	sc.mutex.Lock()
	sc.elementsDone++
	sc.mutex.Unlock()
}

func (sc *ScreenCounter) Reduce() {
	sc.mutex.Lock()
	sc.totalElements--
	sc.mutex.Unlock()
}
