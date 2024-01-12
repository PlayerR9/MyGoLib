package PageInterval

import (
	"fmt"
	"math"
	"slices"
	"sort"
	"strings"
)

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
	intervals [][2]int

	// The 'pageCount' field represents the total number of pages across all
	// intervals.
	pageCount int
}

// String returns a string representation of the PageInterval.
// It concatenates all intervals in the PageInterval and returns them as a
// formatted string.
// Each interval is represented as "[start : end]".
//
// Example:
//
//	pi := PageInterval{
//	    intervals: [][2]int{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	fmt.Println(pi.String()) // Output: "[1 : 5][10 : 15]"
//
// Note: This method does not include the pageCount in the string representation.
func (pi *PageInterval) String() string {
	var builder strings.Builder

	for _, interval := range pi.intervals {
		builder.WriteString(fmt.Sprintf("[%d : %d]", interval[0], interval[1]))
	}

	return builder.String()
}

// ToSlice converts the PageInterval into a slice of integers.
// It iterates through the intervals in the PageInterval and adds each page
// number to the slice.
// The resulting slice is returned.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: [][2]int{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	slice := pi.ToSlice()
//	fmt.Println(slice) // Output: [1 2 3 4 5 10 11 12 13 14 15]
func (pi *PageInterval) ToSlice() []int {
	slice := make([]int, 0, pi.pageCount)

	for _, interval := range pi.intervals {
		if interval[0] == interval[1] {
			slice = append(slice, interval[0])
		} else {
			for j := interval[0]; j <= interval[1]; j++ {
				slice = append(slice, j)
			}
		}
	}

	return slice
}

// PageCount returns the total number of pages across all intervals in the
// PageInterval.
// The pageCount is precalculated and stored in the PageInterval struct.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: [][2]int{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	count := pi.PageCount()
//	fmt.Println(count) // Output: 11
func (pi *PageInterval) PageCount() int {
	return pi.pageCount
}

// Intervals returns the intervals stored in the PageInterval.
// Each interval is represented as a pair of integers, where the first
// integer is the start page number and the second integer is the end page
// number.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: [][2]int{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	intervals := pi.Intervals()
//	fmt.Println(intervals) // Output: [[1 5] [10 15]]
func (pi *PageInterval) Intervals() [][2]int {
	return pi.intervals
}

// NewPageInterval creates a new instance of PageInterval.
// It initializes the 'intervals' field with an empty slice and sets the
// 'pageCount' field to 0.
//
// Example:
//
//	pi := NewPageInterval()
//	fmt.Println(pi.intervals) // Output: []
//	fmt.Println(pi.pageCount) // Output: 0
func NewPageInterval() PageInterval {
	return PageInterval{
		intervals: make([][2]int, 0),
		pageCount: 0,
	}
}

// HasPages checks if the PageInterval has any pages.
// It returns true if the page count is greater than 0, otherwise false.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: [][2]int{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	hasPages := pi.HasPages()
//	fmt.Println(hasPages) // Output: true
func (pi *PageInterval) HasPages() bool {
	return pi.pageCount > 0
}

// GetFirstPage returns the first page of the PageInterval.
// It returns the page number and an error if no pages have been set.
// If pages have been set, it returns the start page number of the first
// interval and nil error.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: [][2]int{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	firstPage, err := pi.GetFirstPage()
//	if err != nil {
//	    fmt.Println(err) // Output: no pages have been set
//	} else {
//	    fmt.Println(firstPage) // Output: 1
//	}
func (pi *PageInterval) GetFirstPage() (int, error) {
	if len(pi.intervals) == 0 {
		return 0, &ErrNoPagesHaveBeenSet{}
	}

	return pi.intervals[0][0], nil
}

// GetLastPage returns the last page of the PageInterval.
// It returns the page number and an error if no pages have been set.
// If pages have been set, it returns the end page number of the last
// interval and nil error.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: [][2]int{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	lastPage, err := pi.GetLastPage()
//	if err != nil {
//	    fmt.Println(err) // Output: no pages have been set
//	} else {
//	    fmt.Println(lastPage) // Output: 15
//	}
func (pi *PageInterval) GetLastPage() (int, error) {
	if len(pi.intervals) == 0 {
		return 0, &ErrNoPagesHaveBeenSet{}
	}

	return pi.intervals[len(pi.intervals)-1][1], nil
}

