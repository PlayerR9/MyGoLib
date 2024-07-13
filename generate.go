//go:generate go run cmd/treenode/main.go -type=TreeNode -fields=a/MyType[T] -g=T/any -output=generic_treenode.go
package MyGoLib

type MyType[T any] struct {
	a T
}
