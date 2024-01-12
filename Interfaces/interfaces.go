package Interfaces

type Cleanable interface {
	// Cleanup is called when the object is no longer needed.
	// This is useful for cleaning up resources such as file handles.
	Cleanup()
}

// Cleanup calls the Cleanup method of the given object if it implements the Cleanable interface.
// It returns a zero value of the given type.
func Cleanup[T any](value any) T {
	if val, ok := value.(Cleanable); ok {
		val.Cleanup()
	}

	return *new(T)
}
