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

// Compare compares two values of the same type that implement the Comparable
// interface. If the values are equal, the function returns 0. If the first
// value is less than the second value, the function returns -1. If the first
// value is greater than the second value, the function returns 1.
//
// Parameters:
//   - a: The first value to compare.
//   - b: The second value to compare.
//
// Returns:
//   - int: -1 if a < b, 0 if a == b, 1 if a > b.
//   - bool: True if the values are comparable.
//
// Behaviors:
//   - If the values are not comparable, the function returns false.
func Compare[T Comparable](a, b T) (int, bool) {
	if a < b {
		return -1, true
	} else if a > b {
		return 1, true
	}

	return 0, true
}

// Compare compares two values of the same type that implement the Comparable
// interface. If the values are equal, the function returns 0. If the first
// value is less than the second value, the function returns -1. If the first
// value is greater than the second value, the function returns 1.
//
// Parameters:
//   - a: The first value to compare.
//   - b: The second value to compare.
//
// Returns:
//   - int: -1 if a < b, 0 if a == b, 1 if a > b.
//   - bool: True if the values are comparable.
//
// Behaviors:
//   - If the values are not comparable, the function returns false.
func CompareAny(a, b any) (int, bool) {
	if a == nil || b == nil {
		return 0, false
	}

	switch a := a.(type) {
	case int:
		valB, ok := b.(int)
		if !ok {
			return 0, false
		}

		return a - valB, true
	case int8:
		valB, ok := b.(int8)
		if !ok {
			return 0, false
		}

		return int(a - valB), true
	case int16:
		valB, ok := b.(int16)
		if !ok {
			return 0, false
		}

		return int(a - valB), true
	case int32:
		valB, ok := b.(int32)
		if !ok {
			return 0, false
		}

		return int(a - valB), true
	case int64:
		valB, ok := b.(int64)
		if !ok {
			return 0, false
		}

		return int(a - valB), true
	case uint:
		valB, ok := b.(uint)
		if !ok {
			return 0, false
		}

		return int(a - valB), true
	case uint8:
		valB, ok := b.(uint8)
		if !ok {
			return 0, false
		}

		return int(a - valB), true
	case uint16:
		valB, ok := b.(uint16)
		if !ok {
			return 0, false
		}

		return int(a - valB), true
	case uint32:
		valB, ok := b.(uint32)
		if !ok {
			return 0, false
		}

		return int(a - valB), true
	case uint64:
		valB, ok := b.(uint64)
		if !ok {
			return 0, false
		}

		return int(a - valB), true
	case float32:
		valB, ok := b.(float32)
		if !ok {
			return 0, false
		}

		return int(a - valB), true
	case float64:
		valB, ok := b.(float64)
		if !ok {
			return 0, false
		}

		return int(a - valB), true
	case string:
		valB, ok := b.(string)
		if !ok {
			return 0, false
		}

		if a < valB {
			return -1, true
		} else if a > valB {
			return 1, true
		} else {
			return 0, true
		}
	default:
		return 0, false
	}
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
