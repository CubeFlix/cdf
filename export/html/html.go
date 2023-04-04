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
	if !settings.UseCustomQuoteBlockClass {
		settings.QuoteBlockClass = DefaultQuoteBlockClass
	}
	if !settings.UseCustomImageBlockClass {
		settings.ImageBlockClass = DefaultImageBlockClass
	}
	if !settings.UseCustomImageCaptionClass {
		settings.ImageCaptionClass = DefaultImageCaptionClass
	}

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
		// TODO
	}

	// Write the content.
	for i := range d.Content {
		err := h.exportBlock(d.Content[i])
		if err != nil {
			return err
		}
	}

	// Write the footer.
	if h.settings.IncludeFooter {
		// TODO
	}

	return nil
}

// Export a block to HTML.
func (h *HTMLExporter) exportBlock(b ast.Block) error {
	alignmentStyle, err := getHTMLAlignmentStyleParameter(b)
	if err != nil {
		return err
	}

	switch b.(type) {
	case *ast.Paragraph:
		// Write the inline block.
		block := b.(*ast.Paragraph)
		h.stream.Write([]byte("<p" + wrapHTMLStyleParameter(alignmentStyle) + ">"))
		for i := range block.Content {
			err := h.exportInlineBlock(block.Content[i])
			if err != nil {
				return err
			}
		}
		h.stream.Write([]byte("</p>\n"))
		break
	case *ast.BasicBlock:
		// Write the basic block.
		block := b.(*ast.BasicBlock)
		h.stream.Write([]byte("<div" + wrapHTMLStyleParameter(alignmentStyle) + ">"))
		for i := range block.Content {
			err := h.exportBlock(block.Content[i])
			if err != nil {
				return err
			}
		}
		h.stream.Write([]byte("</div>\n"))
		break
	case *ast.Quote:
		// Write the quote block.
		block := b.(*ast.Quote)
		h.stream.Write([]byte("<div" + wrapHTMLStyleParameter(alignmentStyle) + " class=\"" + h.settings.QuoteBlockClass + "\">"))
		for i := range block.Content {
			err := h.exportBlock(block.Content[i])
			if err != nil {
				return err
			}
		}
		h.stream.Write([]byte("</div>\n"))
		break
	case *ast.Image:
		// Write the image block.
		block := b.(*ast.Image)
		sizeStyle, err := getHTMLImageWidthHeightStyleParameter(block)
		if err != nil {
			return err
		}
		if block.HasCaption {
			h.stream.Write([]byte("<div" + wrapHTMLStyleParameter(alignmentStyle) + " class=\"" + h.settings.ImageBlockClass + "\"><img" + wrapHTMLStyleParameter(sizeStyle) + " src=\"" + block.Source + "\"><div class=\"" + h.settings.ImageCaptionClass + "\">"))
			for i := range block.Caption {
				err := h.exportInlineBlock(block.Caption[i])
				if err != nil {
					return err
				}
			}
			h.stream.Write([]byte("</div></div>\n"))
		} else {
			h.stream.Write([]byte("<div" + wrapHTMLStyleParameter(alignmentStyle) + " class=\"" + h.settings.ImageBlockClass + "\"><img" + wrapHTMLStyleParameter(sizeStyle) + " src=\"" + block.Source + "\"></div>\n"))
		}
		break
	case *ast.Heading:
		// Write the heading block.
		block := b.(*ast.Heading)
		headingClass, err := getHTMLHeaderType(block)
		if err != nil {
			return err
		}
		h.stream.Write([]byte("<" + headingClass + wrapHTMLStyleParameter(alignmentStyle) + ">"))
		for i := range block.Content {
			err := h.exportInlineBlock(block.Content[i])
			if err != nil {
				return err
			}
		}
		h.stream.Write([]byte("</" + headingClass + ">\n"))
		break
	case *ast.HorizontalRule:
		// Write the horizontal rule.
		h.stream.Write([]byte("<hr>\n"))
		break
	case *ast.List:
		// Write the list block.
		block := b.(*ast.List)
		if block.Ordered {
			h.stream.Write([]byte("<ol" + wrapHTMLStyleParameter(alignmentStyle) + ">\n"))
			for i := range block.Items {
				h.stream.Write([]byte("<li>"))
				err := h.exportBlock(block.Items[i])
				if err != nil {
					return err
				}
				h.stream.Write([]byte("</li>\n"))
			}
			h.stream.Write([]byte("</ol>\n"))
		} else {
			h.stream.Write([]byte("<ul" + wrapHTMLStyleParameter(alignmentStyle) + ">\n"))
			for i := range block.Items {
				h.stream.Write([]byte("<li>"))
				err := h.exportBlock(block.Items[i])
				if err != nil {
					return err
				}
				h.stream.Write([]byte("</li>\n"))
			}
			h.stream.Write([]byte("</ul>\n"))
		}
		break
	default:
		return errors.New("invalid ast")
	}

	return nil
}

