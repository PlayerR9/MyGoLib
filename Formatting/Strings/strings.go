package Strings

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

// AndString concatenates a list of strings using commas and the word "and" before the last string.
//
// Parameters:
//
//   - vals: The list of strings to concatenate.
//
// Returns:
//
//   - string: The concatenated string.
func AndString(vals ...string) string {
	switch len(vals) {
	case 0:
		return ""
	case 1:
		return vals[0]
	case 2:
		var builder strings.Builder

		builder.WriteString(vals[0])
		builder.WriteString(" and ")
		builder.WriteString(vals[1])

		return builder.String()
	default:
		var builder strings.Builder

		builder.WriteString(strings.Join(vals[0:len(vals)-1], ", "))
		builder.WriteString(", and ")
		builder.WriteString(vals[len(vals)-1])

		return builder.String()
	}
}

// OrString concatenates a list of strings using commas and the word "or" before the last string.
//
// Parameters:
//
//   - vals: The list of strings to concatenate.
//
// Returns:
//
//   - string: The concatenated string.
func OrString(vals ...string) string {
	switch len(vals) {
	case 0:
		return ""
	case 1:
		return vals[0]
	case 2:
		var builder strings.Builder

		builder.WriteString(vals[0])
		builder.WriteString(" or ")
		builder.WriteString(vals[1])

		return builder.String()
	default:
		var builder strings.Builder

		builder.WriteString(strings.Join(vals[0:len(vals)-1], ", "))
		builder.WriteString(", or ")
		builder.WriteString(vals[len(vals)-1])

		return builder.String()
	}
}

// DateStringer prints the date in the format "1st January, 2006".
//
// Parameters:
//
//   - date: The date to print.
//
// Returns:
//
//   - string: The date in the format "1st January, 2006".
func DateStringer(date time.Time) string {
	var builder strings.Builder

	builder.WriteString(uc.GetOrdinalSuffix(date.Day()))
	builder.WriteRune(' ')
	builder.WriteString(date.Month().String())
	builder.WriteString(", ")
	builder.WriteString(strconv.Itoa(date.Year()))

	return builder.String()
}

// TimeStringer prints the time in the format "3:04 PM".
//
// Parameters:
//
//   - time: The time to print.
//
// Returns:
//
//   - string: The time in the format "3:04 PM".
func TimeStringer(time time.Time) string {
	return time.Format("3:04 PM")
}

// StringsJoiner joins a list of fmt.Stringer values using a separator.
//
// Parameters:
//   - values: The list of fmt.Stringer values to join.
//   - sep: The separator to use when joining the strings.
//
// Returns:
//   - string: The string representation of the values.
func StringsJoiner[T fmt.Stringer](values []T, sep string) string {
	stringValues := make([]string, 0, len(values))

	for _, value := range values {
		stringValues = append(stringValues, value.String())
	}

	return strings.Join(stringValues, sep)
}

// ArrayFormatter formats a list of strings as an array.
//
// Parameters:
//   - values: The list of strings to format.
//
// Returns:
//   - string: The formatted array.
func ArrayFormatter(values []string) string {
	var builder strings.Builder

	builder.WriteRune('[')
	builder.WriteString(strings.Join(values, ", "))
	builder.WriteRune(']')

	return builder.String()
}

// findTabStop is a helper function to find the next tab stop for a string.
//
// Parameters:
//   - s: The string to find the tab stop for.
//   - tabSize: The size of the tab.
//
// Returns:
//   - int: The tab stop.
func findTabStop(s string, tabSize int) int {
	s = strings.TrimRight(s, " ")

	count := utf8.RuneCountInString(s)

	return tabSize * ((count / tabSize) + 1)
}

// padRight is a helper function to pad a string to the right.
//
// Parameters:
//   - s: The string to pad.
//   - length: The length to pad the string to.
//
// Returns:
//   - string: The padded string.
func padRight(s string, length int) string {
	var builder strings.Builder

	builder.WriteString(s)
	builder.WriteString(strings.Repeat(" ", length-utf8.RuneCountInString(s)))

	return builder.String()
}

// TabAlign aligns the tabs of a table's column.
//
// Parameters:
//   - table: The table to align.
//   - column: The column to align.
//   - tabSize: The size of the tab.
//
// Returns:
//   - [][]string: The aligned table.
//   - error: An error of type *errors.ErrInvalidParameter if the tabSize is less than 1
//     or the column is less than 0.
//
// Behaviors:
//   - If the column is not found in the table, the table is returned as is.
func TabAlign(table [][]string, column int, tabSize int) ([][]string, error) {
	if tabSize < 1 {
		return nil, ue.NewErrInvalidParameter("tabSize", ue.NewErrGT(0))
	} else if column < 0 {
		return nil, ue.NewErrInvalidParameter("column", ue.NewErrGTE(0))
	}

	seen := make(map[int]bool)

	for i := 0; i < len(table); i++ {
		if len(table[i]) > column {
			seen[i] = true
		}
	}

	if len(seen) == 0 {
		return table, nil
	}

	stops := make(map[int]int)

	for k := range seen {
		table[k][column] = strings.TrimRight(table[k][column], " ")

		stops[k] = findTabStop(table[k][column], tabSize)
	}

	max := -1

	for _, val := range stops {
		if max == -1 || val > max {
			max = val
		}
	}

	for k := range seen {
		table[k][column] = padRight(table[k][column], max)
	}

	return table, nil
}

// TableEntriesAlign aligns the entries of a table.
//
// Parameters:
//   - table: The table to align.
//   - tabSize: The size of the tab.
//
// Returns:
//   - [][]string: The aligned table.
//   - error: An error if there was an issue aligning the table.
//
// Errors:
//   - *errors.ErrAt: If there was an issue aligning a specific column.
//   - *errors.ErrInvalidParameter: If the tabSize is less than 1.
func TableEntriesAlign(table [][]string, tabSize int) ([][]string, error) {
	if tabSize < 1 {
		return nil, ue.NewErrInvalidParameter("tabSize", ue.NewErrGT(0))
	}

	width := LongestLine(table)
	if width == -1 {
		return table, nil
	}

	var err error

	for i := 0; i < width; i++ {
		table, err = TabAlign(table, i, tabSize)
		if err != nil {
			return nil, ue.NewErrAt(i+1, "column", err)
		}
	}

	return table, nil
}

// LongestLine finds the longest line in a table.
//
// Parameters:
//   - table: The table to find the longest line in.
//
// Returns:
//   - int: The length of the longest line. -1 if the table is empty.
func LongestLine[T any](table [][]T) int {
	if len(table) == 0 {
		return -1
	}

	max := -1

	for i := 0; i < len(table); i++ {
		if len(table[i]) > max {
			max = len(table[i])
		}
	}

	return max
}
