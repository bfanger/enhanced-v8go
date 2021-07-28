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
