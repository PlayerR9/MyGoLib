package pkg

import (
	"bytes"
	"html/template"
	"strings"
)

type generator struct {
	data_type string
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
		TypeName:    "Linked" + g.data_type + "Stack",
		PackageName: "pkg",
	}

	var buff bytes.Buffer

	err := t.Execute(&buff, data)
	if err != nil {
		return nil, err
	}

	res := buff.Bytes()
	return res, nil
}

func NewGenerator(data_type string) *generator {
	g := &generator{
		data_type: data_type,
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

// {{ .TypeName }} is a generic type that represents a stack data structure with
// or without a limited capacity, implemented using a linked list.
type {{ .TypeName }} struct {
	front *stack_node_{{ .DataType }}
	size int
}

// New{{ .TypeName }} is a function that creates and returns a new instance of a
// {{ .TypeName }}.
//
// Parameters:
//   - values: A variadic parameter of type {{ .DataType }}, which represents the initial values to be
//     stored in the stack.
//
// Returns:
//   - *{{ .TypeName }}: A pointer to the newly created {{ .TypeName }}.
func New{{ .TypeName }}(values ...{{ .DataType }}) *{{ .TypeName }} {
	s := new({{ .TypeName }})
	s.size = len(values)

	if len(values) == 0 {
		return s
	}

	node := &stack_node_{{ .DataType }}{
		value: values[0],
	}

	s.front = node

	for _, element := range values[1:] {
		node := &stack_node_{{ .DataType }}{
			value: element,
			next:  s.front,
		}

		s.front = node
	}

	return s
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

// Pop implements the Stacker interface.
func (s *{{ .TypeName }}) Pop() ({{ .DataType }}, bool) {
	if s.front == nil {
		return *new({{ .DataType }}), false
	}

	toRemove := s.front
	s.front = s.front.next

	s.size--
	toRemove.next = nil

	return toRemove.value, true
}

// Peek implements the Stacker interface.
func (s *{{ .TypeName }}) Peek() ({{ .DataType }}, bool) {
	if s.front == nil {
		return *new({{ .DataType }}), false
	}

	return s.front.value, true
}

// IsEmpty is a method of the {{ .TypeName }} type. It is used to check if the stack
// is empty.
//
// Returns:
//   - bool: true if the stack is empty, and false otherwise.
func (s *{{ .TypeName }}) IsEmpty() bool {
	return s.front == nil
}

// Size is a method of the {{ .TypeName }} type. It is used to return the number of
// elements in the stack.
//
// Returns:
//   - int: The number of elements in the stack.
func (s *{{ .TypeName }}) Size() int {
	return s.size
}

// Iterator is a method of the {{ .TypeName }} type. It is used to return an iterator
// for the elements in the stack.
//
// Returns:
//   - common.Iterater[{{ .DataType }}]: An iterator for the elements in the stack.
func (s *{{ .TypeName }}) Iterator() common.Iterater[{{ .DataType }}] {
	var builder common.Builder[{{ .DataType }}]

	for node := s.front; node != nil; node = node.next {
		builder.Add(node.value)
	}

	return builder.Build()
}

// Clear is a method of the {{ .TypeName }} type. It is used to remove aCommon elements
// from the stack.
func (s *{{ .TypeName }}) Clear() {
	if s.front == nil {
		return // Stack is already empty
	}

	// 1. First node
	prev := s.front

	// 2. Subsequent nodes
	for node := s.front.next; node != nil; node = node.next {
		prev = node
		prev.next = nil
	}

	prev.next = nil

	// 3. Reset list fields
	s.front = nil
	s.size = 0
}

// GoString implements the fmt.GoStringer interface.
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

// CutNilValues is a method of the {{ .TypeName }} type. It is used to remove aCommon nil
// values from the stack.
func (s *{{ .TypeName }}) CutNilValues() {
	{{- if ne .CanBeNil }}
		return
	{{- else }}
		if s.front == nil {
			return // Stack is empty
		}

		top := s.front.value

		if top == nil && s.front.next == nil {
			// Single node
			s.front = nil
			s.size = 0

			return
		}

		var toDelete *stack_node_{{ .DataType }} = nil

		// 1. First node
		if top == nil {
			toDelete = s.front

			s.front = s.front.next

			toDelete.next = nil
			s.size--
		}

		prev := s.front

		// 2. Subsequent nodes (except last)
		node := s.front.next
		for ; node.next != nil; node = node.next {
			if node.value != nil {
				prev = node
			} else {
				prev.next = node.next
				s.size--

				if toDelete != nil {
					toDelete.next = nil
				}

				toDelete = node
			}
		}

		if toDelete != nil {
			toDelete.next = nil
		}

		// 3. Last node
		if s.front.value == nil {
			node = prev
			node.next = nil
			s.size--
		}
	{{- end }}
}

// Slice is a method of the {{ .TypeName }} type. It is used to return a slice of the
// elements in the stack.
//
// Returns:
//   - []{{ .DataType }}: A slice of the elements in the stack.
func (s *{{ .TypeName }}) Slice() []{{ .DataType }} {
	slice := make([]{{ .DataType }}, 0, s.size)

	for node := s.front; node != nil; node = node.next {
		slice = append(slice, node.value)
	}

	return slice
}

// Copy is a method of the {{ .TypeName }} type. It is used to create a shaCommonow copy
// of the stack.
//
// Returns:
//   - common.Copier: A copy of the stack.
func (s *{{ .TypeName }}) Copy() common.Copier {
	stackCopy := &{{ .TypeName }}{
		size: s.size,
	}

	if s.front == nil {
		return stackCopy
	}

	node_copy := &stack_node_{{ .DataType }}{
		value: s.front.value,
	}

	stackCopy.front = node_copy

	// Subsequent nodes
	prev := node_copy

	for node := s.front.next; node != nil; node = node.next {
		node_copy := &stack_node_{{ .DataType }}{
			value: node.value,
		}

		prev.next = node_copy

		prev = node_copy
	}

	return stackCopy
}

// Capacity is a method of the {{ .TypeName }} type. It is used to return the maximum
// number of elements that the stack can store.
//
// Returns:
//   - int: -1
func (s *{{ .TypeName }}) Capacity() int {
	return -1
}

// IsFull is a method of the {{ .TypeName }} type. It is used to check if the stack is
// full.
//
// Returns:
//   - bool: false
func (s *{{ .TypeName }}) IsFull() bool {
	return false
}
`
