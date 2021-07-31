package js

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"rogchap.com/v8go"
)

func Require(ctx *Context, filepath string) (*Module, error) {
	f, err := json.Marshal(filepath)
	if err != nil {
		return nil, err
	}
	v, err := ctx.RunScript(fmt.Sprintf("require(%s)", f), "require.go")
	if err != nil {
		return nil, err
	}
	return &Module{v}, nil
}

func requireFileTemplate(info *v8go.FunctionCallbackInfo) *Value {
	args := info.Args()
	if len(args) != 1 {
		return &Value{}
	}
	filename := args[0].String()
	b, err := os.ReadFile(filename)
	ctx := &Context{info.Context()}
	if err != nil {
		return ctx.dataErrorTuple(nil, err)
	}
	m, err := NewModule(ctx, string(b), filename)
	if err != nil {
		return ctx.dataErrorTuple(nil, err)
	}
	return ctx.dataErrorTuple(m.Value, err)
}

var extensions = []string{".cjs", ".js", ".mjs", ".ts", ".tsx"}

func Resolve(id string, origin string) (string, error) {
	if strings.HasPrefix(id, "/") || strings.HasPrefix(id, "./") || strings.HasPrefix(id, "../") {
		if strings.HasPrefix(id, "/") == false {
			id = path.Join(path.Dir(origin), id)
		}
		info, err := os.Stat(id)
		if err == nil {
			if info.IsDir() {
				return Resolve(id+"/index", "")
			}
			return id, nil // id was a path
		}
		for _, ext := range extensions {
			_, err := os.Stat(id + ext)
			if err == nil {
				return id + ext, nil // id was a path without extension
			}
		}
		return "", errors.Errorf("file '%s' not found", id)
	}
	module := id
	pos := strings.Index(id, "/")
	if pos != -1 {
		module = id[0:pos]
	}
	p, err := resolvePackage(module, origin)
	if err != nil {
		return "", err
	}
	if pos == -1 {
		// require main entrypoint
		mainFile := p.Pkg["main"].(string)
		if mainFile == "" {
			mainFile = "./index"
		}
		fp := path.Join(p.Dir, mainFile)
		return Resolve(fp, "")
	}
	// @todo check package.json[exports]
	submodule := path.Join(p.Dir, id[pos:])
	return Resolve(submodule, "")
}

type packageJSON struct {
	Dir string
	Pkg map[string]interface{}
}

func resolvePackage(module, origin string) (*packageJSON, error) {
	dir := origin
	p := ""
	success := false
	for {
		dir = path.Dir(dir)
		if len(dir) < 4 {
			break
		}
		p = dir + "/node_modules/" + module + "/package.json"
		_, err := os.Stat(p)
		if err == nil {
			success = true
			break
		}
	}
	if success == false {
		return nil, errors.Errorf("module '%s' not found", module)
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	data := &packageJSON{
		Dir: path.Dir(p),
	}
	if err := json.Unmarshal(b, &data.Pkg); err != nil {
		return nil, err
	}
	return data, nil
}

func requireResolveTemplate(info *v8go.FunctionCallbackInfo) *Value {
	args := info.Args()
	ctx := &Context{info.Context()}
	if len(args) != 2 {
		return ctx.dataErrorTuple(nil, errors.New("2 parameters are exepected"))
	}
	f, err := Resolve(args[0].String(), args[1].String())
	if err != nil {
		return ctx.dataErrorTuple(nil, err)
	}
	filepath, err := json.Marshal(f)
	if err != nil {
		panic(err)
	}
	val := ctx.MustEval(fmt.Sprintf("(%s)", filepath))
	return ctx.dataErrorTuple(val, nil)
}
