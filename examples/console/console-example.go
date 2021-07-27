package main

import (
	_ "embed"
	"fmt"

	"github.com/bfanger/v8go-node-polyfills/node"
	"rogchap.com/v8go"
)

//go:embed console-example.js
var source string

func main() {
	fmt.Println("> console-example")
	iso, err := v8go.NewIsolate()
	panicIfError(err)

	polyfill, err := node.NewPolyfill(iso)
	panicIfError(err)

	ctx, _ := v8go.NewContext(iso, polyfill.GlobalTemplate)

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
