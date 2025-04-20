package data

// func TestParseFormData(t *testing.T) {
// 	tests := map[string]struct {
// 		FormValues url.Values
// 		Definition map[string]any
// 		Expected   map[string]any
// 	}{
// 		"flat object": {
// 			FormValues: url.Values{
// 				"firstname":  {"John"},
// 				"email":      {"john@example.com"},
// 				"newsletter": {"true"},
// 			},
// 			Definition: map[string]any{
// 				"firstname":  "string",
// 				"email":      "email",
// 				"newsletter": "bool",
// 			},
// 			Expected: map[string]any{
// 				"firstname":  "John",
// 				"email":      "john@example.com",
// 				"newsletter": true,
// 			},
// 		},
// 		"array of primitives": {
// 			FormValues: url.Values{
// 				"scores.0": {"10"},
// 				"scores.1": {"20"},
// 				"scores.2": {"30"},
// 			},
// 			Definition: map[string]any{
// 				"scores": []any{"number"},
// 			},
// 			Expected: map[string]any{
// 				"scores": []any{10.0, 20.0, 30.0},
// 			},
// 		},
// 		"nested object": {
// 			FormValues: url.Values{
// 				"profile.name.first": {"Alice"},
// 				"profile.name.last":  {"Smith"},
// 				"profile.age":        {"28"},
// 			},
// 			Definition: map[string]any{
// 				"profile": map[string]any{
// 					"name": map[string]any{
// 						"first": "string",
// 						"last":  "string",
// 					},
// 					"age": "number",
// 				},
// 			},
// 			Expected: map[string]any{
// 				"profile": map[string]any{
// 					"name": map[string]any{
// 						"first": "Alice",
// 						"last":  "Smith",
// 					},
// 					"age": 28.0,
// 				},
// 			},
// 		},
// 		"array of objects": {
// 			FormValues: url.Values{
// 				"children.0.firstname": {"Léo"},
// 				"children.0.gender":    {"Mr"},
// 				"children.1.firstname": {"Zoé"},
// 				"children.1.gender":    {"Mme"},
// 			},
// 			Definition: map[string]any{
// 				"children": []any{
// 					map[string]any{
// 						"firstname": "string",
// 						"gender":    "enum:Mr;Mme",
// 					},
// 				},
// 			},
// 			Expected: map[string]any{
// 				"children": []any{
// 					map[string]any{
// 						"firstname": "Léo",
// 						"gender":    "Mr",
// 					},
// 					map[string]any{
// 						"firstname": "Zoé",
// 						"gender":    "Mme",
// 					},
// 				},
// 			},
// 		},
// 		"array with dynamic size": {
// 			FormValues: url.Values{
// 				"tags.0": {"go"},
// 				"tags.1": {"html"},
// 				"tags.2": {"css"},
// 				"tags.3": {"js"},
// 			},
// 			Definition: map[string]any{
// 				"tags": []any{"string"},
// 			},
// 			Expected: map[string]any{
// 				"tags": []any{"go", "html", "css", "js"},
// 			},
// 		},
// 	}

// 	for name, tt := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			body := tt.FormValues.Encode()
// 			req, err := http.NewRequest("POST", "/submit", strings.NewReader(body))
// 			if err != nil {
// 				t.Fatalf("failed to build request: %v", err)
// 			}
// 			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 			root := Parse(tt.Definition, true)

// 			result, err := ParseFormData(req, root)
// 			if err != nil {
// 				t.Fatalf("ParseFormData returned error: %v", err)
// 			}

// 			if !reflect.DeepEqual(result, tt.Expected) {
// 				t.Errorf("Expected:\n%#v\nGot:\n%#v", tt.Expected, result)
// 			}
// 		})
// 	}
// }
