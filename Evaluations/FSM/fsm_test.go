package FSM

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	ub "github.com/PlayerR9/MyGoLib/Units/Debugging"
	ut "github.com/PlayerR9/MyGoLib/Units/Tray"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
	mext "github.com/PlayerR9/MyGoLib/Utility/MathExt"
)

func TestEncode(t *testing.T) {
	const (
		TestNumber int = 3966
	)

	str, err := Encoderer(TestNumber)
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	if str.String() != "uuounouu" {
		t.Fatalf("Expected nuouu, got %s", str)
	}
}

type ZPacket struct {
	header ZSS
	value  int
}

type ParseResult struct {
	binary  []int
	packets []ZPacket
	reason  error
}

func NewParseResult(binary []int) ParseResult {
	return ParseResult{
		binary: binary,
	}
}

func (pr *ParseResult) MergeSolution(packet ParseResult) ParseResult {
	return ParseResult{
		binary:  packet.binary,
		packets: append(pr.packets, packet.packets...),
		reason:  packet.reason,
	}
}

func (pr *ParseResult) IsDone() bool {
	return len(pr.binary) == 0
}

func (pr *ParseResult) HasError() bool {
	return pr.reason != nil
}

func (pr *ParseResult) GetError() error {
	return pr.reason
}

func (pr *ParseResult) GetNBits(n int) ([]int, error) {
	if len(pr.binary) < n {
		return nil, fmt.Errorf("expected at least %d bits, got %d", n, len(pr.binary))
	}

	return pr.binary[:n], nil
}

func Filter(prs []ParseResult) ([]ParseResult, bool) {
	top := 0

	for i := 0; i < len(prs); i++ {
		if !prs[i].HasError() {
			prs[top] = prs[i]
			top++
		}
	}

	if top == 0 {
		return prs, false
	} else {
		return prs[:top], true
	}
}

func Select(left, right []ParseResult) []ParseResult {
	ps := append(left, right...)

	results := []ParseResult{
		ps[0],
	}
	limit := len(ps[0].binary)

	for _, p := range ps[1:] {
		if len(p.binary) < limit {
			results = []ParseResult{
				p,
			}

			limit = len(p.binary)
		} else if len(p.binary) == limit {
			results = append(results, p)
		}
	}

	return results
}

type ParseFunc func(ParseResult) []ParseResult

type Rule struct {
	lhs string
	rhs []string
}

func NewRule(str string) (*Rule, error) {
	fields := strings.Split(str, "->")

	if len(fields) == 0 {
		return nil, fmt.Errorf("missing either lhs or rhs")
	} else if len(fields) > 2 {
		return nil, fmt.Errorf("too many fields: expected 2, got %d", len(fields))
	}

	lhs := strings.TrimSpace(fields[0])

	rhs := strings.Fields(fields[1])
	rhs = us.RemoveEmpty(rhs)
	if len(rhs) == 0 {
		return nil, fmt.Errorf("missing rhs")
	}

	return &Rule{
		lhs: lhs,
		rhs: rhs,
	}, nil
}

type Grammar struct {
	rules []*Rule
}

func NewGrammar() *Grammar {
	return &Grammar{
		rules: make([]*Rule, 0),
	}
}

func (g *Grammar) AddRule(rule *Rule) {
	g.rules = append(g.rules, rule)
}

func (g *Grammar) MakeRule(str string) error {
	rule, err := NewRule(str)
	if err != nil {
		return err
	}

	g.AddRule(rule)

	return nil
}

var (
	ZSGrammar *Grammar
)

func init() {
	ZSGrammar = NewGrammar()

	ZSGrammar.MakeRule("R1a -> 0 0") // R1a -> P0
	ZSGrammar.MakeRule("R1a -> 0 0 R1a") // R1a -> P0 R1a

	return solution
}

func parseR1(pr ParseResult) []ParseResult {
	tmp := parseR1a(pr)

	tmp, ok := Filter(tmp)
	if !ok {
		return tmp
	}

	var results []ParseResult
	var todo []ParseResult

	for _, t := range tmp {
		if t.IsDone() {
			// R1 -> R1a : R1 -> P0
			results = append(results, t)
		} else {
			todo = append(todo, t)
		}
	}

	for _, t := range todo {
		// R1 -> R1a R2
		tmp := parseR2(t)

		for _, t2 := range tmp {
			t2 = t.MergeSolution(t2)

			results = append(results, t2)
		}

		// R1 -> R1a R4
		tmp = parseR4(t)

		for _, t2 := range tmp {
			t2 = t.MergeSolution(t2)

			results = append(results, t2)
		}
	}

	return results
}

