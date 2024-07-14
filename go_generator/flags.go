package go_generator

import (
	"errors"
	"flag"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

var (
	// OutputLoc is a pointer to the output_flag flag.
	OutputLoc *string

	// IsOutputLocRequired is a flag that specifies whether the output_flag flag is required or not.
	IsOutputLocRequired bool
)

// SetOutputFlag sets the flag that specifies the location of the output file.
//
// Parameters:
//   - def_value: The default value of the output_flag flag.
//   - required: Whether the flag is required or not.
//
// Here are all the possible valid calls to this function:
//
//	SetOutputFlag("", false) <-> SetOutputFlag("[no location]", false)
//	SetOutputFlag("path/to/file.go", false)
//	SetOutputFlag("", true) <-> SetOutputFlag("path/to/file.go", true)
func SetOutputFlag(def_value string, required bool) {
	var usage string

	if required {
		var builder strings.Builder

		builder.WriteString("The location of the output file. ")
		builder.WriteString("It must be set and it must specify a .go file.")

		usage = builder.String()
	} else {
		var def_loc string

		if def_value == "" {
			def_loc = "\"[no location]\""
		} else {
			def_loc = strconv.Quote(def_value)
		}

		var builder strings.Builder

		builder.WriteString("The location of the output file. ")

		builder.WriteString("If set, it must specify a .go file. ")
		builder.WriteString("On the other hand, if not set, the default location of ")
		builder.WriteString(def_loc)
		builder.WriteString(" will be used instead.")

		usage = builder.String()
	}

	OutputLoc = flag.String("-o", "", usage)
	IsOutputLocRequired = required
}

type GenericsValue struct {
	letters []rune
	types   []string
}

func (s *GenericsValue) String() string {
	var values []string
	var builder strings.Builder

	for i, letter := range s.letters {
		builder.WriteRune(letter)
		builder.WriteRune(' ')
		builder.WriteString(s.types[i])

		str := builder.String()
		values = append(values, str)

		builder.Reset()
	}

	joined_str := strings.Join(values, ", ")

	builder.WriteRune('[')
	builder.WriteString(joined_str)
	builder.WriteRune(']')

	str := builder.String()
	return str
}

func parse_generics_value(field string) (rune, string, error) {
	uc.Assert(field != "", "field must not be an empty string")

	sub_fields := strings.Split(field, "/")

	if len(sub_fields) == 1 {
		return '\000', "", errors.New("missing type of generic")
	} else if len(sub_fields) > 2 {
		return '\000', "", errors.New("too many fields")
	}

	left := sub_fields[0]

	if left == "" {
		return '\000', "", uc.NewErrEmpty(left)
	}

	size := utf8.RuneCountInString(left)
	if size > 1 {
		err := errors.New("id must be a single character")
		return '\000', "", err
	}

	letter := rune(left[0])

	ok := unicode.IsUpper(letter)
	if !ok {
		return '\000', "", errors.New("id must be an upper case letter")
	}

	right := sub_fields[1]

	return letter, right, nil
}

func (gv *GenericsValue) add(letter rune, g_type string) error {
	uc.AssertParam("letter", unicode.IsLetter(letter) && unicode.IsUpper(letter), errors.New("letter must be an upper case letter"))
	uc.AssertParam("g_type", g_type != "", errors.New("type must be set"))

	pos, ok := slices.BinarySearch(gv.letters, letter)
	if !ok {
		gv.letters = slices.Insert(gv.letters, pos, letter)
		gv.types = slices.Insert(gv.types, pos, g_type)

		return nil
	}

	if gv.types[pos] != g_type {
		err := fmt.Errorf("duplicate definition for generic %q: %s and %s", string(letter), gv.types[pos], g_type)
		return err
	}

	return nil
}

func (s *GenericsValue) Set(value string) error {
	fields := strings.Split(value, ",")

	for i, field := range fields {
		if field == "" {
			continue
		}

		letter, g_type, err := parse_generics_value(field)
		if err != nil {
			err := uc.NewErrAt(i+1, "field", err)
			return err
		}

		err = s.add(letter, g_type)
		if err != nil {
			return err
		}
	}

	return nil
}
