package data

import (
	"fmt"
	"html/template"
	"io"
)

type FormTheme struct {
	FormClass         string
	LabelClass        string
	InputClass        string
	SelectClass       string
	CheckboxClass     string
	FieldWrapper      string
	RowWrapper        string
	FieldsetClass     string
	LegendClass       string
	AddButtonClass    string
	DeleteButtonClass string
	SubmitButtonClass string
}

type BootstrapRenderer struct {
	w     io.Writer
	theme FormTheme
}

func NewBootstrapRenderer(w io.Writer, theme FormTheme) *BootstrapRenderer {
	return &BootstrapRenderer{
		w:     w,
		theme: theme,
	}
}

func (r *BootstrapRenderer) BeginForm(name, action, method string) {
	fmt.Fprintf(r.w, `<form method="%s" action="%s" name="%s" class="%s">`, method, action, name, r.theme.FormClass)
	fmt.Fprintln(r.w)
}

func (r *BootstrapRenderer) EndForm() {
	r.Submit("Submit")
	fmt.Fprintln(r.w, `</form>`)
}

func (r *BootstrapRenderer) BeginFieldset(label string) {
	fmt.Fprintf(r.w, `  <fieldset class="%s">`, r.theme.FieldsetClass)
	fmt.Fprintln(r.w)
	fmt.Fprintf(r.w, `    <legend class="%s">%s</legend>`, r.theme.LegendClass, template.HTMLEscapeString(label))
	fmt.Fprintln(r.w)
}

func (r *BootstrapRenderer) EndFieldset() {
	fmt.Fprintln(r.w, `  </fieldset>`)
}

func (r *BootstrapRenderer) BeginArray(name string, path string) {
	fmt.Fprintf(r.w, `    <div id="container-%s">`, path)
	fmt.Fprintln(r.w)
}

func (r *BootstrapRenderer) EndArray() {
	fmt.Fprintln(r.w, `    </div>`)
}

func (r *BootstrapRenderer) BeginArrayItem(index int) {
	fmt.Fprintf(r.w, `      <div data-item class="row %s">`, r.theme.FieldWrapper)
	fmt.Fprintln(r.w)
	fmt.Fprintln(r.w, `        <div class="col">`)
}

func (r *BootstrapRenderer) EndArrayItem() {
	fmt.Fprintln(r.w, `        </div>`)
	fmt.Fprintln(r.w, `      </div>`)
}

func (r *BootstrapRenderer) Input(label, name, inputType string) {
	fmt.Fprintf(r.w, `    <div class="%s">`, r.theme.FieldWrapper)
	fmt.Fprintln(r.w)
	fmt.Fprintf(r.w, `      <label class="%s">%s: <input type="%s" name="%s" class="%s"/></label>`,
		r.theme.LabelClass,
		template.HTMLEscapeString(label),
		inputType,
		name,
		r.theme.InputClass)
	fmt.Fprintln(r.w)
	fmt.Fprintln(r.w, `    </div>`)
}

func (r *BootstrapRenderer) Checkbox(label, name string) {
	fmt.Fprintf(r.w, `    <div class="%s">`, r.theme.FieldWrapper)
	fmt.Fprintln(r.w)
	fmt.Fprintf(r.w, `      <label class="%s">%s: <input type="checkbox" name="%s" value="true" class="%s"/></label>`,
		r.theme.LabelClass,
		template.HTMLEscapeString(label),
		name,
		r.theme.CheckboxClass)
	fmt.Fprintln(r.w)
	fmt.Fprintln(r.w, `    </div>`)
}

func (r *BootstrapRenderer) Select(label, name string, options []string) {
	fmt.Fprintf(r.w, `    <div class="%s">`, r.theme.FieldWrapper)
	fmt.Fprintln(r.w)
	fmt.Fprintf(r.w, `      <label class="%s">%s: <select name="%s" class="%s">`,
		r.theme.LabelClass,
		template.HTMLEscapeString(label),
		name,
		r.theme.SelectClass)
	fmt.Fprintln(r.w)

	for _, opt := range options {
		fmt.Fprintf(r.w, `        <option value="%s">%s</option>`, opt, opt)
		fmt.Fprintln(r.w)
	}

	fmt.Fprintln(r.w, `      </select></label>`)
	fmt.Fprintln(r.w, `    </div>`)
}

func (r *BootstrapRenderer) Submit(label string) {
	fmt.Fprintf(r.w, `  <button class="%s" type="submit">%s</button>`,
		r.theme.SubmitButtonClass,
		template.HTMLEscapeString(label),
	)
	fmt.Fprintln(r.w)
}
