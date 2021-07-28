package js

import (
	"fmt"
	"os"
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

type registeryEntry struct {
	m *Module
}

// @todo cleanup disposed contexts
var registery = make(map[*Context]map[string]*registeryEntry)

func Require(ctx *Context, filepath string) (*Module, error) {

	entry := registery[ctx][filepath]
	if entry != nil {
		if entry.m == nil {
			// @todo support for cyclic dependencies
			return nil, errors.Errorf("cyclic dependency for %s", filepath)
		}
		return entry.m, nil
	}
	if registery[ctx] == nil {
		registery[ctx] = make(map[string]*registeryEntry)
	}
	registery[ctx][filepath] = &registeryEntry{}
	b, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	code, err := Transform(string(b), api.TransformOptions{Sourcefile: filepath})
	if err != nil {
		return nil, err
	}
	m, err := newModule(ctx, code, filepath)
	if err != nil {
		return nil, err
	}
	registery[ctx][filepath].m = m
	return m, nil
}

func Transform(code string, options api.TransformOptions) (string, error) {
	options.Format = api.FormatCommonJS
	options.Target = api.ESNext
	// @todo Don't use esbuild when options.Sourcefile is a *.cjs file?

	result := api.Transform(code, options)
	if len(result.Errors) != 0 {
		formatted := api.FormatMessages(result.Errors, api.FormatMessagesOptions{TerminalWidth: 80, Kind: api.ErrorMessage, Color: true})
		return "", errors.New(strings.Join(formatted, "\n\n"))
	}
	// Wrap code in an async iife to allow top-level await
	return fmt.Sprintf("(async function () { var exports = {}; var module = { exports };\n%s\nreturn module.exports; })()", result.Code), nil
}

func newModule(ctx *Context, code, filename string) (*Module, error) {
	m, err := ctx.RunScript(code, filename)
	if err != nil {
		return nil, err
	}
	exports, err := Await(m)
	if err != nil {
		return nil, err
	}
	return &Module{exports}, nil
}
