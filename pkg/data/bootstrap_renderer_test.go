package data

import (
	"io"
	"strings"
	"testing"
)

func TestDynamicRenderer_WithWriter(t *testing.T) {
	tests := map[string]struct {
		Render   func(DynamicFormRenderer)
		Renderer func(w io.Writer) DynamicFormRenderer
		Expected []string
	}{
		"basic bootstrap form": {
			Render: func(r DynamicFormRenderer) {
				r.BeginForm("demo", "/submit", "POST")
				r.Input("firstname", "firstname", "text")
				r.Select("gender", "gender", []string{"Mr", "Mme"})
				r.Checkbox("newsletter", "newsletter")
				r.EndForm()
			},
			Renderer: func(w io.Writer) DynamicFormRenderer {
				return NewBootstrapRenderer(w, ThemeBootstrap)
			},
			Expected: []string{
				`<form method="POST" action="/submit" name="demo" class="form-bootstrap">`,
				`  <div class="mb-3">`,
				`    <label class="form-label">firstname: <input type="text" name="firstname" class="form-control"/></label>`,
				`  </div>`,
				`  <div class="mb-3">`,
				`    <label class="form-label">gender: <select name="gender" class="form-select">`,
				`      <option value="Mr">Mr</option>`,
				`      <option value="Mme">Mme</option>`,
				`    </select></label>`,
				`  </div>`,
				`  <div class="mb-3">`,
				`    <label class="form-label">newsletter: <input type="checkbox" name="newsletter" value="true" class="form-check-input"/></label>`,
				`  </div>`,
				`  <button class="btn btn-primary" type="submit">Submit</button>`,
				`</form>`,
			},
		},
		"array of emails": {
			Render: func(r DynamicFormRenderer) {
				r.BeginForm("emailForm", "/submit", "POST")
				r.BeginFieldset("emails")
				r.BeginArray("emails", "emails")
				r.BeginArrayItem(0)
				r.Input("", "emails.0", "email")
				r.EndArrayItem()
				r.EndArray()
				r.EndFieldset()
				r.EndForm()
			},
			Renderer: func(w io.Writer) DynamicFormRenderer {
				return NewBootstrapRenderer(w, ThemeBootstrap)
			},
			Expected: []string{
				`<form method="POST" action="/submit" name="emailForm" class="form-bootstrap">`,
				`  <fieldset class="mb-4">`,
				`    <legend class="fw-bold">emails</legend>`,
				`    <div id="container-emails">`,
				`      <div data-item class="row mb-3">`,
				`        <div class="col">`,
				`          <div class="mb-3">`,
				`            <label class="form-label">: <input type="email" name="emails.0" class="form-control"/></label>`,
				`          </div>`,
				`        </div>`,
				`      </div>`,
				`    </div>`,
				`  </fieldset>`,
				`  <button class="btn btn-primary" type="submit">Submit</button>`,
				`</form>`,
			},
		}, "nested array (children -> pets)": {
			Render: func(r DynamicFormRenderer) {
				r.BeginForm("nestedSlice", "/submit", "POST")
				r.BeginFieldset("children")
				r.BeginArray("children", "children")
				r.BeginArrayItem(0)
				r.Input("name", "children.0.name", "text")
				r.BeginFieldset("pets")
				r.BeginArray("pets", "children.0.pets")
				r.BeginArrayItem(0)
				r.Input("pet", "children.0.pets.0", "text")
				r.EndArrayItem()
				r.EndArray()
				r.EndFieldset()
				r.EndArrayItem()
				r.EndArray()
				r.EndFieldset()
				r.EndForm()
			},
			Renderer: func(w io.Writer) DynamicFormRenderer {
				return NewBootstrapRenderer(w, ThemeBootstrap)
			},
			Expected: []string{
				`<form method="POST" action="/submit" name="nestedSlice" class="form-bootstrap">`,
				`  <fieldset class="mb-4">`,
				`    <legend class="fw-bold">children</legend>`,
				`    <div id="container-children">`,
				`      <div data-item class="row mb-3">`,
				`        <div class="col">`,
				`          <div class="mb-3">`,
				`            <label class="form-label">name: <input type="text" name="children.0.name" class="form-control"/></label>`,
				`          </div>`,
				`          <fieldset class="mb-4">`,
				`            <legend class="fw-bold">pets</legend>`,
				`            <div id="container-children.0.pets">`,
				`              <div data-item class="row mb-3">`,
				`                <div class="col">`,
				`                  <div class="mb-3">`,
				`                    <label class="form-label">pet: <input type="text" name="children.0.pets.0" class="form-control"/></label>`,
				`                  </div>`,
				`                </div>`,
				`              </div>`,
				`            </div>`,
				`          </fieldset>`,
				`        </div>`,
				`      </div>`,
				`    </div>`,
				`  </fieldset>`,
				`  <button class="btn btn-primary" type="submit">Submit</button>`,
				`</form>`,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var buf strings.Builder
			renderer := tt.Renderer(&buf)
			tt.Render(renderer)

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
