package PageInterval

// PageIntervalIterator is a struct that allows iterating over the
// intervals in a PageInterval.
// It contains a slice of intervals, and an index and page number
// to keep track of the current position.
type PageIntervalIterator struct {
	intervals [][2]int
	index     int
	page      int
}

// Next advances the iterator to the next page in the PageInterval.
// It returns true if there is a next page, otherwise false.
// On the first call, it sets the page to the first page of the first
// interval.
// On subsequent calls, it increments the page number and checks if it
// is within the current interval.
// If the page number exceeds the current interval, it moves to the next
// interval.
func (pii *PageIntervalIterator) Next() bool {
	// First call
	if pii.page == -1 {
		if len(pii.intervals) == 0 {
			return false
		}

		pii.page = pii.intervals[0][0]

		return true
	}

	// Subsequent calls
	pii.page++

	if pii.page <= pii.intervals[pii.index][1] {
		return true
	}

	hasNext := pii.index+1 < len(pii.intervals)

	if hasNext {
		pii.index++
		pii.page = pii.intervals[pii.index][0]
	}

	return hasNext
}

// Value returns the current page number in the iterator.
// It should be called after Next to get the current page number.
func (pii *PageIntervalIterator) Value() int {
	return pii.page
}

// PageIntervalReverseIterator is a struct that allows iterating over the intervals
// in a PageInterval in reverse order.
// It contains a slice of intervals, and an index and page number to keep track
// of the current position.
type PageIntervalReverseIterator struct {
	intervals [][2]int
	index     int
	page      int
}

// Previous moves the iterator to the previous page in the PageInterval.
// It returns true if there is a previous page, otherwise false.
// On the first call, it sets the page to the last page of the last interval.
// On subsequent calls, it decrements the page number and checks if it is
// within the current interval.
// If the page number is less than the current interval, it moves to the
// previous interval.
func (piri *PageIntervalReverseIterator) Previous() bool {
	// First call
	if piri.page == -1 {
		if len(piri.intervals) == 0 {
			return false
		}

		piri.index = len(piri.intervals) - 1
		piri.page = piri.intervals[piri.index][1]

		return true
	}

	// Subsequent calls
	piri.page--

	if piri.page >= piri.intervals[piri.index][0] {
		return true
	}

	hasPrevious := piri.index > 0

	if hasPrevious {
		piri.index--
		piri.page = piri.intervals[piri.index][1]
	}

	return hasPrevious
}

// Value returns the current page number in the iterator.
// It should be called after Previous to get the current page number.
func (piri *PageIntervalReverseIterator) Value() int {
	return piri.page
}
