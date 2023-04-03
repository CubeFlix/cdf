// main.go
// testing

package main

import (
	"net/http"

	"github.com/cubeflix/cdf/pages"
)

func main() {
	s, err := pages.LoadServer("testpages")
	if err != nil {
		panic(err)
	}
	http.ListenAndServe(":80", s)
}
