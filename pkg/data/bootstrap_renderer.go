package data

import (
	"fmt"
	"html/template"
	"io"
)

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

func (r *BootstrapRenderer) BeginForm() {
	fmt.Fprintf(r.w, `<form method="%s" action="%s" name="%s" class="%s">`, r.theme.FormMethod, r.theme.FormAction, r.theme.FormName, r.theme.FormClass)
	fmt.Fprintln(r.w)
}

func (r *BootstrapRenderer) EndForm() {
	fmt.Fprintf(r.w, `  <button class="%s" type="submit">%s</button>`,
		r.theme.SubmitButtonClass,
		r.theme.SubmitButtonName,
	)
	fmt.Fprintln(r.w)
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

func (r *BootstrapRenderer) Input(label, name, inputType, value string) {
	fmt.Fprintf(r.w, `    <div class="%s">`, r.theme.FieldWrapper)
	fmt.Fprintln(r.w)
	fmt.Fprintf(r.w, `      <label class="%s">%s: <input type="%s" name="%s" value="%s" class="%s"/></label>`,
		r.theme.LabelClass,
		template.HTMLEscapeString(label),
		inputType,
		name,
		value,
		r.theme.InputClass,
	)
	fmt.Fprintln(r.w)
	fmt.Fprintln(r.w, `    </div>`)
}

func (r *BootstrapRenderer) Checkbox(label, name string, value bool) {
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
