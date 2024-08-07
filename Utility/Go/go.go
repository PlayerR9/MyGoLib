package Go

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	uc "github.com/PlayerR9/lib_units/common"
	luint "github.com/PlayerR9/lib_units/ints"
	utch "github.com/PlayerR9/lib_units/runes"
)

var (
	// GoReservedKeywords is a list of Go reserved keywords.
	GoReservedKeywords []string
)

func init() {
	GoReservedKeywords = []string{
		"break", "case", "chan", "const", "continue", "default", "defer", "else",
		"fallthrough", "for", "func", "go", "goto", "if", "import", "interface",
		"map", "package", "range", "return", "select", "struct", "switch", "type",
		"var",
	}
}

// check_letter checks if the given rune is a letter and if it is upper case.
//
// Parameters:
//   - c: The rune to check.
//
// Returns:
//   - bool: True if the rune is a letter. False otherwise.
//   - bool: True if the letter is upper case. False otherwise.
func check_letter(c rune) (bool, bool) {
	ok := unicode.IsLetter(c)
	if !ok {
		return false, false
	}

	ok = unicode.IsUpper(c)
	return ok, true
}

// IsValidName checks if the given variable name is valid.
//
// This function checks if the variable name is not empty and if it is not a
// Go reserved keyword. It also checks if the variable name is not in the list
// of keywords.
//
// Parameters:
//   - variable_name: The variable name to check.
//   - keywords: The list of keywords to check against.
//
// Returns:
//   - error: An error if the variable name is invalid.
func IsValidName(variable_name string, keywords []string) error {
	if variable_name == "" {
		err := uc.NewErrEmpty(variable_name)
		return err
	}

	ok := slices.Contains(GoReservedKeywords, variable_name)
	if ok {
		err := errors.New("name is a reserved keyword")
		return err
	}

	ok = slices.Contains(keywords, variable_name)
	if ok {
		err := errors.New("name is not allowed")
		return err
	}

	return nil
}

// MakeVariableName converts a type name to a variable name.
//
// Parameters:
//   - type_name: The type name to convert.
//
// Returns:
//   - string: The variable name.
//   - error: An error if the type name is invalid.
//
// Errors:
//   - *common.ErrInvalidParameter: If the type name is empty or not valid UTF-8.
//   - *common.ErrAt: If the type name is invalid at a specific position.
func MakeVariableName(type_name string) (string, error) {
	if type_name == "" {
		return "", uc.NewErrInvalidParameter("type_name", uc.NewErrEmpty(type_name))
	}

	chars, err := utch.StringToUtf8(type_name)
	if err != nil {
		return "", uc.NewErrInvalidParameter("type_name", err)
	}

	var builder strings.Builder

	is_upper, ok := check_letter(chars[0])
	if !ok || !is_upper {
		reason := errors.New("not an upper case letter")
		err := luint.NewErrAt(0, "character", reason)
		return "", err
	}

	c := unicode.ToLower(chars[0])
	builder.WriteRune(c)

	for i, c := range chars[1:] {
		if c == '_' {
			continue
		}

		ok := unicode.IsNumber(c)
		if ok {
			continue
		}

		is_upper, ok := check_letter(c)
		if !ok {
			err := luint.NewErrAt(i+1, "character", errors.New("neither a letter nor a number"))
			return "", err
		}

		if is_upper {
			c = unicode.ToLower(c)
			builder.WriteRune(c)
		}
	}

	str := builder.String()
	return str, nil
}

// fix_variable_name is an helper function that fixes the given variable name by checking whether
// or not the name trimmed to the minimum length is valid. If it is not valid, the function tries
// the remaining characters one by one until a valid name is found.
//
// Parameters:
//   - var_name: The variable name to fix.
//   - keywords: The list of keywords to check against.
//   - min: The minimum length of the variable name.
//
// Returns:
//   - string: The fixed variable name.
//   - error: An error if the variable name is invalid.
func fix_variable_name(var_name string, keywords []string, min int) (string, error) {
	size := utf8.RuneCountInString(var_name)
	if size <= min {
		err := IsValidName(var_name, keywords)
		if err != nil {
			return "", err
		}

		return var_name, nil
	}

	chars, _ := utch.StringToUtf8(var_name)
	// uc.AssertErr(err, "StringToUtf8(%q)", var_name)

	var builder strings.Builder

	for i := 0; i < min; i++ {
		builder.WriteRune(chars[i])
	}

	str := builder.String()

	err := IsValidName(str, keywords)
	if err == nil {
		return str, nil
	}

	for i := min; i < size; i++ {
		builder.WriteRune(chars[i])
		str := builder.String()

		err := IsValidName(str, keywords)
		if err == nil {
			return str, nil
		}
	}

	return "", errors.New("could not fix variable name")
}

// FixVariableName fixes the given variable name by checking whether or not the name
// trimmed to the minimum length is valid. If it is not valid, the function appends
// the given suffix to the name until a valid name is found.
//
// Moreover, if the name is greater than the minimum length, the function first checks whether
// the name can be trimmed to the minimum length and still be valid. If it is not valid, the function
// tries the remaining characters one by one until a valid name is found. Finally, if it is still not
// valid, the function appends the given suffix to the name until a valid name is found.
//
// Parameters:
//   - variable_name: The variable name to fix.
//   - keywords: The list of keywords to check against.
//   - min: The minimum length of the variable name. If less than 1, the function uses 1.
//   - suffix: The suffix to append to the variable name. If empty, the function uses "_".
//
// Returns:
//   - string: The fixed variable name.
//   - bool: True if the variable name is not empty. False otherwise.
func FixVariableName(variable_name string, keywords []string, min int, suffix string) (string, bool) {
	if variable_name == "" {
		return "", false
	}

	if min < 1 {
		min = 1
	}

	if suffix == "" {
		suffix = "_"
	}

	var_name, err := fix_variable_name(variable_name, keywords, min)
	if err == nil {
		return var_name, true
	}

	for {
		variable_name = variable_name + suffix

		err := IsValidName(variable_name, keywords)
		if err == nil {
			return variable_name, true
		}
	}
}

// FixVarNameIncremental works like FixVariableName but it appends an incremental number
// to the variable name instead of a suffix.
//
// Parameters:
//   - variable_name: The variable name to fix.
//   - keywords: The list of keywords to check against.
//   - min: The minimum length of the variable name. If less than 1, the function uses 1.
//   - start: The starting number to append to the variable name. If less than 0, the function uses 0.
//
// Returns:
//   - string: The fixed variable name.
//   - bool: True if the variable name is not empty. False otherwise.
func FixVarNameIncremental(variable_name string, keywords []string, min int, start int) (string, bool) {
	if variable_name == "" {
		return "", false
	}

	if min < 1 {
		min = 1
	}

	if start < 0 {
		start = 0
	}

	var_name, err := fix_variable_name(variable_name, keywords, min)
	if err == nil {
		return var_name, true
	}

	for i := start; ; i++ {
		tmp := variable_name + strconv.Itoa(i)

		err := IsValidName(tmp, keywords)
		if err == nil {
			return tmp, true
		}
	}
}
