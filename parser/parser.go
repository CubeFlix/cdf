// parser/parser.go
// Package parser provides an interface for parsing CDF files into abstract
// syntax trees.

package parser

import (
	"fmt"
	"io"

	"github.com/cubeflix/cdf/ast"
)

// The parser struct.
type Parser struct {
	settings Settings

	// The input data.
	data   []byte
	length int

	// The current position of the cursor.
	cur int

	// The output document AST.
	Tree ast.Document
}

// Create a new parser.
func NewParser(data []byte) *Parser {
	return &Parser{
		data:   data,
		length: len(data),
		Tree:   ast.Document{},
	}
}

// Parse the document.
func (p *Parser) Parse() error {
	// Skip whitespace.
	err := p.skipWhitespace()
	if err != nil {
		return err
	}

	// Read the file's header.
	// err = p.parseHeader()
	// if err != nil {
	// 	return err
	// }

	// Read the first tag.
	ot, err := p.parseOpeningTag()
	if err != nil {
		return err
	}
	fmt.Println(ot)
	return nil

	// for {
	// 	// Read each block.
	// 	err = p.parseBlock()
	// 	if err != nil {
	// 		return err
	// 	}
	// }
}

// Skip whitespace.
func (p *Parser) skipWhitespace() error {
	var b byte
	for {
		if p.cur >= p.length {
			// End of data.
			return io.EOF
		}
		b = p.data[p.cur]
		if b != ' ' && b != '\t' && b != '\n' && b != '\r' {
			return nil
		}
		p.cur++
	}
}

// Parse the header.
// func (p *Parser) parseHeader() error {
// 	if data[0] == '[' {
//
// 	}
// }
