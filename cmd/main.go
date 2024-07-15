// -type=<type> [ -name=<name> ]
//
// **Flag: Type**
//
// The flag "type" is used to specify the data type of the linked stack. This must be set and it specifies the data
// type of the linked stack. For instance, using the flag "type=int" will create a linked stack of integers.
//
// **Flag: Name**
//
// The flag "name" is used to specify a custom name for the linked stack. If set it must be a valid Go identifier that
// starts with an uppercase letter. On the other hand, if not set, the default name of "Linked<DataType>Stack" will
// be used instead; where <DataType> is the data type of the linked stack.
package main

import (
	"flag"
	"log"
	"text/template"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	utgo "github.com/PlayerR9/MyGoLib/Utility/Go"
	ggen "github.com/PlayerR9/MyGoLib/Utility/go_generator"
)

var (
	// Logger is the logger to use.
	Logger *log.Logger

	// t is the template to use.
	t *template.Template
)

func init() {
	Logger = ggen.InitLogger("stack")

	t = template.Must(template.New("").Parse(templ))
}

var (
	// DataType is the data type of the linked stack.
	DataType *string

	// TypeName is the name of the linked stack.
	TypeName *string
)

func init() {
	ggen.SetOutputFlag("<type>_stack.go", false)

	DataType = flag.String("type", "", "the data type of the linked stack. This must be set and it is "+
		"the data type of the linked stack.")

	TypeName = flag.String("name", "", "the name of the linked stack. Must be a valid Go identifier. If not set, "+
		"the default name of 'Linked<DataType>Stack' will be used instead.")
}

type GenData struct {
	TypeName    string
	DataType    string
	PackageName string
}

func (g GenData) SetPackageName(name string) ggen.Generater {
	g.PackageName = name

	return g
}

func main() {
	err := ggen.ParseFlags()
	if err != nil {
		Logger.Fatalf("Invalid flags: %s", err.Error())
	}

	data_type := uc.AssertNil(DataType, "DataType")

	err = ggen.IsValidName(data_type, nil, ggen.Exported)
	if err != nil {
		Logger.Fatalf("Name of the data type is invalid: %s", err.Error())
	}

	data_type = "Linked" + data_type + "Stack"

	dest := "stack.go"

	output_loc, err := utgo.FixImportDir(dest)
	if err != nil {
		Logger.Fatalf("Could not fix import dir: %s", err.Error())
	}

	err = ggen.Generate(output_loc, GenData{}, t,
		func(data *GenData) error {
			type_name := uc.AssertNil(TypeName, "TypeName")
			data.TypeName = type_name

			return nil
		},
		func(data *GenData) error {
			data.DataType = data_type

			return nil
		},
	)
	if err != nil {
		Logger.Fatalf("Could not generate code: %s", err.Error())
	}
}

