package StrExt

type ErrSuffixTooLong struct{}

func (e ErrSuffixTooLong) Error() string {
	return "suffix is too long"
}
