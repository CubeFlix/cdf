// cmd/cdf/main.go
// CDF main executable tool.

package main

import (
	"fmt"
	"os"

	"github.com/cubeflix/cdf/export/html"
	"github.com/cubeflix/cdf/parser"
)

func main() {
	args := os.Args[:1]
	if len(args) != 2 {
		// Invalid number of arguments.
		fmt.Println("cdf is the Cubeflix Document Format converter for HTML")
		fmt.Println("usage: cdf input output")
		os.Exit(1)
	}

	input := args[0]
	output := args[1]
	data, err := os.ReadFile(input)
	if err != nil {
		fmt.Println("cdf:", err.Error())
		os.Exit(1)
	}
	outFile, err := os.Create(output)
	if err != nil {
		fmt.Println("cdf:", err.Error())
		os.Exit(1)
	}
	defer outFile.Close()

	p := parser.NewParser(data)
	err = p.Parse()
	if err != nil {
		fmt.Println("cdf:", err.Error())
		os.Exit(1)
	}
	h := html.NewHTMLExporter(outFile, html.HTMLSettings{})
	err = h.Export(&p.Tree)
	if err != nil {
		fmt.Println("cdf:", err.Error())
		os.Exit(1)
	}
}
