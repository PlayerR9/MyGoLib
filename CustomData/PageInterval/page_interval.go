// Package PageInterval provides a data structure for managing page intervals.
package PageInterval

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	itf "github.com/PlayerR9/MyGoLib/CustomData/Iterators"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
)

// PageRange represents a pair of integers that represent the start and end
// page numbers of an interval.
// The first integer is the start page number and the second integer is the
// end page number of the interval. (both inclusive)
//
// For instance, the PageRange [1, 5] represents the interval from page 1 to
// page 5.
type PageRange [2]int

// PageInterval represents a collection of page intervals, where each
// interval is represented by a pair of integers.
// The PageInterval ensures non-overlapping, non-duplicate intervals and
// reduces the amount of intervals by merging two consecutive intervals
// into one.
//
// Example:
//
// pi := NewPageInterval()
// pi.AddPagesBetween(1, 5)
// pi.AddPagesBetween(10, 15)
//
// fmt.Println(pi.Intervals()) // Output: [[1 5] [10 15]]
// fmt.Println(pi.PageCount()) // Output: 11
type PageInterval struct {
	// The 'intervals' field is a slice of integer pairs, where each pair
	// represents a start and end page number of an interval.
	intervals []PageRange

	// The 'pageCount' field represents the total number of pages across all
	// intervals.
	pageCount int
}

// String is a method of the PageInterval type that returns a string
// representation of the PageInterval.
// Each interval is represented as "start : end" separated by a comma.
//
// Returns:
//
//   - string: A formatted string representation of the PageInterval.
func (pi *PageInterval) String() string {
	if pi.pageCount == 0 {
		return "PageInterval[]"
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "PageInterval[%d : %d", pi.intervals[0][0], pi.intervals[0][1])

	for _, interval := range pi.intervals[1:] {
		builder.WriteRune(',')
		builder.WriteRune(' ')
		fmt.Fprintf(&builder, "%d : %d", interval[0], interval[1])
	}

	builder.WriteRune(']')

	return builder.String()
}

// Iterator is a method of the PageInterval type that returns an iterator for
// iterating over the pages in the PageInterval.
//
// Panics if an error occurs while creating the iterator.
//
// Returns:
//
//   - itf.Iterater[int]: An iterator for iterating over the pages in the PageInterval.
func (pi *PageInterval) Iterator() itf.Iterater[int] {
	iter, err := itf.IteratorFromIterator(
		itf.IteratorFromSlice(pi.intervals),
		func(pr PageRange) itf.Iterater[int] {
			var builder itf.Builder[int]

			for page := pr[0]; page <= pr[1]; page++ {
				builder.Append(page)
			}

			return builder.Build()
		},
	)
	if err != nil {
		panic(err)
	}

	return iter
}

// PageCount is a method of the PageInterval type that returns the total number
// of pages across all intervals in the PageInterval.
//
// Returns:
//
//   - pageCount: The total number of pages across all intervals in the PageInterval.
func (pi *PageInterval) PageCount() int {
	return pi.pageCount
}

// Intervals is a method of the PageInterval type that returns the intervals
// stored in the PageInterval.
// Each interval is represented as a pair of integers, where the first integer
// is the start page number and the second integer is the end page number.
//
// Returns:
//
//   - intervals: A slice of integer pairs representing the intervals in the
//     PageInterval.
func (pi *PageInterval) Intervals() []PageRange {
	return pi.intervals
}

// NewPageInterval creates a new instance of PageInterval with
// empty intervals and a page count of 0.
//
// Returns:
//
//   - PageInterval: The new PageInterval.
func NewPageInterval() PageInterval {
	return PageInterval{
		intervals: make([]PageRange, 0),
		pageCount: 0,
	}
}

// HasPages is a method of the PageInterval type that checks if the PageInterval
// has any pages.
//
// Returns:
//
//   - bool: A boolean value that is true if the PageInterval has pages, and
//     false otherwise.
func (pi *PageInterval) HasPages() bool {
	return pi.pageCount > 0
}

// GetFirstPage is a method of the PageInterval type that returns the first
// page number in the PageInterval.
// It panics with *ers.ErrCallFailed if no pages have been set.
//
// Returns:
//
//   - int: The first page number in the PageInterval.
//   - error: An error of type *ers.ErrNoPagesInInterval if no pages have been set.
func (pi *PageInterval) GetFirstPage() (int, error) {
	if pi.pageCount <= 0 {
		return 0, NewErrNoPagesInInterval()
	}

	return pi.intervals[0][0], nil
}

