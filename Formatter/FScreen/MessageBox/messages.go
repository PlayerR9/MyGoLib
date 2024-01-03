package MessageBox

import (
	"slices"
	"strconv"
	"strings"
)

const (
	Space     = ' '
	Separator = '='
)

type TextMessageType int

const (
	NormalText TextMessageType = iota
	DebugText
	SuccessText
	WarningText
	ErrorText
	FatalText
	SeparatorLine
	BreakLine
)

func (enum TextMessageType) String() string {
	return [...]string{
		"Normal Text",
		"[DEBUG]:",
		"[SUCCESS]:",
		"[WARNING]:",
		"[ERROR]:",
		"[FATAL]:",
		"Separator Line",
		"Break Line",
	}[enum]
}

func (enum TextMessageType) IsInfo() bool {
	return slices.Contains([]TextMessageType{
		DebugText,
		SuccessText,
		WarningText,
		ErrorText,
		FatalText,
	}, enum)
}

type TextMessage struct {
	contents    []string
	messageType TextMessageType
}

func (tm *TextMessage) IsEmpty() bool {
	if len(tm.contents) >= 2 {
		return false
	}

	return len(tm.contents) == 1 && tm.messageType.IsInfo()
}

func (tm *TextMessage) GetType() TextMessageType {
	return tm.messageType
}

func (tm *TextMessage) GetPrefix() string {
	if tm.messageType.IsInfo() {
		return tm.messageType.String()
	}

	return ""
}

func (tm *TextMessage) GetText() string {
	switch tm.messageType {
	case BreakLine:
		return string(Space)
	case SeparatorLine:
		return string(Separator)
	default:
		contents := make([]string, len(tm.contents))
		copy(contents, tm.contents)

		if tm.messageType.IsInfo() {
			contents = contents[1:]
		}

		return strings.Join(contents, " ")
	}
}

func (tm *TextMessage) GetContents() []string {
	return tm.contents
}

func NewTextMessage(messageType TextMessageType, contents ...string) TextMessage {
	if messageType == SeparatorLine || messageType == BreakLine {
		return TextMessage{
			contents:    nil,
			messageType: messageType,
		}
	}

	var toSend []string

	if messageType.IsInfo() {
		toSend = append(toSend, messageType.String())
	}

	for _, content := range contents {
		content = strings.TrimSpace(content)
		if content != "" {
			tmp := strconv.Quote(content)
			toSend = append(toSend, tmp[1:len(tmp)-1])
		}
	}

	return TextMessage{
		contents:    toSend,
		messageType: messageType,
	}
}
