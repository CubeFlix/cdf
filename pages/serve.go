// pages/serve.go
// CDF pages server.

package pages

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
)

// Handle a request to the server.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if path.Dir(r.URL.Path) == "/static" {
		// Serve static file.
		data, err := os.ReadFile(path.Join(s.Path, "static", path.Base(r.URL.Path)))
		if err != nil {
			s.NotFound(w, r, r.URL.Path)
			return
		}
		w.Write(data)
		return
	}

	var page string

	if r.URL.Path == "/" {
		page = "index"
	} else {
		page = path.Base(r.URL.Path)
	}

	// Serve a page.
	if _, err := os.Stat(path.Join(s.Path, "compiled", page+".html")); os.IsNotExist(err) {
		// Compiled page does not exist.
		s.NotFound(w, r, page)
		return
	} else {
		// Compiled page exists.
		s.ServeCompiledPage(w, r, page)
	}
	return
}

// Respond with a compiled HTML page.
func (s *Server) ServeCompiledPage(w http.ResponseWriter, r *http.Request, page string) {
	// Load the page content.
	data, err := os.ReadFile(path.Join(s.Path, "compiled", page+".html"))
	if err != nil {
		s.Error(w, r, err)
		return
	}

	// Get page info.
	pageInfo, ok := s.Pages[page]
	if !ok {
		s.Error(w, r, errors.New("compiled page info not found"))
		return
	}

	// Serve the contents.
	err = s.PageTemplate.Execute(w, PageTemplate{
		Page:     page,
		Title:    pageInfo.Title,
		Subtitle: pageInfo.Subtitle,
		Author:   pageInfo.Author,
		Date:     pageInfo.Date,
		Content:  template.HTML(data)})
	if err != nil {
		s.Error(w, r, err)
		return
	}
}

// Respond with a 404.
func (s *Server) NotFound(w http.ResponseWriter, r *http.Request, page string) {
	// Respond with the 404.
	w.WriteHeader(http.StatusNotFound)
	err := s.NotFoundTemplate.Execute(w, NotFoundTemplate{page})
	if err != nil {
		s.Error(w, r, err)
	}
}

// Respond with a server error.
func (s *Server) Error(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf("internal server error: %s", err.Error())))
}