// GetLastPage is a method of the PageInterval type that returns the last
// page number in the PageInterval.
// It panics with *ers.ErrCallFailed if no pages have been set.
//
// Returns:
//
//   - int: The last page number in the PageInterval.
//   - error: An error of type *ers.ErrNoPagesInInterval if no pages have been set.
func (pi *PageInterval) GetLastPage() (int, error) {
	if pi.pageCount <= 0 {
		return 0, NewErrNoPagesInInterval()
	}

	return pi.intervals[len(pi.intervals)-1][1], nil
}

// AddPage is a method of the PageInterval type that adds a page to the
// PageInterval, maintaining the non-overlapping, non-duplicate intervals.
//
// Parameters:
//
//   - page: The page number to add to the PageInterval.
//
// Returns:
//
//   - error: An error of type *ers.ErrInvalidParameter if the page number is less than 1.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: []PageRange{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	pi.AddPage(6)
//	fmt.Println(pi.intervals) // Output: [[1 6] [10 15]]
//	fmt.Println(pi.pageCount) // Output: 12
func (pi *PageInterval) AddPage(page int) error {
	if page < 1 {
		return ers.NewErrInvalidParameter(
			"page",
			fmt.Errorf("page number (%d) must be greater than 0", page),
		)
	}

	if len(pi.intervals) == 0 {
		pi.intervals = append(pi.intervals, PageRange{page, page})
	} else {
		insertPos := sort.Search(len(pi.intervals), func(i int) bool {
			return pi.intervals[i][0] >= page
		})

		if insertPos > 0 && pi.intervals[insertPos-1][1] >= page-1 {
			insertPos--
			pi.intervals[insertPos][1] = gen.Max(pi.intervals[insertPos][1], page)
		} else if insertPos < len(pi.intervals) && pi.intervals[insertPos][0] <= page+1 {
			pi.intervals[insertPos][0] = gen.Min(pi.intervals[insertPos][0], page)
		} else {
			pi.intervals = append(pi.intervals[:insertPos],
				append([]PageRange{{page, page}}, pi.intervals[insertPos:]...)...,
			)
		}
	}

	pi.pageCount++
	reduce(pi)

	return nil
}

// RemovePage is a method of the PageInterval type that removes the specified
// page from the PageInterval.
// No changes are made if the page number is less than 1 or not found in the
// PageInterval.
//
// Parameters:
//
//   - page: The page number to remove from the PageInterval.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: []PageRange{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	pi.RemovePage(5)
//	fmt.Println(pi.intervals) // Output: [[1 4] [10 15]]
//	fmt.Println(pi.pageCount) // Output: 10
func (pi *PageInterval) RemovePage(page int) {
	if page < 1 {
		return // No-op
	}

	index := findPageInterval(pi, page)
	if index == -1 {
		return
	}

	if pi.intervals[index][0] == pi.intervals[index][1] {
		pi.intervals = append(pi.intervals[:index], pi.intervals[index+1:]...)
	} else if pi.intervals[index][0] == page {
		pi.intervals[index][0]++
	} else if pi.intervals[index][1] == page {
		pi.intervals[index][1]--
	} else {
		newIntervals := make([]PageRange, len(pi.intervals)+1)

		// Copy the intervals before the split
		copy(newIntervals, pi.intervals[:index+1])

		// Modify the interval at the split index
		newIntervals[index] = PageRange{pi.intervals[index][0], page - 1}

		// Add the new interval
		newIntervals[index+1] = PageRange{page + 1, pi.intervals[index][1]}

		// Copy the intervals after the split
		copy(newIntervals[index+2:], pi.intervals[index+1:])

		pi.intervals = newIntervals
	}

	pi.pageCount--

	reduce(pi)
}

// HasPage checks if the given page exists in the PageInterval.
// It returns true if the page is found, otherwise false.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: []PageRange{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	hasPage := pi.HasPage(3)
//	fmt.Println(hasPage) // Output: true

// HasPage is a method of the PageInterval type that checks if the given page
// exists in the PageInterval.
//
// Parameters:
//
//   - page: The page number to check for in the PageInterval.
//
// Returns:
//
//   - bool: A boolean value that is true if the page exists in the PageInterval,
//     and false otherwise.
func (pi *PageInterval) HasPage(page int) bool {
	return findPageInterval(pi, page) != -1
}

