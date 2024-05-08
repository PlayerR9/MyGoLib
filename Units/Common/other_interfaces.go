package Common

// Cleaner is an interface that provides a method to remove all the elements
// from a data structure.
type Cleaner interface {
	// Clean removes all the elements from the data structure.
	Clean()
}

// Clean removes all the elements from the data structure by calling the Clean method if the
// element implements the Cleaner interface.
//
// Parameters:
//   - cleaner: The data structure to clean.
func Clean(elem any) {
	target, ok := elem.(Cleaner)
	if ok {
		target.Clean()
	}
}

// Slicer is an interface that provides a method to convert a data structure to a slice.
type Slicer[T any] interface {
	// The Slice method returns a slice containing all the elements in the data structure.
	Slice() []T
}

// SliceOf converts any type to a slice of elements of the same type.
//
// Parameters:
//   - elem: The element to convert to a slice.
//
// Returns:
//   - []T: The slice representation of the element.
//
// Behaviors:
//   - Nil elements are converted to nil slices.
//   - Slice elements are returned as is.
//   - Slicer elements have their Slice method called.
//   - Other elements are converted to slices containing a single element.
func SliceOf[T any](elem any) []T {
	if elem == nil {
		return nil
	}

	switch elem := elem.(type) {
	case []T:
		return elem
	case Slicer[T]:
		return elem.Slice()
	default:
		return []T{elem.(T)}
	}
}