// MustGetFirstPage returns the first page number in the PageInterval.
// It panics with ErrNoPagesHaveBeenSet if no pages have been set.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: [][2]int{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	firstPage := pi.MustGetFirstPage()
//	fmt.Println(firstPage) // Output: 1
func (pi *PageInterval) MustGetFirstPage() int {
	if len(pi.intervals) == 0 {
		panic(&ErrNoPagesHaveBeenSet{})
	}

	return pi.intervals[0][0]
}

// MustGetLastPage returns the last page number in the PageInterval.
// It panics with ErrNoPagesHaveBeenSet if no pages have been set.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: [][2]int{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	lastPage := pi.MustGetLastPage()
//	fmt.Println(lastPage) // Output: 15
func (pi *PageInterval) MustGetLastPage() int {
	if len(pi.intervals) == 0 {
		panic(&ErrNoPagesHaveBeenSet{})
	}

	return pi.intervals[len(pi.intervals)-1][1]
}

// AddPage adds a page to the PageInterval.
// If the page number is less than 1, it is considered a no-op and the
// function returns.
// If the PageInterval is empty, a new interval is created with the given
// page number.
// If the page number falls within an existing interval, the interval is
// updated accordingly.
// If the page number does not fall within any existing interval, a new
// interval is created.
// After adding the page, the PageInterval is reduced to merge overlapping
// intervals.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: [][2]int{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	pi.AddPage(6)
//	fmt.Println(pi.intervals) // Output: [[1 6] [10 15]]
//	fmt.Println(pi.pageCount) // Output: 12
func (pi *PageInterval) AddPage(page int) {
	if page < 1 {
		return // No-op
	}

	if len(pi.intervals) == 0 {
		pi.intervals = append(pi.intervals, [2]int{page, page})
		pi.pageCount++
		return
	}

	insertPos := sort.Search(len(pi.intervals), func(i int) bool {
		return pi.intervals[i][0] >= page
	})

	if insertPos > 0 && pi.intervals[insertPos-1][1] >= page-1 {
		pi.intervals[insertPos-1][1] = int(math.Max(float64(pi.intervals[insertPos-1][1]), float64(page)))
		pi.pageCount++
		reduce(pi)
		return
	}

	if insertPos < len(pi.intervals) && pi.intervals[insertPos][0] <= page+1 {
		pi.intervals[insertPos][0] = int(math.Min(float64(pi.intervals[insertPos][0]), float64(page)))
		pi.pageCount++
		reduce(pi)
		return
	}

	pi.intervals = append(pi.intervals[:insertPos], append([][2]int{{page, page}}, pi.intervals[insertPos:]...)...)
	pi.pageCount++

	reduce(pi)
}

// RemovePage removes the specified page from the PageInterval.
// If the page number is less than 1, it is considered a no-op and
// no changes are made.
// If the page number is not found in the PageInterval, no changes
// are made.
// After removing the page, the page count is decremented and the intervals
// are updated accordingly.
// Finally, the PageInterval is reduced if necessary.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: [][2]int{{1, 5}, {10, 15}},
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

	updatePages(pi, index, pi.intervals[index], page)
	pi.pageCount--

	reduce(pi)
}

// HasPage checks if the given page exists in the PageInterval.
// It returns true if the page is found, otherwise false.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: [][2]int{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	hasPage := pi.HasPage(3)
//	fmt.Println(hasPage) // Output: true
func (pi *PageInterval) HasPage(page int) bool {
	return findPageInterval(pi, page) != -1
}

// AddPagesBetween adds pages between the first and last page numbers to
// the PageInterval.
// If the first page number is less than 1, it is set to 1 to remove invalid
// pages.
// If the last page number is less than 1, it is set to 1 to remove invalid
// pages.
// If the last page number is less than the first page number, the values
// are swapped.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: [][2]int{{1, 5}, {10, 15}},
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
//	    intervals: [][2]int{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	pi.RemovePagesBetween(3, 4)
//	fmt.Println(pi.intervals) // Output: [[1 2] [5 5] [10 15]]
//	fmt.Println(pi.pageCount) // Output: 9
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

