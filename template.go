package template

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	stdTimeFmt = "Monday 2 January 2006 15:04"
)

// ProviderImpl provides all the templates
// used by the application
// controllers to render responses.
type ProviderImpl struct {
	templates map[string]*template.Template
}

// Cache deals with loading all the templates under the root templates
// directory into memory accessible through a map of {template_name}
// excluding the file extension to template.Template which holds and allows
// execution of a template.
func Cache(fMap ...template.FuncMap) (*ProviderImpl, error) {
	templates, err := buildTemplates(fMap...)
	if err != nil {
		return nil, err
	}
	return &ProviderImpl{templates: templates}, nil
}

func buildTemplates(fMap ...template.FuncMap) (map[string]*template.Template, error) {
	templates := make(map[string]*template.Template)
	templatesDir := "templates"

	layouts, err := filepath.Glob(templatesDir + "/layouts/*.tmpl")
	if err != nil {
		return nil, err
	}

	// Add all the functions from the default funcMap() call to the custom
	// template function map provided.
	if len(fMap) > 0 {
		tFuncMap := funcMap()
		for key, fnc := range tFuncMap {
			fMap[0][key] = fnc
		}
	} else {
		fMap = make([]template.FuncMap, 1)
		fMap[0] = funcMap()
	}

	filepath.Walk(templatesDir+"/includes", func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".tmpl") {
			key := strings.TrimSuffix(f.Name(), ".tmpl")
			files := []string{path}
			for _, layout := range layouts {
				files = append(files, layout)
			}
			templates[key] = template.Must(template.New(key).Funcs(fMap[0]).ParseFiles(files...))
		}
		return nil
	})

	filepath.Walk(templatesDir+"/mail", func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".tmpl") {
			key := strings.TrimSuffix(f.Name(), ".tmpl")
			files := []string{path}
			for _, layout := range layouts {
				files = append(files, layout)
			}
			templates[key] = template.Must(template.New(key).Funcs(fMap[0]).ParseFiles(files...))
		}
		return nil
	})
	return templates, nil
}

func funcMap() template.FuncMap {
	return template.FuncMap{
		"oddoreven": func(i int) string {
			if i%2 == 0 {
				return "even"
			}
			return "odd"
		},
		"fmtdate": func(fmt string, t *time.Time) string {
			formatted := ""
			switch fmt {
			case "std":
				formatted = t.Format(stdTimeFmt)
			}
			return formatted
		},
		"ucfirstwords": func(input string) string {
			input = strings.ToUpper(input)
			words := strings.Fields(input)
			if len(words) == 0 {
				return ""
			}
			output := ""
			for _, word := range words {
				var first rune
				for _, c := range word {
					first = c
					break
				}
				output += string(first)
			}
			return output
		},
	}
}

// Invalidate deals with re-building the templates
// map in the case new templates had been added or removed
// or existing ones being updated. Usually will be called upon on a general
// app clear cache invoked from a web interface or an external app management tool.
func (tp *ProviderImpl) Invalidate() {
	templates, err := buildTemplates()
	if err == nil {
		tp.templates = templates
	}
}

// Render handles rendering the template
// held by the given ProviderImpl with the given
// name to the provided response writer with the given data.
func (tp *ProviderImpl) Render(w io.Writer, tmpl string, data interface{}) error {
	return tp.templates[tmpl].ExecuteTemplate(w, tmpl+".tmpl", data)
}

// RenderWithLayout handles rendering the provided template
// held by the ProviderImpl in the templates map with the given base template
// where the base template needs to be cached as a part of the template to be
// rendered.
func (tp *ProviderImpl) RenderWithLayout(w io.Writer, tmpl string, layout string, data interface{}) error {
	return tp.templates[tmpl].ExecuteTemplate(w, layout, data)
}

func (tp *ProviderImpl) HasTemplate(tmpl string) bool {
	hasTemplate := false
	if _, ok := tp.templates[tmpl]; ok {
		hasTemplate = true
	}
	return hasTemplate
}
