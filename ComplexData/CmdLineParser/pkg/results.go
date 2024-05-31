package pkg

import (
	"fmt"

	evalSlc "github.com/PlayerR9/MyGoLib/Evaluations/Slices"
	uc "github.com/PlayerR9/MyGoLib/Units/Common"
	ue "github.com/PlayerR9/MyGoLib/Units/Errors"
	ui "github.com/PlayerR9/MyGoLib/Units/Iterator"
	up "github.com/PlayerR9/MyGoLib/Units/Pair"
	us "github.com/PlayerR9/MyGoLib/Units/Slices"
)

type resultArg struct {
	args   []string
	parsed []any
}

func newResultArg(args []string, parsed []any) *resultArg {
	return &resultArg{
		args:   args,
		parsed: parsed,
	}
}

func (r *resultArg) Parsed() []any {
	return r.parsed
}

func (r *resultArg) Args() []string {
	return r.args
}

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
		case map[string][]any:
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

type ciEvaluator struct {
	flags           []*FlagInfo
	args            []string
	flagSeen        []*FlagInfo
	pos             int
	flag            *FlagInfo
	startingIndices []int
}

func (inf *ciEvaluator) Iterator() ui.Iterater[int] {
	return ui.NewSimpleIterator(inf.startingIndices)
}

func (inf *ciEvaluator) Init(args []string) (*resultBranch, error) {
	inf.startingIndices = make([]int, 0)
	inf.args = args
	inf.flagSeen = make([]*FlagInfo, 0)
	inf.pos = len(args)
	inf.flag = nil

	for i := len(args) - 1; i >= 0; i-- {
		arg := args[i]

		flag := inf.getFlag(arg)
		if flag == nil {
			continue
		}

		inf.startingIndices = append(inf.startingIndices, i)
		inf.flagSeen = append(inf.flagSeen, flag)
	}

	// Check if all the required flags are seen.
	for _, flag := range inf.flags {
		if !flag.IsRequired() {
			continue
		}

		index := us.FindEquals(inf.flagSeen, flag)
		if index == -1 {
			return nil, fmt.Errorf("missing required flag %q", flag.GetName())
		}
	}

	return newResultBranch(), nil
}

func (inf *ciEvaluator) Core(index int, lp int) (*up.Pair[map[string][]any, error], error) {
	inf.flag = inf.flagSeen[index]

	newArgs := inf.args[lp+1 : inf.pos] // +1 to skip the flag name itself

	branches, err := evalSlc.Evaluate(inf.flag, newArgs)
	if err != nil {
		return nil, fmt.Errorf("error evaluating flag %q: %w", inf.flag.GetName(), err)
	}

	result, err := inf.flag.Parse(branches, newArgs)
	if err == nil {
		inf.pos = lp
	} else if ue.As[*ue.ErrIgnorable](err) {
		err = err.(*ue.ErrIgnorable).Err
	}

	return up.NewPair(result, err), nil
}

func (inf *ciEvaluator) Next(pair *up.Pair[map[string][]any, error], branch *resultBranch) ([]*resultBranch, error) {
	flagName := inf.flag.GetName()

	prev, ok := branch.resultMap[flagName]
	if !ok {
		// At the first evaluation, we have no previous result and so,
		// we can just store the result as is.
		if pair.Second != nil {
			branch.resultMap[flagName] = pair.Second
		} else {
			branch.resultMap[flagName] = pair.First
		}

		return []*resultBranch{branch}, nil
	}

	if pair.Second != nil {
		_, ok := prev.(error)
		if ok {
			// Prioritize the latest error.
			branch.resultMap[flagName] = pair.Second
		}

		return []*resultBranch{branch}, nil
	} else {
		switch prev := prev.(type) {
		case error:
			// Prioritize the latest result over the error.
			branch.resultMap[flagName] = pair.First

			return []*resultBranch{branch}, nil
		case *FlagParseResult:
			// Possible conflict. Duplicate branch.

			rbCopy := branch.Copy().(*resultBranch)

			rbCopy.resultMap[flagName] = pair.First

			return []*resultBranch{branch, rbCopy}, nil
		default:
			return nil, fmt.Errorf("unexpected type %T", prev)
		}
	}
}

func (inf *ciEvaluator) getFlag(arg string) *FlagInfo {
	for _, flag := range inf.flags {
		if flag.GetName() == arg {
			return flag
		}
	}

	return nil
}

// FlagParseResult represents the result of parsing a flag.
type FlagParseResult struct {
	// Parsed arguments.
	Args map[string][]any
}

// Copier implements the Copier interface.
func (fpr *FlagParseResult) Copy() uc.Copier {
	args := make(map[string][]any)
	for k, v := range fpr.Args {
		args[k] = v
	}

	return &FlagParseResult{args}
}

// NewFlagParseResult creates a new FlagParseResult with the given
// arguments, index, and ignorable boolean.
//
// Parameters:
//   - args: The arguments to parse.
//
// Returns:
//   - *FlagParseResult: A pointer to the new FlagParseResult.
func NewFlagParseResult(args map[string][]any) *FlagParseResult {
	return &FlagParseResult{
		Args: args,
	}
}

// Insert inserts the given arguments into the FlagParseResult.
//
// Parameters:
//   - argName: The name of the argument.
//   - argValue: The value of the argument.
//
// Behaviors:
//   - If the argument already exists, the value is overwritten.
func (fpr *FlagParseResult) Insert(argName string, argValue []any) {
	fpr.Args[argName] = argValue
}

// GetResult returns the parsed arguments.
//
// Returns:
//   - map[string][]any: The parsed arguments.
func (fpr *FlagParseResult) GetResult() map[string][]any {
	return fpr.Args
}

type flgEvaluator struct {
	argList        []*ArgInfo
	startIndices   []int
	args           []string
	currentArgName string
}

func (inf *flgEvaluator) Iterator() ui.Iterater[*ArgInfo] {
	return ui.NewSimpleIterator(inf.argList)
}

func (inf *flgEvaluator) Init(args []string) (*FlagParseResult, error) {
	inf.startIndices = make([]int, len(inf.argList))
	inf.args = args

	return NewFlagParseResult(make(map[string][]any)), nil
}

func (inf *flgEvaluator) Core(index int, lp *ArgInfo) (*up.Pair[[]*resultArg, error], error) {
	inf.currentArgName = lp.GetName()

	var newPosition int

	if index == 0 {
		newPosition = 0
	} else {
		newPosition = inf.startIndices[index-1] + 1
	}

	results, err := lp.Parse(inf.args[newPosition:])

	return up.NewPair(results, err), nil
}

func (inf *flgEvaluator) Next(pair *up.Pair[[]*resultArg, error], branch *FlagParseResult) ([]*FlagParseResult, error) {
	if pair.Second != nil {
		// Current branch is invalid.
		return nil, nil
	}

	var newBranches []*FlagParseResult

	for _, result := range pair.First {
		branchCopy := branch.Copy().(*FlagParseResult)

		branchCopy.Insert(inf.currentArgName, result.Parsed())

		newBranches = append(newBranches, branchCopy)
	}

	return newBranches, nil
}
