// parser/block.go
// Parse an entire tag block.

package parser

import (
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/cubeflix/cdf/ast"
	"gopkg.in/go-playground/colors.v1"
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

		// Get the alignment information from the tag.
		alignment, err := p.tagGetAlignment(tag)
		if err != nil {
			return nil, err
		}

		// Check if we should wrap the block.
		_, shouldWrap := tag.Attributes["wrap"]

		// Parse the inner block.
		if tag.Name == "p" {
			// Paragraph block.
			content, err := p.parseParagraphBlockContent()
			if err != nil {
				return nil, err
			}

			blocks = append(blocks, &ast.Paragraph{
				BaseBlock: ast.BaseBlock{Alignment: alignment, Wrap: shouldWrap},
				Content:   content,
			})
		} else if tag.Name == "block" {
			// Basic block.
			content, err := p.parseBlockContent()
			if err != nil {
				return nil, err
			}

			blocks = append(blocks, &ast.BasicBlock{
				BaseBlock: ast.BaseBlock{Alignment: alignment, Wrap: shouldWrap},
				Content:   content,
			})
		} else if tag.Name == "quote" {
			// Quote block.
			content, err := p.parseBlockContent()
			if err != nil {
				return nil, err
			}

			blocks = append(blocks, &ast.Quote{
				BaseBlock: ast.BaseBlock{Alignment: alignment, Wrap: shouldWrap},
				Content:   content,
			})
		} else if tag.Name == "image" {
			// Image block.

			// Get the source.
			imgSrc, ok := tag.Attributes["src"]
			if !ok {
				return nil, errors.New("'image' tag expects a 'src' attribute")
			}

			// Check if we should include the caption.
			_, hasCaption := tag.Attributes["has-caption"]

			a, b, c, d, e, f, err := p.generateImageBlock(tag)
			if err != nil {
				return nil, err
			}

			content, err := p.parseParagraphBlockContent()
			if err != nil {
				return nil, err
			}

			blocks = append(blocks, &ast.Image{
				BaseBlock:          ast.BaseBlock{Alignment: alignment, Wrap: shouldWrap},
				Source:             strings.ReplaceAll(imgSrc, "\n", ""),
				HasCaption:         hasCaption,
				Caption:            content,
				HasWidthParameter:  a,
				WidthValue:         b,
				WidthType:          c,
				HasHeightParameter: d,
				HeightValue:        e,
				HeightType:         f,
			})
		} else if tag.Name == "h" {
			// Heading block.
			class, err := p.generateHeadingClass(tag)
			if err != nil {
				return nil, err
			}
			content, err := p.parseParagraphBlockContent()
			if err != nil {
				return nil, err
			}

			blocks = append(blocks, &ast.Heading{
				BaseBlock: ast.BaseBlock{Alignment: alignment, Wrap: shouldWrap},
				Content:   content,
				Class:     class,
			})
		} else if tag.Name == "hr" {
			// Horizontal rule.
			_, err := p.parseParagraphBlockContent()
			if err != nil {
				return nil, err
			}
			blocks = append(blocks, &ast.HorizontalRule{
				BaseBlock: ast.BaseBlock{Alignment: alignment, Wrap: shouldWrap},
			})
		} else if tag.Name == "list" {
			// List block.
			content, err := p.parseBlockContent()
			if err != nil {
				return nil, err
			}
			_, isOrdered := tag.Attributes["ordered"]

			blocks = append(blocks, &ast.List{
				BaseBlock: ast.BaseBlock{Alignment: alignment, Wrap: shouldWrap},
				Items:     content,
				Ordered:   isOrdered,
			})
		} else if tag.Name == "table" {
			// Table block.
			table, err := p.parseTableContent()
			if err != nil {
				return nil, err
			}
			blocks = append(blocks, &ast.Table{
				BaseBlock: ast.BaseBlock{Alignment: alignment, Wrap: shouldWrap},
				Rows:      table,
			})
		} else if tag.Name == "collapse" {
			// Collapseable block.
			summaryContent, innerContent, err := p.parseCollapse()
			if err != nil {
				return nil, err
			}
			blocks = append(blocks, &ast.Collapse{
				BaseBlock: ast.BaseBlock{Alignment: alignment, Wrap: shouldWrap},
				Summary:   summaryContent,
				Content:   innerContent,
			})
		} else if tag.Name == "break" {
			// Page break.
			_, err := p.parseParagraphBlockContent()
			if err != nil {
				return nil, err
			}
			blocks = append(blocks, &ast.PageBreak{})
		} else {
			// Invalid block type.
			return nil, errors.New("invalid block type")
		}
	}
}

// Parse a table's content.
func (p *Parser) parseTableContent() ([]ast.TableRow, error) {
	// Parse the content in a table.
	rows := make([]ast.TableRow, 0)
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
			return rows, nil
		}

		if tag.Name != "row" {
			return nil, errors.New("table should only contain rows")
		}

		// Parse the inner cells.
		row, err := p.parseTableRow()
		if err != nil {
			return nil, err
		}
		rows = append(rows, ast.TableRow{Cells: row})
	}
}

