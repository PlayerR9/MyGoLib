package StrExt

// ErrSuffixTooLong is a struct that represents an error when a suffix
// is too long.
// It does not have any fields as the error condition is solely based
// on the length of the suffix.
type ErrSuffixTooLong struct{}

// Error is a method of the ErrSuffixTooLong type that implements the
// error interface. It returns a string representation of the error,
// that is, the string "suffix is too long".
func (e ErrSuffixTooLong) Error() string {
	return "suffix is too long"
}

// ErrOpeningTokenEmpty is a struct that represents an error when an
// opening token is empty.
// It does not have any fields as the error condition is solely based
// on the absence of an opening token.
type ErrOpeningTokenEmpty struct{}

// Error is a method of the ErrOpeningTokenEmpty type that implements
// the error interface. It returns a string representation of the error,
// that is, the string "opening token is empty".
func (e ErrOpeningTokenEmpty) Error() string {
	return "opening token is empty"
}

// ErrClosingTokenEmpty is a struct that represents an error when a
// closing token is empty.
// It does not have any fields as the error condition is solely based
// on the absence of a closing token.
type ErrClosingTokenEmpty struct{}

// Error is a method of the ErrClosingTokenEmpty type that implements
// the error interface. It returns a string representation of the error,
// that is, the string "closing token is empty".
func (e ErrClosingTokenEmpty) Error() string {
	return "closing token is empty"
}

// ErrOpeningTokenNotFound is a struct that represents an error when an
// opening token is not found in the content.
// It does not have any fields as the error condition is solely based
// on the absence of an opening token in the content.
type ErrOpeningTokenNotFound struct{}

// Error is a method of the ErrOpeningTokenNotFound type that implements
// the error interface. It returns a string representation of the error,
// that is, the string "opening token not found in content".
func (e ErrOpeningTokenNotFound) Error() string {
	return "opening token not found in content"
}
