package Common

import (
	"fmt"
	"strconv"
)

// Objecter is an interface that defines the behavior of an object that can be
// copied, compared, and converted to a string.
type Objecter interface {
	fmt.Stringer
	Copier
	Equaler
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
	case Objecter:
		return elem.Copy()
	default:
		x := elem

		return x
	}
}

// EqualOf compares two objects of the same type. If any of the objects implements
// the Equaler interface, the Equals method is called. Otherwise, the objects are
// compared using the == operator. However, a is always checked first.
//
// Parameters:
//   - a: The first object to compare.
//   - b: The second object to compare.
//
// Returns:
//   - bool: True if the objects are equal, false otherwise.
//
// Behaviors:
//   - Nil objects are always considered different.
func EqualOf(a, b any) bool {
	if a == nil || b == nil {
		return false
	}

	switch a := a.(type) {
	case Comparer:
		otherB, ok := b.(Comparer)
		if !ok {
			return false
		}

		return a.Compare(otherB) == 0
	case Objecter:
		otherB, ok := b.(Objecter)
		if !ok {
			return false
		}

		return a.Equals(otherB)
	default:
		return a == b
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
	case int:
		return strconv.FormatInt(int64(elem), 10)
	case int8:
		return strconv.FormatInt(int64(elem), 10)
	case int16:
		return strconv.FormatInt(int64(elem), 10)
	case int32:
		return strconv.FormatInt(int64(elem), 10)
	case int64:
		return strconv.FormatInt(elem, 10)
	case uint:
		return strconv.FormatUint(uint64(elem), 10)
	case uint8:
		return strconv.FormatUint(uint64(elem), 10)
	case uint16:
		return strconv.FormatUint(uint64(elem), 10)
	case uint32:
		return strconv.FormatUint(uint64(elem), 10)
	case uint64:
		return strconv.FormatUint(elem, 10)
	case float32:
		return strconv.FormatFloat(float64(elem), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(elem, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(elem)
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