// Parse a table row's cells.
func (p *Parser) parseTableRow() ([]ast.TableCell, error) {
	// Parse the content in a row.
	cells := make([]ast.TableCell, 0)

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
			return cells, nil
		}

		if tag.Name != "cell" {
			return nil, errors.New("row should only contain cells")
		}

		_, isHeader := tag.Attributes["is-header"]

		// Parse the inner content.
		content, err := p.parseBlockContent()
		if err != nil {
			return nil, err
		}
		cells = append(cells, ast.TableCell{
			Content:  content,
			IsHeader: isHeader,
		})
	}
}

// Parse the collapseable block. Returns the summary and content.
func (p *Parser) parseCollapse() ([]ast.InlineBlock, []ast.Block, error) {
	// Skip whitespace.
	err := p.skipWhitespace()
	if err != nil {
		return nil, nil, err
	}

	// Parse the summary tag.
	summaryTag, err := p.parseTag()
	if err != nil {
		return nil, nil, err
	}
	if summaryTag.Name != "summary" {
		return nil, nil, errors.New("collapseable block should contain a summary block")
	}
	summaryContent, err := p.parseParagraphBlockContent()
	if err != nil {
		return nil, nil, err
	}

	err = p.skipWhitespace()
	if err != nil {
		return nil, nil, err
	}

	// Parse the content tag.
	contentTag, err := p.parseTag()
	if err != nil {
		return nil, nil, err
	}
	if contentTag.Name != "content" {
		return nil, nil, errors.New("collapseable block should contain a content block")
	}
	content, err := p.parseBlockContent()
	if err != nil {
		return nil, nil, err
	}

	err = p.skipWhitespace()
	if err != nil {
		return nil, nil, err
	}

	// Closing tag.
	closingTag, err := p.parseTag()
	if err != nil {
		return nil, nil, err
	}
	if !closingTag.IsClosing {
		return nil, nil, errors.New("expected closing tag")
	}

	return summaryContent, content, nil
}

// Get the heading class given a heading tag.
func (p *Parser) generateHeadingClass(t tagItem) (ast.HeadingType, error) {
	if class, ok := t.Attributes["c"]; ok {
		switch class {
		case "1":
			return ast.Heading1Type, nil
		case "2":
			return ast.Heading2Type, nil
		case "3":
			return ast.Heading3Type, nil
		case "4":
			return ast.Heading4Type, nil
		case "5":
			return ast.Heading5Type, nil
		default:
			return 0, errors.New("'h' tag expected a valid class (1-5)")
		}
	}
	return 0, errors.New("'h' tag expected a class ('c')")
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
						return nil, errors.New("'link' tag expected a 'dest' attribute")
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
				} else if tag.Name == "size" {
					// Size block.

					// Get the size information.
					sizeVal, sizeType, err := p.generateSizeBlock(tag)
					if err != nil {
						return nil, err
					}
					blocks = append(blocks, ast.SizeBlock{BaseInlineBlock: ast.BaseInlineBlock{Content: content}, Value: sizeVal, Type: sizeType})
				} else if tag.Name == "font" {
					// Font block.

					// Get the font family.
					if family, ok := tag.Attributes["family"]; ok {
						blocks = append(blocks, ast.FontBlock{
							BaseInlineBlock: ast.BaseInlineBlock{Content: content},
							Family:          strings.ReplaceAll(family, "\n", ""),
						})
					} else {
						return nil, errors.New("'font' tag expected a 'family' attribute")
					}
				} else if tag.Name == "color" {
					// Color block.

					// Get the color information.
					fore, back, err := p.generateColorBlock(tag)
					if err != nil {
						return nil, err
					}
					blocks = append(blocks, ast.ColorBlock{BaseInlineBlock: ast.BaseInlineBlock{Content: content}, ForegroundValue: fore, BackgroundValue: back})
				} else if tag.Name == "inline-image" {
					// Inline image block.
					// Get the source.
					imgSrc, ok := tag.Attributes["src"]
					if !ok {
						return nil, errors.New("'inline-image' tag expects a 'src' attribute")
					}

					a, b, c, d, e, f, err := p.generateImageBlock(tag)
					if err != nil {
						return nil, err
					}

					blocks = append(blocks, ast.InlineImageBlock{
						BaseInlineBlock:    ast.BaseInlineBlock{Content: nil},
						Source:             strings.ReplaceAll(imgSrc, "\n", ""),
						HasWidthParameter:  a,
						WidthValue:         b,
						WidthType:          c,
						HasHeightParameter: d,
						HeightValue:        e,
						HeightType:         f,
					})
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

// Get the size block information from a tag. Returns the size value and the
// size unit.
func (p *Parser) generateSizeBlock(t tagItem) (float32, ast.SizeType, error) {
	// Ensure that there is only one attribute in the tag.
	if len(t.Attributes) != 1 {
		return 0, 0, errors.New("'size' tag should contain one parameter")
	}

	if val, ok := t.Attributes["percent"]; ok {
		// Percentage.
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return 0, 0, err
		}
		return float32(floatVal), ast.PercentageSizeType, nil
	} else if val, ok := t.Attributes["px"]; ok {
		// Pixels.
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return 0, 0, err
		}
		return float32(floatVal), ast.PixelSizeType, nil
	} else if val, ok := t.Attributes["pt"]; ok {
		// Points.
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return 0, 0, err
		}
		return float32(floatVal), ast.PointSizeType, nil
	} else if val, ok := t.Attributes["cm"]; ok {
		// Centimeters.
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return 0, 0, err
		}
		return float32(floatVal), ast.CentimeterSizeType, nil
	} else if val, ok := t.Attributes["mm"]; ok {
		// Millimeters.
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return 0, 0, err
		}
		return float32(floatVal), ast.MillimeterSizeType, nil
	} else {
		// Invalid tag parameters.
		return 0, 0, errors.New("'size' tag should contain a valid size parameter ('percent', 'px', 'pt', 'cm', 'mm')")
	}
}