// AddPagesBetween is a method of the PageInterval type that adds pages between
// the first and last page numbers to the PageInterval.
//
// However, if the first page number is less than 1, it is set to 1 to remove
// invalid pages, same goes for the last page number.
// Finally, if the last page number is less than the first page number, the
// values are swapped.
//
// Parameters:
//
//   - first: The first page number to add to the PageInterval.
//   - last: The last page number to add to the PageInterval.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: []PageRange{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	pi.AddPagesBetween(6, 9)
//	fmt.Println(pi.intervals) // Output: [[1 15]]
//	fmt.Println(pi.pageCount) // Output: 15
func (pi *PageInterval) AddPagesBetween(first, last int) {
	if first < 1 {
		first = 1 // remove invalid pages
	}

	if last < 1 {
		last = 1 // remove invalid pages
	}

	if last < first {
		last, first = first, last // swap values
	}

	for i := first; i <= last; i++ {
		pi.AddPage(i)
	}
}

// RemovePagesBetween removes pages between the specified first and last
// page numbers from the PageInterval.
// If the first page number is less than 1, it is set to 1 to remove invalid
// pages.
// If the last page number is less than 1, it is set to 1 to remove invalid
// pages.
// If the last page number is less than the first page number, the values are
// swapped.
// Pages between the first and last page numbers (inclusive) are removed using
// the RemovePage method.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: []PageRange{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	pi.RemovePagesBetween(3, 4)
//	fmt.Println(pi.intervals) // Output: [[1 2] [5 5] [10 15]]
//	fmt.Println(pi.pageCount) // Output: 9

// RemovePagesBetween is a method of the PageInterval type that removes pages
// between the specified first and last page numbers from the PageInterval.
//
// However, if the first page number is less than 1, it is set to 1 to remove
// invalid pages, same goes for the last page number.
// Finally, if the last page number is less than the first page number, the
// values are swapped.
//
// Parameters:
//
//   - first, last: The first and last page numbers to remove from the PageInterval,
//     respectively.
func (pi *PageInterval) RemovePagesBetween(first, last int) {
	if first < 1 {
		first = 1 // remove invalid pages
	}

	if last < 1 {
		last = 1 // remove invalid pages
	}

	if last < first {
		last, first = first, last // swap values
	}

	for i := first; i <= last; i++ {
		pi.RemovePage(i)
	}
}

// ReverseIterator is a method of the PageInterval type that returns a
// PageIntervalReverseIterator for iterating over the intervals in the
// PageInterval in reverse order.
//
// Panics if an error occurs while creating the iterator.
//
// Returns:
//
//   - itf.Iterater[int]: An iterator for iterating over the intervals in the
//     PageInterval in reverse order.
func (pi *PageInterval) ReverseIterator() itf.Iterater[int] {
	reversed := make([]PageRange, len(pi.intervals))
	copy(reversed, pi.intervals)

	slices.Reverse(reversed)

	iter, err := itf.IteratorFromIterator(
		itf.IteratorFromSlice(reversed),
		func(pr PageRange) itf.Iterater[int] {
			var builder itf.Builder[int]

			for page := pr[1]; page >= pr[0]; page-- {
				builder.Append(page)
			}

			return builder.Build()
		},
	)
	if err != nil {
		panic(err)
	}

	return iter
}

// reduce merges overlapping intervals in the PageInterval.
// It sorts the intervals based on the start value and then merges any
// overlapping intervals.
// The merged intervals are stored in the intervals field of the PageInterval.
// If the PageInterval contains less than two intervals, no operation is
// performed.
//
// Parameters:
//
//   - pi: A pointer to the PageInterval to reduce.
func reduce(pi *PageInterval) {
	if len(pi.intervals) < 2 {
		return
	}

	sort.Slice(pi.intervals, func(i, j int) bool {
		return pi.intervals[i][0] < pi.intervals[j][0]
	})

	mergedIntervals := make([]PageRange, 0, len(pi.intervals))
	currentInterval := pi.intervals[0]

	for i := 1; i < len(pi.intervals); i++ {
		nextInterval := pi.intervals[i]
		if currentInterval[1] >= nextInterval[0]-1 {
			if nextInterval[1] > currentInterval[1] {
				currentInterval[1] = nextInterval[1]
			}
		} else {
			mergedIntervals = append(mergedIntervals, currentInterval)
			currentInterval = nextInterval
		}
	}

	mergedIntervals = append(mergedIntervals, currentInterval)
	pi.intervals = mergedIntervals
}

// findPageInterval searches for the interval that contains the given page
// number in the PageInterval.
//
// Parameters:
//
//   - pi: A pointer to the PageInterval to search in.
//   - page: The page number to search for in the PageInterval.
//
// Returns:
//
//   - int: The index of the interval in the intervals slice if found, otherwise -1.
func findPageInterval(pi *PageInterval, page int) int {
	if page < 1 || pi.pageCount == 0 {
		return -1
	}

	return slices.IndexFunc(pi.intervals, func(interval PageRange) bool {
		return interval[0] <= page && page <= interval[1]
	})
}
