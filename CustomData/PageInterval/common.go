package PageInterval

import (
	itf "github.com/PlayerR9/MyGoLib/Units/Iterator"
)

var (
	// PageRangeIterator returns an iterator that iterates over the pages in the interval.
	//
	// Parameters:
	//   - pr: The page range.
	//
	// Returns:
	//   - itf.Iterater[int]: The iterator that iterates over the pages in the interval.
	PageRangeIterator func(pr *PageRange) itf.Iterater[int] = func(pr *PageRange) itf.Iterater[int] {
		return pr.Iterator()
	}
)
