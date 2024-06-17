package CodeWriter

import (
	"fmt"

	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
)

/*
// WriteSlice writes a slice of CodeWritable values to a string.
//
// Format:
//
//	var name []T
//
//	func init() {
//		name = []T{
//			v1,
//			v2,
//			...,
//			vn,
//		}
//	}
//
// Parameters:
//   - name: The name of the slice.
//   - slice: The slice of CodeWritable values.
//   - opts: Optional options to pass to the CodeWritable values.
//
// Returns:
//   - string: The string representation of the slice.
func WriteSlice[T CodeWritable](name string, slice []T, opts ...Option) string {
	var builder strings.Builder

	builder.WriteString("var ")
	builder.WriteString(name)
	builder.WriteString(" []")
	fmt.Fprintf(&builder, "%T", slice)

	if len(slice) != 0 {
		builder.WriteString("\n\nfunc init() {\n\t")
		builder.WriteString(name)
		builder.WriteString(" = []")
		fmt.Fprintf(&builder, "%T", slice)
		builder.WriteString("{\n\t\t")

		values := make([]string, 0, len(slice))
		for _, v := range slice {
			values = append(values, v.WriteGo(opts...))
		}

		builder.WriteString(strings.Join(values, ",\n\t\t"))

		builder.WriteString(",\n\t}\n}\n")
	}

	return builder.String()
}

// WriteInitFunction writes the init function for a slice of CodeWritable values to a string.
//
// Format:
//
//	func init() {
//		name = []T{
//			v1,
//			v2,
//			...,
//			vn,
//		}
//	}
//
// Parameters:
//   - name: The name of the slice.
//   - slice: The slice of CodeWritable values.
//   - opts: Optional options to pass to the CodeWritable values.
//
// Returns:
//   - string: The string representation of the init function.
func WriteInitFunction[T CodeWritable](name string, slice []T, opts ...Option) string {
	var builder strings.Builder

	builder.WriteString("func init() {\n\t")
	builder.WriteString(name)
	builder.WriteString(" = []")
	fmt.Fprintf(&builder, "%T", slice)
	builder.WriteString("{\n\t\t")

	values := make([]string, 0, len(slice))
	for _, v := range slice {
		values = append(values, v.WriteGo(opts...))
	}

	builder.WriteString(strings.Join(values, ",\n\t\t"))

	builder.WriteString(",\n\t}\n}\n")

	return builder.String()
}
*/

type FunctionDecl[T ffs.FStringer] struct {
	Name string
	Body []T
}

// WriteFunction writes a function with a body of CodeWritable values to a string.
//
// Format:
//
//	func name() {
//		instructions...
//	}
//
// Parameters:
//   - name: The name of the function.
//   - body: The body of the function.
//   - opts: Optional options to pass to the CodeWritable values.
//
// Returns:
//   - string: The string representation of the function.
//
// Behaviors:
//   - If the body is empty, the function will be written without a body.
func (fd *FunctionDecl[T]) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	if trav == nil {
		return nil
	}

	if len(fd.Body) == 0 {
		fmt.Fprintf(trav, "func %s() {}", fd.Name)
		trav.AcceptWord()
		return nil
	}

	fmt.Fprintf(trav, "func %s() {", fd.Name)
	trav.AcceptWord()

	for _, v := range fd.Body {
		err := ffs.ApplyForm(
			trav.GetConfig(
				ffs.WithModifiedIndent(1),
			),
			trav,
			v,
		)
		if err != nil {
			return err
		}

		trav.EmptyLine()
	}

	err := trav.AppendRune('}')
	if err != nil {
		return err
	}

	trav.AcceptWord()

	return nil
}

type SliceDecl[T ffs.FStringer] struct {
	Name  string
	Slice []T
}

func (s *SliceDecl[T]) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	if trav == nil {
		return nil
	}

	fmt.Fprintf(trav, "%s = []%T{", s.Name, s.Slice)
	trav.AcceptWord()

	for _, v := range s.Slice {
		err := ffs.ApplyForm(
			trav.GetConfig(
				ffs.WithModifiedIndent(1),
				ffs.WithRightDelimiter(","),
			),
			trav,
			v,
		)
		if err != nil {
			return err
		}

		err = trav.AppendRune(',')
		if err != nil {
			return err
		}

		trav.AcceptLine()
	}

	err := trav.AppendRune('}')
	if err != nil {
		return err
	}

	trav.AcceptWord()

	return nil
}
