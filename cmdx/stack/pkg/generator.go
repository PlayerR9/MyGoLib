package pkg

import (
	"bytes"
	"html/template"
	"strings"
)

type generator struct {
	type_name string
	data_type string
	pkg_name  string
}

func (g *generator) Generate() ([]byte, error) {
	t := template.Must(
		template.New("stack").Parse(templ),
	)

	type GenData struct {
		TypeName    string
		DataType    string
		CanBeNil    bool
		PackageName string
	}

	can_be_nil := strings.HasPrefix(g.data_type, "*") // Not only this, but interfaces as well

	data := GenData{
		CanBeNil:    can_be_nil,
		DataType:    g.data_type,
		TypeName:    g.type_name,
		PackageName: g.pkg_name,
	}

	var buff bytes.Buffer

	err := t.Execute(&buff, data)
	if err != nil {
		return nil, err
	}

	res := buff.Bytes()
	return res, nil
}

func NewGenerator() *generator {
	g := &generator{
		type_name: *TypeName,
		data_type: *DataType,
		pkg_name:  *OutputDir,
	}

	return g
}

const templ = `// Code generated with go generate. DO NOT EDIT.
package {{ .PackageName }}

import (
	"strconv"
	"strings"

	"github.com/PlayerR9/MyGoLib/Units/common"
	{{- if ne .PackageName "stack" }} "github.com/PlayerR9/stack" {{- end }}
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

{{- if ne .CanBeNil }}
// Copy implements the Stacker interface.
//
// The copy is a deep copy.
{{- else }}
// Copy implements the Stacker interface.
//
// The copy is a shallow copy.
{{- end }}
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