func parseR2a(binary []int) ([]ZPacket, []int, error) {
	// R2a -> P1
	// R2a -> P1 R2a
}

func parseR2(pr ParseResult) []ParseResult {
	tmp := parseR2a(pr)

	tmp, ok := Filter(tmp)
	if !ok {
		return tmp
	}

	var results []ParseResult
	var todo []ParseResult

	for _, t := range tmp {
		if t.IsDone() {
			// R2 -> R2a
			results = append(results, t)
		} else {
			todo = append(todo, t)
		}
	}

	for _, t := range todo {
		// R2 -> R2a R1
		tmp := parseR1(t)

		for _, t2 := range tmp {
			t2 = t.MergeSolution(t2)

			results = append(results, t2)
		}

		// R2 -> R2a R3
		tmp = parseR3(t)

		for _, t2 := range tmp {
			t2 = t.MergeSolution(t2)

			results = append(results, t2)
		}
	}

	return results

}

func parseR3a(binary []int) ([]ZPacket, []int, error) {
	if len(binary) == 0 {
		return nil, binary, fmt.Errorf("expected at least 1 bit, got 0")
	}

	if binary[0] != 0 {
		return nil, binary, fmt.Errorf("expected 0, got %d", binary[0])
	}

	var result []ZPacket

	result = append(result, ZPacket{
		header: NewZss(NoPacket, false),
		value:  1,
	})

	binary = binary[1:]

	if len(binary) == 0 {
		// R3a -> N0 : R3a -> 0
		return result, binary, nil
	}

	if binary[0] != 1 {
		return result, binary, fmt.Errorf("expected 1, got %d", binary[0])
	}

	binary = binary[1:]

	result = []ZPacket{
		{
			header: NewZss(NoPacket, false),
			value:  2,
		},
	}

	if len(binary) == 0 {
		// R3a -> N0 R3a : R3a -> 0 1
		return result, binary, nil
	}

	tmp, end, err := parseR3a(binary)
	if err != nil {
		return result, binary, nil
	}

	// R3a -> N0 R3a : R3a -> 0 1 R3a
	result = append(result, tmp...)

	return result, end, nil

}

func parseR3(pr ParseResult) []ParseResult {
	tmp := parseR3a(pr)

	tmp, ok := Filter(tmp)
	if !ok {
		return tmp
	}

	var results []ParseResult
	var todo []ParseResult

	for _, t := range tmp {
		if t.IsDone() {
			// R3 -> R3a
			results = append(results, t)
		} else {
			todo = append(todo, t)
		}
	}

	for _, t := range todo {
		// R3 -> R3a R5
		tmp1 := parseR5(t)

		for _, t2 := range tmp1 {
			t2 = t.MergeSolution(t2)

			results = append(results, t2)
		}
	}

	return results
}

func parseR4a(pr ParseResult) []ParseResult {
	bits, err := pr.GetNBits(1)
	if err != nil {
		res := ParseResult{
			binary:  pr.binary,
			packets: []ZPacket{},
			reason:  err,
		}

		return []ParseResult{res}
	}

	if bits[0] != 1 {
		res := ParseResult{
			binary:  pr.binary,
			packets: []ZPacket{},
			reason:  fmt.Errorf("expected 1, got %d", bits[0]),
		}

		return []ParseResult{res}
	}

	// R4a -> N1 : R4a -> 1
	res := ParseResult{
		binary: pr.binary[1:],
		packets: []ZPacket{
			{
				header: NewZss(NoPacket, true),
				value:  1,
			},
		},
	}

	solution := []ParseResult{res}

	bits, err = res.GetNBits(1)
	if err != nil {
		return solution
	}

	if bits[0] != 0 {
		return solution
	}

	// R4a -> N1 R4a : R4a -> 1 0
	res = ParseResult{
		binary: res.binary[1:],
		packets: []ZPacket{
			{
				header: NewZss(NoPacket, true),
				value:  2,
			},
		},
	}

	// R4a -> N1 R4a : R4a -> 1 0 R4a
	tmp := parseR4a(res)

	tmp, ok := Filter(tmp)
	if !ok {
		return solution
	}

	for _, t := range tmp {
		t = res.MergeSolution(t)

		solution = append(solution, t)
	}

	return solution
}