// Get the div's optional alignment style. Return an empty string if no
// alignment is provided.
func getHTMLAlignmentStyleParameter(b ast.Block) (string, error) {
	if b.GetAlignment() != ast.NoAlign {
		if b.GetWrap() {
			switch b.GetAlignment() {
			case ast.LeftAlign:
				return "float: left;", nil
			case ast.RightAlign:
				return "float: right;", nil
			case ast.CenterAlign:
				return "float: center;", nil
			default:
				return "", errors.New("invalid ast")
			}
		}
		switch b.GetAlignment() {
		case ast.LeftAlign:
			return "text-align: left;", nil
		case ast.RightAlign:
			return "text-align: right;", nil
		case ast.CenterAlign:
			return "text-align: center;", nil
		default:
			return "", errors.New("invalid ast")
		}
	} else {
		return "", nil
	}
}

// Get the style parameters for width and height from an image block.
func getHTMLImageWidthHeightStyleParameter(b *ast.Image) (string, error) {
	var param string
	if b.HasWidthParameter {
		switch b.WidthType {
		case ast.PercentageSizeType:
			param += fmt.Sprintf("width: %f%%;", b.WidthValue)
		case ast.PixelSizeType:
			param += fmt.Sprintf("width: %fpx;", b.WidthValue)
		case ast.PointSizeType:
			param += fmt.Sprintf("width: %fpt;", b.WidthValue)
		case ast.CentimeterSizeType:
			param += fmt.Sprintf("width: %fcm;", b.WidthValue)
		case ast.MillimeterSizeType:
			param += fmt.Sprintf("width: %fmm;", b.WidthValue)
		default:
			return "", errors.New("invalid ast")
		}
	}
	if b.HasHeightParameter {
		switch b.HeightType {
		case ast.PercentageSizeType:
			param += fmt.Sprintf("height: %f%%;", b.HeightValue)
		case ast.PixelSizeType:
			param += fmt.Sprintf("height: %fpx;", b.HeightValue)
		case ast.PointSizeType:
			param += fmt.Sprintf("height: %fpt;", b.HeightValue)
		case ast.CentimeterSizeType:
			param += fmt.Sprintf("height: %fcm;", b.HeightValue)
		case ast.MillimeterSizeType:
			param += fmt.Sprintf("height: %fmm;", b.HeightValue)
		default:
			return "", errors.New("invalid ast")
		}
	}
	return param, nil
}

// Wrap the style information in " style=\"\"". Returns an empty string if no
// style information is provided.
func wrapHTMLStyleParameter(style string) string {
	if len(strings.TrimSpace(style)) == 0 {
		return ""
	}
	return " style=\"" + style + "\""
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
	case ast.FontBlock:
		// Write the font block.
		block := b.(ast.FontBlock)
		h.stream.Write([]byte("<span style=\"font-family: " + block.Family + "\">"))
		for i := range block.Content {
			err := h.exportInlineBlock(block.Content[i])
			if err != nil {
				return err
			}
		}
		h.stream.Write([]byte("</span>"))
		break
	case ast.ColorBlock:
		// Write the color block.
		block := b.(ast.ColorBlock)
		h.stream.Write([]byte("<span style=\"" + getHTMLColorStyleParameter(&block) + "\">"))
		for i := range block.Content {
			err := h.exportInlineBlock(block.Content[i])
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

// Get the style parameter for a color block.
func getHTMLColorStyleParameter(b *ast.ColorBlock) string {
	var param string
	if b.ForegroundValue != nil {
		param += "color: " + b.ForegroundValue.String() + ";"
	}
	if b.BackgroundValue != nil {
		param += "background-color: " + b.BackgroundValue.String() + ";"
	}
	return param
}

// Get the header type "h1" - "h5" for a header block.
func getHTMLHeaderType(b *ast.Heading) (string, error) {
	switch b.Class {
	case ast.Heading1Type:
		return "h1", nil
	case ast.Heading2Type:
		return "h2", nil
	case ast.Heading3Type:
		return "h3", nil
	case ast.Heading4Type:
		return "h4", nil
	case ast.Heading5Type:
		return "h5", nil
	}
	return "", errors.New("invalid ast")
}
