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
	"log"
	"os"

	"github.com/PlayerR9/MyGoLib/cmdx/stack/pkg"
)

var (
	// Logger is the logger to use.
	Logger *log.Logger
)

func init() {
	Logger = log.New(os.Stdout, "[stack]: ", log.LstdFlags)
}

func main() {
	err := pkg.ParseFlags()
	if err != nil {
		Logger.Fatalf("Invalid flags: %s", err.Error())
	}

	g := pkg.NewGenerator()

	data, err := g.Generate()
	if err != nil {
		Logger.Fatalf("Could not generate code: %s", err.Error())
	}

	const (
		dest string = "stack.go"
	)

	err = os.WriteFile(dest, data, 0644)
	if err != nil {
		Logger.Fatalf("Could not write code: %s", err.Error())
	}
}
