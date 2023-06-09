// parser/parser.go
// Package parser provides an interface for parsing CDF files into abstract
// syntax trees.

package parser

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"

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

	// Read the opening tag.
	openingTag, err := p.parseTag()
	if err != nil {
		return err
	}
	if openingTag.IsClosing {
		return errors.New("expected an opening tag")
	}
	if openingTag.Name != "cdf" {
		return errors.New("expected a 'cdf' tag")
	}
	p.tagGetDocumentFields(openingTag)

	// Read the block's contents.
	content, err := p.parseBlockContent()
	if err != nil {
		return err
	}

	p.Tree.Content = content

	return nil
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

// Clean/escape a block of text.
func escapeText(data []byte) string {
	r := bufio.NewReader(bytes.NewReader(data))
	out := strings.Builder{}
	for {
		chunk, err := r.ReadString('\\')
		if err != nil {
			out.WriteString(chunk)
			return out.String()
		}
		out.WriteString(chunk[:len(chunk)-1])
		b, err := r.ReadByte()
		if err != nil {
			return out.String()
		}
		out.WriteByte(b)
	}
}
