// export/export.go
// Package export provides a generic interface for exporting CBF files into
// alternate formats.

package export

import (
	"github.com/cubeflix/cdf/ast"
)

type Exporter interface {
	Export(ast.Document)
}
