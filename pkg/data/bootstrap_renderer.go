package data

import (
	"fmt"
	"html/template"
	"io"
	"strings"
)

var defaultBootstrapTheme = FormTheme{
	FormClass:         "form-bootstrap",
	LabelClass:        "form-label",
	InputClass:        "form-control",
	SelectClass:       "form-select",
	CheckboxClass:     "form-check-input",
	FieldWrapper:      "mb-3",
	RowWrapper:        "row",
	FieldsetClass:     "mb-4",
	LegendClass:       "fw-bold",
	AddButtonClass:    "btn btn-secondary",
	DeleteButtonClass: "btn btn-outline-danger",
	SubmitButtonClass: "btn btn-primary",
}

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
	Repeat            int // default repeat count for arrays
}

type contextItem struct {
	kind string // "form", "fieldset", "array"
	name string
}

type BootstrapRenderer struct {
	w     io.Writer
	theme FormTheme
	ctx   []contextItem
	jsFns map[string]bool
}

func NewBootstrapRenderer(w io.Writer, theme FormTheme) *BootstrapRenderer {
	return &BootstrapRenderer{
		w:     w,
		theme: theme,
		jsFns: make(map[string]bool),
	}
}

// Context stack management
func (r *BootstrapRenderer) push(kind, name string) {
	r.ctx = append(r.ctx, contextItem{kind, name})
}

func (r *BootstrapRenderer) pop(kind string) string {
	if len(r.ctx) == 0 {
		return ""
	}
	last := r.ctx[len(r.ctx)-1]
	if last.kind != kind {
		return ""
	}
	r.ctx = r.ctx[:len(r.ctx)-1]
	return last.name
}

// Form rendering
func (r *BootstrapRenderer) BeginForm(name, action, method string) {
	fmt.Fprintf(r.w, `<form method="%s" action="%s" name="%s" class="%s">`, method, action, name, r.theme.FormClass)
	fmt.Fprintln(r.w)
	r.push("form", "")
}

func (r *BootstrapRenderer) EndForm() {
	r.Submit("Submit")
	fmt.Fprintln(r.w, `</form>`)
	r.pop("form")

	// Injecter dynamiquement le JS juste après (optionnellement on pourrait buffer à part)
	if len(r.jsFns) > 0 {
		fmt.Fprintln(r.w, `<script>`)
		r.JS()
		fmt.Fprintln(r.w, `</script>`)
	}
}

// Fieldset rendering
func (r *BootstrapRenderer) BeginFieldset(label string) {
	fmt.Fprintf(r.w, `  <fieldset class="%s">`, r.theme.FieldsetClass)
	fmt.Fprintln(r.w)
	fmt.Fprintf(r.w, `    <legend class="%s">%s</legend>`, r.theme.LegendClass, template.HTMLEscapeString(label))
	fmt.Fprintln(r.w)
	r.push("fieldset", label)
}

func (r *BootstrapRenderer) EndFieldset() {
	fmt.Fprintln(r.w, `  </fieldset>`)
	r.pop("fieldset")
}

// Input types
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

// Array rendering
func (r *BootstrapRenderer) BeginArray(name string, path string) {
	fmt.Fprintf(r.w, `    <div id="container-%s">`, path)
	fmt.Fprintln(r.w)
	r.push("array", path)
}

func (r *BootstrapRenderer) EndArray() {
	fmt.Fprintln(r.w, `    </div>`)
	path := r.pop("array")
	r.AddButton(path, "Add")
}

func (r *BootstrapRenderer) BeginArrayItem(index int) {
	fmt.Fprintf(r.w, `      <div data-item class="row %s">`, r.theme.FieldWrapper)
	fmt.Fprintln(r.w)
	fmt.Fprintln(r.w, `        <div class="col">`)
}

func (r *BootstrapRenderer) EndArrayItem() {
	fmt.Fprintln(r.w, `        </div>`)
	fmt.Fprintf(r.w, `        <div class="col-auto align-self-end">`)
	fmt.Fprintln(r.w)
	fmt.Fprintf(r.w, `          <button type="button" class="%s" onclick="this.closest('[data-item]').remove()">Delete</button>`,
		r.theme.DeleteButtonClass)
	fmt.Fprintln(r.w)
	fmt.Fprintln(r.w, `        </div>`)
	fmt.Fprintln(r.w, `      </div>`)
}

func (r *BootstrapRenderer) AddButton(path, label string) {
	safeFuncName := strings.ReplaceAll(path, ".", "_")

	fmt.Fprintf(r.w, `  <button type="button" class="%s mt-2" onclick="add_%s()">%s</button>`,
		r.theme.AddButtonClass,
		safeFuncName,
		template.HTMLEscapeString(label),
	)
	fmt.Fprintln(r.w)

	// Enregistrer le JS à générer
	if !r.jsFns[safeFuncName] {
		r.jsFns[safeFuncName] = true
	}
}

func (r *BootstrapRenderer) Submit(label string) {
	fmt.Fprintf(r.w, `  <button class="%s" type="submit">%s</button>`,
		r.theme.SubmitButtonClass,
		template.HTMLEscapeString(label),
	)
	fmt.Fprintln(r.w)
}

func (r *BootstrapRenderer) JS() {
	for name := range r.jsFns {
		fmt.Fprintf(r.w, `function add_%s() {
  const container = document.getElementById("container-%s");
  const template = container.querySelector("[data-item]");
  const clone = template.cloneNode(true);
  clone.querySelectorAll("input, select").forEach(el => el.value = "");
  container.appendChild(clone);
}`, name, strings.ReplaceAll(name, "_", "."))
		fmt.Fprintln(r.w)
	}
}
