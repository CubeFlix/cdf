// parser/parser.go
// Package parser provides an interface for parsing CDF files into abstract
// syntax trees.

package parser

import (
	"io"

	"github.com/cubeflix/cdf/ast"
)

// The parser struct.
type Parser struct {
	settings Settings

	// The output document AST.
	Tree ast.Document
}

// Create a new parser.
func NewParser() *Parser {
	return &Parser{
		Tree: ast.Document{},
	}
}

// Parse the document. Takes the document as data, and returns the index of the
// end of the document.
func (p *Parser) Parse(data []byte) (int, error) {
	var temp, i int

	// Skip whitespace.
	i, err := p.skipWhitespace(data)
	if err != nil {
		return i, err
	}

	// Read the file's metadata header.
	temp, err = p.parseHeader(data[i:])
	if err != nil {
		return i, err
	}
	i += temp

	for {
		// Read each block.
		temp, err = p.parseBlock(data[i:])
		if err != nil {
			return i, err
		}
		i += temp
	}
}

// Skip whitespace. Returns the length of the whitespace.
func (p *Parser) skipWhitespace(data []byte) (int, error) {
	var i int
	for {
		if i >= len(data) {
			// End of data.
			return i, io.EOF
		}
		b := data[i]
		if b == ' ' || b == '\t' || b == '\n' || b == '\r' {
			return i, io.EOF
		}
		i++
	}

	return i, nil
}

// Parse the header. Returns the length of the header.
func (p *Parser) parseHeader(data []byte) (int, error) {
	if data[0] == '[' {

	}
}
