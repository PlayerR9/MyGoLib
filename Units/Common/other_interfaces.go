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

// Comparer is an interface that defines a method to compare two objects
// of the same type.
type Comparer[T any] interface {
	// Compare returns a negative value if the object is less than the other object,
	// zero if they are equal, and a positive value if the object is greater
	// than the other object.
	//
	// Parameters:
	//   - other: The other object to compare to.
	//
	// Returns:
	//   - int: The result of the comparison.
	Compare(other T) int
}

// CompareOf compares two objects of the same type. If any of the objects implements
// the Comparer interface, the Compare method is called. Otherwise, the objects are
// compared using the < and == operators.
//
// Parameters:
//   - obj1: The first object to compare.
//   - obj2: The second object to compare.
//
// Returns:
//   - int: The result of the comparison.
func CompareOf(obj1, obj2 any) (int, bool) {
	if obj1 == nil || obj2 == nil {
		return 0, false
	}

	switch obj1 := obj1.(type) {
	case Comparer[any]:
		val2, ok := obj2.(Comparer[any])
		if !ok {
			return 0, false
		}

		return obj1.Compare(val2), true
	default:
		return CompareAny(obj1, obj2)
	}
}
