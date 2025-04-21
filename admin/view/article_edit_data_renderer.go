package view

import (
	"app/pkg/data"
	"bytes"
	"html/template"
)

func RenderData(form_data map[string]any) (template.HTML, error) {
	var buf bytes.Buffer
	renderer := data.NewBootstrapRenderer(&buf, data.ThemeBootstrap)
	err := data.Render(form_data, renderer)
	return template.HTML(buf.String()), err
}
