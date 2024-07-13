package main

import (
	"log"
	"os"

	"github.com/PlayerR9/MyGoLib/cmd/stack/pkg"
)

var (
	Logger *log.Logger
)

func init() {
	Logger = log.New(os.Stdout, "[stack]: ", log.LstdFlags)
}

func main() {
	g := pkg.NewGenerator("int")

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
