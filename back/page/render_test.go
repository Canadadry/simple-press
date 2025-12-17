package page

import (
	"bytes"
	"testing"
)

func TestRender(t *testing.T) {
	tests := map[string]struct {
		input Data
		exp   string
		err   string
	}{
		"test if base template is not defined": {
			input: Data{
				Files: map[string]string{},
			},
			err: "base template baseof.html not defined",
		},
		"test minimal success with a simple layout and content": {
			input: Data{
				Files: map[string]string{
					"baseof.html": "<html><body>{{.Content}}</body></html>",
				},
				Content: "Simple Content",
			},
			exp: "<html><body>Simple Content</body></html>",
			err: "",
		},
		"test with two files for layout and body": {
			input: Data{
				Files: map[string]string{
					"baseof.html": `<html><head></head><body>{{template "body" .}}</body></html>`,
					"main_layout": `{{define "body"}}<h1>{{.Title}}</h1><p>{{.Content}}</p>{{end}}`,
				},
				Content: "Test Content",
				Title:   "Test Title",
			},
			exp: `<html><head></head><body><h1>Test Title</h1><p>Test Content</p></body></html>`,
			err: "",
		},
		"test with two files for layout and body with block": {
			input: Data{
				Files: map[string]string{
					"baseof.html": `<html><head></head><body>{{template "body" .}}</body></html>`,
					"main_layout": `{{define "body"}}<h1>{{.Title}}</h1><p>{{.Content}}</p>{{ range $id, $block := .Blocks }}<p>{{partial $block}}</p>{{end}}{{end}}`,
				},
				BlocksContent: map[string]string{
					"basic": "{{.Data.Content}}",
				},
				ArticleBlocks: []ArticleBlock{
					{
						BlockName: "basic",
						Data: map[string]any{
							"Data": map[string]any{
								"Content": "something",
							},
						},
					},
				},
				Content: "Test Content",
				Title:   "Test Title",
			},
			exp: `<html><head></head><body><h1>Test Title</h1><p>Test Content</p><p>something</p></body></html>`,
			err: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			err := Render(&buf, tt.input)
			if err != nil {
				if tt.err == "" {
					t.Fatalf("rendering failed: %v", err)
				} else if err.Error() != tt.err {
					t.Fatalf("expected error: %v\n got %v\n", tt.err, err)
				}
			}
			got := buf.String()
			if got != tt.exp {
				t.Fatalf("\nwant -%v-\ngot  -%v-", tt.exp, got)
			}
		})
	}
}
