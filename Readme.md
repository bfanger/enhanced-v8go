# Node Polyfills for [V8Go](https://github.com/rogchap/v8go)

## Install

```shell
go get -u github.com/bfanger/v8go-node-polyfills
```

> This module uses Golang [embed](https://golang.org/pkg/embed/), so requires Go version 1.16

## Polyfill List

None yet

## Usage

```go
    iso, _ := v8go.NewIsolate()
	polyfill, _ := node.NewPolyfill(iso)
	ctx, _ := v8go.NewContext(iso, polyfill.GlobalTemplate)
    ctx.RunScript("print('hello');", "go")
```

## Approach

### Just in Time (JIT)

The Polyfill.GlobalTemplate injects lightweight placeholders, when a polyfilled object or function is accessed for the first time the real polyfill is loaded and injected into the Context. After that they run "natively" in v8.

### Mostly Javascript

The polyfills are written in Javascript and use some Go api's that are exposed to the v8 runtime.

## Goals

Render a[SvelteKit](https://kit.svelte.dev/) pages.
