package template

import "io"

type Provider interface {
	Render(io.Writer, string, interface{}) error
	RenderWithLayout(io.Writer, string, string, interface{}) error
	RenderMultiple(io.Writer, []string, interface{}) error
	HasTemplate(tmpl string) bool
	Invalidate()
}
