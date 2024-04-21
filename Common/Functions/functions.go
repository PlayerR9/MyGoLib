package Functions

// ObserverFunc is a function type that takes a value of type T and returns an error.
// It is used to observe the values of the nodes during a traversal.
type ObserverFunc[T any] func(T) error
