package Common

import "fmt"

// Comparable is an interface that defines the behavior of a type that can be
// compared with other values of the same type using the < and > operators.
// The interface is implemented by the built-in types int, int8, int16, int32,
// int64, uint, uint8, uint16, uint32, uint64, float32, float64, and string.
type Comparable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 | ~string
}

// StringOf converts any type to a string.
//
// Parameters:
//   - elem: The element to convert to a string.
//
// Returns:
//   - string: The string representation of the element.
//
// Behaviors:
//   - String elements are returned as is.
//   - fmt.Stringer elements have their String method called.
//   - error elements have their Error method called.
//   - []byte and []rune elements are converted to strings.
//   - Other elements are converted to strings using fmt.Sprintf and the %v format.
func StringOf(elem any) string {
	if elem == nil {
		return ""
	}

	switch elem := elem.(type) {
	case string:
		return elem
	case fmt.Stringer:
		return elem.String()
	case error:
		return elem.Error()
	case []byte:
		return string(elem)
	case []rune:
		return string(elem)
	default:
		return fmt.Sprintf("%v", elem)
	}
}

// Equaler is an interface that defines a method to compare two objects
// of the same type.
type Equaler[T any] interface {
	// Equals returns true if the object is equal to the other object.
	//
	// Parameters:
	// 	- other: The other object to compare to.
	//
	// Returns:
	// 	- bool: True if the object is equal to the other object.
	Equals(other T) bool
}

// EqualOf compares two objects of the same type. If any of the objects implements
// the Equaler interface, the Equals method is called. Otherwise, the objects are
// compared using the == operator. However, obj1 is always checked first.
//
// Parameters:
//   - obj1: The first object to compare.
//   - obj2: The second object to compare.
//
// Returns:
//   - bool: True if the objects are equal, false otherwise.
//
// Behaviors:
//   - Nil objects are always considered different.
func EqualOf(obj1, obj2 any) bool {
	if obj1 == nil || obj2 == nil {
		return false
	}

	a, ok := obj1.(Equaler[any])
	if ok {
		return a.Equals(obj2)
	}

	b, ok := obj2.(Equaler[any])
	if ok {
		return b.Equals(obj1)
	}

	return obj1 == obj2
}

// Copier is an interface that provides a method to create a deep copy of an object.
type Copier interface {
	// Copy creates a shallow copy of the object.
	//
	// Returns:
	//   - Copier: A shallow copy or a deep copy of the object.
	Copy() Copier
}

// CopyOf creates a copy of the element by either calling the Copy method if the
// element implements the Copier interface or returning the element as is.
//
// Parameters:
//   - elem: The element to copy.
//
// Returns:
//   - any: A copy of the element.
func CopyOf(elem any) any {
	if elem == nil {
		return nil
	}

	switch elem := elem.(type) {
	case Copier:
		return elem.Copy()
	default:
		return elem
	}
}
