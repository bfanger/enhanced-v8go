# Enhanced [V8Go](https://github.com/rogchap/v8go)

Build on top of https://github.com/rogchap/v8go and improves the developer experience:

## Features

- Module system
- shorter "js" package alias
- Typescript support (via esbuild)
- console.log (wip)
- require() (wip)

## Install

```shell
go get -u github.com/bfanger/enhanced-v8go
```

> This module uses Golang [embed](https://golang.org/pkg/embed/), so requires Go version 1.16

## Polyfills

- (wip) console

## Usage

```go
package main

import (
	"fmt"

	"github.com/bfanger/enhanced-v8go/js"
)

func main() {
	ctx, err := js.NewContext()
	if err != nil {
		panic(err)
	}
	result := ctx.MustEval("1 + 1")
	fmt.Println(result)
}
```

## Approach

### Just in Time (JIT)

The Polyfill.GlobalTemplate injects lightweight placeholders, when a polyfilled object or function is accessed for the first time the real polyfill is loaded and injected into the Context. After that they run "natively" in v8.

### Mostly Javascript

The polyfills are written in Javascript and use some Go api's that are exposed to the v8 runtime.

## Goals

Render a[SvelteKit](https://kit.svelte.dev/) pages.
