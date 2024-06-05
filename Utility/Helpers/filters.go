package Helpers

// FilterIsSuccess filters any helper that is not successful.
//
// Parameters:
//   - h: The helper to filter.
//
// Returns:
//   - bool: True if the helper is successful, false otherwise.
//
// Behaviors:
//   - It assumes that the h is not nil.
func FilterIsSuccess[T Helperer[O], O any](h T) bool {
	return h.GetData().Second == nil
}
