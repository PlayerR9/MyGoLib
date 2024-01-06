package StrExt

type ErrSuffixTooLong struct{}

func (e ErrSuffixTooLong) Error() string {
	return "suffix is too long"
}

type ErrOpeningTokenEmpty struct{}

func (e ErrOpeningTokenEmpty) Error() string {
	return "opening token is empty"
}

type ErrClosingTokenEmpty struct{}

func (e ErrClosingTokenEmpty) Error() string {
	return "closing token is empty"
}

type ErrOpeningTokenNotFound struct{}

func (e ErrOpeningTokenNotFound) Error() string {
	return "opening token not found in content"
}
