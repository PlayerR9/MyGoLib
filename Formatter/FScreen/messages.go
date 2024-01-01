package FScreen

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/gdamore/tcell"
)

func ShortFormOf(line string) string {
	lineLen := len(line)

	// If the line fits on the screen without modification
	if lineLen <= innerScreenWidth {
		return line
	}

	// Otherwise, reduce it
	fields := strings.Fields(line)
	fieldsLen := len(fields)

	// If the line is a single word
	if fieldsLen == 1 {
		return line[:innerScreenWidth-HellipLen] + Hellip
	}

	var builder strings.Builder

	for i, field := range fields {
		fieldLen := len(field)
		builderLen := builder.Len()

		// If adding this field would exceed the screen width
		if builderLen+fieldLen+HellipLen+1 > innerScreenWidth {
			break
		}

		builder.WriteString(field)

		// Only add a space if this is not the last word that fits
		if i < fieldsLen-1 && builderLen+len(fields[i+1])+HellipLen+1 <= innerScreenWidth {
			builder.WriteString(" ")
		}
	}

	builder.WriteString(Hellip)

	return builder.String()
}

func validateContents(prefix, content string, optionals []string) ([]string, bool) {
	var toSend []string

	if prefix != "" {
		toSend = append(toSend, prefix)
	}

	if content != "" {
		x := strconv.Quote(content)
		toSend = append(toSend, x[1:len(x)-1])
	}

	for _, optional := range optionals {
		if optional != "" {
			x := strconv.Quote(optional)
			toSend = append(toSend, x[1:len(x)-1])
		}
	}

	return toSend, (prefix == "" && len(toSend) == 0) || (prefix != "" && len(toSend) == 1)
}

type Messager interface {
	Apply()
	GetStyle() tcell.Style
}

type ClearScreenMSG struct {
}

func (msg ClearScreenMSG) Apply() {
	switch shouldPrint {
	case YetToPrint, NotPrinted:
		printScreen()
	}

	clearScreen()

	shouldPrint = MustPrint
}

func (msg ClearScreenMSG) GetStyle() tcell.Style {
	return NormalStyle
}

func MessageClearScreen() ClearScreenMSG {
	return ClearScreenMSG{}
}

type CloseMSG struct {
}

func (msg CloseMSG) Apply() {
	close(messageChannel)

	switch shouldPrint {
	case YetToPrint, NotPrinted:
		printScreen()
		shouldPrint = JustPrinted
	default:
		shouldPrint = NotPrinted
	}
}

func (msg CloseMSG) GetStyle() tcell.Style {
	return NormalStyle
}

func MessageClose() CloseMSG {
	return CloseMSG{}
}

type ModifyCurrentProcessMSG struct {
	title string
}

func (msg ModifyCurrentProcessMSG) Apply() {
	currentProcess = msg.title

	shouldPrint = YetToPrint
}

func (msg ModifyCurrentProcessMSG) GetStyle() tcell.Style {
	return NormalStyle
}

func MessageSetCurrentProcess(title string) Messager {
	if title == "" {
		return MessageFatal("title must not be empty")
	}

	return ModifyCurrentProcessMSG{
		title: title,
	}
}

func MessageDesetCurrentProcess() ModifyCurrentProcessMSG {
	return ModifyCurrentProcessMSG{}
}

type SetCounterMSG struct {
	counter  ScreenCounter
	ifExists ImportantMSG
}

func (msg SetCounterMSG) Apply() {
	index := slices.IndexFunc(counters, func(c ScreenCounter) bool {
		return msg.counter.Equal(c)
	})

	if index != -1 {
		msg.ifExists.Apply()
	} else {
		counters[index] = msg.counter

		shouldPrint = YetToPrint
	}
}

func (msg SetCounterMSG) GetStyle() tcell.Style {
	return NormalStyle
}

func MessageSetCounter(label string, totalElements int) Messager {
	if strings.TrimSpace(label) == "" {
		return MessageError("trying to set counter with empty label")
	} else if totalElements < 0 {
		return MessageError(
			fmt.Sprintf("cannot set counter \"%s\":", label),
			"Negative 'totalElements'",
			fmt.Sprintf("(got %d instead)", totalElements),
		)
	}

	return SetCounterMSG{
		counter: MakeCounter(label, totalElements),
		ifExists: MessageWarning(
			fmt.Sprintf("counter \"%s\" already exists", label),
			"skipping set",
		).(ImportantMSG),
	}
}

type DesetCounterMSG struct {
	label       string
	ifNotExists ImportantMSG
}

func (msg DesetCounterMSG) Apply() {
	index := slices.IndexFunc(counters, func(c ScreenCounter) bool {
		return c.HasLabel(msg.label)
	})

	if index != -1 {
		counters = append(counters[:index], counters[index+1:]...)

		shouldPrint = YetToPrint
	} else {
		msg.ifNotExists.Apply()
	}
}

func (msg DesetCounterMSG) GetStyle() tcell.Style {
	return NormalStyle
}

func MessageDesetCounter(label string) Messager {
	if strings.TrimSpace(label) == "" {
		return MessageWarning(
			"trying to deset counter with empty label",
			"skipping deset",
		).(ImportantMSG)
	}

	return DesetCounterMSG{
		label: label,
		ifNotExists: MessageWarning(
			fmt.Sprintf("counter \"%s\" does not exist", label),
			"skipping deset",
		).(ImportantMSG),
	}
}

type IncrementMSG struct {
	label       string
	ifNotExists ImportantMSG
}

