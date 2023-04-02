// main.go
// testing

package main

import (
	"os"

	"github.com/cubeflix/cdf/export/html"
	"github.com/cubeflix/cdf/parser"
)

var code string = `
[[cdf title=hello]]
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
	h := html.NewHTMLExporter(os.Stdout, html.HTMLSettings{})
	err = h.Export(&p.Tree)
	if err != nil {
		panic(err)
	}
}