func parseR4(pr ParseResult) []ParseResult {
	tmp1 := parseR4a(pr)

	tmp1, ok := Filter(tmp1)
	if !ok {
		return tmp1
	}

	var results []ParseResult
	var todo []ParseResult

	for _, t := range tmp1 {
		if t.IsDone() {
			// R4 -> R4a
			results = append(results, t)
		} else {
			todo = append(todo, t)
		}
	}

	for _, t := range todo {
		// R4 -> R4a R5
		tmp2 := parseR5(t)

		for _, t2 := range tmp2 {
			t2 = t.MergeSolution(t2)

			results = append(results, t2)
		}
	}

	return results
}

func parseR5(pr ParseResult) []ParseResult {
	tmp1 := parseR1(pr)
	tmp2 := parseR2(pr)

	tmp1, ok1 := Filter(tmp1)
	tmp2, ok2 := Filter(tmp2)

	if ok1 == ok2 {
		res := Select(tmp1, tmp2)

		return res
	}

	if ok1 {
		// R5 -> R1
		return tmp1
	} else {
		// R5 -> R2
		return tmp2
	}
}
}

func parseR1a(pr ParseResult) []ParseResult {
	bits, err := pr.GetNBits(2)
	if err != nil {
		res := ParseResult{
			binary:  pr.binary,
			packets: []ZPacket{},
			reason:  err,
		}

		return []ParseResult{res}
	}

	if bits[0] != 0 || bits[1] != 0 {
		res := ParseResult{
			binary:  pr.binary,
			packets: []ZPacket{},
			reason:  fmt.Errorf("expected '0 0', got %d %d", bits[0], bits[1]),
		}

		return []ParseResult{res}
	}

	res := ParseResult{
		binary: pr.binary[2:],
		packets: []ZPacket{
			{
				header: NewZss(Packet, false),
				value:  1,
			},
		},
	}

	solution := []ParseResult{res}

	pr.binary = res.binary

	if pr.IsDone() {
		// R1a -> P0 : R1a -> 0 0
		return solution
	}

	// R1a -> P0 R1a : R1a -> 0 0 R1a
	tmp := parseR1a(pr)

	tmp, ok := Filter(tmp)
	if !ok {
		return solution
	}

	for _, t := range tmp {
		t = res.MergeSolution(t)

		solution = append(solution, t)
	}

	return solution
}

func parseR1(pr ParseResult) []ParseResult {
	tmp := parseR1a(pr)

	tmp, ok := Filter(tmp)
	if !ok {
		return tmp
	}

	var results []ParseResult
	var todo []ParseResult

	for _, t := range tmp {
		if t.IsDone() {
			// R1 -> R1a : R1 -> P0
			results = append(results, t)
		} else {
			todo = append(todo, t)
		}
	}

	for _, t := range todo {
		// R1 -> R1a R2
		tmp := parseR2(t)

		for _, t2 := range tmp {
			t2 = t.MergeSolution(t2)

			results = append(results, t2)
		}

		// R1 -> R1a R4
		tmp = parseR4(t)

		for _, t2 := range tmp {
			t2 = t.MergeSolution(t2)

			results = append(results, t2)
		}
	}

	return results
}

func parseR2a(binary []int) ([]ZPacket, []int, error) {
	// R2a -> P1
	// R2a -> P1 R2a
}

func parseR2(pr ParseResult) []ParseResult {
	tmp := parseR2a(pr)

	tmp, ok := Filter(tmp)
	if !ok {
		return tmp
	}

	var results []ParseResult
	var todo []ParseResult

	for _, t := range tmp {
		if t.IsDone() {
			// R2 -> R2a
			results = append(results, t)
		} else {
			todo = append(todo, t)
		}
	}

	for _, t := range todo {
		// R2 -> R2a R1
		tmp := parseR1(t)

		for _, t2 := range tmp {
			t2 = t.MergeSolution(t2)

			results = append(results, t2)
		}

		// R2 -> R2a R3
		tmp = parseR3(t)

		for _, t2 := range tmp {
			t2 = t.MergeSolution(t2)

			results = append(results, t2)
		}
	}

	return results

}

