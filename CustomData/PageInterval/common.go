package PageInterval

import (
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
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

type settingsTable struct {
	ws  string
	sep string
}

func WithWS(ws string) ffs.Option {
	return func(s ffs.Settinger) {
		set, ok := s.(*settingsTable)
		if !ok {
			return
		}

		set.ws = ws
	}
}

func WithSep(sep string) ffs.Option {
	return func(s ffs.Settinger) {
		set, ok := s.(*settingsTable)
		if !ok {
			return
		}

		if sep == "" {
			set.sep = ":"
		} else {
			set.sep = sep
		}
	}
}
