package OrderedMap

// ModifyValueFunc is a function that modifies a value.
//
// Parameters:
//   - V: The value to modify.
//
// Returns:
//   - V: The modified value.
type ModifyValueFunc[V any] func(V) (V, error)
