// ast/ast.go
// Package ast provides the generic CDF abstract syntax tree.

package ast

// CDF document.
type Document struct {
	// Document metadata.
	Title     string
	Subtitle  string
	Date      string
	Author    string
	Copyright string
	Header    string
	Footer    string
}

// Paragraph alignment types.
type AlignmentType int64

const (
	FullAlign AlignmentType = iota
	LeftAlign
	RightAlign
	CenterAlign
)

// Paragraph struct for AST.
type Paragraph struct {
	Alignment AlignmentType
	Content   []InlineBlock
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
}

// Inline formatting types.
type FormattingType int64

const (
	BoldFormatting FormattingType = iota
	ItalicFormatting
	StrikethroughFormatting
	UnderlineFormatting
	TeletypeFormatting
	ColorFormatting
	SizeFormatting
)
