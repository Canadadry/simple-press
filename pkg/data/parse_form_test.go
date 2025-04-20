package data

import (
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestParseFormData(t *testing.T) {
	tests := map[string]struct {
		FormValues url.Values
		Definition map[string]any
		Expected   map[string]any
	}{
		"flat object": {
			FormValues: url.Values{
				"firstname":  {"John"},
				"email":      {"john@example.com"},
				"newsletter": {"true"},
			},
			Definition: map[string]any{
				"firstname":  "Alice",
				"email":      "Alice@example.com",
				"newsletter": false,
			},
			Expected: map[string]any{
				"firstname":  "John",
				"email":      "john@example.com",
				"newsletter": true,
			},
		},
		"nested object": {
			FormValues: url.Values{
				"profile.name.first": {"Alice"},
				"profile.name.last":  {"Smith"},
				"profile.age":        {"28"},
			},
			Definition: map[string]any{
				"profile": map[string]any{
					"name": map[string]any{
						"first": "Bob",
						"last":  "Morane",
					},
					"age": 10,
				},
			},
			Expected: map[string]any{
				"profile": map[string]any{
					"name": map[string]any{
						"first": "Alice",
						"last":  "Smith",
					},
					"age": 28.0,
				},
			},
		},
		"missing field default to old value": {
			FormValues: url.Values{
				"profile.name.first": {"Alice"},
				"profile.name.last":  {"Smith"},
				"profile.age":        {"28"},
				"profile.gender":     {"Mme"},
			},
			Definition: map[string]any{
				"profile": map[string]any{
					"name": map[string]any{
						"first": "Bob",
						"last":  "Morane",
					},
					"age": 10,
				},
			},
			Expected: map[string]any{
				"profile": map[string]any{
					"name": map[string]any{
						"first": "Alice",
						"last":  "Smith",
					},
					"age": 28.0,
				},
			},
		},

		"extra field ignored": {
			FormValues: url.Values{
				"profile.name.first": {"Alice"},
				"profile.name.last":  {"Smith"},
			},
			Definition: map[string]any{
				"profile": map[string]any{
					"name": map[string]any{
						"first": "Bob",
						"last":  "Morane",
					},
					"age": 10,
				},
			},
			Expected: map[string]any{
				"profile": map[string]any{
					"name": map[string]any{
						"first": "Alice",
						"last":  "Smith",
					},
					"age": 10,
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			body := tt.FormValues.Encode()
			req, err := http.NewRequest("POST", "/submit", strings.NewReader(body))
			if err != nil {
				t.Fatalf("failed to build request: %v", err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			result, err := ParseFormData(req, tt.Definition)
			if err != nil {
				t.Fatalf("ParseFormData returned error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.Expected) {
				t.Errorf("Expected:\n%#v\nGot:\n%#v", tt.Expected, result)
			}
		})
	}
}
