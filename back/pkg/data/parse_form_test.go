package data

import (
	"fmt"
	"testing"
)

func TestParseFormData(t *testing.T) {
	tests := map[string]struct {
		FormValues map[string]any
		Definition map[string]any
		Expected   map[string]any
	}{
		"flat object": {
			FormValues: map[string]any{
				"firstname":  "John",
				"email":      "john@example.com",
				"newsletter": "true",
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
			FormValues: map[string]any{
				"profile": map[string]any{
					"name": map[string]any{
						"first": "Alice",
						"last":  "Smith",
					},
					"age": 28,
				},
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
			FormValues: map[string]any{
				"profile": map[string]any{
					"name": map[string]any{
						"first": "Alice",
						"last":  "Smith",
					},
					"age":    28,
					"gender": "Mme",
				},
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
			FormValues: map[string]any{
				"profile": map[string]any{
					"name": map[string]any{
						"first": "Alice",
						"last":  "Smith",
					},
				},
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

		"false value": {
			FormValues: map[string]any{
				"boolean1": false,
				"boolean2": "false",
				"boolean3": "0",
				"boolean4": 0,
			},
			Definition: map[string]any{
				"boolean1": true,
				"boolean2": true,
				"boolean3": true,
				"boolean4": true,
			},
			Expected: map[string]any{
				"boolean1": false,
				"boolean2": false,
				"boolean3": false,
				"boolean4": false,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := ParseFormData(tt.FormValues, tt.Definition)
			if err != nil {
				t.Fatalf("ParseFormData returned error: %v", err)
			}

			if fmt.Sprintf("%#v", tt.Expected) != fmt.Sprintf("%#v", result) {
				t.Errorf("Expected:\n%#v\nGot:\n%#v", tt.Expected, result)
			}
		})
	}
}
