// export/html/html.go
// Package html provides functionality for exporting into HTML.

package html

import (
	"io"
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
