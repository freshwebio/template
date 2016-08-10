package template

import "io"

type Provider interface {
	Render(io.Writer, string, interface{})
	RenderWithLayout(io.Writer, string, string, interface{})
	HasTemplate(tmpl string) bool
	Invalidate()
}
