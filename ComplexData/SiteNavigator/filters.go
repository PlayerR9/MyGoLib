package SiteNavigator

import (
	us "github.com/PlayerR9/MyGoLib/Units/Slice"
)

var (
	FilterNilFEFuncs us.PredicateFilter[FilterErrFunc]
)

func init() {
	FilterNilFEFuncs = func(fef FilterErrFunc) bool {
		return fef != nil
	}
}