func parseR3a(binary []int) ([]ZPacket, []int, error) {
	if len(binary) == 0 {
		return nil, binary, fmt.Errorf("expected at least 1 bit, got 0")
	}

	if binary[0] != 0 {
		return nil, binary, fmt.Errorf("expected 0, got %d", binary[0])
	}

	var result []ZPacket

	result = append(result, ZPacket{
		header: NewZss(NoPacket, false),
		value:  1,
	})

	binary = binary[1:]

	if len(binary) == 0 {
		// R3a -> N0 : R3a -> 0
		return result, binary, nil
	}

	if binary[0] != 1 {
		return result, binary, fmt.Errorf("expected 1, got %d", binary[0])
	}

	binary = binary[1:]

	result = []ZPacket{
		{
			header: NewZss(NoPacket, false),
			value:  2,
		},
	}

	if len(binary) == 0 {
		// R3a -> N0 R3a : R3a -> 0 1
		return result, binary, nil
	}

	tmp, end, err := parseR3a(binary)
	if err != nil {
		return result, binary, nil
	}

	// R3a -> N0 R3a : R3a -> 0 1 R3a
	result = append(result, tmp...)

	return result, end, nil

}

func parseR3(pr ParseResult) []ParseResult {
	tmp := parseR3a(pr)

	tmp, ok := Filter(tmp)
	if !ok {
		return tmp
	}

	var results []ParseResult
	var todo []ParseResult

	for _, t := range tmp {
		if t.IsDone() {
			// R3 -> R3a
			results = append(results, t)
		} else {
			todo = append(todo, t)
		}
	}

	for _, t := range todo {
		// R3 -> R3a R5
		tmp1 := parseR5(t)

		for _, t2 := range tmp1 {
			t2 = t.MergeSolution(t2)

			results = append(results, t2)
		}
	}

	return results
}

func parseR4a(pr ParseResult) []ParseResult {
	bits, err := pr.GetNBits(1)
	if err != nil {
		res := ParseResult{
			binary:  pr.binary,
			packets: []ZPacket{},
			reason:  err,
		}

		return []ParseResult{res}
	}

	if bits[0] != 1 {
		res := ParseResult{
			binary:  pr.binary,
			packets: []ZPacket{},
			reason:  fmt.Errorf("expected 1, got %d", bits[0]),
		}

		return []ParseResult{res}
	}

	// R4a -> N1 : R4a -> 1
	res := ParseResult{
		binary: pr.binary[1:],
		packets: []ZPacket{
			{
				header: NewZss(NoPacket, true),
				value:  1,
			},
		},
	}

	solution := []ParseResult{res}

	bits, err = res.GetNBits(1)
	if err != nil {
		return solution
	}

	if bits[0] != 0 {
		return solution
	}

	// R4a -> N1 R4a : R4a -> 1 0
	res = ParseResult{
		binary: res.binary[1:],
		packets: []ZPacket{
			{
				header: NewZss(NoPacket, true),
				value:  2,
			},
		},
	}

	// R4a -> N1 R4a : R4a -> 1 0 R4a
	tmp := parseR4a(res)

	tmp, ok := Filter(tmp)
	if !ok {
		return solution
	}

	for _, t := range tmp {
		t = res.MergeSolution(t)

		solution = append(solution, t)
	}

	return solution
}

func parseR4(pr ParseResult) []ParseResult {
	tmp1 := parseR4a(pr)

	tmp1, ok := Filter(tmp1)
	if !ok {
		return tmp1
	}

	var results []ParseResult
	var todo []ParseResult

	for _, t := range tmp1 {
		if t.IsDone() {
			// R4 -> R4a
			results = append(results, t)
		} else {
			todo = append(todo, t)
		}
	}

	for _, t := range todo {
		// R4 -> R4a R5
		tmp2 := parseR5(t)

		for _, t2 := range tmp2 {
			t2 = t.MergeSolution(t2)

			results = append(results, t2)
		}
	}

	return results
}

func parseR5(pr ParseResult) []ParseResult {
	tmp1 := parseR1(pr)
	tmp2 := parseR2(pr)

	tmp1, ok1 := Filter(tmp1)
	tmp2, ok2 := Filter(tmp2)

	if ok1 == ok2 {
		res := Select(tmp1, tmp2)

		return res
	}

	if ok1 {
		// R5 -> R1
		return tmp1
	} else {
		// R5 -> R2
		return tmp2
	}
}

