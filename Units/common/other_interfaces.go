package common

import (
	"fmt"
)

// Enumer is an interface representing an enumeration.
type Enumer interface {
	~int

	fmt.Stringer
}

// Copier is an interface that provides a method to create a copy of an element.
type Copier interface {
	// Copy creates a copy of the element.
	//
	// Returns:
	//   - Copier: The copy of the element.
	Copy() Copier
}

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

// Runer is an interface that provides a method to get the runes of a string.
type Runer interface {
	// Runes returns the runes of the string.
	//
	// Returns:
	//   - []rune: The runes of the string.
	Runes() []rune
}

// RunesOf returns the runes of a string. If the string implements the Runer interface,
// the Runes method is called. Otherwise, the string is converted to a slice of runes.
//
// Parameters:
//   - str: The string to get the runes from.
//
// Returns:
//   - []rune: The runes of the string.
func RunesOf(str any) []rune {
	if str == nil {
		return nil
	}

	switch str := str.(type) {
	case Runer:
		return str.Runes()
	case []rune:
		return str
	case string:
		return []rune(str)
	case []byte:
		return []rune(string(str))
	case error:
		return []rune(str.Error())
	case fmt.Stringer:
		return []rune(str.String())
	default:
		return []rune(fmt.Sprintf("%v", str))
	}
}

// Hashable is an interface that provides a method to get the hash code of an element.
type Hashable interface {
	// Hash returns the hash code of the element.
	//
	// Returns:
	//   - int: The hash code of the element.
	Hash() int
}

// HashCode returns the hash code of an element. If the element implements the Hashable interface,
// the Hash method is called. Otherwise, the hash code is calculated using the default hash function.
//
// Parameters:
//   - elem: The element to get the hash code from.
//
// Returns:
//   - int: The hash code of the element.
func HashCode(elem any) int {
	if elem == nil {
		return 0
	}

	switch elem := elem.(type) {
	case Hashable:
		return elem.Hash()
	case int:
		return elem
	case int8:
		return int(elem)
	case int16:
		return int(elem)
	case int32:
		return int(elem)
	case int64:
		return int(elem)
	case uint:
		return int(elem)
	case uint8:
		return int(elem)
	case uint16:
		return int(elem)
	case uint32:
		return int(elem)
	case uint64:
		return int(elem)
	case float32:
		return int(elem)
	case float64:
		return int(elem)
	case bool:
		if elem {
			return 1
		} else {
			return 0
		}
	case string:
		return hashString(elem)
	case []byte:
		return hashBytes(elem)
	case error:
		return hashString(elem.Error())
	case fmt.Stringer:
		return hashString(elem.String())
	default:
		return hashString(fmt.Sprintf("%v", elem))
	}
}

// hashString is a helper function that calculates the hash code of a string.
//
// Parameters:
//   - str: The string to calculate the hash code from.
//
// Returns:
//   - int: The hash code of the string.
func hashString(str string) int {
	hash := 0

	for _, r := range str {
		hash = 31*hash + int(r)
	}

	return hash
}

// hashBytes is a helper function that calculates the hash code of a byte slice.
//
// Parameters:
//   - bytes: The byte slice to calculate the hash code from.
//
// Returns:
//   - int: The hash code of the byte slice.
func hashBytes(bytes []byte) int {
	hash := 0

	for _, b := range bytes {
		hash = 31*hash + int(b)
	}

	return hash
}
