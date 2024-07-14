package pkg

import (
	"flag"

	ggen "github.com/PlayerR9/MyGoLib/go_generator"
)

var (
	TypeNameFlag *string
)

func init() {
	TypeNameFlag = flag.String("name", "", "The name of the generated type. It must be set.")

	ggen.SetOutputFlag("<type>_table.go", false)

	ggen.SetTypeListValFlag("type", true, 1, "The type of each table's cell.")
	ggen.SetGenericsSignFlag("g", false, 1)
}

// GenData is the data struct for the generator.
type GenData struct {
	// PackageName is the package name of the generated code.
	PackageName string

	// TypeName is the name of the generated type.
	TypeName string

	// TypeSig is the signature of the generated type.
	TypeSig string

	// GenericsSign is the signature of the generics.
	GenericsSign string

	// CellType is the type of each table's cell.
	CellType string
}

// SetPackageName implements the go_generator.Generater interface.
func (g GenData) SetPackageName(pkg_name string) ggen.Generater {
	g.PackageName = pkg_name

	return g
}
