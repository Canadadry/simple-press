package data

import (
	"errors"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

func ParseFormData(r *http.Request, schema Field) (map[string]any, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	return extractFormData(r.Form, schema)
}

func extractFormData(values url.Values, f Field) (map[string]any, error) {
	if f.Type != "object" {
		return nil, errors.New("schema root must be an object")
	}

	result := map[string]any{}

	for _, child := range f.Children {
		val, err := extractValue(values, child)
		if err != nil {
			return nil, err
		}
		result[child.Key] = val
	}

	return result, nil
}

func extractValue(values url.Values, f Field) (any, error) {
	switch f.Type {
	case "string", "email", "enum", "date":
		return values.Get(f.Path), nil

	case "number":
		raw := values.Get(f.Path)
		if raw == "" {
			return 0.0, nil
		}
		return strconv.ParseFloat(raw, 64)

	case "bool":
		return values.Get(f.Path) == "true", nil

	case "object":
		obj := map[string]any{}
		for _, child := range f.Children {
			val, err := extractValue(values, child)
			if err != nil {
				return nil, err
			}
			obj[child.Key] = val
		}
		return obj, nil

	case "array":
		var arr []any
		indices := extractArrayIndices(values, f.Path)

		for _, i := range indices {
			itemSchema := f.Children[0]
			clone := updatePathForArrayIndex(itemSchema, i)

			val, err := extractValue(values, clone)
			if err != nil {
				return nil, err
			}

			if isEmptyValue(val) {
				continue
			}
			arr = append(arr, val)
		}

		return arr, nil

	default:
		return nil, errors.New("unsupported type: " + f.Type)
	}
}

func isEmptyValue(v any) bool {
	switch vv := v.(type) {
	case string:
		return vv == ""
	case map[string]any:
		return len(vv) == 0
	default:
		return false
	}
}

func extractArrayIndices(values url.Values, prefix string) []int {
	indexSet := map[int]struct{}{}

	for key := range values {
		if strings.HasPrefix(key, prefix+".") {
			parts := strings.Split(key[len(prefix)+1:], ".") // after "prefix."
			if len(parts) == 0 {
				continue
			}
			if idx, err := strconv.Atoi(parts[0]); err == nil {
				indexSet[idx] = struct{}{}
			}
		}
	}

	// Convert map to sorted slice
	var indices []int
	for i := range indexSet {
		indices = append(indices, i)
	}
	sort.Ints(indices)

	return indices
}
