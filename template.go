package template

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Provider provides all the templates
// used by the application
// controllers to render responses.
type Provider struct {
	templates map[string]*template.Template
}

// Cache deals with loading all the templates under the root templates
// directory into memory accessible through a map of {template_name}
// excluding the file extension to template.Template which holds and allows
// execution of a template.
func Cache() *Provider {
	templates := buildTemplates()
	return &Provider{templates: templates}
}

func buildTemplates() map[string]*template.Template {
	templates := make(map[string]*template.Template)
	templatesDir := "templates"

	layouts, err := filepath.Glob(templatesDir + "/layouts/*.tmpl")
	if err != nil {
		panic(err)
	}

	filepath.Walk(templatesDir+"/includes", func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".tmpl") {
			key := strings.TrimSuffix(f.Name(), ".tmpl")
			files := []string{path}
			for _, layout := range layouts {
				files = append(files, layout)
			}
			templates[key] = template.Must(template.ParseFiles(files...))
		}
		return nil
	})
	return templates
}

// Invalidate deals with re-building the templates
// map in the case new templates had been added or removed
// or existing ones being updated.
func (tp *Provider) Invalidate() {
	tp.templates = buildTemplates()
}

// Render handles rendering the template
// held by the given provider with the given
// name to the provided response writer with the given data.
func (tp *Provider) Render(w http.ResponseWriter, tmpl string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tp.templates[tmpl].ExecuteTemplate(w, tmpl+".tmpl", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RenderWithLayout handles rendering the provided template
// held by the provider in the templates map with the given base template
// where the base template needs to be cached as a part of the template to be
// rendered.
func (tp *Provider) RenderWithLayout(w http.ResponseWriter, tmpl string, layout string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tp.templates[tmpl].ExecuteTemplate(w, layout, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (tp *Provider) HasTemplate(tmpl string) bool {
	hasTemplate := false
	if _, ok := tp.templates[tmpl]; ok {
		hasTemplate = true
	}
	return hasTemplate
}
