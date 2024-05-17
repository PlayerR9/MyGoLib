package PageInterval

import (
	"fmt"
	"slices"

	itf "github.com/PlayerR9/MyGoLib/ListLike/Iterator"

	cdp "github.com/PlayerR9/MyGoLib/Units/Pair"
)

// PageRange represents a pair of integers that represent the start and end
// page numbers of an interval.
// The first integer is the start page number and the second integer is the
// end page number of the interval. (both inclusive)
//
// For instance, the PageRange [1, 5] represents the interval from page 1 to
// page 5.
type PageRange cdp.Pair[int, int]

// String returns the string representation of the PageRange.
//
// Returns:
//   - string: The string representation of the PageRange.
func (pr *PageRange) String() string {
	return fmt.Sprintf("[%d : %d]", pr.First, pr.Second)
}

// Iterator returns an iterator that iterates over the pages in the interval.
//
// Returns:
//   - itf.Iterater[int]: The iterator that iterates over the pages in the interval.
func (pr *PageRange) Iterator() itf.Iterater[int] {
	var builder itf.Builder[int]

	for page := pr.First; page <= pr.Second; page++ {
		builder.Append(page)
	}

	return builder.Build()
}

// newPageRange creates a new instance of PageRange with the given start and
// end page numbers.
//
// Parameters:
//
//   - start: The start page number of the interval.
//   - end: The end page number of the interval.
//
// Returns:
//
//   - *PageRange: The new PageRange.
func newPageRange(start, end int) *PageRange {
	return &PageRange{start, end}
}

// findPageInterval searches for the interval that contains the given page
// number in the PageInterval.
//
// Parameters:
//   - pi: A pointer to the PageInterval to search in.
//   - page: The page number to search for in the PageInterval.
//
// Returns:
//   - int: The index of the interval in the intervals slice if found, otherwise -1.
func (pi *PageInterval) findPageInterval(page int) int {
	if page < 1 || pi.pageCount == 0 {
		return -1
	}

	isPageBetween := func(interval *PageRange) bool {
		return interval.First <= page && page <= interval.Second
	}

	return slices.IndexFunc(pi.intervals, isPageBetween)
}
