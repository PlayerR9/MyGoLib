package Go

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	utse "github.com/PlayerR9/MyGoLib/Utility/StringExt"
)

type Parser struct {
	chars []rune
	idx   int
}

// isDone returns true if the parser has reached the end of the character stream.
func (p *Parser) isDone() bool {
	return p.idx >= len(p.chars)
}

// is_valid_identifier_char checks if the given character is a valid identifier character.
//
// Parameters:
//   - char: The character to check.
//
// Returns:
//   - bool: True if the char is a digit. False otherwise.
//   - bool: True if the character is a valid identifier character. False otherwise.
func is_valid_identifier_char(char rune) (bool, bool) {
	if char == '_' {
		return false, true
	}

	ok := unicode.IsLetter(char)
	if ok {
		return false, true
	}

	ok = unicode.IsDigit(char)
	if ok {
		return true, true
	}

	return false, false
}

// parseID parses an identifier from the character stream.
//
// Returns:
//   - string: The parsed identifier.
//   - bool: True if an identifier is successfully parsed, false otherwise.
func (p *Parser) parseID() (string, bool) {
	if p.idx >= len(p.chars) {
		return "", false
	}

	char := p.chars[p.idx]

	is_num, ok := is_valid_identifier_char(char)
	if !ok || is_num {
		return "", false
	}

	var i int
	for i = p.idx + 1; i < len(p.chars); i++ {
		char := p.chars[i]

		_, ok := is_valid_identifier_char(char)
		if !ok {
			break
		}
	}

	id := string(p.chars[p.idx:i])
	p.idx = i

	return id, true
}

// skipWS skips whitespace characters in the character stream.
func (p *Parser) skipWS() {
	if p.idx >= len(p.chars) {
		return
	}

	char := p.chars[p.idx]

	ok := unicode.IsSpace(char)
	if !ok {
		return
	}

	for i := p.idx + 1; i < len(p.chars); i++ {
		char := p.chars[i]

		ok := unicode.IsSpace(char)
		if !ok {
			p.idx = i
			return
		}
	}

	p.idx = len(p.chars)
}

// sub_parse_fields is a helper function that parses the fields of a struct according
// to the following simplified EBNF rule:
//
//	Field1 = name .
//	Field1 = name { "," name } .
//
// Parameters:
//   - chars: The characters to parse.
//
// Returns:
//   - []string: The parsed fields.
//   - int: The length of the parsed fields. -1 if no fields are found.
func sub_parse_fields(p *Parser) []string {
	uc.AssertParam("p", p != nil, errors.New("p cannot be nil"))

	if p.idx >= len(p.chars) {
		return nil
	}

	var fields []string

	for {
		id, ok := p.parseID()
		if !ok {
			return fields
		}

		fields = append(fields, id)

		p.skipWS()

		ok = p.isDone()
		if ok {
			break
		}

		char := p.chars[p.idx]
		if char != ',' {
			break
		}

		p.idx++

		p.skipWS()
	}

	return fields
}

// parse_type_name parses and returns the type name from the input Parser.
//
// Parameter:
//   - p: the Parser containing the type name to be parsed.
//
// Return:
//   - string: the parsed type name.
//   - error: an error if the parsing fails.
func parse_type_name(p *Parser) (string, error) {
	uc.AssertParam("p", p != nil, errors.New("p cannot be nil"))

	var builder strings.Builder

	id, ok := p.parseID()
	if !ok {
		return "", errors.New("no identifier found")
	}

	builder.WriteString(id)

	ok = p.isDone()
	if ok {
		id := builder.String()

		return id, nil
	}

	char := p.chars[p.idx]
	if char != '.' {
		id := builder.String()

		return id, nil
	}

	p.idx++
	builder.WriteRune('.')

	id, ok = p.parseID()
	if !ok {
		return "", errors.New("after '.' expected an identifier")
	}

	builder.WriteString(id)

	id = builder.String()

	return id, nil
}

// ParseField parses the given field.
//
// A field is defined by the following EBNF rule:
//
//	Field = name { "," name } type .
//
// where:
//   - name: The a valid Go identifier for the field.
//   - type: The type of the field.
//
// When two or more names are present, then they all share the specified type.
func parse_field(p *Parser) (string, []string, error) {
	uc.AssertParam("p", p != nil, errors.New("p cannot be nil"))

	// Field = Field1 type .
	//
	// Field1 = name .
	// Field1 = name { "," name } .

	ids := sub_parse_fields(p)
	if len(ids) == 0 {
		return "", nil, errors.New("no fields were found")
	}

	p.skipWS()

	if p.isDone() {
		return "", nil, errors.New("expected a type name")
	}

	type_name, err := parse_type_name(p)
	if err != nil {
		return "", nil, fmt.Errorf("while parsing type name: %w", err)
	}

	return type_name, ids, nil
}

// ParseFields parses the given string to extract fields and their types.
//
// Parameters:
//   - str: The string containing fields to parse.
//
// Returns:
//   - map[string]string: A map of field names to their corresponding types.
//   - error: An error if parsing encounters any issues.
func ParseFields(str string) (map[string]string, error) {
	if str == "" {
		return nil, nil
	}

	chars, err := utse.ToUTF8Runes(str)
	if err != nil {
		return nil, err
	}

	p := &Parser{
		chars: chars,
		idx:   0,
	}

	field_map := make(map[string]string)

	for {
		p.skipWS()

		ok := p.isDone()
		if ok {
			break
		}

		typename, ids, err := parse_field(p)
		if err != nil {
			return nil, fmt.Errorf("while parsing field: %w", err)
		}

		for _, id := range ids {
			_, ok := field_map[id]
			if ok {
				return nil, fmt.Errorf("duplicate field: %s", id)
			}

			field_map[id] = typename
		}

		p.skipWS()

		ok = p.isDone()
		if ok {
			break
		}

		char := p.chars[p.idx]

		if char == ',' {
			p.idx++
		}
	}

	return field_map, nil
}

// ParseGenerics parses a string representing a list of generic types enclosed in square brackets.
//
// Parameters:
//   - str: The string to parse.
//
// Return types:
//   - []rune: An array of runes representing the parsed generic types.
//   - error: An error if the parsing fails.
func ParseGenerics(str string) ([]rune, error) {
	ok := strings.HasSuffix(str, "]")
	if !ok {
		return nil, nil
	}

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

	var letters []rune

	for i, field := range fields {
		letter, err := IsGenericsID(field)
		if err != nil {
			err := uc.NewErrAt(i+1, "field", err)
			return nil, err
		}

		letters = append(letters, letter)
	}

	return letters, nil
}
