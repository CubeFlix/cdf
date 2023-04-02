// export/html/settings.go
// HTML export settings.

package html

import "github.com/cubeflix/cdf/export"

// HTML export settings.
type HTMLSettings struct {
	export.Settings

	IncludeHeader bool
	IncludeFooter bool
}
