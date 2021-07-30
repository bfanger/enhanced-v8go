package main

import (
	"fmt"
	"os"

	"github.com/bfanger/enhanced-v8go/js"
)

func main() {
	cwd, err := os.Getwd()
	panicIfError(err)
	ctx, err := js.NewContext()
	panicIfError(err)
	m, err := js.Require(ctx, cwd+"/examples/svelte/hello-world.ssr.mjs")
	panicIfError(err)
	defaultExport, err := m.Default()
	panicIfError(err)
	component, err := defaultExport.AsObject()
	panicIfError(err)
	render, err := component.Get("render")
	panicIfError(err)
	renderFn, err := render.AsFunction()
	panicIfError(err)
	result, err := renderFn.Call()
	panicIfError(err)
	resultObj, err := result.AsObject()
	panicIfError(err)
	html, err := resultObj.Get("html")
	fmt.Println(html)
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
