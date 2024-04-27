package StringExt

// TokenType is an enum that represents the type of token in a string.
type TokenType int8

const (
	// OpToken represents an opening token.
	OpToken TokenType = iota

	// ClToken represents a closing token.
	ClToken
)

// String is a method of fmt.Stringer interface that returns
// the string representation of the token type.
//
// Returns:
// 	- string: the string representation of the token type.
func (t TokenType) String() string {
	return [...]string{"opening", "closing"}[t]
}
