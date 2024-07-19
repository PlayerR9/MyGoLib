package Go

import (
	"slices"
	"strings"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

var (
	// NonNilTypeList is a list of non-nil types.
	NonNilTypeList []string

	// NillablePrefix is a list of prefixes that indicate a type is nillable.
	NillablePrefix []string
)

func init() {
	for _, elem := range []string{
		"bool", "byte", "complex64", "complex128", "float32", "float64",
		"int", "int8", "int16", "int32", "int64", "rune", "string", "uint",
		"uint8", "uint16", "uint32", "uint64", "uintptr",
	} {
		pos, ok := slices.BinarySearch(NonNilTypeList, elem)
		uc.Assert(ok, "duplicate type in NonNilTypeList")

		NonNilTypeList = slices.Insert(NonNilTypeList, pos, elem)
	}

	NillablePrefix = []string{
		"[]",
		"map",
		"*",
		"chan",
		"func",
		"interface",
		"<-",
	}
}

// TypeNillability is an enum for the nillability of a type.
type TypeNillability int8

const (
	// IsNillable indicates that the type is nillable such as a pointer or an interface.
	IsNillable TypeNillability = iota

	// IsNotNillable indicates that the type is not nillable such as an int or a string.
	IsNotNillable

	// IsUnknown indicates that the type is unknown (i.e., it can either be nillable or not).
	IsUnknown
)

// String implements the fmt.Stringer interface.
func (t TypeNillability) String() string {
	return [...]string{
		"is nillable",
		"is not nillable",
		"is unknown",
	}[t]
}

// IsNonNilTypeID checks if the given type ID is a non-nil type.
//
// Parameters:
//   - id: The type ID to check.
//
// Returns:
//   - TypeNillability: The nillability of the type ID.
func IsNonNilTypeID(id string) TypeNillability {
	if id == "" {
		return IsUnknown
	}

	_, ok := slices.BinarySearch(NonNilTypeList, id)
	if ok {
		return IsNotNillable
	}

	if id == "error" {
		return IsNillable
	}

	ok = strings.HasPrefix(id, "struct")
	if ok {
		return IsNotNillable
	}

	for _, prefix := range NillablePrefix {
		ok := strings.HasPrefix(id, prefix)
		if ok {
			return IsNillable
		}
	}

	return IsUnknown
}
