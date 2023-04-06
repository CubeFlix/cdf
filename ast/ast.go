// ast/ast.go
// Package ast provides the generic CDF abstract syntax tree.

package ast

import "gopkg.in/go-playground/colors.v1"

// CDF document.
type Document struct {
	// Document information.
	Title    string
	Subtitle string
	Date     string
	Author   string

	// The blocks in the document.
	Content []Block
}

// Block for AST.
type Block interface {
	GetAlignment() AlignmentType
	GetWrap() bool
}

// Base block.
type BaseBlock struct {
	Alignment AlignmentType
	Wrap      bool
}

// Get alignment.
func (b *BaseBlock) GetAlignment() AlignmentType {
	return b.Alignment
}

// Get wrap settings.
func (b *BaseBlock) GetWrap() bool {
	return b.Wrap
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

// Basic block struct for AST.
type BasicBlock struct {
	BaseBlock

	Content []Block
}

// Quote block struct for AST.
type Quote struct {
	BaseBlock

	Content []Block
}

// Image block struct for AST.
type Image struct {
	BaseBlock

	Source     string
	HasCaption bool
	Caption    []InlineBlock

	// Image size information.
	HasWidthParameter  bool
	WidthValue         float32
	WidthType          SizeType
	HasHeightParameter bool
	HeightValue        float32
	HeightType         SizeType
}

// Heading block struct for AST.
type Heading struct {
	BaseBlock

	Class   HeadingType
	Content []InlineBlock
}

// Heading type.
type HeadingType int64

const (
	Heading1Type HeadingType = iota
	Heading2Type
	Heading3Type
	Heading4Type
	Heading5Type
)

// Horizontal rule block struct for AST.
type HorizontalRule struct {
	BaseBlock
}

// List block.
type List struct {
	BaseBlock

	Items   []Block
	Ordered bool
}

// Table block.
type Table struct {
	BaseBlock

	Rows []TableRow
}

// Table row block.
type TableRow struct {
	Cells []TableCell
}

// Table row cell block.
type TableCell struct {
	Content  []Block
	IsHeader bool
}

// Collapsed block.
type Collapse struct {
	BaseBlock

	Summary []InlineBlock
	Content []Block
}

// Page break block.
type PageBreak struct {
	BaseBlock
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

	// nil values indicate unspecified value.
	ForegroundValue colors.Color
	BackgroundValue colors.Color
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

// Inline image block.
type InlineImageBlock struct {
	BaseInlineBlock

	Source string

	// Image size information.
	HasWidthParameter  bool
	WidthValue         float32
	WidthType          SizeType
	HasHeightParameter bool
	HeightValue        float32
	HeightType         SizeType
}
