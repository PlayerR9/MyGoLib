package Stack

type ErrEmptyStack struct{}

func (e ErrEmptyStack) Error() string {
	return "Empty stack"
}

type ErrFullStack struct{}

func (e ErrFullStack) Error() string {
	return "Full stack"
}
