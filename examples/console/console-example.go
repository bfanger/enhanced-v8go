package main

import (
	_ "embed"
	"fmt"

	"github.com/bfanger/enhanced-v8go/js"
)

//go:embed console-example.js
var source string

func main() {
	fmt.Println("> console-example")

	ctx, _ := js.NewContext()

	if source == "" {
		panic("source was empty")
	}

	val, err := ctx.RunScript(source, "console-example.go")
	panicIfError(err)
	fmt.Printf("Output: %+v\n\n", val)
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