func (msg IncrementMSG) Apply() {
	index := slices.IndexFunc(counters, func(c ScreenCounter) bool {
		return c.HasLabel(msg.label)
	})

	if index != -1 {
		x := counters[index]
		x.Increment()
		counters[index] = x

		shouldPrint = YetToPrint
	} else {
		msg.ifNotExists.Apply()
	}
}

func (msg IncrementMSG) GetStyle() tcell.Style {
	return NormalStyle
}

func MessageIncrement(label string) IncrementMSG {
	return IncrementMSG{
		label: label,
		ifNotExists: MessageWarning(
			fmt.Sprintf("counter \"%s\" does not exist", label),
			"skipping increment",
		).(ImportantMSG),
	}
}

type ReduceMSG struct {
	label       string
	ifNotExists ImportantMSG
}

func (msg ReduceMSG) Apply() {
	index := slices.IndexFunc(counters, func(c ScreenCounter) bool {
		return c.HasLabel(msg.label)
	})

	if index != -1 {
		x := counters[index]
		x.Reduce()
		counters[index] = x

		shouldPrint = YetToPrint
	} else {
		msg.ifNotExists.Apply()
	}
}

func (msg ReduceMSG) GetStyle() tcell.Style {
	return NormalStyle
}

func MessageReduce(label string) ReduceMSG {
	return ReduceMSG{
		label: label,
		ifNotExists: MessageWarning(
			fmt.Sprintf("counter \"%s\" does not exist", label),
			"skipping reduce",
		).(ImportantMSG),
	}
}

type EmptyMSG struct{}

func (msg EmptyMSG) Apply() {

}

func (msg EmptyMSG) GetStyle() tcell.Style {
	return NormalStyle
}

type LineMSG struct {
	contents []string
	stlye    tcell.Style
}

func (msg LineMSG) Apply() {
	writeLines(msg.contents)
}

func (msg LineMSG) GetStyle() tcell.Style {
	return msg.stlye
}

func MessageDebug(content string, optionals ...string) Messager {
	toSend, isEmpty := validateContents("[DEBUG]:", content, optionals)
	if isEmpty {
		return EmptyMSG{}
	}

	return LineMSG{
		contents: toSend,
		stlye:    DebugStyle,
	}
}

func MessagePrint(content string, optionals ...string) Messager {
	toSend, isEmpty := validateContents("", content, optionals)
	if isEmpty {
		return EmptyMSG{}
	}

	return LineMSG{
		contents: toSend,
		stlye:    NormalStyle,
	}
}

func MessageCustomLine(label, content string, optionals ...string) Messager {
	if label == "" {
		return MessageError(
			"trying to create custom line with empty label",
		)
	}

	toSend, isEmpty := validateContents(label, content, optionals)
	if isEmpty {
		return EmptyMSG{}
	}

	return LineMSG{
		contents: toSend,
		stlye:    NormalStyle,
	}
}

func MessageSeparator(count int) Messager {
	if count < 0 {
		return MessageError(
			"trying to create separator with negative count",
		)
	}

	msg := LineMSG{
		contents: make([]string, 0, count),
		stlye:    NormalStyle,
	}

	for i := 0; i < count; i++ {
		msg.contents = append(msg.contents, separator)
	}

	return msg
}

func MessageLineBreak(count int) Messager {
	if count < 0 {
		return MessageError(
			"trying to create line break with negative count",
		)
	}

	msg := LineMSG{
		contents: make([]string, 0, count),
		stlye:    NormalStyle,
	}

	for i := 0; i < count; i++ {
		msg.contents = append(msg.contents, lineBreak)
	}

	return msg
}

func MessageSuccess(content string, optionals ...string) Messager {
	toSend, isEmpty := validateContents("SUCCESS:", content, optionals)
	if isEmpty {
		return EmptyMSG{}
	}

	return ImportantMSG{
		contents: toSend,
		style:    SuccessStyle,
	}
}

type ImportantMSG struct {
	contents []string
	style    tcell.Style
}

func (msg ImportantMSG) Apply() {
	writeLines(msg.contents)
}

func (msg ImportantMSG) GetStyle() tcell.Style {
	return msg.style
}

func MessageWarning(content string, optionals ...string) Messager {
	toSend, isEmpty := validateContents("WARNING:", content, optionals)
	if isEmpty {
		return EmptyMSG{}
	}

	return ImportantMSG{
		contents: toSend,
		style:    WarningStyle,
	}
}

func MessageError(content string, optionals ...string) Messager {
	toSend, isEmpty := validateContents("ERROR:", content, optionals)
	if isEmpty {
		return EmptyMSG{}
	}

	return ImportantMSG{
		contents: toSend,
		style:    ErrorStyle,
	}
}

type FatalMSG struct {
	contents []string
}

func (msg FatalMSG) Apply() {
	writeLines(msg.contents)
}

func (msg FatalMSG) GetStyle() tcell.Style {
	return FatalStyle
}

func MessageFatal(content string, optionals ...string) Messager {
	toSend, isEmpty := validateContents("FATAL:", content, optionals)
	if isEmpty {
		return EmptyMSG{}
	}

	return FatalMSG{
		contents: toSend,
	}
}

type ForcePrintMSG struct {
}

func (msg ForcePrintMSG) Apply() {
	shouldPrint = MustPrint
}

func (msg ForcePrintMSG) GetStyle() tcell.Style {
	return NormalStyle
}

func MessageForcePrint() ForcePrintMSG {
	return ForcePrintMSG{}
}
