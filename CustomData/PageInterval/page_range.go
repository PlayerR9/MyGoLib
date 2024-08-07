package PageInterval

import (
	"slices"
	"strconv"
	"strings"

	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	uc "github.com/PlayerR9/lib_units/common"
	lup "github.com/PlayerR9/lib_units/pair"
)

// PageRange represents a pair of integers that represent the start and end
// page numbers of an interval.
// The first integer is the start page number and the second integer is the
// end page number of the interval. (both inclusive)
//
// For instance, the PageRange [1, 5] represents the interval from page 1 to
// page 5.
type PageRange lup.Pair[int, int]

// FString returns the string representation of the PageRange using the given
// traversor and options.
//
// Parameters:
//   - trav: The traversor to use for printing.
//   - ws: The whitespace to use between the elements. By default, it is a single space.
//   - sep: The separator to use between the key and value. By default, it is a colon.
//
// Behaviors:
//   - If sep is an empty string, it is set to a colon.
//   - ws can be empty. The default value is a single space.
//   - The default call for AString is: AString(trav, " ", "").
//   - If trav is empty, the function does nothing.
func (pr *PageRange) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	if trav == nil {
		return nil
	}

	settings := &settingsTable{
		ws:  " ",
		sep: ":",
	}

	for _, opt := range opts {
		opt(settings)
	}

	var err error

	if pr.First == pr.Second {
		err = trav.AppendString(strconv.Itoa(pr.First))
	} else {
		err = trav.AppendJoinedString(settings.ws, strconv.Itoa(pr.First), settings.sep, strconv.Itoa(pr.Second))
	}
	if err != nil {
		return err
	}

	trav.AcceptWord()

	return nil
}

// String returns the string representation of the PageRange.
//
// Returns:
//   - string: The string representation of the PageRange.
func (pr *PageRange) String() string {
	if pr.First == pr.Second {
		return strconv.Itoa(pr.First)
	} else {
		var builder strings.Builder

		builder.WriteString(strconv.Itoa(pr.First))
		builder.WriteRune(':')
		builder.WriteString(strconv.Itoa(pr.Second))

		return builder.String()
	}
}

// Iterator returns an iterator that iterates over the pages in the interval.
//
// Returns:
//   - uc.Iterater[int]: The iterator that iterates over the pages in the interval.
func (pr *PageRange) Iterator() uc.Iterater[int] {
	var builder uc.Builder[int]

	for page := pr.First; page <= pr.Second; page++ {
		builder.Add(page)
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
