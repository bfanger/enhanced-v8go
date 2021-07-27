package node

import (
	"embed"
	"fmt"
	"strings"

	"rogchap.com/v8go"
)

//go:embed js/*.js
var js embed.FS

type Polyfill struct {
	GlobalTemplate *v8go.ObjectTemplate
}

func NewPolyfill(iso *v8go.Isolate) (*Polyfill, error) {

	// @todo ...ConfigOptions
	global, err := v8go.NewObjectTemplate(iso)
	if err != nil {
		return nil, err
	}

	println, err := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		var args []string
		for _, arg := range info.Args() {
			args = append(args, fmt.Sprintf("%v", arg))
		}
		fmt.Printf("%v\n", strings.Join(args, " "))
		return nil
	})
	if err != nil {
		return nil, err
	}
	global.Set("println", println)

	console, err := v8go.NewObjectTemplate(iso)
	if err != nil {
		return nil, err
	}
	console.Set("log", println)
	global.Set("console2", console)
	global.Set("console", console)

	err = lazyGlobalFunction(iso, global, "require", "js/require.js")
	if err != nil {
		return nil, err
	}
	return &Polyfill{GlobalTemplate: global}, nil
}

func lazyGlobalFunction(iso *v8go.Isolate, global *v8go.ObjectTemplate, name string, script string) error {
	template, err := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		source, err := js.ReadFile(script)
		if err != nil {
			panic(err)
		}
		val, err := info.Context().RunScript(string(source), script)
		if err != nil {
			panic(err)
		}
		fn, err := val.AsFunction()
		if err != nil {
			panic(err)
		}
		info.Context().Global().Set(name, fn)
		var args []v8go.Valuer
		for _, arg := range info.Args() {
			args = append(args, arg)
		}
		result, err := fn.Call(args...)
		if err != nil {
			panic(err) // @todo throw?
		}
		return result
	})
	if err != nil {
		return err
	}
	err = global.Set(name, template)
	return err
}
