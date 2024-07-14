package go_generator

import (
	"errors"
	"go/build"
	"path/filepath"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
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

// IsGenericsID checks if the input string is a valid single upper case letter and returns it as a rune.
//
// Parameters:
//   - id: The id to check.
//
// Returns:
//   - rune: The valid single upper case letter.
//   - error: An error of type *ErrInvalidID if the input string is not a valid identifier.
func IsGenericsID(id string) (rune, error) {
	if id == "" {
		return '\000', NewErrInvalidID(id, uc.NewErrEmpty(id))
	}

	size := utf8.RuneCountInString(id)
	if size > 1 {
		return '\000', NewErrInvalidID(id, errors.New("value must be a single character"))
	}

	letter := rune(id[0])

	ok := unicode.IsUpper(letter)
	if !ok {
		return '\000', NewErrInvalidID(id, errors.New("value must be an upper case letter"))
	}

	return letter, nil
}

// ParseGenerics parses a string representing a list of generic types enclosed in square brackets.
//
// Parameters:
//   - str: The string to parse.
//
// Returns:
//   - []rune: An array of runes representing the parsed generic types.
//   - error: An error if the parsing fails.
//
// Errors:
//   - *ErrNotGeneric: The string is not a valid list of generic types.
//   - error: An error if the string is a possibly valid list of generic types but fails to parse.
func ParseGenerics(str string) ([]rune, error) {
	if str == "" {
		return nil, NewErrNotGeneric(uc.NewErrEmpty(str))
	}

	var letters []rune

	ok := strings.HasSuffix(str, "]")
	if ok {
		idx := strings.Index(str, "[")
		if idx == -1 {
			err := errors.New("missing opening square bracket")
			return nil, err
		}

		generic := str[idx+1 : len(str)-1]
		if generic == "" {
			err := errors.New("empty generic type")
			return nil, err
		}

		fields := strings.Split(generic, ",")

		for i, field := range fields {
			letter, err := IsGenericsID(field)
			if err != nil {
				err := uc.NewErrAt(i+1, "field", err)
				return nil, err
			}

			letters = append(letters, letter)
		}
	} else {
		letter, err := IsGenericsID(str)
		if err != nil {
			err := NewErrNotGeneric(err)
			return nil, err
		}

		letters = append(letters, letter)
	}

	return letters, nil
}

// FixImportDir takes a destination string and manipulates it to get the correct import path.
//
// Parameters:
//   - dest: The destination path.
//
// Returns:
//   - string: The correct import path.
//   - error: An error if there is any.
func FixImportDir(dest string) (string, error) {
	if dest == "" {
		dest = "."
	}

	dir := filepath.Dir(dest)
	if dir == "." {
		pkg, err := build.ImportDir(".", 0)
		if err != nil {
			return "", err
		}

		return pkg.Name, nil
	}

	_, right := filepath.Split(dir)
	return right, nil
}

// MakeTypeSig creates a type signature from a type name and a suffix.
//
// It also adds the generic signature if it exists.
//
// Parameters:
//   - type_name: The name of the type.
//   - suffix: The suffix of the type.
//
// Returns:
//   - string: The type signature.
//   - error: An error if the type signature cannot be created. (i.e., the type name is empty)
func MakeTypeSig(type_name string, suffix string) (string, error) {
	if type_name == "" {
		return "", uc.NewErrInvalidParameter("type_name", uc.NewErrEmpty(type_name))
	}

	var builder strings.Builder

	builder.WriteString(type_name)
	builder.WriteString(suffix)

	if GenericsSigFlag == nil {
		return builder.String(), nil
	}

	if len(GenericsSigFlag.letters) > 0 {
		str := GenericsSigFlag.GetSignature()
		builder.WriteString(str)
	}

	return builder.String(), nil
}

// FixOutputLoc fixes the output location.
//
// Parameters:
//   - type_name: The name of the type.
//   - suffix: The suffix of the type.
//
// Returns:
//   - string: The output location.
//   - error: An error if any.
//
// Errors:
//   - *common.ErrInvalidParameter: If the type name is empty.
//   - *common.ErrInvalidUsage: If the OutputLoc flag was not set.
//   - error: Any other error that may have occurred.
func FixOutputLoc(type_name, suffix string) (string, error) {
	output_loc, err := GetOutputLoc()
	if err != nil {
		return "", err
	}

	if type_name == "" {
		return "", uc.NewErrInvalidParameter("type_name", uc.NewErrEmpty(type_name))
	}

	var filename string

	if output_loc == "" {
		var builder strings.Builder

		str := strings.ToLower(type_name)
		builder.WriteString(str)
		builder.WriteString(suffix)

		filename = builder.String()
	} else {
		filename = output_loc
	}

	if output_loc == "" {
		if IsOutputLocRequiredFlag {
			return "", errors.New("flag must be set")
		}

		output_loc = filename
	}

	ext := filepath.Ext(output_loc)
	if ext == "" {
		return "", errors.New("location cannot be a directory")
	} else if ext != ".go" {
		return "", errors.New("location must be a .go file")
	}

	return output_loc, nil
}

// GoExport is an enum that represents whether a variable is exported or not.
type GoExport int

const (
	// NotExported represents a variable that is not exported.
	NotExported GoExport = iota

	// Exported represents a variable that is exported.
	Exported

	// Either represents a variable that is either exported or not exported.
	Either
)

// IsValidName checks if the given variable name is valid.
//
// This function checks if the variable name is not empty and if it is not a
// Go reserved keyword. It also checks if the variable name is not in the list
// of keywords.
//
// Parameters:
//   - variable_name: The variable name to check.
//   - keywords: The list of keywords to check against.
//   - exported: Whether the variable is exported or not.
//
// Returns:
//   - error: An error if the variable name is invalid.
//
// If the variable is exported, the function checks if the variable name starts
// with an uppercase letter. If the variable is not exported, the function checks
// if the variable name starts with a lowercase letter. Any other case, the
// function does not perform any checks.
func IsValidName(variable_name string, keywords []string, exported GoExport) error {
	if variable_name == "" {
		err := uc.NewErrEmpty(variable_name)
		return err
	}

	switch exported {
	case NotExported:
		r, _ := utf8.DecodeRuneInString(variable_name)
		if r == utf8.RuneError {
			return errors.New("invalid UTF-8 encoding")
		}

		ok := unicode.IsLower(r)
		if !ok {
			return errors.New("identifier must start with a lowercase letter")
		}

		ok = slices.Contains(GoReservedKeywords, variable_name)
		if ok {
			err := errors.New("name is a reserved keyword")
			return err
		}
	case Exported:
		r, _ := utf8.DecodeRuneInString(variable_name)
		if r == utf8.RuneError {
			return errors.New("invalid UTF-8 encoding")
		}

		ok := unicode.IsUpper(r)
		if !ok {
			return errors.New("identifier must start with an uppercase letter")
		}
	}

	ok := slices.Contains(keywords, variable_name)
	if ok {
		err := errors.New("name is not allowed")
		return err
	}

	return nil
}