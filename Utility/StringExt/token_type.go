package StringExt

type TokenType int8

const (
	OpToken TokenType = iota
	ClToken
)

func (t TokenType) String() string {
	return [...]string{"opening", "closing"}[t]
}