// Get the color block information from a tag. Returns the foreground and
// background color values. Values are nil if unspecified.
func (p *Parser) generateColorBlock(t tagItem) (colors.Color, colors.Color, error) {
	// Ensure that there is only either one or two parameters.
	if len(t.Attributes) != 1 && len(t.Attributes) != 2 {
		return nil, nil, errors.New("'color' tag should contain either one or two parameters ('fg', 'bg')")
	}

	var fore, back colors.Color
	var err error

	if val, ok := t.Attributes["fg"]; ok {
		// Foreground color specified.
		fore, err = colors.Parse(val)
		if err != nil {
			return nil, nil, err
		}
	}
	if val, ok := t.Attributes["bg"]; ok {
		// Background color specified.
		back, err = colors.Parse(val)
		if err != nil {
			return nil, nil, err
		}
	}
	if fore == nil && back == nil {
		// Invalid tag parameters.
		return nil, nil, errors.New("'color' tag should contain a valid color parameter ('fg', 'bg')")
	}

	return fore, back, nil
}

// Get the size information for the image block. Returns the width size
// information and the height size information.
func (p *Parser) generateImageBlock(t tagItem) (bool, float32, ast.SizeType, bool, float32, ast.SizeType, error) {
	var hasWidth, hasHeight bool
	var widthVal, heightVal float32
	var widthType, heightType ast.SizeType

	if val, ok := t.Attributes["width-percent"]; ok {
		// Percentage.
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return false, 0, 0, false, 0, 0, err
		}
		hasWidth = true
		widthVal = float32(floatVal)
		widthType = ast.PercentageSizeType
	} else if val, ok := t.Attributes["width-px"]; ok {
		// Pixels.
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return false, 0, 0, false, 0, 0, err
		}
		hasWidth = true
		widthVal = float32(floatVal)
		widthType = ast.PixelSizeType
	} else if val, ok := t.Attributes["width-pt"]; ok {
		// Points.
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return false, 0, 0, false, 0, 0, err
		}
		hasWidth = true
		widthVal = float32(floatVal)
		widthType = ast.PointSizeType
	} else if val, ok := t.Attributes["width-cm"]; ok {
		// Centimeters.
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return false, 0, 0, false, 0, 0, err
		}
		hasWidth = true
		widthVal = float32(floatVal)
		widthType = ast.CentimeterSizeType
	} else if val, ok := t.Attributes["width-mm"]; ok {
		// Millimeters.
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return false, 0, 0, false, 0, 0, err
		}
		hasWidth = true
		widthVal = float32(floatVal)
		widthType = ast.MillimeterSizeType
	}

	if val, ok := t.Attributes["height-percent"]; ok {
		// Percentage.
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return false, 0, 0, false, 0, 0, err
		}
		hasHeight = true
		heightVal = float32(floatVal)
		heightType = ast.PercentageSizeType
	} else if val, ok := t.Attributes["height-px"]; ok {
		// Pixels.
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return false, 0, 0, false, 0, 0, err
		}
		hasHeight = true
		heightVal = float32(floatVal)
		heightType = ast.PixelSizeType
	} else if val, ok := t.Attributes["height-pt"]; ok {
		// Points.
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return false, 0, 0, false, 0, 0, err
		}
		hasHeight = true
		heightVal = float32(floatVal)
		heightType = ast.PointSizeType
	} else if val, ok := t.Attributes["height-cm"]; ok {
		// Centimeters.
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return false, 0, 0, false, 0, 0, err
		}
		hasHeight = true
		heightVal = float32(floatVal)
		heightType = ast.CentimeterSizeType
	} else if val, ok := t.Attributes["height-mm"]; ok {
		// Millimeters.
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return false, 0, 0, false, 0, 0, err
		}
		hasHeight = true
		heightVal = float32(floatVal)
		heightType = ast.MillimeterSizeType
	}

	return hasWidth, widthVal, widthType, hasHeight, heightVal, heightType, nil
}
