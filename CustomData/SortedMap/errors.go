package SortedMap

import (
	"fmt"

	com "github.com/PlayerR9/MyGoLib/Units/Common"
)

// ErrKeyNotFound is an error type that represents a key not found error.
type ErrKeyNotFound struct {
	// The key that was not found.
	Key string
}

// Error returns the error message: "key %q not found".
//
// Returns:
//   - string: The error message.
func (e *ErrKeyNotFound) Error() string {
	return fmt.Sprintf("key %q not found", e.Key)
}

// NewErrKeyNotFound creates a new ErrKeyNotFound with the provided key.
//
// Parameters:
//   - key: The key that was not found.
//
// Returns:
//   - *ErrKeyNotFound: A pointer to the new ErrKeyNotFound.
func NewErrKeyNotFound[K comparable](key K) *ErrKeyNotFound {
	return &ErrKeyNotFound{Key: com.StringOf(key)}
}
