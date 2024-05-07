package Interfaces

import (
	"fmt"
)

// StringOf converts any type to a string.
//
// The function converts the input element to a string using the following rules:
//   - If the element is a string, it is returned as is.
//   - If the element implements the fmt.Stringer interface, the String method is called.
//   - If the element is an error, the Error method is called.
//   - If the element is a byte slice, it is converted to a string.
//   - If the element is a rune slice, it is converted to a string.
//   - If the element implements the ffs.FStringer interface, the FString function is called.
//   - Otherwise, the element is converted to a string using the fmt.Sprintf function.
//
// Parameters:
//   - elem: The element to convert to a string.
//
// Returns:
//   - string: The string representation of the element.
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
