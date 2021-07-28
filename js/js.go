package js

import "rogchap.com/v8go"

type Object = v8go.Object
type Value = v8go.Value
type ObjectTemplate = v8go.ObjectTemplate
type Isolate = v8go.Isolate

func NewIsolate() (*Isolate, error) {
	return v8go.NewIsolate()
}