func decodeSequence(binary []int) {
	isPacket := binary[0] == binary[1]

	/*

		0 -> (N0, 1)
		1 -> (N1, 1)

		00 -> (P0, 1)
		01 -> (N0, 2)
		10 -> (N1, 2)
		11 -> (P1, 1)

		000 -> (N0, 1), (P0, 1)
		001 -> (P0, 1), (N1, 1)
		010 -> (N0, 3)
		011 -> (N0, 1), (P1, 1)
		100 -> (N1, 1), (P0, 1)
		101 -> (N1, 3)
		110 -> (P1, 1), (N0, 1)
		111 -> (N1, 1), (P1, 1)

		S -> R1 | R2 | R3 | R4

		R1 -> R1a (R2 | R4)?
		R1a -> P0 | P0 R1a
		R2 -> P1 P1* (R1 | R3)?
		R3 -> N0 N0* R5?
		R4 -> N1 N1* R5?

		R5 -> R1 | R2

		N0 :
			- 0
			- 0, 0

		N0: 0, 01
			- 0:
				0 ->
				00 ->
				1 ->
				11 ->
			- 01:
				0 ->
				00 ->
				1 ->
				11 ->
		N1: 1, 10
			- 1:
				0 ->
				00 ->
				1 ->
				11 ->
			- 10:
				0 ->
				00 ->
				1 ->
				11 ->
		P0: 00
			- 0 -> invalid. It should have been 0, 00
			- 00 -> P0
			- 1 -> N1
			- 11 -> P1
		P1: 11
			- 0 -> N0
			- 00 -> P0
			- 1 -> invalid. It should have been 1, 11
			- 11 -> P1

	*/
}

func Encoderer(num int) (*Number, error) {
	num += 2 // offset

	binary, err := mext.DecToBase(num, Base)
	if err != nil {
		return nil, fmt.Errorf("error converting number to binary: %w", err)
	}

	binary = binary[:len(binary)-1] // remove the last bit

	slices.Reverse(binary) // reverse the binary

	stream, err := NewEncoderStream(num)
	if err != nil {
		return nil, fmt.Errorf("error creating encoder stream: %w", err)
	}

	digits, err := Encoder.Run(stream)
	if err != nil {
		return nil, fmt.Errorf("error running encoder: %w", err)
	}

	encoded := NewNumber()

	for _, digit := range digits {
		encoded.AppendDigit(digit)
	}

	/*

		if num < 0 {
			return nil, ue.NewErrInvalidParameter("number", ue.NewErrGTE(0))
		}

		encoded := NewNumber()

		// convert <num> to binary and remove the last bit
		// add 2 to the number (offset)
		binary, err := mext.DecToBase(num+2, Base)
		if err != nil {
			return nil, fmt.Errorf("error converting number to binary: %w", err)
		}

		binary = binary[:len(binary)-1]

		state := InitialState

		state = state.SetBit(binary[0] == 1) // initialize the state with the first bit

		for len(binary) > 0 {
			header, err := DetermineHeader(state, binary)
			if err != nil {
				return nil, fmt.Errorf("error determining header: %w", err)
			}

			nextState, err := DetermineNextState(state, binary)
			if err != nil {
				return nil, fmt.Errorf("error getting next state: %w", err)
			}

			tmp = nextState

			packetType := state.GetState()

			half, C, err := GetPacket(packetType, binary)
			if err != nil {
				return nil, fmt.Errorf("error getting packet: %w", err)
			}

			binary = half[C:]

			digit, err := NewDigit(header, C-1)
			if err != nil {
				return nil, fmt.Errorf("error creating digit: %w", err)
			}

			encoded.AppendDigit(digit)

			state = tmp
		}
	*/

	return encoded, nil
}

// given a slice of boolean values checks if it starts with a packet
// panics if the slice is empty
func StartsWithPacket(binary []int) (bool, error) {
	if len(binary) == 0 {
		return false, ue.NewErrInvalidParameter("binary", ue.NewErrEmpty(binary))
	}

	ok := len(binary) != 1 && binary[0] == binary[1]

	return ok, nil
}

// given a slice of boolean values returns how many bits are in the packet
// panics if the slice is empty
func GetPacket(packetType ZssState, binary []int) ([]int, int, error) {
	if len(binary) == 0 {
		return nil, 0, ue.NewErrInvalidParameter("binary", ue.NewErrEmpty(binary))
	}

	var counter int

	switch packetType {
	case Packet:
		if len(binary) == 1 {
			return nil, 0, ue.NewErrInvalidParameter("binary", errors.New("only one bit 1"))
		}

		for i := 0; i < len(binary)-1; i++ {
			if binary[i] != binary[i+1] {
				break
			}

			counter++
		}

		binary = binary[1:]
	case NoPacket:
		var abruptStop bool

		for i := 0; i < len(binary)-1; i++ {
			if binary[i] == binary[i+1] {
				abruptStop = true
				break
			}

			counter++
		}

		if !abruptStop {
			counter++
		}
	default:
		return nil, 0, fmt.Errorf("state %d is not valid", packetType)
	}

	return binary, counter, nil
}

