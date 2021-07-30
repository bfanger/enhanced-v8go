package js

import (
	"testing"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestTransformCJS(t *testing.T) {
	assert := assert.New(t)
	code, err := forceCJS("module.exports=  123", api.TransformOptions{})
	if assert.NoError(err) == false {
		return
	}
	assert.Equal("module.exports = 123;\n", code) // processed by esbuild to force cjs
	ctx, err := NewContext()
	if assert.NoError(err) == false {
		return
	}
	m, err := NewModule(ctx, code, "module_test.go")
	if assert.NoError(err) == false {
		return
	}
	defaultExport, err := m.Default()
	if assert.NoError(err) == false {
		return
	}
	assert.Equal(int64(123), defaultExport.Integer())
}
func TestTransformKnownCJS(t *testing.T) {
	code, err := forceCJS("module.exports=  123", api.TransformOptions{Sourcefile: "test.cjs"})
	if assert.NoError(t, err) {
		assert.Equal(t, "module.exports=  123", code) // not processed by esbuild
	}
}
func TestTransformESM(t *testing.T) {
	assert := assert.New(t)
	ctx, err := NewContext()
	if assert.NoError(err) == false {
		return
	}
	m, err := NewModule(ctx, "export default 123", "module_test.go")
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
	_, err := forceCJS("oops!?!", api.TransformOptions{})
	// allow esbuild to prevent executing "invalid" code
	assert.Error(t, err)
}

func TestTransformTypescript(t *testing.T) {
	_, err := forceCJS("const val: number = 123; export default n", api.TransformOptions{Sourcefile: "file.ts"})
	assert.NoError(t, err)
}
