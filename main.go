// main.go
// testing

package main

import (
	"fmt"

	"github.com/cubeflix/cdf/parser"
)

var code string = `
[[cdf]]
[[p align=left]]ageni:
[[link dest=https://google.com]]hello paragraph[[/]] gensp gn[[/]]
[[/]]
`

func main() {
	p := parser.NewParser([]byte(code))
	err := p.Parse()
	if err != nil {
		panic(err)
	}
	fmt.Println(p.Tree.Content[0])
}
