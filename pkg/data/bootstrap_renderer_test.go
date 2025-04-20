package data

import (
	"app/pkg/scrapper"
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

func TestScrapGeneratedForms(t *testing.T) {
	tests := map[string]struct {
		Input          map[string]any
		FormName       string
		ExpectedMethod string
		ExpectedAction string
		ExpectedFields map[string]string
	}{

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
					"departments": map[string]any{
						"name": "string",
						"manager": map[string]any{
							"firstname": "string",
							"gender":    "enum:Mr;Mme",
						},
					},
				},
			},
			FormName:       "company-form",
			ExpectedMethod: "POST",
			ExpectedAction: "/submit",
			ExpectedFields: map[string]string{
				"company.name":                          "input",
				"company.address.street":                "input",
				"company.address.city":                  "input",
				"company.address.coords.lat":            "input",
				"company.address.coords.long":           "input",
				"company.departments.name":              "input",
				"company.departments.manager.firstname": "input",
				"company.departments.manager.gender":    "input",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			theme := ThemeBootstrap
			theme.FormName = tt.FormName
			var buf strings.Builder
			err := Render(tt.Input, NewBootstrapRenderer(&buf, theme))
			if err != nil {
				t.Fatalf("error rendering HTML: %v", err)
			}

			doc, err := scrapper.NewDocumentFromReader(strings.NewReader(buf.String()))
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
