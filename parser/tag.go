// parser/tag.go
// Parse an opening or closing tag.

package parser

import (
	"errors"
	"io"
)

// Parse an opening tag.
func (p *Parser) parseOpeningTag(data []byte) (int, error) {
	var i int
	

	// Expect the opening '[['.
	if len(data) < 2 {
		return i, io.EOF
	}
	i += 
	if data[i] != '[' || data[i] != '[' {
		// Invalid opening tags.
		return i, errors.New("expected '[[' for opening tag")
	}

	// Read each tag.
	for {
		// If we get a ']]', that means the opening tag is over.

	}
}
