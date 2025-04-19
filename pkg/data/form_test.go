package data

import (
	"app/pkg/scrapper"
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

func TestScrapGeneratedForms(t *testing.T) {
	tests := map[string]struct {
		Input          map[string]any
		FormName       string
		ExpectedMethod string
		ExpectedAction string
		ExpectedFields map[string]string // fieldName -> elementType (input/select/etc)
	}{
		"simple contact form": {
			Input: map[string]any{
				"email":     "email",
				"firstname": "string",
				"gender":    "enum:Mr;Mme",
			},
			FormName:       "contact",
			ExpectedMethod: "POST",
			ExpectedAction: "/submit",
			ExpectedFields: map[string]string{
				"email":     "input",
				"firstname": "input",
				"gender":    "select",
			},
		},
		"array with nested fields": {
			Input: map[string]any{
				"children": []any{
					map[string]any{
						"name": "string",
						"age":  "number",
					},
				},
			},
			FormName:       "family",
			ExpectedMethod: "POST",
			ExpectedAction: "/submit",
			ExpectedFields: map[string]string{
				"children.0.name": "input",
				"children.0.age":  "input",
				"children.1.name": "input",
				"children.1.age":  "input",
			},
		},
		"deeply nested and mixed structure": {
			Input: map[string]any{
				"company": map[string]any{
					"name": "string",
					"address": map[string]any{
						"street": "string",
						"city":   "string",
						"coords": map[string]any{
							"lat":  "number",
							"long": "number",
						},
					},
					"departments": []any{
						map[string]any{
							"name": "string",
							"manager": map[string]any{
								"firstname": "string",
								"gender":    "enum:Mr;Mme",
							},
						},
					},
				},
			},
			FormName:       "company-form",
			ExpectedMethod: "POST",
			ExpectedAction: "/submit",
			ExpectedFields: map[string]string{
				"company.name":                "input",
				"company.address.street":      "input",
				"company.address.city":        "input",
				"company.address.coords.lat":  "input",
				"company.address.coords.long": "input",

				"company.departments.0.name":              "input",
				"company.departments.0.manager.firstname": "input",
				"company.departments.0.manager.gender":    "select",
				"company.departments.1.name":              "input",
				"company.departments.1.manager.firstname": "input",
				"company.departments.1.manager.gender":    "select",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			field := Parse(tt.Input, true)
			html := GenerateFormHTMLWithName(field, ThemeNoStyle, tt.FormName)

			doc, err := scrapper.NewDocumentFromReader(strings.NewReader(html))
			if err != nil {
				t.Fatalf("error parsing HTML: %v", err)
			}

			form, err := scrapper.GetForm(doc, tt.FormName)
			if err != nil {
				t.Fatalf("GetForm error: %v", err)
			}

			if form.Method != tt.ExpectedMethod {
				t.Errorf("Expected method %s, got %s", tt.ExpectedMethod, form.Method)
			}
			if form.Action != tt.ExpectedAction {
				t.Errorf("Expected action %s, got %s", tt.ExpectedAction, form.Action)
			}

			for fieldName, expectedType := range tt.ExpectedFields {
				typ, ok := form.Attribute[fieldName]
				if !ok {
					t.Errorf("Missing field '%s' in form", fieldName)
					continue
				}
				if typ != expectedType {
					t.Errorf("Expected type '%s' for field '%s', got '%s'", expectedType, fieldName, typ)
				}
			}
		})
	}
}