// Iterator returns a PageIntervalIterator for iterating over the
// intervals in the PageInterval.
// The iterator starts from the first interval and moves forward.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: [][2]int{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	iterator := pi.Iterator()
//	for iterator.HasNext() {
//	    fmt.Println(iterator.Next()) // Output: 1, 2, 3, 4, 5, 10, 11, 12, 13, 14, 15
//	}
func (pi *PageInterval) Iterator() PageIntervalIterator {
	return PageIntervalIterator{
		intervals: pi.intervals,
		index:     0,
		page:      -1,
	}
}

// ReverseIterator returns a reverse iterator for the PageInterval.
// The reverse iterator allows iterating over the intervals in
// reverse order.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: [][2]int{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	reverseIterator := pi.ReverseIterator()
//	for reverseIterator.HasPrev() {
//	    fmt.Println(reverseIterator.Prev()) // Output: 15, 14, 13, 12, 11, 10, 5, 4, 3, 2, 1
//	}
func (pi *PageInterval) ReverseIterator() PageIntervalReverseIterator {
	return PageIntervalReverseIterator{
		intervals: pi.intervals,
		index:     len(pi.intervals) - 1,
		page:      -1,
	}
}

// reduce merges overlapping intervals in the PageInterval.
// It sorts the intervals based on the start value and then merges any
// overlapping intervals.
// The merged intervals are stored in the intervals field of the PageInterval.
// If the PageInterval contains less than two intervals, no operation is
// performed.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: [][2]int{{1, 5}, {10, 15}, {4, 7}},
//	    pageCount: 13,
//	}
//
//	reduce(&pi)
//	fmt.Println(pi.intervals) // Output: [[1 7] [10 15]]
func reduce(pi *PageInterval) {
	if len(pi.intervals) < 2 {
		return
	}

	sort.Slice(pi.intervals, func(i, j int) bool {
		return pi.intervals[i][0] < pi.intervals[j][0]
	})

	mergedIntervals := make([][2]int, 0, len(pi.intervals))
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

// findPageInterval searches for the interval that contains the given
// page number in the PageInterval.
// It returns the index of the interval in the intervals slice if found,
// otherwise -1.
// If the page number is less than 1 or the PageInterval is empty (pageCount
// is 0), it returns -1.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: [][2]int{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	index := findPageInterval(&pi, 3)
//	fmt.Println(index) // Output: 0
func findPageInterval(pi *PageInterval, page int) int {
	if page < 1 || pi.pageCount == 0 {
		return -1
	}

	return slices.IndexFunc(pi.intervals, func(interval [2]int) bool {
		return interval[0] <= page && page <= interval[1]
	})
}

// updatePages updates the PageInterval by modifying the intervals based on
// the given parameters.
// If p[0] is equal to p[1], the interval at the specified index is removed
// from the intervals slice.
// If p[0] is equal to page, the starting page of the interval at the specified
// index is incremented by 1.
// If p[1] is equal to page, the ending page of the interval at the specified
// index is decremented by 1.
// Otherwise, a new interval is created by splitting the existing interval
// at the specified index.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: [][2]int{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	updatePages(&pi, 0, [2]int{1, 5}, 3)
//	fmt.Println(pi.intervals) // Output: [[1 2] [4 5] [10 15]]
func updatePages(pi *PageInterval, index int, p [2]int, page int) {
	if p[0] == p[1] {
		pi.intervals = append(pi.intervals[:index], pi.intervals[index+1:]...)
	} else if p[0] == page {
		pi.intervals[index][0]++
	} else if p[1] == page {
		pi.intervals[index][1]--
	} else {
		splitPageInterval(pi, index, p, page)
	}
}

// splitPageInterval splits the given page interval at the specified index
// and inserts a new interval.
// It modifies the intervals slice of the PageInterval struct.
// The pi parameter is a pointer to the PageInterval struct.
// The index parameter specifies the index at which the interval should be
// split.
// The p parameter is an array representing the original interval to be split.
// The page parameter is the page number at which the split should occur.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: [][2]int{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	splitPageInterval(&pi, 0, [2]int{1, 5}, 3)
//	fmt.Println(pi.intervals) // Output: [[1 2] [4 5] [10 15]]
func splitPageInterval(pi *PageInterval, index int, p [2]int, page int) {
	pi.intervals[index] = [2]int{p[0], page - 1}
	pi.intervals = append(append(pi.intervals[:index+1], [2]int{page + 1, p[1]}), pi.intervals[index+1:]...)
}
