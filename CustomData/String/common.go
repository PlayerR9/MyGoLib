package String

import "strings"

// Join joins a slice of strings with a separator.
//
// Parameters:
//   - elems: The slice of strings to join.
//   - sep: The separator to join the strings with.
//
// Returns:
//   - *String: The joined string.
//
// Behaviors:
//   - If the slice of strings is empty, a nil string is returned.
//   - The style of the first string is used for the joined string.
//   - The separator is not added to the end of the joined string.
func Join(elems []*String, sep string) *String {
	switch len(elems) {
	case 0:
		return nil
	case 1:
		return elems[0].Copy()
	default:
		var builder strings.Builder

		builder.WriteString(elems[0].content)

		for _, elem := range elems[1:] {
			builder.WriteString(sep)
			builder.WriteString(elem.content)
		}

		return NewString(builder.String())
	}
}

// FieldsFunc splits a string into fields using a function to determine the separator.
//
// Parameters:
//   - s: The string to split.
//   - sep: The function to determine the separator.
//
// Returns:
//   - []*String: The fields of the string.
func FieldsFunc(s *String, sep string) []*String {
	if sep == "" {
		return []*String{s.Copy()}
	} else {
		fields := make([]*String, 0)
		var builder strings.Builder

		runes := []rune(sep)
		counter := 0

		for _, r := range s.content {
			if r != runes[counter] {
				counter = 0
				builder.WriteRune(r)

				continue
			}

			counter++

			if counter == len(runes) {
				fields = append(fields, NewString(builder.String()))
				builder.Reset()
				counter = 0
			}
		}

		if builder.Len() > 0 {
			fields = append(fields, NewString(builder.String()))
		}

		return fields
	}
}

// Repeat repeats a string count times.
//
// Parameters:
//   - s: The string to repeat.
//   - count: The number of times to repeat the string.
//
// Returns:
//   - *String: The repeated string.
//
// Behaviors:
//   - If the count is less than or equal to 0, nil is returned.
func Repeat(s *String, count int) *String {
	if count <= 0 {
		return nil
	}

	var builder strings.Builder

	for i := 0; i < count; i++ {
		builder.WriteString(s.content)
	}

	return NewString(builder.String())
}
