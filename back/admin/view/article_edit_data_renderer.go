package view

import (
	"app/pkg/data"
	"bytes"
	"fmt"
	"html/template"
	"io"
)

func RenderData(form_data map[string]any) (template.HTML, error) {
	var buf bytes.Buffer
	renderer := &MyRenderer{w: &buf}
	err := data.Render(form_data, renderer)
	return template.HTML(buf.String()), err
}

type MyRenderer struct {
	w io.Writer
}

func (r *MyRenderer) BeginForm() {}
func (r *MyRenderer) EndForm()   {}

func (r *MyRenderer) BeginObject(label string) {
	fmt.Fprintf(r.w, `  <fieldset class="mb-4">`)
	fmt.Fprintln(r.w)
	fmt.Fprintf(r.w, `    <legend class="fw-bold">%s</legend>`, template.HTMLEscapeString(label))
	fmt.Fprintln(r.w)
}

func (r *MyRenderer) EndObject() {
	fmt.Fprintln(r.w, `  </fieldset>`)
}

func (r *MyRenderer) Input(label, name, inputType, value string) {
	fmt.Fprintf(r.w, `    <div class="mb-3">`)
	fmt.Fprintln(r.w)
	fmt.Fprintf(r.w, `      <label class="form-label">%s: <input type="%s" name="%s" value="%s" class="form-control"/></label>`,
		template.HTMLEscapeString(label),
		inputType,
		name,
		value,
	)
	fmt.Fprintln(r.w)
	fmt.Fprintln(r.w, `    </div>`)
}

func (r *MyRenderer) Checkbox(label, name string, value bool) {
	fmt.Fprintf(r.w, `    <div class="mb-3">`)
	fmt.Fprintln(r.w)
	fmt.Fprintf(r.w, `      <label class="form-label">%s: <input type="checkbox" name="%s" value="true" class="form-check-input"/></label>`,
		template.HTMLEscapeString(label),
		name)
	fmt.Fprintln(r.w)
	fmt.Fprintln(r.w, `    </div>`)
}
