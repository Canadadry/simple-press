package data

import (
	"strings"
	"testing"
)

func TestGenerateFormDynamicHTMLWithName(t *testing.T) {
	tests := map[string]struct {
		Input    map[string]any
		FormName string
		Contains []string
	}{
		"simple dynamic array": {
			Input: map[string]any{
				"children": []any{
					map[string]any{
						"firstname": "string",
						"gender":    "enum:Mr;Mme",
					},
				},
			},
			FormName: "family",
			Contains: []string{
				`<form`,
				`name="children.0.firstname"`,
				`name="children.0.gender"`,
				`<select name="children.0.gender"`,
				`<template id="template-children">`,
				`name="children.__INDEX__.firstname"`,
				`name="children.__INDEX__.gender"`,
			},
		},
		"nested dynamic object with array": {
			Input: map[string]any{
				"profile": map[string]any{
					"firstname": "string",
					"emails":    []any{"email"},
				},
			},
			FormName: "nested",
			Contains: []string{
				`name="profile.firstname"`,
				`name="profile.emails.0"`,
				`<template id="template-profile.emails">`,
				`name="profile.emails.__INDEX__"`,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			root := Parse(tt.Input, true)
			html := GenerateFormDynamicHTMLWithName(root, ThemeNoStyle, tt.FormName)

			for _, frag := range tt.Contains {
				if !strings.Contains(html, frag) {
					t.Errorf("Expected HTML to contain:\n%s\nGot:\n%s", frag, html)
				}
			}
		})
	}
}

func TestGenerateDynamicJS(t *testing.T) {
	tests := map[string]struct {
		Input      map[string]any
		FormName   string
		ContainsJS []string
	}{
		"simple array js": {
			Input: map[string]any{
				"children": []any{
					map[string]any{
						"firstname": "string",
					},
				},
			},
			FormName: "family",
			ContainsJS: []string{
				"let currentIndex_children = 1",
				"document.getElementById(\"template-children\")",
				"document.getElementById(\"container-children\")",
				".replaceAll(\"__INDEX__\", currentIndex_children)",
				"currentIndex_children++",
			},
		},
		"nested array js": {
			Input: map[string]any{
				"profile": map[string]any{
					"emails": []any{"string"},
				},
			},
			FormName: "nested",
			ContainsJS: []string{
				"let currentIndex_profile_emails = 1",
				"document.getElementById(\"template-profile.emails\")",
				"document.getElementById(\"container-profile.emails\")",
				".replaceAll(\"__INDEX__\", currentIndex_profile_emails)",
				"currentIndex_profile_emails++",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			field := Parse(tt.Input, true)
			js := GenerateDynamicJS(field)

			for _, frag := range tt.ContainsJS {
				if !strings.Contains(js, frag) {
					t.Errorf("Expected JS to contain:\n%s\nGot:\n%s", frag, js)
				}
			}
		})
	}
}

func TestGenerateDynamicFormWithJS_ExactMatch(t *testing.T) {
	tests := map[string]struct {
		Input        map[string]any
		FormName     string
		ExpectedHTML string
		ExpectedJS   string
	}{
		"simple children array": {
			Input: map[string]any{
				"children": []any{
					map[string]any{
						"name": "string",
					},
				},
			},
			FormName: "family",
			ExpectedHTML: `
<form method="POST" action="/submit" name="family" class="">
  <fieldset class="">
    <legend class="">children</legend>
    <div id="container-children">
      <div>
        <div class="">
          <label class="">name: <input type="text" name="children.0.name" class=""/></label><br/>
        </div>
      </div>
    </div>
    <button type="button" class="" onclick="add_children()">Add</button>
    <template id="template-children">
      <div>
        <div class="">
          <label class="">name: <input type="text" name="children.__INDEX__.name" class=""/></label><br/>
        </div>
      </div>
    </template>
  </fieldset>
  <button class="" type="submit">Submit</button>
</form>
`,

			ExpectedJS: `
let currentIndex_children = 1;

function add_children() {
  const template = document.getElementById("template-children").innerHTML;
  const container = document.getElementById("container-children");
  const html = template.replaceAll("__INDEX__", currentIndex_children);
  const temp = document.createElement("div");
  temp.innerHTML = html;
  container.appendChild(temp.firstElementChild);
  currentIndex_children++;
}
`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			root := Parse(tt.Input, true)

			actualHTML := GenerateFormDynamicHTMLWithName(root, ThemeNoStyle, tt.FormName)
			actualJS := GenerateDynamicJS(root)

			// Normalisation HTML et JS
			normalize := func(s string) string {
				s = strings.ReplaceAll(s, "\n", "")
				s = strings.ReplaceAll(s, "\t", "")
				s = strings.ReplaceAll(s, "  ", "")
				return strings.TrimSpace(s)
			}

			expectedHTML := normalize(tt.ExpectedHTML)
			expectedJS := normalize(tt.ExpectedJS)
			actualHTML = normalize(actualHTML)
			actualJS = normalize(actualJS)

			if actualHTML != expectedHTML {
				t.Errorf("HTML output mismatch.\nExpected:\n%s\n\nGot:\n%s", tt.ExpectedHTML, actualHTML)
			}

			if actualJS != expectedJS {
				t.Errorf("JS output mismatch.\nExpected:\n%s\n\nGot:\n%s", tt.ExpectedJS, actualJS)
			}
		})
	}
}
