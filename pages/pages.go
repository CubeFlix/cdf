// pages/pages.go
// CDF pages.

package pages

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cubeflix/cdf/export/html"
	"github.com/cubeflix/cdf/parser"
)

// The CDF pages server.
type Server struct {
	Path string

	Pages map[string]PageInfo

	PageTemplate        *template.Template
	NotFoundTemplate    *template.Template
	InvalidPageTemplate *template.Template

	staticHandler http.Handler
}

// Page info.
type PageInfo struct {
	Title    string
	Subtitle string
	Date     string
	Author   string

	DidError bool
	Error    string
}

// Page content template.
type PageTemplate struct {
	Page     string
	Title    string
	Subtitle string
	Date     string
	Author   string
	Content  template.HTML
}

// Not found template.
type NotFoundTemplate struct {
	Page string
}

// Invalid page template.
type InvalidPageTemplate struct {
	Page  string
	Error string
}

// Get a filename without extension. Source: https://gist.github.com/ivanzoid/129460aa08aff72862a534ebe0a9ae30
func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

// Load a new server from a path.
func LoadServer(p string) (*Server, error) {
	s := &Server{Path: p, Pages: map[string]PageInfo{}}

	// Load the templates.
	pageTemplate, err := os.ReadFile(path.Join(p, "template.html"))
	if err != nil {
		return nil, err
	}
	s.PageTemplate = template.Must(template.New("page").Parse(string(pageTemplate)))
	notFoundTemplate, err := os.ReadFile(path.Join(p, "404.html"))
	if err != nil {
		return nil, err
	}
	s.NotFoundTemplate = template.Must(template.New("404").Parse(string(notFoundTemplate)))
	invalidPageTemplate, err := os.ReadFile(path.Join(p, "invalid.html"))
	if err != nil {
		return nil, err
	}
	s.InvalidPageTemplate = template.Must(template.New("invalid").Parse(string(invalidPageTemplate)))

	// Recompile the files.
	os.RemoveAll(path.Join(p, "compiled"))
	err = os.Mkdir(path.Join(p, "compiled"), 0777)
	if err != nil {
		return nil, err
	}
	stat, err := os.ReadDir(path.Join(p, "pages"))
	if err != nil {
		return nil, err
	}
	for i := range stat {
		if path.Ext(stat[i].Name()) != ".cdf" {
			continue
		}
		if err := s.CompilePage(fileNameWithoutExtension(stat[i].Name())); err != nil {
			return nil, err
		}
	}

	// Create the static file handler.
	s.staticHandler = http.StripPrefix("/static/", http.FileServer(http.Dir(path.Join(p, "static"))))

	return s, nil
}

// Compile a single page.
func (s *Server) CompilePage(page string) error {
	// Compile the page.
	source, err := os.ReadFile(path.Join(s.Path, "pages", page+".cdf"))
	if err != nil {
		return err
	}
	outFile, err := os.Create(path.Join(s.Path, "compiled", page+".html"))
	if err != nil {
		return err
	}
	defer outFile.Close()

	parser := parser.NewParser(source)
	exporter := html.NewHTMLExporter(outFile, html.HTMLSettings{})

	// Parse the page.
	if err := parser.Parse(); err != nil {
		// Parsing failed.
		s.Pages[page] = PageInfo{
			Title:    parser.Tree.Title,
			Subtitle: parser.Tree.Subtitle,
			Author:   parser.Tree.Author,
			Date:     parser.Tree.Date,
			DidError: true,
			Error:    err.Error(),
		}
		return nil
	}

	// Export the page.
	if err := exporter.Export(&parser.Tree); err != nil {
		// Exporting failed.
		s.Pages[page] = PageInfo{
			Title:    parser.Tree.Title,
			Subtitle: parser.Tree.Subtitle,
			Author:   parser.Tree.Author,
			Date:     parser.Tree.Date,
			DidError: true,
			Error:    err.Error(),
		}
		return nil
	}

	s.Pages[page] = PageInfo{
		Title:    parser.Tree.Title,
		Subtitle: parser.Tree.Subtitle,
		Author:   parser.Tree.Author,
		Date:     parser.Tree.Date,
	}

	return nil
}
