// export/html/html.go
// Package html provides functionality for exporting into HTML.

package html

import (
	"errors"
	"fmt"
	"html"
	"io"
	"strings"

	"github.com/cubeflix/cdf/ast"
)

// HTML exporter.
type HTMLExporter struct {
	stream   io.Writer
	settings HTMLSettings
}

// Create a new HTML exporter.
func NewHTMLExporter(stream io.Writer, settings HTMLSettings) *HTMLExporter {
	return &HTMLExporter{
		stream:   stream,
		settings: settings,
	}
}

// Export the document to HTML.
func (h *HTMLExporter) Export(d *ast.Document) error {
	var hasTitle bool

	// Write the title.
	if !h.settings.OmitTitle && d.Title != "" {
		h.stream.Write([]byte("<h1>"))
		h.stream.Write([]byte(html.EscapeString(d.Title)))
		h.stream.Write([]byte("</h1>\n"))
		hasTitle = true
	}

	// Write the subtitle.
	if !h.settings.OmitSubtitle && d.Subtitle != "" {
		h.stream.Write([]byte("<h2>"))
		h.stream.Write([]byte(html.EscapeString(d.Subtitle)))
		h.stream.Write([]byte("</h2>\n"))
		hasTitle = true
	}

	// Write the date.
	if !h.settings.OmitDate && d.Date != "" {
		h.stream.Write([]byte("<h3>"))
		h.stream.Write([]byte(html.EscapeString(d.Date)))
		h.stream.Write([]byte("</h3>\n"))
		hasTitle = true
	}

	// Write the author.
	if !h.settings.OmitAuthor && d.Author != "" {
		h.stream.Write([]byte("<h3>"))
		h.stream.Write([]byte(html.EscapeString(d.Author)))
		h.stream.Write([]byte("</h3>\n"))
		hasTitle = true
	}

	if hasTitle {
		h.stream.Write([]byte("<hr>\n"))
	}

	// Write the header.
	if h.settings.IncludeHeader {
	}

	// Write the content.
	for i := range d.Content {
		err := h.exportBlock(d.Content[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// Export a block to HTML.
func (h *HTMLExporter) exportBlock(b ast.Block) error {
	if b.GetAlignment() != ast.NoAlign {
		switch b.GetAlignment() {
		case ast.LeftAlign:
			h.stream.Write([]byte("<div style=\"text-align: left\">"))
			break
		case ast.RightAlign:
			h.stream.Write([]byte("<div style=\"text-align: right\">"))
			break
		case ast.CenterAlign:
			h.stream.Write([]byte("<div style=\"text-align: center\">"))
			break
		default:
			return errors.New("invalid ast")
		}
	} else {
		h.stream.Write([]byte("<div>"))
	}

	switch b.(type) {
	case *ast.Paragraph:
		// Write the inline block.
		block := b.(*ast.Paragraph)
		h.stream.Write([]byte("<p>"))
		for i := range block.Content {
			err := h.exportInlineBlock(block.Content[i])
			if err != nil {
				return err
			}
		}
		h.stream.Write([]byte("</p>"))
	default:
		return errors.New("invalid ast")
	}

	h.stream.Write([]byte("</div>\n"))
	return nil
}

// Export an inline block to HTMl.
func (h *HTMLExporter) exportInlineBlock(b ast.InlineBlock) error {
	switch b.(type) {
	case string:
		// Write the content.
		h.stream.Write([]byte(strings.ReplaceAll(html.EscapeString(b.(string)), "\n", "<br>")))
		break
	case ast.HyperlinkBlock:
		// Write the hyperlink.
		block := b.(ast.HyperlinkBlock)
		h.stream.Write([]byte("<a href=\"" + block.Destination + "\">"))
		for i := range block.Content {
			err := h.exportInlineBlock(block.Content[i])
			if err != nil {
				return err
			}
		}
		h.stream.Write([]byte("</a>"))
		break
	case ast.FormattingBlock:
		// Write the formatting.
		block := b.(ast.FormattingBlock)
		attrName, err := getHTMLFormattingTagName(&block)
		if err != nil {
			return err
		}
		h.stream.Write([]byte("<" + attrName + ">"))
		for i := range block.Content {
			err = h.exportInlineBlock(block.Content[i])
			if err != nil {
				return err
			}
		}
		h.stream.Write([]byte("</" + attrName + ">"))
		break
	case ast.SizeBlock:
		// Write the size block.
		block := b.(ast.SizeBlock)
		paramValue, err := getHTMLFontSizeParameter(&block)
		if err != nil {
			return err
		}
		h.stream.Write([]byte("<span style=\"font-size: " + paramValue + "\">"))
		for i := range block.Content {
			err = h.exportInlineBlock(block.Content[i])
			if err != nil {
				return err
			}
		}
		h.stream.Write([]byte("</span>"))
		break
	default:
		return errors.New("invalid ast")
	}

	return nil
}

// Get the tag name for a formatting block.
func getHTMLFormattingTagName(b *ast.FormattingBlock) (string, error) {
	switch b.Attribute {
	case ast.BoldFormatting:
		return "b", nil
	case ast.ItalicFormatting:
		return "i", nil
	case ast.StrikethroughFormatting:
		return "s", nil
	case ast.UnderlineFormatting:
		return "u", nil
	case ast.TeletypeFormatting:
		return "code", nil
	}
	return "", errors.New("invalid ast")
}

// Get the font-size parameter for a size block.
func getHTMLFontSizeParameter(b *ast.SizeBlock) (string, error) {
	switch b.Type {
	case ast.PercentageSizeType:
		return fmt.Sprintf("%f%%", b.Value), nil
	case ast.PixelSizeType:
		return fmt.Sprintf("%fpx", b.Value), nil
	case ast.PointSizeType:
		return fmt.Sprintf("%fpt", b.Value), nil
	case ast.CentimeterSizeType:
		return fmt.Sprintf("%fcm", b.Value), nil
	case ast.MillimeterSizeType:
		return fmt.Sprintf("%fmm", b.Value), nil
	}
	return "", errors.New("invalid ast")
}
