// parser/tag.go
// Parse an opening or closing tag.

package parser

import (
	"errors"
	"fmt"
	"io"
	"unicode"
)

// An opening tag item.
type tagItem struct {
	// The tag type/name.
	Name string

	// The tag's attributes.
	Attributes map[string]string
}

// Parse an opening tag item.
func (p *Parser) parseOpeningTag() (tagItem, error) {
	// Expect the opening '[['.
	if p.cur+1 >= p.length {
		return tagItem{}, io.EOF
	}
	if p.data[p.cur] != '[' || p.data[p.cur+1] != '[' {
		// Invalid opening tags.
		return tagItem{}, errors.New("expected '[[' for opening tag")
	}
	p.cur += 2

	// Skip whitespace.
	err := p.skipWhitespace()
	if err != nil {
		return tagItem{}, err
	}

	// Parse the tag name.
	var nameLen int
	for {
		if p.cur+nameLen >= p.length {
			return tagItem{}, io.EOF
		}
		if !unicode.IsLetter(rune(p.data[p.cur+nameLen])) {
			// End of name.
			break
		}
		nameLen++
	}
	name := p.data[p.cur : p.cur+nameLen]
	p.cur += nameLen

	// Skip whitespace.
	err = p.skipWhitespace()
	if err != nil {
		return tagItem{}, err
	}

	// If we get a ']]', that means the opening tag is over.
	if p.cur+1 >= p.length {
		return tagItem{}, io.EOF
	}
	if p.data[p.cur] == ']' && p.data[p.cur+1] == ']' {
		// Close the tag item.
		p.cur += 2
		return tagItem{
			Name:       string(name),
			Attributes: map[string]string{},
		}, nil
	}

	// Read each tag attribute.
	attributes := map[string]string{}
	for {
		attrName, attrValue, expectMore, err := p.parseOpeningTagAttribute()
		if err != nil {
			return tagItem{}, err
		}
		fmt.Println(attrName, attrValue, expectMore)
		attributes[attrName] = attrValue

		if !expectMore {
			// Close the tag item.
			break
		}

		// Skip whitespace.
		err = p.skipWhitespace()
		if err != nil {
			return tagItem{}, err
		}
	}

	return tagItem{
		Name:       string(name),
		Attributes: attributes,
	}, nil
}

// Parse a single tag attribute in an opening tag item. Returns the attribute
// name and value, along with if the parser should expect another attribute.
func (p *Parser) parseOpeningTagAttribute() (string, string, bool, error) {
	// Parse the attribute name.
	var nameLen int
	for {
		if p.cur+nameLen >= p.length {
			return "", "", false, io.EOF
		}
		if !unicode.IsLetter(rune(p.data[p.cur+nameLen])) {
			// End of name.
			break
		}
		nameLen++
	}
	name := p.data[p.cur : p.cur+nameLen]
	p.cur += nameLen

	// Skip whitespace.
	err := p.skipWhitespace()
	if err != nil {
		return "", "", false, err
	}

	// Expect an '='.
	if p.cur >= p.length {
		return "", "", false, io.EOF
	}
	if p.data[p.cur] != '=' {
		return "", "", false, errors.New("expected '='")
	}
	p.cur++

	// Skip whitespace.
	err = p.skipWhitespace()
	if err != nil {
		return "", "", false, err
	}

	// Parse the attribute value.
	var valueLen int
	for {
		if p.cur+valueLen >= p.length {
			return "", "", false, io.EOF
		}

		// Check for a '\'
		if p.data[p.cur+valueLen] == '\\' {
			// Skip the next value.
			valueLen += 2
			continue
		}

		// TODO: NEED TO CLEAN/ESCAPE THE FINAL VALUE

		// Check for a ']]'.
		if p.cur+valueLen+1 >= p.length {
			return "", "", false, io.EOF
		}
		if p.data[p.cur+valueLen] == ']' && p.data[p.cur+valueLen+1] == ']' {
			// End the value.
			value := p.data[p.cur : p.cur+valueLen]
			p.cur += valueLen + 2

			return string(name), string(value), false, nil
		}

		// Check for a '|'.
		if p.data[p.cur+valueLen] == '|' {
			// End the value.
			value := p.data[p.cur : p.cur+valueLen]
			p.cur += valueLen + 1

			return string(name), string(value), true, nil
		}

		// Make sure the value is printable.
		if !unicode.IsPrint(rune(p.data[p.cur+valueLen])) {
			return "", "", false, errors.New("non-printable character")
		}

		valueLen++
	}
}
