package js

import (
	"encoding/json"
	"errors"
	"fmt"

	"go.kuoruan.net/v8go-polyfills/console"
	"rogchap.com/v8go"
)

type Context struct {
	*v8go.Context
}

func NewContext(opt ...v8go.ContextOption) (*Context, error) {
	var isolate *Isolate
	for _, o := range opt {
		iso, ok := o.(*Isolate)
		if ok {
			isolate = iso
			break
		}
		if _, ok := o.(*ObjectTemplate); ok == true {
			return nil, errors.New("(multiple) object templates not implemented") // @todo apply polyfills to this template?
		}
	}
	if isolate == nil {
		var err error
		isolate, err = v8go.NewIsolate()
		if err != nil {
			return nil, err
		}
		opt = append(opt, isolate)
	}
	polyfill, err := NewPolyfill(isolate)
	if err != nil {
		return nil, err
	}
	opt = append(opt, polyfill.GlobalTemplate)
	ctx, err := v8go.NewContext(opt...)

	console.InjectTo(ctx)
	return &Context{Context: ctx}, err
}

func (ctx *Context) Eval(code string) (*Value, error) {
	return ctx.RunScript(code, "eval")
}
func (ctx *Context) MustEval(code string) *Value {
	val, err := ctx.Eval(code)
	if err != nil {
		panic(err)
	}
	return val
}
func (ctx *Context) Error(message string) *Value {
	m, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	return ctx.MustEval(fmt.Sprintf("new Error(%s)", m))
}

func (ctx *Context) dataErrorTuple(data *v8go.Value, err error) *Value {
	var tuple *Value
	if err != nil {
		message, err := json.Marshal(err.Error())
		if err != nil {
			panic(err)
		}
		tuple = ctx.MustEval(fmt.Sprintf("([ null, new Error(%s) ])", message))
	} else {
		tuple = ctx.MustEval("([ null, null ])")
	}
	if data != nil {
		val, err := tuple.AsObject()
		if err != nil {
			panic(nil)
		}
		val.SetIdx(0, data)
	}
	return tuple
}
