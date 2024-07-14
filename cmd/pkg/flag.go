package pkg

import (
	"flag"

	ggen "github.com/PlayerR9/MyGoLib/go_generator"
)

func init() {
	ggen.SetOutputFlag("<type>_stack.go", false)
}

func ParseFlags() string {
	flag.Parse()

	return ""
}