func GetSect(zss_text []rune) []int {
	var binary []int
	var next rune

	for i, r := range zss_text {
		if i < len(zss_text)-1 {
			next = zss_text[i+1]
		} else {
			next = 'o'
		}

		if r == 'n' {
			binary = append(binary, 1)
		} else {
			binary = append(binary, 0)
		}

		if next == 'o' {
			break
		}
	}

	/*

		fmt.Println("binary:")

		for _, r := range binary {
			fmt.Printf("%t ", r)
		}

		fmt.Println()*/

	return binary
}

type ZssState int

const (
	ZssStateZero ZssState = iota
	Packet
	NoPacket
)

var InitialState ZSS

func init() {
	InitialState = ZSS{
		state: ZssStateZero,
		bit:   false,
	}
}

type ZSS struct {
	state ZssState
	bit   bool
}

func (o ZSS) String() string {
	var builder strings.Builder

	switch o.state {
	case ZssStateZero:
		builder.WriteRune('Z')
	case Packet:
		builder.WriteRune('P')
	case NoPacket:
		builder.WriteRune('N')
	}

	if o.bit {
		builder.WriteRune('1')
	} else {
		builder.WriteRune('0')
	}

	return builder.String()
}

func NewZss(state ZssState, bit bool) ZSS {
	return ZSS{
		state: state,
		bit:   bit,
	}
}

func (o ZSS) SetBit(bit bool) ZSS {
	return ZSS{
		state: o.state,
		bit:   bit,
	}
}

func (o ZSS) GetState() ZssState {
	return o.state
}

func (o ZSS) GetBit() bool {
	return o.bit
}

func (o ZSS) IsPacket() bool {
	return o.state == Packet
}

type ZssElement int

const (
	ElHeader ZssElement = iota
	ElNextState
	ElEndOfPacket
)

func (e ZssElement) String() string {
	return [...]string{
		"Header",
		"NextState",
		"EndOfPacket",
	}[e]
}

var (
	DetermineHeader      DetFunc[int, ZSS, ZssElement]
	DetermineNextState   DetFunc[int, ZSS, ZssElement]
	DetermineEndOfPacket DetFunc[int, ZSS, ZssElement]
)

