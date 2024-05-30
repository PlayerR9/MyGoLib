// Package CnsPanel provides a structure and functions for handling
// console command flags.
package pkg

import (
	"fmt"

	uc "github.com/PlayerR9/MyGoLib/Units/Common"
)

type resultBranch struct {
	resultMap map[string]any // *FlagParseResult or error
}

func (rb *resultBranch) Copy() uc.Copier {
	resultMap := make(map[string]any)
	for k, v := range rb.resultMap {
		resultMap[k] = v
	}

	return &resultBranch{resultMap}
}

func newResultBranch() *resultBranch {
	return &resultBranch{make(map[string]any)}
}

func (rb *resultBranch) size() int {
	return len(rb.resultMap)
}

func (rb *resultBranch) evaluate(flagName string, result map[string]any, err error) []*resultBranch {
	prev, ok := rb.resultMap[flagName]
	if !ok {
		// At the first evaluation, we have no previous result and so,
		// we can just store the result as is.
		if err != nil {
			rb.resultMap[flagName] = err
		} else {
			rb.resultMap[flagName] = result
		}

		return []*resultBranch{rb}
	}

	if err != nil {
		_, ok := prev.(error)
		if ok {
			// Prioritize the latest error.
			rb.resultMap[flagName] = err
		}

		return []*resultBranch{rb}
	} else {
		switch prev := prev.(type) {
		case error:
			// Prioritize the latest result over the error.
			rb.resultMap[flagName] = result

			return []*resultBranch{rb}
		case *FlagParseResult:
			// Possible conflict. Duplicate branch.

			rbCopy := rb.Copy().(*resultBranch)

			rbCopy.resultMap[flagName] = result

			return []*resultBranch{rb, rbCopy}
		default:
			panic(fmt.Sprintf("unexpected type %T", prev))
		}
	}
}

func (rb *resultBranch) errIfInvalidRequiredFlags(flags []*FlagInfo) error {
	for _, flag := range flags {
		if !flag.IsRequired() {
			continue
		}

		val, ok := rb.resultMap[flag.GetName()]
		if !ok {
			return fmt.Errorf("missing required flag %q", flag.GetName())
		}

		switch val := val.(type) {
		case *FlagParseResult:
			// Do nothing.
		case error:
			return fmt.Errorf("invalid required flag %q: %w", flag.GetName(), val)
		default:
			return fmt.Errorf("required flag %q has an unexpected type %T", flag.GetName(), val)
		}
	}

	return nil
}

func errIfAnyError(rb *resultBranch) (*resultBranch, error) {
	for arg, val := range rb.resultMap {
		reason, ok := val.(error)
		if ok {
			return rb, fmt.Errorf("flag %q has an error: %w", arg, reason)
		}
	}

	return rb, nil
}
