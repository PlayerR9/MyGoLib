// This generates a tree node that uses first child and next sibling pointers.
//
// To use it, run the following command:
//
//	//go:generate treenode -type=<type_name> -data=<data_file> [ -output=<output_file> ]
//
// More specifically, the <data> field must be set and it specifies the type(s) that you want to
// create a tree node for. The syntax of this argument is:
//
//	Argument = "\"" Field { "," Field } "\"".
//	Field = name { "," name } type .
//
// For instance, running the following command:
//
//	//go:generate treenode -type="TreeNode" -data="a, b int, name string"
//
// will generate a tree node with the following fields:
//
//		type TreeNode struct {
//			// Node pointers.
//
//			a, b int
//			name string
//	}
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"
	"strings"
	"text/template"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	utgo "github.com/PlayerR9/MyGoLib/Utility/Go"
)

var (
	// InvalidVarNames is a list of invalid variable names.
	InvalidVarNames []string

	// Logger is the logger for the package.
	Logger *log.Logger
)

func init() {
	InvalidVarNames = []string{
		"ffs",
		"lls",
		"tr",
		"uc",
		"us",
		"trav",
		"err",
		"c",
		"child",
		"parent",
		"p",
		"ok",
		"children",
		"stack",
		"leaves",
		"top",
		"node",
		"h",
		"prev",
		"next",
		"ancestors",
		"slices",
		"target",
		"n",
		"size",
		"parents",
	}

	Logger = log.New(os.Stderr, "[TreeNode]: ", log.LstdFlags)
}

var (
	// TypeNameFlag is the flag for the type name.
	TypeNameFlag *string

	// OutputFileFlag is the flag for the output file.
	OutputFileFlag *string

	// DataFileFlag is the flag for the data file.
	DataFileFlag *string
)

func init() {
	TypeNameFlag = flag.String("type", "",
		"The type name to generate the tree node for. It must be set."+
			" Must start with an upper case letter and must be a valid Go identifier.",
	)

	OutputFileFlag = flag.String("output", "",
		"The output file to write the generated code to. If not set, the default file name is used."+
			" That is \"<type_name>_treenode.go\".",
	)

	DataFileFlag = flag.String("data", "",
		"The data file to read the type names from. It must be set."+
			" The syntax of the data file is described in the documentation.",
	)
}

// generator is a code generator for a tree node.
type generator struct {
	// package_name is the name of the package.
	package_name string

	// type_name is the name of the type.
	type_name string

	// variable_name is the name of the variable.
	variable_name string

	// data is the data to generate the code for.
	data map[string]string
}

// generate generates the code for the tree node.
//
// Returns:
//   - []byte: The generated code.
//   - error: An error if the code could not be generated.
func (g *generator) generate() ([]byte, error) {
	if g.package_name == "" {
		return nil, errors.New("package name must be set")
	}

	if g.type_name == "" {
		return nil, errors.New("type name must be set")
	}

	if g.variable_name == "" {
		return nil, errors.New("variable name must be set")
	}

	t := template.Must(
		template.New("").Parse(templ),
	)

	type GenData struct {
		PackageName  string
		TypeName     string
		OutputType   string
		VariableName string
		Data         map[string]string
	}

	data := GenData{
		PackageName:  g.package_name,
		TypeName:     g.type_name,
		OutputType:   "*" + g.type_name,
		VariableName: g.variable_name,
		Data:         g.data,
	}

	var buf bytes.Buffer

	err := t.Execute(&buf, data)
	if err != nil {
		return nil, err
	}

	result := buf.Bytes()

	return result, nil
}

func parse_flags() (string, string, map[string]string, error) {
	flag.Parse()

	if *TypeNameFlag == "" {
		return "", "", nil, errors.New("the type name must be set")
	}

	if *DataFileFlag == "" {
		return "", "", nil, errors.New("the data file must be set")
	}

	type_name := *TypeNameFlag

	var filename string

	if *OutputFileFlag == "" {
		var builder strings.Builder

		str := strings.ToLower(type_name)
		builder.WriteString(str)
		builder.WriteString("_treenode.go")

		filename = builder.String()
	} else {
		filename = *OutputFileFlag
	}

	dff := *DataFileFlag

	dff = strings.TrimPrefix(dff, "\"")
	dff = strings.TrimSuffix(dff, "\"")

	parsed, err := utgo.ParseFields(dff)
	if err != nil {
		return "", "", nil, fmt.Errorf("could not parse data file: %s", err.Error())
	}

	return type_name, filename, parsed, nil
}