func init() {
	DetermineHeader = func(fsm *ActiveFSM[int, ZSS, ZssElement]) (any, error) {
		debugLines := DebugPrintTray(fsm.Tray)

		fmt.Println(strings.Join(debugLines, "\n"))

		var header Header

		f := func(t ut.Trayer[int]) (bool, error) {
			firstDigit, err := fsm.Tray.Read()
			if err != nil {
				return false, fmt.Errorf("error reading first digit: %w", err)
			}

			var secondDigit *int

			rem := fsm.Tray.Move(1)
			if rem == 1 {
				secondDigit = nil
			} else {
				tmp, err := fsm.Tray.Read()
				if err != nil {
					secondDigit = nil
				} else {
					secondDigit = &tmp
				}
			}

			currState := fsm.GetState()
			packetType := currState.GetState()

			header.isPacket = packetType != ZssStateZero || (secondDigit != nil && firstDigit == *secondDigit)

			bit := packetType != NoPacket && currState.GetBit()

			var startsWithPacket bool

			switch packetType {
			case ZssStateZero:
				startsWithPacket = false
			case Packet:
				startsWithPacket = secondDigit != nil && firstDigit == *secondDigit
			case NoPacket:
				startsWithPacket = firstDigit == 1
			default:
				return false, fmt.Errorf("state %d is not valid", packetType)
			}

			if bit != startsWithPacket {
				header.digit = UpDigit
			} else {
				header.digit = DownDigit
			}

			return true, nil
		}

		err := ub.DoWithBackup(fsm.Tray, f)
		if err != nil {
			return header, fmt.Errorf("error determining header: %w", err)
		}

		return header, nil
	}

	DetermineNextState = func(fsm *ActiveFSM[int, ZSS, ZssElement]) (any, error) {
		var nextState ZSS

		f := func(tray ut.Trayer[int]) (bool, error) {
			firstDigit, err := fsm.Tray.Read()
			if err != nil {
				return false, fmt.Errorf("error reading first digit: %w", err)
			}

			var secondDigit *int

			rem := fsm.Tray.Move(1)
			if rem == 1 {
				secondDigit = nil
			} else {
				tmp, err := fsm.Tray.Read()
				if err != nil {
					secondDigit = nil
				} else {
					secondDigit = &tmp
				}
			}

			currState := fsm.GetState()
			packetType := currState.GetState()

			var ns ZssState

			startsWithPacket := packetType == NoPacket || (secondDigit != nil && firstDigit == *secondDigit)
			if startsWithPacket {
				ns = Packet
			} else {
				ns = NoPacket
			}

			var bit bool

			switch packetType {
			case ZssStateZero:
				bit = currState.GetBit()
			case Packet:
				bit = currState.GetBit()
				bit = !bit
			case NoPacket:
				bit = firstDigit == 1
			default:
				return false, fmt.Errorf("state %d is not valid", packetType)
			}

			nextState = NewZss(ns, bit)

			return true, nil
		}

		err := ub.DoWithBackup(fsm.Tray, f)
		if err != nil {
			return nextState, fmt.Errorf("error determining next state: %w", err)
		}

		return nextState, nil
	}

	DetermineEndOfPacket = func(fsm *ActiveFSM[int, ZSS, ZssElement]) (any, error) {
		var counter int

		f := func(tray ut.Trayer[int]) (bool, error) {
			val, ok := fsm.GetValue(ElNextState)
			if !ok {
				return false, fmt.Errorf("no value for element %s", ElNextState.String())
			}

			currState := val.(ZSS)
			packetType := currState.GetState()

			var counter int

			prevDigit, err := fsm.Tray.Read()
			if err != nil {
				return false, fmt.Errorf("error reading first digit: %w", err)
			}

			switch packetType {
			case Packet:
				rem := fsm.Tray.Move(1)
				if rem == 1 {
					return false, ue.NewErrInvalidParameter("binary", errors.New("only one bit 1"))
				}

				for {
					// Check if we are at the end of the tray. If so,
					// ignore the last element.
					rem := fsm.Tray.Move(1)
					if rem == 1 {
						break
					}

					fsm.Tray.Move(-1) // move back to the original position

					currDigit, err := fsm.Tray.Read()
					if err != nil {
						return false, fmt.Errorf("error reading current digit: %w", err)
					}

					if prevDigit != currDigit {
						break
					}

					counter++
					prevDigit = currDigit

					fsm.Tray.Delete(1) // delete the current element
				}

				fsm.Tray.Move(-1)
				fsm.Tray.Delete(2) // delete the first  two elements
			case NoPacket:
				var abruptStop bool

				rem := fsm.Tray.Move(1)
				if rem == 1 {
					return true, nil
				}

				for {
					// Check if we are at the end of the tray. If so,
					// ignore the last element.
					rem = fsm.Tray.Move(1)
					if rem == 1 {
						break
					}
					fsm.Tray.Move(-1) // move back to the original position

					currDigit, err := fsm.Tray.Read()
					if err != nil {
						return false, fmt.Errorf("error reading current digit: %w", err)
					}

					if prevDigit == currDigit {
						abruptStop = true
						break
					}

					counter++
					prevDigit = currDigit

					fsm.Tray.Delete(1) // delete the current element
				}

				// Move back to the original position
				fsm.Tray.Move(-1)
				fsm.Tray.Delete(1) // delete the first element

				if !abruptStop {
					counter++
				}
			default:
				return false, fmt.Errorf("state %d is not valid", packetType)
			}

			return true, nil
		}

		err := ub.DoWithBackup(fsm.Tray, f)
		if err != nil {
			return counter, fmt.Errorf("error determining end of packet: %w", err)
		}

		return counter, nil
	}

}

var (
	Init          InitFunc[int, ZSS]
	EvaluateDigit EvalFunc[int, ZSS, Digit, ZssElement]
	Next          TransFunc[int, ZSS, ZssElement]
)

