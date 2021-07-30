package js

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
