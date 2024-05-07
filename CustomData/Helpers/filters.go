package Helpers

// FilterNilHResult filters the HResult slice to remove nil values.
//
// Parameters:
//   - hr: The HResult slice to filter.
//
// Returns:
//   - bool: True if the HResult is not nil, false otherwise.
func FilterNilHResult[T any](hr HResult[T]) bool {
	return hr.Second != nil
}