const templ = `// Code generated with go generate. DO NOT EDIT.
package {{ .PackageName }}

import (
	"strconv"
	"strings"

	"github.com/PlayerR9/MyGoLib/Units/common"

	{{- if ne .PackageName "stack" }}
		"github.com/PlayerR9/stack"
	{{- end }}
)

type stack_node_{{ .DataType }} struct {
	value {{ .DataType }}
	next *stack_node_{{ .DataType }}
}

// {{ .TypeName }} is a stack of {{ .DataType }} values implemented without a maximum capacity
// and using a linked list.
type {{ .TypeName }} struct {
	front *stack_node_{{ .DataType }}
	size int
}

// New{{ .TypeName }} creates a new linked stack.
//
// Returns:
//   - *{{ .TypeName }}: A pointer to the newly created stack.
func New{{ .TypeName }}() *{{ .TypeName }} {
	return &{{ .TypeName }}{
		size: 0,
	}
}

// Push implements the Stacker interface.
//
// Always returns true.
func (s *{{ .TypeName }}) Push(value {{ .DataType }}) bool {
	node := &stack_node_{{ .DataType }}{
		value: value,
	}

	if s.front != nil {
		node.next = s.front
	}

	s.front = node
	s.size++

	return true
}

// PushMany implements the Stacker interface.
//
// Always returns the number of values pushed onto the stack.
func (s *{{ .TypeName }}) PushMany(values []{{ .DataType }}) int {
	if len(values) == 0 {
		return 0
	}

	node := &stack_node_{{ .DataType }}{
		value: value,
	}

	if s.front != nil {
		node.next = s.front
	}

	s.front = node

	for i := 1; i < len(values); i++ {
		node := &stack_node_{{ .DataType }}{
			value: value,
			next:  s.front,
		}

		s.front = node
	}

	s.size += len(values)
	
	return len(values)
}

// Pop implements the Stacker interface.
func (s *{{ .TypeName }}) Pop() ({{ .DataType }}, bool) {
	if s.front == nil {
		return *new({{ .DataType }}), false
	}

	to_remove := s.front
	s.front = s.front.next

	s.size--
	to_remove.next = nil

	return to_remove.value, true
}

// Peek implements the Stacker interface.
func (s *{{ .TypeName }}) Peek() ({{ .DataType }}, bool) {
	if s.front == nil {
		return *new({{ .DataType }}), false
	}

	return s.front.value, true
}

// IsEmpty implements the Stacker interface.
func (s *{{ .TypeName }}) IsEmpty() bool {
	return s.front == nil
}

// Size implements the Stacker interface.
func (s *{{ .TypeName }}) Size() int {
	return s.size
}

// Iterator implements the Stacker interface.
func (s *{{ .TypeName }}) Iterator() common.Iterater[{{ .DataType }}] {
	var builder common.Builder[{{ .DataType }}]

	for node := s.front; node != nil; node = node.next {
		builder.Add(node.value)
	}

	return builder.Build()
}

// Clear implements the Stacker interface.
func (s *{{ .TypeName }}) Clear() {
	if s.front == nil {
		return
	}

	prev := s.front

	for node := s.front.next; node != nil; node = node.next {
		prev = node
		prev.next = nil
	}

	prev.next = nil

	s.front = nil
	s.size = 0
}

// GoString implements the Stacker interface.
func (s *{{ .TypeName }}) GoString() string {
	values := make([]string, 0, s.size)
	for node := s.front; node != nil; node = node.next {
		values = append(values, common.StringOf(node.value))
	}

	var builder strings.Builder

	builder.WriteString("{{ .TypeName }}[size=")
	builder.WriteString(strconv.Itoa(s.size))
	builder.WriteString(", values=[")
	builder.WriteString(strings.Join(values, ", "))
	builder.WriteString(" →]]")

	return builder.String()
}

// Slice implements the Stacker interface.
//
// The 0th element is the top of the stack.
func (s *{{ .TypeName }}) Slice() []{{ .DataType }} {
	slice := make([]{{ .DataType }}, 0, s.size)

	for node := s.front; node != nil; node = node.next {
		slice = append(slice, node.value)
	}

	return slice
}

// Copy implements the Stacker interface.
//
// The copy is a shallow copy.
func (s *{{ .TypeName }}) Copy() common.Copier {
	if s.front == nil {
		return &{{ .TypeName }}{}
	}

	s_copy := &{{ .TypeName }}{
		size: s.size,
	}

	node_copy := &stack_node_{{ .DataType }}{
		value: s.front.value,
	}

	s_copy.front = node_copy

	prev := node_copy

	for node := s.front.next; node != nil; node = node.next {
		node_copy := &stack_node_{{ .DataType }}{
			value: node.value,
		}

		prev.next = node_copy

		prev = node_copy
	}

	return s_copy
}

// Capacity implements the Stacker interface.
//
// Always returns -1.
func (s *{{ .TypeName }}) Capacity() int {
	return -1
}

// IsFull implements the Stacker interface.
//
// Always returns false.
func (s *{{ .TypeName }}) IsFull() bool {
	return false
}
`
