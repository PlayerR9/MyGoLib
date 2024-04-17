package FString

import "strings"

const (
	DefaultIndentation string = "   "
	DefaultSeparator   string = ", "
)

type Options []BuildOption

var (
	ArrayDefault Options = Options{
		WithIndentation(NewIndentConfig(DefaultIndentation, 0, false, true)),
		WithDelimiterLeft(NewDelimiterConfig("[", false)),
		WithDelimiterRight(NewDelimiterConfig("]", false)),
		WithSeparator(NewSeparator(DefaultSeparator, false)),
	}
)

type FStringer interface {
	FString(int) []string
}

func FString(obj FStringer) string {
	return strings.Join(obj.FString(0), "\n")
}
