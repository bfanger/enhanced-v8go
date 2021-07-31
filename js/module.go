package js

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/pkg/errors"
)

type Module struct {
	*Value
}

func (m *Module) Default() (*Value, error) {
	exports, err := m.AsObject()
	if err != nil {
		return m.Value, nil // CommonJS
	}
	esm, err := exports.Get("__esModule")
	if err != nil {
		return nil, err
	}
	if esm.IsBoolean() && esm.Boolean() {
		return exports.Get("default")
	}
	return m.Value, nil // CommonJS
}

func NewModule(ctx *Context, code string, origin string) (*Module, error) {
	// Wrap the commonjs code in an async iife exposing the exports as a promise
	filename, err := json.Marshal(origin)
	if err != nil {
		return nil, err
	}
	dirname, err := json.Marshal(path.Dir(origin))
	if err != nil {
		return nil, err
	}
	code, err = forceCJS(code, api.TransformOptions{Sourcefile: origin})
	if err != nil {
		return nil, err
	}
	augmented := fmt.Sprintf(`
const __filename = %s;
const __dirname = %s;
var module =  {
	 exports: {}, 
	 id: __filename,
	 filename: __filename,
	 path: __dirname,
	 loaded: false,
};
var exports = module.exports;
require.main = module;
`, filename, dirname)
	oneliner := strings.Join(strings.Split(augmented, "\n"), " ")
	wrapped := fmt.Sprintf("(async function () { %s %s\nmodule.loaded = true; return module.exports; })()", oneliner, code)
	m, err := ctx.RunScript(wrapped, origin)
	if err != nil {
		return nil, err
	}
	exports, err := Await(m)
	if err != nil {
		return nil, err
	}
	return &Module{exports}, nil
}
func forceCJS(code string, options api.TransformOptions) (string, error) {
	if strings.HasSuffix(options.Sourcefile, ".cjs") {
		// Already commonjs? skip esbuild
		return code, nil
	}
	if strings.HasSuffix(options.Sourcefile, ".ts") {
		options.Loader = api.LoaderTS
	}
	if strings.HasSuffix(options.Sourcefile, ".tsx") {
		options.Loader = api.LoaderTSX
	}
	if strings.HasSuffix(options.Sourcefile, ".json") {
		options.Loader = api.LoaderJSON
	}
	options.Format = api.FormatCommonJS
	options.Target = api.ESNext
	options.Engines = []api.Engine{
		{Name: api.EngineChrome, Version: "92"}, // @todo auto update to latest?
	}
	result := api.Transform(code, options)
	if len(result.Errors) != 0 {
		formatted := api.FormatMessages(result.Errors, api.FormatMessagesOptions{TerminalWidth: 80, Kind: api.ErrorMessage, Color: true})
		return "", errors.New(strings.Join(formatted, "\n\n"))
	}
	// Wrap code in an async iife to allow top-level await
	return string(result.Code), nil
}
