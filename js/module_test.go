package js

import (
	"testing"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestTransformCJS(t *testing.T) {
	assert := assert.New(t)
	code, err := Transform("module.exports = 123;", api.TransformOptions{})
	if assert.NoError(err) == false {
		return
	}
	assert.Equal("(async function () { var exports = {}; var module = { exports };\nmodule.exports = 123;\n\nreturn module.exports; })()", code)
	ctx, err := NewContext()
	if assert.NoError(err) == false {
		return
	}
	m, err := newModule(ctx, code, "module_test.go")
	if assert.NoError(err) == false {
		return
	}
	defaultExport, err := m.Default()
	if assert.NoError(err) == false {
		return
	}
	assert.Equal(int64(123), defaultExport.Integer())
}

func TestTransformESM(t *testing.T) {
	assert := assert.New(t)
	code, err := Transform("export default 123", api.TransformOptions{})
	if assert.NoError(err) == false {
		return
	}
	ctx, err := NewContext()
	if assert.NoError(err) == false {
		return
	}
	m, err := newModule(ctx, code, "module_test.go")
	if assert.NoError(err) == false {
		return
	}
	defaultExport, err := m.Default()
	if assert.NoError(err) == false {
		return
	}
	assert.Equal(int64(123), defaultExport.Integer())
}

func TestTransformInvalidJavascript(t *testing.T) {
	_, err := Transform("oops!?!", api.TransformOptions{})
	assert.Error(t, err)
}

func TestRequire(t *testing.T) {
	assert := assert.New(t)
	ctx, err := NewContext()
	if assert.NoError(err) == false {
		return
	}
	exports, err := Require(ctx, "/Volumes/Sites/v8go-node-polyfills/node/js/console.js")
	if assert.NoError(err) == false {
		return
	}
	defaultExport, err := exports.Default()
	if assert.NoError(err) == false {
		return
	}
	console, err := defaultExport.AsObject()
	if assert.NoError(err) == false {
		return
	}
	consoleLog, err := console.Get("log")
	if assert.NoError(err) == false {
		return
	}
	assert.True(consoleLog.IsFunction())
}
