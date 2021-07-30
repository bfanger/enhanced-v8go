package js

import (
	"encoding/json"
	"fmt"
	"os"

	"rogchap.com/v8go"
)

func Require(ctx *Context, filepath string) (*Module, error) {
	f, err := json.Marshal(filepath)
	if err != nil {
		return nil, err
	}
	v, err := ctx.RunScript(fmt.Sprintf("require(%s)", f), "js.go")
	if err != nil {
		return nil, err
	}
	return &Module{v}, nil
}
func requireTemplate(info *v8go.FunctionCallbackInfo) *Value {
	args := info.Args()
	if len(args) != 1 {
		return &Value{}
	}
	modulename := args[0].String()
	filename := modulename
	// @todo resolve module
	if modulename == "svelte/internal" {
		filename = "/Volumes/Sites/enhanced-v8go/examples/svelte/node_modules/svelte/internal/index.mjs"
	}

	b, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	m, err := NewModule(&Context{info.Context()}, string(b), filename)
	if err != nil {
		panic(err)
	}
	return m.Value
}
