// ast/ast.go
// Package ast provides the generic CDF abstract syntax tree.

package ast

import "image/color"

// CDF document.
type Document struct {
	// Document information.
	Title    string
	Subtitle string
	Date     string
	Author   string

	Header *InlineBlock
	Footer *InlineBlock

	// The blocks in the document.
	Content []Block
}

// Block for AST.
type Block interface {
	GetAlignment() AlignmentType
}

// Base block.
type BaseBlock struct {
	Alignment AlignmentType
}

// Get alignment.
func (b *BaseBlock) GetAlignment() AlignmentType {
	return b.Alignment
}

// Block alignment types.
type AlignmentType int64

const (
	NoAlign AlignmentType = iota
	LeftAlign
	RightAlign
	CenterAlign
)

// Paragraph struct for AST.
type Paragraph struct {
	BaseBlock

	Content []InlineBlock
}

// Inline block for AST. An inline block may be a base inline block or string.
type InlineBlock interface{}

// Base inline block.
type BaseInlineBlock struct {
	// The block's content. A slice of more inline blocks.
	Content []InlineBlock
}

// Hyperlink block.
type HyperlinkBlock struct {
	BaseInlineBlock

	Destination string
}

// Formatting block.
type FormattingBlock struct {
	BaseInlineBlock

	Attribute FormattingType
}

// Inline formatting types.
type FormattingType int64

const (
	BoldFormatting FormattingType = iota
	ItalicFormatting
	StrikethroughFormatting
	UnderlineFormatting
	TeletypeFormatting
)

// Color block.
type ColorBlock struct {
	BaseInlineBlock

	Value color.Color
}

// Size block.
type SizeBlock struct {
	BaseInlineBlock

	Value float32
	Type  SizeType
}

// Size value type.
type SizeType int64

const (
	PercentageSizeType SizeType = iota
	PixelSizeType
	PointSizeType
	CentimeterSizeType
	MillimeterSizeType
)

// Font block.
type FontBlock struct {
	BaseInlineBlock

	Family string
}
