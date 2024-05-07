package FString

const (
	// DefaultIndentation is the default indentation string.
	DefaultIndentation string = "   "

	// DefaultSeparator is the default separator string.
	DefaultSeparator string = ", "
)

// Options is a type that represents the options that can be passed to the builder.
type Options []BuildOption

var (
	// ArrayDefault is the default options for an array.
	// [1, 2, 3]
	ArrayDefault Options = Options{
		WithIndentation(NewIndentConfig(DefaultIndentation, 0, false, true)),
		WithDelimiterLeft(NewDelimiterConfig("[", false)),
		WithDelimiterRight(NewDelimiterConfig("]", false)),
		WithSeparator(NewSeparator(DefaultSeparator, false)),
	}
)
