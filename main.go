// main.go
// testing

package main

import "github.com/cubeflix/cdf/parser"

var code string = `
[[hello tag=12|a=1357 gnio afpi\|gr]]
`

func main() {
	p := parser.NewParser([]byte(code))
	err := p.Parse()
	if err != nil {
		panic(err)
	}
}
