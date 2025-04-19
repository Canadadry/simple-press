package data

import (
	"reflect"
	"testing"
)

func TestParseDefinition(t *testing.T) {
	tests := map[string]struct {
		Input    any
		Expected Field
	}{
		"first example": {
			Input: map[string]any{
				"firstname": "string",
				"lastname":  "string",
				"email":     "email",
				"children": []any{
					map[string]any{
						"firstname": "string",
						"gender":    "enum:Mr;Mme",
					},
				},
			},
			Expected: Field{
				IsRoot: true,
				Type:   "object",
				Path:   "",
				Children: []Field{
					{
						Key:    "children",
						Path:   "children",
						Type:   "array",
						Repeat: 5,
						Children: []Field{
							{
								Key:  "",
								Path: "children.0",
								Type: "object",
								Children: []Field{
									{Key: "firstname", Path: "children.0.firstname", Type: "string"},
									{Key: "gender", Path: "children.0.gender", Type: "enum", EnumVals: []string{"Mr", "Mme"}},
								},
							},
						},
					},
					{Key: "email", Path: "email", Type: "email"},
					{Key: "firstname", Path: "firstname", Type: "string"},
					{Key: "lastname", Path: "lastname", Type: "string"},
				},
			},
		},
		"array at root (of objects)": {
			Input: []any{
				map[string]any{
					"title": "string",
					"views": "number",
				},
			},
			Expected: Field{
				IsRoot: true,
				Type:   "array",
				Path:   "",
				Repeat: 5,
				Children: []Field{
					{
						Type: "object",
						Path: "0",
						Children: []Field{
							{Key: "title", Path: "0.title", Type: "string"},
							{Key: "views", Path: "0.views", Type: "number"},
						},
					},
				},
			},
		},
		"array of primitives (numbers)": {
			Input: map[string]any{
				"scores": []any{"number"},
			},
			Expected: Field{
				IsRoot: true,
				Type:   "object",
				Path:   "",
				Children: []Field{
					{
						Key:    "scores",
						Path:   "scores",
						Type:   "array",
						Repeat: 5,
						Children: []Field{
							{Path: "scores.0", Type: "number"},
						},
					},
				},
			},
		},
		"nested object": {
			Input: map[string]any{
				"profile": map[string]any{
					"name": map[string]any{
						"first": "string",
						"last":  "string",
					},
					"age": "number",
				},
			},
			Expected: Field{
				IsRoot: true,
				Type:   "object",
				Path:   "",
				Children: []Field{
					{
						Key:  "profile",
						Path: "profile",
						Type: "object",
						Children: []Field{
							{Key: "age", Path: "profile.age", Type: "number"},
							{
								Key:  "name",
								Path: "profile.name",
								Type: "object",
								Children: []Field{
									{Key: "first", Path: "profile.name.first", Type: "string"},
									{Key: "last", Path: "profile.name.last", Type: "string"},
								},
							},
						},
					},
				},
			},
		},
		"enum only": {
			Input: map[string]any{
				"civility": "enum:Mr;Mme;Dr",
			},
			Expected: Field{
				IsRoot: true,
				Type:   "object",
				Path:   "",
				Children: []Field{
					{Key: "civility", Path: "civility", Type: "enum", EnumVals: []string{"Mr", "Mme", "Dr"}},
				},
			},
		},
		"booleans and dates": {
			Input: map[string]any{
				"active": "bool",
				"dob":    "date",
			},
			Expected: Field{
				IsRoot: true,
				Type:   "object",
				Path:   "",
				Children: []Field{
					{Key: "active", Path: "active", Type: "bool"},
					{Key: "dob", Path: "dob", Type: "date"},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := Parse(tt.Input, true)
			if !reflect.DeepEqual(result, tt.Expected) {
				t.Errorf("Expected:\n%#v\nGot:\n%#v", tt.Expected, result)
			}
		})
	}
}
