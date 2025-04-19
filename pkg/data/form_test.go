package data

import (
	"app/pkg/scrapper"
	"strings"
	"testing"
)

func collectPaths(f Field, out *[]string) {
	*out = append(*out, f.Path)
	for _, child := range f.Children {
		collectPaths(child, out)
	}
}

func TestUpdatePathForArrayIndex(t *testing.T) {

	tests := map[string]struct {
		Input         Field
		Index         int
		ExpectedPaths []string
	}{
		"flat array": {
			Input: Field{
				Path: "children",
				Type: "array",
				Children: []Field{
					{Key: "firstname", Path: "children.0.firstname"},
					{Key: "age", Path: "children.0.age"},
				},
			},
			Index: 2,
			ExpectedPaths: []string{
				"children.2.firstname",
				"children.2.age",
			},
		},
		"deep nested array": {
			Input: Field{
				Path: "company.departments",
				Type: "array",
				Children: []Field{
					{
						Key:  "",
						Path: "company.departments.0",
						Type: "object",
						Children: []Field{
							{Key: "name", Path: "company.departments.0.name"},
							{Key: "manager", Path: "company.departments.0.manager", Type: "object", Children: []Field{
								{Key: "firstname", Path: "company.departments.0.manager.firstname"},
								{Key: "gender", Path: "company.departments.0.manager.gender"},
							}},
						},
					},
				},
			},
			Index: 3,
			ExpectedPaths: []string{
				"company.departments.3",
				"company.departments.3.name",
				"company.departments.3.manager",
				"company.departments.3.manager.firstname",
				"company.departments.3.manager.gender",
			},
		},
		"array in array": {
			Input: Field{
				Path: "matrix",
				Type: "array",
				Children: []Field{
					{
						Key:  "",
						Path: "matrix.0",
						Type: "array",
						Children: []Field{
							{Key: "", Path: "matrix.0.0"},
							{Key: "", Path: "matrix.0.1"},
						},
					},
				},
			},
			Index: 4,
			ExpectedPaths: []string{
				"matrix.4",
				"matrix.4.0",
				"matrix.4.1",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var collected []string

			// On met à jour tous les enfants de l'array
			for _, child := range tt.Input.Children {
				updated := updatePathForArrayIndex(child, tt.Index)
				collectPaths(updated, &collected)
			}

			if len(collected) != len(tt.ExpectedPaths) {
				t.Fatalf("Expected %d paths, got %d\nGot: %#v", len(tt.ExpectedPaths), len(collected), collected)
			}

			for i, expected := range tt.ExpectedPaths {
				if collected[i] != expected {
					t.Errorf("Expected path %q, got %q", expected, collected[i])
				}
			}
		})
	}

}

func TestGenerateFormHTML(t *testing.T) {
	tests := map[string]struct {
		Input       map[string]any
		Contains    []string
		NotContains []string
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
				`name="children.1.firstname"`,
				`name="children.2.firstname"`,
				`name="children.3.firstname"`,
				`name="children.4.firstname"`,
			},
			NotContains: []string{
				`name="children.5.firstname"`,
				`name="children.6.firstname"`,
			},
		},
		"array of primitives": {
			Input: map[string]any{
				"scores": []any{"number"},
			},
			Contains: []string{
				`name="scores.0"`,
				`name="scores.1"`,
				`name="scores.2"`,
				`name="scores.3"`,
				`name="scores.4"`,
			},
			NotContains: []string{
				`name="scores.5"`,
				`name="scores.6"`,
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

			for _, frag := range tt.NotContains {
				if strings.Contains(html, frag) {
					t.Errorf("Expected HTML to NOT contain:\n%s\nGot:\n%s", frag, html)
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
		Input             map[string]any
		FormName          string
		ExpectedMethod    string
		ExpectedAction    string
		ExpectedFields    map[string]string // fieldName -> elementType (input/select/etc)
		NotExpectedFields []string          // champs qui ne doivent pas exister
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
			NotExpectedFields: []string{
				"lastname", "phone",
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
				"children.2.name": "input",
				"children.2.age":  "input",
				"children.3.name": "input",
				"children.3.age":  "input",
				"children.4.name": "input",
				"children.4.age":  "input",
			},
			NotExpectedFields: []string{
				"children.5.name", "children.5.age", "children.6.name",
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
				"company.departments.2.name":              "input",
				"company.departments.2.manager.firstname": "input",
				"company.departments.2.manager.gender":    "select",
				"company.departments.3.name":              "input",
				"company.departments.3.manager.firstname": "input",
				"company.departments.3.manager.gender":    "select",
				"company.departments.4.name":              "input",
				"company.departments.4.manager.firstname": "input",
				"company.departments.4.manager.gender":    "select",
			},
			NotExpectedFields: []string{
				"company.departments.5.name",
				"company.departments.5.manager.firstname",
				"company.departments.5.manager.gender",
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

			for _, fieldName := range tt.NotExpectedFields {
				if _, ok := form.Attribute[fieldName]; ok {
					t.Errorf("Field '%s' should NOT be present in form", fieldName)
				}
			}
		})
	}
}
