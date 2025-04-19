package data

import (
	"strings"
	"testing"
)

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
				`<input type="text" name="firstname" />`,
				`<input type="email" name="email" />`,
				`<input type="checkbox" name="newsletter" value="true" />`,
			},
		},
		"object with enum and date": {
			Input: map[string]any{
				"civility": "enum:Mr;Mme",
				"dob":      "date",
			},
			Contains: []string{
				`<select name="civility">`,
				`<option value="Mr">Mr</option>`,
				`<option value="Mme">Mme</option>`,
				`<input type="date" name="dob" />`,
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
				`<input type="text" name="profile.name.first" />`,
				`<input type="text" name="profile.name.last" />`,
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
				`<input type="text" name="children.0.firstname" />`,
				`<select name="children.0.gender">`,
				`<input type="text" name="children.1.firstname" />`,
				`<select name="children.1.gender">`,
			},
		},
		"array of primitives": {
			Input: map[string]any{
				"scores": []any{"number"},
			},
			Contains: []string{
				`<input type="number" name="scores.0" />`,
				`<input type="number" name="scores.1" />`,
				`<input type="number" name="scores.2" />`,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			root := Parse(tt.Input, true)
			html := GenerateFormHTML(root)

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
		ExpectedHTML string
	}{
		"simple object with enum": {
			Input: map[string]any{
				"firstname": "string",
				"email":     "email",
				"civility":  "enum:Mr;Mme",
			},
			ExpectedHTML: `<form method="POST" action="/submit">` +
				`<label>civility: <select name="civility"><option value="Mr">Mr</option><option value="Mme">Mme</option></select></label><br/>` +
				`<label>email: <input type="email" name="email" /></label><br/>` +
				`<label>firstname: <input type="text" name="firstname" /></label><br/>` +
				`<button type="submit">Submit</button></form>`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			root := Parse(tt.Input, true)
			html := GenerateFormHTML(root)

			// Remove newlines/spaces for stable match (if needed)
			if html != tt.ExpectedHTML {
				t.Errorf("HTML output mismatch.\nExpected:\n%s\n\nGot:\n%s", tt.ExpectedHTML, html)
			}
		})
	}
}