func init() {
	Init = func(t ut.Trayer[int]) (ZSS, error) {
		state := InitialState

		value, err := t.Read()
		if err != nil {
			return state, fmt.Errorf("error reading first digit: %w", err)
		}

		state = state.SetBit(value == 1)

		return state, nil
	}

	EvaluateDigit = func(fsm *ActiveFSM[int, ZSS, ZssElement]) (Digit, error) {
		headerVal, ok := fsm.GetValue(ElHeader)
		if !ok {
			return Digit{}, fmt.Errorf("no value for element %s", ElHeader.String())
		}

		eopVal, ok := fsm.GetValue(ElEndOfPacket)
		if !ok {
			return Digit{}, fmt.Errorf("no value for element %s", ElEndOfPacket.String())
		}

		header := headerVal.(Header)
		eop := eopVal.(int)

		digit, err := NewDigit(header, eop-1)
		if err != nil {
			return digit, fmt.Errorf("error creating digit: %w", err)
		}

		return digit, nil
	}

	Next = func(fsm *ActiveFSM[int, ZSS, ZssElement]) (ZSS, error) {
		value, ok := fsm.GetValue(ElNextState)
		if !ok {
			return ZSS{}, fmt.Errorf("no value for element %s", ElNextState.String())
		}

		return value.(ZSS), nil
	}
}

var (
	Encoder *FSM[int, ZSS, Digit, ZssElement]
)

func init() {
	builder := FsmBuilder[int, ZSS, Digit, ZssElement]{
		InitFn: Init,
		ShouldEndFn: func(fsm *ActiveFSM[int, ZSS, ZssElement]) bool {
			_, err := fsm.Tray.Read()

			return err != nil
		},
		GetResFn: EvaluateDigit,
		NextFn:   Next,
	}

	builder.AddDetFn(ElHeader, DetermineHeader)
	builder.AddDetFn(ElNextState, DetermineNextState)
	builder.AddDetFn(ElEndOfPacket, DetermineEndOfPacket)

	ec, err := builder.Build()
	if err != nil {
		panic(fmt.Errorf("error building FSM: %w", err))
	}

	Encoder = ec
}

type EncoderStream struct {
	binary []int
}

func (es *EncoderStream) ToTray() ut.Trayer[int] {
	return ut.NewSimpleTray(es.binary)
}

func NewEncoderStream(number int) (*EncoderStream, error) {
	if number < 0 {
		return nil, fmt.Errorf("invalid number: %d", number)
	}

	number += 2

	binary, err := mext.DecToBase(number, Base)
	if err != nil {
		return nil, fmt.Errorf("error converting number to binary: %w", err)
	}

	binary = binary[:len(binary)-1]

	return &EncoderStream{
		binary: binary,
	}, nil
}

const (
	Base int = 2
)

type ZssDigit int

const (
	DownDigit ZssDigit = iota
	UpDigit
	ZeroDigit
)

func (d ZssDigit) String() string {
	return [...]string{
		"u",
		"n",
		"o",
	}[d]
}

type Header struct {
	isPacket bool
	digit    ZssDigit
}

func (h *Header) String() string {
	var builder strings.Builder

	if h.isPacket {
		builder.WriteString("o")
	}

	builder.WriteString(h.digit.String())

	return builder.String()
}

func NewHeader(isPacket bool, digit ZssDigit) Header {
	return Header{
		isPacket: isPacket,
		digit:    digit,
	}
}

type Digit struct {
	header Header
	values []ZssDigit
}

func (d *Digit) String() string {
	var builder strings.Builder

	builder.WriteString(d.header.String())

	for _, value := range d.values {
		builder.WriteString(value.String())
	}

	return builder.String()
}

func NewDigit(header Header, num int) (Digit, error) {
	digit := Digit{
		header: header,
		values: make([]ZssDigit, 0),
	}

	if num < 0 {
		return digit, ue.NewErrInvalidParameter("number", ue.NewErrGTE(0))
	}

	// add 1 to the number (offset)
	num += 2

	binary, err := mext.DecToBase(num, Base)
	if err != nil {
		return digit, fmt.Errorf("error converting number to binary: %w", err)
	}

	// remove the last bit
	binary = binary[:len(binary)-1]

	for _, bit := range binary {
		if bit == 1 {
			digit.values = append(digit.values, UpDigit)
		} else {
			digit.values = append(digit.values, DownDigit)
		}
	}

	return digit, nil
}

type Number struct {
	digits []Digit
}

func (n *Number) String() string {
	var builder strings.Builder

	for _, digit := range n.digits {
		builder.WriteString(digit.String())
	}

	return builder.String()
}

func NewNumber() *Number {
	return &Number{
		digits: make([]Digit, 0),
	}
}

func (n *Number) AppendDigit(digit Digit) {
	n.digits = append(n.digits, digit)
}
