package pkg

import (
	uts "github.com/PlayerR9/MyGoLib/Utility/Sorting"
)

var (
	ResultBranchSortFunc    uts.SortFunc[*resultBranch]
	FlagParseResultSortFunc uts.SortFunc[*FlagParseResult]
)

func init() {
	ResultBranchSortFunc = func(a, b *resultBranch) int {
		diff := len(b.resultMap) - len(a.resultMap)

		if diff != 0 {
			return diff
		}

		var size1, size2 int

		for _, v := range a.resultMap {
			size1 += v.size()
		}

		for _, v := range b.resultMap {
			size2 += v.size()
		}

		return size2 - size1
	}

	FlagParseResultSortFunc = func(a, b *FlagParseResult) int {
		diff := len(b.argMap) - len(a.argMap)

		if diff != 0 {
			return diff
		}

		var size1, size2 int

		for _, v := range a.argMap {
			size1 += len(v)
		}

		for _, v := range b.argMap {
			size2 += len(v)
		}

		return size2 - size1
	}
}
