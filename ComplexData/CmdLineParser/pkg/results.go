package pkg

import (
	"fmt"

	evalSlc "github.com/PlayerR9/MyGoLib/Evaluations/Slices"
	ui "github.com/PlayerR9/MyGoLib/Units/Iterators"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
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
	// map[FLAG]map[ARG]parsed
	resultMap map[string]*FlagParseResult

	reason error

	// argsDone is a list of arguments that have been parsed so far.
	argsDone []string
}

func (rb *resultBranch) Copy() uc.Copier {
	argsDone := make([]string, len(rb.argsDone))
	copy(argsDone, rb.argsDone)

	resultMap := make(map[string]*FlagParseResult)

	for k, v := range rb.resultMap {
		values := v.Copy().(*FlagParseResult)

		resultMap[k] = values
	}

	return &resultBranch{
		resultMap: resultMap,
		argsDone:  argsDone,
		reason:    rb.reason,
	}
}

func newResultBranch(argMap map[string]*FlagParseResult, reason error, argsDone []string) *resultBranch {
	rb := &resultBranch{
		reason: reason,
	}

	if len(argsDone) > 0 {
		rb.argsDone = argsDone
	} else {
		rb.argsDone = make([]string, 0)
	}

	if reason == nil && argMap != nil {
		rb.resultMap = argMap
	} else {
		rb.resultMap = make(map[string]*FlagParseResult)
	}

	return rb
}

func (rb *resultBranch) getResultMap() (map[string]*FlagParseResult, error) {
	if rb.reason != nil {
		return nil, rb.reason
	} else {
		return rb.resultMap, nil
	}
}

func (rb *resultBranch) getArgumentsDone() []string {
	return rb.argsDone
}

/*
func (rb *resultBranch) changeReason(err error) {
	if err != nil {
		rb.reason = err
		rb.resultMap = make(map[string]*FlagParseResult)
	} else {
		rb.reason = nil
	}
}

func (rb *resultBranch) setResultMap(result map[string]*FlagParseResult) {
	rb.resultMap = result
}
*/

func (rb *resultBranch) hasFlag(flagName string) *FlagParseResult {
	values, ok := rb.resultMap[flagName]
	if !ok {
		return nil
	}

	return values
}

func (rb *resultBranch) errIfInvalidRequiredFlags(flags []*FlagInfo) error {
	for _, flag := range flags {
		if !flag.IsRequired() {
			continue
		}

		flagName := flag.GetName()

		val := rb.hasFlag(flagName)
		if val == nil {
			return fmt.Errorf("missing required flag %q", flagName)
		}

		if rb.reason != nil {
			return fmt.Errorf("invalid required flag %q: %w", flagName, rb.reason)
		}
	}

	return nil
}

func (rb *resultBranch) merge(key string, value *FlagParseResult) {
	rb.resultMap[key] = value

	rb.argsDone = append(rb.argsDone, value.GetArgumentsDone()...)
}

func errIfAnyError(rb *resultBranch) (*resultBranch, error) {
	return rb, rb.reason
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

	return newResultBranch(nil, nil, nil), nil
}

func (inf *ciEvaluator) Core(index int, lp int) (*uc.Pair[[]*FlagParseResult, error], error) {
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

	p := uc.NewPair(result, err)

	return &p, nil

}

func (inf *ciEvaluator) Next(pair *uc.Pair[[]*FlagParseResult, error], branch *resultBranch) ([]*resultBranch, error) {
	var newBranches []*resultBranch

	flagName := inf.flag.GetName()

	argumentsDone := branch.getArgumentsDone()
	size := len(argumentsDone)
	argumentsDone = append(argumentsDone, flagName)

	if size == 0 {
		var newBranch *resultBranch

		// At the first evaluation, we have no previous result and so,
		// we can just store the result as is.
		if pair.Second != nil {
			newBranch = newResultBranch(nil, pair.Second, argumentsDone)

			newBranches = append(newBranches, newBranch)
		} else {
			for _, result := range pair.First {
				initialMap := map[string]*FlagParseResult{
					flagName: result,
				}

				newBranch = newResultBranch(initialMap, nil, argumentsDone)

				newBranches = append(newBranches, newBranch)
			}
		}
	} else {
		for _, result := range pair.First {
			branchCopy := branch.Copy().(*resultBranch)

			branchCopy.merge(flagName, result)

			newBranches = append(newBranches, branchCopy)
		}

		if pair.Second == nil && branch.reason == nil {
			// Possible conflict. Duplicate branch.

			newBranches = append(newBranches, branch)
		}
	}

	return newBranches, nil
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
	// argMap is a map of argument names to their values.
	argMap map[string][]any

	// argumentsDone is a list of arguments that have been parsed so far.
	argumentsDone []string
}

