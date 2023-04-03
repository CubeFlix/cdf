// parser/block.go
// Parse an entire tag block.

package parser

import (
	"errors"
	"io"
	"strings"

	"github.com/cubeflix/cdf/ast"
)

// Parse the content in a block.
func (p *Parser) parseBlockContent() ([]ast.Block, error) {
	blocks := make([]ast.Block, 0)
	// Parse the inner blocks.
	for {
		// Skip whitespace.
		err := p.skipWhitespace()
		if err != nil {
			return nil, err
		}

		// Parse the block's tag.
		tag, err := p.parseTag()
		if err != nil {
			return nil, err
		}

		// Check for a closing tag.
		if tag.IsClosing {
			return blocks, nil
		}

		// Parse the inner block.
		if tag.Name == "p" {
			// Paragraph block.
			content, err := p.parseParagraphBlockContent()
			if err != nil {
				return nil, err
			}

			// Get the alignment information from the tag.
			alignment, err := p.tagGetAlignment(tag)
			if err != nil {
				return nil, err
			}

			blocks = append(blocks, &ast.Paragraph{
				BaseBlock: ast.BaseBlock{Alignment: alignment},
				Content:   content,
			})
		} else {
			// Invalid block type.
			return nil, errors.New("invalid block type")
		}
	}
}

// Parse the content in a paragraph block.
func (p *Parser) parseParagraphBlockContent() ([]ast.InlineBlock, error) {
	blocks := make([]ast.InlineBlock, 0)

	// Parse the block's content.
	var chunkLen int
	for {
		// Parse a single chunk of content.
		chunkLen = 0
		for {
			if p.cur+chunkLen+1 >= p.length {
				return nil, io.EOF
			}

			// Check for a '\'
			if p.data[p.cur+chunkLen] == '\\' {
				// Skip the next value.
				chunkLen += 2
				continue
			}

			// Check for a '[['.
			if p.data[p.cur+chunkLen] == '[' && p.data[p.cur+chunkLen+1] == '[' {
				// End the chunk.
				if chunkLen != 0 {
					chunk := p.data[p.cur : p.cur+chunkLen]
					p.cur += chunkLen
					blocks = append(blocks, escapeText(chunk))
				}

				// Parse the tag.
				tag, err := p.parseTag()
				if err != nil {
					return nil, err
				}

				// Check for a closing tag.
				if tag.IsClosing {
					return blocks, nil
				}

				// Parse the inner block.
				content, err := p.parseParagraphBlockContent()
				if err != nil {
					return nil, err
				}

				if tag.Name == "link" {
					// Hyperlink.

					// Get the hyperlink destination.
					if dest, ok := tag.Attributes["dest"]; ok {
						blocks = append(blocks, ast.HyperlinkBlock{
							BaseInlineBlock: ast.BaseInlineBlock{Content: content},
							Destination:     strings.ReplaceAll(dest, "\n", ""),
						})
					} else {
						return nil, errors.New("expected a 'dest' attribute")
					}
				} else if tag.Name == "b" {
					// Bold text.
					blocks = append(blocks, ast.FormattingBlock{BaseInlineBlock: ast.BaseInlineBlock{Content: content}, Attribute: ast.BoldFormatting})
				} else if tag.Name == "i" {
					// Italic text.
					blocks = append(blocks, ast.FormattingBlock{BaseInlineBlock: ast.BaseInlineBlock{Content: content}, Attribute: ast.ItalicFormatting})
				} else if tag.Name == "s" {
					// Strikethrough text.
					blocks = append(blocks, ast.FormattingBlock{BaseInlineBlock: ast.BaseInlineBlock{Content: content}, Attribute: ast.StrikethroughFormatting})
				} else if tag.Name == "u" {
					// Underline text.
					blocks = append(blocks, ast.FormattingBlock{BaseInlineBlock: ast.BaseInlineBlock{Content: content}, Attribute: ast.UnderlineFormatting})
				} else if tag.Name == "t" {
					// Teletype text.
					blocks = append(blocks, ast.FormattingBlock{BaseInlineBlock: ast.BaseInlineBlock{Content: content}, Attribute: ast.TeletypeFormatting})
				} else {
					return nil, errors.New("invalid tag type")
				}
				break
			} else {
				chunkLen++
			}
		}
	}
}