func main() {
	type_name, filename, fields, err := parse_flags()
	if err != nil {
		Logger.Fatalf("Could not parse flags: %s", err.Error())
	}

	// Check if the type name is valid.

	var_name, err := utgo.MakeVariableName(type_name)
	if err != nil {
		Logger.Fatalf("Could not make variable name: %s", err.Error())
	}

	var_name, ok := utgo.FixVarNameIncremental(var_name, InvalidVarNames, 2, 1)
	uc.Assert(ok, "FixVarNameIncremental should not return false")

	pkg, err := build.Default.ImportDir(".", 0)
	if err != nil {
		Logger.Fatalf("Could not import directory: %s", err.Error())
	}

	// Generate the code.

	g := &generator{
		package_name:  pkg.Name,
		type_name:     type_name,
		variable_name: var_name,
		data:          fields,
	}

	generated_data, err := g.generate()
	if err != nil {
		Logger.Fatalf("Could not generate data: %s", err.Error())
	}

	// Write the code to the file.

	err = os.WriteFile(filename, generated_data, 0644)
	if err != nil {
		Logger.Fatalf("Could not write file: %s", err.Error())
	}
}

// templ is the template for the tree node.
const templ = `// Code generated by go generate; EDIT THIS FILE DIRECTLY

package {{ .PackageName }}

import (
	"slices"

	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// {{ .TypeName }}Iterator is a pull-based iterator that iterates
// over the children of a {{ .TypeName }}.
type {{ .TypeName }}Iterator struct {
	parent, current {{ .OutputType }}
}

// Consume implements the common.Iterater interface.
//
// *common.ErrExhaustedIter is the only error returned by this function and the returned
// node is never nil.
func (iter *{{ .TypeName }}Iterator) Consume() (tr.Noder, error) {
	if iter.current == nil {
		return nil, uc.NewErrExhaustedIter()
	}

	node := iter.current
	iter.current = iter.current.NextSibling

	return node, nil
}

// Restart implements the common.Iterater interface.
func (iter *{{ .TypeName }}Iterator) Restart() {
	iter.current = iter.parent.FirstChild
}

// {{ .TypeName }} is a node in a tree.
type {{ .TypeName }} struct {
	Parent, FirstChild, NextSibling, LastChild, PrevSibling {{ .OutputType }}

	{{- range $key, $value := .Data }}
	{{ $key }} {{ $value }}
	{{- end }}
}

// Iterator implements the Tree.Noder interface.
//
// This function iterates over the children of the node, it is a pull-based iterator,
// and never returns nil.
func ({{ .VariableName }} {{ .OutputType }}) Iterator() uc.Iterater[tr.Noder] {
	return &{{ .TypeName }}Iterator{
		parent: {{ .VariableName }},
		current: {{ .VariableName }}.FirstChild,
	}
}

// FString implements the Tree.Noder interface.
func ({{ .VariableName }} {{ .OutputType }}) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	if trav == nil {
		return nil
	}

	// Write here the code to format the node through the traversor.
	err := trav.AddLine("[WRITE HERE THE CODE TO FORMAT THE NODE]")
	if err != nil {
		return err
	}

	// This is safe to modify.
	err = ffs.ApplyFormMany(
		trav.GetConfig(
			ffs.WithModifiedIndent(1),
		),
		trav,
		{{ .VariableName }}.GetChildren(),
	)
	if err != nil {
		return err
	}

	return nil
}

// Copy implements the Noder interface.
//
// It never returns nil and it does not copy the parent or the sibling pointers.
func ({{ .VariableName }} {{ .OutputType }}) Copy() uc.Copier {
	var child_copy []tr.Noder	

	for c := {{ .VariableName }}.FirstChild; c != nil; c = c.NextSibling {
		child_copy = append(child_copy, c.Copy().(tr.Noder))
	}

	// Copy here the data of the node.

	{{ .VariableName}}_copy := &{{ .TypeName }}{
	 	// Add here the copied data of the node.
	}

	{{ .VariableName }}_copy.LinkChildren(child_copy)

	return {{ .VariableName }}_copy
}

// SetParent implements the Noder interface.
func ({{ .VariableName }} {{ .OutputType }}) SetParent(parent tr.Noder) bool {
	if parent == nil {
		{{ .VariableName }}.Parent = nil
		return true
	}

	p, ok := parent.({{ .OutputType }})
	if !ok {
		return false
	}

	{{ .VariableName }}.Parent = p

	return true
}

// GetParent implements the Noder interface.
func ({{ .VariableName }} {{ .OutputType }}) GetParent() tr.Noder {
	return {{ .VariableName }}.Parent
}

// LinkWithParent implements the Noder interface.
//
// Children that are not of type {{ .OutputType }} or nil are ignored.
func ({{ .VariableName }} {{ .OutputType }}) LinkChildren(children []tr.Noder) {
	if len(children) == 0 {
		return
	}

	var valid_children []{{ .OutputType }}

	for _, child := range children {
		if child == nil {
			continue
		}

		c, ok := child.({{ .OutputType }})
		if ok {
			c.Parent = {{ .VariableName }}
			valid_children = append(valid_children, c)
		}		
	}
	
	if len(valid_children) == 0 {
		return
	}

	valid_children[0].PrevSibling = nil
	valid_children[len(valid_children)-1].NextSibling = nil

	if len(valid_children) == 1 {
		return
	}

	for i := 0; i < len(valid_children)-1; i++ {
		valid_children[i].NextSibling = valid_children[i+1]
	}

	for i := 1; i < len(valid_children); i++ {
		valid_children[i].PrevSibling = valid_children[i-1]
	}

	{{ .VariableName }}.FirstChild, {{ .VariableName }}.LastChild = valid_children[0], valid_children[len(valid_children)-1]
}

// GetLeaves implements the Noder interface.
//
// This is expensive as leaves are not stored and so, every time this function is called,
// it has to do a DFS traversal to find the leaves. Thus, it is recommended to call
// this function once and then store the leaves somewhere if needed.
//
// Despite the above, this function does not use recursion and is safe to use.
//
// Finally, no nil nodes are returned.
func ({{ .VariableName }} {{ .OutputType }}) GetLeaves() []tr.Noder {
	// It is safe to change the stack implementation as long as
	// it is not limited in size. If it is, make sure to check the error
	// returned by the Push and Pop methods.
	stack := lls.NewLinkedStack[tr.Noder]({{ .VariableName }})

	var leaves []tr.Noder

	for {
		top, ok := stack.Pop()
		if !ok {
			break
		}

		node := top.({{ .OutputType }})
		if node.FirstChild == nil {
			leaves = append(leaves, top)
		} else {
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				stack.Push(c)
			}
		}
	}

	return leaves
}

// Cleanup implements the Noder interface.
//
// This is expensive as it has to traverse the whole tree to clean up the nodes, one
// by one. While this is useful for freeing up memory, for large enough trees, it is
// recommended to let the garbage collector handle the cleanup.
//
// Despite the above, this function does not use recursion and is safe to use (but
// make sure goroutines are not running on the tree while this function is called).
//
// Finally, it also logically removes the node from the siblings and the parent.
func ({{ .VariableName }} {{ .OutputType }}) Cleanup() {
	type Helper struct {
		previous, current {{ .OutputType }}
	}

	stack := lls.NewLinkedStack[*Helper]()

	// Free the first node.
	for c := {{ .VariableName }}.FirstChild; c != nil; c = c.NextSibling {
		h := &Helper{
			previous:	c.PrevSibling,
			current: 	c,
		}

		stack.Push(h)
	}

	{{ .VariableName }}.FirstChild = nil
	{{ .VariableName }}.LastChild = nil
	{{ .VariableName }}.Parent = nil

	// Free the rest of the nodes.
	for {
		h, ok := stack.Pop()
		if !ok {
			break
		}

		for c := h.current.FirstChild; c != nil; c = c.NextSibling {
			h := &Helper{
				previous:	c.PrevSibling,
				current: 	c,
			}

			stack.Push(h)
		}

		h.previous.NextSibling = nil
		h.previous.PrevSibling = nil

		h.current.FirstChild = nil
		h.current.LastChild = nil
		h.current.Parent = nil
	}

	prev := {{ .VariableName }}.PrevSibling
	next := {{ .VariableName }}.NextSibling

	if prev != nil {
		prev.NextSibling = next
	}

	if next != nil {
		next.PrevSibling = prev
	}

	{{ .VariableName }}.PrevSibling = nil
	{{ .VariableName }}.NextSibling = nil
}

// GetAncestors implements the Noder interface.
//
// This is expensive since ancestors are not stored and so, every time this
// function is called, it has to traverse the tree to find the ancestors. Thus, it is
// recommended to call this function once and then store the ancestors somewhere if needed.
//
// Despite the above, this function does not use recursion and is safe to use.
//
// Finally, no nil nodes are returned.
func ({{ .VariableName }} {{ .OutputType }}) GetAncestors() []tr.Noder {
	var ancestors []tr.Noder

	for node := {{ .VariableName }}; node.Parent != nil; node = node.Parent {
		ancestors = append(ancestors, node.Parent)
	}

	slices.Reverse(ancestors)

	return ancestors
}

// IsLeaf implements the Noder interface.
func ({{ .VariableName }} {{ .OutputType }}) IsLeaf() bool {
	return {{ .VariableName }}.FirstChild == nil
}

// IsSingleton implements the Noder interface.
func ({{ .VariableName }} {{ .OutputType }}) IsSingleton() bool {
	return {{ .VariableName }}.FirstChild != nil && {{ .VariableName }}.FirstChild == {{ .VariableName }}.LastChild
}

// GetFirstChild implements the Noder interface.
func ({{ .VariableName }} {{ .OutputType }}) GetFirstChild() tr.Noder {
	return {{ .VariableName }}.FirstChild
}

// DeleteChild implements the Noder interface.
//
// No nil nodes are returned.
func ({{ .VariableName }} {{ .OutputType }}) DeleteChild(target tr.Noder) []tr.Noder {
	if target == nil {
		return nil
	}

	n, ok := target.({{ .OutputType }})
	if !ok {
		return nil
	}

	children := {{ .VariableName }}.delete_child(n)

	delink_with_parent({{ .VariableName }}, children)

	return children
}

// Size implements the Noder interface.
//
// This is expensive as it has to traverse the whole tree to find the size of the tree.
// Thus, it is recommended to call this function once and then store the size somewhere if needed.
//
// Despite the above, this function does not use recursion and is safe to use.
//
// Finally, the traversal is done in a depth-first manner.
func ({{ .VariableName }} {{ .OutputType }}) Size() int {
	// It is safe to change the stack implementation as long as
	// it is not limited in size. If it is, make sure to check the error
	// returned by the Push and Pop methods.
	stack := lls.NewLinkedStack({{ .VariableName }})

	var size int

	for {
		top, ok := stack.Pop()
		if !ok {
			break
		}

		size++

		for c := top.FirstChild; c != nil; c = c.NextSibling {
			stack.Push(c)
		}
	}

	return size
}

// AddChild adds a new child to the node. If the child is nil or it is not of type
// {{ .OutputType }}, it does nothing.
//
// This function clears the parent and sibling pointers of the child and so, it
// does not add relatives to the child.
//
// Parameters:
//   - child: The child to add.
func ({{ .VariableName }} {{ .OutputType }}) AddChild(child tr.Noder) {
	if child == nil {
		return
	}

	c, ok := child.({{ .OutputType }})
	if !ok {
		return
	}
	
	c.NextSibling = nil
	c.PrevSibling = nil

	last_child := {{ .VariableName }}.LastChild

	if last_child == nil {
		{{ .VariableName }}.FirstChild = c
	} else {
		last_child.NextSibling = c
		c.PrevSibling = last_child
	}

	c.Parent = {{ .VariableName }}
	{{ .VariableName }}.LastChild = c
}

// RemoveNode removes the node from the tree while shifting the children up one level to
// maintain the tree structure.
//
// Also, the returned children can be used to create a forest of trees if the root node
// is removed.
//
// Returns:
//   - []Noder: A slice of pointers to the children of the node iff the node is the root.
//     Nil otherwise.
//
// Example:
//
//	// Given the tree:
//	1
//	├── 2
//	└── 3
//		├── 4
//		└── 5
//	└── 6
//
//	// The tree after removing node 3:
//
//	1
//	├── 2
//	└── 4
//	└── 5
//	└── 6
func ({{ .VariableName }} {{ .OutputType }}) RemoveNode() []tr.Noder {
	prev := {{ .VariableName }}.PrevSibling
	next := {{ .VariableName }}.NextSibling
	parent := {{ .VariableName }}.Parent

	var sub_roots []tr.Noder

	if parent == nil {
		for c := {{ .VariableName }}.FirstChild; c != nil; c = c.NextSibling {
			sub_roots = append(sub_roots, c)
		}
	} else {
		children := parent.delete_child({{ .VariableName }})

		for _, child := range children {
			child.SetParent(parent)
		}
	}

	if prev != nil {
		prev.NextSibling = next
	} else {
		parent.FirstChild = next
	}

	if next != nil {
		next.PrevSibling = prev
	} else {
		parent.Parent.LastChild = prev
	}

	{{ .VariableName }}.Parent = nil
	{{ .VariableName }}.PrevSibling = nil
	{{ .VariableName }}.NextSibling = nil

	delink_with_parent({{ .VariableName }}, sub_roots)

	return sub_roots
}

// New{{ .TypeName }} creates a new node with the given data.
//
// Returns:
//   - *{{ .TypeName }}: A pointer to the newly created node. It is
//   never nil.
func New{{ .TypeName }}() {{ .OutputType }} {
	return &{{ .TypeName }}{
		// Put here the data of the node.
	}
}

// GetLastSibling returns the last sibling of the node. If it has a parent,
// it returns the last child of the parent. Otherwise, it returns the last
// sibling of the node.
//
// As an edge case, if the node has no parent and no next sibling, it returns
// the node itself. Thus, this function never returns nil.
//
// Returns:
//   - {{ .OutputType }}: A pointer to the last sibling.
func ({{ .VariableName }} {{ .OutputType }}) GetLastSibling() {{ .OutputType }} {
	if {{ .VariableName }}.Parent != nil {
		return {{ .VariableName }}.Parent.LastChild
	} else if {{ .VariableName }}.NextSibling == nil {
		return {{ .VariableName }}
	}

	last_sibling := {{ .VariableName }}

	for last_sibling.NextSibling != nil {
		last_sibling = last_sibling.NextSibling
	}

	return last_sibling
}

// GetFirstSibling returns the first sibling of the node. If it has a parent,
// it returns the first child of the parent. Otherwise, it returns the first
// sibling of the node.
//
// As an edge case, if the node has no parent and no previous sibling, it returns
// the node itself. Thus, this function never returns nil.
//
// Returns:
//   - {{ .OutputType }}: A pointer to the first sibling.
func ({{ .VariableName }} {{ .OutputType }}) GetFirstSibling() {{ .OutputType }} {
	if {{ .VariableName }}.Parent != nil {
		return {{ .VariableName }}.Parent.FirstChild
	} else if {{ .VariableName }}.PrevSibling == nil {
		return {{ .VariableName }}
	}

	first_sibling := {{ .VariableName }}

	for first_sibling.PrevSibling != nil {
		first_sibling = first_sibling.PrevSibling
	}

	return first_sibling
}

// IsRoot returns true if the node does not have a parent.
//
// Returns:
//   - bool: True if the node is the root, false otherwise.
func ({{ .VariableName }} {{ .OutputType }}) IsRoot() bool {
	return {{ .VariableName }}.Parent == nil
}

// AddChildren is a convenience function to add multiple children to the node at once.
// It is more efficient than adding them one by one. Therefore, the behaviors are the
// same as the behaviors of the {{ .TypeName }}.AddChild function.
//
// Parameters:
//   - children: The children to add.
func ({{ .VariableName }} {{ .OutputType }}) AddChildren(children []{{ .OutputType }}) {
	children = us.FilterNilValues(children)
	if len(children) == 0 {
		return
	}

	// Deal with the first child
	first_child := children[0]

	first_child.NextSibling = nil
	first_child.PrevSibling = nil

	last_child := {{ .VariableName }}.LastChild

	if last_child == nil {
		{{ .VariableName }}.FirstChild = first_child
	} else {
		last_child.NextSibling = first_child
		first_child.PrevSibling = last_child
	}

	first_child.Parent = {{ .VariableName }}
	{{ .VariableName }}.LastChild = first_child

	// Deal with the rest of the children
	for i := 1; i < len(children); i++ {
		child := children[i]

		child.NextSibling = nil
		child.PrevSibling = nil

		last_child := {{ .VariableName }}.LastChild
		last_child.NextSibling = child
		child.PrevSibling = last_child

		child.Parent = {{ .VariableName }}
		{{ .VariableName }}.LastChild = child
	}
}

// GetChildren returns the immediate children of the node.
//
// The returned nodes are never nil and are not copied. Thus, modifying the returned
// nodes will modify the tree.
//
// Returns:
//   - []tr.Noder: A slice of pointers to the children of the node.
func ({{ .VariableName }} {{ .OutputType }}) GetChildren() []tr.Noder {
	var children []tr.Noder

	for c := {{ .VariableName }}.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}

	return children
}

// HasChild returns true if the node has the given child.
//
// Because children of a node cannot be nil, a nil target will always return false.
//
// Parameters:
//   - target: The child to check for.
//
// Returns:
//   - bool: True if the node has the child, false otherwise.
func ({{ .VariableName }} {{ .OutputType }}) HasChild(target {{ .OutputType }}) bool {
	if target == nil || {{ .VariableName }}.FirstChild == nil {
		return false
	}

	for c := {{ .VariableName }}.FirstChild; c != nil; c = c.NextSibling {
		if c == target {
			return true
		}
	}

	return false
}

// delete_child is a helper function to delete the child from the children of the node.
//
// No nil nodes are returned.
//
// Parameters:
//   - target: The child to remove.
//
// Returns:
//   - []tr.Noder: A slice of pointers to the children of the node.
func ({{ .VariableName }} {{ .OutputType }}) delete_child(target {{ .OutputType}}) []tr.Noder {
	ok := {{ .VariableName }}.HasChild(target)
	if !ok {
		return nil
	}

	prev := target.PrevSibling
	next := target.NextSibling

	if prev != nil {
		prev.NextSibling = next
	}

	if next != nil {
		next.PrevSibling = prev
	}

	if target == {{ .VariableName }}.FirstChild {
		{{ .VariableName }}.FirstChild = next

		if next == nil {
			{{ .VariableName }}.LastChild = nil
		}
	} else if target == {{ .VariableName }}.LastChild {
		{{ .VariableName }}.LastChild = prev
	}

	target.Parent = nil
	target.PrevSibling = nil
	target.NextSibling = nil

	children := target.GetChildren()

	return children
}

// IsChildOf returns true if the node is a child of the parent. If target is nil,
// it returns false.
//
// Parameters:
//   - target: The target parent to check for.
//
// Returns:
//   - bool: True if the node is a child of the parent, false otherwise.
func ({{ .VariableName }} {{ .OutputType }}) IsChildOf(target {{ .OutputType }}) bool {
	if target == nil {
		return false
	}

	parents := target.GetAncestors()

	for node := {{ .VariableName }}; node.Parent != nil; node = node.Parent {
		parent := tr.Noder(node.Parent)

		ok := slices.Contains(parents, parent)
		if ok {
			return true
		}
	}

	return false
}

// delink_with_parent is a helper function to delink the parent with the children.
//
// Parameters:
//   - parent: The parent node. Must not be nil.
//   - children: The children nodes. Must not be nil and must be of type {{ .OutputType }}.
func delink_with_parent(parent {{ .OutputType }}, children []tr.Noder) {
	if len(children) == 0 {
		return
	}

	for _, child := range children {
		c := child.({{ .OutputType }})

		c.PrevSibling = nil
		c.NextSibling = nil
		c.Parent = nil
	}

	parent.FirstChild = nil
	parent.LastChild = nil
}
`
