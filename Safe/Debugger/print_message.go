package Debugger

type PMType int8

const (
	PM_Println PMType = iota
	PM_Printf
	PM_Print
)

type PrintMessage struct {
	pmType PMType
	format string
	v      []any
}

func NewPrintMessage(pmType PMType, format string, v ...any) *PrintMessage {
	return &PrintMessage{
		pmType: pmType,
		format: format,
		v:      v,
	}
}
