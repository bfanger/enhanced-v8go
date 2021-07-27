package promise

import (
	"errors"
	"sync"

	"rogchap.com/v8go"
)

func Await(promise *v8go.Value) (*v8go.Value, error) {
	p, err := promise.AsPromise()
	if err != nil {
		return nil, err
	}
	if p.State() == v8go.Pending {
		wg := sync.WaitGroup{}
		wg.Add(1)
		p.Then(func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			wg.Done()
			return nil
		}, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			wg.Done()
			return nil
		})
		wg.Wait()
	}
	r := p.Result()
	if p.State() == v8go.Rejected {
		if r.IsNativeError() || r.IsString() {
			return r, errors.New(r.String())
		}
		return r, errors.New("promise was rejected")
	}
	return r, nil
}
