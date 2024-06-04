package Runner

import (
	ers "github.com/PlayerR9/MyGoLib/Units/errors"
)

// NewRedirector creates a new redirector that redirects messages
// from one channel to another.
//
// Parameters:
//   - from: The channel to redirect messages from.
//   - to: The channel to redirect messages to.
//
// Returns:
//   - ExecuteMsgFunc: The function that redirects messages.
//   - error: An error of type *ers.ErrInvalidParameter if from or to are nil.
func NewRedirector[T any](from <-chan T, to chan<- T) (ExecuteMsgFunc[T], error) {
	if from == nil {
		return nil, ers.NewErrNilParameter("from channel")
	} else if to == nil {
		return nil, ers.NewErrNilParameter("to channel")
	}

	return func(msg T) error {
		to <- msg

		return nil
	}, nil
}
