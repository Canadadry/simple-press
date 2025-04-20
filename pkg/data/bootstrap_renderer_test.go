package data

import (
	"io"
	"strings"
	"testing"
)

func TestDynamicRenderer_WithWriter(t *testing.T) {
	tests := map[string]struct {
		Input    any
		Renderer func(w io.Writer) DynamicFormRenderer
		Expected []string
	}{
		"nested object": {
			Input: map[string]any{
				"profile": map[string]any{
					"name": map[string]any{
						"first": "Jane",
						"last":  "Doe",
					},
					"age": 42,
				},
			},
			Renderer: func(w io.Writer) DynamicFormRenderer {
				return NewBootstrapRenderer(w, ThemeBootstrap)
			},
			Expected: []string{
				`<form method="POST" action="/submit" name="form" class="form-bootstrap">`,
				`  <fieldset class="mb-4">`,
				`    <legend class="fw-bold">profile</legend>`,
				`    <div class="mb-3">`,
				`      <label class="form-label">age: <input type="number" name="profile.age" value="42" class="form-control"/></label>`,
				`    </div>`,
				`    <fieldset class="mb-4">`,
				`      <legend class="fw-bold">name</legend>`,
				`      <div class="mb-3">`,
				`        <label class="form-label">first: <input type="text" name="profile.name.first" value="Jane" class="form-control"/></label>`,
				`      </div>`,
				`      <div class="mb-3">`,
				`        <label class="form-label">last: <input type="text" name="profile.name.last" value="Doe" class="form-control"/></label>`,
				`      </div>`,
				`    </fieldset>`,
				`  </fieldset>`,
				`  <button class="btn btn-primary" type="submit">Submit</button>`,
				`</form>`,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var buf strings.Builder
			Render(tt.Input, tt.Renderer(&buf))
			lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
			compareCallStacks(t, lines, tt.Expected)
		})
	}
}

// func TestScrapGeneratedForms(t *testing.T) {
// 	tests := map[string]struct {
// 		Input             map[string]any
// 		FormName          string
// 		ExpectedMethod    string
// 		ExpectedAction    string
// 		ExpectedFields    map[string]string // fieldName -> elementType (input/select/etc)
// 		NotExpectedFields []string          // champs qui ne doivent pas exister
// 	}{
// 		"simple contact form": {
// 			Input: map[string]any{
// 				"email":     "email",
// 				"firstname": "string",
// 				"gender":    "enum:Mr;Mme",
// 			},
// 			FormName:       "contact",
// 			ExpectedMethod: "POST",
// 			ExpectedAction: "/submit",
// 			ExpectedFields: map[string]string{
// 				"email":     "input",
// 				"firstname": "input",
// 				"gender":    "select",
// 			},
// 			NotExpectedFields: []string{
// 				"lastname", "phone",
// 			},
// 		},
// 		"array with nested fields": {
// 			Input: map[string]any{
// 				"children": []any{
// 					map[string]any{
// 						"name": "string",
// 						"age":  "number",
// 					},
// 				},
// 			},
// 			FormName:       "family",
// 			ExpectedMethod: "POST",
// 			ExpectedAction: "/submit",
// 			ExpectedFields: map[string]string{
// 				"children.0.name": "input",
// 				"children.0.age":  "input",
// 				"children.1.name": "input",
// 				"children.1.age":  "input",
// 				"children.2.name": "input",
// 				"children.2.age":  "input",
// 				"children.3.name": "input",
// 				"children.3.age":  "input",
// 				"children.4.name": "input",
// 				"children.4.age":  "input",
// 			},
// 			NotExpectedFields: []string{
// 				"children.5.name", "children.5.age", "children.6.name",
// 			},
// 		},
// 		"deeply nested and mixed structure": {
// 			Input: map[string]any{
// 				"company": map[string]any{
// 					"name": "string",
// 					"address": map[string]any{
// 						"street": "string",
// 						"city":   "string",
// 						"coords": map[string]any{
// 							"lat":  "number",
// 							"long": "number",
// 						},
// 					},
// 					"departments": []any{
// 						map[string]any{
// 							"name": "string",
// 							"manager": map[string]any{
// 								"firstname": "string",
// 								"gender":    "enum:Mr;Mme",
// 							},
// 						},
// 					},
// 				},
// 			},
// 			FormName:       "company-form",
// 			ExpectedMethod: "POST",
// 			ExpectedAction: "/submit",
// 			ExpectedFields: map[string]string{
// 				"company.name":                "input",
// 				"company.address.street":      "input",
// 				"company.address.city":        "input",
// 				"company.address.coords.lat":  "input",
// 				"company.address.coords.long": "input",

// 				"company.departments.0.name":              "input",
// 				"company.departments.0.manager.firstname": "input",
// 				"company.departments.0.manager.gender":    "select",
// 				"company.departments.1.name":              "input",
// 				"company.departments.1.manager.firstname": "input",
// 				"company.departments.1.manager.gender":    "select",
// 				"company.departments.2.name":              "input",
// 				"company.departments.2.manager.firstname": "input",
// 				"company.departments.2.manager.gender":    "select",
// 				"company.departments.3.name":              "input",
// 				"company.departments.3.manager.firstname": "input",
// 				"company.departments.3.manager.gender":    "select",
// 				"company.departments.4.name":              "input",
// 				"company.departments.4.manager.firstname": "input",
// 				"company.departments.4.manager.gender":    "select",
// 			},
// 			NotExpectedFields: []string{
// 				"company.departments.5.name",
// 				"company.departments.5.manager.firstname",
// 				"company.departments.5.manager.gender",
// 			},
// 		},
// 	}

// 	for name, tt := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			field := Parse(tt.Input, true)

// 			var buf strings.Builder
// 			mock := &MockWriterRenderer{w: &buf}
// 			GenerateFormDynamicHTMLWithName(field, mock, tt.FormName)
// 			doc, err := scrapper.NewDocumentFromReader(strings.NewReader(buf.String()))
// 			if err != nil {
// 				t.Fatalf("error parsing HTML: %v", err)
// 			}

// 			form, err := scrapper.GetForm(doc, tt.FormName)
// 			if err != nil {
// 				t.Fatalf("GetForm error: %v", err)
// 			}

// 			if form.Method != tt.ExpectedMethod {
// 				t.Errorf("Expected method %s, got %s", tt.ExpectedMethod, form.Method)
// 			}
// 			if form.Action != tt.ExpectedAction {
// 				t.Errorf("Expected action %s, got %s", tt.ExpectedAction, form.Action)
// 			}

// 			for fieldName, expectedType := range tt.ExpectedFields {
// 				typ, ok := form.Attribute[fieldName]
// 				if !ok {
// 					t.Errorf("Missing field '%s' in form", fieldName)
// 					continue
// 				}
// 				if typ != expectedType {
// 					t.Errorf("Expected type '%s' for field '%s', got '%s'", expectedType, fieldName, typ)
// 				}
// 			}

// 			for _, fieldName := range tt.NotExpectedFields {
// 				if _, ok := form.Attribute[fieldName]; ok {
// 					t.Errorf("Field '%s' should NOT be present in form", fieldName)
// 				}
// 			}
// 		})
// 	}
// }
