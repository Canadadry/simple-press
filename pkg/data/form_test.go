package data

import (
	"strings"
	"testing"
)

var ThemeNoStyle = FormTheme{
	FormClass:     "",
	LabelClass:    "",
	InputClass:    "",
	SelectClass:   "",
	CheckboxClass: "",
	FieldWrapper:  "",
	RowWrapper:    "",
	FieldsetClass: "",
	LegendClass:   "",
	Repeat:        5,
}

func TestGenerateFormHTML(t *testing.T) {
	tests := map[string]struct {
		Input    map[string]any
		Contains []string
	}{
		"flat object": {
			Input: map[string]any{
				"firstname":  "string",
				"email":      "email",
				"newsletter": "bool",
			},
			Contains: []string{
				`name="firstname"`,
				`name="email"`,
				`name="newsletter"`,
				`type="text"`,
				`type="email"`,
				`type="checkbox"`,
			},
		},
		"object with enum and date": {
			Input: map[string]any{
				"civility": "enum:Mr;Mme",
				"dob":      "date",
			},
			Contains: []string{
				`<select name="civility"`,
				`<option value="Mr"`,
				`<option value="Mme"`,
				`name="dob"`,
				`type="date"`,
			},
		},
		"nested object": {
			Input: map[string]any{
				"profile": map[string]any{
					"name": map[string]any{
						"first": "string",
						"last":  "string",
					},
				},
			},
			Contains: []string{
				`name="profile.name.first"`,
				`name="profile.name.last"`,
			},
		},
		"array of objects": {
			Input: map[string]any{
				"children": []any{
					map[string]any{
						"firstname": "string",
						"gender":    "enum:Mr;Mme",
					},
				},
			},
			Contains: []string{
				`name="children.0.firstname"`,
				`name="children.0.gender"`,
				`name="children.1.firstname"`,
				`name="children.1.gender"`,
			},
		},
		"array of primitives": {
			Input: map[string]any{
				"scores": []any{"number"},
			},
			Contains: []string{
				`name="scores.0"`,
				`name="scores.1"`,
				`type="number"`,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			root := Parse(tt.Input, true)
			html := GenerateFormHTML(root, ThemeNoStyle)

			for _, frag := range tt.Contains {
				if !strings.Contains(html, frag) {
					t.Errorf("Expected HTML to contain:\n%s\nGot:\n%s", frag, html)
				}
			}
		})
	}
}

func TestGenerateFormHTML_ExactMatch(t *testing.T) {
	tests := map[string]struct {
		Input        map[string]any
		Theme        FormTheme
		ExpectedHTML string
	}{
		"simple_object_with_enum": {
			Input: map[string]any{
				"civility":  "enum:Mr;Mme",
				"email":     "email",
				"firstname": "string",
			},
			Theme: ThemeNoStyle,
			ExpectedHTML: `<form method="POST" action="/submit" class="">` +
				`<div class=""><label class="">civility: <select name="civility" class="">` +
				`<option value="Mr">Mr</option><option value="Mme">Mme</option>` +
				`</select></label><br/></div>` +
				`<div class=""><label class="">email: <input type="email" name="email" class=""/></label><br/></div>` +
				`<div class=""><label class="">firstname: <input type="text" name="firstname" class=""/></label><br/></div>` +
				`<button type="submit">Submit</button></form>`,
		},
		"simple_object_with_enum_and_boostratp_theme": {
			Input: map[string]any{
				"civility":  "enum:Mr;Mme",
				"email":     "email",
				"firstname": "string",
			},
			Theme: FormTheme{
				FormClass:     "needs-validation",
				LabelClass:    "form-label",
				InputClass:    "form-control",
				SelectClass:   "form-select",
				CheckboxClass: "form-check-input",
				FieldWrapper:  "mb-3",
				RowWrapper:    "row",
				FieldsetClass: "mb-4",
				LegendClass:   "fw-bold",
				Repeat:        3,
			},
			ExpectedHTML: `<form method="POST" action="/submit" class="needs-validation">` +
				`<div class="mb-3"><label class="form-label">civility: <select name="civility" class="form-select">` +
				`<option value="Mr">Mr</option><option value="Mme">Mme</option>` +
				`</select></label><br/></div>` +
				`<div class="mb-3"><label class="form-label">email: <input type="email" name="email" class="form-control"/></label><br/></div>` +
				`<div class="mb-3"><label class="form-label">firstname: <input type="text" name="firstname" class="form-control"/></label><br/></div>` +
				`<button type="submit">Submit</button></form>`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			root := Parse(tt.Input, true)
			html := GenerateFormHTML(root, tt.Theme)

			if html != tt.ExpectedHTML {
				t.Errorf("HTML output mismatch.\nExpected:\n%s\n\nGot:\n%s", tt.ExpectedHTML, html)
			}
		})
	}
}
