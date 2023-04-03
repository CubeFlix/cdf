// export/html/settings.go
// HTML export settings.

package html

import "github.com/cubeflix/cdf/export"

const (
	DefaultQuoteBlockClass   = "quote"
	DefaultImageBlockClass   = "image-block"
	DefaultImageCaptionClass = "image-caption"
)

// HTML export settings.
type HTMLSettings struct {
	export.Settings

	IncludeHeader bool
	IncludeFooter bool

	UseCustomQuoteBlockClass bool
	QuoteBlockClass          string

	UseCustomImageBlockClass bool
	ImageBlockClass          string

	UseCustomImageCaptionClass bool
	ImageCaptionClass          string
}
