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
	code := "1 + 1"
	result := ctx.MustEval(code)
	fmt.Printf("%s = %s\n", code, result)
}
