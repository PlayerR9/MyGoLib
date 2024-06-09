package common

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

// Equaler is an interface that defines a method to compare two objects of the
// same type for equality.
type Equaler interface {
	// Equals returns true if the object is equal to the other object.
	//
	// Parameters:
	//   - other: The other object to compare to.
	//
	// Returns:
	//   - bool: True if the objects are equal, false otherwise.
	Equals(other Equaler) bool
}