// Copier implements the Copier interface.
func (fpr *FlagParseResult) Copy() uc.Copier {
	argMap := make(map[string][]any)
	for k, v := range fpr.argMap {
		argMap[k] = v
	}
	argsDone := make([]string, len(fpr.argumentsDone))
	copy(argsDone, fpr.argumentsDone)

	return &FlagParseResult{
		argMap:        argMap,
		argumentsDone: argsDone,
	}
}

// NewFlagParseResult creates a new FlagParseResult with the given
// arguments, index, and ignorable boolean.
//
// Parameters:
//   - argMap: A map of argument names to their values.
//   - argsDone: A list of arguments that have been parsed so far.
//
// Returns:
//   - *FlagParseResult: A pointer to the new FlagParseResult.
func NewFlagParseResult() *FlagParseResult {
	return &FlagParseResult{
		argMap:        make(map[string][]any),
		argumentsDone: make([]string, 0),
	}
}

/*
func (fpr *FlagParseResult) hasFlag(argName string) bool {
	_, ok := fpr.argMap[argName]
	return ok
}
*/

// insert inserts the given arguments into the FlagParseResult.
//
// Parameters:
//   - argName: The name of the argument.
//   - argValue: The value of the argument.
//
// Returns:
//   - bool: true if the argument was inserted, false otherwise.
func (fpr *FlagParseResult) insert(argName string, res *resultArg) bool {
	_, ok := fpr.argMap[argName]
	if ok {
		return false
	}

	fpr.argMap[argName] = res.Parsed()
	fpr.argumentsDone = append(fpr.argumentsDone, res.Args()...)

	return true
}

// GetResult returns the parsed arguments.
//
// Returns:
//   - map[string][]any: The parsed arguments.
func (fpr *FlagParseResult) GetResult() map[string][]any {
	return fpr.argMap
}

// GetArgumentsDone returns the arguments that have been parsed so far.
//
// Returns:
//   - []string: The arguments that have been parsed so far.
func (fpr *FlagParseResult) GetArgumentsDone() []string {
	return fpr.argumentsDone
}

func (fpr *FlagParseResult) size() int {
	var size int

	for _, v := range fpr.argMap {
		size += len(v)
	}

	return size
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

	return NewFlagParseResult(), nil
}

func (inf *flgEvaluator) Core(index int, lp *ArgInfo) (*uc.Pair[[]*resultArg, error], error) {
	inf.currentArgName = lp.GetName()

	var newPosition int

	if index == 0 {
		newPosition = 0
	} else {
		newPosition = inf.startIndices[index-1] + 1
	}

	results, err := lp.Parse(inf.args[newPosition:])

	p := uc.NewPair(results, err)

	return &p, nil
}

func (inf *flgEvaluator) Next(pair *uc.Pair[[]*resultArg, error], branch *FlagParseResult) ([]*FlagParseResult, error) {
	if pair.Second != nil {
		// Current branch is invalid.
		return nil, nil
	}

	var newBranches []*FlagParseResult

	for _, result := range pair.First {
		branchCopy := branch.Copy().(*FlagParseResult)

		ok := branchCopy.insert(inf.currentArgName, result)
		if !ok {
			return nil, fmt.Errorf("duplicate argument %q", inf.currentArgName)
		}

		newBranches = append(newBranches, branchCopy)
	}

	return newBranches, nil
}
