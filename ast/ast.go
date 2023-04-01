// ast/ast.go
// Package ast provides the generic CDF abstract syntax tree.

package ast

import "image/color"

// CDF document.
type Document struct {
	// Document information.
	Title     string
	Subtitle  string
	Date      string
	Author    string
	Copyright string
	Header    string
	Footer    string
}

// Block for AST.
type Block interface {
	// Reserved for future use.
}

// Base block.
type BaseBlock struct {
	Alignment AlignmentType
	Children  []Block
}

// Block alignment types.
type AlignmentType int64

const (
	FullAlign AlignmentType = iota
	LeftAlign
	RightAlign
	CenterAlign
)

// Paragraph struct for AST.
type Paragraph struct {
	BaseBlock

	Content []InlineBlock
}

// Inline block for AST.
type InlineBlock interface {
	// Reserved for future use.
}

// Base inline block.
type BaseInlineBlock struct {
	Children []InlineBlock
}

// Content block.
type ContentBlock struct {
	BaseInlineBlock

	Content []byte
}

// Hyperlink block.
type HyperlinkBlock struct {
	BaseInlineBlock

	Content     []byte
	Destination []byte
}

// Formatting block.
type FormattingBlock struct {
	BaseInlineBlock

	Attributes []FormattingBlock
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
	PercentageType SizeType = iota
	PixelType
)
